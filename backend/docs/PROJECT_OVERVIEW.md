# 📋 Tổng Quan Dự Án Foodie Backend

## 🎯 Giới Thiệu

**Foodie Backend** là hệ thống backend cho nền tảng giao đồ ăn được xây dựng bằng **Golang**, tuân theo kiến trúc **Modular Monolith** với các nguyên tắc **Clean Architecture**.

### Đặc Điểm Chính

- ✅ **Modular Monolith**: Tất cả modules trong một codebase, dễ dàng tách thành microservices sau này
- ✅ **Clean Architecture**: Tách biệt rõ ràng giữa domain, application và infrastructure layers
- ✅ **Hexagonal Architecture**: Adapters cho các dependencies bên ngoài (database, messaging, APIs)
- ✅ **Repository Pattern**: Truy cập dữ liệu trừu tượng với SQL implementations
- ✅ **SQL Database Support**: Hỗ trợ PostgreSQL và MySQL

---

## 📁 Cấu Trúc Thư Mục và Chức Năng

### 🔵 `/api` - API Definitions

**Chức năng:**

- Định nghĩa các API contracts (OpenAPI/Swagger, gRPC proto)
- Tài liệu API specifications
- Chứa schema definitions cho request/response

**Nhiệm vụ:**

- `openapi.yaml`: Định nghĩa REST API endpoints, schemas, responses
- `README.md`: Hướng dẫn về API documentation

**Ví dụ sử dụng:**

```yaml
# Định nghĩa endpoint /orders
/orders:
  get:
    summary: List orders
    responses:
      "200":
        description: List of orders
```

---

### 🔵 `/cmd` - Application Entry Points

**Chức năng:**

- Chứa các entry points của ứng dụng (main functions)
- Mỗi thư mục con là một ứng dụng độc lập

**Cấu trúc:**

#### `/cmd/server/`

- **Chức năng**: Entry point chính của ứng dụng
- **Nhiệm vụ**:
  - Khởi tạo HTTP server
  - Setup routing, middleware
  - Dependency injection (repositories, services, controllers)
  - Graceful shutdown handling

#### `/cmd/migrate/`

- **Chức năng**: CLI tool để chạy database migrations
- **Nhiệm vụ**:
  - Apply migrations (up/down)
  - Quản lý version database schema

#### `/cmd/examples/`

- **Chức năng**: Các ví dụ demo
- **Nhiệm vụ**: Các ví dụ sử dụng API và patterns

---

### 🔵 `/configs` - Configuration Templates

**Chức năng:**

- Chứa các file config mẫu (JSON, YAML)
- Template để tạo `.env` file

**Nhiệm vụ:**

- Cung cấp examples cho developers
- Document các config options có sẵn

---

### 🔵 `/internal` - Internal Application Code

> ⚠️ **Quan trọng**: Code trong `/internal` không thể được import bởi các project khác (Go's internal package rule).

#### `/internal/interfaces/` - **Interface/Adapter Layer (Transport)**

**Chức năng:** Lớp giao tiếp với bên ngoài (HTTP, gRPC)

##### `/internal/interfaces/http/`

- **Chức năng**: HTTP handlers, routing, middleware
- **Nhiệm vụ**:
  - `controller/`: HTTP handlers cho các endpoints
    - `health_controller.go`: Health check endpoint
    - `order_controller.go`: Order CRUD operations
    - `product_controller.go`: Product listing
    - `response.go`: Standard response helpers
  - `dto/`: Data Transfer Objects (request/response mapping)
  - `middleware/`: HTTP middleware (auth, logging, CORS, recovery)
  - `router.go`: Route definitions và setup
  - `router_group.go`: Route grouping utilities

**Flow hoạt động:**

```
Request → Router → Middleware → Controller → UseCase → Repository → Database
```

##### `/internal/interfaces/grpc/`

- **Chức năng**: gRPC service definitions
- **Nhiệm vụ**:
  - `service.proto`: Protocol Buffer definitions
  - Chuẩn bị cho microservices communication

---

#### `/internal/domain/` - **Domain Layer (Business Logic)**

**Chức năng:** Core business logic, không phụ thuộc vào infrastructure

**Nguyên tắc:**

- ✅ Không import packages từ `infrastructure` hoặc `app`
- ✅ Chỉ chứa business entities, interfaces, domain services

**Cấu trúc mỗi module:**

##### `/internal/domain/order/`

```
order/
├── entity.go          # Domain entity (Order struct)
├── repository.go      # Repository interface (không có implementation)
├── service.go         # Service interface
└── service_impl.go    # Service implementation (business logic)
```

**Nhiệm vụ:**

- `entity.go`: Định nghĩa Order struct với business rules
- `repository.go`: Interface cho data access (abstraction)
- `service.go`: Interface cho business operations
- `service_impl.go`: Business logic (tính tổng, validation, workflows)

##### Các module khác (tương tự):

- `/internal/domain/product/`: Product/menu management
- `/internal/domain/user/`: User authentication (sẽ implement)
- `/internal/domain/restaurant/`: Restaurant info (sẽ implement)
- `/internal/domain/payment/`: Payment processing (sẽ implement)
- `/internal/domain/delivery/`: Delivery tracking (sẽ implement)
- `/internal/domain/notification/`: Notifications (sẽ implement)

---

#### `/internal/infrastructure/` - **Infrastructure Layer (Adapters)**

**Chức năng:** Implementations cho các external dependencies

##### `/internal/infrastructure/database/` - **Database Layer**

**Chức năng:** Repository implementations và orchestration

**Cấu trúc:**

```
database/
├── order/
│   └── repository.go    # Order repository implementation
├── product/
│   └── repository.go    # Product repository implementation
└── repositories.go      # Repositories aggregate + factory
```

**Nguyên tắc:**

- Mỗi domain module có folder riêng trong `database/`
- Mỗi folder chứa `repository.go` implement interface từ domain layer
- `repositories.go` chỉ chứa aggregate struct và factory function
- Tách biệt rõ ràng giữa domain logic (domain layer) và data access (infrastructure layer)
- Sử dụng PostgreSQL/MySQL driver
- SQL queries, prepared statements

##### `/internal/infrastructure/database/repositories.go`

- **Chức năng**: Aggregate tất cả repositories và factory function
- **Logic**: Quản lý và khởi tạo tất cả repositories cho application

```go
// Repository implementations
type OrderRepository struct { ... }
type ProductRepository struct { ... }

// Aggregate all repositories
type Repositories struct {
    Order   order.Repository
    Product product.Repository
    // Future: User, Restaurant, Payment...
}

func NewRepositories(sqlDB *sql.DB) (*Repositories, error) {
    return &Repositories{
        Order:   NewOrderRepository(sqlDB),
        Product: NewProductRepository(sqlDB),
    }, nil
}
```

##### `/internal/infrastructure/messaging/`

- **Chức năng**: Event/message brokers (RabbitMQ, Kafka, etc.)
- `event_publisher.go`: Publish events cho domain events
- **Use case**: Order created → publish event → trigger notifications

##### `/internal/infrastructure/external/`

- **Chức năng**: External service clients
- `payment_gateway.go`: Integration với payment providers (Stripe, PayPal)
- `maps_service.go`: Integration với mapping APIs (Google Maps, Mapbox)

---

### 🔵 `/pkg` - Shared Packages

> ✅ **Lưu ý**: Code trong `/pkg` có thể được import bởi các project khác.

**Chức năng:** Utilities và shared code

##### `/pkg/config/`

- `config.go`: Load configuration từ `.env` file
- Helper functions: `Get()`, `MustLoad()`

##### `/pkg/dbtypes/`

- `types.go`: Type definitions cho database (kept for backward compatibility)

##### `/pkg/logger/`

- `logger.go`: Logging utilities (structured logging)

##### `/pkg/migrate/`

- `migrate.go`: Database migration utilities

---

### 🔵 `/migrations` - Database Migrations

**Chức năng:**

- SQL migration files (up/down)
- Version control cho database schema

**Format:**

```
000001_create_orders_table.up.sql    # Create table
000001_create_orders_table.down.sql  # Drop table
```

---

### 🔵 `/tmp` - Temporary Files

**Chức năng:**

- Build artifacts
- Temporary files (thường được ignore trong git)

---

## 🏗️ Kiến Trúc và Dependency Flow

### Dependency Rule (Quy Tắc Phụ Thuộc)

```
Domain ← Application ← Infrastructure
```

**Giải thích:**

1. **Domain** không phụ thuộc gì cả (zero dependencies)
2. **Application** chỉ phụ thuộc vào Domain
3. **Infrastructure** phụ thuộc vào Domain và Application

**Ví dụ:**

```go
// ✅ ĐÚNG: Application import Domain
// internal/interfaces/http/controller/order_controller.go
import "foodie/backend/internal/domain/order"

// ✅ ĐÚNG: Infrastructure import Domain
// internal/infrastructure/database/repositories.go
import "github.com/example/foodie/backend/internal/domain/order"

// ❌ SAI: Domain import Infrastructure (KHÔNG ĐƯỢC)
// internal/domain/order/entity.go
import "github.com/example/foodie/backend/internal/infrastructure/database"
```

---

## 🔄 Database Switching Mechanism

### Cách Hoạt Động

1. **Domain Layer** định nghĩa interface:

```go
// internal/domain/order/repository.go
type Repository interface {
    Save(ctx context.Context, order *Order) error
    FindByID(ctx context.Context, id string) (*Order, error)
}
```

2. **Infrastructure Layer** cung cấp SQL implementations:

   - `internal/infrastructure/database/order/repository.go`: Order repository implementation
   - `internal/infrastructure/database/product/repository.go`: Product repository implementation
   - `internal/infrastructure/database/repositories.go`: Repositories aggregate và factory

3. **Repositories** aggregate tất cả repositories:

```go
repos, err := database.NewRepositories(sqlDB)
// Tất cả modules dùng SQL database
```

---

## 🚀 Khả Năng Mở Rộng (Extensibility)

### 1. Thêm Module Mới (Domain)

**Ví dụ: Thêm User Module**

#### Bước 1: Tạo Domain Layer

```go
// internal/domain/user/entity.go
package user

type User struct {
    ID    string
    Email string
    Name  string
}

// internal/domain/user/repository.go
type Repository interface {
    Save(ctx context.Context, user *User) error
    FindByEmail(ctx context.Context, email string) (*User, error)
}

// internal/domain/user/service.go
type Service interface {
    Register(ctx context.Context, email, password string) (*User, error)
}

// internal/domain/user/service_impl.go
type service struct {
    repo Repository
}

func NewService(repo Repository) Service {
    return &service{repo: repo}
}
```

#### Bước 2: Tạo Infrastructure Implementation

Tạo folder và file mới: `internal/infrastructure/database/user/repository.go`

```go
// internal/infrastructure/database/user/repository.go
package user

import (
    "context"
    "database/sql"
    "github.com/example/foodie/backend/internal/domain/user"
)

type Repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) user.Repository {
    return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, u *user.User) error {
    // Implementation...
}
```

#### Bước 3: Thêm vào Repositories

```go
// internal/infrastructure/database/repositories.go
import (
    userrepo "github.com/example/foodie/backend/internal/infrastructure/database/user"
    // ...
)

type Repositories struct {
    Order   order.Repository
    Product product.Repository
    User    user.Repository  // ← Thêm vào đây
}

func NewRepositories(sqlDB *sql.DB) (*Repositories, error) {
    if sqlDB == nil {
        return nil, fmt.Errorf("sql database connection is nil")
    }
    return &Repositories{
        Order:   orderrepo.NewRepository(sqlDB),
        Product: productrepo.NewRepository(sqlDB),
        User:    userrepo.NewRepository(sqlDB),  // ← Thêm vào đây
    }, nil
}
```

#### Bước 4: Tạo HTTP Layer

```go
// internal/interfaces/http/controller/user_controller.go
type UserController struct {
    service user.Service
}

// internal/interfaces/http/router/router.go
func (r *Router) setupPrivateRoutes(private *RouteGroup) {
    private.POST("/users/register", r.userController.Register)
}
```

**Kết quả:** Module mới tự động hỗ trợ cả SQL và MongoDB! 🎉

---

### 2. Thêm External Service Integration

**Ví dụ: Thêm SMS Service**

```go
// internal/infrastructure/external/sms_service.go
type SMSService interface {
    SendSMS(ctx context.Context, phone, message string) error
}

type twilioSMSService struct {
    client *twilio.Client
}

// internal/domain/notification/service.go
type Service interface {
    SendOrderConfirmation(ctx context.Context, userID, orderID string) error
}
```

---

### 3. Thêm Message Broker

**Ví dụ: Event-Driven Notifications**

```go
// internal/infrastructure/messaging/event_publisher.go
func (p *EventPublisher) PublishOrderCreated(ctx context.Context, order *order.Order) error {
    event := OrderCreatedEvent{
        OrderID: order.ID,
        UserID:  order.UserID,
    }
    return p.publisher.Publish("orders.created", event)
}

// internal/domain/order/service_impl.go
func (s *service) CreateOrder(...) (*Order, error) {
    // ... create order logic ...
    eventPublisher.PublishOrderCreated(ctx, newOrder)  // Publish event
    return newOrder, nil
}
```

---

### 4. Thêm gRPC Support

```go
// internal/interfaces/grpc/service.proto
service OrderService {
    rpc CreateOrder(CreateOrderRequest) returns (Order);
}

// internal/interfaces/grpc/order_handler.go
func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
    order, err := s.orderService.CreateOrder(...)
    // Convert to protobuf
}
```

---

## 📈 Khả Năng Scale Up

### 1. Modular Monolith → Microservices

**Hiện tại:** Tất cả modules trong một service

**Khi cần scale:**

#### Bước 1: Tách theo Domain Boundaries

```
foodie-backend/
├── order-service/     # Chỉ Order module
├── product-service/   # Chỉ Product module
├── user-service/      # Chỉ User module
└── api-gateway/       # Routing layer
```

#### Bước 2: Mỗi Service độc lập

- Mỗi service có `cmd/server/main.go` riêng
- Mỗi service có database riêng
- Communication qua gRPC hoặc message queue

#### Bước 3: Deploy độc lập

```yaml
# docker-compose.yml
services:
  order-service:
    image: foodie/order-service:latest
  product-service:
    image: foodie/product-service:latest
```

**Lợi ích:**

- ✅ Scale từng service độc lập
- ✅ Deploy độc lập (không ảnh hưởng lẫn nhau)
- ✅ Technology stack riêng cho từng service (nếu cần)

---

### 2. Database Scaling Strategies

#### Option A: Read Replicas

```go
// internal/infrastructure/database/repositories.go
type Repositories struct {
    Order   order.Repository
    OrderRead order.Repository  // Read replica
}

// Routing: Write → Master, Read → Replica
```

#### Option B: Database Sharding

```go
// internal/infrastructure/database/sql/repositories.go
type orderRepository struct {
    shards []*sql.DB  // Multiple databases
}

func (r *orderRepository) getShard(orderID string) *sql.DB {
    // Hash orderID to determine shard
    return r.shards[hash(orderID) % len(r.shards)]
}
```

#### Option C: CQRS Pattern

```
Command Side (Write): PostgreSQL
Query Side (Read): MongoDB/Elasticsearch (optimized for reads)
```

---

### 3. Horizontal Scaling (Load Balancing)

#### Architecture:

```
                 Load Balancer (nginx/traefik)
                      /    |    \
              Service 1  Service 2  Service 3
                  |         |          |
              Database (Primary) + Read Replicas
```

#### Implementation:

```yaml
# Kubernetes deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: foodie-backend
spec:
  replicas: 5 # 5 instances
  template:
    spec:
      containers:
        - name: server
          image: foodie/backend:latest
```

---

### 4. Caching Layer

#### Redis Integration:

```go
// internal/infrastructure/cache/redis.go
type Cache interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}

// internal/domain/product/service_impl.go
func (s *service) GetProduct(ctx context.Context, id string) (*Product, error) {
    // Check cache first
    cached, _ := s.cache.Get(ctx, "product:"+id)
    if cached != nil {
        return unmarshal(cached), nil
    }

    // Cache miss: fetch from DB
    product, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Cache result
    s.cache.Set(ctx, "product:"+id, marshal(product), 5*time.Minute)
    return product, nil
}
```

---

### 5. Message Queue for Async Processing

#### RabbitMQ/Kafka Integration:

```go
// internal/infrastructure/messaging/rabbitmq.go
type MessageQueue interface {
    Publish(ctx context.Context, queue string, message []byte) error
    Consume(ctx context.Context, queue string, handler Handler) error
}

// Use case: Order processing
// 1. API nhận order request → return immediately
// 2. Publish "order.created" event
// 3. Worker process order async (calculate total, send email, etc.)
```

**Lợi ích:**

- ✅ API response nhanh hơn (không chờ processing)
- ✅ Retry failed operations
- ✅ Decouple services

---

### 6. API Gateway Pattern

#### Với API Gateway:

```
Client → API Gateway → [Order Service, Product Service, User Service]
```

**Chức năng Gateway:**

- Authentication/Authorization
- Rate limiting
- Request routing
- Aggregation (combine multiple services)
- Load balancing

---

## 🎯 Best Practices cho Scaling

### 1. Stateless Services

- ✅ Không lưu state trong memory
- ✅ Sử dụng shared cache (Redis) cho session
- ✅ Mọi instance có thể xử lý mọi request

### 2. Database Connection Pooling

```go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

### 3. Health Checks

```go
// internal/interfaces/http/controller/health_controller.go
func (c *HealthController) Check(w http.ResponseWriter, r *http.Request) {
    // Check database connection
    // Check external services
    // Return status
}
```

### 4. Monitoring & Observability

- Logging: Structured logs (JSON format)
- Metrics: Prometheus
- Tracing: OpenTelemetry
- APM: New Relic, Datadog

### 5. Graceful Shutdown

```go
// cmd/server/main.go
ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
// ... wait for in-flight requests to complete
server.Shutdown(shutdownCtx)
```

---

## 📊 Tóm Tắt

### Điểm Mạnh của Kiến Trúc

1. **Clean Architecture**: Dễ test, dễ maintain
2. **Modular**: Dễ thêm module mới
3. **Flexible Database**: Switch SQL/MongoDB dễ dàng
4. **Scalable**: Sẵn sàng tách thành microservices
5. **Extensible**: Dễ thêm external services, message queues, caching

### Roadmap Scaling

```
Phase 1: Monolith (Hiện tại)
    ↓
Phase 2: Monolith + Caching + Read Replicas
    ↓
Phase 3: Modular Monolith (tách deployment)
    ↓
Phase 4: Microservices (tách service riêng)
    ↓
Phase 5: Multi-region deployment
```

---

## 📚 Tài Liệu Tham Khảo

- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Modular Monolith](https://www.kamilgrzybek.com/blog/posts/modular-monolith-primer)

---

**Tác giả:** Foodie Backend Team  
**Cập nhật:** 2024
