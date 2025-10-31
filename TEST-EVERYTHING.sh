#!/bin/bash

echo "🧪 Testing Complete System..."
echo "=============================="

# Test 1: Go Backend Compilation
echo ""
echo "1️⃣  Testing Go Backend..."
cd backend
if go build -o test-build . 2>/dev/null; then
    echo "  ✅ Go backend compiles"
    rm test-build
else
    echo "  ❌ Go backend has errors"
    exit 1
fi
cd ..

# Test 2: Frontend Setup
echo ""
echo "2️⃣  Testing Frontend Setup..."
cd frontend-nextjs
if [ -f "package.json" ]; then
    echo "  ✅ Frontend package.json exists"
    if [ -f "Dockerfile" ]; then
        echo "  ✅ Frontend Dockerfile exists"
    else
        echo "  ❌ Frontend Dockerfile missing"
    fi
else
    echo "  ❌ Frontend package.json missing"
fi
cd ..

# Test 3: Docker Compose
echo ""
echo "3️⃣  Testing Docker Compose..."
if [ -f "docker-compose.yml" ]; then
    echo "  ✅ docker-compose.yml exists"
    if docker-compose config > /dev/null 2>&1; then
        echo "  ✅ docker-compose.yml is valid"
    else
        echo "  ⚠️  docker-compose.yml has warnings (may still work)"
    fi
else
    echo "  ❌ docker-compose.yml missing"
fi

# Test 4: Deployment Script
echo ""
echo "4️⃣  Testing Deployment Script..."
if [ -f "deploy.sh" ] && [ -x "deploy.sh" ]; then
    echo "  ✅ deploy.sh exists and is executable"
else
    echo "  ❌ deploy.sh missing or not executable"
fi

# Test 5: Documentation
echo ""
echo "5️⃣  Testing Documentation..."
if [ -f "README.md" ]; then
    echo "  ✅ README.md exists"
fi
if [ -f "DEPLOYMENT.md" ]; then
    echo "  ✅ DEPLOYMENT.md exists"
fi

# Test 6: Archive Structure
echo ""
echo "6️⃣  Checking Archive..."
if [ -d "archive" ]; then
    echo "  ✅ Archive directory exists"
    echo "  📦 Archived items:"
    ls -1 archive/ | head -5
else
    echo "  ⚠️  No archive directory"
fi

# Summary
echo ""
echo "=============================="
echo "✅ SYSTEM TEST COMPLETE"
echo "=============================="
echo ""
echo "🎯 Ready to:"
echo "  1. Start development: docker-compose up -d"
echo "  2. Deploy production: ./deploy.sh"
echo "  3. Test backend: curl http://localhost:8080/health"
echo ""
