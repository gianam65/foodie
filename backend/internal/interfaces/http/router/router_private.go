package router

import (
	"net/http"
	"strings"
)

// setupPrivateRoutes registers private routes that require authentication.
func (r *Router) setupPrivateRoutes(private *RouteGroup) {
	// Order routes (require authentication)
	// Use single handler for GET /orders that handles both list and get-by-ID
	private.GET("/orders", r.handleOrders)
	// POST /api/v1/orders - Create order
	private.POST("/orders", r.orderController.CreateOrder)
}

// handleOrders routes GET requests to /api/v1/orders
// Handles both list orders and get order by ID based on the path
func (r *Router) handleOrders(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract path parts: /api/v1/orders or /api/v1/orders/{id}
	pathParts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")

	// If path has 4+ parts and the 4th part (index 3) is not empty, it's a get-by-ID request
	// Example: /api/v1/orders/123 -> pathParts = ["api", "v1", "orders", "123"]
	if len(pathParts) >= 4 && pathParts[2] == "orders" && pathParts[3] != "" {
		// Get order by ID
		r.orderController.GetOrder(w, req)
		return
	}

	// Otherwise, it's a list request
	// Example: /api/v1/orders -> pathParts = ["api", "v1", "orders"]
	r.orderController.ListOrders(w, req)
}
