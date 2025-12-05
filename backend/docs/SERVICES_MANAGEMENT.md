# 🚀 Services Management Guide

## 📋 Overview

Project có nhiều services cần chạy:

- **HTTP Server** - API server
- **Scheduler** - Scheduled tasks (cron jobs)
- **Workers** - Background job processors

---

## 🎯 Cách Chạy Services

### Option 1: Chạy Tất Cả Services (Recommended)

#### **Development Mode** (logs visible)

```bash
# Sử dụng tmux để quản lý multiple terminals
./scripts/dev-all.sh

# Hoặc với Makefile
make dev-all
```

**Features:**

- ✅ Tất cả services chạy trong tmux session
- ✅ Logs visible trong terminal
- ✅ Dễ dàng monitor tất cả services
- ✅ Dễ dàng stop tất cả (Ctrl+C)

**Nếu không có tmux:**

- Services sẽ chạy trong background
- Logs được redirect vào `logs/` directory

#### **Production Mode** (background)

```bash
# Chạy tất cả services trong background
./scripts/start-all.sh

# Hoặc với Makefile
make start-all
```

**Features:**

- ✅ Services chạy trong background
- ✅ Logs được redirect vào files
- ✅ PIDs được lưu trong `.pids` file
- ✅ Dễ dàng stop với `make stop-all`

---

### Option 2: Chạy Từng Service Riêng Lẻ

#### HTTP Server

```bash
# Development (hot reload)
make dev

# Production
make run
```

#### Scheduler

```bash
go run ./cmd/scheduler
# hoặc
make scheduler
```

#### Workers

```bash
# Order worker
go run ./cmd/worker order
# hoặc
make worker-order

# Notification worker
go run ./cmd/worker notification
# hoặc
make worker-notification
```

---

## 📊 Services Overview

| Service                 | Command                    | Port/Description   |
| ----------------------- | -------------------------- | ------------------ |
| **HTTP Server**         | `make dev`                 | `:8080`            |
| **Scheduler**           | `make scheduler`           | Background service |
| **Order Worker**        | `make worker-order`        | Background service |
| **Notification Worker** | `make worker-notification` | Background service |

---

## 🛠️ Development Workflow

### 1. **Start Infrastructure Services** (Docker)

```bash
# Start PostgreSQL, Redis, RabbitMQ
make docker-up

# Hoặc chỉ start một service
make rabbitmq-up
```

### 2. **Run Database Migrations**

```bash
make migrate-up
```

### 3. **Start Application Services**

```bash
# Development mode (recommended)
./dev-all.sh

# Hoặc start từng service riêng
```

---

## 📝 Logs Management

### Logs Location

```
logs/
├── server.log              # HTTP server logs
├── scheduler.log           # Scheduler logs
├── worker-order.log        # Order worker logs
└── worker-notification.log # Notification worker logs
```

### View Logs

```bash
# View all logs
tail -f logs/*.log

# View specific service
tail -f logs/server.log
tail -f logs/scheduler.log
```

### With tmux (dev-all.sh)

Khi dùng `./scripts/dev-all.sh` với tmux, logs hiển thị trực tiếp trong terminal. Có thể:

- Scroll trong mỗi pane
- Switch giữa panes với arrow keys
- Detach và reattach sau

---

## 🛑 Stop Services

### Stop All Services

```bash
# Nếu dùng start-all.sh
make stop-all
# hoặc
./scripts/stop-all.sh

# Nếu dùng dev-all.sh
# Press Ctrl+C trong terminal
# Hoặc nếu đã detach tmux:
tmux kill-session -t foodie-backend
```

### Stop Individual Service

```bash
# Find PID
ps aux | grep "go run ./cmd/server"

# Kill process
kill <PID>
```

---

## 🔧 Build Binaries

### Build All Services

```bash
make build-all
```

Output:

```
bin/
├── server     # HTTP server binary
├── scheduler  # Scheduler binary
└── worker     # Worker binary
```

### Run Built Binaries

```bash
# Start với binaries (faster startup)
./bin/server
./bin/scheduler
./bin/worker order
```

---

## 📋 Complete Startup Sequence

### Full Development Setup

```bash
# 1. Start infrastructure
make docker-up

# 2. Run migrations
make migrate-up

# 3. Start all application services
./dev-all.sh
```

### Full Production Setup

```bash
# 1. Start infrastructure
make docker-up

# 2. Run migrations
make migrate-up

# 3. Build binaries
make build-all

# 4. Start all services
./start-all.sh
```

---

## 🐛 Troubleshooting

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

### Services Not Starting

1. **Check logs:**

   ```bash
   tail -f logs/*.log
   ```

2. **Check environment:**

   ```bash
   # Ensure .env file exists
   cp .env.example .env
   ```

3. **Check dependencies:**
   ```bash
   # Ensure database is running
   docker ps | grep postgres
   ```

### tmux Not Available

Install tmux:

```bash
# macOS
brew install tmux

# Linux
sudo apt-get install tmux
```

Nếu không có tmux, `scripts/dev-all.sh` sẽ fallback về background mode với log files.

---

## 📚 Related Commands

Xem tất cả available commands:

```bash
make help
```

---

## 💡 Tips

1. **Development**: Dùng `./scripts/dev-all.sh` với tmux để monitor tất cả services
2. **Testing**: Dùng `make dev` để test HTTP server với hot reload
3. **Production**: Build binaries trước với `make build-all`, sau đó dùng `./scripts/start-all.sh`
4. **Debugging**: Check logs trong `logs/` directory hoặc view trong tmux panes
