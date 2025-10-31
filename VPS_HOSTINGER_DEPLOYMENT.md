# VPS Hostinger Deployment Guide

## Server Information
- **Domain**: ieltspass1.com
- **Email**: Khaledalzayat278@gmail.com
- **VPS Provider**: Hostinger

## Prerequisites

### 1. Server Setup
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install required packages
sudo apt install -y nginx docker.io docker-compose git certbot python3-certbot-nginx ufw

# Start and enable services
sudo systemctl start nginx
sudo systemctl enable nginx
sudo systemctl start docker
sudo systemctl enable docker

# Add user to docker group
sudo usermod -aG docker $USER
```

### 2. Firewall Configuration
```bash
# Configure UFW firewall
sudo ufw allow OpenSSH
sudo ufw allow 'Nginx Full'
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 8080/tcp  # Backend API
sudo ufw --force enable
```

### 3. Domain Configuration
- Point your domain `ieltspass1.com` to your VPS IP address
- Add both `ieltspass1.com` and `www.ieltspass1.com` A records

## Deployment Steps

### 1. Clone Repository
```bash
cd /opt
sudo git clone https://github.com/yourusername/nutrition-platform.git
sudo chown -R $USER:$USER nutrition-platform
cd nutrition-platform
```

### 2. Environment Configuration
```bash
# Copy environment file
cp .env.docker .env

# Edit environment variables for production
nano .env
```

### 3. Build and Deploy with Docker
```bash
# Build images
docker-compose build

# Start services
docker-compose up -d

# Check status
docker-compose ps
```

### 4. Nginx Configuration
```bash
# Remove default nginx config
sudo rm /etc/nginx/sites-enabled/default

# Create new site configuration
sudo nano /etc/nginx/sites-available/ieltspass1.com
```

Add the following configuration:
```nginx
# HTTP server - redirect to HTTPS
server {
    listen 80;
    server_name ieltspass1.com www.ieltspass1.com;
    return 301 https://$server_name$request_uri;
}

# HTTPS server
server {
    listen 443 ssl http2;
    server_name ieltspass1.com www.ieltspass1.com;
    
    # SSL configuration
    ssl_certificate /etc/letsencrypt/live/ieltspass1.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/ieltspass1.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    
    # Enable gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_comp_level 6;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/javascript
        application/xml+rss
        application/json
        application/manifest+json
        image/svg+xml;
    
    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;
    add_header Content-Security-Policy "default-src 'self' https: data: blob: 'unsafe-inline' 'unsafe-eval'; connect-src 'self' https: wss:" always;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    
    # PWA headers
    add_header Service-Worker-Allowed "/" always;
    
    # Frontend (served by Docker container)
    location / {
        proxy_pass http://localhost:8081;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
    
    # Backend API
    location /api/ {
        proxy_pass http://localhost:8080/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
    
    # Health check
    location /health {
        access_log off;
        return 200 "healthy\n";
        add_header Content-Type text/plain;
    }
}
```

### 5. Enable Site and SSL
```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/ieltspass1.com /etc/nginx/sites-enabled/

# Test nginx configuration
sudo nginx -t

# Reload nginx
sudo systemctl reload nginx

# Get SSL certificate
sudo certbot --nginx -d ieltspass1.com -d www.ieltspass1.com
```

### 6. Auto-renewal Setup
```bash
# Test auto-renewal
sudo certbot renew --dry-run

# Add cron job for auto-renewal
sudo crontab -e
# Add this line:
0 12 * * * /usr/bin/certbot renew --quiet
```

## Docker Compose Configuration

Create or update `docker-compose.yml`:
```yaml
version: '3.8'

services:
  backend:
    build: ./backend
    container_name: nutrition-backend
    ports:
      - "8080:8080"
    environment:
      - ENV=production
      - PORT=8080
      - CORS_ORIGINS=https://ieltspass1.com,https://www.ieltspass1.com
    volumes:
      - ./data:/app/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  frontend:
    build: ./frontend
    container_name: nutrition-frontend
    ports:
      - "8081:80"
    depends_on:
      - backend
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:80/health"]
      interval: 30s
      timeout: 10s
      retries: 3

networks:
  default:
    name: nutrition_network
```

## Monitoring and Maintenance

### 1. Log Monitoring
```bash
# View container logs
docker-compose logs -f

# View nginx logs
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

### 2. System Monitoring
```bash
# Check system resources
htop
df -h
free -h

# Check docker status
docker ps
docker stats
```

### 3. Backup Strategy
```bash
# Create backup script
sudo nano /opt/backup.sh
```

Add backup script:
```bash
#!/bin/bash
BACKUP_DIR="/opt/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup directory
mkdir -p $BACKUP_DIR

# Backup application data
tar -czf $BACKUP_DIR/nutrition-platform-$DATE.tar.gz /opt/nutrition-platform

# Backup nginx configuration
tar -czf $BACKUP_DIR/nginx-config-$DATE.tar.gz /etc/nginx

# Keep only last 7 days of backups
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete

echo "Backup completed: $DATE"
```

```bash
# Make executable and add to cron
sudo chmod +x /opt/backup.sh
sudo crontab -e
# Add: 0 2 * * * /opt/backup.sh
```

## Troubleshooting

### Common Issues

1. **SSL Certificate Issues**
   ```bash
   sudo certbot certificates
   sudo certbot renew --force-renewal
   ```

2. **Docker Container Issues**
   ```bash
   docker-compose down
   docker-compose up -d --force-recreate
   ```

3. **Nginx Configuration Issues**
   ```bash
   sudo nginx -t
   sudo systemctl reload nginx
   ```

4. **Port Conflicts**
   ```bash
   sudo netstat -tulpn | grep :80
   sudo netstat -tulpn | grep :443
   ```

## Performance Optimization

### 1. Nginx Optimization
Add to `/etc/nginx/nginx.conf`:
```nginx
worker_processes auto;
worker_connections 1024;
keepalive_timeout 65;
client_max_body_size 50M;
```

### 2. Docker Resource Limits
Add to `docker-compose.yml`:
```yaml
services:
  backend:
    deploy:
      resources:
        limits:
          memory: 512M
        reservations:
          memory: 256M
```

## Security Checklist

- [x] Firewall configured (UFW)
- [x] SSL/TLS enabled (Let's Encrypt)
- [x] Security headers configured
- [x] Non-root user for services
- [x] Regular security updates
- [x] Strong passwords
- [x] SSH key authentication (recommended)
- [x] Fail2ban (optional but recommended)

## Support

For issues or questions:
1. Check logs: `docker-compose logs`
2. Verify nginx config: `sudo nginx -t`
3. Check SSL status: `sudo certbot certificates`
4. Monitor resources: `htop`, `df -h`

---

**Deployment Date**: $(date)
**Domain**: ieltspass1.com
**Status**: Ready for deployment