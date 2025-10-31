#!/bin/bash

echo "Running smoke tests..."

# Test 1: Backend health
echo "Test 1: Backend health check..."
response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)
if [ "$response" != "200" ]; then
    echo "❌ Backend health check failed (HTTP $response)"
    exit 1
fi
echo "✅ Backend health check passed"

# Test 2: Frontend loads
echo "Test 2: Frontend loads..."
response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:3000)
if [ "$response" != "200" ]; then
    echo "❌ Frontend load failed (HTTP $response)"
    exit 1
fi
echo "✅ Frontend loads"

# Test 3: API responds
echo "Test 3: API info endpoint..."
response=$(curl -s http://localhost:8080/api/v1/info 2>/dev/null | grep -o '"status":"active"' || echo "")
if [ -z "$response" ]; then
    echo "❌ API info endpoint failed"
    exit 1
fi
echo "✅ API responds correctly"

# Test 4: Database connection
echo "Test 4: Database connection..."
docker-compose exec -T postgres pg_isready -U nutrition_user > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "❌ Database connection failed"
    exit 1
fi
echo "✅ Database connected"

# Test 5: Redis connection
echo "Test 5: Redis connection..."
docker-compose exec -T redis redis-cli ping > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "❌ Redis connection failed"
    exit 1
fi
echo "✅ Redis connected"

echo ""
echo "All smoke tests passed! ✅"
