package redisadaptor

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisSubscriber struct {
	client *redis.Client
}

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
			}
		case <-ctx.Done():
			log.Println("Subscription cancelled")
			return ctx.Err()
		}
	}
}

// Close closes the Redis connection
func (r *RedisSubscriber) Close() error {
	return r.client.Close()
}
