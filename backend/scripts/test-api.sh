#!/bin/bash
# Quick API Testing Script
# Usage: ./scripts/test-api.sh [token]

BASE_URL="http://localhost:8080"
TOKEN="${1:-}"

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo -e "\n${YELLOW}Testing: $description${NC}"
    echo "  $method $endpoint"
    
    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X $method \
            -H "Authorization: Bearer $TOKEN" \
            "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X $method \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $TOKEN" \
            -d "$data" \
            "$BASE_URL$endpoint")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo -e "${GREEN}✅ Success ($http_code)${NC}"
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
    else
        echo -e "${RED}❌ Failed ($http_code)${NC}"
        echo "$body"
    fi
}

# Health check (no auth needed)
echo -e "${YELLOW}=== Health Check ===${NC}"
test_endpoint "GET" "/health" "" "Health check"

# If no token provided, try to get one
if [ -z "$TOKEN" ]; then
    echo -e "\n${YELLOW}=== Getting Auth Token ===${NC}"
    
    # Register test user
    register_response=$(curl -s -X POST "$BASE_URL/api/v1/auth/register" \
        -H "Content-Type: application/json" \
        -d '{
            "email": "test'$(date +%s)'@example.com",
            "password": "password123",
            "first_name": "Test",
            "last_name": "User"
        }')
    
    # Login
    login_response=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d '{
            "email": "test'$(date +%s)'@example.com",
            "password": "password123"
        }')
    
    TOKEN=$(echo "$login_response" | jq -r '.data.access_token' 2>/dev/null)
    
    if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
        echo -e "${RED}❌ Failed to get token${NC}"
        echo "Response: $login_response"
        exit 1
    fi
    
    echo -e "${GREEN}✅ Token obtained${NC}"
fi

# Test action endpoints
echo -e "\n${YELLOW}=== Testing Action Endpoints ===${NC}"

# Track measurement
test_endpoint "POST" "/api/v1/actions/track-measurement" \
    '{"waist": 85.5, "measurement_date": "'$(date +%Y-%m-%d)'"}' \
    "Track measurement"

# Progress summary
test_endpoint "GET" "/api/v1/actions/progress-summary?days=30" \
    "" "Get progress summary"

# Generate meal plan
test_endpoint "POST" "/api/v1/actions/generate-meal-plan" \
    '{"goal": "weight_loss", "target_calories": 2000, "duration": 7}' \
    "Generate meal plan"

# Generate workout
test_endpoint "POST" "/api/v1/actions/generate-workout" \
    '{"goal": "weight_loss", "duration": 30, "difficulty": "intermediate"}' \
    "Generate workout"

echo -e "\n${GREEN}✅ API testing complete!${NC}"