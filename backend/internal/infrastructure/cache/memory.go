package cache

import (
	"context"
	"sync"
	"time"
)

// MemoryCache is an in-memory cache implementation for development/testing.
// NOT suitable for production or distributed systems.
type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]*cacheItem
}

type cacheItem struct {
	value      []byte
	expiration time.Time
}

// NewMemoryCache creates a new in-memory cache.
func NewMemoryCache() *MemoryCache {
	c := &MemoryCache{
		items: make(map[string]*cacheItem),
	}

	// Start background cleanup goroutine
	go c.cleanup()

	return c
}

// Get retrieves a value by key.
func (m *MemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.items[key]
	if !exists {
		return nil, nil
	}

	// Check expiration
	if time.Now().After(item.expiration) {
		delete(m.items, key)
		return nil, nil
	}

	// Return a copy to prevent external modification
	result := make([]byte, len(item.value))
	copy(result, item.value)
	return result, nil
}

// Set stores a value with expiration time.
func (m *MemoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[key] = &cacheItem{
		value:      append([]byte(nil), value...), // Copy value
		expiration: time.Now().Add(ttl),
	}
	return nil
}

// Delete removes a key from cache.
func (m *MemoryCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, key)
	return nil
}

// Exists checks if a key exists in cache.
func (m *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.items[key]
	if !exists {
		return false, nil
	}

	// Check expiration
	if time.Now().After(item.expiration) {
		return false, nil
	}

	return true, nil
}

// SetNX sets a key only if it doesn't exist.
func (m *MemoryCache) SetNX(ctx context.Context, key string, value []byte, ttl time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if exists and not expired
	if item, exists := m.items[key]; exists {
		if time.Now().Before(item.expiration) {
			return false, nil // Key already exists
		}
	}

	m.items[key] = &cacheItem{
		value:      append([]byte(nil), value...),
		expiration: time.Now().Add(ttl),
	}
	return true, nil
}

// Close closes the cache (no-op for in-memory cache).
func (m *MemoryCache) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items = make(map[string]*cacheItem)
	return nil
}

// cleanup periodically removes expired items.
func (m *MemoryCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for key, item := range m.items {
			if now.After(item.expiration) {
				delete(m.items, key)
			}
		}
		m.mu.Unlock()
	}
}
