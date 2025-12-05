package order

import "context"

// Repository defines the storage operations required by the Order use cases.
type Repository interface {
	Save(ctx context.Context, order *Order) error
	FindByID(ctx context.Context, id string) (*Order, error)
	FindByUserID(ctx context.Context, userID string, limit, offset int) ([]Order, error)
	CountByUserID(ctx context.Context, userID string) (int, error)
}
