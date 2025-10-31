#!/bin/bash
set -e

echo "🚀 Deploying to Coolify..."

# Load environment variables
set -a
source .env.production
set +a

# Build and start services
docker-compose -f docker-compose.production.yml up -d --build

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 30

# Health check
echo "🔍 Checking application health..."
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Application is healthy!"
    echo "🌐 Access your application at: https://super.doctorhealthy1.com"
else
    echo "❌ Health check failed. Check the logs."
    exit 1
fi
