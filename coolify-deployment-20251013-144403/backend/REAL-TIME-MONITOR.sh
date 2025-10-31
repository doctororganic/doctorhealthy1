#!/bin/bash

################################################################################
# REAL-TIME MONITOR
# Monitors application health and performance in real-time
################################################################################

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m'

API_URL="${API_URL:-http://localhost:8080}"
CHECK_INTERVAL="${CHECK_INTERVAL:-5}"

clear

echo -e "${CYAN}╔═══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║         REAL-TIME APPLICATION MONITOR                         ║${NC}"
echo -e "${CYAN}╚═══════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Statistics
total_checks=0
successful_checks=0
failed_checks=0

while true; do
    # Move cursor to top
    tput cup 4 0
    
    total_checks=$((total_checks + 1))
    
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} Monitoring $API_URL"
    echo "─────────────────────────────────────────────────────────────────"
    
    # Health Check
    if response=$(curl -s -w "\n%{http_code}" "$API_URL/health" 2>/dev/null); then
        http_code=$(echo "$response" | tail -n1)
        body=$(echo "$response" | head -n-1)
        
        if [ "$http_code" = "200" ]; then
            successful_checks=$((successful_checks + 1))
            echo -e "${GREEN}✓${NC} Health Check: ${GREEN}HEALTHY${NC} (HTTP $http_code)"
            
            # Parse JSON response
            if command -v jq &> /dev/null && [ -n "$body" ]; then
                status=$(echo "$body" | jq -r '.status // "unknown"' 2>/dev/null)
                uptime=$(echo "$body" | jq -r '.uptime // "unknown"' 2>/dev/null)
                echo -e "  Status: ${GREEN}$status${NC}"
                echo -e "  Uptime: ${CYAN}$uptime${NC}"
            fi
        else
            failed_checks=$((failed_checks + 1))
            echo -e "${RED}✗${NC} Health Check: ${RED}UNHEALTHY${NC} (HTTP $http_code)"
        fi
    else
        failed_checks=$((failed_checks + 1))
        echo -e "${RED}✗${NC} Health Check: ${RED}FAILED${NC} (Connection error)"
    fi
    
    echo ""
    
    # API Endpoints
    echo "API Endpoints:"
    
    endpoints=(
        "/api/v1/users"
        "/api/v1/foods"
        "/api/v1/workouts"
        "/api/v1/recipes"
        "/api/v1/health"
    )
    
    for endpoint in "${endpoints[@]}"; do
        if curl -s -f -o /dev/null "$API_URL$endpoint" 2>/dev/null; then
            echo -e "  ${GREEN}✓${NC} $endpoint"
        else
            echo -e "  ${RED}✗${NC} $endpoint"
        fi
    done
    
    echo ""
    
    # Statistics
    success_rate=0
    if [ $total_checks -gt 0 ]; then
        success_rate=$((successful_checks * 100 / total_checks))
    fi
    
    echo "Statistics:"
    echo -e "  Total Checks: ${CYAN}$total_checks${NC}"
    echo -e "  Successful: ${GREEN}$successful_checks${NC}"
    echo -e "  Failed: ${RED}$failed_checks${NC}"
    echo -e "  Success Rate: ${MAGENTA}$success_rate%${NC}"
    
    echo ""
    echo "─────────────────────────────────────────────────────────────────"
    echo -e "Refreshing in ${CHECK_INTERVAL}s... (Press Ctrl+C to stop)"
    
    sleep "$CHECK_INTERVAL"
done
