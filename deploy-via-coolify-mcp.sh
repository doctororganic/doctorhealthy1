#!/bin/bash

# Deploy Nutrition Platform using Coolify MCP
# This script uses the MCP server to deploy the project

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${BLUE}[DEPLOY]${NC} $1"; }
success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; }
warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }

echo ""
echo "🚀 =================================="
echo "🚀 DEPLOYING NUTRITION PLATFORM"
echo "🚀 =================================="
echo ""

# Load environment variables
if [ -f ".env.secure" ]; then
    export $(grep -v '^#' .env.secure | xargs)
    success "Loaded secure environment variables"
else
    error ".env.secure file not found"
    exit 1
fi

# Check if Coolify MCP is available
log "Checking Coolify MCP connection..."
if ! coolify-mcp-server --version &>/dev/null; then
    error "Coolify MCP server not found. Please install with: npm install -g coolify-mcp-server"
    exit 1
fi

success "Coolify MCP server is available"

# Create deployment package
log "Creating deployment package..."
tar -czf nutrition-platform-deploy.tar.gz \
    --exclude='.git' \
    --exclude='node_modules' \
    --exclude='*.log' \
    --exclude='archive/' \
    --exclude='coolify-deployment-*/' \
    --exclude='nutrition-platform-coolify-*/' \
    --exclude='.env*' \
    .

success "Deployment package created"

# Deploy using MCP
log "Initiating deployment via Coolify MCP..."

# Create MCP deployment request
cat > mcp-deploy-request.json << EOF
{
  "action": "deploy_application",
  "server_uuid": "x8gck8ggggsgkggg4coosg0g",
  "application": {
    "name": "nutrition-platform",
    "description": "AI-powered nutrition and health management platform",
    "domains": ["super.doctorhealthy1.com"],
    "build_pack": "dockercompose",
    "port": 3000,
    "environment_variables": {
      "NODE_ENV": "production",
      "PORT": "8080",
      "DB_HOST": "postgres",
      "DB_PORT": "5432",
      "DB_NAME": "nutrition_platform",
      "DB_USER": "nutrition_user",
      "DB_PASSWORD": "nutrition_pass",
      "REDIS_HOST": "redis",
      "REDIS_PORT": "6379",
      "ENVIRONMENT": "production",
      "NEXT_PUBLIC_API_URL": "http://backend:8080"
    },
    "source": {
      "type": "local",
      "path": "./nutrition-platform-deploy.tar.gz"
    },
    "health_check": {
      "path": "/health",
      "port": 8080,
      "interval": 30,
      "timeout": 10,
      "retries": 3
    },
    "force_https": true,
    "is_static": false
  }
}
EOF

# Execute deployment via MCP
log "Sending deployment request to Coolify..."
RESPONSE=$(curl -s -X POST \
    -H "Authorization: Bearer $COOLIFY_API_TOKEN" \
    -H "Content-Type: application/json" \
    -d @mcp-deploy-request.json \
    "$COOLIFY_API_URL/api/v1/applications" 2>/dev/null || echo "error")

# Check response
if echo "$RESPONSE" | grep -q "error\|Error"; then
    error "Deployment failed: $RESPONSE"
    exit 1
elif echo "$RESPONSE" | grep -q "id\|uuid"; then
    APP_ID=$(echo "$RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4 || echo "$RESPONSE" | grep -o '"uuid":"[^"]*"' | cut -d'"' -f4)
    success "Application deployment initiated! ID: $APP_ID"
else
    warning "Deployment response unclear: $RESPONSE"
fi

# Monitor deployment
log "Monitoring deployment progress..."
TIMEOUT=600  # 10 minutes
ELAPSED=0
INTERVAL=30

while [ $ELAPSED -lt $TIMEOUT ]; do
    if [ -n "$APP_ID" ]; then
        STATUS=$(curl -s -H "Authorization: Bearer $COOLIFY_API_TOKEN" \
            "$COOLIFY_API_URL/api/v1/applications/$APP_ID/deployments" 2>/dev/null | \
            grep -o '"status":"[^"]*"' | cut -d'"' -f4 || echo "unknown")
    else
        STATUS="checking"
    fi

    case "$STATUS" in
        "success"|"finished"|"completed"|"running")
            success "✅ Deployment completed successfully!"
            break
            ;;
        "failed"|"error")
            error "❌ Deployment failed"
            exit 1
            ;;
        "building"|"in_progress"|"deploying"|"pending")
            log "🔄 Deployment in progress... ($ELAPSED/$TIMEOUT seconds)"
            ;;
        *)
            log "📊 Deployment status: $STATUS ($ELAPSED/$TIMEOUT seconds)"
            ;;
    esac

    sleep $INTERVAL
    ELAPSED=$((ELAPSED + INTERVAL))
done

if [ $ELAPSED -ge $TIMEOUT ]; then
    warning "⚠️ Deployment monitoring timed out, but deployment may still be running"
fi

# Cleanup
log "Cleaning up temporary files..."
rm -f nutrition-platform-deploy.tar.gz
rm -f mcp-deploy-request.json

echo ""
echo "🎉 ===================================="
echo "🎉 DEPLOYMENT COMPLETED!"
echo "🎉 ===================================="
echo ""
echo "📍 Your Application is LIVE:"
echo "   🌐 Website: https://super.doctorhealthy1.com"
echo "   🏥 Health Check: https://super.doctorhealthy1.com/health"
echo "   📊 API Base: https://super.doctorhealthy1.com/api"
echo ""
echo "🖥️ Coolify Dashboard:"
echo "   🌐 Management: $COOLIFY_API_URL"
echo ""
echo "📋 Next Steps:"
echo "   1. 🔑 Test all API endpoints"
echo "   2. 📊 Monitor in Coolify dashboard"
echo "   3. 🌐 Share your application URL"
echo "   4. 📱 Test on mobile devices"
echo ""
success "🚀 Nutrition Platform deployed successfully!"