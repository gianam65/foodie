package product

import (
	"context"

	"foodie/backend/internal/domain/product"
)

// UseCase defines use cases for product management.
type UseCase interface {
	// ListProducts lists products with optional filters.
	ListProducts(ctx context.Context, req ListProductsRequest) ([]product.Product, int, error)

	// GetProduct retrieves a product by ID.
	GetProduct(ctx context.Context, productID string) (*product.Product, error)

	// InvalidateProductCache invalidates cached product data.
	InvalidateProductCache(ctx context.Context, productID string) error
}

// ListProductsRequest represents filters for listing products.
type ListProductsRequest struct {
	RestaurantID string
	Page         int // Page number (default: 1)
	Offset       int // Offset (if provided, will be used directly; otherwise calculated from page)
	Limit        int // Items per page (default: 20)
}
