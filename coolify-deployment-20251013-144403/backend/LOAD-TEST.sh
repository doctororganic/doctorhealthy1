#!/bin/bash

################################################################################
# LOAD TEST
# Performance and load testing for the API
################################################################################

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

API_URL="${API_URL:-http://localhost:8080}"
CONCURRENT_USERS="${CONCURRENT_USERS:-10}"
REQUESTS_PER_USER="${REQUESTS_PER_USER:-100}"

log() { echo -e "${GREEN}[$(date +'%H:%M:%S')]${NC} $1"; }
log_info() { echo -e "${BLUE}[$(date +'%H:%M:%S')] INFO:${NC} $1"; }
log_success() { echo -e "${GREEN}[$(date +'%H:%M:%S')] ✓${NC} $1"; }

log "═══════════════════════════════════════════════════════════════"
log "LOAD TEST - STARTING"
log "═══════════════════════════════════════════════════════════════"
log_info "API URL: $API_URL"
log_info "Concurrent Users: $CONCURRENT_USERS"
log_info "Requests per User: $REQUESTS_PER_USER"
echo ""

# Test endpoints
endpoints=(
    "/health"
    "/api/v1/users"
    "/api/v1/foods"
    "/api/v1/workouts"
)

total_requests=0
successful_requests=0
failed_requests=0
total_time=0

for endpoint in "${endpoints[@]}"; do
    log "Testing endpoint: $endpoint"
    
    start_time=$(date +%s)
    
    for ((i=1; i<=$REQUESTS_PER_USER; i++)); do
        if curl -s -f "$API_URL$endpoint" > /dev/null 2>&1; then
            successful_requests=$((successful_requests + 1))
        else
            failed_requests=$((failed_requests + 1))
        fi
        total_requests=$((total_requests + 1))
        
        # Progress indicator
        if [ $((i % 10)) -eq 0 ]; then
            echo -n "."
        fi
    done
    
    end_time=$(date +%s)
    duration=$((end_time - start_time))
    total_time=$((total_time + duration))
    
    echo ""
    log_success "Completed $REQUESTS_PER_USER requests in ${duration}s"
    echo ""
done

# Calculate statistics
success_rate=$((successful_requests * 100 / total_requests))
avg_time=$((total_time / ${#endpoints[@]}))

log "═══════════════════════════════════════════════════════════════"
log "LOAD TEST RESULTS"
log "═══════════════════════════════════════════════════════════════"
echo ""
log_info "Total Requests: $total_requests"
log_success "Successful: $successful_requests"
echo -e "${RED}Failed: $failed_requests${NC}"
echo -e "${CYAN}Success Rate: $success_rate%${NC}"
log_info "Total Time: ${total_time}s"
log_info "Average Time per Endpoint: ${avg_time}s"
echo ""

if [ $success_rate -ge 95 ]; then
    log_success "Load test PASSED! ✓"
    exit 0
else
    echo -e "${RED}Load test FAILED!${NC}"
    exit 1
fi
