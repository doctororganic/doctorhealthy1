# 🍎 Nutrition Platform

AI-powered nutrition and health management platform.

## Quick Start

```bash
# Start everything
docker-compose up -d

# Or use deployment script
./deploy.sh
```

## Services

- **Backend:** Go API on port 8080
- **Frontend:** Next.js on port 3000
- **Database:** PostgreSQL on port 5432
- **Cache:** Redis on port 6379

## Development

### Backend (Go)
```bash
cd backend
go run main.go
```

### Frontend (Next.js)
```bash
cd frontend-nextjs
npm install
npm run dev
```

## API Endpoints

- `GET /health` - Health check
- `GET /api/v1/info` - API information
- `POST /api/v1/nutrition/analyze` - Nutrition analysis
- `GET /api/v1/recipes` - Recipe management
- `GET /api/v1/workouts` - Workout plans
- `POST /api/v1/generate-meal-plan` - Meal plan generation

## Documentation

- [Backend README](backend/README.md)
- [API Documentation](backend/docs/)
- [Deployment Guide](DEPLOYMENT.md)

## Architecture

```
┌─────────────┐      ┌─────────────┐
│   Next.js   │─────▶│   Go API    │
│  Frontend   │      │   Backend   │
└─────────────┘      └─────────────┘
                            │
                     ┌──────┴──────┐
                     ▼              ▼
              ┌───────────┐  ┌──────────┐
              │ PostgreSQL│  │  Redis   │
              └───────────┘  └──────────┘
```

## License

MIT
