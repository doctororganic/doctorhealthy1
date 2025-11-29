#!/bin/bash

BASE_URL="http://localhost:8080"
FAILED=0
PASSED=0

echo "🧪 Testing All JSON Data Endpoints"
echo "===================================="
echo ""

# Function to test endpoint
test_endpoint() {
    local name=$1
    local url=$2
    local expected_status=${3:-200}
    
    echo "Testing: $name"
    echo "  URL: $url"
    
    RESPONSE=$(curl -s -w "\n%{http_code}" "$url" 2>&1)
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | sed '$d')
    
    if [ "$HTTP_CODE" = "$expected_status" ]; then
        echo "  ✅ Status: $HTTP_CODE"
        # Try to parse JSON and show status
        STATUS=$(echo "$BODY" | jq -r '.status // empty' 2>/dev/null)
        if [ -n "$STATUS" ]; then
            echo "  📊 Response status: $STATUS"
        fi
        # Show data count if available
        COUNT=$(echo "$BODY" | jq -r '.total // .data | length // empty' 2>/dev/null)
        if [ -n "$COUNT" ] && [ "$COUNT" != "null" ]; then
            echo "  📈 Items found: $COUNT"
        fi
        PASSED=$((PASSED + 1))
    else
        echo "  ❌ Status: $HTTP_CODE (expected $expected_status)"
        echo "  Response: $(echo "$BODY" | head -c 200)"
        FAILED=$((FAILED + 1))
    fi
    echo ""
}

# Test health endpoint first
test_endpoint "Health Check" "$BASE_URL/health"

# Test disease endpoints
test_endpoint "Diseases List" "$BASE_URL/api/v1/diseases/"
test_endpoint "Disease Categories" "$BASE_URL/api/v1/diseases/categories"
test_endpoint "Disease Search" "$BASE_URL/api/v1/diseases/search?search=diabetes"

# Test injury endpoints
test_endpoint "Injuries List" "$BASE_URL/api/v1/injuries/"
test_endpoint "Injury Categories" "$BASE_URL/api/v1/injuries/categories"
test_endpoint "Injury Search" "$BASE_URL/api/v1/injuries/search?search=sprain"

# Test vitamins/minerals endpoints
test_endpoint "Vitamins List" "$BASE_URL/api/v1/vitamins-minerals/vitamins"
test_endpoint "Supplements List" "$BASE_URL/api/v1/vitamins-minerals/supplements"
test_endpoint "Vitamins Search" "$BASE_URL/api/v1/vitamins-minerals/search?q=vitamin"
test_endpoint "Weight Loss Drugs" "$BASE_URL/api/v1/vitamins-minerals/weight-loss-drugs"
test_endpoint "Drug Categories" "$BASE_URL/api/v1/vitamins-minerals/drug-categories"

# Test nutrition data endpoints
test_endpoint "Recipes" "$BASE_URL/api/v1/nutrition-data/recipes"
test_endpoint "Workouts" "$BASE_URL/api/v1/nutrition-data/workouts"
test_endpoint "Complaints" "$BASE_URL/api/v1/nutrition-data/complaints"
test_endpoint "Metabolism" "$BASE_URL/api/v1/nutrition-data/metabolism"
test_endpoint "Drugs-Nutrition" "$BASE_URL/api/v1/nutrition-data/drugs-nutrition"

# Test legacy nutrition endpoints
test_endpoint "Metabolism (legacy)" "$BASE_URL/api/v1/metabolism"
test_endpoint "Workout Techniques" "$BASE_URL/api/v1/workout-techniques"
test_endpoint "Meal Plans" "$BASE_URL/api/v1/meal-plans"
test_endpoint "Drugs-Nutrition (legacy)" "$BASE_URL/api/v1/drugs-nutrition"

# Test validation endpoints
test_endpoint "Validate All Files" "$BASE_URL/api/v1/validation/all"
test_endpoint "Validate Recipes" "$BASE_URL/api/v1/validation/file/qwen-recipes.json"

echo "===================================="
echo "📊 Test Summary:"
echo "  ✅ Passed: $PASSED"
echo "  ❌ Failed: $FAILED"
echo "  📈 Total: $((PASSED + FAILED))"
echo ""

if [ $FAILED -eq 0 ]; then
    echo "🎉 All tests passed!"
    exit 0
else
    echo "⚠️  Some tests failed"
    exit 1
fi

