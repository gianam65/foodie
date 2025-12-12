package router

import (
	"net/http"
	"strings"

	"foodie/backend/internal/interfaces/http/handler"
	"foodie/backend/internal/interfaces/http/middleware"
	"foodie/backend/pkg/config"
)

var swaggerHandler = handler.NewSwaggerHandler()

// setupPublicRoutes registers public routes that don't require authentication.
func (r *Router) setupPublicRoutes(public *RouteGroup) {
	// Health check endpoints
	public.GET("/health", r.healthController.Check)
	public.GET("/ping", r.healthController.Ping)

	// Swagger/OpenAPI documentation (only if enabled)
	if config.GetBool("ENABLE_SWAGGER", true) {
		swaggerPath := config.Get("SWAGGER_PATH", "/swagger/")
		// Ensure path ends with /
		if swaggerPath != "" && swaggerPath[len(swaggerPath)-1] != '/' {
			swaggerPath += "/"
		}

		// Apply Swagger auth middleware if credentials are set
		swaggerAuth := middleware.SwaggerAuthMiddleware
		if config.Get("SWAGGER_USERNAME", "") == "" && config.Get("SWAGGER_PASSWORD", "") == "" {
			// No auth required if credentials not set
			swaggerAuth = func(next http.Handler) http.Handler { return next }
		}

		// Remove trailing slash for base path
		swaggerBasePath := strings.TrimSuffix(swaggerPath, "/")

		// Create wrapped Swagger UI handler with auth
		swaggerUIHandler := swaggerAuth(swaggerHandler.ServeSwaggerUI())

		// Register Swagger UI handler directly on Router's mux to handle all sub-paths
		// http-swagger needs to handle paths like /swagger/, /swagger/index.html, etc.
		// Must register on the main mux, not RouteGroup, to handle sub-path matching correctly
		swaggerPattern := swaggerBasePath + "/"
		r.mux.HandleFunc(swaggerPattern, func(w http.ResponseWriter, req *http.Request) {
			swaggerUIHandler.ServeHTTP(w, req)
		})

		// Register OpenAPI spec endpoints
		public.HandleFunc(swaggerBasePath+"/doc.json", func(w http.ResponseWriter, r *http.Request) {
			swaggerAuth(swaggerHandler.ServeOpenAPISpec()).ServeHTTP(w, r)
		})
		public.HandleFunc(swaggerBasePath+"/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
			swaggerAuth(swaggerHandler.ServeOpenAPISpecYAML()).ServeHTTP(w, r)
		})
		public.HandleFunc(swaggerBasePath+"/openapi.yml", func(w http.ResponseWriter, r *http.Request) {
			swaggerAuth(swaggerHandler.ServeOpenAPISpecYAML()).ServeHTTP(w, r)
		})
	}

	// Public product listing (anyone can view products)
	public.GET("/api/v1/products", r.productController.ListProducts)
}
