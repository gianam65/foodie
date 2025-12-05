package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements Cache interface using Redis.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache instance.
func NewRedisCache(addr, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
	}, nil
}

// NewRedisCacheFromClient creates a Redis cache from an existing Redis client.
// Useful when you already have a Redis client configured.
func NewRedisCacheFromClient(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

// Get retrieves a value by key.
func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil // Key doesn't exist, return nil without error
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Set stores a value with expiration time.
func (r *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Delete removes a key from cache.
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in cache.
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SetNX sets a key only if it doesn't exist (useful for distributed locks).
func (r *RedisCache) SetNX(ctx context.Context, key string, value []byte, ttl time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, ttl).Result()
}

// Close closes the Redis connection.
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// Helper functions for common operations

// GetJSON retrieves and unmarshals a JSON value.
func (r *RedisCache) GetJSON(ctx context.Context, key string, v interface{}) error {
	data, err := r.Get(ctx, key)
	if err != nil {
		return err
	}
	if data == nil {
		return nil // Key doesn't exist
	}
	return json.Unmarshal(data, v)
}

// SetJSON marshals and stores a JSON value.
func (r *RedisCache) SetJSON(ctx context.Context, key string, v interface{}, ttl time.Duration) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return r.Set(ctx, key, data, ttl)
}
