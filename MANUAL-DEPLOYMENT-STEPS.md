# 🚀 MANUAL COOLIFY DEPLOYMENT STEPS

## 📋 Step-by-Step Instructions

### Step 1: Access Coolify Dashboard
1. Open your web browser
2. Go to: **https://api.doctorhealthy1.com**
3. Login with your Coolify credentials

### Step 2: Navigate to Applications
1. In the left sidebar, click **"Applications"**
2. Click the **"Add Application"** button (usually at the top right)

### Step 3: Upload Deployment Package
1. Select **"Upload ZIP file"** as the source type
2. Click **"Choose file"** and select: `nutrition-platform-coolify-20251013-164858.zip`
3. Application Name: `nutrition-platform-secure`
4. Description: `AI-powered nutrition platform with enterprise security`
5. Click **"Next"** or **"Continue"**

### Step 4: Configure Source Settings
1. Source Type: **Archive**
2. Archive Type: **ZIP file**
3. Repository URL: (leave empty)
4. Root Directory: **/**** (root of archive)
5. Click **"Next"**

### Step 5: Configure Build Settings
1. Build Pack: **Dockerfile**
2. Dockerfile Location: **backend/Dockerfile**
3. Build Context: **./**
4. Install Command: (leave blank)
5. Build Command: (leave blank)
6. Start Command: (leave blank)
7. Click **"Next"**

### Step 6: Configure Deployment Settings
1. Domain: **super.doctorhealthy1.com**
2. Port: **8080**
3. Health Check Path: **/health**
4. Health Check Interval: **30s**
5. Health Check Timeout: **5s**
6. Health Check Retries: **3**
7. Auto Deploy: **Enabled** (toggle on)
8. Click **"Next"**

### Step 7: Add Environment Variables (CRITICAL)
1. Click the **"Environment Variables"** tab
2. Click **"Bulk Import"** or **"Add Variable"**
3. Copy and paste ALL variables from below:

```
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=nutrition_user
DB_PASSWORD=ac287cc0e30f54afad53c6dc7e02fd0cccad979d62b75d75d97b1ede12daf8d5
DB_SSL_MODE=require
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0
REDIS_PASSWORD=f606b2d16d6697e666ce78a8685574d042df15484ca8f18f39f2e67bf38dc09a

# Security Configuration
JWT_SECRET=9a00511e8e23764f8f4524c02f1db9eccc1923208c02fb36cb758d874d8d569bce9ea1b24ac18a958334abe15ef89e09d6010fe64a1d1ffc02a45b07898b2473
API_KEY_SECRET=5d2763e839f7e71b90ff88bef12f690a41802635aa131f6bc7160056ef0aeb7dc9caaeb07dbe0028128e617529a48903f8d01c6cc64ce61419eb7f309fdfc8bc
ENCRYPTION_KEY=cc1574e486b2f5abd69d86537079ba928974cc463e36ff410647b15b15533d23
SESSION_SECRET=f40776484ee20b35e4f754909fb3067cef2a186d0da7c4c24f1bcd54870d9fba

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://super.doctorhealthy1.com,https://my.doctorhealthy1.com

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Features
RELIGIOUS_FILTER_ENABLED=true
FILTER_ALCOOL=true
FILTER_PORK=true
FILTER_STRICT_MODE=false

# Internationalization
DEFAULT_LANGUAGE=en
SUPPORTED_LANGUAGES=en,ar
RTL_LANGUAGES=ar

# Health Check
HEALTH_CHECK_ENABLED=true
HEALTH_CHECK_INTERVAL=30
HEALTH_CHECK_TIMEOUT=5
```

4. Click **"Save"** or **"Next"**

### Step 8: Add Database Services
1. Click the **"Services"** tab
2. Click **"Add Service"**
3. Select **"PostgreSQL"**
4. Service Name: `nutrition-postgres`
5. Version: **15**
6. Database: `nutrition_platform`
7. Username: `nutrition_user`
8. Password: `ac287cc0e30f54afad53c6dc7e02fd0cccad979d62b75d75d97b1ede12daf8d5`
9. Click **"Add Service"**

10. Click **"Add Another Service"**
11. Select **"Redis"**
12. Service Name: `nutrition-redis`
13. Version: **7-alpine**
14. Password: `f606b2d16d6697e666ce78a8685574d042df15484ca8f18f39f2e67bf38dc09a`
15. Click **"Add Service"**

### Step 9: Deploy Application
1. Click the **"Deploy"** button (usually at the top right)
2. Wait for the deployment to start (you'll see progress indicators)
3. Monitor the deployment in the **"Deployments"** tab
4. Wait 5-10 minutes for the deployment to complete

## 🔍 Post-Deployment Verification

Once deployment is complete, test these URLs:

1. **Main Website:** https://super.doctorhealthy1.com
2. **Health Check:** https://super.doctorhealthy1.com/health
3. **API Info:** https://super.doctorhealthy1.com/api/v1/info

### Test API Endpoint
1. Use Postman or curl to test the nutrition API:
2. URL: `POST https://super.doctorhealthy1.com/api/v1/nutrition/analyze`
3. Method: POST
4. Headers: Content-Type: application/json
5. Body:
```json
{
  "food": "chicken breast",
  "quantity": 100,
  "unit": "grams",
  "checkHalal": true,
  "language": "en"
}
```

### Expected Health Response
```json
{
  "status": "healthy",
  "timestamp": "2025-10-13T...",
  "version": "1.0.0"
}
```

## 🚨 Troubleshooting

### If Deployment Fails:
1. Check the **"Deployments"** tab for error messages
2. Verify all environment variables are set correctly
3. Ensure database services are running
4. Check if the ZIP file was uploaded correctly

### If Application Won't Start:
1. Verify database connection strings
2. Check Redis password matches environment variable
3. Confirm all required environment variables are present
4. Check if port 8080 is available

### If URLs Not Accessible:
1. Wait 5-15 minutes for SSL certificate provisioning
2. Check DNS configuration for the domain
3. Verify firewall rules allow traffic on ports 80 and 443

## 📊 Monitoring

### In Coolify Dashboard:
1. Click **"Applications"** → **nutrition-platform-secure**
2. View **"Logs"** for real-time application logs
3. Check **"Metrics"** for performance data
4. Monitor **"Resources"** for resource usage

### External Monitoring:
1. **Health Check:** https://super.doctorhealthy1.com/health
2. **SSL Certificate:** Check certificate validity in browser
3. **Performance:** Test application response times

## 🎯 Success Indicators

Your deployment is successful when:
- ✅ Application loads at https://super.doctorhealthy1.com
- ✅ Health check returns 200 OK
- ✅ API endpoints respond correctly
- ✅ SSL certificate is valid and trusted
- ✅ All security headers are present
- ✅ Database connections are encrypted
- ✅ CORS is properly configured

## 📞 Support

If you encounter issues:
1. Check the deployment logs in Coolify
2. Verify environment variables match exactly
3. Ensure database services are running
4. Check the troubleshooting guide in the documentation

## 🎊 Congratulations!

Once deployed successfully, your AI-powered nutrition platform will be live with:
- ✅ Real-time nutrition analysis
- ✅ 10 evidence-based diet plans
- ✅ Recipe management system
- ✅ Health tracking and analytics
- ✅ Medication management
- ✅ Workout programs
- ✅ Multi-language support (EN/AR)
- ✅ Religious dietary filtering
- ✅ SSL secured with HTTPS

---
**Last Updated:** October 13, 2025  
**Deployment Status:** ✅ READY FOR MANUAL DEPLOYMENT  
**Security Level:** 🔒 PRODUCTION READY