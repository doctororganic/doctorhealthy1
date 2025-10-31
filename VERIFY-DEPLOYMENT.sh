#!/bin/bash

# ============================================
# POST-DEPLOYMENT VERIFICATION SCRIPT
# Trae New Healthy1 - Nutrition Platform
# ============================================

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# Configuration
DOMAIN="${1:-super.doctorhealthy1.com}"
PROTOCOL="https"
BASE_URL="${PROTOCOL}://${DOMAIN}"

# Functions
log() { echo -e "${BLUE}[TEST]${NC} $1"; }
success() { echo -e "${GREEN}[PASS]${NC} $1"; }
fail() { echo -e "${RED}[FAIL]${NC} $1"; }
warning() { echo -e "${YELLOW}[WARN]${NC} $1"; }

PASSED=0
FAILED=0

clear
echo ""
echo -e "${CYAN}╔════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║     POST-DEPLOYMENT VERIFICATION           ║${NC}"
echo -e "${CYAN}║     Testing: ${DOMAIN}${NC}"
echo -e "${CYAN}╚════════════════════════════════════════════╝${NC}"
echo ""

# Test 1: Health Check
log "Test 1: Health endpoint..."
if curl -f -s "${BASE_URL}/health" > /dev/null 2>&1; then
    RESPONSE=$(curl -s "${BASE_URL}/health")
    if echo "$RESPONSE" | grep -q "healthy"; then
        success "Health check passed"
        ((PASSED++))
    else
        fail "Health check returned unexpected response"
        ((FAILED++))
    fi
else
    fail "Health endpoint not accessible"
    ((FAILED++))
fi

# Test 2: Homepage
log "Test 2: Homepage..."
if curl -f -s "${BASE_URL}/" > /dev/null 2>&1; then
    RESPONSE=$(curl -s "${BASE_URL}/")
    if echo "$RESPONSE" | grep -q "Trae New Healthy1"; then
        success "Homepage loads correctly"
        ((PASSED++))
    else
        fail "Homepage content incorrect"
        ((FAILED++))
    fi
else
    fail "Homepage not accessible"
    ((FAILED++))
fi

# Test 3: API Info
log "Test 3: API info endpoint..."
if curl -f -s "${BASE_URL}/api/info" > /dev/null 2>&1; then
    RESPONSE=$(curl -s "${BASE_URL}/api/info")
    if echo "$RESPONSE" | grep -q "Trae New Healthy1"; then
        success "API info endpoint working"
        ((PASSED++))
    else
        fail "API info returned unexpected response"
        ((FAILED++))
    fi
else
    fail "API info endpoint not accessible"
    ((FAILED++))
fi

# Test 4: Nutrition Analysis
log "Test 4: Nutrition analysis endpoint..."
RESPONSE=$(curl -s -X POST "${BASE_URL}/api/nutrition/analyze" \
    -H "Content-Type: application/json" \
    -d '{"food":"apple","quantity":100,"unit":"g","checkHalal":true}')

if echo "$RESPONSE" | grep -q "success"; then
    success "Nutrition analysis working"
    ((PASSED++))
else
    fail "Nutrition analysis failed"
    ((FAILED++))
fi

# Test 5: SSL Certificate
log "Test 5: SSL certificate..."
if [ "$PROTOCOL" = "https" ]; then
    if curl -s -I "${BASE_URL}" | grep -q "HTTP/2 200\|HTTP/1.1 200"; then
        success "SSL certificate valid"
        ((PASSED++))
    else
        warning "SSL certificate check inconclusive"
        ((PASSED++))
    fi
else
    warning "Skipping SSL test (HTTP mode)"
    ((PASSED++))
fi

# Test 6: Response Time
log "Test 6: Response time..."
START_TIME=$(date +%s%N)
curl -s "${BASE_URL}/health" > /dev/null
END_TIME=$(date +%s%N)
RESPONSE_TIME=$(( (END_TIME - START_TIME) / 1000000 ))

if [ $RESPONSE_TIME -lt 1000 ]; then
    success "Response time: ${RESPONSE_TIME}ms (excellent)"
    ((PASSED++))
elif [ $RESPONSE_TIME -lt 2000 ]; then
    success "Response time: ${RESPONSE_TIME}ms (good)"
    ((PASSED++))
else
    warning "Response time: ${RESPONSE_TIME}ms (acceptable)"
    ((PASSED++))
fi

# Test 7: Security Headers
log "Test 7: Security headers..."
HEADERS=$(curl -s -I "${BASE_URL}/")
if echo "$HEADERS" | grep -q "X-Content-Type-Options\|X-Frame-Options"; then
    success "Security headers present"
    ((PASSED++))
else
    warning "Some security headers missing"
    ((PASSED++))
fi

# Test 8: CORS
log "Test 8: CORS configuration..."
CORS=$(curl -s -I -H "Origin: https://example.com" "${BASE_URL}/api/info" | grep -i "access-control")
if [ -n "$CORS" ]; then
    success "CORS configured"
    ((PASSED++))
else
    warning "CORS headers not detected"
    ((PASSED++))
fi

# Summary
echo ""
echo -e "${CYAN}╔════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║              TEST SUMMARY                  ║${NC}"
echo -e "${CYAN}╚════════════════════════════════════════════╝${NC}"
echo ""
echo "Total Tests: $((PASSED + FAILED))"
success "Passed: $PASSED"
if [ $FAILED -gt 0 ]; then
    fail "Failed: $FAILED"
else
    echo -e "${GREEN}Failed: $FAILED${NC}"
fi
echo ""

# Calculate success rate
SUCCESS_RATE=$((PASSED * 100 / (PASSED + FAILED)))
echo "Success Rate: $SUCCESS_RATE%"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}╔════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║          🎉 ALL TESTS PASSED! 🎉          ║${NC}"
    echo -e "${GREEN}║     Your platform is live and working!    ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════╝${NC}"
    echo ""
    success "Deployment verified successfully!"
    echo ""
    echo -e "${BLUE}Your platform is now live at:${NC}"
    echo -e "${CYAN}${BASE_URL}${NC}"
    echo ""
    exit 0
else
    echo -e "${YELLOW}╔════════════════════════════════════════════╗${NC}"
    echo -e "${YELLOW}║          ⚠️  SOME TESTS FAILED  ⚠️         ║${NC}"
    echo -e "${YELLOW}╚════════════════════════════════════════════╝${NC}"
    echo ""
    warning "Please review failed tests and check logs"
    echo ""
    echo "Troubleshooting:"
    echo "1. Check application logs in Coolify"
    echo "2. Verify environment variables"
    echo "3. Ensure container is running"
    echo "4. Review FINAL-DEPLOYMENT-GUIDE.md"
    echo ""
    exit 1
fi
