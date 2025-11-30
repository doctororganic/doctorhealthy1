# 🚀 Quick Deployment Guide

## For Server Deployment

### Option 1: Using Docker Compose (Recommended)

```bash
# 1. Clone repository
git clone https://github.com/doctororganic/doctorhealthy1.git
cd doctorhealthy1/nutrition-platform

# 2. Copy environment template
cp .env.production.example .env.production

# 3. Edit .env.production with your settings
nano .env.production

# 4. Run deployment script
chmod +x deploy.sh
./deploy.sh
```

### Option 2: Manual Docker Compose

```bash
# Build and start
docker-compose -f docker-compose.production.yml up -d --build

# Check status
docker-compose -f docker-compose.production.yml ps

# View logs
docker-compose -f docker-compose.production.yml logs -f
```

## Services Included

- **Backend**: Go API server (port 8080)
- **Frontend**: Next.js application (port 3000)
- **PostgreSQL**: Database (internal)
- **Redis**: Cache/Session store (internal)
- **Nginx**: Reverse proxy (ports 80, 443)

## Access Points

- Frontend: http://your-server:3000
- Backend API: http://your-server:8080
- Health Check: http://your-server:8080/health

## Important Files

- `docker-compose.production.yml` - Main deployment configuration
- `.env.production` - Environment variables (create from .env.production.example)
- `deploy.sh` - Automated deployment script
- `SERVER_DEPLOYMENT_GUIDE.md` - Detailed deployment guide

## Quick Commands

```bash
# Start all services
docker-compose -f docker-compose.production.yml up -d

# Stop all services
docker-compose -f docker-compose.production.yml down

# Restart a service
docker-compose -f docker-compose.production.yml restart backend

# View logs
docker-compose -f docker-compose.production.yml logs -f backend

# Update application
git pull
docker-compose -f docker-compose.production.yml up -d --build
```

## Troubleshooting

See `SERVER_DEPLOYMENT_GUIDE.md` for detailed troubleshooting steps.

