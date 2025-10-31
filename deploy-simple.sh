#!/bin/bash

# Simple Deployment Script
# Deploys the nutrition platform using Docker

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
DOMAIN="super.doctorhealthy1.com"
NGINX_CONF="./nginx/conf.d/default.conf"

# Function to print colored output
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log "🚀 Starting Simple Docker Deployment for Nutrition Platform"
echo ""

# Step 1: Build Docker image
log "📦 Building Docker image..."
if docker build -t nutrition-platform:latest .; then
    success "✅ Docker image built successfully"
else
    error "❌ Failed to build Docker image"
fi

# Step 2: Stop any existing container
log "🛑 Stopping any existing container..."
docker stop nutrition-platform 2>/dev/null || true
docker rm nutrition-platform 2>/dev/null || true

# Step 3: Run the container
log "🚀 Starting container..."
docker run -d \
    --name nutrition-platform \
    -p 80:80 \
    -p 443:443 \
    -e DOMAIN="$DOMAIN" \
    -e SERVER_PORT=80 \
    -e DB_HOST=localhost \
    -e DB_SSL_MODE=disable \
    -e CORS_ALLOWED_ORIGINS="https://$DOMAIN,http://localhost" \
    -v $(pwd)/nginx/conf.d:/etc/nginx/conf.d \
    nutrition-platform:latest

if [ $? -eq 0 ]; then
    success "✅ Container started successfully"
else
    error "❌ Failed to start container"
fi

# Step 4: Wait for container to be ready
log "⏳ Waiting for container to be ready..."
sleep 10

# Step 5: Check if container is running
if docker ps | grep -q nutrition-platform; then
    success "✅ Container is running"
else
    error "❌ Container is not running"
fi

# Step 6: Health check
log "🏥 Performing health check..."
sleep 5

# Check HTTP
if curl -f -s http://localhost/health > /dev/null; then
    success "✅ HTTP health check passed"
else
    warning "⚠️ HTTP health check failed"
fi

# Step 7: Display access information
echo ""
echo "🎉 ==================================="
echo "🎉 DEPLOYMENT COMPLETED!"
echo "🎉 ==================================="
echo ""
echo "📍 Local Access:"
echo "   🌐 HTTP: http://localhost"
echo "   🔒 HTTPS: http://localhost (SSL not configured locally)"
echo ""
echo "📊 Container Status:"
docker ps | grep nutrition-platform
echo ""
echo "📋 Container Logs:"
echo "   docker logs nutrition-platform"
echo ""
echo "🔧 Management:"
echo "   Stop: docker stop nutrition-platform"
echo "   Restart: docker restart nutrition-platform"
echo "   Remove: docker rm nutrition-platform"
echo ""

success "🚀 Nutrition Platform is now running locally!"
echo ""
echo "⚠️ Note: For production deployment with SSL and domain, configure"
echo "   1. DNS for $DOMAIN to point to this server"
echo "   2. SSL certificates (Let's Encrypt recommended)"
echo "   3. Firewall rules for ports 80 and 443"