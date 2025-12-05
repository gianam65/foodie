# 📊 Structured Logging Setup for Grafana/Loki

## 🎯 Overview

Project sử dụng **structured logging** với **JSON format** để dễ dàng tích hợp với Grafana/Loki hoặc các log aggregation systems khác (ELK, Splunk, etc.).

## 🏗️ Architecture

### Logger Implementation

- **Library**: `go.uber.org/zap` - High-performance structured logger
- **Format**: JSON (mặc định) hoặc Text
- **Output**: stdout (mặc định) - được capture bởi container log drivers

### Features

- ✅ **Structured JSON logs** - Dễ parse và query
- ✅ **Correlation ID** - Trace requests across services
- ✅ **Log levels** - debug, info, warn, error
- ✅ **Request tracing** - HTTP method, path, status, duration
- ✅ **Grafana/Loki compatible** - Ready for log aggregation

---

## ⚙️ Configuration

### Environment Variables

Thêm vào file `.env`:

```env
# Logging Configuration
LOG_LEVEL=info          # debug, info, warn, error (default: info)
LOG_FORMAT=json         # json, text (default: json)
LOG_OUTPUT=stdout       # stdout, stderr, or file path (default: stdout)
```

### Log Levels

| Level | Usage |
|-------|-------|
| `debug` | Development - verbose logging |
| `info` | Production default - general information |
| `warn` | Warnings - potential issues |
| `error` | Errors only - critical issues |

---

## 📝 Log Format

### JSON Format (Default)

```json
{
  "level": "info",
  "timestamp": "2024-12-05T10:30:00Z",
  "caller": "middleware/logging.go:38",
  "message": "http_request",
  "method": "GET",
  "path": "/api/v1/orders",
  "query": "user_id=123&page=1",
  "status": 200,
  "duration": "45.2ms",
  "correlation_id": "550e8400-e29b-41d4-a716-446655440000",
  "remote_addr": "192.168.1.1:54321",
  "user_agent": "Mozilla/5.0..."
}
```

### Text Format (Development)

```
2024-12-05T10:30:00Z	INFO	middleware/logging.go:38	http_request	{"method": "GET", "path": "/api/v1/orders", "status": 200, ...}
```

---

## 🔍 Correlation ID

Mỗi HTTP request tự động được gán một **correlation ID** để trace requests:

- **Header**: `X-Correlation-ID`
- **Auto-generated**: Nếu request không có header này
- **Propagation**: ID được include trong tất cả logs liên quan đến request đó

### Usage Example

```bash
# Request với correlation ID
curl -H "X-Correlation-ID: my-custom-id" http://localhost:8080/api/v1/orders

# Response sẽ có header
X-Correlation-ID: my-custom-id
```

Tất cả logs từ request này sẽ có `correlation_id: "my-custom-id"`, giúp query dễ dàng trong Grafana.

---

## 📊 Grafana/Loki Integration

### 1. Log Collection

Logs được output ra `stdout` (JSON format), được capture bởi:

- **Docker**: `docker logs` hoặc log drivers
- **Kubernetes**: Container logs → log collector
- **Systemd**: Journald hoặc file logging

### 2. Promtail Configuration

Promtail (Loki log shipper) có thể scrape logs:

```yaml
# promtail-config.yaml
server:
  http_listen_port: 9080

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://loki:3100/loki/api/v1/push

scrape_configs:
  - job_name: foodie-backend
    static_configs:
      - targets:
          - localhost
        labels:
          job: foodie-backend
          __path__: /var/log/foodie/*.log
    
  # Docker logs
  - job_name: docker
    docker_sd_configs:
      - host: unix:///var/run/docker.sock
        refresh_interval: 5s
    relabel_configs:
      - source_labels: ['__meta_docker_container_name']
        regex: '/foodie.*'
        action: keep
```

### 3. Loki Query Examples

#### Find all errors

```logql
{job="foodie-backend"} |= "error"
```

#### Find requests by correlation ID

```logql
{job="foodie-backend"} |= "correlation_id" | json | correlation_id="550e8400-e29b-41d4-a716-446655440000"
```

#### Find slow requests (>1s)

```logql
{job="foodie-backend"} | json | duration > 1s
```

#### HTTP requests by status code

```logql
{job="foodie-backend"} | json | message="http_request" | status >= 400
```

#### Group by endpoint

```logql
sum by (path) (count_over_time({job="foodie-backend"} | json | message="http_request" [5m]))
```

---

## 🔧 Usage in Code

### Initialize Logger

```go
import "foodie/backend/pkg/logger"

// Default logger (JSON, info level, stdout)
log := logger.NewDefault()

// Custom configuration
logConfig := logger.Config{
    Level:  "debug",
    Format: "json",
    Output: "stdout",
}
log, err := logger.New(logConfig)
if err != nil {
    panic(err)
}
defer log.Sync()
```

### Log Messages

```go
import "go.uber.org/zap"

// Info log with fields
log.Info("order_created",
    zap.String("order_id", orderID),
    zap.String("user_id", userID),
    zap.Float64("total", total),
)

// Error log
log.Error("failed_to_create_order",
    zap.Error(err),
    zap.String("user_id", userID),
)

// With correlation ID from request
correlationID := middleware.GetCorrelationID(r)
log.WithRequestID(correlationID).Info("processing_order",
    zap.String("order_id", orderID),
)
```

### Available Helper Functions

```go
logger.String("key", "value")
logger.Int("key", 123)
logger.Int64("key", 123456)
logger.Float64("key", 123.45)
logger.Bool("key", true)
logger.Error(err)
logger.Duration("key", time.Since(start))
```

---

## 📈 Grafana Dashboard Queries

### Request Rate

```logql
rate({job="foodie-backend"} | json | message="http_request" [5m])
```

### Error Rate

```logql
rate({job="foodie-backend"} | json | status >= 400 [5m])
```

### Average Response Time

```logql
avg_over_time({job="foodie-backend"} | json | message="http_request" | duration [5m])
```

### Top Endpoints

```logql
topk(10, sum by (path) (count_over_time({job="foodie-backend"} | json | message="http_request" [5m])))
```

### Requests by Method

```logql
sum by (method) (count_over_time({job="foodie-backend"} | json | message="http_request" [5m]))
```

---

## 🐳 Docker Logging

### Docker Compose

```yaml
services:
  backend:
    image: foodie-backend:latest
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
        labels: "app=foodie-backend"
```

### Kubernetes

Logs tự động được capture vào container logs, có thể forward đến Loki:

```yaml
# Deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: foodie-backend
spec:
  template:
    metadata:
      labels:
        app: foodie-backend
    spec:
      containers:
      - name: backend
        image: foodie-backend:latest
        env:
        - name: LOG_LEVEL
          value: "info"
        - name: LOG_FORMAT
          value: "json"
```

---

## ✅ Best Practices

1. **Always include correlation_id** trong logs liên quan đến request
2. **Use structured fields** thay vì string concatenation
3. **Set appropriate log levels** - không log quá nhiều ở production
4. **Include context** - user_id, order_id, etc. trong logs
5. **Log errors with stack traces** - đã được handle bởi recovery middleware
6. **Use meaningful message keys** - `order_created`, `payment_failed`, etc.

---

## 🔗 References

- [Uber Zap Documentation](https://github.com/uber-go/zap)
- [Grafana Loki Documentation](https://grafana.com/docs/loki/latest/)
- [LogQL Query Language](https://grafana.com/docs/loki/latest/logql/)

