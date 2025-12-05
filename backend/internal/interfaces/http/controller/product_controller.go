package controller

import (
	"net/http"
	"time"

	productusecase "foodie/backend/internal/application/usecase/product"
	"foodie/backend/internal/domain/product"
	"foodie/backend/internal/interfaces/http/dto"
	httputils "foodie/backend/pkg/utils/http"
	"foodie/backend/pkg/utils/pagination"
)

// ProductController handles HTTP requests for product operations.
type ProductController struct {
	productUseCase productusecase.UseCase
}

// NewProductController creates a new product controller.
func NewProductController(productUseCase productusecase.UseCase) *ProductController {
	return &ProductController{
		productUseCase: productUseCase,
	}
}

// ListProducts handles GET /api/v1/products
func (c *ProductController) ListProducts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page := pagination.ParsePage(r.URL.Query().Get("page"))
	offset := pagination.ParseOffset(r.URL.Query().Get("offset"))
	limit := pagination.ParseLimit(r.URL.Query().Get("limit"), 20, 1, 100)

	// Calculate actual offset used (may be calculated from page)
	actualOffset := offset
	if actualOffset == 0 {
		actualOffset = pagination.CalculateOffset(page, limit)
	}

	// Convert DTO request to use case request
	useCaseReq := productusecase.ListProductsRequest{
		RestaurantID: r.URL.Query().Get("restaurant_id"),
		Page:         page,
		Offset:       actualOffset,
		Limit:        limit,
	}

	// Call use case
	products, total, err := c.productUseCase.ListProducts(r.Context(), useCaseReq)
	if err != nil {
		httputils.InternalServerError(w, "Failed to list products", err)
		return
	}

	// Calculate pagination metadata
	paginationMeta := pagination.CalculateMeta(page, limit, total)

	// Convert domain entities to DTOs
	productDTOs := make([]dto.ProductResponse, 0, len(products))
	for _, p := range products {
		productDTOs = append(productDTOs, c.productToDTO(&p))
	}

	// Respond with pagination
	httputils.Success(w, dto.ListProductsResponse{
		Data: productDTOs,
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

// productToDTO converts domain Product entity to DTO.
func (c *ProductController) productToDTO(p *product.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:           p.ID,
		RestaurantID: p.RestaurantID,
		Name:         p.Name,
		Price:        p.Price,
		CreatedAt:    p.CreatedAt.Format(time.RFC3339),
	}
}
