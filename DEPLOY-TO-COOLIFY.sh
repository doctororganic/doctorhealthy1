#!/bin/bash
set -e

echo "🚀 COOLIFY PRODUCTION DEPLOYMENT - SECURE EDITION"
echo "================================================"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}📋 COOLIFY DEPLOYMENT CHECKLIST${NC}"
echo "================================="

# Check if .env.production exists
if [ ! -f ".env.production" ]; then
    echo -e "${YELLOW}⚠️  .env.production not found. Creating with secure credentials...${NC}"
    # Generate secure credentials
    DB_PASSWORD=$(openssl rand -hex 32)
    JWT_SECRET=$(openssl rand -hex 64)
    API_KEY_SECRET=$(openssl rand -hex 64)
    REDIS_PASSWORD=$(openssl rand -hex 32)
    ENCRYPTION_KEY=$(openssl rand -hex 32)

    # Create .env.production
    cat > .env.production << EOF
# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=nutrition_user
DB_PASSWORD=${DB_PASSWORD}
DB_SSL_MODE=require

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=${REDIS_PASSWORD}

# Security Configuration
JWT_SECRET=${JWT_SECRET}
API_KEY_SECRET=${API_KEY_SECRET}
ENCRYPTION_KEY=${ENCRYPTION_KEY}

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
ENVIRONMENT=production
DEBUG=false

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://super.doctorhealthy1.com,https://www.super.doctorhealthy1.com

# Performance & Monitoring
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60s
SECURITY_HEADERS_ENABLED=true
COMPRESSION_ENABLED=true
METRICS_ENABLED=true
LOG_LEVEL=info
LOG_FORMAT=json

# Feature Flags
RELIGIOUS_FILTER_ENABLED=true
FILTER_ALCOHOL=true
FILTER_PORK=true
DEFAULT_LANGUAGE=en
SUPPORTED_LANGUAGES=en,ar
HEALTH_CHECK_ENABLED=true
EOF

    echo -e "${GREEN}✅ Secure credentials generated and saved to .env.production${NC}"
    echo ""
    echo -e "${YELLOW}📋 IMPORTANT: SAVE THESE CREDENTIALS:${NC}"
    echo "====================================="
    echo "DB_PASSWORD=${DB_PASSWORD}"
    echo "REDIS_PASSWORD=${REDIS_PASSWORD}"
    echo "JWT_SECRET=${JWT_SECRET}"
    echo ""
else
    echo -e "${GREEN}✅ .env.production already exists${NC}"
fi

# Create deployment package
echo -e "${BLUE}📦 Creating deployment package...${NC}"

# Create coolify deployment directory
COOLIFY_DIR="coolify-deployment-$(date +%Y%m%d-%H%M%S)"
mkdir -p "$COOLIFY_DIR"

# Copy necessary files
cp -r backend "$COOLIFY_DIR/"
cp -r frontend "$COOLIFY_DIR/"
cp -r nginx "$COOLIFY_DIR/"
cp docker-compose.production.yml "$COOLIFY_DIR/"
cp .env.production "$COOLIFY_DIR/"
cp coolify-step-by-step.md "$COOLIFY_DIR/"

# Create Coolify-specific files
cat > "$COOLIFY_DIR/coolify-config.json" << EOF
{
  "name": "nutrition-platform-secure",
  "domain": "super.doctorhealthy1.com",
  "buildCommand": "docker-compose -f docker-compose.production.yml build",
  "startCommand": "docker-compose -f docker-compose.production.yml up -d",
  "environmentVariables": [
    "ENVIRONMENT=production",
    "DEBUG=false"
  ],
  "healthCheckPath": "/health",
  "port": 8080
}
EOF

# Create deployment script for Coolify
cat > "$COOLIFY_DIR/deploy.sh" << 'EOF'
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
EOF

chmod +x "$COOLIFY_DIR/deploy.sh"

# Create ZIP file for Coolify upload
echo -e "${BLUE}📦 Creating ZIP package for Coolify...${NC}"
ZIP_FILE="../nutrition-platform-coolify-$(date +%Y%m%d-%H%M%S).zip"
cd "$COOLIFY_DIR"
zip -r "$ZIP_FILE" .
cd ..

echo -e "${GREEN}✅ Deployment package created: $ZIP_FILE${NC}"
echo ""
echo -e "${BLUE}🎯 COOLIFY DEPLOYMENT INSTRUCTIONS${NC}"
echo "=================================="
echo ""
echo -e "${YELLOW}📋 MANUAL DEPLOYMENT STEPS:${NC}"
echo ""
echo "1. 🌐 Go to your Coolify dashboard:"
echo "   https://api.doctorhealthy1.com"
echo ""
echo "2. 📁 Upload the deployment package:"
echo "   - Go to 'Applications' → 'Add Application'"
echo "   - Upload file: $ZIP_FILE"
echo "   - Application Name: nutrition-platform-secure"
echo ""
echo "3. ⚙️  Configure Application Settings:"
echo "   - Build Pack: Dockerfile"
echo "   - Dockerfile Location: backend/Dockerfile"
echo "   - Port: 8080"
echo "   - Domain: super.doctorhealthy1.com"
echo ""
echo "4. 🔐 Add Environment Variables:"
echo "   Copy all variables from .env.production to Coolify"
echo ""
echo "5. 🗄️  Add Database Services:"
echo "   - PostgreSQL 15 (nutrition-postgres)"
echo "   - Redis 7 (nutrition-redis)"
echo ""
echo "6. 🚀 Deploy:"
echo "   - Click 'Deploy' button"
echo "   - Wait 5-10 minutes for deployment"
echo ""
echo "7. ✅ Verify:"
echo "   - Check: https://super.doctorhealthy1.com/health"
echo "   - Test: https://super.doctorhealthy1.com/api/v1/info"
echo ""
echo -e "${GREEN}✅ COOLIFY DEPLOYMENT PACKAGE READY!${NC}"
echo ""
echo -e "${YELLOW}📋 Quick Deploy Command:${NC}"
echo "cd '$COOLIFY_DIR'"
echo "./deploy.sh"
echo ""
echo -e "${BLUE}🎉 Your secure nutrition platform is ready for Coolify deployment!${NC}"