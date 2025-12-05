#!/bin/bash

# Foodie Backend - Start All Services
# Script để chạy tất cả services: HTTP server, workers, scheduler

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# PIDs file để track running processes
PIDS_FILE="$SCRIPT_DIR/.pids"

# Function to cleanup processes on exit
cleanup() {
    echo -e "\n${YELLOW}🛑 Stopping all services...${NC}"
    if [ -f "$PIDS_FILE" ]; then
        while read pid; do
            if kill -0 "$pid" 2>/dev/null; then
                kill "$pid" 2>/dev/null || true
            fi
        done < "$PIDS_FILE"
        rm -f "$PIDS_FILE"
    fi
    echo -e "${GREEN}✅ All services stopped${NC}"
    exit 0
}

# Trap signals to cleanup
trap cleanup SIGINT SIGTERM

# Remove old PIDs file
rm -f "$PIDS_FILE"

echo -e "${BLUE}🚀 Foodie Backend - Starting All Services${NC}"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Error: Go is not installed${NC}"
    exit 1
fi

# Check if .env file exists
if [ ! -f .env ]; then
    echo -e "${YELLOW}⚠️  Warning: .env file not found${NC}"
    echo "The application will continue with default values."
    echo ""
fi

# Build binaries (optional, for faster startup)
echo -e "${BLUE}📦 Building services...${NC}"
go build -o bin/server ./cmd/server 2>/dev/null || true
go build -o bin/scheduler ./cmd/scheduler 2>/dev/null || true
go build -o bin/worker ./cmd/worker 2>/dev/null || true
echo -e "${GREEN}✅ Build completed${NC}"
echo ""

# Start HTTP Server
echo -e "${GREEN}▶️  Starting HTTP Server...${NC}"
if [ -f bin/server ]; then
    ./bin/server > logs/server.log 2>&1 &
else
    go run ./cmd/server > logs/server.log 2>&1 &
fi
SERVER_PID=$!
echo "$SERVER_PID" >> "$PIDS_FILE"
echo -e "${GREEN}✅ HTTP Server started (PID: $SERVER_PID)${NC}"
echo "   Access: http://localhost:8080"
echo "   Swagger: http://localhost:8080/swagger/"
sleep 2

# Start Scheduler
echo -e "${GREEN}▶️  Starting Scheduler...${NC}"
if [ -f bin/scheduler ]; then
    ./bin/scheduler > logs/scheduler.log 2>&1 &
else
    go run ./cmd/scheduler > logs/scheduler.log 2>&1 &
fi
SCHEDULER_PID=$!
echo "$SCHEDULER_PID" >> "$PIDS_FILE"
echo -e "${GREEN}✅ Scheduler started (PID: $SCHEDULER_PID)${NC}"
sleep 1

# Start Workers (optional - comment out if not needed)
echo -e "${GREEN}▶️  Starting Workers...${NC}"

# Order Worker
if [ -f bin/worker ]; then
    ./bin/worker order > logs/worker-order.log 2>&1 &
else
    go run ./cmd/worker order > logs/worker-order.log 2>&1 &
fi
WORKER_ORDER_PID=$!
echo "$WORKER_ORDER_PID" >> "$PIDS_FILE"
echo -e "${GREEN}✅ Order Worker started (PID: $WORKER_ORDER_PID)${NC}"

# Notification Worker
if [ -f bin/worker ]; then
    ./bin/worker notification > logs/worker-notification.log 2>&1 &
else
    go run ./cmd/worker notification > logs/worker-notification.log 2>&1 &
fi
WORKER_NOTIFICATION_PID=$!
echo "$WORKER_NOTIFICATION_PID" >> "$PIDS_FILE"
echo -e "${GREEN}✅ Notification Worker started (PID: $WORKER_NOTIFICATION_PID)${NC}"

echo ""
# Create logs directory if it doesn't exist
mkdir -p logs

echo ""
echo -e "${GREEN}✅ All services started successfully!${NC}"
echo ""
echo -e "${BLUE}📊 Running Services:${NC}"
echo "   - HTTP Server: http://localhost:8080 (PID: $SERVER_PID)"
echo "   - Swagger UI: http://localhost:8080/swagger/"
echo "   - Scheduler: (PID: $SCHEDULER_PID)"
echo "   - Order Worker: (PID: $WORKER_ORDER_PID)"
echo "   - Notification Worker: (PID: $WORKER_NOTIFICATION_PID)"
echo ""
echo -e "${BLUE}📝 Logs:${NC}"
echo "   - Server: logs/server.log"
echo "   - Scheduler: logs/scheduler.log"
echo "   - Workers: logs/worker-*.log"
echo ""
echo -e "${BLUE}📋 View Logs:${NC}"
echo "   tail -f logs/server.log"
echo "   tail -f logs/scheduler.log"
echo "   tail -f logs/*.log"
echo ""
echo -e "${BLUE}🛑 Stop Services:${NC}"
echo "   make stop-all"
echo "   hoặc: ./scripts/stop-all.sh"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"
echo ""

# Wait for all background processes
wait

