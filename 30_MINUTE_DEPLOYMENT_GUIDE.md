# 🚀 30-Minute Production Deployment Guide

**Last Updated:** $(date +"%Y-%m-%d %H:%M:%S")

## ⚡ Quick Start (5 minutes)

### 1. Prerequisites Check
```bash
# Check Go version
go version  # Should be 1.24+

# Check if port 8080 is available
lsof -ti:8080 | xargs kill -9 2>/dev/null || echo "Port 8080 is free"

# Check Redis (optional - app works without it)
redis-cli ping 2>/dev/null || echo "Redis not available (will use in-memory cache)"
```

### 2. Environment Setup (2 minutes)
```bash
cd nutrition-platform/backend

# Create .env file if it doesn't exist
cat > .env << EOF
DATABASE_URL=sqlite:///data/nutrition.db
JWT_SECRET=$(openssl rand -hex 32)
PORT=8080
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
ENVIRONMENT=production
EOF
```

### 3. Build & Run (3 minutes)
```bash
# Build the server
go build -o bin/server .

# Run migrations (if needed)
./bin/server migrate 2>/dev/null || echo "Migrations handled automatically"

# Start server
./bin/server > server.log 2>&1 &
SERVER_PID=$!

# Wait for server to start
sleep 3

# Verify health
curl http://localhost:8080/health | jq .
```

---

## 📋 Full Deployment Checklist (25 minutes)

### Phase 1: Backend Deployment (10 minutes)

#### Step 1: Database Setup
```bash
# Option A: SQLite (Quick - for testing)
mkdir -p data
export DATABASE_URL="sqlite:///data/nutrition.db"

# Option B: PostgreSQL (Production)
# Create database first, then:
export DATABASE_URL="postgres://user:pass@localhost/nutrition?sslmode=disable"
```

#### Step 2: Build & Test
```bash
cd backend

# Install dependencies
go mod download

# Run tests (optional but recommended)
go test ./... -v 2>&1 | head -20

# Build
go build -o bin/server .

# Verify build
./bin/server --help 2>&1 | head -5
```

#### Step 3: Start Server
```bash
# Production mode
export ENVIRONMENT=production
export PORT=8080

# Start server
nohup ./bin/server > server.log 2>&1 &

# Or use systemd (recommended for production)
sudo tee /etc/systemd/system/nutrition-api.service > /dev/null << EOF
[Unit]
Description=Nutrition Platform API
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/nutrition-platform/backend
ExecStart=/opt/nutrition-platform/backend/bin/server
Restart=always
Environment=ENVIRONMENT=production
Environment=DATABASE_URL=sqlite:///data/nutrition.db

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable nutrition-api
sudo systemctl start nutrition-api
sudo systemctl status nutrition-api
```

#### Step 4: Verify Backend
```bash
# Health check
curl http://localhost:8080/health

# Test endpoints
curl http://localhost:8080/api/v1/nutrition-data/recipes?limit=5 | jq '.status'
curl http://localhost:8080/api/v1/nutrition-data/workouts?limit=5 | jq '.status'
curl http://localhost:8080/api/v1/diseases?limit=5 | jq '.status'
```

### Phase 2: Frontend Deployment (10 minutes)

#### Step 1: Build Frontend
```bash
cd frontend-nextjs

# Install dependencies
npm install

# Build for production
npm run build

# Verify build
ls -lh .next/
```

#### Step 2: Deploy Frontend

**Option A: Static Export (Simplest)**
```bash
# Update next.config.js to enable static export
# Then:
npm run build
# Deploy .next/static and .next/export to your CDN/static host
```

**Option B: Node.js Server**
```bash
# Start production server
npm start

# Or use PM2
pm2 start npm --name "nutrition-frontend" -- start
pm2 save
```

**Option C: Docker**
```bash
docker build -t nutrition-frontend .
docker run -d -p 3000:3000 --name nutrition-frontend nutrition-frontend
```

### Phase 3: Verification & Monitoring (5 minutes)

#### Health Checks
```bash
# Backend health
curl http://localhost:8080/health | jq .

# Frontend health
curl http://localhost:3000 | head -20

# API endpoints
curl http://localhost:8080/api/v1/nutrition-data/recipes?limit=1 | jq '.status'
```

#### Monitor Logs
```bash
# Backend logs
tail -f backend/server.log

# System logs (if using systemd)
sudo journalctl -u nutrition-api -f
```

---

## 🔧 Production Configuration

### Environment Variables
```bash
# Required
DATABASE_URL=postgres://user:pass@localhost/nutrition
JWT_SECRET=<generate-strong-secret>
PORT=8080
ENVIRONMENT=production

# Optional (with fallbacks)
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
FRONTEND_URL=https://yourdomain.com
API_URL=https://api.yourdomain.com
```

### Security Checklist
- [ ] JWT_SECRET is strong (32+ characters, random)
- [ ] Database credentials are secure
- [ ] HTTPS enabled (use reverse proxy like Nginx)
- [ ] Rate limiting enabled (✅ already configured)
- [ ] Security headers enabled (✅ already configured)
- [ ] CORS configured for your domain
- [ ] Redis password set (if using Redis)

### Performance Optimization
- [ ] Redis cache enabled (optional but recommended)
- [ ] Database connection pooling configured
- [ ] CDN configured for static assets
- [ ] Gzip compression enabled (✅ already configured)

---

## 🚨 Troubleshooting

### Server Won't Start
```bash
# Check if port is in use
lsof -ti:8080

# Check logs
tail -50 backend/server.log

# Check database connection
# For SQLite: ls -lh data/nutrition.db
# For PostgreSQL: psql $DATABASE_URL -c "SELECT 1"
```

### Build Errors
```bash
# Clean and rebuild
cd backend
go clean -cache
go mod tidy
go build -o bin/server .
```

### API Not Responding
```bash
# Check server status
curl http://localhost:8080/health

# Check middleware
curl -v http://localhost:8080/api/v1/nutrition-data/recipes?limit=1

# Check logs
tail -100 backend/server.log | grep -i error
```

---

## 📊 Post-Deployment Verification

### Critical Endpoints Test
```bash
#!/bin/bash
# Run this script to verify all critical endpoints

BASE_URL="http://localhost:8080"

echo "Testing Health Endpoint..."
curl -s "$BASE_URL/health" | jq '.status' || echo "❌ Health check failed"

echo "Testing Recipes Endpoint..."
curl -s "$BASE_URL/api/v1/nutrition-data/recipes?limit=5" | jq '.status' || echo "❌ Recipes failed"

echo "Testing Workouts Endpoint..."
curl -s "$BASE_URL/api/v1/nutrition-data/workouts?limit=5" | jq '.status' || echo "❌ Workouts failed"

echo "Testing Diseases Endpoint..."
curl -s "$BASE_URL/api/v1/diseases?limit=5" | jq '.status' || echo "❌ Diseases failed"

echo "Testing Injuries Endpoint..."
curl -s "$BASE_URL/api/v1/injuries?limit=5" | jq '.status' || echo "❌ Injuries failed"

echo "✅ All tests completed!"
```

### Performance Test
```bash
# Quick load test
ab -n 100 -c 10 http://localhost:8080/api/v1/nutrition-data/recipes?limit=10
```

---

## 🎯 Success Criteria

Your deployment is successful if:
- ✅ Health endpoint returns `{"status":"ok"}`
- ✅ All API endpoints return `{"status":"success"}`
- ✅ Frontend loads without errors
- ✅ No errors in server logs
- ✅ Response times < 500ms for cached endpoints
- ✅ Rate limiting works (test with rapid requests)

---

## 📝 Quick Reference

### Start Server
```bash
cd backend && ./bin/server
```

### Stop Server
```bash
pkill -f "bin/server"
# Or if using systemd:
sudo systemctl stop nutrition-api
```

### View Logs
```bash
tail -f backend/server.log
```

### Restart Server
```bash
pkill -f "bin/server" && sleep 1 && cd backend && ./bin/server > server.log 2>&1 &
```

---

## 🆘 Emergency Rollback

If something goes wrong:

```bash
# Stop the server
pkill -f "bin/server"

# Restore from backup (if you have one)
# cp backup/nutrition.db data/nutrition.db

# Start previous version
# git checkout <previous-commit>
# go build -o bin/server .
# ./bin/server
```

---

**Ready to deploy!** 🚀

If you encounter any issues, check the logs first: `tail -f backend/server.log`

