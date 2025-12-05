package order

import (
	"context"
	"fmt"
	"time"

	"foodie/backend/internal/domain/order"
	productrepo "foodie/backend/internal/domain/product"

	"github.com/google/uuid"
)

// useCaseImpl implements the UseCase interface.
type useCaseImpl struct {
	orderRepo   order.Repository
	productRepo productrepo.Repository
}

// NewUseCase creates a new order use case.
func NewUseCase(orderRepo order.Repository, productRepo productrepo.Repository) UseCase {
	return &useCaseImpl{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

// CreateOrder creates a new order.
func (uc *useCaseImpl) CreateOrder(ctx context.Context, cmd CreateOrderCommand) (*order.Order, error) {
	// 1. Validate command
	if err := uc.validateCreateCommand(cmd); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Fetch products and calculate total
	items := make([]order.OrderItem, 0, len(cmd.Items))
	var total float64

	for _, itemCmd := range cmd.Items {
		// Fetch product to get price and name
		product, err := uc.productRepo.FindByID(ctx, itemCmd.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product not found: %s: %w", itemCmd.ProductID, err)
		}

		itemTotal := product.Price * float64(itemCmd.Quantity)
		total += itemTotal

		items = append(items, order.OrderItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    itemCmd.Quantity,
			Price:       product.Price,
		})
	}

	// 3. Create order entity
	now := time.Now()
	orderEntity := &order.Order{
		ID:              uuid.New().String(),
		UserID:          cmd.UserID,
		RestaurantID:    cmd.RestaurantID,
		Status:          order.StatusPending,
		Items:           items,
		Total:           total,
		PaymentMethod:   cmd.PaymentMethod,
		DeliveryAddress: cmd.DeliveryAddress,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// 4. Save via repository
	if err := uc.orderRepo.Save(ctx, orderEntity); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	// 5. TODO: Emit domain events (order created)

	return orderEntity, nil
}

// validateCreateCommand validates the create order command.
func (uc *useCaseImpl) validateCreateCommand(cmd CreateOrderCommand) error {
	if cmd.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if cmd.RestaurantID == "" {
		return fmt.Errorf("restaurant_id is required")
	}
	if len(cmd.Items) == 0 {
		return fmt.Errorf("order must have at least one item")
	}
	for i, item := range cmd.Items {
		if item.ProductID == "" {
			return fmt.Errorf("items[%d].product_id is required", i)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("items[%d].quantity must be greater than 0", i)
		}
	}
	if cmd.PaymentMethod == "" {
		return fmt.Errorf("payment_method is required")
	}
	if cmd.DeliveryAddress == "" {
		return fmt.Errorf("delivery_address is required")
	}
	return nil
}

// GetOrder retrieves an order by ID.
func (uc *useCaseImpl) GetOrder(ctx context.Context, orderID string) (*order.Order, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order_id is required")
	}

	o, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	return o, nil
}

// ListOrders lists orders with optional filters.
func (uc *useCaseImpl) ListOrders(ctx context.Context, req ListOrdersRequest) ([]order.Order, int, error) {
	// Validate and set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}

	// Calculate offset: use provided offset if specified, otherwise calculate from page
	offset := req.Offset
	if offset == 0 && req.Page > 0 {
		offset = (req.Page - 1) * req.Limit
	}

	var orders []order.Order
	var total int
	var err error

	// Filter by user_id if provided
	if req.UserID != "" {
		orders, err = uc.orderRepo.FindByUserID(ctx, req.UserID, req.Limit, offset)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to fetch orders: %w", err)
		}

		total, err = uc.orderRepo.CountByUserID(ctx, req.UserID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count orders: %w", err)
		}

		// Filter by status if provided
		if req.Status != "" {
			filtered := make([]order.Order, 0)
			for _, o := range orders {
				if string(o.Status) == req.Status {
					filtered = append(filtered, o)
				}
			}
			orders = filtered
		}
	} else {
		// TODO: Implement FindAll if needed
		return nil, 0, fmt.Errorf("user_id filter is required for listing orders")
	}

	return orders, total, nil
}
