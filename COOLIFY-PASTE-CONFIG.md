# 🚀 COOLIFY - READY TO PASTE CONFIGURATION

## Copy and paste these into Coolify Dashboard

---

## 1. DOCKERFILE (Copy entire content)

```dockerfile
# Production-Ready Nutrition Platform - Optimized for Coolify
FROM node:18-alpine AS base
WORKDIR /app
RUN apk add --no-cache dumb-init curl

FROM base AS dependencies
COPY production-nodejs/package*.json ./
RUN npm ci --only=production && npm cache clean --force

FROM base AS production
ENV NODE_ENV=production PORT=3000 HOST=0.0.0.0

RUN addgroup -g 1001 -S nodejs && adduser -S nodejs -u 1001

COPY --from=dependencies --chown=nodejs:nodejs /app/node_modules ./node_modules
COPY --chown=nodejs:nodejs production-nodejs/ ./

RUN mkdir -p logs data && chown -R nodejs:nodejs logs data

USER nodejs
EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD node -e "require('http').get('http://localhost:3000/health', (r) => process.exit(r.statusCode === 200 ? 0 : 1))"

ENTRYPOINT ["dumb-init", "--"]
CMD ["node", "server.js"]
```

---

## 2. ENVIRONMENT VARIABLES (Copy and paste)

```
NODE_ENV=production
PORT=3000
HOST=0.0.0.0
ALLOWED_ORIGINS=https://super.doctorhealthy1.com,https://www.super.doctorhealthy1.com
```

---

## 3. APPLICATION SETTINGS

**Basic Configuration:**
```
Name: trae-new-healthy1
Type: Dockerfile
Port: 3000
Domain: super.doctorhealthy1.com
```

**Build Configuration:**
```
Dockerfile Path: ./Dockerfile
Context: .
Build Args: (none)
```

**Health Check:**
```
Path: /health
Port: 3000
Interval: 30s
Timeout: 10s
Retries: 3
```

**SSL/HTTPS:**
```
Enable SSL: Yes
Force HTTPS: Yes
Auto-renew: Yes
```

---

## 4. DEPLOYMENT STEPS

### Step 1: Login to Coolify
- URL: https://api.doctorhealthy1.com
- Use your credentials

### Step 2: Navigate to Project
- Go to "Projects"
- Select "new doctorhealthy1"

### Step 3: Create Application
- Click "New Application"
- Choose "Dockerfile" type
- Name: trae-new-healthy1

### Step 4: Configure Source
**Option A: Git Repository**
- Connect your Git repository
- Branch: main or master
- Dockerfile path: ./Dockerfile

**Option B: Simple Dockerfile**
- Paste the Dockerfile from section 1 above
- Coolify will build from this

### Step 5: Set Environment Variables
- Click "Environment Variables"
- Add each variable from section 2 above
- Save

### Step 6: Configure Domain
- Click "Domains"
- Add domain: super.doctorhealthy1.com
- Enable SSL
- Save

### Step 7: Deploy
- Click "Deploy" button
- Monitor build logs
- Wait 5-10 minutes

### Step 8: Verify
```bash
# Test health endpoint
curl https://super.doctorhealthy1.com/health

# Expected response:
# {"status":"healthy","message":"Trae New Healthy1 is running successfully"}
```

---

## 5. QUICK VERIFICATION CHECKLIST

After deployment, verify:

- [ ] Build completed successfully
- [ ] Container is running
- [ ] Health check passes
- [ ] Domain resolves correctly
- [ ] SSL certificate is valid
- [ ] Homepage loads: https://super.doctorhealthy1.com
- [ ] API responds: https://super.doctorhealthy1.com/api/info
- [ ] No errors in logs

---

## 6. TROUBLESHOOTING

### Build Fails
- Check Dockerfile syntax
- Verify all files are in repository
- Check build logs for specific error

### Container Won't Start
- Check environment variables
- Review application logs
- Verify port 3000 is correct

### Domain Not Working
- Wait 5-10 minutes for SSL
- Check DNS settings
- Verify domain configuration in Coolify

### Health Check Fails
- Check if app is listening on port 3000
- Verify /health endpoint exists
- Increase health check timeout

---

## 7. POST-DEPLOYMENT TESTS

```bash
# 1. Health Check
curl https://super.doctorhealthy1.com/health

# 2. API Info
curl https://super.doctorhealthy1.com/api/info

# 3. Nutrition Analysis
curl -X POST https://super.doctorhealthy1.com/api/nutrition/analyze \
  -H "Content-Type: application/json" \
  -d '{"food":"apple","quantity":100,"unit":"g","checkHalal":true}'

# 4. Homepage (open in browser)
open https://super.doctorhealthy1.com
```

---

## 8. EXPECTED RESULTS

✅ **Build Time:** 3-5 minutes  
✅ **Deployment Time:** 1-2 minutes  
✅ **SSL Generation:** 2-5 minutes  
✅ **Total Time:** 5-10 minutes  

✅ **Success Indicators:**
- Green status in Coolify
- Health check passing
- Domain accessible via HTTPS
- No errors in logs

---

## 🎉 READY TO DEPLOY

Everything is prepared. Just follow the steps above and your app will be live in 10 minutes!

**Good luck!** 🚀
