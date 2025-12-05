package cache

import (
	"context"
	"time"
)

// Cache defines the interface for caching operations.
// This abstraction allows swapping implementations (Redis, in-memory, etc.)
type Cache interface {
	// Get retrieves a value by key. Returns nil if key doesn't exist.
	Get(ctx context.Context, key string) ([]byte, error)

	// Set stores a value with expiration time.
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Delete removes a key from cache.
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists in cache.
	Exists(ctx context.Context, key string) (bool, error)

	// SetNX sets a key only if it doesn't exist (useful for distributed locks).
	SetNX(ctx context.Context, key string, value []byte, ttl time.Duration) (bool, error)

	// Close closes the cache connection.
	Close() error
}
