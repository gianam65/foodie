# ⏰ Scheduler Setup Guide

## 📋 Tổng Quan

Scheduler service quản lý các scheduled tasks (cron jobs) chạy định kỳ trong hệ thống. Sử dụng thư viện `github.com/robfig/cron/v3` để thực thi các tasks theo lịch.

## 🏗️ Kiến Trúc

```
cmd/scheduler/
└── main.go                    # Entry point cho scheduler service

internal/infrastructure/scheduler/
├── scheduler.go               # Scheduler core (quản lý tasks)
└── tasks/
    ├── order_cleanup_task.go  # Order cleanup tasks
    └── health_check_task.go   # Health check task
```

## 🚀 Sử Dụng

### Chạy Scheduler

#### Cách 1: Sử dụng Makefile

```bash
make scheduler
```

#### Cách 2: Chạy trực tiếp

```bash
go run ./cmd/scheduler
```

#### Cách 3: Build và chạy binary

```bash
make scheduler-build
./bin/scheduler
```

## 📅 Cron Expression Format

Scheduler sử dụng cron expression với format 6 fields (bao gồm seconds):

```
┌───────────── second (0-59)
│ ┌─────────── minute (0-59)
│ │ ┌───────── hour (0-23)
│ │ │ ┌─────── day of month (1-31)
│ │ │ │ ┌───── month (1-12)
│ │ │ │ │ ┌─── day of week (0-6) (Sunday=0)
│ │ │ │ │ │
* * * * * *
```

### Ví dụ Cron Expressions

| Expression          | Mô tả                           |
| ------------------- | ------------------------------- |
| `0 * * * * *`       | Mỗi phút                        |
| `0 */5 * * * *`     | Mỗi 5 phút                      |
| `0 0 * * * *`       | Mỗi giờ (vào đầu giờ)           |
| `0 0 0 * * *`       | Mỗi ngày lúc 00:00              |
| `0 0 9 * * MON-FRI` | Mỗi ngày trong tuần lúc 9:00 AM |
| `0 0 2 * * *`       | Mỗi ngày lúc 2:00 AM            |
| `0 0 3 * * 0`       | Mỗi Chủ Nhật lúc 3:00 AM        |
| `0 30 14 * * *`     | Mỗi ngày lúc 14:30              |

## 📝 Scheduled Tasks

### 1. Health Check Task

- **Tên**: `health_check`
- **Schedule**: Mỗi 5 phút (`0 */5 * * * *`)
- **Mô tả**: Kiểm tra health của các services (database, cache, external APIs)

### 2. Order Cleanup Task

- **Tên**: `order_cleanup`
- **Schedule**: Mỗi ngày lúc 2:00 AM (`0 0 2 * * *`)
- **Mô tả**: Xóa các orders cũ (ví dụ: cancelled orders > 30 ngày)

### 3. Cleanup Completed Orders Task

- **Tên**: `cleanup_completed_orders`
- **Schedule**: Mỗi Chủ Nhật lúc 3:00 AM (`0 0 3 * * 0`)
- **Mô tả**: Archive các completed orders cũ (> 90 ngày)

## ➕ Thêm Task Mới

### Bước 1: Tạo Task Implementation

Tạo file mới trong `internal/infrastructure/scheduler/tasks/`:

```go
// internal/infrastructure/scheduler/tasks/my_task.go
package tasks

import (
	"context"
	"log"
)

type MyTask struct {
	logger *log.Logger
	// Add dependencies here
}

func NewMyTask(logger *log.Logger) *MyTask {
	return &MyTask{
		logger: logger,
	}
}

func (t *MyTask) Name() string {
	return "my_task"
}

func (t *MyTask) Run(ctx context.Context) error {
	t.logger.Printf("Running my task...")

	// Implement task logic here

	return nil
}
```

### Bước 2: Register Task trong Scheduler

Cập nhật `cmd/scheduler/main.go`:

```go
func registerTasks(sched *scheduler.Scheduler, repos *database.Repositories, logger *log.Logger) {
	// ... existing tasks ...

	// Add your new task
	myTask := tasks.NewMyTask(logger)
	if err := sched.AddTask("0 */10 * * * *", myTask); err != nil {
		logger.Printf("Failed to register my task: %v", err)
	}
}
```

## 🔧 Cấu Hình

Scheduler sử dụng cùng cấu hình database như API server:

```env
# Database Configuration
SQL_DSN=postgres://user:password@localhost:5432/foodie?sslmode=disable
SQL_DRIVER=postgres
```

## 🎯 Use Cases

### 1. Order Management

- **Auto-complete orders**: Tự động đánh dấu delivered sau thời gian delivery
- **Cleanup old orders**: Xóa/archive orders cũ
- **Update order status**: Cập nhật status dựa trên thời gian

### 2. Notifications

- **Send reminders**: Gửi nhắc nhở cho orders chưa confirm
- **Delivery notifications**: Thông báo delivery status
- **Promotional emails**: Gửi email khuyến mãi định kỳ

### 3. Reports & Analytics

- **Daily reports**: Tạo báo cáo hàng ngày
- **Weekly summaries**: Tổng hợp tuần
- **Monthly analytics**: Phân tích tháng

### 4. Data Maintenance

- **Cache cleanup**: Xóa cache cũ
- **Log rotation**: Xoay log files
- **Database optimization**: Tối ưu database

## 🔍 Monitoring

Scheduler log tất cả task executions:

```
scheduler 2024/12/05 14:30:00 main.go:45: Running scheduled task: health_check
scheduler 2024/12/05 14:30:00 main.go:49: Task health_check completed successfully
```

## 🛠️ Best Practices

1. **Idempotency**: Đảm bảo tasks có thể chạy lại mà không gây side effects
2. **Error Handling**: Log errors nhưng không crash scheduler
3. **Context**: Sử dụng context để handle cancellation
4. **Timeout**: Set timeout cho long-running tasks
5. **Locking**: Sử dụng distributed locks nếu chạy multiple scheduler instances

## 📦 Deployment

### Standalone Service

Chạy scheduler như một service riêng:

```bash
# Development
make scheduler

# Production
./bin/scheduler
```

### Docker

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o scheduler ./cmd/scheduler

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/scheduler .
CMD ["./scheduler"]
```

### Systemd Service

```ini
[Unit]
Description=Foodie Scheduler Service
After=network.target postgresql.service

[Service]
Type=simple
User=foodie
WorkingDirectory=/opt/foodie/backend
ExecStart=/opt/foodie/backend/bin/scheduler
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

## 🔒 Production Considerations

1. **Multiple Instances**: Chỉ chạy 1 scheduler instance hoặc dùng distributed locks
2. **Monitoring**: Monitor task execution times và failures
3. **Alerting**: Alert khi tasks fail liên tục
4. **Logging**: Log đầy đủ để debug
5. **Graceful Shutdown**: Handle shutdown signals properly
