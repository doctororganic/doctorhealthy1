#!/usr/bin/env bash
# Test script for Phase 1: Backend Performance & Security
# Tests caching and rate limiting functionality

set -euo pipefail

BASE_URL="${1:-http://localhost:8080}"
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=========================================="
echo "Phase 1 Testing: Performance & Security"
echo "=========================================="
echo ""

# Test 1: Check if server is running
echo "Test 1: Server Health Check"
if curl -s -f "${BASE_URL}/health" > /dev/null; then
    echo -e "${GREEN}✅ Server is running${NC}"
else
    echo -e "${RED}❌ Server is not running. Please start the server first.${NC}"
    exit 1
fi
echo ""

# Test 2: Test Caching (Cache MISS then HIT)
echo "Test 2: Response Caching"
echo "Making first request (should be MISS)..."
RESPONSE1=$(curl -s -w "\n%{http_code}" "${BASE_URL}/api/v1/nutrition-data/recipes?limit=5")
HTTP_CODE1=$(echo "$RESPONSE1" | tail -n1)
CACHE_STATUS1=$(curl -s -I "${BASE_URL}/api/v1/nutrition-data/recipes?limit=5" | grep -i "X-Cache" | cut -d' ' -f2 | tr -d '\r')

if [ "$HTTP_CODE1" = "200" ]; then
    if [ -n "$CACHE_STATUS1" ]; then
        echo -e "  First request cache status: ${YELLOW}${CACHE_STATUS1}${NC}"
    else
        echo -e "  ${YELLOW}⚠️  Cache headers not present (cache may be disabled)${NC}"
    fi
    
    echo "Making second request (should be HIT if cache enabled)..."
    sleep 1
    CACHE_STATUS2=$(curl -s -I "${BASE_URL}/api/v1/nutrition-data/recipes?limit=5" | grep -i "X-Cache" | cut -d' ' -f2 | tr -d '\r')
    
    if [ -n "$CACHE_STATUS2" ]; then
        echo -e "  Second request cache status: ${YELLOW}${CACHE_STATUS2}${NC}"
        if [ "$CACHE_STATUS2" = "HIT" ]; then
            echo -e "${GREEN}✅ Caching is working correctly${NC}"
        else
            echo -e "${YELLOW}⚠️  Cache not hitting (may need more time or cache disabled)${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  Cache headers not present${NC}"
    fi
else
    echo -e "${RED}❌ Request failed with status ${HTTP_CODE1}${NC}"
fi
echo ""

# Test 3: Test Rate Limiting Headers
echo "Test 3: Rate Limiting Headers"
HEADERS=$(curl -s -I "${BASE_URL}/api/v1/nutrition-data/recipes?limit=5")
RATE_LIMIT=$(echo "$HEADERS" | grep -i "X-RateLimit" | head -1)

if [ -n "$RATE_LIMIT" ]; then
    echo -e "${GREEN}✅ Rate limit headers present:${NC}"
    echo "$HEADERS" | grep -i "X-RateLimit" | sed 's/^/  /'
else
    echo -e "${YELLOW}⚠️  Rate limit headers not present${NC}"
fi
echo ""

# Test 4: Test Rate Limiting Behavior (make many requests)
echo "Test 4: Rate Limiting Behavior"
echo "Making 10 rapid requests to test rate limiting..."
SUCCESS_COUNT=0
RATE_LIMITED=0

for i in {1..10}; do
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "${BASE_URL}/api/v1/nutrition-data/recipes?limit=1")
    if [ "$HTTP_CODE" = "200" ]; then
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    elif [ "$HTTP_CODE" = "429" ]; then
        RATE_LIMITED=$((RATE_LIMITED + 1))
    fi
    sleep 0.1
done

echo "  Successful requests: ${SUCCESS_COUNT}/10"
echo "  Rate limited requests: ${RATE_LIMITED}/10"

if [ $RATE_LIMITED -gt 0 ]; then
    echo -e "${GREEN}✅ Rate limiting is working (blocked ${RATE_LIMITED} requests)${NC}"
elif [ $SUCCESS_COUNT -eq 10 ]; then
    echo -e "${YELLOW}⚠️  All requests succeeded (rate limit may be high or disabled)${NC}"
else
    echo -e "${YELLOW}⚠️  Some requests failed for other reasons${NC}"
fi
echo ""

# Test 5: Test Security Headers
echo "Test 5: Security Headers"
SECURITY_HEADERS=("X-Frame-Options" "X-Content-Type-Options" "X-XSS-Protection" "Referrer-Policy")
MISSING_HEADERS=()

for header in "${SECURITY_HEADERS[@]}"; do
    if echo "$HEADERS" | grep -qi "$header"; then
        echo -e "  ${GREEN}✅${NC} $header"
    else
        echo -e "  ${RED}❌${NC} $header (missing)"
        MISSING_HEADERS+=("$header")
    fi
done

if [ ${#MISSING_HEADERS[@]} -eq 0 ]; then
    echo -e "${GREEN}✅ All security headers present${NC}"
else
    echo -e "${YELLOW}⚠️  Some security headers missing${NC}"
fi
echo ""

# Test 6: Performance Test (response time)
echo "Test 6: Response Time Performance"
echo "Measuring response time for cached vs uncached requests..."

# Uncached (first request)
TIME1=$(curl -s -o /dev/null -w "%{time_total}" "${BASE_URL}/api/v1/nutrition-data/recipes?limit=5&_nocache=$(date +%s)")
echo "  First request (uncached): ${TIME1}s"

# Cached (second request)
sleep 1
TIME2=$(curl -s -o /dev/null -w "%{time_total}" "${BASE_URL}/api/v1/nutrition-data/recipes?limit=5")
echo "  Second request (cached): ${TIME2}s"

if (( $(echo "$TIME2 < $TIME1" | bc -l) )); then
    IMPROVEMENT=$(echo "scale=2; (($TIME1 - $TIME2) / $TIME1) * 100" | bc)
    echo -e "${GREEN}✅ Performance improvement: ${IMPROVEMENT}% faster${NC}"
else
    echo -e "${YELLOW}⚠️  No significant performance improvement (cache may not be enabled)${NC}"
fi
echo ""

echo "=========================================="
echo "Testing Complete!"
echo "=========================================="

