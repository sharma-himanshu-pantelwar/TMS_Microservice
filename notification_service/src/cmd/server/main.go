package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	redisadaptor "notification_service/src/internal/adaptors/redis"
)

func main() {
	fmt.Println("Notification Service...")

	// Create Redis subscriber
	redisSubscriber := redisadaptor.NewRedisSubscriber("localhost:6379", "", 0)
	defer redisSubscriber.Close()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nReceived shutdown signal, gracefully shutting down...")
		cancel()
	}()

	// Start subscribing to Redis channel
	fmt.Println(" Listening for task notifications...")

	err := redisSubscriber.SubscribeToChannel(ctx, "task_notifications")
	if err != nil && err != context.Canceled {
		log.Fatalf("Redis subscription error: %v", err)
	}

	fmt.Println("Notification Service shutdown complete")
}
