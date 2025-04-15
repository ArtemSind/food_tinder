package persistence

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPostgresPool(ctx context.Context) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Retry logic for database connection
	var dbPool *pgxpool.Pool
	var err error
	maxRetries := 5
	retryDelay := 5 * time.Second

	for i := range maxRetries {
		dbPool, err = pgxpool.Connect(ctx, connString)
		if err == nil {
			// Test the connection
			if err = dbPool.Ping(ctx); err == nil {
				log.Println("Successfully connected to database")
				return dbPool, nil
			}
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			log.Printf("Retrying in %v...", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}
