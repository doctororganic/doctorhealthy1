# 🚀 Coolify Step-by-Step Deployment Guide

## Current Status Check

Our automated deployment started successfully, but let's complete it manually to ensure everything works perfectly.

## Step 1: Access Your Coolify Application

**Go to:** https://api.doctorhealthy1.com/project/us4gwgo8o4o4wocgo0k80kg0/environment/w8ksg0gk8sg8ogckwg4ggsc8/application/hcw0gc8wcwk440gw4c88408o

## Step 2: Check Current Deployment Status

1. **Click on your application** in Coolify
2. **Check the "Deployments" tab** - you should see a recent deployment
3. **Look at the logs** to see if there are any errors

## Step 3: Upload Your Application Code

### Option A: Upload ZIP File
1. **Create a ZIP file** of the `nutrition-platform` folder
2. **In Coolify**, go to "Source" tab
3. **Upload the ZIP file**

### Option B: Connect Git Repository (Recommended)
1. **Push your code to GitHub/GitLab**
2. **In Coolify**, go to "Source" tab
3. **Connect your Git repository**
4. **Set branch to "main"**

## Step 4: Configure Application Settings

### Basic Configuration:
```
Name: trae-new-healthy1
Domain: super.doctorhealthy1.com
Port: 8080
Build Pack: dockerfile
Dockerfile Location: backend/Dockerfile
```

## Step 5: Add Environment Variables

**Go to "Environment Variables" tab and add these:**

```bash
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
ENVIRONMENT=production
DEBUG=false

# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=nutrition_user
DB_PASSWORD=a898142298f6dd1f9a385bdcfb7e5cd6854642fbfadace6a2192eed818a6f218
DB_SSL_MODE=require

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=acadf682123532aef5c7a1ab8b01ca2003864d3055aba4bebe4a234161c75018

# Security (Use the SECURE generated values)
JWT_SECRET=4f920673e4f5e21ebdef27c631fa43fbdfb35ded892808b5d0bf05dcd317166c4e724a0988d4b6d624341a62b15f0dea9b6fd380ccc09eb3ed0f99e6f315f13a
API_KEY_SECRET=a898142298f6dd1f9a385bdcfb7e5cd6854642fbfadace6a2192eed818a6f218
ENCRYPTION_KEY=a898142298f6dd1f9a385bdcfb7e5cd6854642fbfadace6a2192eed818a6f218

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://super.doctorhealthy1.com,https://www.super.doctorhealthy1.com

# Performance & Features
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60s
SECURITY_HEADERS_ENABLED=true
COMPRESSION_ENABLED=true
METRICS_ENABLED=true
LOG_LEVEL=info
LOG_FORMAT=json
DATA_PATH=./data
NUTRITION_DATA_PATH=./
UPLOAD_PATH=./uploads
RELIGIOUS_FILTER_ENABLED=true
FILTER_ALCOHOL=true
FILTER_PORK=true
DEFAULT_LANGUAGE=en
SUPPORTED_LANGUAGES=en,ar
HEALTH_CHECK_ENABLED=true
```

## Step 6: Create Database Services

### PostgreSQL Database:
1. **In Coolify**, go to "Services" → "Add Service"
2. **Select "PostgreSQL"**
3. **Configure:**
   - Name: `nutrition-postgres`
   - Version: `15`
   - Database Name: `nutrition_platform`
   - Username: `nutrition_user`
   - Password: `secure_db_password_123`

### Redis Cache:
1. **Add another service**
2. **Select "Redis"**
3. **Configure:**
   - Name: `nutrition-redis`
   - Version: `7-alpine`
   - Password: `secure_redis_password_123`

## Step 7: Deploy Application

1. **Click "Deploy"** in Coolify
2. **Wait for deployment** to complete (5-10 minutes)
3. **Check logs** for any errors

## Step 8: Run Database Setup

**After deployment completes:**

1. **Go to "Terminal" tab** in Coolify
2. **Run these commands:**

```bash
# Navigate to backend directory
cd backend

# Run database migrations
go run cmd/migrate/main.go -direction up

# Seed initial data
go run cmd/seed/main.go
```

## Step 9: Test Your Application

**Test these URLs:**
- https://super.doctorhealthy1.com/health
- https://super.doctorhealthy1.com/api/info
- https://super.doctorhealthy1.com/api/nutrition/analyze

## Step 10: Configure SSL (Auto-handled by Coolify)

Coolify should automatically configure SSL for your domain.

---

## 🆘 Troubleshooting

### If Application Won't Start:
1. **Check Coolify logs** for error messages
2. **Verify all environment variables** are set correctly
3. **Ensure database services** are running
4. **Check domain DNS** settings

### If Database Connection Fails:
1. **Verify PostgreSQL service** is running
2. **Check database credentials** match environment variables
3. **Look at database logs** in Coolify

### If SSL Issues:
1. **Verify domain** points to your server
2. **Check Coolify SSL** configuration
3. **Wait for certificate** generation (can take 5-10 minutes)

---

## 📞 Need Help?

If you encounter any issues:

1. **Share the Coolify logs** with me
2. **Tell me which step** you're having trouble with
3. **I can help troubleshoot** specific errors

Let me know when you've completed each step and I'll help you with the next one!