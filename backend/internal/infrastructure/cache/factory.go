package cache

import (
	"fmt"
	"os"
	"strconv"

	"foodie/backend/pkg/config"
)

// NewCache creates a cache instance based on configuration.
// It reads from environment variables:
//   - CACHE_TYPE: "redis" or "memory" (default: "memory")
//   - REDIS_ADDR: Redis address (default: "localhost:6379")
//   - REDIS_PASSWORD: Redis password (optional)
//   - REDIS_DB: Redis database number (default: 0)
func NewCache() (Cache, error) {
	cacheType := config.Get("CACHE_TYPE", "memory")

	switch cacheType {
	case "redis":
		return newRedisCacheFromConfig()
	case "memory":
		return NewMemoryCache(), nil
	default:
		return nil, fmt.Errorf("unsupported cache type: %s", cacheType)
	}
}

// newRedisCacheFromConfig creates a Redis cache from environment variables.
func newRedisCacheFromConfig() (*RedisCache, error) {
	addr := config.Get("REDIS_ADDR", "localhost:6379")
	password := config.Get("REDIS_PASSWORD", "")

	dbStr := os.Getenv("REDIS_DB")
	db := 0
	if dbStr != "" {
		var err error
		db, err = strconv.Atoi(dbStr)
		if err != nil {
			return nil, fmt.Errorf("invalid REDIS_DB value: %s", dbStr)
		}
	}

	return NewRedisCache(addr, password, db)
}
