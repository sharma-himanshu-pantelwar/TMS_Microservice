package redisadaptor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	notifications "notification_service/src/internal/core/notification"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisSubscriber struct {
	client *redis.Client
}

// NotificationData represents the structure of a notification to be stored

func NewRedisSubscriber(addr, password string, db int) *RedisSubscriber {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisSubscriber{client: rdb}
}

// SubscribeToChannel subscribes to a Redis channel and prints messages
func (r *RedisSubscriber) SubscribeToChannel(ctx context.Context, channel string) error {
	pubsub := r.client.Subscribe(ctx, channel)
	defer pubsub.Close()

	// Wait for confirmation that subscription is created
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return fmt.Errorf("failed to confirm subscription: %w", err)
	}

	log.Printf(" Successfully subscribed to channel: %s", channel)
	log.Println("Listening for notifications...")

	// Listen for messages and print them
	ch := pubsub.Channel()
	for {
		select {
		case msg := <-ch:
			if msg != nil {
				fmt.Printf("Notification [%s] %s\n", msg.Channel, msg.Payload)

				// Store notification in Redis using SET
				err := r.StoreNotification(ctx, msg.Channel, msg.Payload)
				if err != nil {
					log.Printf("Failed to store notification: %v", err)
				} else {
					log.Printf("Notification stored successfully in Redis")
				}
			}
		case <-ctx.Done():
			log.Println("Subscription end")
			return ctx.Err()
		}
	}
}

// StoreNotification stores a notification in Redis using SET command
func (r *RedisSubscriber) StoreNotification(ctx context.Context, channel, payload string) error {
	// Generate unique key for the notification
	notificationID := fmt.Sprintf("%s:%d", channel, time.Now().UnixNano())

	// Create notification data structure
	notification := notifications.NotificationData{
		ID:        notificationID,
		Channel:   channel,
		Payload:   payload,
		Timestamp: time.Now(),
		Status:    "received",
	}

	// Convert to JSON
	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	// Store in Redis with TTL (24 hours)
	err = r.client.Set(ctx, notificationID, notificationJSON, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to store notification in Redis: %w", err)
	}

	// Also add to a sorted set for easy retrieval by timestamp
	err = r.client.ZAdd(ctx, fmt.Sprintf("notifications:%s", channel), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: notificationID,
	}).Err()
	if err != nil {
		log.Printf("Warning: failed to add to sorted set: %v", err)
	}

	return nil
}

// Close closes the Redis connection
func (r *RedisSubscriber) Close() error {
	return r.client.Close()
}
