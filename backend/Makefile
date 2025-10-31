# Makefile for Nutrition Platform Backend

.PHONY: build test test-coverage test-integration clean run dev docker-build docker-run lint security-audit

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=nutrition-platform
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...

# Run tests
test:
	$(GOTEST) -v ./tests/...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./tests/...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run integration tests
test-integration:
	$(GOTEST) -v -tags=integration ./tests/...

# Run all tests with race detection
test-race:
	$(GOTEST) -v -race ./tests/...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out coverage.html

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

# Run in development mode with hot reload
dev:
	air

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

# Docker build
docker-build:
	docker build -t nutrition-platform-backend .

# Docker run
docker-run:
	docker run -p 8080:8080 nutrition-platform-backend

# Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) verify

# Update dependencies
deps-update:
	$(GOMOD) tidy
	$(GOGET) -u ./...

# Lint code
lint:
	golangci-lint run

# Security audit
security-audit:
	gosec ./...
	nancy sleuth

# Format code
fmt:
	$(GOCMD) fmt ./...

# Vet code
vet:
	$(GOCMD) vet ./...

# Generate API documentation
docs:
	swag init -g main.go

# Database migrations
migrate-up:
	migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

# Create new migration
migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(name)

# Benchmark tests
benchmark:
	$(GOTEST) -bench=. -benchmem ./tests/...

# Profile application
profile:
	$(GOTEST) -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. ./tests/...

# Install development tools
install-tools:
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) github.com/securecodewarrior/nancy@latest
	$(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	$(GOGET) github.com/swaggo/swag/cmd/swag@latest
	$(GOGET) github.com/cosmtrek/air@latest

# Check test coverage threshold
check-coverage:
	@echo "Checking test coverage..."
	@$(GOTEST) -coverprofile=coverage.out ./tests/... > /dev/null
	@coverage=$$($(GOCMD) tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$coverage < 70" | bc) -eq 1 ]; then \
		echo "❌ Test coverage is $$coverage%, which is below the required 70%"; \
		exit 1; \
	else \
		echo "✅ Test coverage is $$coverage%, which meets the required 70%"; \
	fi

# Pre-commit checks
pre-commit: fmt vet lint test check-coverage security-audit
	@echo "✅ All pre-commit checks passed!"

# CI pipeline
ci: deps fmt vet lint test-coverage security-audit
	@echo "✅ CI pipeline completed successfully!"

# Help
help:
	@echo "Available commands:"
	@echo "  build           - Build the application"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  test-integration- Run integration tests"
	@echo "  test-race       - Run tests with race detection"
	@echo "  clean           - Clean build artifacts"
	@echo "  run             - Build and run the application"
	@echo "  dev             - Run in development mode with hot reload"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run Docker container"
	@echo "  lint            - Lint code"
	@echo "  security-audit  - Run security audit"
	@echo "  fmt             - Format code"
	@echo "  vet             - Vet code"
	@echo "  docs            - Generate API documentation"
	@echo "  pre-commit      - Run all pre-commit checks"
	@echo "  ci              - Run CI pipeline"
	@echo "  help            - Show this help message"