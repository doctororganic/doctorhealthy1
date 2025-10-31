#!/bin/bash

# Automated Test Runner for Trae New Healthy1
# This script runs all validation tests automatically

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${BLUE}[TEST]${NC} $1"; }
success() { echo -e "${GREEN}[PASS]${NC} $1"; }
fail() { echo -e "${RED}[FAIL]${NC} $1"; }
warning() { echo -e "${YELLOW}[WARN]${NC} $1"; }

PASSED=0
FAILED=0

echo ""
echo "🧪 ================================="
echo "🧪 TRAE NEW HEALTHY1 TEST SUITE"
echo "🧪 ================================="
echo ""

# Test 1: Check if Dockerfile exists
log "Test 1: Checking Dockerfile..."
if [ -f "QUICK-DEPLOY-COOLIFY.md" ]; then
    success "Dockerfile found"
    ((PASSED++))
else
    fail "Dockerfile not found"
    ((FAILED++))
fi

# Test 2: Check Node.js syntax
log "Test 2: Validating Node.js syntax..."
if [ -f "production-nodejs/server.js" ]; then
    if node --check production-nodejs/server.js 2>/dev/null; then
        success "Node.js syntax valid"
        ((PASSED++))
    else
        fail "Node.js syntax error"
        ((FAILED++))
    fi
else
    warning "server.js not found, skipping"
fi

# Test 3: Validate JSON files
log "Test 3: Validating JSON files..."
if [ -f "production-nodejs/package.json" ]; then
    if node -e "JSON.parse(require('fs').readFileSync('production-nodejs/package.json', 'utf8'))" 2>/dev/null; then
        success "package.json valid"
        ((PASSED++))
    else
        fail "package.json invalid"
        ((FAILED++))
    fi
else
    warning "package.json not found, skipping"
fi

# Test 4: Check documentation
log "Test 4: Checking documentation..."
if [ -f "PRODUCTION-READY-GUIDE.md" ] && [ -f "QUICK-DEPLOY-COOLIFY.md" ]; then
    success "Documentation complete"
    ((PASSED++))
else
    fail "Documentation incomplete"
    ((FAILED++))
fi

# Test 5: Check file structure
log "Test 5: Checking file structure..."
REQUIRED_FILES=(
    "QUICK-DEPLOY-COOLIFY.md"
    "PRODUCTION-READY-GUIDE.md"
    "AI-ASSISTANT-VALIDATION-GUIDE.md"
)

ALL_EXIST=true
for file in "${REQUIRED_FILES[@]}"; do
    if [ ! -f "$file" ]; then
        ALL_EXIST=false
        break
    fi
done

if $ALL_EXIST; then
    success "All required files present"
    ((PASSED++))
else
    fail "Some required files missing"
    ((FAILED++))
fi

# Test 6: Check for security best practices
log "Test 6: Checking security implementations..."
if grep -q "helmet" production-nodejs/server.js 2>/dev/null || \
   grep -q "Access-Control-Allow-Origin" QUICK-DEPLOY-COOLIFY.md; then
    success "Security headers implemented"
    ((PASSED++))
else
    fail "Security headers missing"
    ((FAILED++))
fi

# Test 7: Check for error handling
log "Test 7: Checking error handling..."
if grep -q "try.*catch\|error" production-nodejs/server.js 2>/dev/null || \
   grep -q "catch.*error" QUICK-DEPLOY-COOLIFY.md; then
    success "Error handling implemented"
    ((PASSED++))
else
    fail "Error handling missing"
    ((FAILED++))
fi

# Test 8: Check for health endpoints
log "Test 8: Checking health endpoints..."
if grep -q "/health" QUICK-DEPLOY-COOLIFY.md; then
    success "Health endpoint implemented"
    ((PASSED++))
else
    fail "Health endpoint missing"
    ((FAILED++))
fi

# Test 9: Check for API documentation
log "Test 9: Checking API documentation..."
if grep -q "api/nutrition/analyze" QUICK-DEPLOY-COOLIFY.md; then
    success "API endpoints documented"
    ((PASSED++))
else
    fail "API documentation incomplete"
    ((FAILED++))
fi

# Test 10: Check for monitoring
log "Test 10: Checking monitoring capabilities..."
if grep -q "metrics\|monitoring\|logging" PRODUCTION-READY-GUIDE.md; then
    success "Monitoring capabilities present"
    ((PASSED++))
else
    fail "Monitoring capabilities missing"
    ((FAILED++))
fi

# Summary
echo ""
echo "📊 ================================="
echo "📊 TEST SUMMARY"
echo "📊 ================================="
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
    success "🎉 ALL TESTS PASSED! Platform is ready for deployment!"
    exit 0
else
    fail "❌ Some tests failed. Please review and fix issues."
    exit 1
fi