# 🚀 Quick Start Guide

## Chạy Backend Server

### Sử dụng Makefile (Khuyên dùng)

```bash
# Chạy bình thường
make run

# Chạy với hot reload
make dev

# Xem tất cả commands
make help
```

---

### Chạy trực tiếp với Go

```bash
# Chạy server
go run ./cmd/server

# Hoặc build trước rồi chạy
go build -o bin/server ./cmd/server
./bin/server
```

---

## ⚙️ Cấu hình

### Tạo file .env

Copy file `.env.example` và tùy chỉnh:

```bash
cp .env.example .env
```

File `.env.example` bao gồm:

- **Server Configuration**: PORT, HOST
- **Database Configuration**: SQL (PostgreSQL/MySQL)
- **Redis Cache Configuration**: Cache type và Redis settings
- **Message Broker Configuration**: RabbitMQ settings

**Ví dụ .env cho development:**

```env
# Server
SERVER_PORT=8080

# Database (PostgreSQL/MySQL)
SQL_DSN=postgres://user:password@localhost:5432/foodie?sslmode=disable

# Cache (dùng memory cache - không cần Redis)
CACHE_TYPE=memory

# Message Broker (dùng memory - không cần RabbitMQ)
MESSAGE_BROKER_TYPE=memory
```

**Ví dụ .env cho production với Redis và RabbitMQ:**

```env
# Server
SERVER_PORT=8080

# Database (PostgreSQL/MySQL)
SQL_DSN=postgres://user:password@localhost:5432/foodie?sslmode=disable

# Cache (dùng Redis)
CACHE_TYPE=redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=your_redis_password
REDIS_DB=0
```

**Lưu ý:**

- Nếu không có file `.env`, server vẫn chạy được với các giá trị mặc định
- `CACHE_TYPE=memory` không cần Redis (phù hợp cho development)
- `CACHE_TYPE=redis` cần Redis server đang chạy

---

## 🧪 Kiểm tra Server

Sau khi chạy server, mở browser và truy cập:

- **Health Check:** http://localhost:8080/health
- **API Base URL:** http://localhost:8080/api/v1

---

## 🛠️ Các Commands Khác

```bash
# Chạy migrations
make migrate-up

# Rollback migration
make migrate-down

# Build binary
make build

# Run tests
make test

# Format code
make fmt

# Clean build artifacts
make clean
```

---

## ❓ Troubleshooting

### Lỗi: "Go is not installed"

Cài đặt Go từ https://golang.org/dl/

### Lỗi: "air is not installed"

Cài đặt air:

```bash
go install github.com/air-verse/air@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

---

## 📝 Ghi chú

- Server mặc định chạy trên port **8080**
- Để dừng server, nhấn `Ctrl+C`
- `make dev` sử dụng hot reload (air), tự động restart khi bạn sửa code
- Xem [Services Management](./SERVICES_MANAGEMENT.md) để chạy tất cả services cùng lúc
