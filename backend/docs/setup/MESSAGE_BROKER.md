# Message Broker Setup Guide

## 📋 Tổng Quan

Hướng dẫn setup và sử dụng message broker cho project. Project hỗ trợ RabbitMQ và in-memory (development). Có thể mở rộng sang Kafka trong tương lai.

---

## 📊 RabbitMQ vs Kafka: So Sánh

| Đặc Điểm              | RabbitMQ                                 | Apache Kafka                                    |
| --------------------- | ---------------------------------------- | ----------------------------------------------- |
| **Kiến Trúc**         | Message Queue (Push-based)               | Distributed Log/Event Stream (Pull-based)       |
| **Message Model**     | Queue-based, Point-to-point hoặc Pub/Sub | Topic-based, Pub/Sub với partitioning           |
| **Throughput**        | Medium (10K-100K messages/sec)           | Very High (1M+ messages/sec)                    |
| **Latency**           | Low (sub-millisecond)                    | Low-Medium (milliseconds)                       |
| **Message Retention** | Xóa sau khi consumer ACK (hoặc TTL)      | Giữ lại theo thời gian/bytes (retention period) |
| **Replay Messages**   | Không hỗ trợ                             | Có thể replay (giữ trong log)                   |
| **Complexity**        | Đơn giản hơn                             | Phức tạp hơn (cần quản lý offsets)              |
| **Use Cases**         | Task queues, RPC, routing                | Event streaming, log aggregation, analytics     |

### 🎯 Recommendation cho Project này: **RabbitMQ**

**Lý do:**

- ✅ Phù hợp với use cases (task queues, notifications)
- ✅ Đơn giản, dễ maintain
- ✅ Đủ mạnh cho scale hiện tại (100K-500K orders/day)
- ✅ Cost-effective
- ✅ Team có thể implement nhanh

**Chọn Kafka nếu:**

- Scale > 1M events/day
- Cần event sourcing
- Cần real-time analytics
- Có data engineering team

---

## 🏗️ Kiến Trúc

```
Application → Publisher → RabbitMQ Exchange → Queue → Consumer/Worker
```

### Components

- **Publisher**: Publish events vào RabbitMQ exchange
- **Exchange**: Route messages đến queues dựa trên routing keys
- **Queue**: Lưu trữ messages chờ consumer xử lý
- **Consumer**: Workers xử lý messages từ queues

---

## 🚀 Cài Đặt RabbitMQ

### Option 1: Docker (Khuyên dùng)

```bash
# Start với docker-compose (đã có sẵn)
make docker-up
# hoặc chỉ RabbitMQ
make rabbitmq-up

# Access management UI: http://localhost:15672
# Username: admin (hoặc từ .env)
# Password: admin123 (hoặc từ .env)
```

### Option 2: Standalone Docker

```bash
docker run -d \
  --name rabbitmq \
  -p 5672:5672 \
  -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=admin \
  -e RABBITMQ_DEFAULT_PASS=admin123 \
  rabbitmq:3-management-alpine
```

### Option 3: Homebrew (macOS)

```bash
brew install rabbitmq
brew services start rabbitmq
```

---

## ⚙️ Cấu Hình

### Environment Variables

Thêm vào `.env`:

```env
# Message Broker Type
MESSAGE_BROKER_TYPE=rabbitmq  # "rabbitmq" or "memory"

# RabbitMQ Configuration
RABBITMQ_URL=amqp://admin:admin123@localhost:5672/
RABBITMQ_EXCHANGE=foodie_events

# For docker-compose
RABBITMQ_USER=admin
RABBITMQ_PASSWORD=admin123
```

**Development mode** (không cần RabbitMQ):

```env
MESSAGE_BROKER_TYPE=memory
```

---

## 💻 Sử Dụng

### 1. Publish Events

```go
import "foodie/backend/internal/infrastructure/messaging"

// Tự động detect từ environment variables
publisher, err := messaging.NewPublisher()
if err != nil {
    log.Fatal("Failed to create publisher:", err)
}
defer publisher.Close()

// Publish event
event := messaging.Event{
    Type:        "order.created",
    AggregateID: "order-123",
    Payload: map[string]interface{}{
        "order_id": "order-123",
        "user_id":  "user-456",
        "total":    50.00,
    },
    Timestamp: time.Now().Unix(),
}

err = publisher.Publish(ctx, event)
```

### 2. Consume Events (Worker)

```go
consumer, err := messaging.NewConsumer()
if err != nil {
    log.Fatal("Failed to create consumer:", err)
}
defer consumer.Close()

// Define handler
handler := func(ctx context.Context, event messaging.Event) error {
    log.Printf("Processing: %s", event.Type)

    switch event.Type {
    case "order.created":
        // Process order
        return nil
    default:
        return fmt.Errorf("unknown event type: %s", event.Type)
    }
}

// Start consuming
err = consumer.Consume(
    ctx,
    "order_processing_queue",  // Queue name
    "order.*",                 // Routing pattern
    handler,
)
```

### 3. Routing Patterns

- `order.*` - Matches `order.created`, `order.delivered`, `order.cancelled`
- `order.created` - Exact match
- `notification.#` - Matches `notification.email`, `notification.sms.push`
- `*.created` - Matches any `.created` event

---

## 👷 Worker Services

### Chạy Workers

```bash
# Order worker
make worker-order

# Notification worker
make worker-notification

# Email worker
make worker-email

# SMS worker
make worker-sms
```

### Worker Architecture

```
cmd/
├── server/      # API Server (publishes events)
└── worker/      # Worker Service (consumes events)
    └── main.go  # Worker entry point với routing
```

### Worker Types

1. **Order Worker**: Process order events

   - Queue: `order_queue`
   - Pattern: `order.*`

2. **Notification Worker**: Send emails, SMS
   - Queue: `notification_queue`
   - Pattern: `notification.*`

---

## 📦 Production Deployment

### Docker Compose (Recommended)

Project đã có `docker-compose.yml` với RabbitMQ, PostgreSQL, Redis:

```bash
# Start all infrastructure services
make docker-up

# Start only RabbitMQ
make rabbitmq-up

# View logs
make rabbitmq-logs
```

### Production Configuration

File `configs/rabbitmq/rabbitmq.conf`:

- Memory limit: 60% of available RAM
- Disk limit: 2GB
- Heartbeat: 60 seconds

---

## 📊 Monitoring

### RabbitMQ Management UI

- URL: http://localhost:15672
- Monitor: Queues, Exchanges, Connections, Message rates

### Metrics to Watch

- **Queue depth**: Messages chờ xử lý
- **Consumer lag**: Thời gian messages chờ
- **Publish rate**: Số events publish mỗi giây
- **Consume rate**: Số messages processed mỗi giây
- **Error rate**: Số messages failed

---

## 🔧 Best Practices

### 1. Event Naming Convention

```
{entity}.{action}
- order.created
- order.delivered
- notification.email
```

### 2. Idempotent Handlers

Handlers phải idempotent (có thể chạy nhiều lần):

```go
handler := func(ctx context.Context, event messaging.Event) error {
    // Check if already processed
    if alreadyProcessed(event.AggregateID) {
        return nil // Skip
    }
    // Process event
    processEvent(event)
    markAsProcessed(event.AggregateID)
    return nil
}
```

### 3. Error Handling

```go
handler := func(ctx context.Context, event messaging.Event) error {
    if err := processEvent(event); err != nil {
        if isTransientError(err) {
            return err // Return error to requeue
        }
        // Log permanent errors and skip
        log.Printf("Permanent error, skipping: %v", err)
        return nil
    }
    return nil
}
```

---

## 🐛 Troubleshooting

### Connection Issues

```bash
# Check if RabbitMQ is running
docker ps | grep rabbitmq

# Test connection
docker exec -it foodie-rabbitmq rabbitmq-diagnostics ping

# View logs
make rabbitmq-logs
```

### Queue Backup

```bash
# Check queue depth
docker exec -it foodie-rabbitmq rabbitmqctl list_queues name messages

# Scale up workers
make worker-order &  # Run multiple instances

# Purge queue (⚠️ deletes all messages)
docker exec -it foodie-rabbitmq rabbitmqctl purge_queue order_queue
```

---

## 📚 Tài Liệu Tham Khảo

- [RabbitMQ Documentation](https://www.rabbitmq.com/documentation.html)
- [RabbitMQ Best Practices](https://www.rabbitmq.com/best-practices.html)
- [AMQP Concepts](https://www.rabbitmq.com/tutorials/amqp-concepts.html)

---

## ✅ Quick Start Checklist

- [ ] Start RabbitMQ: `make rabbitmq-up`
- [ ] Configure `.env`: `MESSAGE_BROKER_TYPE=rabbitmq`
- [ ] Test connection: Check management UI
- [ ] Start workers: `make worker-order`
- [ ] Monitor: http://localhost:15672
