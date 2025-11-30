#!/bin/bash
# Complete Deployment Script for Nutrition Platform

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}🚀 Nutrition Platform - Production Deployment${NC}"
echo "=============================================="

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker is not installed. Please install Docker first.${NC}"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ Docker Compose is not installed. Please install Docker Compose first.${NC}"
    exit 1
fi

# Check if .env.production exists
if [ ! -f .env.production ]; then
    echo -e "${YELLOW}⚠️  .env.production not found. Creating from template...${NC}"
    
    # Generate secure passwords
    DB_PASSWORD=$(openssl rand -hex 32)
    JWT_SECRET=$(openssl rand -hex 32)
    API_KEY_SECRET=$(openssl rand -hex 32)
    SESSION_SECRET=$(openssl rand -hex 32)
    
    cat > .env.production << EOF
# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=postgres
DB_PASSWORD=${DB_PASSWORD}

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=

# Server Configuration
PORT=8080
ENV=production
DOMAIN=localhost

# Security
JWT_SECRET=${JWT_SECRET}
API_KEY_SECRET=${API_KEY_SECRET}
SESSION_SECRET=${SESSION_SECRET}

# CORS
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
EOF
    
    echo -e "${GREEN}✅ Created .env.production${NC}"
    echo -e "${YELLOW}⚠️  Please edit .env.production and update DOMAIN and ALLOWED_ORIGINS${NC}"
fi

# Build images
echo -e "${GREEN}📦 Building Docker images...${NC}"
docker-compose -f docker-compose.production.yml build

# Start services
echo -e "${GREEN}🚀 Starting services...${NC}"
docker-compose -f docker-compose.production.yml up -d

# Wait for services to be ready
echo -e "${GREEN}⏳ Waiting for services to start...${NC}"
sleep 15

# Health checks
echo -e "${GREEN}🏥 Running health checks...${NC}"

# Check backend
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Backend is healthy${NC}"
else
    echo -e "${YELLOW}⚠️  Backend health check failed (may need more time)${NC}"
fi

# Check frontend
if curl -f http://localhost:3000 > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Frontend is healthy${NC}"
else
    echo -e "${YELLOW}⚠️  Frontend health check failed (may need more time)${NC}"
fi

# Show status
echo -e "${GREEN}📊 Service Status:${NC}"
docker-compose -f docker-compose.production.yml ps

echo ""
echo -e "${GREEN}✅ Deployment complete!${NC}"
echo ""
echo "Access your application:"
echo "  Frontend: http://localhost:3000"
echo "  Backend API: http://localhost:8080"
echo "  Health Check: http://localhost:8080/health"
echo ""
echo "View logs: docker-compose -f docker-compose.production.yml logs -f"
echo "Stop services: docker-compose -f docker-compose.production.yml down"

