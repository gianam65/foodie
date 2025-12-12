.PHONY: help run dev migrate-up migrate-down build clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application
	@go run ./cmd/server

dev: ## Run with hot reload (requires air: go install github.com/air-verse/air@latest)
	@if command -v air > /dev/null || [ -f "$(shell go env GOPATH)/bin/air" ]; then \
		$(shell go env GOPATH)/bin/air 2>/dev/null || air; \
	else \
		echo "Error: air is not installed."; \
		echo "Install it with: go install github.com/air-verse/air@latest"; \
		echo "Then add \$$(go env GOPATH)/bin to your PATH"; \
		exit 1; \
	fi

migrate-up: ## Run database migrations
	@go run ./cmd/migrate -command=up -driver=postgres

migrate-down: ## Rollback last migration
	@go run ./cmd/migrate -command=down -driver=postgres

schema-check: ## Check database schema for a table (usage: make schema-check TABLE=products)
	@if [ -z "$(TABLE)" ]; then \
		echo "Error: TABLE is required. Usage: make schema-check TABLE=products"; \
		exit 1; \
	fi
	@go run ./cmd/schema-check -table=$(TABLE)

schema-sync: ## Auto-detect and create migration for missing columns (usage: make schema-sync TABLE=products)
	@if [ -z "$(TABLE)" ]; then \
		echo "Error: TABLE is required. Usage: make schema-sync TABLE=products"; \
		exit 1; \
	fi
	@go run ./cmd/schema-sync -table=$(TABLE)

migrate-create: ## Create a new migration file (usage: make migrate-create NAME=add_user_table)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=add_user_table"; \
		exit 1; \
	fi
	@if command -v migrate > /dev/null; then \
		migrate create -ext sql -dir migrations -seq $(NAME); \
	else \
		echo "Installing golang-migrate CLI..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
		$$(go env GOPATH)/bin/migrate create -ext sql -dir migrations -seq $(NAME); \
	fi
	@echo "✅ Migration files created! Edit migrations/*_$(NAME).up.sql and migrations/*_$(NAME).down.sql"

build: ## Build the application
	@go build -o bin/server ./cmd/server

clean: ## Clean build artifacts
	@rm -rf tmp/ bin/ *.log

test: ## Run tests
	@go test ./...

fmt: ## Format code
	@go fmt ./...

lint: ## Run linter (if golangci-lint is installed)
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install: https://golangci-lint.run/usage/install/"; \
	fi

# Docker commands
docker-up: ## Start all services (RabbitMQ, PostgreSQL, Redis)
	@docker-compose up -d

docker-down: ## Stop all services
	@docker-compose down

docker-logs: ## View docker logs
	@docker-compose logs -f

rabbitmq-up: ## Start only RabbitMQ
	@docker-compose up -d rabbitmq

rabbitmq-logs: ## View RabbitMQ logs
	@docker-compose logs -f rabbitmq

# Worker commands
worker-order: ## Run order worker
	@go run ./cmd/worker order

worker-notification: ## Run notification worker
	@go run ./cmd/worker notification

worker-email: ## Run email worker
	@go run ./cmd/worker email

worker-sms: ## Run SMS worker
	@go run ./cmd/worker sms

# Scheduler commands
scheduler: ## Run scheduler service
	@go run ./cmd/scheduler

scheduler-build: ## Build scheduler binary
	@go build -o bin/scheduler ./cmd/scheduler

# Build all binaries
build-all: ## Build all service binaries
	@echo "Building all services..."
	@mkdir -p bin
	@go build -o bin/server ./cmd/server
	@go build -o bin/scheduler ./cmd/scheduler
	@go build -o bin/worker ./cmd/worker
	@echo "✅ All binaries built in bin/"

# Start all services
start-all: ## Start all services (server, scheduler, workers) in background
	@./scripts/start-all.sh

# Development mode: Start all services with visible logs
dev-all: ## Start all services in development mode (with tmux if available)
	@./scripts/dev-all.sh

# Stop all services (if using start-all.sh)
stop-all: ## Stop all running services
	@./scripts/stop-all.sh

