#!/bin/bash
set -e

echo "🚀 PRODUCTION DEPLOYMENT - SECURE CREDENTIALS"
echo "=============================================="

# Generate secure credentials
DB_PASSWORD=$(openssl rand -hex 32)
JWT_SECRET=$(openssl rand -hex 64)
API_KEY_SECRET=$(openssl rand -hex 64)
REDIS_PASSWORD=$(openssl rand -hex 32)
ENCRYPTION_KEY=$(openssl rand -hex 32)
SESSION_SECRET=$(openssl rand -hex 32)

# Create .env.production
cat > .env.production << EOF
# Database
DB_HOST=postgres
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=nutrition_user
DB_PASSWORD=${DB_PASSWORD}
DB_SSL_MODE=require

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=${REDIS_PASSWORD}

# Security
JWT_SECRET=${JWT_SECRET}
API_KEY_SECRET=${API_KEY_SECRET}
ENCRYPTION_KEY=${ENCRYPTION_KEY}
SESSION_SECRET=${SESSION_SECRET}

# Server
PORT=8080
ENVIRONMENT=production
DOMAIN=super.doctorhealthy1.com
ALLOWED_ORIGINS=https://super.doctorhealthy1.com,https://www.super.doctorhealthy1.com

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60s
EOF

echo "✅ Credentials generated and saved to .env.production"
echo ""
echo "📋 SAVE THESE CREDENTIALS SECURELY:"
echo "===================================="
echo "DB_PASSWORD=${DB_PASSWORD}"
echo "REDIS_PASSWORD=${REDIS_PASSWORD}"
echo "JWT_SECRET=${JWT_SECRET}"
echo ""

# Deploy
echo "🐳 Starting Docker deployment..."
docker-compose -f docker-compose.production.yml up -d --build

echo ""
echo "✅ DEPLOYMENT COMPLETE!"
echo "🌐 Access: https://super.doctorhealthy1.com"
