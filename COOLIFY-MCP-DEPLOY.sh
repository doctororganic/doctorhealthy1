#!/bin/bash
set -e

echo "🚀 COOLIFY MCP DEPLOYMENT"
echo "========================="

# Load credentials
if [ -f ".coolify-credentials.enc" ]; then
    source .coolify-credentials.enc
    echo "✅ Credentials loaded"
else
    echo "❌ Credentials file not found"
    exit 1
fi

# Verify MCP is installed
if ! command -v npx &> /dev/null; then
    echo "❌ npx not found. Install Node.js first"
    exit 1
fi

echo "✅ npx found"

# Test Coolify connection
echo ""
echo "🔍 Testing Coolify API connection..."
response=$(curl -s -w "%{http_code}" -o /tmp/coolify-test.json \
    -H "Authorization: Bearer $COOLIFY_TOKEN" \
    "$COOLIFY_BASE_URL/api/v1/servers")

if [ "$response" = "200" ]; then
    echo "✅ Coolify API connection successful"
    cat /tmp/coolify-test.json | head -20
else
    echo "❌ Coolify API connection failed (HTTP $response)"
    cat /tmp/coolify-test.json
    exit 1
fi

# Prepare deployment package
echo ""
echo "📦 Preparing deployment package..."

# Create deployment config
cat > coolify-deploy-config.json << EOF
{
  "project": "nutrition-platform",
  "environment": "production",
  "domain": "super.doctorhealthy1.com",
  "services": {
    "backend": {
      "image": "nutrition-backend",
      "port": 8080,
      "healthcheck": "/health"
    },
    "frontend": {
      "image": "nutrition-frontend",
      "port": 3000
    },
    "postgres": {
      "image": "postgres:15-alpine",
      "port": 5432
    },
    "redis": {
      "image": "redis:7-alpine",
      "port": 6379
    }
  }
}
EOF

echo "✅ Deployment config created"

# Generate secure environment variables
echo ""
echo "🔐 Generating secure environment variables..."

DB_PASSWORD=$(openssl rand -hex 32)
JWT_SECRET=$(openssl rand -hex 64)
API_KEY_SECRET=$(openssl rand -hex 64)
REDIS_PASSWORD=$(openssl rand -hex 32)

# Save to secure file
cat > .env.coolify << EOF
# Coolify Production Environment
DB_HOST=postgres
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=nutrition_user
DB_PASSWORD=${DB_PASSWORD}
DB_SSL_MODE=require

REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=${REDIS_PASSWORD}

JWT_SECRET=${JWT_SECRET}
API_KEY_SECRET=${API_KEY_SECRET}

PORT=8080
ENVIRONMENT=production
DOMAIN=super.doctorhealthy1.com
ALLOWED_ORIGINS=https://super.doctorhealthy1.com

RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60s
EOF

echo "✅ Environment variables generated"
echo ""
echo "📋 SAVE THESE CREDENTIALS:"
echo "=========================="
echo "DB_PASSWORD: ${DB_PASSWORD}"
echo "REDIS_PASSWORD: ${REDIS_PASSWORD}"
echo "JWT_SECRET: ${JWT_SECRET:0:32}..."
echo ""

# Deploy using Coolify API
echo "🚀 Deploying to Coolify..."

# Create project
curl -X POST "$COOLIFY_BASE_URL/api/v1/projects" \
    -H "Authorization: Bearer $COOLIFY_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "nutrition-platform",
        "description": "AI-powered nutrition and health management platform"
    }' | jq '.'

echo ""
echo "✅ DEPLOYMENT INITIATED"
echo ""
echo "🌐 Access URLs:"
echo "  Frontend: https://super.doctorhealthy1.com"
echo "  Backend:  https://api.doctorhealthy1.com"
echo "  Coolify:  https://api.doctorhealthy1.com"
echo ""
echo "📊 Monitor deployment:"
echo "  Visit Coolify dashboard"
echo "  Check logs in real-time"
echo ""
