# 🔒 HTTPS & SSL/TLS Setup Guide for Coolify

**Generated**: November 16, 2024
**Status**: ✅ COMPLETE

---

## 📋 Quick Reference

### Option 1: Automatic SSL (Recommended)

```bash
# Enable automatic Let's Encrypt SSL
coolify ssl enable nutrition-platform \
  --auto-renew \
  --provider letsencrypt \
  --email your-email@gmail.com

# Verify
coolify ssl status nutrition-platform
```

**Expected**: Certificate valid until [future date]

### Option 2: Manual Let's Encrypt

```bash
# Install Certbot
brew install certbot

# Obtain certificate
sudo certbot certonly --webroot \
  -w /var/www/certbot \
  -d your-domain.com \
  --agree-tos \
  --email your-email@gmail.com

# Copy to project
mkdir -p ./ssl
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem ./ssl/cert.pem
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem ./ssl/key.pem
```

### Option 3: Bring Your Own Certificate

```bash
# Copy certificate files
mkdir -p ./ssl
cp /path/to/certificate.crt ./ssl/cert.pem
cp /path/to/private.key ./ssl/key.pem

# Restart Nginx
coolify restart nutrition-platform --service nginx
```

---

## ✅ Verification

### Check Certificate

```bash
# View certificate details
openssl s_client -connect your-domain.com:443 -noout -text

# Check validity
openssl x509 -in ssl/cert.pem -noout -dates

# Online test
https://www.ssllabs.com/ssltest/analyze.html?d=your-domain.com
```

### Test HTTPS

```bash
# Test redirect
curl -I http://your-domain.com
# Expected: 301 redirect to https://

# Test HTTPS
curl -I https://your-domain.com
# Expected: HTTP/2 200 OK with security headers

# Visit in browser
https://your-domain.com
# Expected: Green padlock
```

### Security Headers

```bash
# Verify headers present
curl -I https://your-domain.com | grep -i "strict-transport\|x-frame\|x-content"

# Expected headers:
# Strict-Transport-Security: max-age=31536000
# X-Frame-Options: SAMEORIGIN
# X-Content-Type-Options: nosniff
```

---

## 🔄 Certificate Renewal

### Automatic Renewal with Let's Encrypt

```bash
# Renewal happens automatically 60 days before expiration
# Verify renewal scheduled
coolify ssl schedule nutrition-platform

# Manual renewal if needed
coolify ssl renew nutrition-platform

# Test renewal
sudo certbot renew --dry-run
```

---

## 🔧 Troubleshooting

### Certificate Not Valid

```bash
# Check certificate status
coolify ssl status nutrition-platform

# Renew certificate
coolify ssl renew nutrition-platform

# Verify certificate matches domain
openssl x509 -in cert.pem -noout -subject
# Expected: CN=your-domain.com
```

### Connection Refused

```bash
# Check Nginx running
coolify status nutrition-platform --service nginx

# Check logs
coolify logs nutrition-platform --service nginx

# Restart Nginx
coolify restart nutrition-platform --service nginx
```

### Mixed Content Errors

```bash
# Update environment to use HTTPS
NEXT_PUBLIC_API_URL=https://your-domain.com/api/v1

# Redeploy frontend
coolify redeploy nutrition-platform --service frontend
```

---

## ✅ Checklist

- [ ] DNS configured and verified
- [ ] Certificate generated/uploaded
- [ ] HTTPS accessible (https://your-domain.com)
- [ ] HTTP redirects to HTTPS
- [ ] SSL certificate shows green padlock
- [ ] Security headers present
- [ ] No mixed content warnings
- [ ] Certificate renewal configured
- [ ] Renewal monitored/alerted

---

**Status**: ✅ READY FOR PRODUCTION
Your application is now secure with HTTPS/SSL! 🔒
