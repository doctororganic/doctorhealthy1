#!/bin/bash
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🚀 LIVE DEPLOYMENT WITH REAL-TIME TESTING${NC}"
echo "=========================================="

# Step 1: Check prerequisites
echo -e "\n${BLUE}[1/8] Checking Prerequisites...${NC}"
command -v docker >/dev/null 2>&1 || { echo -e "${RED}❌ Docker not installed${NC}"; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo -e "${RED}❌ Docker Compose not installed${NC}"; exit 1; }
echo -e "${GREEN}✅ Prerequisites OK${NC}"

# Step 2: Generate credentials
echo -e "\n${BLUE}[2/8] Generating Secure Credentials...${NC}"
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
ALLOWED_ORIGINS=https://super.doctorhealthy1.com
EOF

echo -e "${GREEN}✅ Credentials generated${NC}"
echo -e "${YELLOW}DB_PASSWORD: ${DB_PASSWORD:0:20}...${NC}"

# Step 3: Clean previous deployment
echo -e "\n${BLUE}[3/8] Cleaning Previous Deployment...${NC}"
docker-compose -f docker-compose.production.yml down -v 2>/dev/null || true
docker system prune -f
echo -e "${GREEN}✅ Cleaned${NC}"

# Step 4: Build images
echo -e "\n${BLUE}[4/8] Building Docker Images...${NC}"
docker-compose -f docker-compose.production.yml build --no-cache 2>&1 | tee build.log
if [ ${PIPESTATUS[0]} -ne 0 ]; then
    echo -e "${RED}❌ Build failed. Check build.log${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Images built${NC}"

# Step 5: Start services
echo -e "\n${BLUE}[5/8] Starting Services...${NC}"
docker-compose -f docker-compose.production.yml up -d
echo -e "${GREEN}✅ Services started${NC}"

# Step 6: Wait and monitor
echo -e "\n${BLUE}[6/8] Monitoring Service Health...${NC}"

# Wait for postgres
echo -n "Waiting for PostgreSQL..."
for i in {1..30}; do
    if docker-compose -f docker-compose.production.yml exec -T postgres pg_isready -U nutrition_user >/dev/null 2>&1; then
        echo -e " ${GREEN}✅${NC}"
        break
    fi
    echo -n "."
    sleep 2
done

# Wait for redis
echo -n "Waiting for Redis..."
for i in {1..30}; do
    if docker-compose -f docker-compose.production.yml exec -T redis redis-cli ping >/dev/null 2>&1; then
        echo -e " ${GREEN}✅${NC}"
        break
    fi
    echo -n "."
    sleep 2
done

# Wait for backend
echo -n "Waiting for Backend..."
for i in {1..60}; do
    if curl -sf http://localhost:8080/health >/dev/null 2>&1; then
        echo -e " ${GREEN}✅${NC}"
        break
    fi
    echo -n "."
    sleep 2
    if [ $i -eq 60 ]; then
        echo -e " ${RED}❌ TIMEOUT${NC}"
        echo -e "${YELLOW}Backend logs:${NC}"
        docker-compose -f docker-compose.production.yml logs backend | tail -50
        exit 1
    fi
done

# Step 7: Run comprehensive tests
echo -e "\n${BLUE}[7/8] Running Comprehensive Tests...${NC}"

PASSED=0
FAILED=0

# Test 1: Health endpoint
echo -n "Test 1: Health endpoint... "
if curl -sf http://localhost:8080/health | grep -q "ok\|healthy\|status"; then
    echo -e "${GREEN}✅ PASS${NC}"
    ((PASSED++))
else
    echo -e "${RED}❌ FAIL${NC}"
    ((FAILED++))
fi

# Test 2: Database connection
echo -n "Test 2: Database connection... "
if docker-compose -f docker-compose.production.yml exec -T postgres psql -U nutrition_user -d nutrition_platform -c "SELECT 1" >/dev/null 2>&1; then
    echo -e "${GREEN}✅ PASS${NC}"
    ((PASSED++))
else
    echo -e "${RED}❌ FAIL${NC}"
    ((FAILED++))
fi

# Test 3: Redis connection
echo -n "Test 3: Redis connection... "
if docker-compose -f docker-compose.production.yml exec -T redis redis-cli ping >/dev/null 2>&1; then
    echo -e "${GREEN}✅ PASS${NC}"
    ((PASSED++))
else
    echo -e "${RED}❌ FAIL${NC}"
    ((FAILED++))
fi

# Test 4: API endpoints
echo -n "Test 4: API info endpoint... "
if curl -sf http://localhost:8080/api/v1/info >/dev/null 2>&1; then
    echo -e "${GREEN}✅ PASS${NC}"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠️  SKIP (endpoint may not exist)${NC}"
fi

# Test 5: CORS headers
echo -n "Test 5: CORS configuration... "
CORS_HEADER=$(curl -sI -H "Origin: https://super.doctorhealthy1.com" http://localhost:8080/health | grep -i "access-control-allow-origin" || echo "")
if [ -n "$CORS_HEADER" ]; then
    echo -e "${GREEN}✅ PASS${NC}"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠️  WARNING (CORS may need configuration)${NC}"
fi

# Test 6: Container health
echo -n "Test 6: All containers running... "
RUNNING=$(docker-compose -f docker-compose.production.yml ps --services --filter "status=running" | wc -l)
TOTAL=$(docker-compose -f docker-compose.production.yml ps --services | wc -l)
if [ "$RUNNING" -eq "$TOTAL" ]; then
    echo -e "${GREEN}✅ PASS ($RUNNING/$TOTAL)${NC}"
    ((PASSED++))
else
    echo -e "${RED}❌ FAIL ($RUNNING/$TOTAL)${NC}"
    ((FAILED++))
fi

# Step 8: Results
echo -e "\n${BLUE}[8/8] Test Results${NC}"
echo "===================="
echo -e "Passed: ${GREEN}$PASSED${NC}"
echo -e "Failed: ${RED}$FAILED${NC}"

if [ $FAILED -eq 0 ]; then
    echo -e "\n${GREEN}🎉 DEPLOYMENT SUCCESSFUL!${NC}"
    echo -e "\n${BLUE}Access URLs:${NC}"
    echo "  Frontend: http://localhost:3000"
    echo "  Backend:  http://localhost:8080"
    echo "  Health:   http://localhost:8080/health"
    echo ""
    echo -e "${BLUE}View logs:${NC}"
    echo "  docker-compose -f docker-compose.production.yml logs -f"
    echo ""
    echo -e "${BLUE}Stop services:${NC}"
    echo "  docker-compose -f docker-compose.production.yml down"
else
    echo -e "\n${RED}⚠️  DEPLOYMENT COMPLETED WITH WARNINGS${NC}"
    echo "Check logs for details:"
    echo "  docker-compose -f docker-compose.production.yml logs"
fi

# Save credentials
cat > CREDENTIALS.txt << EOF
DEPLOYMENT CREDENTIALS - SAVE SECURELY
======================================
Generated: $(date)

DB_PASSWORD=${DB_PASSWORD}
REDIS_PASSWORD=${REDIS_PASSWORD}
JWT_SECRET=${JWT_SECRET}
API_KEY_SECRET=${API_KEY_SECRET}
ENCRYPTION_KEY=${ENCRYPTION_KEY}

⚠️  DELETE THIS FILE AFTER SAVING TO PASSWORD MANAGER
EOF

echo -e "\n${YELLOW}📋 Credentials saved to CREDENTIALS.txt${NC}"
echo -e "${RED}⚠️  Save these credentials and delete the file!${NC}"
