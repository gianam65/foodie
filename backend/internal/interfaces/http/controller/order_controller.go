package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	orderusecase "foodie/backend/internal/application/usecase/order"
	"foodie/backend/internal/domain/order"
	"foodie/backend/internal/interfaces/http/dto"
	httputils "foodie/backend/pkg/utils/http"
	"foodie/backend/pkg/utils/pagination"
)

// OrderController handles HTTP requests for order operations.
type OrderController struct {
	orderUseCase orderusecase.UseCase
}

// NewOrderController creates a new order controller.
func NewOrderController(orderUseCase orderusecase.UseCase) *OrderController {
	return &OrderController{
		orderUseCase: orderUseCase,
	}
}

// CreateOrder handles POST /api/v1/orders
func (c *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.BadRequest(w, "Invalid request body", err)
		return
	}

	// Convert DTO to use case command
	cmd := orderusecase.CreateOrderCommand{
		UserID:          req.UserID,
		RestaurantID:    req.RestaurantID,
		PaymentMethod:   req.PaymentMethod,
		DeliveryAddress: req.DeliveryAddress,
	}

	// Convert order items
	for _, item := range req.Items {
		cmd.Items = append(cmd.Items, orderusecase.OrderItemCommand{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	// Call use case
	createdOrder, err := c.orderUseCase.CreateOrder(r.Context(), cmd)
	if err != nil {
		// Check error type and return appropriate status code
		if strings.Contains(err.Error(), "validation failed") {
			httputils.BadRequest(w, "Validation failed", err)
			return
		}
		if strings.Contains(err.Error(), "not found") {
			httputils.NotFound(w, "Product not found")
			return
		}
		httputils.InternalServerError(w, "Failed to create order", err)
		return
	}

	// Convert domain entity to DTO
	response := c.orderToDTO(createdOrder)

	httputils.Created(w, response)
}

// GetOrder handles GET /api/v1/orders/{id}
func (c *OrderController) GetOrder(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path: /api/v1/orders/{id}
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		httputils.BadRequest(w, "Invalid order ID", nil)
		return
	}
	orderID := pathParts[3]

	// Call use case
	o, err := c.orderUseCase.GetOrder(r.Context(), orderID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httputils.NotFound(w, "Order not found")
			return
		}
		httputils.InternalServerError(w, "Failed to get order", err)
		return
	}

	// Convert to DTO and respond
	response := c.orderToDTO(o)
	httputils.Success(w, response)
}

// ListOrders handles GET /api/v1/orders
func (c *OrderController) ListOrders(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page := pagination.ParsePage(r.URL.Query().Get("page"))
	offset := pagination.ParseOffset(r.URL.Query().Get("offset"))
	limit := pagination.ParseLimit(r.URL.Query().Get("limit"), 20, 1, 100)

	// Convert DTO request to use case request
	useCaseReq := orderusecase.ListOrdersRequest{
		UserID: r.URL.Query().Get("user_id"),
		Status: r.URL.Query().Get("status"),
		Page:   page,
		Offset: offset,
		Limit:  limit,
	}

	// Call use case
	orders, total, err := c.orderUseCase.ListOrders(r.Context(), useCaseReq)
	if err != nil {
		if strings.Contains(err.Error(), "required") {
			httputils.BadRequest(w, "user_id filter is required", err)
			return
		}
		httputils.InternalServerError(w, "Failed to list orders", err)
		return
	}

	// Calculate actual offset used (may be calculated from page)
	actualOffset := offset
	if actualOffset == 0 {
		actualOffset = pagination.CalculateOffset(page, limit)
	}

	// Calculate pagination metadata
	paginationMeta := pagination.CalculateMeta(page, limit, total)

	// Convert domain entities to DTOs
	orderDTOs := make([]dto.OrderResponse, 0, len(orders))
	for _, o := range orders {
		orderDTOs = append(orderDTOs, c.orderToDTO(&o))
	}

	// Respond with pagination
	httputils.Success(w, dto.ListOrdersResponse{
		Data: orderDTOs,
		Pagination: dto.PaginationMeta{
			CurrentPage: paginationMeta.CurrentPage,
			PerPage:     paginationMeta.PerPage,
			Offset:      actualOffset,
			Total:       paginationMeta.Total,
			TotalPages:  paginationMeta.TotalPages,
			HasNext:     paginationMeta.HasNext,
			HasPrev:     paginationMeta.HasPrev,
		},
	})
}

// orderToDTO converts domain Order entity to DTO.
func (c *OrderController) orderToDTO(o *order.Order) dto.OrderResponse {
	items := make([]dto.OrderItemResponse, 0, len(o.Items))
	for _, item := range o.Items {
		items = append(items, dto.OrderItemResponse{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
		})
	}

	return dto.OrderResponse{
		ID:        o.ID,
		UserID:    o.UserID,
		Status:    string(o.Status),
		Total:     o.Total,
		CreatedAt: o.CreatedAt.Format(time.RFC3339),
		Items:     items,
	}
}
