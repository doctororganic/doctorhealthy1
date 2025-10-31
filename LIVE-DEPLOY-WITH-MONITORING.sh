#!/bin/bash
set -e

echo "🚀 LIVE DEPLOYMENT WITH REAL-TIME MONITORING"
echo "============================================="

# Run tests first
echo "Running pre-deployment tests..."
if ! ./LIVE-TEST-AND-FIX.sh; then
    echo "❌ Tests failed. Fix errors first."
    exit 1
fi

echo ""
echo "✅ Tests passed. Starting deployment..."
echo ""

# Generate credentials
echo "🔑 Generating secure credentials..."
DB_PASSWORD=$(openssl rand -hex 32)
JWT_SECRET=$(openssl rand -hex 64)
API_KEY_SECRET=$(openssl rand -hex 64)
REDIS_PASSWORD=$(openssl rand -hex 32)
ENCRYPTION_KEY=$(openssl rand -hex 32)

# Create .env.production
cat > .env.production << EOF
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
ALLOWED_ORIGINS=https://super.doctorhealthy1.com,https://www.super.doctorhealthy1.com
EOF

echo "✅ Credentials generated"

# Save credentials
cat > .credentials-backup.txt << EOF
SAVE THESE CREDENTIALS SECURELY:
================================
DB_PASSWORD=${DB_PASSWORD}
REDIS_PASSWORD=${REDIS_PASSWORD}
JWT_SECRET=${JWT_SECRET}
API_KEY_SECRET=${API_KEY_SECRET}
ENCRYPTION_KEY=${ENCRYPTION_KEY}
Generated: $(date)
EOF

echo "📋 Credentials saved to .credentials-backup.txt"

# Start deployment
echo ""
echo "🐳 Starting Docker deployment..."
docker-compose -f docker-compose.production.yml down -v 2>/dev/null || true
docker-compose -f docker-compose.production.yml up -d --build

# Monitor startup
echo ""
echo "⏳ Waiting for services to start..."

# Wait for postgres
echo "Waiting for PostgreSQL..."
for i in {1..30}; do
    if docker-compose -f docker-compose.production.yml exec -T postgres pg_isready -U nutrition_user > /dev/null 2>&1; then
        echo "✅ PostgreSQL ready"
        break
    fi
    echo "  Attempt $i/30..."
    sleep 2
done

# Wait for redis
echo "Waiting for Redis..."
for i in {1..30}; do
    if docker-compose -f docker-compose.production.yml exec -T redis redis-cli ping > /dev/null 2>&1; then
        echo "✅ Redis ready"
        break
    fi
    echo "  Attempt $i/30..."
    sleep 2
done

# Wait for backend
echo "Waiting for Backend..."
for i in {1..60}; do
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        echo "✅ Backend ready"
        break
    fi
    echo "  Attempt $i/60..."
    sleep 2
done

# Wait for frontend
echo "Waiting for Frontend..."
for i in {1..60}; do
    if curl -f http://localhost:3000 > /dev/null 2>&1; then
        echo "✅ Frontend ready"
        break
    fi
    echo "  Attempt $i/60..."
    sleep 2
done

# Run live tests
echo ""
echo "🧪 Running live integration tests..."

# Test backend health
if curl -f http://localhost:8080/health; then
    echo "✅ Backend health check passed"
else
    echo "❌ Backend health check failed"
    docker-compose -f docker-compose.production.yml logs backend
    exit 1
fi

# Test frontend
if curl -f http://localhost:3000 > /dev/null 2>&1; then
    echo "✅ Frontend accessible"
else
    echo "❌ Frontend not accessible"
    docker-compose -f docker-compose.production.yml logs frontend
    exit 1
fi

# Show status
echo ""
echo "=================================="
echo "✅ DEPLOYMENT SUCCESSFUL!"
echo "=================================="
echo ""
echo "🌐 Services:"
echo "  Frontend: http://localhost:3000"
echo "  Backend:  http://localhost:8080"
echo "  Health:   http://localhost:8080/health"
echo ""
echo "📊 Monitor logs:"
echo "  docker-compose -f docker-compose.production.yml logs -f"
echo ""
echo "🔍 Check status:"
echo "  docker-compose -f docker-compose.production.yml ps"
echo ""
echo "📋 Credentials saved in: .credentials-backup.txt"
echo ""

# Start live monitoring in background
./monitor-deployment.sh &
echo "🔄 Live monitoring started (PID: $!)"
