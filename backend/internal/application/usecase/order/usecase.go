package order

import (
	"context"

	"foodie/backend/internal/domain/order"
)

// UseCase defines use cases for order management.
// This is the application layer that orchestrates business logic.
type UseCase interface {
	// CreateOrder creates a new order.
	CreateOrder(ctx context.Context, cmd CreateOrderCommand) (*order.Order, error)

	// GetOrder retrieves an order by ID.
	GetOrder(ctx context.Context, orderID string) (*order.Order, error)

	// ListOrders lists orders with optional filters.
	ListOrders(ctx context.Context, req ListOrdersRequest) ([]order.Order, int, error)
}

// CreateOrderCommand represents the command to create an order.
type CreateOrderCommand struct {
	UserID          string
	RestaurantID    string
	Items           []OrderItemCommand
	PaymentMethod   string
	DeliveryAddress string
}

// OrderItemCommand represents an item in the create order command.
type OrderItemCommand struct {
	ProductID string
	Quantity  int
}

// ListOrdersRequest represents filters for listing orders.
type ListOrdersRequest struct {
	UserID string
	Status string
	Page   int // Page number (default: 1)
	Offset int // Offset (if provided, will be used directly; otherwise calculated from page)
	Limit  int // Items per page (default: 20)
}
