# Developer Onboarding Guide

Welcome to the Nutrition Platform development team! This guide will help you get set up and productive quickly.

## Prerequisites

Before you start, make sure you have the following installed:

### Required Tools
- **Go 1.21+** - Backend development
- **Node.js 18+** - Frontend development
- **npm 9+** - Package management
- **Git** - Version control
- **VS Code** (recommended) - Code editor

### Databases & Services
- **PostgreSQL 13+** (optional for production)
- **SQLite** (for development)
- **Redis 6+** (optional for caching)

### Recommended Tools
- **Docker** - Containerization
- **Make** - Build automation
- **curl** - API testing
- **Postman** - API development

## Quick Start (5 minutes)

### 1. Clone Repository
```bash
git clone https://github.com/DrKhaled123/kiro-nutrition.git
cd nutrition-platform
```

### 2. One-Command Setup
```bash
make quick-start
```

This will:
- Install backend dependencies
- Install frontend dependencies  
- Create environment files
- Start both development servers

### 3. Verify Installation
Open your browser and navigate to:
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **API Health**: http://localhost:8080/health

## Detailed Setup

### Backend Setup

#### Environment Configuration
```bash
cd backend
cp .env.example .env
```

Edit `.env` with your configuration:
```env
# Database
DB_TYPE=sqlite
DB_PATH=./nutrition.db

# Server
PORT=8080
HOST=localhost

# Redis (optional)
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your-secret-key-here
JWT_EXPIRES_IN=24h

# External APIs (optional)
NUTRITION_API_KEY=your-api-key
WEATHER_API_KEY=your-weather-key
```

#### Dependencies & Database
```bash
# Install Go dependencies
go mod download
go mod verify

# Run database migrations
make db-migrate
```

#### Start Backend Server
```bash
# Development mode
make run-backend

# Or directly
go run main.go
```

### Frontend Setup

#### Environment Configuration
```bash
cd frontend-nextjs
cp .env.local.example .env.local
```

Edit `.env.local`:
```env
# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_APP_NAME=Nutrition Platform

# Feature Flags
NEXT_PUBLIC_ENABLE_ANALYTICS=false
NEXT_PUBLIC_ENABLE_CACHE=true

# Development
NEXT_PUBLIC_DEV_MODE=true
```

#### Dependencies
```bash
# Install dependencies
npm install

# Or with exact versions (recommended for CI)
npm ci
```

#### Start Frontend Server
```bash
# Development mode
make run-frontend

# Or directly
npm run dev
```

## Development Workflow

### Daily Development Commands

```bash
# Start both servers
make dev

# Run tests
make test

# Format code
make format

# Lint code
make lint

# Clean build artifacts
make clean
```

### Making Changes

1. **Create a feature branch**
```bash
git checkout -b feature/your-feature-name
```

2. **Make your changes**
   - Backend: Add/edit handlers, services, models
   - Frontend: Add/edit components, pages, hooks

3. **Test your changes**
```bash
# Backend tests
make test-backend

# Frontend tests
make test-frontend
```

4. **Commit your changes**
```bash
# Pre-commit hooks will run automatically
git add .
git commit -m "feat: add new nutrition calculator"
```

5. **Push and create PR**
```bash
git push origin feature/your-feature-name
# Create Pull Request on GitHub
```

## Project Structure

```
nutrition-platform/
├── backend/                    # Go backend API
│   ├── handlers/               # HTTP handlers
│   ├── services/               # Business logic
│   ├── repositories/           # Data access layer
│   ├── models/                 # Data models
│   ├── middleware/             # HTTP middleware
│   ├── utils/                 # Utility functions
│   ├── tests/                  # Test files
│   ├── migrations/             # Database migrations
│   └── main.go                 # Application entry
├── frontend-nextjs/            # Next.js frontend
│   ├── src/
│   │   ├── app/               # App router pages
│   │   ├── components/         # React components
│   │   ├── hooks/             # Custom hooks
│   │   ├── lib/               # Utilities and API
│   │   ├── types/             # TypeScript types
│   │   └── utils/             # Helper functions
│   ├── public/                # Static assets
│   └── package.json
├── docs/                      # Documentation
├── scripts/                    # Deployment and utility scripts
├── .github/workflows/          # CI/CD pipelines
└── Makefile                   # Build automation
```

## Common Tasks

### Adding a New API Endpoint

1. **Define the model** (`backend/models/`)
2. **Create repository** (`backend/repositories/`)
3. **Implement service** (`backend/services/`)
4. **Add handler** (`backend/handlers/`)
5. **Register route** (`backend/main.go`)
6. **Write tests** (`backend/tests/`)

### Adding a New Frontend Page

1. **Create page component** (`frontend-nextjs/src/app/`)
2. **Add routing** (automatic with App Router)
3. **Create components** (`frontend-nextjs/src/components/`)
4. **Add API integration** (`frontend-nextjs/src/lib/`)
5. **Write tests** (`frontend-nextjs/__tests__/`)

### Database Changes

1. **Create migration** (`backend/migrations/`)
2. **Update model** (`backend/models/`)
3. **Update repository** (`backend/repositories/`)
4. **Run migration**
```bash
make db-migrate
```

## Testing

### Backend Tests

```bash
# Run all tests
make test-backend

# Run with coverage
go test ./... -cover

# Run specific test
go test ./handlers -run TestGetRecipes

# Run integration tests
go test ./tests/integration/... -tags=integration

# Run contract tests
go test ./tests/contract/... -tags=contract
```

### Frontend Tests

```bash
# Run all tests
make test-frontend

# Run with coverage
npm test -- --coverage

# Run in watch mode
npm test -- --watch

# Run specific test file
npm test -- RecipePage.test.tsx
```

### API Testing

```bash
# Run API smoke tests
./scripts/smoke-test.sh

# Generate Postman collection
./scripts/generate-postman-collection.sh

# Manual API testing
curl http://localhost:8080/health
```

## Debugging

### Backend Debugging

```bash
# Enable debug mode
DEBUG=true go run main.go

# Use delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug main.go
```

### Frontend Debugging

- Use VS Code debugger with launch configurations
- Browser DevTools for React component inspection
- Network tab for API calls

### Common Issues

**Backend won't start**:
```bash
# Check port is available
lsof -ti:8080

# Check database file permissions
ls -la nutrition.db

# Check environment variables
cat .env
```

**Frontend build errors**:
```bash
# Clear Next.js cache
rm -rf .next

# Reinstall dependencies
rm -rf node_modules package-lock.json
npm install
```

## Code Standards

### Backend (Go)

```go
// Use proper error handling
if err != nil {
    return nil, fmt.Errorf("failed to process: %w", err)
}

// Follow Go conventions
type UserService struct {
    repo UserRepository
}

// Add context to functions
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    // implementation
}
```

### Frontend (TypeScript/React)

```typescript
// Use proper types
interface Recipe {
    id: string;
    name: string;
    calories: number;
}

// Use custom hooks
const { data, loading, error } = useRecipes();

// Follow React patterns
const RecipeComponent: React.FC<RecipeProps> = ({ recipe }) => {
    return <div>{recipe.name}</div>;
};
```

## Performance

### Backend Optimization

- Use database indexes
- Implement caching with Redis
- Use connection pooling
- Monitor memory usage

### Frontend Optimization

- Use React.memo for expensive components
- Implement code splitting
- Optimize images and assets
- Use lazy loading

## Security

### Backend Security

- Validate all inputs
- Use parameterized queries
- Implement rate limiting
- Secure JWT tokens

### Frontend Security

- Sanitize user inputs
- Use HTTPS in production
- Implement CSP headers
- Validate API responses

## Deployment

### Development Deployment
```bash
# Deploy to staging
make deploy-staging
```

### Production Deployment
```bash
# Deploy to production (with confirmation)
make deploy-prod

# Manual deployment with version
./scripts/deploy.sh production v1.2.0
```

## Getting Help

### Documentation
- **API Reference**: [API_REFERENCE.md](API_REFERENCE.md)
- **Development Setup**: [DEVELOPMENT_SETUP.md](DEVELOPMENT_SETUP.md)
- **Troubleshooting**: [TROUBLESHOOTING.md](TROUBLESHOOTING.md)

### Team Communication
- **Slack**: #nutrition-platform-dev
- **GitHub Issues**: For bugs and feature requests
- **Code Review**: All PRs require review

### Office Hours
- **Daily Standup**: 9:00 AM UTC
- **Weekly Planning**: Monday 2:00 PM UTC
- **Retrospective**: Friday 4:00 PM UTC

## Next Steps

1. Complete the Quick Start above
2. Explore the codebase
3. Make a small change (fix a typo, add a comment)
4. Submit your first PR
5. Join team communication channels

## Resources

### Learning Resources
- [Go Documentation](https://golang.org/doc/)
- [React Documentation](https://react.dev/)
- [Next.js Documentation](https://nextjs.org/docs)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)

### Tools
- **API Testing**: Postman, Insomnia
- **Database**: DBeaver, pgAdmin
- **Git**: GitKraken, SourceTree
- **Code Editor**: VS Code with extensions

### Extensions (VS Code)
- Go extension for Go development
- ES7+ React/Redux/React-Native snippets
- Prettier - Code formatter
- ESLint - JavaScript linter
- Thunder Client - API testing

---

Welcome aboard! We're excited to have you on the team. Don't hesitate to ask questions and contribute ideas!
