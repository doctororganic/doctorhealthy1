# Vultr VPS Deployment Script for Nutrition Platform
# Subdomain: love.doctorhealthy1.com
# IPv6: 2001:19f0:1000:6d30:5400:05ff:fe9f:f721

# Set error action preference to stop on errors
$ErrorActionPreference = "Stop"

# Colors for output
$Green = [ConsoleColor]::Green
$Yellow = [ConsoleColor]::Yellow
$Red = [ConsoleColor]::Red
$NC = [ConsoleColor]::White

# Configuration
$SUBDOMAIN = "love.doctorhealthy1.com"
$IPV6_ADDRESS = "2001:19f0:1000:6d30:5400:05ff:fe9f:f721"
$VPS_USER = "root"
$APP_DIR = "/opt/nutrition-platform"
$SERVICE_NAME = "nutrition-platform"

Write-Host "# Starting Vultr VPS deployment for ${SUBDOMAIN}..." -ForegroundColor $Yellow

# Check if SSH key exists
if (-not (Test-Path "~/.ssh/id_rsa")) {
    Write-Host "# Generating SSH key..." -ForegroundColor $Yellow
    & ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa -N ""
}

# Build Docker image locally
Write-Host "# Building Docker image..." -ForegroundColor $Yellow
& docker build -f Dockerfile.simple -t nutrition-platform:latest .

# Save Docker image to tar file
Write-Host "# Saving Docker image..." -ForegroundColor $Yellow
& docker save nutrition-platform:latest > nutrition-platform.tar

# Create deployment script for VPS
$deployScript = @'
#!/bin/bash
set -e

# Update system
apt update && apt upgrade -y

# Install Docker
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    systemctl enable docker
    systemctl start docker
fi

# Install Docker Compose
if ! command -v docker-compose &> /dev/null; then
    curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
fi

# Create app directory
mkdir -p /opt/nutrition-platform
cd /opt/nutrition-platform

# Load Docker image
if [ -f nutrition-platform.tar ]; then
    docker load < nutrition-platform.tar
    rm nutrition-platform.tar
fi

# Create docker-compose.yml
cat > docker-compose.yml << 'COMPOSE_EOF'
version: '3.8'
services:
  app:
    image: nutrition-platform:latest
    container_name: nutrition-platform
    restart: unless-stopped
    ports:
      - "80:8080"
      - "443:8080"
    environment:
      - ENVIRONMENT=production
      - SERVER_HOST=0.0.0.0
      - SERVER_PORT=8080
      - LOG_LEVEL=info
      - DOMAIN=love.doctorhealthy1.com
    volumes:
      - ./data:/app/data
      - ./logs:/app/logs
      - ./uploads:/app/uploads
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
COMPOSE_EOF

# Create systemd service
cat > /etc/systemd/system/nutrition-platform.service << 'SERVICE_EOF'
[Unit]
Description=Nutrition Platform
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/nutrition-platform
ExecStart=/usr/local/bin/docker-compose up -d
ExecStop=/usr/local/bin/docker-compose down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
SERVICE_EOF

# Enable and start service
systemctl daemon-reload
systemctl enable nutrition-platform
systemctl start nutrition-platform

# Install and configure Nginx
apt install -y nginx certbot python3-certbot-nginx

# Create Nginx configuration
cat > /etc/nginx/sites-available/love.doctorhealthy1.com << 'NGINX_EOF'
server {
    listen 80;
    listen [::]:80;
    server_name love.doctorhealthy1.com;
    
    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name love.doctorhealthy1.com;
    
    # SSL configuration will be added by certbot
    
    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    
    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/javascript application/xml+rss application/json;
    
    # Proxy to application
    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        proxy_read_timeout 300s;
        proxy_connect_timeout 75s;
    }
    
    # Health check endpoint
    location /health {
        proxy_pass http://localhost:8080/health;
        access_log off;
    }
}
NGINX_EOF

# Enable site
ln -sf /etc/nginx/sites-available/love.doctorhealthy1.com /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default

# Test Nginx configuration
nginx -t
systemctl reload nginx

# Get SSL certificate
certbot --nginx -d love.doctorhealthy1.com --non-interactive --agree-tos --email admin@doctorhealthy1.com

# Setup automatic certificate renewal
echo "0 12 * * * /usr/bin/certbot renew --quiet" | crontab -

echo "Deployment completed successfully!"
echo "Application is available at: https://love.doctorhealthy1.com"
'@
Set-Content -Path "deploy-on-vps.sh" -Value $deployScript

& chmod +x deploy-on-vps.sh

Write-Host "# Deployment files created successfully!" -ForegroundColor $Green
Write-Host "# Next steps:" -ForegroundColor $Yellow
Write-Host "# 1. Upload files to your Vultr VPS:" -ForegroundColor $Yellow
Write-Host "#    scp nutrition-platform.tar deploy-on-vps.sh root@[VPS_IP]:/root/" -ForegroundColor $Yellow
Write-Host "# 2. SSH into your VPS and run:" -ForegroundColor $Yellow
Write-Host "#    ssh root@[VPS_IP]" -ForegroundColor $Yellow
Write-Host "#    chmod +x deploy-on-vps.sh" -ForegroundColor $Yellow
Write-Host "#    ./deploy-on-vps.sh" -ForegroundColor $Yellow
Write-Host "# 3. Configure DNS:" -ForegroundColor $Yellow
Write-Host "#    Add A record: love.doctorhealthy1.com → [VPS_IPv4]" -ForegroundColor $Yellow
Write-Host "#    Add AAAA record: love.doctorhealthy1.com → ${IPV6_ADDRESS}" -ForegroundColor $Yellow
Write-Host "# Your application will be available at: https://love.doctorhealthy1.com" -ForegroundColor $Green