package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Create a background context for the Redis operations
	ctx := context.Background()

	// Initialize the Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set by default
		DB:       0,                // Use default DB 0
	})

	// Safely close the client connection pool when main() finishes
	defer rdb.Close()

	// Test the connection using Ping
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Printf("Connected successfully! Redis responded: %s\n", pong)

	// Example: Set a key with a 10-minute expiration
	err = rdb.Set(ctx, "user:100:name", "Alice", 10*time.Minute).Err()
	if err != nil {
		log.Fatalf("Failed to set key: %v", err)
	}

	// Example: Get the key back
	val, err := rdb.Get(ctx, "user:100:name").Result()
	if err != nil {
		log.Fatalf("Failed to get key: %v", err)
	}
	fmt.Printf("Retrieved value for 'user:100:name': %s\n", val)
}
