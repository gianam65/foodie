package product

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"foodie/backend/internal/domain/product"
	"foodie/backend/internal/infrastructure/cache"
)

// useCaseImpl implements the UseCase interface.
type useCaseImpl struct {
	productRepo product.Repository
	cache       cache.Cache // Optional cache, can be nil
}

// NewUseCase creates a new product use case.
func NewUseCase(productRepo product.Repository) UseCase {
	return &useCaseImpl{
		productRepo: productRepo,
	}
}

// NewUseCaseWithCache creates a new product use case with cache support.
func NewUseCaseWithCache(productRepo product.Repository, cache cache.Cache) UseCase {
	return &useCaseImpl{
		productRepo: productRepo,
		cache:       cache,
	}
}

// ListProducts lists products with optional filters.
func (uc *useCaseImpl) ListProducts(ctx context.Context, req ListProductsRequest) ([]product.Product, int, error) {
	// TODO: Implement filtering and pagination
	if req.RestaurantID != "" {
		products, err := uc.productRepo.FindByRestaurant(ctx, req.RestaurantID)
		if err != nil {
			return nil, 0, err
		}
		return products, len(products), nil
	}
	return nil, 0, fmt.Errorf("not implemented: need to implement FindAll or similar")
}

// GetProduct retrieves a product by ID.
// If cache is available, it will check cache first before querying database.
func (uc *useCaseImpl) GetProduct(ctx context.Context, productID string) (*product.Product, error) {
	// Try cache first if available
	if uc.cache != nil {
		cacheKey := "product:" + productID

		// Try to get from cache
		cachedData, err := uc.cache.Get(ctx, cacheKey)
		if err == nil && cachedData != nil {
			var cached product.Product
			if json.Unmarshal(cachedData, &cached) == nil && cached.ID != "" {
				return &cached, nil
			}
		}

		// Cache miss or error - fetch from database
		p, err := uc.productRepo.FindByID(ctx, productID)
		if err != nil {
			return nil, err
		}

		// Store in cache with 5 minute TTL
		if p != nil {
			if data, err := json.Marshal(p); err == nil {
				_ = uc.cache.Set(ctx, cacheKey, data, 5*time.Minute)
			}
		}

		return p, nil
	}

	// No cache - direct database query
	return uc.productRepo.FindByID(ctx, productID)
}

// InvalidateProductCache invalidates cached product data.
// This should be called when product is updated or deleted.
func (uc *useCaseImpl) InvalidateProductCache(ctx context.Context, productID string) error {
	if uc.cache == nil {
		return nil
	}

	cacheKey := "product:" + productID
	return uc.cache.Delete(ctx, cacheKey)
}
