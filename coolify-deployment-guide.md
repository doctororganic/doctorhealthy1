# 🚀 COOLIFY DEPLOYMENT GUIDE

## 📋 STEP-BY-STEP COOLIFY DEPLOYMENT

### Step 1: Access Coolify Dashboard

🌐 Go to: https://api.doctorhealthy1.com
🔑 Login with your Coolify credentials

### Step 2: Navigate to Applications

📍 Click "Applications" in left sidebar
📋 Click "Add Application" button

### Step 3: Upload Deployment Package

📦 Select "Upload ZIP file"
📁 Choose file: `nutrition-platform-coolify-20251013-164858.zip`
📝 Application Name: `nutrition-platform-secure`
📋 Description: `AI-powered nutrition platform with enterprise security`

### Step 4: Configure Source Settings

🔧 Source Type: Archive
📦 Archive Type: ZIP file
🌐 Repository URL: (leave empty)
📁 Root Directory: `/` (root of archive)

### Step 5: Configure Build Settings

🔨 Build Pack: Dockerfile
📄 Dockerfile Location: `backend/Dockerfile`
🏗️ Build Context: `./`
🚀 Install Command: (blank)
🏗️ Build Command: (blank)
🚀 Start Command: (blank)

### Step 6: Configure Deployment Settings

🌐 Domain: `super.doctorhealthy1.com`
🔌 Port: `8080`
📊 Health Check Path: `/health`
⏱️ Health Check Interval: `30s`
🔄 Auto Deploy: ✅ Enabled

### Step 7: Add Environment Variables ⚠️ CRITICAL STEP

⚙️ Click "Environment Variables" tab
📋 Click "Bulk Import" or add individually
📝 Copy ALL variables from: `nutrition-platform/.env.production`

🔐 CRITICAL VARIABLES:
```
DB_PASSWORD=ac287cc0e30f54afad53c6dc7e02fd0cccad979d62b75d75d97b1ede12daf8d5
REDIS_PASSWORD=f606b2d16d6697e666ce78a8685574d042df15484ca8f18f39f2e67bf38dc09a
JWT_SECRET=9a00511e8e23764f8f4524c02f1db9eccc1923208c02fb36cb758d874d8d569bce9ea1b24ac18a958334abe15ef89e09d6010fe64a1d1ffc02a45b07898b2473
API_KEY_SECRET=5d2763e839f7e71b90ff88bef12f690a41802635aa131f6bc7160056ef0aeb7dc9caaeb07dbe0028128e617529a48903f8d01c6cc64ce61419eb7f309fdfc8bc
ENCRYPTION_KEY=cc1574e486b2f5abd69d86537079ba928974cc463e36ff410647b15b15533d23
SESSION_SECRET=f40776484ee20b35e4f754909fb3067cef2a186d0da7c4c24f1bcd54870d9fba
```

### Step 8: Add Database Services 🗄️ IMPORTANT

🗄️ Click "Services" tab
➕ Click "Add Service"
📦 Select "PostgreSQL"
📝 Name: `nutrition-postgres`
📊 Version: `15`
🗃️ Database: `nutrition_platform`
👤 Username: `nutrition_user`
🔑 Password: `ac287cc0e30f54afad53c6dc7e02fd0cccad979d62b75d75d97b1ede12daf8d5`

➕ Click "Add Another Service"
📦 Select "Redis"
📝 Name: `nutrition-redis`
📊 Version: `7-alpine`
🔑 Password: `f606b2d16d6697e666ce78a8685574d042df15484ca8f18f39f2e67bf38dc09a`

### Step 9: Deploy Application 🚀 FINAL STEP

🚀 Click "Deploy" button (top right)
⏳ Wait 5-10 minutes for deployment
📊 Monitor: Click "Deployments" tab → Watch real-time logs
✅ Success: Green checkmark appears
❌ Error: Red X with error message

## 🔍 POST-DEPLOYMENT VERIFICATION

Test These URLs After Deployment:

🌐 Main Site: `https://super.doctorhealthy1.com`
🔍 Health Check: `https://super.doctorhealthy1.com/health`
📚 API Info: `https://super.doctorhealthy1.com/api/v1/info`
🧪 Nutrition Test: `POST https://super.doctorhealthy1.com/api/v1/nutrition/analyze`

### Test Payload for Nutrition API:
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

## 🐛 TROUBLESHOOTING

### If Deployment Fails:
❌ Check Coolify logs for specific errors
❌ Verify all environment variables are set
❌ Ensure database services are running
❌ Check SSL certificate generation

### If Application Won't Start:
❌ Check: Database connection strings
❌ Verify: Redis password matches
❌ Confirm: All environment variables present
❌ Check: Port 8080 is available

## 📋 DEPLOYMENT SUMMARY

| Component | Status | Configuration |
|----------|--------|-------------|
| 📦 Deployment Package | ✅ Ready | 98MB ZIP file |
| 🔐 Environment Variables | ✅ Configured | All secure credentials |
| 🗄️ Database Services | ✅ Ready | PostgreSQL + Redis |
| 🌐 Domain | ✅ Ready | super.doctorhealthy1.com |
| 🔒 Security | ✅ Ready | Enterprise-grade |
| 📊 Monitoring | ✅ Ready | Health checks active |

## 🎯 FINAL ACTION REQUIRED

🚀 YOU NEED TO:

1. 🌐 Access Coolify: `https://api.doctorhealthy1.com`
2. 📦 Upload ZIP: `nutrition-platform-coolify-20251013-164858.zip`
3. ⚙️ Configure settings as shown above
4. 🔐 Add environment variables from `.env.production`
5. 🗄️ Add database services (PostgreSQL + Redis)
6. 🚀 Click Deploy and wait 5-10 minutes

## 📞 SUPPORT

For any issues:
1. Check the test suites in `/tests/` directory
2. Review the deployment logs in Coolify
3. Verify environment variables are set correctly

## 🎊 CONGRATULATIONS!

Once deployed, your AI-powered nutrition platform will be live with:
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
**Deployment Status:** ✅ READY FOR DEPLOYMENT  
**Security Level:** 🔒 PRODUCTION READY