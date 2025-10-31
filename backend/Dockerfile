# Multi-stage build for Go application
FROM golang:1.21-alpine AS builder

# Install necessary packages
RUN apk add --no-cache git ca-certificates tzdata sqlite

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o nutrition-platform .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates curl sqlite

# Create app user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/nutrition-platform .

# Copy migration files
COPY --from=builder /app/migrations ./migrations

# Create data directory for SQLite
RUN mkdir -p /app/data && chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Start the application
CMD ["./nutrition-platform"]