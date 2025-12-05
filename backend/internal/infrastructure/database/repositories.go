package database

import (
	"database/sql"

	"foodie/backend/internal/domain/order"
	"foodie/backend/internal/domain/product"
	orderrepo "foodie/backend/internal/infrastructure/database/order"
	productrepo "foodie/backend/internal/infrastructure/database/product"
)

// Repositories bundles every repository implementation the application needs.
// As new modules (user, restaurant, payment, etc.) are implemented,
// add them here and initialize them in NewRepositories.
type Repositories struct {
	Order   order.Repository
	Product product.Repository
	// User       user.Repository
	// Restaurant restaurant.Repository
	// Payment    payment.Repository
}

// NewRepositories creates and initializes all repositories for the application.
//
// When adding a new module (e.g., user, restaurant), follow these steps:
// 1. Create domain/<module> with Repository interface
// 2. Create infrastructure/database/<module>/repository.go with implementation
// 3. Add field to Repositories struct above
// 4. Initialize it here: repos.Module = modulerepo.NewRepository(sqlDB)
//
// Note: Database connection should be created separately using NewConnection() or NewConnectionFromEnv()
func NewRepositories(sqlDB *sql.DB) (*Repositories, error) {
	if sqlDB == nil {
		return nil, &ErrInvalidDatabase{}
	}

	return &Repositories{
		Order:   orderrepo.NewRepository(sqlDB),
		Product: productrepo.NewRepository(sqlDB),
	}, nil
}

// ErrInvalidDatabase is returned when SQL database connection is nil.
type ErrInvalidDatabase struct{}

func (e *ErrInvalidDatabase) Error() string {
	return "sql database connection is nil"
}
