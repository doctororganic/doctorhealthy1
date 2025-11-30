# 🚀 Coolify Deployment Guide

Complete guide for deploying the Nutrition Platform on Coolify.

---

## 📋 Prerequisites

- Coolify instance running
- GitHub repository access
- Domain name configured

---

## 🎯 Quick Deployment Steps

### Step 1: Create New Application in Coolify

1. Login to your Coolify dashboard
2. Navigate to **Projects** → Select or create project
3. Click **New Application**
4. Choose **Dockerfile** type

### Step 2: Connect GitHub Repository

1. **Repository**: `doctororganic/doctorhealthy1`
2. **Branch**: `main`
3. **Dockerfile Path**: `nutrition-platform/docker-compose.production.yml` (or use individual Dockerfiles)
4. **Build Pack**: Docker

### Step 3: Configure Application Settings

**Basic Settings:**
- **Name**: `nutrition-platform`
- **Type**: Docker Compose (or Dockerfile)
- **Port**: `8080` (Backend) or `3000` (Frontend)

**For Backend:**
- **Dockerfile Path**: `nutrition-platform/backend/Dockerfile.go`
- **Context**: `nutrition-platform/backend`
- **Port**: `8080`

**For Frontend:**
- **Dockerfile Path**: `nutrition-platform/frontend-nextjs/Dockerfile`
- **Context**: `nutrition-platform/frontend-nextjs`
- **Port**: `3000`

### Step 4: Environment Variables

Add these environment variables in Coolify:

```env
# Database (if using external database)
DB_HOST=your-db-host
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=postgres
DB_PASSWORD=your-secure-password

# Redis (if using external Redis)
REDIS_HOST=your-redis-host
REDIS_PORT=6379
REDIS_PASSWORD=

# Server
PORT=8080
ENV=production
DOMAIN=yourdomain.com

# Security (Generate secure values)
JWT_SECRET=generate-random-64-chars
API_KEY_SECRET=generate-random-64-chars
SESSION_SECRET=generate-random-32-chars

# CORS
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# Frontend API URL
NEXT_PUBLIC_API_URL=https://api.yourdomain.com
```

### Step 5: Health Check Configuration

**Backend Health Check:**
- **Path**: `/health`
- **Port**: `8080`
- **Interval**: `30s`
- **Timeout**: `10s`
- **Retries**: `3`

**Frontend Health Check:**
- **Path**: `/`
- **Port**: `3000`
- **Interval**: `30s`
- **Timeout**: `10s`
- **Retries**: `3`

### Step 6: Domain Configuration

1. Add your domain in Coolify
2. Enable SSL (Let's Encrypt)
3. Configure subdomains:
   - `api.yourdomain.com` → Backend (port 8080)
   - `yourdomain.com` → Frontend (port 3000)

### Step 7: Deploy

1. Click **Deploy** button
2. Monitor build logs
3. Wait for deployment to complete
4. Verify health checks pass

---

## 🐳 Docker Compose Deployment (Recommended)

For full stack deployment, use Docker Compose:

### Configuration

1. **Type**: Docker Compose
2. **Compose File**: `nutrition-platform/docker-compose.production.yml`
3. **Environment File**: Create `.env.production` in Coolify

### Services Included

- `backend` - Go API server
- `frontend` - Next.js application
- `postgres` - Database (or use external)
- `redis` - Cache (or use external)
- `nginx` - Reverse proxy

---

## 📝 Coolify-Specific Files

### Dockerfile for Backend (Coolify)

```dockerfile
# File: backend/Dockerfile.go
FROM golang:1.21-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nutrition-platform .

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates curl
COPY --from=builder /app/nutrition-platform .
COPY --from=builder /app/config ./config
COPY --from=builder /app/data ./data
RUN addgroup -g 1001 -S appuser && adduser -S appuser -u 1001
RUN chown -R appuser:appuser /app
USER appuser
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1
CMD ["./nutrition-platform"]
```

### Dockerfile for Frontend (Coolify)

```dockerfile
# File: frontend-nextjs/Dockerfile
FROM node:18-alpine AS base
FROM base AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm ci

FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
ENV NEXT_TELEMETRY_DISABLED 1
ENV NODE_ENV production
RUN npm run build

FROM base AS runner
WORKDIR /app
ENV NODE_ENV production
ENV NEXT_TELEMETRY_DISABLED 1
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs
COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static
USER nextjs
EXPOSE 3000
ENV PORT 3000
ENV HOSTNAME "0.0.0.0"
CMD ["node", "server.js"]
```

---

## 🔧 Post-Deployment Configuration

### 1. Verify Deployment

```bash
# Check backend
curl https://api.yourdomain.com/health

# Check frontend
curl https://yourdomain.com
```

### 2. Database Setup

If using internal PostgreSQL:
- Database will be created automatically
- Run migrations if needed: `docker exec -it container_name ./migrations`

### 3. SSL Configuration

- Coolify handles SSL automatically via Let's Encrypt
- Ensure domain DNS is configured correctly
- SSL will auto-renew

---

## 📊 Monitoring

### Health Checks

- Backend: `https://api.yourdomain.com/health`
- Frontend: `https://yourdomain.com`

### Logs

View logs in Coolify dashboard:
- Application logs
- Build logs
- Deployment logs

---

## 🔄 Update Application

### Automatic Updates

1. Push changes to `main` branch
2. Coolify will detect changes
3. Automatic rebuild and redeploy

### Manual Update

1. Go to application in Coolify
2. Click **Redeploy**
3. Monitor deployment logs

---

## 🆘 Troubleshooting

### Build Fails

- Check Dockerfile path
- Verify build context
- Check build logs in Coolify

### Application Not Starting

- Check environment variables
- Verify health check configuration
- Review application logs

### Database Connection Issues

- Verify database credentials
- Check network connectivity
- Ensure database is accessible

---

## ✅ Deployment Checklist

- [ ] Repository connected
- [ ] Dockerfile configured
- [ ] Environment variables set
- [ ] Health checks configured
- [ ] Domain configured
- [ ] SSL enabled
- [ ] Database configured
- [ ] Redis configured (if needed)
- [ ] CORS configured
- [ ] Deployment successful
- [ ] Health checks passing

---

## 🎯 Quick Reference

**Repository**: `https://github.com/doctororganic/doctorhealthy1`

**Backend Dockerfile**: `nutrition-platform/backend/Dockerfile.go`

**Frontend Dockerfile**: `nutrition-platform/frontend-nextjs/Dockerfile`

**Docker Compose**: `nutrition-platform/docker-compose.production.yml`

**Environment Template**: `.env.production.example`

---

**Your application is ready for Coolify deployment! 🚀**

