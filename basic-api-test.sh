#!/bin/bash

# Basic API Testing Script for Nutrition Platform
# Tests API endpoints and captures application status

echo "🚀 Starting Basic API Testing..."
echo "================================="

BASE_URL="http://localhost:8080"
FRONTEND_URL="http://localhost:3000"

# Function to test an endpoint
test_endpoint() {
    local endpoint=$1
    local description=$2
    
    echo "📝 Testing: $description ($endpoint)"
    
    # Use curl to test the endpoint
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" "$BASE_URL$endpoint" -o "temp_response.json" 2>/dev/null)
    http_code=$(echo $response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    
    if [ "$http_code" = "200" ]; then
        echo "✅ PASSED: $description - HTTP $http_code"
        
        # Count items if it's a JSON array
        if command -v jq >/dev/null 2>&1; then
            count=$(jq 'length' temp_response.json 2>/dev/null || echo "N/A")
            echo "   📊 Data: $count items"
        fi
        
        return 0
    else
        echo "❌ FAILED: $description - HTTP $http_code"
        return 1
    fi
}

# Function to test frontend
test_frontend() {
    echo "📝 Testing: Frontend Accessibility ($FRONTEND_URL)"
    
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" "$FRONTEND_URL" -o "temp_frontend.html" 2>/dev/null)
    http_code=$(echo $response | tr -d '\n' | sed -e 's/HTTPSTATUS://')
    
    if [ "$http_code" = "200" ]; then
        echo "✅ PASSED: Frontend is accessible - HTTP $http_code"
        
        # Check if HTML contains expected content
        if grep -q "html\|HTML" temp_frontend.html; then
            echo "   📄 Valid HTML content detected"
        fi
        
        return 0
    else
        echo "❌ FAILED: Frontend not accessible - HTTP $http_code"
        return 1
    fi
}

# Function to test response time
test_response_time() {
    local endpoint=$1
    local description=$2
    
    echo "⏱️  Testing response time: $description"
    
    start_time=$(date +%s%3N)
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" "$BASE_URL$endpoint" -o /dev/null 2>/dev/null)
    end_time=$(date +%s%3N)
    
    response_time=$((end_time - start_time))
    http_code=$(echo $response | tr -d '\n' | sed -e 's/HTTPSTATUS://')
    
    if [ "$http_code" = "200" ]; then
        if [ "$response_time" -lt 3000 ]; then
            echo "✅ $description: ${response_time}ms (Good)"
        else
            echo "⚠️  $description: ${response_time}ms (Slow)"
        fi
        return 0
    else
        echo "❌ $description: Failed (HTTP $http_code)"
        return 1
    fi
}

# Initialize counters
total_tests=0
passed_tests=0
failed_tests=0

echo ""
echo "🔍 API HEALTH CHECKS"
echo "===================="

# Test basic health endpoints
total_tests=$((total_tests + 1))
if test_endpoint "/health" "Health Check"; then
    passed_tests=$((passed_tests + 1))
else
    failed_tests=$((failed_tests + 1))
fi

total_tests=$((total_tests + 1))
if test_endpoint "/api/info" "API Info"; then
    passed_tests=$((passed_tests + 1))
else
    failed_tests=$((failed_tests + 1))
fi

echo ""
echo "📊 NUTRITION DATA ENDPOINTS"
echo "============================"

# Test nutrition data endpoints
endpoints=(
    "/api/v1/metabolism:Metabolism Guide"
    "/api/v1/meal-plans:Meal Plans"
    "/api/v1/vitamins-minerals:Vitamins & Minerals"
    "/api/v1/workout-techniques:Workout Techniques"
    "/api/v1/calories:Calories Data"
    "/api/v1/skills:Skills Data"
    "/api/v1/diseases:Disease Data"
    "/api/v1/type-plans:Type Plans"
)

for endpoint_info in "${endpoints[@]}"; do
    endpoint="${endpoint_info%%:*}"
    description="${endpoint_info##*:}"
    
    total_tests=$((total_tests + 1))
    if test_endpoint "$endpoint" "$description"; then
        passed_tests=$((passed_tests + 1))
    else
        failed_tests=$((failed_tests + 1))
    fi
done

echo ""
echo "🏥 HEALTH SERVICE ENDPOINTS"
echo "==========================="

# Test health service endpoints
health_endpoints=(
    "/api/v1/health/conditions:Health Conditions"
    "/api/v1/health/tips:Health Tips"
)

for endpoint_info in "${health_endpoints[@]}"; do
    endpoint="${endpoint_info%%:*}"
    description="${endpoint_info##*:}"
    
    total_tests=$((total_tests + 1))
    if test_endpoint "$endpoint" "$description"; then
        passed_tests=$((passed_tests + 1))
    else
        failed_tests=$((failed_tests + 1))
    fi
done

echo ""
echo "🎨 FRONTEND ACCESSIBILITY"
echo "========================"

total_tests=$((total_tests + 1))
if test_frontend; then
    passed_tests=$((passed_tests + 1))
else
    failed_tests=$((failed_tests + 1))
fi

echo ""
echo "⏱️  PERFORMANCE TESTS"
echo "===================="

# Test response times
performance_endpoints=(
    "/health:Health Check Performance"
    "/api/info:API Info Performance"
    "/api/v1/calories:Calories Performance"
)

for endpoint_info in "${performance_endpoints[@]}"; do
    endpoint="${endpoint_info%%:*}"
    description="${endpoint_info##*:}"
    
    test_response_time "$endpoint" "$description"
done

echo ""
echo "🛡️  ERROR HANDLING"
echo "=================="

# Test 404 handling
echo "📝 Testing: 404 Error Handling"
response=$(curl -s -w "HTTPSTATUS:%{http_code}" "$BASE_URL/non-existent-endpoint" -o /dev/null 2>/dev/null)
http_code=$(echo $response | tr -d '\n' | sed -e 's/HTTPSTATUS://')

if [ "$http_code" = "404" ]; then
    echo "✅ PASSED: 404 Error Handling - Correctly returns 404"
    passed_tests=$((passed_tests + 1))
else
    echo "❌ FAILED: 404 Error Handling - Expected 404, got $http_code"
    failed_tests=$((failed_tests + 1))
fi
total_tests=$((total_tests + 1))

# Cleanup temporary files
rm -f temp_response.json temp_frontend.html

echo ""
echo "📊 TEST EXECUTION REPORT"
echo "======================="
echo "Total Tests: $total_tests"
echo "Passed: $passed_tests ✅"
echo "Failed: $failed_tests ❌"

if [ $total_tests -gt 0 ]; then
    success_rate=$(echo "scale=1; ($passed_tests * 100) / $total_tests" | bc 2>/dev/null || echo "N/A")
    echo "Success Rate: ${success_rate}%"
else
    echo "Success Rate: N/A"
fi

echo "================================="

# Create JSON report
timestamp=$(date -Iseconds)
json_report="{
  \"timestamp\": \"$timestamp\",
  \"summary\": {
    \"total\": $total_tests,
    \"passed\": $passed_tests,
    \"failed\": $failed_tests,
    \"successRate\": $success_rate
  },
  \"environment\": {
    \"frontend\": \"$FRONTEND_URL\",
    \"backend\": \"$BASE_URL\",
    \"platform\": \"$(uname -s)\",
    \"node_version\": \"$(node --version 2>/dev/null || echo 'N/A')\"
  }
}"

echo "$json_report" > basic-api-test-report.json
echo "📄 Detailed report saved to: basic-api-test-report.json"

# Final verdict
echo ""
if [ $failed_tests -eq 0 ]; then
    echo "🏆 SUCCESS: All tests passed! The application is working correctly."
    exit 0
elif [ $passed_tests -gt $failed_tests ]; then
    echo "⚠️  PARTIAL SUCCESS: Most tests passed, but some issues detected."
    exit 1
else
    echo "🚨 CRITICAL ISSUES: Many tests failed. Please check the application."
    exit 2
fi
