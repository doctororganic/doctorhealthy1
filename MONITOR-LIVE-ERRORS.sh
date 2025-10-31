#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}📊 LIVE ERROR MONITORING${NC}"
echo "========================"

# Create monitoring log
MONITOR_LOG="monitoring-$(date +%Y%m%d-%H%M%S).log"

echo "Monitoring started at $(date)" > $MONITOR_LOG
echo "Press Ctrl+C to stop monitoring"
echo ""

# Function to check service health
check_health() {
    local service=$1
    local url=$2
    
    if curl -sf "$url" >/dev/null 2>&1; then
        echo -e "${GREEN}✅ $service OK${NC}"
        echo "$(date) - $service: OK" >> $MONITOR_LOG
    else
        echo -e "${RED}❌ $service FAILED${NC}"
        echo "$(date) - $service: FAILED" >> $MONITOR_LOG
        return 1
    fi
}

# Function to check container status
check_containers() {
    local failed=0
    while IFS= read -r line; do
        if echo "$line" | grep -q "Up"; then
            echo -e "${GREEN}✅ $(echo $line | awk '{print $1}')${NC}"
        else
            echo -e "${RED}❌ $(echo $line | awk '{print $1}')${NC}"
            ((failed++))
        fi
    done < <(docker-compose -f docker-compose.production.yml ps --format "table {{.Name}}\t{{.Status}}" | tail -n +2)
    return $failed
}

# Function to check errors in logs
check_errors() {
    local service=$1
    local errors=$(docker-compose -f docker-compose.production.yml logs --tail=50 $service 2>&1 | grep -iE "error|fatal|panic|exception" | wc -l)
    
    if [ $errors -gt 0 ]; then
        echo -e "${RED}⚠️  $service: $errors errors found${NC}"
        echo "$(date) - $service: $errors errors" >> $MONITOR_LOG
        docker-compose -f docker-compose.production.yml logs --tail=10 $service | grep -iE "error|fatal|panic|exception"
    else
        echo -e "${GREEN}✅ $service: No errors${NC}"
    fi
}

# Monitoring loop
while true; do
    clear
    echo -e "${BLUE}📊 LIVE MONITORING - $(date)${NC}"
    echo "========================================"
    
    echo -e "\n${BLUE}Service Health:${NC}"
    check_health "Backend" "http://localhost:8080/health"
    check_health "Frontend" "http://localhost:3000"
    
    echo -e "\n${BLUE}Container Status:${NC}"
    check_containers
    
    echo -e "\n${BLUE}Recent Errors:${NC}"
    check_errors "backend"
    check_errors "frontend"
    check_errors "postgres"
    check_errors "redis"
    
    echo -e "\n${BLUE}Resource Usage:${NC}"
    docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"
    
    echo -e "\n${BLUE}Database Connections:${NC}"
    docker-compose -f docker-compose.production.yml exec -T postgres psql -U nutrition_user -d nutrition_platform -c "SELECT count(*) as active_connections FROM pg_stat_activity;" 2>/dev/null || echo "Unable to check"
    
    echo -e "\n${YELLOW}Monitoring log: $MONITOR_LOG${NC}"
    echo "Refreshing in 5 seconds... (Ctrl+C to stop)"
    
    sleep 5
done
