package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"

	"foodie/backend/pkg/config"
)

// SwaggerAuthMiddleware provides basic authentication for Swagger UI.
// It reads SWAGGER_USERNAME and SWAGGER_PASSWORD from environment variables.
func SwaggerAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := config.Get("SWAGGER_USERNAME", "")
		password := config.Get("SWAGGER_PASSWORD", "")

		// If no credentials are set, allow access (for development)
		if username == "" && password == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Get Authorization header
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Swagger UI"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse Basic Auth
		if !strings.HasPrefix(auth, "Basic ") {
			w.Header().Set("WWW-Authenticate", `Basic realm="Swagger UI"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Decode credentials
		encoded := strings.TrimPrefix(auth, "Basic ")
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Swagger UI"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Split username:password
		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Swagger UI"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Verify credentials
		if parts[0] != username || parts[1] != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Swagger UI"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Authentication successful
		next.ServeHTTP(w, r)
	})
}
