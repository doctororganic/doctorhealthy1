# Deployment Guide

## Local Development

```bash
docker-compose up -d
```

Visit:
- Frontend: http://localhost:3000
- Backend: http://localhost:8080
- Health: http://localhost:8080/health

## Production Deployment

### Option 1: Docker Compose (Recommended)

```bash
# On your server
git clone <repo>
cd nutrition-platform
./deploy.sh
```

### Option 2: Coolify

1. Create new project in Coolify
2. Connect Git repository
3. Set environment variables
4. Deploy

### Option 3: Manual

```bash
# Backend
cd backend
go build -o nutrition-platform
./nutrition-platform

# Frontend
cd frontend-nextjs
npm install
npm run build
npm start
```

## Environment Variables

### Backend
```
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=nutrition_user
DB_PASSWORD=your_password
REDIS_HOST=localhost
REDIS_PORT=6379
```

### Frontend
```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Monitoring

- Health: `curl http://localhost:8080/health`
- Metrics: `curl http://localhost:8080/metrics`
- Logs: `docker-compose logs -f`

## Troubleshooting

### Backend won't start
```bash
cd backend
go build
# Check for errors
```

### Frontend won't connect
- Check NEXT_PUBLIC_API_URL
- Verify backend is running
- Check CORS settings

### Database issues
```bash
docker-compose down -v
docker-compose up -d
```
