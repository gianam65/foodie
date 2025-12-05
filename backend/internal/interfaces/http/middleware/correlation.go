package middleware

import (
	"net/http"

	"github.com/google/uuid"
)

const CorrelationIDHeader = "X-Correlation-ID"
const CorrelationIDContextKey = "correlation_id"

// CorrelationIDMiddleware adds a correlation ID to each request for tracing.
// This ID is used to correlate logs across services and is essential for Grafana/Loki querying.
func CorrelationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if correlation ID exists in request header
		correlationID := r.Header.Get(CorrelationIDHeader)

		// If not present, generate a new one
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		// Add correlation ID to response header
		w.Header().Set(CorrelationIDHeader, correlationID)

		// Add to request context for use in handlers
		ctx := r.Context()
		ctx = setCorrelationID(ctx, correlationID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// GetCorrelationID extracts correlation ID from request context.
func GetCorrelationID(r *http.Request) string {
	return getCorrelationID(r.Context())
}
