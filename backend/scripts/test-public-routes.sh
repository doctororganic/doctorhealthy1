#!/bin/bash

BASE_URL="http://localhost:8080"

echo "🧪 Testing public routes (no auth required)..."
echo ""

# Test disease routes
echo "1. Testing /api/v1/diseases"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/v1/diseases")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
if [ "$HTTP_CODE" = "200" ]; then
    echo "   ✅ Status: $HTTP_CODE"
    echo "$BODY" | jq -r '.status // "OK"' 2>/dev/null || echo "   Response received"
else
    echo "   ❌ Status: $HTTP_CODE"
    echo "$BODY"
fi
echo ""

# Test injury routes
echo "2. Testing /api/v1/injuries"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/v1/injuries")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
if [ "$HTTP_CODE" = "200" ]; then
    echo "   ✅ Status: $HTTP_CODE"
    echo "$BODY" | jq -r '.status // "OK"' 2>/dev/null || echo "   Response received"
else
    echo "   ❌ Status: $HTTP_CODE"
    echo "$BODY"
fi
echo ""

# Test vitamins routes
echo "3. Testing /api/v1/vitamins-minerals/vitamins"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/v1/vitamins-minerals/vitamins")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
if [ "$HTTP_CODE" = "200" ]; then
    echo "   ✅ Status: $HTTP_CODE"
    echo "$BODY" | jq -r '.status // "OK"' 2>/dev/null || echo "   Response received"
else
    echo "   ❌ Status: $HTTP_CODE"
    echo "$BODY"
fi
echo ""

# Test nutrition data routes
echo "4. Testing /api/v1/nutrition-data/recipes"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/v1/nutrition-data/recipes")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
if [ "$HTTP_CODE" = "200" ]; then
    echo "   ✅ Status: $HTTP_CODE"
    echo "$BODY" | jq -r '.status // "OK"' 2>/dev/null || echo "   Response received"
else
    echo "   ❌ Status: $HTTP_CODE"
    echo "$BODY"
fi
echo ""

# Test metabolism route
echo "5. Testing /api/v1/metabolism"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/v1/metabolism")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
if [ "$HTTP_CODE" = "200" ]; then
    echo "   ✅ Status: $HTTP_CODE"
    echo "$BODY" | jq -r '.status // "OK"' 2>/dev/null || echo "   Response received"
else
    echo "   ❌ Status: $HTTP_CODE"
    echo "$BODY"
fi
echo ""

echo "✅ Public route tests complete"

