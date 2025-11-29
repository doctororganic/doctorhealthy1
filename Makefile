.PHONY: help build test clean deploy setup-dev setup-prod lint format security docker-build

# Default target
help:
	@echo "Nutrition Platform - Available Commands:"
	@echo ""
	@echo "Development:"
	@echo "  make setup-dev      - Set up development environment"
	@echo "  make setup-prod     - Set up production environment"
	@echo "  make dev            - Run development servers"
	@echo "  make build          - Build backend and frontend"
	@echo "  make test           - Run all tests"
	@echo ""
	@echo "Backend:"
	@echo "  make test-backend   - Run backend tests"
	@echo "  make run-backend    - Run backend server"
	@echo "  make backend-build   - Build backend binary"
	@echo "  make backend-lint    - Lint backend code"
	@echo "  make backend-format  - Format backend code"
	@echo ""
	@echo "Frontend:"
	@echo "  make test-frontend  - Run frontend tests"
	@echo "  make run-frontend   - Run frontend dev server"
	@echo "  make frontend-build  - Build frontend"
	@echo "  make frontend-lint   - Lint frontend code"
	@echo "  make frontend-format - Format frontend code"
	@echo ""
	@echo "Deployment:"
	@echo "  make deploy-staging - Deploy to staging"
	@echo "  make deploy-prod    - Deploy to production"
	@echo "  make docker-build   - Build Docker images"
	@echo "  make security       - Run security scans"
	@echo ""
	@echo "Maintenance:"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make lint          - Lint all code"
	@echo "  make format        - Format all code"
	@echo "  make db-migrate    - Run database migrations"
	@echo "  make db-backup     - Backup database"

# Development setup
setup-dev:
	@echo "Setting up development environment..."
	cd backend && cp .env.example .env || echo "No .env.example found in backend"
	cd frontend-nextjs && cp .env.local.example .env.local || echo "No .env.local.example found in frontend"
	cd backend && go mod download
	cd frontend-nextjs && npm install
	@echo "✅ Development environment setup complete"

setup-prod:
	@echo "Setting up production environment..."
	@echo "⚠️  This is for production setup - ensure you have proper credentials"
	cd backend && go mod download
	cd frontend-nextjs && npm ci --production
	@echo "✅ Production environment setup complete"

# Development servers
dev:
	@echo "Starting development servers..."
	@echo "Backend will run on :8080"
	@echo "Frontend will run on :3000"
	@echo "Press Ctrl+C to stop both servers"
	@make -j2 run-backend run-frontend

# Build commands
build: backend-build frontend-build

backend-build:
	@echo "Building backend..."
	cd backend && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nutrition-platform .
	@echo "✅ Backend built"

frontend-build:
	@echo "Building frontend..."
	cd frontend-nextjs && npm run build
	@echo "✅ Frontend built"

# Test commands
test: test-backend test-frontend

test-backend:
	@echo "Running backend tests..."
	cd backend && go test ./... -v -race -cover

test-frontend:
	@echo "Running frontend tests..."
	cd frontend-nextjs && npm test -- --coverage --watchAll=false

# Backend commands
run-backend:
	@echo "Starting backend server..."
	cd backend && go run main.go

backend-lint:
	@echo "Linting backend..."
	cd backend && go vet ./...
	cd backend && golint ./... || echo "golint not installed, skipping"

backend-format:
	@echo "Formatting backend..."
	cd backend && go fmt ./...

# Frontend commands
run-frontend:
	@echo "Starting frontend dev server..."
	cd frontend-nextjs && npm run dev

frontend-lint:
	@echo "Linting frontend..."
	cd frontend-nextjs && npm run lint

frontend-format:
	@echo "Formatting frontend..."
	cd frontend-nextjs && npm run format || echo "No format script found"

# Deployment commands
deploy-staging:
	@echo "Deploying to staging..."
	./scripts/deploy.sh staging latest

deploy-prod:
	@echo "Deploying to production..."
	@echo "⚠️  Press Ctrl+C to cancel within 10 seconds..."
	@sleep 10
	./scripts/deploy.sh production latest

# Docker commands
docker-build:
	@echo "Building Docker images..."
	docker build -t nutrition-platform/backend ./backend
	docker build -t nutrition-platform/frontend ./frontend-nextjs
	@echo "✅ Docker images built"

# Security commands
security:
	@echo "Running security scans..."
	@command -v gosec >/dev/null 2>&1 || { echo "gosec not found. Installing..."; go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; }
	cd backend && gosec ./...
	cd frontend-nextjs && npm audit --audit-level moderate
	@echo "✅ Security scans completed"

# Linting and formatting
lint: backend-lint frontend-lint

format: backend-format frontend-format

# Database commands
db-migrate:
	@echo "Running database migrations..."
	cd backend && ./run_migrations.sh

db-backup:
	@echo "Backing up database..."
	cd backend && ./scripts/db-backup.sh

# Cleanup
clean:
	@echo "Cleaning build artifacts..."
	rm -f backend/nutrition-platform
	rm -rf frontend-nextjs/.next
	rm -rf frontend-nextjs/node_modules/.cache
	rm -f frontend-nextjs/.next/server/.next/static/chunks/pages/
	@echo "✅ Clean completed"

# Quick health checks
health:
	@echo "Running health checks..."
	@curl -f http://localhost:8080/health 2>/dev/null && echo "✅ Backend healthy" || echo "❌ Backend not responding"
	@curl -f http://localhost:3000 2>/dev/null && echo "✅ Frontend healthy" || echo "❌ Frontend not responding"

# Generate documentation
docs:
	@echo "Generating documentation..."
	cd backend && go run docs/generate_docs.go
	@echo "✅ Documentation generated"

# Performance testing
perf-test:
	@echo "Running performance tests..."
	@command -v wrk >/dev/null 2>&1 || { echo "wrk not found. Install with: brew install wrk"; exit 1; }
	wrk -t12 -c400 -d30s http://localhost:8080/api/v1/health

# Load testing
load-test:
	@echo "Running load tests..."
	@command -v hey >/dev/null 2>&1 || { echo "hey not found. Install with: go install github.com/rakyll/hey@latest"; exit 1; }
	hey -n 1000 -c 10 http://localhost:8080/api/v1/nutrition-data/recipes

# Database reset (development only)
db-reset:
	@echo "⚠️  This will reset the database. Press Ctrl+C to cancel within 5 seconds..."
	@sleep 5
	cd backend && rm -f nutrition.db && ./run_migrations.sh
	@echo "✅ Database reset"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/rakyll/hey@latest
	@echo "✅ Development tools installed"

# CI/CD helpers
ci-test:
	@echo "Running CI tests..."
	make test
	make lint
	make security
	@echo "✅ CI tests completed"

ci-build:
	@echo "Running CI build..."
	make clean
	make build
	make docker-build
	@echo "✅ CI build completed"

# Version management
version:
	@echo "Nutrition Platform Version Info:"
	@echo "Backend: $$(cd backend && go version | awk '{print $$3}')"
	@echo "Frontend: $$(cd frontend-nextjs && node --version)"
	@echo "Go: $$(cd backend && go version | awk '{print $$3}')"
	@echo "Node: $$(cd frontend-nextjs && node --version)"
	@echo "npm: $$(cd frontend-nextjs && npm --version)"

# Quick start (for new developers)
quick-start: setup-dev dev

# Production deployment pipeline
deploy-pipeline: ci-test ci-build deploy-staging
	@echo "✅ Deployment pipeline completed"
