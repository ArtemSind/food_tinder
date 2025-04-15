package persistence

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

// NewRedisClient creates a new Redis client
func NewRedisClient(ctx context.Context) (*redis.Client, error) {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Try pinging the Redis server with retries
	maxRetries := 5
	retryDelay := 5 * time.Second

	for i := range maxRetries {
		if err := client.Ping(ctx).Err(); err == nil {
			log.Println("Successfully connected to Redis")
			return client, nil
		} else {
			log.Printf("Failed to connect to Redis (attempt %d/%d): %v", i+1, maxRetries, err)
			if i < maxRetries-1 {
				log.Printf("Retrying in %v...", retryDelay)
				time.Sleep(retryDelay)
			} else {
				return nil, fmt.Errorf("failed to connect to Redis after %d attempts", maxRetries)
			}
		}
	}

	// This line should not be reached due to the for loop above
	return client, nil
}
