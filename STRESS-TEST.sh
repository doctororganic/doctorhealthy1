#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}⚡ HIGH-PERFORMANCE STRESS TEST${NC}"
echo "================================"

# Check if ab (Apache Bench) is installed
if ! command -v ab &> /dev/null; then
    echo -e "${YELLOW}Installing Apache Bench...${NC}"
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS - ab comes with Apache
        echo "Apache Bench should be pre-installed on macOS"
    else
        sudo apt-get install -y apache2-utils
    fi
fi

# Test 1: Health endpoint stress test
echo -e "\n${BLUE}Test 1: Health Endpoint (1000 requests, 10 concurrent)${NC}"
ab -n 1000 -c 10 http://localhost:8080/health 2>&1 | grep -E "Requests per second|Time per request|Failed requests"

# Test 2: Concurrent connections
echo -e "\n${BLUE}Test 2: High Concurrency (5000 requests, 50 concurrent)${NC}"
ab -n 5000 -c 50 http://localhost:8080/health 2>&1 | grep -E "Requests per second|Time per request|Failed requests"

# Test 3: Sustained load
echo -e "\n${BLUE}Test 3: Sustained Load (10000 requests, 100 concurrent)${NC}"
ab -n 10000 -c 100 http://localhost:8080/health 2>&1 | grep -E "Requests per second|Time per request|Failed requests"

# Test 4: Database stress
echo -e "\n${BLUE}Test 4: Database Connection Pool${NC}"
for i in {1..10}; do
    docker-compose -f docker-compose.production.yml exec -T postgres psql -U nutrition_user -d nutrition_platform -c "SELECT COUNT(*) FROM pg_stat_activity;" &
done
wait
echo -e "${GREEN}✅ Database connection test complete${NC}"

# Test 5: Memory leak check
echo -e "\n${BLUE}Test 5: Memory Usage Before/After Load${NC}"
echo "Memory before:"
docker stats --no-stream --format "{{.Name}}: {{.MemUsage}}" | grep backend

echo "Running load test..."
ab -n 5000 -c 50 http://localhost:8080/health >/dev/null 2>&1

sleep 5
echo "Memory after:"
docker stats --no-stream --format "{{.Name}}: {{.MemUsage}}" | grep backend

echo -e "\n${GREEN}✅ Stress tests complete${NC}"
