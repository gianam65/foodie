package product

import "time"

// Product represents a product/item in the restaurant menu.
type Product struct {
	ID           string
	RestaurantID string
	Name         string
	Price        float64
	CreatedAt    time.Time
}
