# Foodie Backend

Backend system for food delivery platform built with **Golang**, following **Modular Monolith** architecture with **Clean Architecture** principles.

## 📚 Documentation

**👉 Xem tài liệu đầy đủ tại [docs/](./docs/)**

### Quick Links

- 📖 [Tổng quan dự án](./docs/PROJECT_OVERVIEW.md)
- 🚀 [Quick Start Guide](./docs/QUICKSTART.md)
- 🔧 [Redis Cache Guide](./docs/setup/REDIS_CACHE.md)
- 📡 [API Documentation](./docs/api/README.md)

## 🏗️ Architecture

This project uses:

- **Modular Monolith**: All modules in one codebase, easy to split into microservices later
- **Clean Architecture**: Clear separation between domain, application, and infrastructure layers
- **Hexagonal Architecture**: Adapters for external dependencies (database, messaging, APIs)
- **Repository Pattern**: Abstract data access with SQL implementations

## 🚀 Quick Start

```bash
# Run server
make run

# Run with hot reload
make dev

# Run all services
make dev-all
```

See [docs/QUICKSTART.md](./docs/QUICKSTART.md) for detailed instructions.

## 📁 Project Structure

```
backend/
├── cmd/              # Application entry points
├── internal/         # Internal application code
│   ├── interfaces/  # Interface layer (HTTP, gRPC)
│   ├── domain/      # Domain layer (business logic)
│   └── infrastructure/ # Infrastructure layer (DB, cache, external services)
├── pkg/             # Shared packages
├── docs/            # Documentation
└── migrations/      # Database migrations
```

## 🔧 Features

- ✅ SQL database support (PostgreSQL/MySQL)
- ✅ Redis cache integration
- ✅ Clean architecture with dependency injection
- ✅ HTTP and gRPC support
- ✅ Modular design for easy scaling

## 📖 Learn More

- [Project Overview](./docs/PROJECT_OVERVIEW.md) - Detailed architecture and folder structure
- [Redis Cache](./docs/setup/REDIS_CACHE.md) - How to use Redis caching
- [API Docs](./docs/api/README.md) - API specifications

## 📄 License

[Your License Here]
