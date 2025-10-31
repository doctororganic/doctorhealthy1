#!/bin/bash

echo "🧪 TESTING ALL ENDPOINTS"
echo "========================"

BASE_URL="http://localhost:8080"
PASSED=0
FAILED=0

test_endpoint() {
    local method=$1
    local endpoint=$2
    local expected=$3
    local name=$4
    
    response=$(curl -s -o /dev/null -w "%{http_code}" -X $method "$BASE_URL$endpoint")
    if [ "$response" = "$expected" ]; then
        echo "✅ $name ($response)"
        ((PASSED++))
    else
        echo "❌ $name (Expected $expected, got $response)"
        ((FAILED++))
    fi
}

# Test endpoints
test_endpoint "GET" "/health" "200" "Health Check"
test_endpoint "GET" "/api/v1/info" "200" "API Info"
test_endpoint "GET" "/api/v1/foods" "200" "Foods List"
test_endpoint "GET" "/api/v1/diseases" "200" "Diseases List"

echo ""
echo "Results: $PASSED passed, $FAILED failed"
