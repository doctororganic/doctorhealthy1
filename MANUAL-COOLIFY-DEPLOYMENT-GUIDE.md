# 🔗 MANUAL COOLIFY DEPLOYMENT GUIDE

## 🌐 Coolify Dashboard Link
**Access URL:** https://api.doctorhealthy1.com

## 📋 Step-by-Step Deployment Instructions

### Step 1: Access Coolify Dashboard
1. Open your web browser
2. Go to: **https://api.doctorhealthy1.com**
3. Login with your credentials

### Step 2: Navigate to Applications
1. In the left sidebar, click **"Applications"**
2. Click the **"Add Application"** button

### Step 3: Create New Application
1. **Application Name:** `nutrition-platform-secure`
2. **Description:** `AI-powered nutrition platform with enterprise security`
3. **Project:** `new doctorhealthy1`
4. **Environment:** `production`

### Step 4: Configure Source
1. **Source Type:** **GitHub Repository** (if using Git) or **Upload ZIP file**
2. **If using ZIP:**
   - Select **"Upload ZIP file"**
   - Choose: `nutrition-platform-coolify-20251013-164858.zip`
   - Root Directory: `/`

### Step 5: Configure Build Settings
1. **Build Pack:** **Dockerfile**
2. **Dockerfile Location:** `backend/Dockerfile`
3. **Build Context:** `./`
4. **Port:** `8080`

### Step 6: Configure Deployment
1. **Domain:** `super.doctorhealthy1.com`
2. **Health Check Path:** `/health`
3. **Health Check Interval:** `30s`
4. **Auto Deploy:** **Enabled**

### Step 7: Add Environment Variables (CRITICAL)
1. Click **"Environment Variables"** tab
2. Click **"Add Variable"** and add each variable:

#### Database Configuration:
```
DB_HOST=localhost
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=nutrition_user
DB_PASSWORD=ac287cc0e30f54afad53c6dc7e02fd0cccad979d62b75d75d97b1ede12daf8d5
DB_SSL_MODE=require
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5
```

#### Redis Configuration:
```
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0
REDIS_PASSWORD=f606b2d16d6697e666ce78a8685574d042df15484ca8f18f39f2e67bf38dc09a
```

#### Security Configuration:
```
JWT_SECRET=9a00511e8e23764f8f4524c02f1db9eccc1923208c02fb36cb758d874d8d569bce9ea1b24ac18a958334abe15ef89e09d6010fe64a1d1ffc02a45b07898b2473
API_KEY_SECRET=5d2763e839f7e71b90ff88bef12f690a41802635aa131f6bc7160056ef0aeb7dc9caaeb07dbe0028128e617529a48903f8d01c6cc64ce61419eb7f309fdfc8bc
ENCRYPTION_KEY=cc1574e486b2f5abd69d86537079ba928974cc463e36ff410647b15b15533d23
SESSION_SECRET=f40776484ee20b35e4f754909fb3067cef2a186d0da7c4c24f1bcd54870d9fba
```

#### Server Configuration:
```
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
CORS_ALLOWED_ORIGINS=https://super.doctorhealthy1.com,https://my.doctorhealthy1.com
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60
LOG_LEVEL=info
LOG_FORMAT=json
```

#### Features:
```
RELIGIOUS_FILTER_ENABLED=true
FILTER_ALCOOL=true
FILTER_PORK=true
FILTER_STRICT_MODE=false
DEFAULT_LANGUAGE=en
SUPPORTED_LANGUAGES=en,ar
RTL_LANGUAGES=ar
HEALTH_CHECK_ENABLED=true
HEALTH_CHECK_INTERVAL=30
HEALTH_CHECK_TIMEOUT=5
```

### Step 8: Add Database Services
1. Click **"Services"** tab
2. Click **"Add Service"**

#### Add PostgreSQL:
1. **Service Name:** `nutrition-postgres`
2. **Type:** **PostgreSQL**
3. **Version:** `15`
4. **Database:** `nutrition_platform`
5. **Username:** `nutrition_user`
6. **Password:** `ac287cc0e30f54afad53c6dc7e02fd0cccad979d62b75d75d97b1ede12daf8d5`
7. Click **"Add Service"**

#### Add Redis:
1. **Service Name:** `nutrition-redis`
2. **Type:** **Redis**
3. **Version:** `7-alpine`
4. **Password:** `f606b2d16d6697e666ce78a8685574d042df15484ca8f18f39f2e67bf38dc09a`
5. Click **"Add Service"**

### Step 9: Deploy Application
1. Click the **"Deploy"** button
2. Wait for deployment to complete (5-10 minutes)
3. Monitor progress in the **"Deployments"** tab

## 🔍 Post-Deployment Verification

### Test URLs After Deployment:
1. **Main Website:** https://super.doctorhealthy1.com
2. **Health Check:** https://super.doctorhealthy1.com/health
3. **API Info:** https://super.doctorhealthy1.com/api/v1/info

### Test API Endpoint:
1. Use Postman or curl:
2. **URL:** `POST https://super.doctorhealthy1.com/api/v1/nutrition/analyze`
3. **Method:** POST
4. **Headers:** `Content-Type: application/json`
5. **Body:**
```json
{
  "food": "chicken breast",
  "quantity": 100,
  "unit": "grams",
  "checkHalal": true,
  "language": "en"
}
```

### Expected Health Response:
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

## 📊 Monitoring in Coolify

### Application Logs:
1. Go to **Applications** → **nutrition-platform-secure**
2. Click **"Logs"** tab
3. View real-time application logs

### Resource Usage:
1. Click **"Resources"** tab
2. Monitor CPU, memory, and disk usage
3. Set up alerts for high resource usage

### Health Checks:
1. Go to **"Health Checks"** tab
2. View application health status
3. Configure notification settings

## 📋 Deployment Checklist

- [ ] Application created successfully
- [ ] Environment variables configured
- [ ] Database services added
- [ ] Deployment started
- [ ] Application is accessible at https://super.doctorhealthy1.com
- [ ] Health check endpoint is working
- [ ] API endpoints are responding
- [ ] SSL certificate is valid
- [ ] Database connections are encrypted

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
4. Review the troubleshooting guide

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

## 🔗 Important Links

- **Coolify Dashboard:** https://api.doctorhealthy1.com
- **Project:** new doctorhealthy1
- **Environment:** production
- **Server:** 128.140.111.171
- **Domain:** super.doctorhealthy1.com

---
**Last Updated:** October 13, 2025  
**Deployment Status:** ✅ READY FOR MANUAL DEPLOYMENT  
**Security Level:** 🔒 PRODUCTION READY