package router

import (
	"foodie/backend/internal/interfaces/http/handler"
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

		public.HandleFunc(swaggerPath, swaggerHandler.ServeSwaggerUI())
		public.HandleFunc(swaggerPath+"doc.json", swaggerHandler.ServeOpenAPISpec())
		public.HandleFunc(swaggerPath+"openapi.yaml", swaggerHandler.ServeOpenAPISpecYAML())
		public.HandleFunc(swaggerPath+"openapi.yml", swaggerHandler.ServeOpenAPISpecYAML())
	}

	// Public product listing (anyone can view products)
	public.GET("/api/v1/products", r.productController.ListProducts)
}
