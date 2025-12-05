package middleware

import (
	"context"
	"net/http"
	"strings"
)

// ContextKey is a type for context keys to avoid collisions.
type ContextKey string

const (
	// UserIDKey is the context key for user ID.
	UserIDKey ContextKey = "user_id"
	// UserRoleKey is the context key for user role.
	UserRoleKey ContextKey = "user_role"
)

// AuthMiddleware validates JWT tokens and extracts user information.
// This is a placeholder implementation - replace with actual JWT validation.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		if token == "" {
			http.Error(w, "Token is required", http.StatusUnauthorized)
			return
		}

		// TODO: Validate JWT token
		// For now, this is a placeholder that accepts any non-empty token
		// In production, you should:
		// 1. Parse and validate the JWT
		// 2. Extract user ID and role from claims
		// 3. Verify token signature and expiration
		// 4. Check if user exists and is active

		// Placeholder: Extract user info from token
		// In real implementation, decode JWT and extract claims
		userID := extractUserIDFromToken(token)     // Placeholder
		userRole := extractUserRoleFromToken(token) // Placeholder

		if userID == "" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user info to request context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserRoleKey, userRole)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// extractUserIDFromToken is a placeholder - replace with actual JWT parsing.
func extractUserIDFromToken(token string) string {
	// TODO: Parse JWT and extract user_id claim
	// For now, return empty string to force authentication
	return ""
}

// extractUserRoleFromToken is a placeholder - replace with actual JWT parsing.
func extractUserRoleFromToken(token string) string {
	// TODO: Parse JWT and extract role claim
	return "user"
}

// OptionalAuthMiddleware allows requests with or without authentication.
// If a valid token is present, user info is added to context.
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				if token != "" {
					// Try to extract user info, but don't fail if invalid
					userID := extractUserIDFromToken(token)
					if userID != "" {
						userRole := extractUserRoleFromToken(token)
						ctx := context.WithValue(r.Context(), UserIDKey, userID)
						ctx = context.WithValue(ctx, UserRoleKey, userRole)
						r = r.WithContext(ctx)
					}
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

// RoleMiddleware restricts access to specific roles.
func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(UserRoleKey).(string)
			if !ok {
				http.Error(w, "Forbidden: role required", http.StatusForbidden)
				return
			}

			allowed := false
			for _, role := range allowedRoles {
				if userRole == role {
					allowed = true
					break
				}
			}

			if !allowed {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserID extracts user ID from request context.
func GetUserID(r *http.Request) string {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}

// GetUserRole extracts user role from request context.
func GetUserRole(r *http.Request) string {
	userRole, ok := r.Context().Value(UserRoleKey).(string)
	if !ok {
		return ""
	}
	return userRole
}
