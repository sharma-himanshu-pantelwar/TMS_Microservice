package redisadaptor

//publish notification here whenever the task is created or updated
import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisPublisher struct {
	client *redis.Client
}

func NewRedisPublisher(addr, password string, db int) *RedisPublisher {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisPublisher{client: rdb}
}

// PublishTaskNotification publishes a notification to a Redis channel when a task is created or updated
func (r *RedisPublisher) PublishTaskNotification(ctx context.Context, channel string, message string) error {
	err := r.client.Publish(ctx, channel, message).Err()
	if err != nil {
		return fmt.Errorf("failed to publish notification: %w", err)
	}
	return nil
}
