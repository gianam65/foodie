package order

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"foodie/backend/internal/domain/order"
)

// Repository is a Postgres/MySQL-style implementation of order.Repository.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new SQL-based order repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Save inserts a new order row.
func (r *Repository) Save(ctx context.Context, o *order.Order) error {
	// Serialize items to JSON
	itemsJSON, err := json.Marshal(o.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}

	const query = `INSERT INTO orders (
		id, user_id, restaurant_id, status, items, total, 
		payment_method, delivery_address, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err = r.db.ExecContext(ctx, query,
		o.ID, o.UserID, o.RestaurantID, string(o.Status),
		itemsJSON, o.Total, o.PaymentMethod, o.DeliveryAddress,
		o.CreatedAt, o.UpdatedAt,
	)
	return err
}

// FindByID loads an order by its ID.
func (r *Repository) FindByID(ctx context.Context, id string) (*order.Order, error) {
	const query = `SELECT 
		id, user_id, restaurant_id, status, items, total,
		payment_method, delivery_address, created_at, updated_at
		FROM orders WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	var o order.Order
	var statusStr string
	var itemsJSON []byte
	var createdAt, updatedAt time.Time

	err := row.Scan(
		&o.ID, &o.UserID, &o.RestaurantID, &statusStr,
		&itemsJSON, &o.Total, &o.PaymentMethod, &o.DeliveryAddress,
		&createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	o.Status = order.OrderStatus(statusStr)
	o.CreatedAt = createdAt
	o.UpdatedAt = updatedAt

	// Deserialize items
	if err := json.Unmarshal(itemsJSON, &o.Items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	return &o, nil
}

// FindByUserID loads orders by user ID with pagination.
func (r *Repository) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]order.Order, error) {
	const query = `SELECT 
		id, user_id, restaurant_id, status, items, total,
		payment_method, delivery_address, created_at, updated_at
		FROM orders 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []order.Order
	for rows.Next() {
		var o order.Order
		var statusStr string
		var itemsJSON []byte
		var createdAt, updatedAt time.Time

		if err := rows.Scan(
			&o.ID, &o.UserID, &o.RestaurantID, &statusStr,
			&itemsJSON, &o.Total, &o.PaymentMethod, &o.DeliveryAddress,
			&createdAt, &updatedAt,
		); err != nil {
			return nil, err
		}

		o.Status = order.OrderStatus(statusStr)
		o.CreatedAt = createdAt
		o.UpdatedAt = updatedAt

		if err := json.Unmarshal(itemsJSON, &o.Items); err != nil {
			return nil, fmt.Errorf("failed to unmarshal items: %w", err)
		}

		orders = append(orders, o)
	}

	return orders, rows.Err()
}

// CountByUserID counts orders by user ID.
func (r *Repository) CountByUserID(ctx context.Context, userID string) (int, error) {
	const query = `SELECT COUNT(*) FROM orders WHERE user_id = $1`
	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}
