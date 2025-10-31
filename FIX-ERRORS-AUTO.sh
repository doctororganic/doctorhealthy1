#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔧 AUTOMATIC ERROR DETECTION & FIXING${NC}"
echo "======================================"

# Check backend logs for errors
echo -e "\n${BLUE}Checking for errors...${NC}"

BACKEND_ERRORS=$(docker-compose -f docker-compose.production.yml logs backend 2>&1 | grep -i "error\|fatal\|panic" | tail -10)

if [ -n "$BACKEND_ERRORS" ]; then
    echo -e "${RED}❌ Found backend errors:${NC}"
    echo "$BACKEND_ERRORS"
    
    # Common fixes
    echo -e "\n${YELLOW}Applying automatic fixes...${NC}"
    
    # Fix 1: Database connection issues
    if echo "$BACKEND_ERRORS" | grep -qi "database\|postgres\|connection"; then
        echo "🔧 Restarting database..."
        docker-compose -f docker-compose.production.yml restart postgres
        sleep 5
    fi
    
    # Fix 2: Redis connection issues
    if echo "$BACKEND_ERRORS" | grep -qi "redis"; then
        echo "🔧 Restarting Redis..."
        docker-compose -f docker-compose.production.yml restart redis
        sleep 5
    fi
    
    # Fix 3: Backend crash
    if echo "$BACKEND_ERRORS" | grep -qi "panic\|fatal"; then
        echo "🔧 Restarting backend..."
        docker-compose -f docker-compose.production.yml restart backend
        sleep 10
    fi
    
    # Verify fixes
    echo -e "\n${BLUE}Verifying fixes...${NC}"
    sleep 5
    
    if curl -sf http://localhost:8080/health >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Backend is healthy now${NC}"
    else
        echo -e "${RED}❌ Backend still has issues${NC}"
        echo "Recent logs:"
        docker-compose -f docker-compose.production.yml logs --tail=20 backend
    fi
else
    echo -e "${GREEN}✅ No errors found${NC}"
fi

# Check container health
echo -e "\n${BLUE}Container Health Status:${NC}"
docker-compose -f docker-compose.production.yml ps

# Show resource usage
echo -e "\n${BLUE}Resource Usage:${NC}"
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"
