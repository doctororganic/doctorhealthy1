# 🌐 Custom Domain Configuration Guide

**Generated**: November 16, 2024
**Status**: ✅ COMPLETE

---

## 📋 Overview

```
Your Domain → DNS Registrar → A Record → Coolify Server IP → Nginx → Your App
```

---

## Step 1: Get Your Coolify Server IP

```bash
# Get IP from Coolify
coolify server ip nutrition-platform

# Note this IP (e.g., 192.168.1.100)
```

---

## Step 2: Add DNS Records (Choose Your Registrar)

### GoDaddy

1. Login to GoDaddy.com
2. Go to "My Products" → Your Domain
3. Click "Manage DNS"
4. Add A Record:
   - Name: **@**
   - Type: **A**
   - Value: **[Your Coolify IP]**
   - TTL: **3600**
5. Add CNAME for www:
   - Name: **www**
   - Type: **CNAME**
   - Value: **your-domain.com**

### Namecheap

1. Login to Namecheap.com
2. Click "Manage" next to your domain
3. Click "Advanced DNS"
4. Add A Record:
   - Host: **@**
   - Type: **A Record**
   - Value: **[Your Coolify IP]**
   - TTL: **3600**
5. Add CNAME:
   - Host: **www**
   - Type: **CNAME Record**
   - Value: **your-domain.com**

### Google Domains

1. Login to domains.google
2. Select your domain
3. Click "DNS" → "Custom Records"
4. Create A Record:
   - Name: **[blank]** or **@**
   - Type: **A**
   - Value: **[Your Coolify IP]**
5. Create CNAME Record:
   - Name: **www**
   - Type: **CNAME**
   - Value: **your-domain.com**

### Route 53 (AWS)

1. Go to Route 53 dashboard
2. Create hosted zone for your domain
3. Create A Record:
   - Name: **your-domain.com**
   - Type: **A**
   - Value: **[Your Coolify IP]**
4. Create CNAME Record:
   - Name: **www.your-domain.com**
   - Type: **CNAME**
   - Value: **your-domain.com**
5. Update your registrar to use Route 53 nameservers

---

## Step 3: Wait for DNS Propagation

```bash
# Wait 5-30 minutes for DNS to propagate
# Check status:
nslookup your-domain.com

# Expected: Should show your Coolify IP
```

---

## Step 4: Add Domain to Coolify

```bash
# Add domain to project
coolify domain add nutrition-platform \
  --domain your-domain.com \
  --primary

# Add www subdomain
coolify domain add nutrition-platform \
  --domain www.your-domain.com \
  --alias

# Verify
coolify domain list nutrition-platform
```

---

## Step 5: Update Application Configuration

```bash
# Update environment variables
# Edit .env.coolify.production:
DOMAIN=your-domain.com
NEXT_PUBLIC_API_URL=https://your-domain.com/api/v1
NEXT_PUBLIC_APP_URL=https://your-domain.com
CORS_ORIGIN=https://your-domain.com

# Apply changes
coolify env set nutrition-platform \
  --from-file .env.coolify.production

# Redeploy
coolify redeploy nutrition-platform --service frontend
coolify restart nutrition-platform --service nginx
```

---

## ✅ Verification

```bash
# Test domain resolves
nslookup your-domain.com
# Expected: Shows your Coolify IP

# Test HTTP access
curl -I http://your-domain.com
# Expected: Redirects to HTTPS

# Test HTTPS access
curl -I https://your-domain.com
# Expected: HTTP/2 200 OK

# Visit in browser
https://your-domain.com
# Expected: Application loads with green padlock
```

---

## 🔧 Troubleshooting

### Domain Doesn't Resolve

```bash
# Check DNS records
dig your-domain.com

# DNS propagation checker
# https://www.dnschecker.org/

# If using Cloudflare:
# Allow it 24-48 hours for nameserver update

# Force local DNS refresh (macOS)
sudo dscacheutil -flushcache
```

### Connection Refused

```bash
# Check Coolify running
coolify status nutrition-platform

# Verify domain added to Coolify
coolify domain list nutrition-platform

# Check firewall
ping your-domain.com
```

### Wrong IP

```bash
# Go back to registrar
# Update A record with correct Coolify IP
# Wait 5-30 minutes for propagation
# Verify: nslookup your-domain.com
```

---

## ✅ Checklist

- [ ] Domain registered and owned
- [ ] Coolify IP obtained
- [ ] A record added (@ → IP)
- [ ] CNAME record added (www → domain)
- [ ] DNS propagation verified (5-30 min)
- [ ] Domain added to Coolify
- [ ] Environment variables updated
- [ ] Application redeployed
- [ ] HTTPS working
- [ ] Green padlock showing

---

**Status**: ✅ READY FOR PRODUCTION
Your application is now accessible via your custom domain! 🌐
