#!/bin/bash
set -e

echo "🚀 DEPLOYING TO COOLIFY VIA MCP"
echo "================================"

# Load credentials
source .coolify-credentials.enc

# Server info from MCP
SERVER_UUID="x8gck8ggggsgkggg4coosg0g"
SERVER_IP="128.140.111.171"
TEAM_ID="0"

echo "✅ Coolify Version: 4.0.0-beta.434"
echo "✅ Server: localhost ($SERVER_IP)"
echo "✅ Team: Root Team"
echo ""

# Generate deployment secrets
echo "🔐 Generating secure credentials..."
DB_PASSWORD=$(openssl rand -hex 32)
JWT_SECRET=$(openssl rand -hex 64)
API_KEY_SECRET=$(openssl rand -hex 64)
REDIS_PASSWORD=$(openssl rand -hex 32)
ENCRYPTION_KEY=$(openssl rand -hex 32)

# Save credentials
cat > .env.coolify.secure << EOF
# COOLIFY DEPLOYMENT CREDENTIALS
# Generated: $(date)

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

# Server
PORT=8080
ENVIRONMENT=production
DOMAIN=super.doctorhealthy1.com
API_DOMAIN=api.super.doctorhealthy1.com
ALLOWED_ORIGINS=https://super.doctorhealthy1.com,https://www.super.doctorhealthy1.com

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60s
EOF

echo "✅ Credentials generated"
echo ""
echo "📋 SAVE THESE CREDENTIALS SECURELY:"
echo "===================================="
echo "DB_PASSWORD=${DB_PASSWORD}"
echo "REDIS_PASSWORD=${REDIS_PASSWORD}"
echo "JWT_SECRET=${JWT_SECRET:0:32}..."
echo ""

# Create deployment package
echo "📦 Creating deployment package..."

# Create docker-compose for Coolify
cat > docker-compose.coolify.yml << 'EOF'
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    command: postgres -c ssl=on
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 3s
      retries: 5

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile.secure
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_NAME=${DB_NAME}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_SSL_MODE=require
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - JWT_SECRET=${JWT_SECRET}
      - API_KEY_SECRET=${API_KEY_SECRET}
      - PORT=8080
      - ENVIRONMENT=production
      - DOMAIN=${DOMAIN}
      - ALLOWED_ORIGINS=${ALLOWED_ORIGINS}
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    ports:
      - "8080:8080"
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  frontend:
    build:
      context: ./frontend-nextjs
      dockerfile: Dockerfile.secure
    environment:
      - NEXT_PUBLIC_API_URL=https://${API_DOMAIN}
      - NODE_ENV=production
    depends_on:
      - backend
    ports:
      - "3000:3000"

volumes:
  postgres_data:
  redis_data:
EOF

echo "✅ Docker Compose created"

# Create Coolify configuration
cat > coolify.json << EOF
{
  "project": {
    "name": "nutrition-platform",
    "description": "AI-powered nutrition and health management platform"
  },
  "services": {
    "backend": {
      "type": "application",
      "build_pack": "dockerfile",
      "dockerfile_location": "backend/Dockerfile.secure",
      "ports_exposes": "8080",
      "health_check_enabled": true,
      "health_check_path": "/health",
      "health_check_port": "8080",
      "health_check_interval": 30,
      "health_check_timeout": 10,
      "health_check_retries": 3,
      "domains": ["api.super.doctorhealthy1.com"]
    },
    "frontend": {
      "type": "application",
      "build_pack": "dockerfile",
      "dockerfile_location": "frontend-nextjs/Dockerfile.secure",
      "ports_exposes": "3000",
      "domains": ["super.doctorhealthy1.com"]
    },
    "postgres": {
      "type": "database",
      "image": "postgres:15-alpine",
      "ports_exposes": "5432"
    },
    "redis": {
      "type": "database",
      "image": "redis:7-alpine",
      "ports_exposes": "6379"
    }
  }
}
EOF

echo "✅ Coolify config created"

# Deploy using Coolify API
echo ""
echo "🚀 Deploying to Coolify..."
echo ""

# Create project
echo "Creating project..."
project_response=$(curl -s -X POST "$COOLIFY_BASE_URL/api/v1/projects" \
    -H "Authorization: Bearer $COOLIFY_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "nutrition-platform",
        "description": "AI-powered nutrition and health management platform"
    }')

echo "$project_response" | jq '.' || echo "$project_response"

# Instructions for manual deployment
echo ""
echo "================================"
echo "✅ DEPLOYMENT PACKAGE READY"
echo "================================"
echo ""
echo "📁 Files created:"
echo "  - docker-compose.coolify.yml"
echo "  - coolify.json"
echo "  - .env.coolify.secure"
echo ""
echo "🌐 Next steps in Coolify Dashboard:"
echo ""
echo "1. Go to: https://api.doctorhealthy1.com"
echo "2. Create new project: 'nutrition-platform'"
echo "3. Add Git repository or upload files"
echo "4. Configure environment variables from .env.coolify.secure"
echo "5. Set domains:"
echo "   - Backend: api.super.doctorhealthy1.com"
echo "   - Frontend: super.doctorhealthy1.com"
echo "6. Deploy!"
echo ""
echo "📊 Monitor deployment:"
echo "  - Check build logs"
echo "  - Verify health checks"
echo "  - Test endpoints"
echo ""
echo "🎯 Access URLs (after deployment):"
echo "  Frontend: https://super.doctorhealthy1.com"
echo "  Backend:  https://api.super.doctorhealthy1.com"
echo "  Health:   https://api.super.doctorhealthy1.com/health"
echo ""
echo "🔐 Credentials saved in: .env.coolify.secure"
echo ""
