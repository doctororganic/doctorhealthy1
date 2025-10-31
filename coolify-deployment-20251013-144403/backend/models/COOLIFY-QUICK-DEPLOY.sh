#!/bin/bash

################################################################################
# COOLIFY QUICK DEPLOY
# Prepares everything for Coolify deployment
################################################################################

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

log() { echo -e "${GREEN}[$(date +'%H:%M:%S')]${NC} $1"; }
log_success() { echo -e "${GREEN}[$(date +'%H:%M:%S')] ✓${NC} $1"; }
log_info() { echo -e "${BLUE}[$(date +'%H:%M:%S')] INFO:${NC} $1"; }

clear

echo -e "${CYAN}"
cat << 'EOF'
╔═══════════════════════════════════════════════════════════════╗
║                                                               ║
║         COOLIFY DEPLOYMENT PREPARATION                        ║
║                                                               ║
╚═══════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

cd ..

log "═══════════════════════════════════════════════════════════════"
log "Step 1: Building Application"
log "═══════════════════════════════════════════════════════════════"

go build -o bin/server ./cmd/server
log_success "Application built"

log "═══════════════════════════════════════════════════════════════"
log "Step 2: Creating Deployment Package"
log "═══════════════════════════════════════════════════════════════"

tar -czf nutrition-platform-coolify.tar.gz \
  bin/ \
  migrations/ \
  handlers/ \
  models/ \
  services/ \
  middleware/ \
  config/ \
  cmd/ \
  go.mod \
  go.sum \
  Dockerfile \
  .env.example

log_success "Deployment package created: nutrition-platform-coolify.tar.gz"
ls -lh nutrition-platform-coolify.tar.gz

log "═══════════════════════════════════════════════════════════════"
log "Step 3: Creating Coolify Configuration"
log "═══════════════════════════════════════════════════════════════"

cat > coolify-config.json << 'EOFCONFIG'
{
  "name": "nutrition-platform",
  "type": "application",
  "build_pack": "dockerfile",
  "port": 8080,
  "health_check": {
    "path": "/health",
    "interval": 30,
    "timeout": 10,
    "retries": 3
  },
  "environment": {
    "PORT": "8080",
    "ENVIRONMENT": "production"
  },
  "services": {
    "postgres": {
      "image": "postgres:15-alpine",
      "environment": {
        "POSTGRES_DB": "nutrition_platform",
        "POSTGRES_USER": "postgres",
        "POSTGRES_PASSWORD": "CHANGE_ME"
      }
    },
    "redis": {
      "image": "redis:7-alpine"
    }
  }
}
EOFCONFIG

log_success "Coolify configuration created"

log "═══════════════════════════════════════════════════════════════"
log "Step 4: Creating Environment Template"
log "═══════════════════════════════════════════════════════════════"

cat > .env.coolify << 'EOFENV'
# Copy these to Coolify Environment Variables

PORT=8080
ENVIRONMENT=production

# Database (Coolify will provide these)
DB_HOST=postgres
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=postgres
DB_PASSWORD=CHANGE_ME_IN_COOLIFY

# Redis (Coolify will provide these)
REDIS_HOST=redis
REDIS_PORT=6379

# Security (Generate secure values in Coolify)
JWT_SECRET=GENERATE_32_CHAR_SECRET_IN_COOLIFY
API_KEY_SECRET=GENERATE_32_CHAR_SECRET_IN_COOLIFY

# CORS
ALLOWED_ORIGINS=https://yourdomain.com

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60
EOFENV

log_success "Environment template created: .env.coolify"

log "═══════════════════════════════════════════════════════════════"
log "DEPLOYMENT READY!"
log "═══════════════════════════════════════════════════════════════"

echo ""
echo -e "${GREEN}✅ All files prepared for Coolify deployment${NC}"
echo ""
echo "📦 Deployment Package:"
echo "   nutrition-platform-coolify.tar.gz ($(du -h nutrition-platform-coolify.tar.gz | cut -f1))"
echo ""
echo "📋 Configuration Files:"
echo "   ✓ coolify-config.json"
echo "   ✓ .env.coolify"
echo "   ✓ Dockerfile"
echo ""
echo -e "${CYAN}Next Steps:${NC}"
echo ""
echo "1. Go to your Coolify dashboard"
echo "   https://your-coolify-instance.com"
echo ""
echo "2. Create new project: 'nutrition-platform'"
echo ""
echo "3. Add application:"
echo "   - Type: Application"
echo "   - Build Pack: Dockerfile"
echo "   - Port: 8080"
echo "   - Health Check: /health"
echo ""
echo "4. Add services:"
echo "   - PostgreSQL 15"
echo "   - Redis 7"
echo ""
echo "5. Copy environment variables from .env.coolify"
echo ""
echo "6. Deploy!"
echo ""
echo -e "${YELLOW}📖 Full guide: COOLIFY-DEPLOYMENT-GUIDE.md${NC}"
echo ""
