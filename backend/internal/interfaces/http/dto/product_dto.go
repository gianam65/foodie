package dto

// ProductResponse represents a product in the API response.
type ProductResponse struct {
	ID           string  `json:"id"`
	RestaurantID string  `json:"restaurant_id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	CreatedAt    string  `json:"created_at"`
}

// ListProductsRequest represents query parameters for listing products.
type ListProductsRequest struct {
	RestaurantID string `json:"restaurant_id,omitempty"`
	Page         int    `json:"page,omitempty"`   // Page number (default: 1)
	Offset       int    `json:"offset,omitempty"` // Offset (default: 0, calculated from page if page provided)
	Limit        int    `json:"limit,omitempty"`  // Items per page (default: 20)
}

// ListProductsResponse represents the response for listing products with pagination.
type ListProductsResponse struct {
	Data       []ProductResponse `json:"data"`
	Pagination PaginationMeta    `json:"pagination"`
}
