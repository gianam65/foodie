package middleware

import (
	"net/http"
	"runtime/debug"

	"foodie/backend/pkg/logger"

	"go.uber.org/zap"
)

// RecoveryMiddleware recovers from panics and logs with structured format.
func RecoveryMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					correlationID := GetCorrelationID(r)
					log.Error("panic_recovered",
						zap.Any("error", err),
						zap.String("correlation_id", correlationID),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.ByteString("stack", debug.Stack()),
					)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
