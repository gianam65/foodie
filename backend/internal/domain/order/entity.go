package order

import "time"

// OrderStatus represents the status of an order.
type OrderStatus string

const (
	StatusPending    OrderStatus = "pending"
	StatusConfirmed  OrderStatus = "confirmed"
	StatusPreparing  OrderStatus = "preparing"
	StatusReady      OrderStatus = "ready"
	StatusDelivering OrderStatus = "delivering"
	StatusCompleted  OrderStatus = "completed"
	StatusCancelled  OrderStatus = "cancelled"
)

// OrderItem represents an item in an order.
type OrderItem struct {
	ProductID   string
	ProductName string
	Quantity    int
	Price       float64
}

// Order represents a food order aggregate in the domain layer.
type Order struct {
	ID              string
	UserID          string
	RestaurantID    string
	Status          OrderStatus
	Items           []OrderItem
	Total           float64
	PaymentMethod   string
	DeliveryAddress string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
