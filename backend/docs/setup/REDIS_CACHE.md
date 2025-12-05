# Redis Cache Implementation Guide

## 📋 Tổng Quan

Dự án đã được tích hợp Redis cache để cải thiện performance và giảm tải cho database. Cache được implement theo **Clean Architecture** với abstraction layer cho phép dễ dàng swap implementation.

## 🏗️ Kiến Trúc

### Cấu Trúc

```
internal/infrastructure/cache/
├── cache.go        # Cache interface (abstraction)
├── redis.go        # Redis implementation
├── memory.go       # In-memory implementation (dev/testing)
└── factory.go      # Factory để tạo cache instance
```

### Dependency Flow

```
Domain Layer (Service)
    ↓ uses interface
Infrastructure Layer (Cache)
    ↓ implements
Redis/Memory Cache
```

## 🔌 Cache Interface

```go
type Cache interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
    SetNX(ctx context.Context, key string, value []byte, ttl time.Duration) (bool, error)
    Close() error
}
```

## 🚀 Cài Đặt và Cấu Hình

### 1. Cài Đặt Redis

#### Local Development (Docker)

```bash
docker run -d -p 6379:6379 --name redis redis:latest
```

#### Hoặc dùng Homebrew (macOS)

```bash
brew install redis
brew services start redis
```

### 2. Cấu Hình Environment Variables

Thêm vào file `.env`:

```env
# Cache Configuration
CACHE_TYPE=redis              # "redis" or "memory" (default: "memory")

# Redis Configuration (required if CACHE_TYPE=redis)
REDIS_ADDR=localhost:6379     # Redis server address
REDIS_PASSWORD=               # Redis password (optional)
REDIS_DB=0                    # Redis database number (default: 0)
```

**Lưu ý:**

- `CACHE_TYPE=memory`: Dùng in-memory cache (phù hợp cho dev/testing, không cần Redis)
- `CACHE_TYPE=redis`: Dùng Redis cache (phù hợp cho production)

## 💻 Sử Dụng

### 1. Tạo Cache Instance

#### Cách 1: Dùng Factory (Khuyên dùng)

```go
import "github.com/example/foodie/backend/internal/infrastructure/cache"

// Tự động detect từ environment variables
cacheInstance, err := cache.NewCache()
if err != nil {
    log.Fatal("Failed to create cache:", err)
}
defer cacheInstance.Close()
```

#### Cách 2: Tạo trực tiếp

```go
import "github.com/example/foodie/backend/internal/infrastructure/cache"

// Redis
redisCache, err := cache.NewRedisCache("localhost:6379", "", 0)
if err != nil {
    log.Fatal(err)
}

// Hoặc Memory (for testing)
memoryCache := cache.NewMemoryCache()
```

### 2. Inject Cache vào Service

```go
import (
    "github.com/example/foodie/backend/internal/domain/product"
    "github.com/example/foodie/backend/internal/infrastructure/cache"
)

// Tạo cache
cacheInstance, _ := cache.NewCache()

// Tạo repository (SQL hoặc MongoDB)
repo := // ... create repository

// Tạo service với cache
productService := product.NewServiceWithCache(repo, cacheInstance)
```

### 3. Sử Dụng trong Service

Ví dụ với Product Service (đã được implement):

```go
// GetProduct tự động check cache trước khi query database
product, err := productService.GetProduct(ctx, "product-123")

// Invalidate cache khi update/delete
err = productService.InvalidateProductCache(ctx, "product-123")
```

## 📝 Example: Thêm Cache vào Service Mới

### Bước 1: Cập nhật Service Constructor

```go
// internal/domain/order/service_impl.go
type serviceImpl struct {
    repo  Repository
    cache cache.Cache // Optional
}

// Constructor không có cache
func NewService(repo Repository) Service {
    return &serviceImpl{repo: repo}
}

// Constructor có cache
func NewServiceWithCache(repo Repository, cache cache.Cache) Service {
    return &serviceImpl{
        repo:  repo,
        cache: cache,
    }
}
```

### Bước 2: Implement Cache Logic

```go
func (s *serviceImpl) GetOrder(ctx context.Context, orderID string) (*Order, error) {
    // Check cache first
    if s.cache != nil {
        cacheKey := "order:" + orderID

        var cached Order
        if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil && cached.ID != "" {
            return &cached, nil
        }

        // Cache miss - fetch from DB
        order, err := s.repo.FindByID(ctx, orderID)
        if err != nil {
            return nil, err
        }

        // Store in cache
        if order != nil {
            _ = s.cache.SetJSON(ctx, cacheKey, order, 10*time.Minute)
        }

        return order, nil
    }

    // No cache - direct query
    return s.repo.FindByID(ctx, orderID)
}

// Invalidate cache when order is updated
func (s *serviceImpl) UpdateOrder(ctx context.Context, order *Order) error {
    // Update in database
    err := s.repo.Update(ctx, order)
    if err != nil {
        return err
    }

    // Invalidate cache
    if s.cache != nil {
        _ = s.cache.Delete(ctx, "order:"+order.ID)
    }

    return nil
}
```

### Bước 3: Cập nhật main.go

```go
// cmd/server/main.go
func main() {
    // ... existing code ...

    // Create cache
    cacheInstance, err := cache.NewCache()
    if err != nil {
        log.Fatal("Failed to create cache:", err)
    }
    defer cacheInstance.Close()

    // Create services with cache
    orderService := order.NewServiceWithCache(orderRepo, cacheInstance)
    productService := product.NewServiceWithCache(productRepo, cacheInstance)

    // ... rest of code ...
}
```

## 🎯 Best Practices

### 1. Cache Key Naming Convention

Sử dụng format nhất quán: `{entity}:{id}`

```go
"product:123"
"order:456"
"user:789"
"restaurant:abc"
```

### 2. TTL (Time To Live) Recommendations

- **Products**: 5-15 minutes (ít thay đổi)
- **Orders**: 10-30 minutes (thay đổi theo status)
- **User sessions**: 1-24 hours
- **Static data**: 1 hour - 1 day

```go
// Short TTL for frequently changing data
s.cache.SetJSON(ctx, key, data, 5*time.Minute)

// Long TTL for static data
s.cache.SetJSON(ctx, key, data, 1*time.Hour)
```

### 3. Cache Invalidation Strategy

#### Immediate Invalidation (Khuyên dùng)

```go
// Khi update/delete, invalidate ngay
func (s *service) UpdateProduct(ctx context.Context, product *Product) error {
    err := s.repo.Update(ctx, product)
    if err != nil {
        return err
    }

    // Invalidate cache
    _ = s.cache.Delete(ctx, "product:"+product.ID)

    // Also invalidate related caches
    _ = s.cache.Delete(ctx, "restaurant:"+product.RestaurantID+":products")

    return nil
}
```

#### Pattern-based Invalidation

```go
// Invalidate all products of a restaurant
func invalidateRestaurantProducts(ctx context.Context, cache cache.Cache, restaurantID string) {
    // Note: Redis supports pattern deletion, but our interface doesn't yet
    // For now, maintain a list of keys or use Redis directly
    pattern := "product:*:restaurant:" + restaurantID
}
```

### 4. Error Handling

Cache failures should not break the application:

```go
func (s *serviceImpl) GetProduct(ctx context.Context, id string) (*Product, error) {
    if s.cache != nil {
        var cached Product
        if err := s.cache.GetJSON(ctx, "product:"+id, &cached); err == nil && cached.ID != "" {
            return &cached, nil
        }
        // Continue to DB if cache miss or error (don't fail)
    }

    // Always fallback to database
    return s.repo.FindByID(ctx, id)
}
```

### 5. Distributed Locks (SetNX)

Sử dụng `SetNX` để tránh cache stampede (thundering herd):

```go
func (s *serviceImpl) GetProduct(ctx context.Context, id string) (*Product, error) {
    cacheKey := "product:" + id
    lockKey := "lock:" + cacheKey

    // Try cache first
    var product Product
    if s.cache != nil {
        if err := s.cache.GetJSON(ctx, cacheKey, &product); err == nil && product.ID != "" {
            return &product, nil
        }

        // Try to acquire lock
        lockAcquired, _ := s.cache.SetNX(ctx, lockKey, []byte("1"), 10*time.Second)
        if !lockAcquired {
            // Another goroutine is fetching, wait a bit and retry cache
            time.Sleep(100 * time.Millisecond)
            if err := s.cache.GetJSON(ctx, cacheKey, &product); err == nil && product.ID != "" {
                return &product, nil
            }
        }
        defer s.cache.Delete(ctx, lockKey)
    }

    // Fetch from DB
    product, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Store in cache
    if s.cache != nil {
        _ = s.cache.SetJSON(ctx, cacheKey, product, 5*time.Minute)
    }

    return product, nil
}
```

## 🧪 Testing

### Unit Test với Memory Cache

```go
func TestProductService_GetProduct_WithCache(t *testing.T) {
    // Setup
    mockRepo := &MockRepository{}
    memoryCache := cache.NewMemoryCache()
    service := product.NewServiceWithCache(mockRepo, memoryCache)

    // Test cache hit
    // ...

    // Test cache miss
    // ...
}
```

### Integration Test với Redis

```go
func TestProductService_GetProduct_RedisIntegration(t *testing.T) {
    // Skip if Redis not available
    redisCache, err := cache.NewRedisCache("localhost:6379", "", 0)
    if err != nil {
        t.Skip("Redis not available")
    }
    defer redisCache.Close()

    // Run tests
    // ...
}
```

## 🔧 Advanced Usage

### Custom JSON Marshaling

Nếu cần custom JSON format:

```go
// Use Get/Set with manual marshaling
data, _ := json.Marshal(product)
_ = cache.Set(ctx, key, data, ttl)

// Or use GetJSON/SetJSON helpers
_ = cache.SetJSON(ctx, key, product, ttl)
```

### Cache Warming

Pre-populate cache on startup:

```go
func warmupCache(ctx context.Context, cache cache.Cache, repo Repository) {
    products, _ := repo.FindAll(ctx)
    for _, product := range products {
        key := "product:" + product.ID
        _ = cache.SetJSON(ctx, key, product, 1*time.Hour)
    }
}
```

## 📊 Monitoring

### Cache Metrics (Future Enhancement)

Có thể thêm metrics để monitor cache performance:

- Cache hit rate
- Cache miss rate
- Average response time
- Memory usage (Redis)

## 🐛 Troubleshooting

### Redis Connection Issues

```bash
# Check if Redis is running
redis-cli ping
# Should return: PONG

# Check connection from Go
redis-cli -h localhost -p 6379 ping
```

### Cache Not Working

1. Check environment variables:

   ```bash
   echo $CACHE_TYPE
   echo $REDIS_ADDR
   ```

2. Verify cache is being injected:

   ```go
   if service.cache == nil {
       log.Println("Cache not initialized")
   }
   ```

3. Check Redis logs:
   ```bash
   docker logs redis
   ```

## 📚 Tài Liệu Tham Khảo

- [Redis Go Client](https://github.com/redis/go-redis)
- [Redis Best Practices](https://redis.io/docs/manual/patterns/)
- [Cache Patterns](https://docs.aws.amazon.com/AmazonElastiCache/latest/mem-ug/best-practices.html)
