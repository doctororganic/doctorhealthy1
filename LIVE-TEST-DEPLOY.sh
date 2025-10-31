#!/bin/bash
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🚀 LIVE DEPLOYMENT WITH REAL-TIME TESTING${NC}"
echo "=========================================="

# Step 1: Generate credentials
echo -e "\n${BLUE}[1/8] Generating secure credentials...${NC}"
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

# Step 2: Clean previous deployment
echo -e "\n${BLUE}[2/8] Cleaning previous deployment...${NC}"
docker-compose -f docker-compose.production.yml down -v 2>/dev/null || true
echo -e "${GREEN}✅ Cleaned${NC}"

# Step 3: Build images
echo -e "\n${BLUE}[3/8] Building Docker images...${NC}"
docker-compose -f docker-compose.production.yml build --no-cache 2>&1 | tee build.log
if [ ${PIPESTATUS[0]} -eq 0 ]; then
    echo -e "${GREEN}✅ Build successful${NC}"
else
    echo -e "${RED}❌ Build failed - check build.log${NC}"
    exit 1
fi

# Step 4: Start services
echo -e "\n${BLUE}[4/8] Starting services...${NC}"
docker-compose -f docker-compose.production.yml up -d
echo -e "${GREEN}✅ Services started${NC}"

# Step 5: Wait for services
echo -e "\n${BLUE}[5/8] Waiting for services to be ready...${NC}"

echo "Waiting for PostgreSQL..."
for i in {1..30}; do
    if docker-compose -f docker-compose.production.yml exec -T postgres pg_isready -U nutrition_user >/dev/null 2>&1; then
        echo -e "${GREEN}✅ PostgreSQL ready${NC}"
        break
    fi
    [ $i -eq 30 ] && echo -e "${RED}❌ PostgreSQL timeout${NC}" && exit 1
    sleep 2
done

echo "Waiting for Redis..."
for i in {1..30}; do
    if docker-compose -f docker-compose.production.yml exec -T redis redis-cli -a "${REDIS_PASSWORD}" ping >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Redis ready${NC}"
        break
    fi
    [ $i -eq 30 ] && echo -e "${RED}❌ Redis timeout${NC}" && exit 1
    sleep 2
done

echo "Waiting for Backend..."
for i in {1..60}; do
    if curl -sf http://localhost:8080/health >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Backend ready${NC}"
        break
    fi
    [ $i -eq 60 ] && echo -e "${RED}❌ Backend timeout${NC}" && docker-compose -f docker-compose.production.yml logs backend && exit 1
    sleep 2
done

echo "Waiting for Frontend..."
for i in {1..60}; do
    if curl -sf http://localhost:3000 >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Frontend ready${NC}"
        break
    fi
    [ $i -eq 60 ] && echo -e "${RED}❌ Frontend timeout${NC}" && docker-compose -f docker-compose.production.yml logs frontend && exit 1
    sleep 2
done

# Step 6: Run comprehensive tests
echo -e "\n${BLUE}[6/8] Running comprehensive tests...${NC}"

PASSED=0
FAILED=0

test_endpoint() {
    local name=$1
    local url=$2
    local expected=$3
    
    response=$(curl -s -o /dev/null -w "%{http_code}" "$url" 2>/dev/null || echo "000")
    if [ "$response" = "$expected" ]; then
        echo -e "${GREEN}✅ $name${NC}"
        ((PASSED++))
    else
        echo -e "${RED}❌ $name (got $response, expected $expected)${NC}"
        ((FAILED++))
    fi
}

test_endpoint "Backend Health" "http://localhost:8080/health" "200"
test_endpoint "Frontend" "http://localhost:3000" "200"
test_endpoint "API Info" "http://localhost:8080/api/v1/info" "200"

# Step 7: Performance test
echo -e "\n${BLUE}[7/8] Running performance tests...${NC}"

if command -v ab >/dev/null 2>&1; then
    echo "Testing backend performance (100 requests, 10 concurrent)..."
    ab -n 100 -c 10 -q http://localhost:8080/health 2>&1 | grep -E "Requests per second|Time per request|Failed requests" || true
    echo -e "${GREEN}✅ Performance test complete${NC}"
else
    echo -e "${YELLOW}⚠️  Apache Bench not installed, skipping performance test${NC}"
fi

# Step 8: Display results
echo -e "\n${BLUE}[8/8] Deployment Summary${NC}"
echo "========================================"
echo -e "Tests Passed: ${GREEN}${PASSED}${NC}"
echo -e "Tests Failed: ${RED}${FAILED}${NC}"
echo ""
echo "Services Status:"
docker-compose -f docker-compose.production.yml ps
echo ""
echo "Access URLs:"
echo "  Frontend: http://localhost:3000"
echo "  Backend:  http://localhost:8080"
echo "  Health:   http://localhost:8080/health"
echo ""
echo "Credentials (SAVE THESE):"
echo "  DB_PASSWORD: ${DB_PASSWORD}"
echo "  REDIS_PASSWORD: ${REDIS_PASSWORD}"
echo "  JWT_SECRET: ${JWT_SECRET}"
echo ""
echo "View logs: docker-compose -f docker-compose.production.yml logs -f"
echo "Stop: docker-compose -f docker-compose.production.yml down"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 DEPLOYMENT SUCCESSFUL!${NC}"
    exit 0
else
    echo -e "${RED}⚠️  DEPLOYMENT COMPLETED WITH ERRORS${NC}"
    exit 1
fi
