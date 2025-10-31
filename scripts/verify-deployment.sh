#!/bin/bash

echo "Verifying deployment..."

PASSED=0
FAILED=0

# Function to test endpoint
test_endpoint() {
    local url=$1
    local expected=$2
    local name=$3
    
    response=$(curl -s -o /dev/null -w "%{http_code}" "$url")
    if [ "$response" = "$expected" ]; then
        echo "✅ $name"
        ((PASSED++))
    else
        echo "❌ $name (Expected $expected, got $response)"
        ((FAILED++))
    fi
}

# Test endpoints
test_endpoint "http://localhost:8080/health" "200" "Backend health"
test_endpoint "http://localhost:3000" "200" "Frontend"
test_endpoint "http://localhost:8080/api/v1/info" "200" "API info"

# Test CORS
echo "Testing CORS..."
cors_response=$(curl -s -H "Origin: http://localhost:3000" \
    -H "Access-Control-Request-Method: POST" \
    -X OPTIONS http://localhost:8080/api/v1/nutrition/analyze \
    -I 2>/dev/null | grep -i "access-control-allow-origin" || echo "")

if [ -n "$cors_response" ]; then
    echo "✅ CORS configured"
    ((PASSED++))
else
    echo "❌ CORS not configured"
    ((FAILED++))
fi

# Summary
echo ""
echo "Verification complete:"
echo "  Passed: $PASSED"
echo "  Failed: $FAILED"

if [ $FAILED -eq 0 ]; then
    echo "✅ Deployment verified successfully"
    exit 0
else
    echo "❌ Deployment verification failed"
    exit 1
fi
