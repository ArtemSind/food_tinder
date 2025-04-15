package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArtemSind/food_tinder/internal/application"
	"github.com/ArtemSind/food_tinder/internal/infrastructure/external/foodji"
	httpRouter "github.com/ArtemSind/food_tinder/internal/infrastructure/http"
	"github.com/ArtemSind/food_tinder/internal/infrastructure/persistence"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Create context that listens for signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalCh
		log.Println("Shutdown signal received")
		cancel()
	}()

	// Init closer slice for gracefull shutdown
	var closers []func() error

	// Connect to database
	db, err := persistence.NewPostgresPool(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	closers = append(closers, func() error {
		db.Close()
		return nil
	})

	// Ping the database to verify connection
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database")

	// Connect to Redis
	redisClient, err := persistence.NewRedisClient(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	closers = append(closers, redisClient.Close)

	// Create http clients
	foodjiClient := foodji.NewClient()

	// Create repositories
	sessionRepo := persistence.NewSessionRepository(db)
	voteRepo := persistence.NewVoteRepository(db)
	productRepo := persistence.NewProductRepository(redisClient, foodjiClient)

	closers = append(closers, func() error {
		productRepo.Close()
		return nil
	})

	// Create services
	sessionService := application.NewSessionService(sessionRepo)
	voteService := application.NewVoteService(voteRepo, productRepo)

	// Create router
	router := httpRouter.NewRouter(sessionService, voteService)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server in a goroutine so it doesn't block shutdown
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for context cancellation (from signal)
	<-ctx.Done()
	log.Println("Server shutting down")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	// Gracefully shutdown the server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	// Close repositories
	var shutdownErrs []error
	for i := len(closers) - 1; i >= 0; i-- {
		if err := closers[i](); err != nil {
			log.Printf("Error closing resource %d: %v", i, err)
			shutdownErrs = append(shutdownErrs, err)
		}
	}

	if len(shutdownErrs) > 0 {
		log.Printf("Encountered %d errors during shutdown", len(shutdownErrs))
	} else {
		log.Println("All resources closed successfully")
	}

	log.Println("Server stopped gracefully")
}
