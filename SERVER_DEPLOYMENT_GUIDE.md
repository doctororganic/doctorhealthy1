# 🚀 Complete Server Deployment Guide

## Quick Start - Deploy Full Stack Application

This guide will help you deploy the complete nutrition platform to your server.

---

## Prerequisites

- Server with Ubuntu 20.04+ or similar Linux distribution
- Docker and Docker Compose installed
- Git installed
- At least 4GB RAM and 20GB disk space
- Domain name (optional, for SSL)

---

## Step 1: Clone Repository

```bash
# Clone the repository
git clone https://github.com/doctororganic/doctorhealthy1.git
cd doctorhealthy1/nutrition-platform

# Or if you have the code locally, upload it to your server
```

---

## Step 2: Install Docker & Docker Compose

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verify installation
docker --version
docker-compose --version
```

---

## Step 3: Configure Environment Variables

```bash
# Create production environment file
cat > .env.production << 'EOF'
# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=postgres
DB_PASSWORD=CHANGE_THIS_SECURE_PASSWORD

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=

# Server Configuration
PORT=8080
ENV=production
DOMAIN=yourdomain.com

# Security (Generate secure random strings)
JWT_SECRET=GENERATE_RANDOM_64_CHAR_STRING
API_KEY_SECRET=GENERATE_RANDOM_64_CHAR_STRING
SESSION_SECRET=GENERATE_RANDOM_32_CHAR_STRING

# CORS
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
EOF

# Generate secure passwords (run these commands)
echo "DB_PASSWORD=$(openssl rand -hex 32)" >> .env.production
echo "JWT_SECRET=$(openssl rand -hex 32)" >> .env.production
echo "API_KEY_SECRET=$(openssl rand -hex 32)" >> .env.production
echo "SESSION_SECRET=$(openssl rand -hex 32)" >> .env.production
```

---

## Step 4: Update Docker Compose Configuration

Edit `docker-compose.production.yml` and update:
- Database passwords
- Domain names
- Port mappings (if needed)

---

## Step 5: Build and Start Services

```bash
# Build all services
docker-compose -f docker-compose.production.yml build

# Start all services
docker-compose -f docker-compose.production.yml up -d

# Check status
docker-compose -f docker-compose.production.yml ps

# View logs
docker-compose -f docker-compose.production.yml logs -f
```

---

## Step 6: Verify Deployment

```bash
# Check backend health
curl http://localhost:8080/health

# Check frontend
curl http://localhost:3000

# Check all containers
docker ps
```

---

## Step 7: Configure Nginx (Optional - for production)

If you want to use Nginx as reverse proxy:

```bash
# Install Nginx
sudo apt install nginx

# Create Nginx configuration
sudo nano /etc/nginx/sites-available/nutrition-platform
```

Add this configuration:

```nginx
server {
    listen 80;
    server_name yourdomain.com www.yourdomain.com;

    # Frontend
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # Backend API
    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable site:
```bash
sudo ln -s /etc/nginx/sites-available/nutrition-platform /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

---

## Step 8: SSL Certificate (Let's Encrypt)

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx

# Get SSL certificate
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# Auto-renewal is set up automatically
```

---

## Management Commands

### Start Services
```bash
docker-compose -f docker-compose.production.yml start
```

### Stop Services
```bash
docker-compose -f docker-compose.production.yml stop
```

### Restart Services
```bash
docker-compose -f docker-compose.production.yml restart
```

### View Logs
```bash
# All services
docker-compose -f docker-compose.production.yml logs -f

# Specific service
docker-compose -f docker-compose.production.yml logs -f backend
docker-compose -f docker-compose.production.yml logs -f frontend
```

### Update Application
```bash
# Pull latest code
git pull

# Rebuild and restart
docker-compose -f docker-compose.production.yml up -d --build
```

### Backup Database
```bash
# Create backup
docker-compose -f docker-compose.production.yml exec postgres pg_dump -U postgres nutrition_platform > backup_$(date +%Y%m%d).sql

# Restore backup
docker-compose -f docker-compose.production.yml exec -T postgres psql -U postgres nutrition_platform < backup_20241130.sql
```

---

## Troubleshooting

### Check Container Logs
```bash
docker-compose -f docker-compose.production.yml logs backend
docker-compose -f docker-compose.production.yml logs frontend
```

### Restart Specific Service
```bash
docker-compose -f docker-compose.production.yml restart backend
```

### Check Container Status
```bash
docker-compose -f docker-compose.production.yml ps
```

### Access Container Shell
```bash
docker-compose -f docker-compose.production.yml exec backend sh
docker-compose -f docker-compose.production.yml exec postgres psql -U postgres
```

---

## Production Checklist

- [ ] Environment variables configured
- [ ] Secure passwords generated
- [ ] Docker Compose file updated
- [ ] All services running
- [ ] Health checks passing
- [ ] Nginx configured (if using)
- [ ] SSL certificate installed
- [ ] Firewall configured (ports 80, 443)
- [ ] Database backups scheduled
- [ ] Monitoring set up

---

## Support

For issues or questions:
- Check logs: `docker-compose logs`
- GitHub Issues: https://github.com/doctororganic/doctorhealthy1/issues
- Documentation: See `DEPLOYMENT_COMPLETE.md`

---

## Quick Deploy Script

Save this as `deploy.sh`:

```bash
#!/bin/bash
set -e

echo "🚀 Starting deployment..."

# Pull latest code
git pull

# Build and start
docker-compose -f docker-compose.production.yml up -d --build

# Wait for services
sleep 10

# Check health
curl -f http://localhost:8080/health || echo "Backend health check failed"
curl -f http://localhost:3000 || echo "Frontend health check failed"

echo "✅ Deployment complete!"
```

Make it executable:
```bash
chmod +x deploy.sh
./deploy.sh
```

