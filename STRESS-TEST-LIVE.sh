#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}⚡ HIGH-PERFORMANCE STRESS TEST${NC}"
echo "================================"

# Check if services are running
if ! curl -sf http://localhost:8080/health >/dev/null 2>&1; then
    echo -e "${RED}❌ Backend not running. Start with ./LIVE-TEST-DEPLOY.sh${NC}"
    exit 1
fi

echo -e "\n${BLUE}Test 1: Concurrent Requests (1000 requests, 50 concurrent)${NC}"
if command -v ab >/dev/null 2>&1; then
    ab -n 1000 -c 50 http://localhost:8080/health
else
    echo -e "${YELLOW}Installing Apache Bench...${NC}"
    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo "Apache Bench should be pre-installed on macOS"
    else
        sudo apt-get install -y apache2-utils
    fi
fi

echo -e "\n${BLUE}Test 2: Memory Usage${NC}"
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}"

echo -e "\n${BLUE}Test 3: Response Time Test${NC}"
for i in {1..10}; do
    start=$(date +%s%N)
    curl -sf http://localhost:8080/health >/dev/null
    end=$(date +%s%N)
    duration=$(( (end - start) / 1000000 ))
    echo "Request $i: ${duration}ms"
done

echo -e "\n${BLUE}Test 4: Database Connection Pool${NC}"
docker-compose -f docker-compose.production.yml exec -T postgres psql -U nutrition_user -d nutrition_platform -c "SELECT count(*) FROM pg_stat_activity;"

echo -e "\n${BLUE}Test 5: Redis Performance${NC}"
docker-compose -f docker-compose.production.yml exec -T redis redis-cli -a "$(grep REDIS_PASSWORD .env.production | cut -d'=' -f2)" INFO stats | grep -E "total_commands_processed|instantaneous_ops_per_sec"

echo -e "\n${GREEN}✅ Stress test complete${NC}"
