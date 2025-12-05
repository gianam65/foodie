package router

import (
	"net/http"

	"foodie/backend/internal/interfaces/http/controller"
	"foodie/backend/internal/interfaces/http/middleware"
	"foodie/backend/pkg/logger"
)

// Router sets up HTTP routes and delegates to controllers.
type Router struct {
	mux               *http.ServeMux
	logger            *logger.Logger
	healthController  *controller.HealthController
	orderController   *controller.OrderController
	productController *controller.ProductController
}

// NewRouter creates a new HTTP router with controllers and logger.
func NewRouter(
	logger *logger.Logger,
	healthController *controller.HealthController,
	orderController *controller.OrderController,
	productController *controller.ProductController,
) *Router {
	return &Router{
		mux:               http.NewServeMux(),
		logger:            logger,
		healthController:  healthController,
		orderController:   orderController,
		productController: productController,
	}
}

// SetupRoutes configures all HTTP routes with public and private route groups.
func (r *Router) SetupRoutes() {
	// Correlation ID middleware should be first to ensure all logs have request ID
	correlationMiddleware := middleware.CorrelationIDMiddleware

	// Apply global middleware to all routes
	globalMiddleware := []middleware.Middleware{
		correlationMiddleware,
		middleware.RecoveryMiddleware(r.logger),
		middleware.LoggingMiddleware(r.logger),
		middleware.CORSMiddleware,
	}

	// Public routes (no authentication required)
	public := r.RouteGroup("", globalMiddleware...)
	r.setupPublicRoutes(public)

	// Private routes (authentication required)
	private := r.RouteGroup("/api/v1", append(globalMiddleware, middleware.AuthMiddleware)...)
	r.setupPrivateRoutes(private)
}

// ServeHTTP implements http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
