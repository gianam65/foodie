#!/bin/bash

# Foodie Backend - Development Mode (All Services with Visible Logs)
# Script để chạy tất cả services với logs visible trong terminal

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Error: Go is not installed${NC}"
    exit 1
fi

# Check if tmux is available (for better terminal management)
if command -v tmux &> /dev/null; then
    USE_TMUX=true
else
    USE_TMUX=false
    echo -e "${YELLOW}⚠️  tmux not found. Services will run in background.${NC}"
    echo -e "${YELLOW}   Install tmux for better terminal management: brew install tmux${NC}"
    echo ""
fi

echo -e "${BLUE}🔥 Foodie Backend - Development Mode (All Services)${NC}"
echo ""

if [ "$USE_TMUX" = true ]; then
    # Use tmux to run all services in separate panes
    SESSION_NAME="foodie-backend"
    
    # Kill existing session if exists
    tmux kill-session -t "$SESSION_NAME" 2>/dev/null || true
    
    # Create new tmux session
    tmux new-session -d -s "$SESSION_NAME" -n "server"
    
    # Split window and start services
    tmux send-keys -t "$SESSION_NAME:server" "echo 'Starting HTTP Server...' && go run ./cmd/server" C-m
    sleep 1
    
    tmux split-window -h -t "$SESSION_NAME:server"
    tmux send-keys -t "$SESSION_NAME:server" "echo 'Starting Scheduler...' && go run ./cmd/scheduler" C-m
    sleep 1
    
    tmux split-window -v -t "$SESSION_NAME:server"
    tmux send-keys -t "$SESSION_NAME:server" "echo 'Starting Order Worker...' && go run ./cmd/worker order" C-m
    sleep 1
    
    tmux split-window -v -t "$SESSION_NAME:server"
    tmux send-keys -t "$SESSION_NAME:server" "echo 'Starting Notification Worker...' && go run ./cmd/worker notification" C-m
    
    # Select layout
    tmux select-layout -t "$SESSION_NAME:server" tiled
    
    echo -e "${GREEN}✅ All services started in tmux session: $SESSION_NAME${NC}"
    echo ""
    echo -e "${BLUE}📊 Services:${NC}"
    echo "   - HTTP Server: http://localhost:8080"
    echo "   - Swagger UI: http://localhost:8080/swagger/"
    echo "   - Scheduler: Running"
    echo "   - Workers: Running"
    echo ""
    echo -e "${YELLOW}Attach to tmux session:${NC}"
    echo "   tmux attach -t $SESSION_NAME"
    echo ""
    echo -e "${YELLOW}Detach from session:${NC}"
    echo "   Press Ctrl+B, then D"
    echo ""
    echo -e "${YELLOW}Kill session:${NC}"
    echo "   tmux kill-session -t $SESSION_NAME"
    echo ""
    
    # Attach to session
    tmux attach -t "$SESSION_NAME"
else
    # Fallback: run in background with logs
    echo -e "${YELLOW}Running services in background mode...${NC}"
    echo "   (Install tmux for better terminal management)"
    echo ""
    
    mkdir -p logs
    
    # Start services in background
    go run ./cmd/server > logs/server.log 2>&1 &
    SERVER_PID=$!
    echo -e "${GREEN}✅ HTTP Server started (PID: $SERVER_PID)${NC}"
    echo "   Logs: logs/server.log"
    echo "   Access: http://localhost:8080/swagger/"
    
    go run ./cmd/scheduler > logs/scheduler.log 2>&1 &
    SCHEDULER_PID=$!
    echo -e "${GREEN}✅ Scheduler started (PID: $SCHEDULER_PID)${NC}"
    echo "   Logs: logs/scheduler.log"
    
    go run ./cmd/worker order > logs/worker-order.log 2>&1 &
    WORKER_ORDER_PID=$!
    echo -e "${GREEN}✅ Order Worker started (PID: $WORKER_ORDER_PID)${NC}"
    echo "   Logs: logs/worker-order.log"
    
    go run ./cmd/worker notification > logs/worker-notification.log 2>&1 &
    WORKER_NOTIF_PID=$!
    echo -e "${GREEN}✅ Notification Worker started (PID: $WORKER_NOTIF_PID)${NC}"
    echo "   Logs: logs/worker-notification.log"
    
    echo ""
    echo -e "${BLUE}All services running!${NC}"
    echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"
    echo ""
    echo "View logs:"
    echo "   tail -f logs/server.log"
    echo "   tail -f logs/scheduler.log"
    echo "   tail -f logs/worker-*.log"
    echo ""
    
    # Wait for Ctrl+C
    trap "kill $SERVER_PID $SCHEDULER_PID $WORKER_ORDER_PID $WORKER_NOTIF_PID 2>/dev/null; exit" SIGINT SIGTERM
    wait
fi

