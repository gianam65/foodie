package middleware

import "net/http"

// Middleware is a function that wraps an HTTP handler.
type Middleware func(http.Handler) http.Handler

// Chain chains multiple middlewares together.
// Middlewares are applied in reverse order (last middleware wraps first).
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
