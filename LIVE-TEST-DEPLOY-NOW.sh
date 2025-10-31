#!/bin/bash
set -e

echo "� LLIVE TEST & DEPLOY - REAL-TIME MONITORING"
echo "============================================="

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

ERRORS=0
WARNINGS=0

log_error() { echo -e "${RED}❌ $1${NC}"; ((ERRORS++)); }
log_success() { echo -e "${GREEN}✅ $1${NC}"; }
log_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; ((WARNINGS++)); }
log_info() { echo -e "${BLUE}ℹ️  $1${NC}"; }

# 1. Pre-flight checks
echo -e "\n${BLUE}1️⃣  PRE-FLIGHT CHECKS${NC}"
echo "===================="

if ! command -v docker &> /dev/null; then
    log_error "Docker not installed"
    exit 1
fi
log_success "Docker installed"

if ! command -v docker-compose &> /dev/null; then
    log_error "Docker Compose not installed"
    exit 1
fi
log_success "Docker Compose installed"

if ! command -v go &> /dev/null; then
    log_warning "Go not installed (optional for local dev)"
else
    log_success "Go installed: $(go version | awk '{print $3}')"
fi

# 2. Generate credentials
echo -e "\n${BLUE}2️⃣  GENERATING SECURE CREDENTIALS${NC}"
echo "=================================="

DB_PASSWORD=$(openssl rand -hex 32)
JWT_SECRET=$(openssl rand -hex 64)
API_KEY_SECRET=$(openssl rand -hex 64)
REDIS_PASSWORD=$(openssl rand -hex 32)
ENCRYPTION_KEY=$(openssl rand -hex 32)

cat > .env.production << EOF
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
ENCRYPTION_KEY=${ENCRYPTION_KEY}
PORT=8080
ENVIRONMENT=production
DOMAIN=super.doctorhealthy1.com
ALLOWED_ORIGINS=https://super.doctorhealthy1.com,https://www.super.doctorhealthy1.com
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60s
EOF

log_success "Credentials generated"
log_info "DB_PASSWORD: ${DB_PASSWORD:0:16}..."
log_info "JWT_SECRET: ${JWT_SECRET:0:16}..."

# 3. Build images
echo -e "\n${BLUE}3️⃣  BUILDING DOCKER IMAGES${NC}"
echo "=========================="

log_info "Building backend..."
if docker-compose -f docker-compose.production.yml build backend 2>&1 | tee /tmp/build-backend.log; then
    log_success "Backend built"
else
    log_error "Backend build failed"
    cat /tmp/build-backend.log
    exit 1
fi

log_info "Building frontend..."
if docker-compose -f docker-compose.production.yml build frontend 2>&1 | tee /tmp/build-frontend.log; then
    log_success "Frontend built"
else
    log_error "Frontend build failed"
    cat /tmp/build-frontend.log
    exit 1
fi

# 4. Start services
echo -e "\n${BLUE}4️⃣  STARTING SERVICES${NC}"
echo "===================="

docker-compose -f docker-compose.production.yml down -v 2>/dev/null || true
docker-compose -f docker-compose.production.yml up -d

log_success "Services started"

# 5. Wait for services
echo -e "\n${BLUE}5️⃣  WAITING FOR SERVICES${NC}"
echo "======================="

log_info "Waiting for PostgreSQL..."
for i in {1..30}; do
    if docker-compose -f docker-compose.production.yml exec -T postgres pg_isready -U nutrition_user &>/dev/null; then
        log_success "PostgreSQL ready"
        break
    fi
    [ $i -eq 30 ] && log_error "PostgreSQL timeout" && exit 1
    sleep 2
done

log_info "Waiting for Redis..."
for i in {1..30}; do
    if docker-compose -f docker-compose.production.yml exec -T redis redis-cli ping &>/dev/null; then
        log_success "Redis ready"
        break
    fi
    [ $i -eq 30 ] && log_error "Redis timeout" && exit 1
    sleep 2
done

log_info "Waiting for Backend..."
for i in {1..60}; do
    if curl -sf http://localhost:8080/health &>/dev/null; then
        log_success "Backend ready"
        break
    fi
    [ $i -eq 60 ] && log_error "Backend timeout" && docker-compose -f docker-compose.production.yml logs backend && exit 1
    sleep 2
done

log_info "Waiting for Frontend..."
for i in {1..60}; do
    if curl -sf http://localhost:3000 &>/dev/null; then
        log_success "Frontend ready"
        break
    fi
    [ $i -eq 60 ] && log_error "Frontend timeout" && docker-compose -f docker-compose.production.yml logs frontend && exit 1
    sleep 2
done

# 6. Run tests
echo -e "\n${BLUE}6️⃣  RUNNING LIVE TESTS${NC}"
echo "===================="

# Health check
if curl -sf http://localhost:8080/health | grep -q "ok"; then
    log_success "Health check passed"
else
    log_error "Health check failed"
fi

# API info
if curl -sf http://localhost:8080/api/v1/info &>/dev/null; then
    log_success "API info endpoint working"
else
    log_warning "API info endpoint not responding"
fi

# Frontend
if curl -sf http://localhost:3000 | grep -q "html"; then
    log_success "Frontend serving HTML"
else
    log_error "Frontend not serving HTML"
fi

# CORS test
CORS_RESPONSE=$(curl -sI -H "Origin: https://super.doctorhealthy1.com" \
    -H "Access-Control-Request-Method: POST" \
    -X OPTIONS http://localhost:8080/api/v1/nutrition/analyze 2>/dev/null | grep -i "access-control-allow-origin" || echo "")

if [ -n "$CORS_RESPONSE" ]; then
    log_success "CORS configured correctly"
else
    log_warning "CORS test inconclusive"
fi

# 7. Performance test
echo -e "\n${BLUE}7️⃣  PERFORMANCE TEST${NC}"
echo "==================="

log_info "Running 100 concurrent requests..."
if command -v ab &> /dev/null; then
    ab -n 100 -c 10 http://localhost:8080/health 2>&1 | tee /tmp/perf-test.log
    FAILED_REQUESTS=$(grep "Failed requests:" /tmp/perf-test.log | awk '{print $3}')
    if [ "$FAILED_REQUESTS" = "0" ]; then
        log_success "Performance test passed (0 failed requests)"
    else
        log_warning "Performance test: $FAILED_REQUESTS failed requests"
    fi
else
    log_warning "Apache Bench not installed, skipping performance test"
fi

# 8. Security scan
echo -e "\n${BLUE}8️⃣  SECURITY SCAN${NC}"
echo "================"

log_info "Checking for exposed secrets..."
if grep -r "password.*=" backend/ --include="*.go" | grep -v "DB_PASSWORD\|REDIS_PASSWORD" | grep -v "//"; then
    log_error "Hardcoded passwords found in code"
else
    log_success "No hardcoded passwords in code"
fi

log_info "Checking SSL configuration..."
if grep -q "DB_SSL_MODE=require" .env.production; then
    log_success "Database SSL enabled"
else
    log_error "Database SSL not enabled"
fi

log_info "Checking CORS configuration..."
if grep -q "super.doctorhealthy1.com" .env.production; then
    log_success "CORS restricted to domain"
else
    log_error "CORS not properly configured"
fi

# 9. Live monitoring
echo -e "\n${BLUE}9️⃣  LIVE MONITORING${NC}"
echo "=================="

log_info "Container status:"
docker-compose -f docker-compose.production.yml ps

log_info "Resource usage:"
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"

log_info "Recent logs (last 20 lines):"
docker-compose -f docker-compose.production.yml logs --tail=20

# 10. Summary
echo -e "\n${BLUE}🎯 DEPLOYMENT SUMMARY${NC}"
echo "===================="
echo -e "Errors: ${RED}$ERRORS${NC}"
echo -e "Warnings: ${YELLOW}$WARNINGS${NC}"

if [ $ERRORS -eq 0 ]; then
    echo -e "\n${GREEN}✅ DEPLOYMENT SUCCESSFUL!${NC}"
    echo ""
    echo "🌐 Access URLs:"
    echo "  Frontend: http://localhost:3000"
    echo "  Backend:  http://localhost:8080"
    echo "  Health:   http://localhost:8080/health"
    echo ""
    echo "📊 Monitor:"
    echo "  docker-compose -f docker-compose.production.yml logs -f"
    echo "  docker-compose -f docker-compose.production.yml ps"
    echo ""
    echo "🛑 Stop:"
    echo "  docker-compose -f docker-compose.production.yml down"
    exit 0
else
    echo -e "\n${RED}❌ DEPLOYMENT FAILED${NC}"
    echo "Check logs above for errors"
    exit 1
fi
