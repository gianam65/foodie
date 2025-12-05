package router

import (
	"net/http"
	"sync"

	"foodie/backend/internal/interfaces/http/middleware"
)

// methodHandlers stores handlers for different HTTP methods on the same pattern
type methodHandlers struct {
	mu         sync.RWMutex
	handlers   map[string]http.Handler
	registered sync.Once // Ensure dispatcher is registered only once
}

// RouteGroup represents a group of routes with shared middleware.
type RouteGroup struct {
	prefix            string
	middleware        []middleware.Middleware
	mux               *http.ServeMux
	parent            *Router
	handlersByPattern map[string]*methodHandlers // Pattern -> method -> handler
	mu                sync.Mutex
}

// RouteGroup creates a new route group with a prefix and optional middleware.
func (r *Router) RouteGroup(prefix string, mw ...middleware.Middleware) *RouteGroup {
	return &RouteGroup{
		prefix:            prefix,
		middleware:        mw,
		mux:               r.mux,
		parent:            r,
		handlersByPattern: make(map[string]*methodHandlers),
	}
}

// GET registers a GET route in the group.
func (rg *RouteGroup) GET(pattern string, handler http.HandlerFunc) {
	rg.Handle(http.MethodGet, pattern, handler)
}

// POST registers a POST route in the group.
func (rg *RouteGroup) POST(pattern string, handler http.HandlerFunc) {
	rg.Handle(http.MethodPost, pattern, handler)
}

// PUT registers a PUT route in the group.
func (rg *RouteGroup) PUT(pattern string, handler http.HandlerFunc) {
	rg.Handle(http.MethodPut, pattern, handler)
}

// DELETE registers a DELETE route in the group.
func (rg *RouteGroup) DELETE(pattern string, handler http.HandlerFunc) {
	rg.Handle(http.MethodDelete, pattern, handler)
}

// Handle registers a route with a specific HTTP method in the group.
func (rg *RouteGroup) Handle(method, pattern string, handler http.HandlerFunc) {
	fullPattern := rg.prefix + pattern

	// Apply group middleware
	var h http.Handler = http.HandlerFunc(handler)
	for i := len(rg.middleware) - 1; i >= 0; i-- {
		h = rg.middleware[i](h)
	}

	// Store handler in map by pattern and method
	rg.mu.Lock()
	mh, exists := rg.handlersByPattern[fullPattern]
	if !exists {
		mh = &methodHandlers{
			handlers: make(map[string]http.Handler),
		}
		rg.handlersByPattern[fullPattern] = mh
	}
	rg.mu.Unlock()

	// Add handler for this method
	mh.mu.Lock()
	mh.handlers[method] = h
	mh.mu.Unlock()

	// Register wrapper handler only once per pattern using sync.Once
	mh.registered.Do(func() {
		// Register a dispatcher that routes based on HTTP method
		rg.mux.HandleFunc(fullPattern, func(w http.ResponseWriter, r *http.Request) {
			mh.mu.RLock()
			handler, ok := mh.handlers[r.Method]
			mh.mu.RUnlock()

			if ok {
				handler.ServeHTTP(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})
	})
}

// HandleFunc registers a route that accepts any HTTP method in the group.
func (rg *RouteGroup) HandleFunc(pattern string, handler http.HandlerFunc) {
	fullPattern := rg.prefix + pattern

	// Apply group middleware
	var h http.Handler = http.HandlerFunc(handler)
	for i := len(rg.middleware) - 1; i >= 0; i-- {
		h = rg.middleware[i](h)
	}

	rg.mux.HandleFunc(fullPattern, func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}
