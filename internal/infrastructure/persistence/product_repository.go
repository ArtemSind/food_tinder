package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ArtemSind/food_tinder/internal/infrastructure/external/foodji"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const (
	// ProductCacheKeyPrefix is the prefix for product cache keys in Redis
	ProductCacheKeyPrefix = "product:"
	// ProductListKey is the key for storing the list of all product IDs
	ProductListKey = "products:list"
	// UpdateInterval is how often to update the product data
	UpdateInterval = 24 * time.Hour
	// DefaultMachineID is the default machine ID to use for API requests
	DefaultMachineID = "4bf115ee-303a-4089-a3ea-f6e7aae0ab94"
	// ProductCacheTimeout is the expiration time for cached products
	ProductCacheTimeout = 24 * time.Hour
	// UpdateContextTimeout is the timeout for each update operation
	UpdateContextTimeout = 5 * time.Minute
)

// ProductRepository handles product data persistence and caching
type ProductRepository struct {
	redisClient  *redis.Client
	foodjiClient *foodji.Client
	cancel       context.CancelFunc
}

// NewProductRepository creates a new product repository
func NewProductRepository(redisClient *redis.Client, foodjiClient *foodji.Client) *ProductRepository {
	ctx, cancel := context.WithCancel(context.Background())
	repo := &ProductRepository{
		redisClient:  redisClient,
		foodjiClient: foodjiClient,
		cancel:       cancel,
	}

	// Start the periodic update in a goroutine
	go repo.startPeriodicUpdate(ctx)

	return repo
}

// Close stops the periodic update goroutine
func (r *ProductRepository) Close() error {
	if r.cancel != nil {
		log.Println("Shutting down product repository")
		r.cancel()
	}
	return nil
}

// startPeriodicUpdate starts a goroutine to update product data periodically
func (r *ProductRepository) startPeriodicUpdate(ctx context.Context) {
	ticker := time.NewTicker(UpdateInterval)
	defer ticker.Stop()

	// Initial update
	if err := r.updateProducts(ctx); err != nil {
		log.Printf("Error in initial product update: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			// Create a timeout context for each update operation
			updateCtx, cancel := context.WithTimeout(ctx, UpdateContextTimeout)
			if err := r.updateProducts(updateCtx); err != nil {
				log.Printf("Error in periodic product update: %v", err)
			}
			cancel()
		case <-ctx.Done():
			return
		}
	}
}

// updateProducts fetches products from the Foodji API and updates Redis
func (r *ProductRepository) updateProducts(ctx context.Context) error {
	log.Println("Starting product update from Foodji API")

	// Get products from the Foodji API
	products, err := r.foodjiClient.GetMachineProducts(ctx, DefaultMachineID)
	if err != nil {
		return fmt.Errorf("failed to fetch products from Foodji API: %w", err)
	}

	log.Printf("Retrieved %d products from API", len(*products))

	// Start a Redis transaction
	pipe := r.redisClient.TxPipeline()

	// Clear existing product list
	pipe.Del(ctx, ProductListKey)

	// Store each product
	for _, product := range *products {
		productID := product.ID.String()
		productKey := fmt.Sprintf("%s%s", ProductCacheKeyPrefix, productID)

		// Store product data
		productJSON, err := json.Marshal(product)
		if err != nil {
			return fmt.Errorf("failed to marshal product: %w", err)
		}

		pipe.Set(ctx, productKey, productJSON, ProductCacheTimeout)
		pipe.SAdd(ctx, ProductListKey, productID)
	}

	// Execute the transaction
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update Redis: %w", err)
	}

	log.Println("Product cache successfully updated")
	return nil
}

// GetProduct retrieves a product by ID from Redis
func (r *ProductRepository) GetProduct(ctx context.Context, id uuid.UUID) (*foodji.Product, error) {
	productKey := fmt.Sprintf("%s%s", ProductCacheKeyPrefix, id.String())

	// Create a timeout context for the Redis operation
	redisCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	productJSON, err := r.redisClient.Get(redisCtx, productKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product from Redis: %w", err)
	}

	var product foodji.Product
	if err := json.Unmarshal([]byte(productJSON), &product); err != nil {
		return nil, fmt.Errorf("failed to unmarshal product: %w", err)
	}

	return &product, nil
}

// ListProducts retrieves all product IDs from Redis
func (r *ProductRepository) ListProducts(ctx context.Context) ([]uuid.UUID, error) {
	// Create a timeout context for the Redis operation
	redisCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	productIDs, err := r.redisClient.SMembers(redisCtx, ProductListKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get product list from Redis: %w", err)
	}

	if len(productIDs) == 0 {
		log.Println("No products found in Redis, triggering update")
		// If no products are found, try to update them
		updateCtx, updateCancel := context.WithTimeout(ctx, UpdateContextTimeout)
		defer updateCancel()

		if err := r.updateProducts(updateCtx); err != nil {
			log.Printf("Failed to update products: %v", err)
			return []uuid.UUID{}, nil
		}

		// Try again after update
		productIDs, err = r.redisClient.SMembers(redisCtx, ProductListKey).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get product list from Redis after update: %w", err)
		}
	}

	ids := make([]uuid.UUID, 0, len(productIDs))
	for _, idStr := range productIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			log.Printf("Invalid UUID in product list: %s", idStr)
			continue // Skip invalid UUIDs
		}
		ids = append(ids, id)
	}

	return ids, nil
}
