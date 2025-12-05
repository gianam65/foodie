# API Definitions

This directory contains API contract definitions for the Food Delivery service.

## Files

- `openapi.yaml` - OpenAPI 3.0 specification for REST API
- `*.proto` - Protocol Buffer definitions for gRPC services (to be added)

## Usage

### OpenAPI

View and test the API:

```bash
# Using Swagger UI
docker run -p 8081:8080 -e SWAGGER_JSON=/api/openapi.yaml -v $(pwd)/api:/api swaggerapi/swagger-ui

# Or use online editor
# https://editor.swagger.io/
```

Generate client code:

```bash
# Using openapi-generator
openapi-generator-cli generate -i api/openapi.yaml -g go -o client/go
```

### gRPC

Generate Go code from proto files:

```bash
protoc --go_out=. --go-grpc_out=. api/*.proto
```

## API Versioning

APIs are versioned in the URL path:

- `/api/v1/orders` - Version 1
- `/api/v2/orders` - Version 2 (future)

## Documentation

API documentation is automatically served from OpenAPI spec when the server is running:

- **Swagger UI**: `http://localhost:8080/swagger/`
- **OpenAPI Spec (JSON)**: `http://localhost:8080/swagger/doc.json`
- **OpenAPI Spec (YAML)**: `http://localhost:8080/swagger/openapi.yaml`

### Accessing Swagger UI

1. Start the server:

   ```bash
   make dev
   # or
   go run ./cmd/server
   ```

2. Open browser and navigate to:

   ```
   http://localhost:8080/swagger/
   ```

3. You can now:
   - View all API endpoints
   - See request/response schemas
   - Test APIs directly from the UI
   - View authentication requirements

### Integration

Swagger UI is automatically integrated and served from the `/swagger/` endpoint. The OpenAPI spec file (`api/openapi.yaml`) is automatically loaded when the server starts.
