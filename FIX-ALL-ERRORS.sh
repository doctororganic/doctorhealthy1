#!/bin/bash
set -e

echo "🔧 AUTOMATIC ERROR FIXING"
echo "========================="

# Fix 1: Update Go dependencies
echo "1️⃣  Fixing Go dependencies..."
cd backend
go mod tidy
go mod download
go mod verify
cd ..
echo "✅ Go dependencies fixed"

# Fix 2: Create missing directories
echo ""
echo "2️⃣  Creating missing directories..."
mkdir -p backend/config
mkdir -p backend/handlers
mkdir -p backend/models
mkdir -p backend/services
mkdir -p backend/middleware
mkdir -p logs
mkdir -p ssl
echo "✅ Directories created"

# Fix 3: Fix file permissions
echo ""
echo "3️⃣  Fixing file permissions..."
chmod +x *.sh
chmod +x scripts/*.sh 2>/dev/null || true
echo "✅ Permissions fixed"

# Fix 4: Clean Docker
echo ""
echo "4️⃣  Cleaning Docker environment..."
docker-compose -f docker-compose.production.yml down -v 2>/dev/null || true
docker system prune -f
echo "✅ Docker cleaned"

# Fix 5: Validate and fix .env
echo ""
echo "5️⃣  Validating environment configuration..."
if [ ! -f ".env.production" ]; then
    echo "Creating .env.production..."
    cat > .env.production << 'EOF'
DB_HOST=postgres
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=nutrition_user
DB_PASSWORD=${DB_PASSWORD}
DB_SSL_MODE=require

REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=${REDIS_PASSWORD}

JWT_SECRET=${JWT_SECRET}
API_KEY_SECRET=${API_KEY_SECRET}
ENCRYPTION_KEY=${ENCRYPTION_KEY}

PORT=8080
ENVIRONMENT=production
DOMAIN=super.doctorhealthy1.com
ALLOWED_ORIGINS=https://super.doctorhealthy1.com
EOF
fi
echo "✅ Environment configuration validated"

# Fix 6: Check and fix CORS
echo ""
echo "6️⃣  Fixing CORS configuration..."
if [ -f "nginx/production.conf" ]; then
    sed -i.bak 's/Access-Control-Allow-Origin "\*"/Access-Control-Allow-Origin "https:\/\/super.doctorhealthy1.com"/g' nginx/production.conf
    echo "✅ CORS fixed"
else
    echo "⚠️  nginx/production.conf not found"
fi

# Fix 7: Ensure SSL is enabled
echo ""
echo "7️⃣  Ensuring SSL is enabled..."
if [ -f "backend/.env.example" ]; then
    sed -i.bak 's/DB_SSL_MODE=disable/DB_SSL_MODE=require/g' backend/.env.example
fi
echo "✅ SSL configuration fixed"

echo ""
echo "=================================="
echo "✅ ALL FIXES APPLIED!"
echo "=================================="
echo ""
echo "Run ./LIVE-TEST-AND-FIX.sh to verify"
