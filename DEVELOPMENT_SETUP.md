# Development Setup Guide

## Overview

This guide will help you set up the DoctorHealthy nutrition platform for local development. The platform consists of a Go backend and a Next.js frontend.

## System Requirements

### Prerequisites
- **Git**: For version control
- **Go 1.21+**: Backend development
- **Node.js 18+**: Frontend development
- **npm or pnpm**: Package management
- **SQLite 3**: Default database (development)
- **Redis 6+**: Caching (optional but recommended)
- **Docker**: Optional, for containerized development

### Development Tools (Recommended)
- **VS Code**: Preferred IDE with extensions:
  - Go extension (golang.go)
  - Tailwind CSS IntelliSense
  - TypeScript Hero
  - GitLens
- **Postman** or **Insomnia**: API testing
- **DBeaver** or **TablePlus**: Database management

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/DrKhaled123/kiro-nutrition.git
cd kiro-nutrition
```

### 2. Backend Setup

```bash
cd nutrition-platform/backend

# Install Go dependencies
go mod download

# Copy environment configuration
cp .env.example .env

# Edit environment variables
nano .env
```

### 3. Frontend Setup

```bash
cd nutrition-platform/frontend-nextjs

# Install dependencies
npm install
# or
pnpm install

# Copy environment configuration
cp .env.local.example .env.local

# Edit environment variables
nano .env.local
```

### 4. Database Setup

```bash
cd nutrition-platform/backend

# Run database migrations
make migrate
# or
go run main.go migrate

# Seed data (optional)
make seed
# or
go run main.go seed
```

### 5. Start Development Servers

**Backend:**
```bash
cd nutrition-platform/backend
make run
# or
go run main.go
```

**Frontend:**
```bash
cd nutrition-platform/frontend-nextjs
npm run dev
# or
pnpm dev
```

## Detailed Setup

### Backend Configuration

#### Environment Variables (.env)

```bash
# Server Configuration
PORT=8080
HOST=localhost
ENV=development

# Database Configuration
DB_TYPE=sqlite
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=doctorhealthy_dev
DB_PATH=./data/doctorhealthy.db

# Redis Configuration (Optional)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRES_IN=24h

# API Configuration
API_BASE_URL=http://localhost:8080
CORS_ORIGINS=http://localhost:3000,http://localhost:3001

# File Upload Configuration
UPLOAD_DIR=./uploads
MAX_FILE_SIZE=10MB

# Logging Configuration
LOG_LEVEL=debug
LOG_FORMAT=json

# External Services
NUTRITION_DATA_PATH=../nutrition data json
```

#### Database Setup

**SQLite (Default - Development):**
```bash
# Create data directory
mkdir -p data

# Run migrations
make migrate

# Check database
sqlite3 data/doctorhealthy.db ".tables"
```

**PostgreSQL (Optional - Production-like):**
```bash
# Install PostgreSQL
brew install postgresql  # macOS
sudo apt install postgresql  # Ubuntu

# Start PostgreSQL
brew services start postgresql
sudo systemctl start postgresql

# Create database
createdb doctorhealthy_dev

# Update .env
DB_TYPE=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=doctorhealthy_dev
```

#### Redis Setup (Optional but Recommended)

```bash
# Install Redis
brew install redis  # macOS
sudo apt install redis-server  # Ubuntu

# Start Redis
brew services start redis
sudo systemctl start redis

# Test connection
redis-cli ping
```

### Frontend Configuration

#### Environment Variables (.env.local)

```bash
# API Configuration
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080

# Feature Flags
NEXT_PUBLIC_ENABLE_ANALYTICS=false
NEXT_PUBLIC_ENABLE_DEBUG=true

# Environment
NODE_ENV=development
```

#### TypeScript Configuration

The project uses strict TypeScript configuration:

```json
{
  "compilerOptions": {
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "exactOptionalPropertyTypes": true
  }
}
```

## Development Workflow

### 1. Code Structure

```
nutrition-platform/
├── backend/                    # Go backend
│   ├── handlers/               # HTTP handlers
│   ├── models/                 # Data models
│   ├── repositories/           # Data access layer
│   ├── services/              # Business logic
│   ├── middleware/            # Echo middleware
│   ├── utils/                 # Utilities
│   ├── migrations/            # Database migrations
│   ├── tests/                 # Tests
│   └── main.go               # Application entry
├── frontend-nextjs/           # Next.js frontend
│   ├── src/
│   │   ├── app/              # App Router (Next.js 13+)
│   │   ├── components/       # React components
│   │   ├── hooks/            # Custom hooks
│   │   ├── lib/              # Utilities
│   │   ├── types/            # TypeScript types
│   │   └── utils/            # Helper functions
│   ├── public/               # Static assets
│   └── tests/                # Tests
└── nutrition data json/       # Sample data files
```

### 2. Available Scripts

**Backend Scripts:**
```bash
# Development
make run              # Start development server
make build            # Build for production
make test             # Run tests
make test-coverage    # Run tests with coverage
make lint             # Run linter
make fmt              # Format code
make migrate          # Run database migrations
make seed             # Seed database with sample data
make clean            # Clean build artifacts

# Database
make db-backup        # Backup database
make db-restore       # Restore database
make db-reset         # Reset database

# Testing
make test-integration # Run integration tests
make test-unit        # Run unit tests
make test-api         # Run API tests
```

**Frontend Scripts:**
```bash
# Development
npm run dev           # Start development server
npm run build         # Build for production
npm run start         # Start production server
npm run test          # Run tests
npm run test:watch    # Run tests in watch mode
npm run test:coverage # Run tests with coverage
npm run lint          # Run linter
npm run lint:fix      # Fix linting issues
npm run type-check    # Type checking
```

### 3. Database Migrations

**Creating a New Migration:**
```bash
cd nutrition-platform/backend

# Create migration file
make migration-create name=add_new_table
# or
goose -dir "./migrations" create add_new_table sql

# Edit migration file
# migrations/001_add_new_table.sql
```

**Migration File Format:**
```sql
-- +goose Up
CREATE TABLE new_table (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE new_table;
```

**Running Migrations:**
```bash
make migrate-up     # Run all pending migrations
make migrate-down   # Rollback last migration
make migrate-status  # Check migration status
```

### 4. Testing

**Backend Testing:**
```bash
# Run all tests
make test

# Run specific test types
make test-unit          # Unit tests only
make test-integration    # Integration tests only
make test-api          # API tests only

# Run with coverage
make test-coverage

# Run specific test file
go test ./tests/integration -v

# Run specific test function
go test ./tests/integration -run TestNutritionDataHandler_GetRecipes -v
```

**Frontend Testing:**
```bash
# Run all tests
npm test

# Run tests in watch mode
npm run test:watch

# Run tests with coverage
npm run test:coverage

# Run specific test file
npm test -- RecipeCard.test.tsx

# Run tests matching pattern
npm test -- --testNamePattern="should render"
```

### 5. API Development

**Adding New Endpoint:**
1. Create handler in `handlers/`
2. Register route in `main.go`
3. Add tests in `tests/`
4. Update API documentation

**Example Handler:**
```go
// handlers/new_handler.go
package handlers

import (
    "net/http"
    "github.com/labstack/echo/v4"
)

type NewHandler struct {
    // Dependencies here
}

func NewNewHandler() *NewHandler {
    return &NewHandler{}
}

func (h *NewHandler) GetItems(c echo.Context) error {
    // Handler logic here
    return c.JSON(http.StatusOK, map[string]interface{}{
        "status": "success",
        "data": []string{"item1", "item2"},
    })
}
```

**Register Route:**
```go
// main.go
newHandler := handlers.NewNewHandler()

// API routes
api := e.Group("/api/v1")
api.GET("/items", newHandler.GetItems)
```

### 6. Frontend Development

**Creating New Component:**
```bash
# Use component generator
cd frontend-nextjs
npm run generate-component ComponentName

# Or create manually
mkdir -p src/components/new-feature
touch src/components/new-feature/ComponentName.tsx
touch src/components/new-feature/ComponentName.test.tsx
touch src/components/new-feature/index.ts
```

**Component Template:**
```tsx
// src/components/new-feature/ComponentName.tsx
'use client';

interface ComponentNameProps {
  title: string;
  onAction?: () => void;
}

export function ComponentName({ title, onAction }: ComponentNameProps) {
  return (
    <div className="p-4 border rounded-lg">
      <h2 className="text-lg font-semibold">{title}</h2>
      {onAction && (
        <button
          onClick={onAction}
          className="mt-2 px-4 py-2 bg-blue-500 text-white rounded"
        >
          Action
        </button>
      )}
    </div>
  );
}

export default ComponentName;
```

## Code Quality

### 1. Linting and Formatting

**Backend:**
```bash
# Format code
make fmt
# or
go fmt ./...

# Lint code
make lint
# or
golangci-lint run

# Fix issues
golangci-lint run --fix
```

**Frontend:**
```bash
# Lint code
npm run lint

# Fix linting issues
npm run lint:fix

# Type checking
npm run type-check

# Format code (with Prettier)
npx prettier --write .
```

### 2. Pre-commit Hooks

Install pre-commit hooks for code quality:

```bash
# Backend
cd nutrition-platform/backend
go install github.com/pre-commit/pre-commit@latest
pre-commit install

# Frontend
cd frontend-nextjs
npm install husky lint-staged --save-dev
npx husky install
```

### 3. Code Standards

**Go Guidelines:**
- Follow Go conventions and idioms
- Use meaningful variable names
- Add comments for complex logic
- Handle errors properly
- Write unit tests for all functions

**TypeScript Guidelines:**
- Use strict TypeScript
- Define interfaces for all props
- Use proper typing for API responses
- Avoid `any` type
- Use generic types where appropriate

## Debugging

### 1. Backend Debugging

**Using Delve Debugger:**
```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug with breakpoints
dlv debug main.go

# Debug tests
dlv test ./tests/integration
```

**Logging:**
```bash
# Enable debug logging
export LOG_LEVEL=debug

# View logs in development
make run 2>&1 | tee -a app.log

# Search logs
grep "ERROR" app.log
grep "recipe.*not.*found" app.log
```

### 2. Frontend Debugging

**Browser DevTools:**
- Use React DevTools for component debugging
- Use Redux DevTools for state debugging
- Check Network tab for API calls
- Use Console for error messages

**VS Code Debugging:**
```json
// .vscode/launch.json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Next.js: debug server-side",
      "type": "node-terminal",
      "request": "launch",
      "command": "npm run dev"
    },
    {
      "name": "Next.js: debug client-side",
      "type": "chrome",
      "request": "launch",
      "url": "http://localhost:3000"
    }
  ]
}
```

## Performance Optimization

### 1. Backend Optimization

- Use connection pooling for database
- Implement Redis caching
- Use middleware for compression
- Monitor memory usage
- Profile with `pprof`

### 2. Frontend Optimization

- Use React.memo for components
- Implement code splitting
- Optimize images and assets
- Use Next.js Image component
- Monitor bundle size

## Common Issues and Solutions

### Backend Issues

**Port Already in Use:**
```bash
# Kill process on port 8080
./scripts/kill-port.sh 8080
# or
lsof -ti:8080 | xargs kill -9
```

**Database Connection Issues:**
```bash
# Check SQLite file
ls -la data/doctorhealthy.db

# Test connection
sqlite3 data/doctorhealthy.db ".tables"

# Reset database
make db-reset
```

### Frontend Issues

**Module Resolution:**
```bash
# Clear node_modules
rm -rf node_modules package-lock.json
npm install

# Clear Next.js cache
rm -rf .next
npm run dev
```

**TypeScript Errors:**
```bash
# Rebuild types
npm run type-check

# Clear TypeScript cache
rm -rf .tsbuildinfo
```

## Deployment

### 1. Backend Deployment

```bash
# Build binary
make build

# Run with production config
./bin/doctorhealthy -env production
```

### 2. Frontend Deployment

```bash
# Build for production
npm run build

# Start production server
npm run start
```

## Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/new-feature`
3. Make changes and test
4. Commit changes: `git commit -m "Add new feature"`
5. Push to branch: `git push origin feature/new-feature`
6. Create Pull Request

## Support

- **Documentation**: Check `/docs` directory
- **API Documentation**: http://localhost:8080/docs
- **Health Check**: http://localhost:8080/health
- **Issues**: Create GitHub issue
- **Discussions**: GitHub Discussions
