package product

import (
	"context"
	"database/sql"

	"foodie/backend/internal/domain/product"
)

// Repository implements product.Repository using SQL.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new SQL-based product repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Save inserts a new product row.
func (r *Repository) Save(ctx context.Context, p *product.Product) error {
	const query = `INSERT INTO products (id, restaurant_id, name, price, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.RestaurantID, p.Name, p.Price, p.CreatedAt)
	return err
}

// FindByID loads a product by its ID.
func (r *Repository) FindByID(ctx context.Context, id string) (*product.Product, error) {
	const query = `SELECT id, restaurant_id, name, price, created_at FROM products WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var p product.Product
	if err := row.Scan(&p.ID, &p.RestaurantID, &p.Name, &p.Price, &p.CreatedAt); err != nil {
		return nil, err
	}
	return &p, nil
}

// FindByRestaurant loads all products for a restaurant.
func (r *Repository) FindByRestaurant(ctx context.Context, restaurantID string) ([]product.Product, error) {
	const query = `SELECT id, restaurant_id, name, price, created_at FROM products WHERE restaurant_id = $1`
	rows, err := r.db.QueryContext(ctx, query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []product.Product
	for rows.Next() {
		var p product.Product
		if err := rows.Scan(&p.ID, &p.RestaurantID, &p.Name, &p.Price, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}
