package controller

import (
	"net/http"
	"time"

	httputils "foodie/backend/pkg/utils/http"
	timutils "foodie/backend/pkg/utils/time"
)

// HealthController handles health check requests.
type HealthController struct{}

// NewHealthController creates a new health controller.
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Check handles GET /health
func (c *HealthController) Check(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "ok",
		"service":   "foodie-backend",
		"timestamp": timutils.FormatRFC3339(time.Now().UTC()),
	}

	httputils.Success(w, response)
}

// Ping handles GET /ping (simpler health check, returns "pong")
func (c *HealthController) Ping(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"pong": timutils.FormatRFC3339(time.Now().UTC()),
	}

	httputils.Success(w, response)
}
