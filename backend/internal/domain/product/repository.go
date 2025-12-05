package product

import "context"

// Repository defines storage operations for products.
type Repository interface {
	Save(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, id string) (*Product, error)
	FindByRestaurant(ctx context.Context, restaurantID string) ([]Product, error)
}
