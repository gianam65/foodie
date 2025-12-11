# Multi-stage Dockerfile for Go services
# This Dockerfile builds all services and can be used to build specific services

# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata wget

# Set Go toolchain to auto (will download required version if needed)
ENV GOTOOLCHAIN=auto
ENV GOPROXY=https://proxy.golang.org,direct

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build arguments for service selection
ARG SERVICE=server
ARG VERSION=dev
ARG COMMIT_SHA=unknown

# Build the service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.CommitSHA=${COMMIT_SHA}" \
    -o /app/${SERVICE} \
    ./cmd/${SERVICE}

# Stage 2: Runtime
FROM alpine:latest

# Install ca-certificates for HTTPS and wget for healthcheck
RUN apk --no-cache add ca-certificates tzdata wget

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app /app

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port (default, can be overridden)
EXPOSE 8080

# Health check (only for server service)
# Note: This will be overridden in docker-compose for scheduler/worker
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health 2>/dev/null || exit 1

# Default command (can be overridden in docker-compose)
CMD ["/app/server"]

