#!/bin/bash

# ============================================
# FINAL DEPLOYMENT SCRIPT
# Trae New Healthy1 - Nutrition Platform
# ============================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

# Functions
log() { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; }
warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
step() { echo -e "${PURPLE}[STEP]${NC} $1"; }

clear
echo ""
echo -e "${CYAN}╔════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║                                            ║${NC}"
echo -e "${CYAN}║     🚀 FINAL DEPLOYMENT SCRIPT 🚀         ║${NC}"
echo -e "${CYAN}║     Trae New Healthy1 Platform            ║${NC}"
echo -e "${CYAN}║                                            ║${NC}"
echo -e "${CYAN}╚════════════════════════════════════════════╝${NC}"
echo ""

# Step 1: Pre-deployment checks
step "1/6: Running pre-deployment checks..."
sleep 1

log "Checking required files..."
REQUIRED_FILES=(
    "Dockerfile"
    ".dockerignore"
    "production-nodejs/server.js"
    "production-nodejs/package.json"
)

ALL_EXIST=true
for file in "${REQUIRED_FILES[@]}"; do
    if [ -f "$file" ]; then
        success "✓ $file"
    else
        error "✗ $file missing"
        ALL_EXIST=false
    fi
done

if [ "$ALL_EXIST" = false ]; then
    error "Missing required files. Cannot proceed."
    exit 1
fi

success "All required files present"
echo ""

# Step 2: Validate code
step "2/6: Validating code..."
sleep 1

if command -v node &> /dev/null; then
    log "Checking Node.js syntax..."
    if node --check production-nodejs/server.js 2>/dev/null; then
        success "✓ Node.js syntax valid"
    else
        error "✗ Node.js syntax error"
        exit 1
    fi
else
    warning "Node.js not found, skipping syntax check"
fi

success "Code validation complete"
echo ""

# Step 3: Test Docker build
step "3/6: Testing Docker build..."
sleep 1

log "Building Docker image (this may take 2-3 minutes)..."
if docker build -t trae-healthy1-test -f Dockerfile . > /tmp/docker-build.log 2>&1; then
    success "✓ Docker build successful"
else
    error "✗ Docker build failed"
    echo "Check /tmp/docker-build.log for details"
    exit 1
fi

success "Docker build test complete"
echo ""

# Step 4: Test container
step "4/6: Testing container..."
sleep 1

log "Starting test container..."
docker run -d -p 8080:3000 --name trae-test trae-healthy1-test > /dev/null 2>&1

log "Waiting for container to start..."
sleep 5

log "Testing health endpoint..."
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    success "✓ Health check passed"
else
    error "✗ Health check failed"
    docker logs trae-test
    docker stop trae-test > /dev/null 2>&1
    docker rm trae-test > /dev/null 2>&1
    exit 1
fi

log "Cleaning up test container..."
docker stop trae-test > /dev/null 2>&1
docker rm trae-test > /dev/null 2>&1
docker rmi trae-healthy1-test > /dev/null 2>&1

success "Container test complete"
echo ""

# Step 5: Deployment summary
step "5/6: Deployment summary..."
sleep 1

echo ""
echo -e "${CYAN}╔════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║         DEPLOYMENT READY SUMMARY           ║${NC}"
echo -e "${CYAN}╚════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}✓${NC} Code validation: PASSED"
echo -e "${GREEN}✓${NC} Docker build: PASSED"
echo -e "${GREEN}✓${NC} Container test: PASSED"
echo -e "${GREEN}✓${NC} Health check: PASSED"
echo -e "${GREEN}✓${NC} All systems: GO"
echo ""
echo -e "${BLUE}Platform:${NC} Coolify"
echo -e "${BLUE}Domain:${NC} super.doctorhealthy1.com"
echo -e "${BLUE}Port:${NC} 3000 (internal)"
echo -e "${BLUE}SSL:${NC} Auto-configured"
echo ""

# Step 6: Deployment instructions
step "6/6: Next steps..."
sleep 1

echo ""
echo -e "${CYAN}╔════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║           DEPLOYMENT OPTIONS               ║${NC}"
echo -e "${CYAN}╚════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${YELLOW}Option 1: Coolify Dashboard (Recommended)${NC}"
echo "  1. Login to https://api.doctorhealthy1.com"
echo "  2. Navigate to 'new doctorhealthy1' project"
echo "  3. Create/update application"
echo "  4. Set domain: super.doctorhealthy1.com"
echo "  5. Set port: 3000"
echo "  6. Add environment variables:"
echo "     - NODE_ENV=production"
echo "     - PORT=3000"
echo "     - HOST=0.0.0.0"
echo "  7. Click 'Deploy'"
echo ""
echo -e "${YELLOW}Option 2: Manual Docker Deploy${NC}"
echo "  docker build -t trae-healthy1 ."
echo "  docker run -d -p 3000:3000 trae-healthy1"
echo ""

# Final message
echo ""
echo -e "${CYAN}╔════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║              🎉 SUCCESS! 🎉                ║${NC}"
echo -e "${CYAN}╚════════════════════════════════════════════╝${NC}"
echo ""
success "All pre-deployment checks passed!"
success "Your platform is ready to deploy!"
echo ""
echo -e "${GREEN}Next:${NC} Choose a deployment option above"
echo -e "${GREEN}Docs:${NC} See FINAL-DEPLOYMENT-GUIDE.md"
echo -e "${GREEN}Help:${NC} Check troubleshooting guides"
echo ""
echo -e "${PURPLE}Your months of work are about to pay off!${NC}"
echo -e "${PURPLE}Deploy with confidence! 🚀${NC}"
echo ""
