#!/bin/bash

# Foodie Backend - Stop All Services
# Script để dừng tất cả services đã được start bởi start-all.sh

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

PIDS_FILE="$SCRIPT_DIR/.pids"

echo -e "${BLUE}🛑 Stopping all services...${NC}"
echo ""

if [ ! -f "$PIDS_FILE" ]; then
    echo -e "${YELLOW}⚠️  No .pids file found. Services may not be running.${NC}"
    echo ""
    echo "Trying to find and stop processes manually..."
    
    # Try to find processes by name
    pkill -f "go run ./cmd/server" 2>/dev/null || true
    pkill -f "go run ./cmd/scheduler" 2>/dev/null || true
    pkill -f "go run ./cmd/worker" 2>/dev/null || true
    pkill -f "bin/server" 2>/dev/null || true
    pkill -f "bin/scheduler" 2>/dev/null || true
    pkill -f "bin/worker" 2>/dev/null || true
    
    echo -e "${GREEN}✅ Cleaned up processes${NC}"
    exit 0
fi

# Read and kill all PIDs
STOPPED=0
while read pid; do
    if kill -0 "$pid" 2>/dev/null; then
        echo -e "${YELLOW}Stopping process $pid...${NC}"
        kill "$pid" 2>/dev/null || true
        sleep 1
        # Force kill if still running
        if kill -0 "$pid" 2>/dev/null; then
            kill -9 "$pid" 2>/dev/null || true
        fi
        STOPPED=$((STOPPED + 1))
    fi
done < "$PIDS_FILE"

# Remove PIDs file
rm -f "$PIDS_FILE"

if [ $STOPPED -gt 0 ]; then
    echo ""
    echo -e "${GREEN}✅ Stopped $STOPPED service(s)${NC}"
else
    echo ""
    echo -e "${YELLOW}⚠️  No running services found${NC}"
fi

