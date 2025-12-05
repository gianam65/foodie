package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	httpSwagger "github.com/swaggo/http-swagger"
	"gopkg.in/yaml.v3"
)

// SwaggerHandler handles Swagger UI and OpenAPI spec serving.
type SwaggerHandler struct {
	specPath    string
	openAPIPath string
}

// NewSwaggerHandler creates a new Swagger handler.
// It automatically finds the OpenAPI spec file (api/openapi.yaml) from the project root.
func NewSwaggerHandler() *SwaggerHandler {
	// Get current working directory
	wd, _ := os.Getwd()
	var openAPIPath string

	// Try to find api/openapi.yaml by walking up from current directory
	currentDir := wd
	for i := 0; i < 10; i++ { // Max 10 levels up
		testPath := filepath.Join(currentDir, "api", "openapi.yaml")
		if _, err := os.Stat(testPath); err == nil {
			absPath, _ := filepath.Abs(testPath)
			openAPIPath = absPath
			break
		}
		// Go up one level
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break // Reached root
		}
		currentDir = parent
	}

	return &SwaggerHandler{
		specPath:    "/swagger/doc.json",
		openAPIPath: openAPIPath,
	}
}

// ServeSwaggerUI serves Swagger UI at /swagger/
func (h *SwaggerHandler) ServeSwaggerUI() http.HandlerFunc {
	// Configure http-swagger to use our custom spec endpoint
	handler := httpSwagger.Handler(
		httpSwagger.URL(h.specPath), // URL to OpenAPI spec
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)
	return handler
}

// ServeOpenAPISpec serves the OpenAPI specification as JSON.
// Converts YAML to JSON for Swagger UI compatibility.
func (h *SwaggerHandler) ServeOpenAPISpec() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h.openAPIPath == "" {
			http.Error(w, "OpenAPI spec file not found. Please ensure api/openapi.yaml exists.", http.StatusNotFound)
			return
		}

		// Read OpenAPI spec file
		specData, err := os.ReadFile(h.openAPIPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to read OpenAPI spec: %v", err), http.StatusInternalServerError)
			return
		}

		// Parse YAML
		var data interface{}
		if err := yaml.Unmarshal(specData, &data); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse OpenAPI spec: %v", err), http.StatusInternalServerError)
			return
		}

		// Convert to JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to convert to JSON: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow CORS for Swagger UI
		w.Write(jsonData)
	}
}

// ServeOpenAPISpecYAML serves the OpenAPI specification as YAML.
func (h *SwaggerHandler) ServeOpenAPISpecYAML() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h.openAPIPath == "" {
			http.Error(w, "OpenAPI spec file not found. Please ensure api/openapi.yaml exists.", http.StatusNotFound)
			return
		}

		// Read OpenAPI spec file
		specData, err := os.ReadFile(h.openAPIPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to read OpenAPI spec: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/x-yaml")
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow CORS for Swagger UI
		w.Write(specData)
	}
}
