package dto

// CreateOrderRequest represents the request to create an order.
type CreateOrderRequest struct {
	UserID          string             `json:"user_id" validate:"required"`
	RestaurantID    string             `json:"restaurant_id" validate:"required"`
	Items           []OrderItemRequest `json:"items" validate:"required,min=1"`
	PaymentMethod   string             `json:"payment_method" validate:"required"`
	DeliveryAddress string             `json:"delivery_address" validate:"required"`
}

// OrderItemRequest represents an item in the order request.
type OrderItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

// OrderResponse represents an order in the API response.
type OrderResponse struct {
	ID        string              `json:"id"`
	UserID    string              `json:"user_id"`
	Status    string              `json:"status"`
	Total     float64             `json:"total,omitempty"`
	CreatedAt string              `json:"created_at"`
	Items     []OrderItemResponse `json:"items,omitempty"`
}

// OrderItemResponse represents an item in the order response.
type OrderItemResponse struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

// ListOrdersRequest represents query parameters for listing orders.
type ListOrdersRequest struct {
	UserID string `json:"user_id,omitempty"`
	Status string `json:"status,omitempty"`
	Page   int    `json:"page,omitempty"`   // Page number (default: 1)
	Offset int    `json:"offset,omitempty"` // Offset (default: 0, calculated from page if page provided)
	Limit  int    `json:"limit,omitempty"`  // Items per page (default: 20)
}

// ListOrdersResponse represents the response for listing orders with pagination.
type ListOrdersResponse struct {
	Data       []OrderResponse `json:"data"`
	Pagination PaginationMeta  `json:"pagination"`
}

// PaginationMeta contains pagination metadata.
type PaginationMeta struct {
	CurrentPage int  `json:"current_page"`
	PerPage     int  `json:"per_page"`
	Offset      int  `json:"offset"`
	Total       int  `json:"total"`
	TotalPages  int  `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrev     bool `json:"has_prev"`
}
