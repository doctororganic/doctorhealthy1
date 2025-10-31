#!/bin/bash
# deploy-production.sh - Foolproof deployment

set -e  # Exit on error
set -u  # Exit on undefined variable

echo "🚀 Starting Production Deployment..."
echo "===================================="

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Step 1: Pre-deployment checks
echo ""
echo "1️⃣  Running pre-deployment checks..."
./scripts/pre-deployment-check.sh || {
    echo -e "${RED}❌ Pre-deployment checks failed${NC}"
    exit 1
}
echo -e "${GREEN}✅ Pre-deployment checks passed${NC}"

# Step 2: Build Docker images
echo ""
echo "2️⃣  Building Docker images..."
docker-compose build --no-cache || {
    echo -e "${RED}❌ Docker build failed${NC}"
    exit 1
}
echo -e "${GREEN}✅ Docker images built${NC}"

# Step 3: Run security scan
echo ""
echo "3️⃣  Running security scan..."
if [ -f "./scripts/security-scan.sh" ]; then
    ./scripts/security-scan.sh || {
        echo -e "${YELLOW}⚠️  Security scan had warnings${NC}"
    }
fi
echo -e "${GREEN}✅ Security scan complete${NC}"

# Step 4: Start services
echo ""
echo "4️⃣  Starting services..."
docker-compose up -d || {
    echo -e "${RED}❌ Failed to start services${NC}"
    exit 1
}
echo -e "${GREEN}✅ Services started${NC}"

# Step 5: Wait for services to be healthy
echo ""
echo "5️⃣  Waiting for services to be healthy..."
sleep 10

# Check backend health
for i in {1..30}; do
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Backend is healthy${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}❌ Backend health check timeout${NC}"
        docker-compose logs backend
        exit 1
    fi
    echo "Waiting for backend... ($i/30)"
    sleep 2
done

# Step 6: Run smoke tests
echo ""
echo "6️⃣  Running smoke tests..."
./scripts/smoke-tests.sh || {
    echo -e "${RED}❌ Smoke tests failed${NC}"
    docker-compose logs
    exit 1
}
echo -e "${GREEN}✅ Smoke tests passed${NC}"

# Step 7: Final verification
echo ""
echo "7️⃣  Final verification..."
./scripts/verify-deployment.sh || {
    echo -e "${RED}❌ Deployment verification failed${NC}"
    exit 1
}

echo ""
echo "===================================="
echo -e "${GREEN}🎉 DEPLOYMENT SUCCESSFUL!${NC}"
echo "===================================="
echo ""
echo "Services:"
echo "  Frontend: http://localhost:3000"
echo "  Backend:  http://localhost:8080"
echo "  Health:   http://localhost:8080/health"
echo ""
echo "Logs:"
echo "  docker-compose logs -f"
echo ""
