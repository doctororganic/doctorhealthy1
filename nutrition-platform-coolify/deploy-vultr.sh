#!/bin/bash
# Vultr VPS Deployment Script for Nutrition Platform
# Server Details:
# IP: 64.176.212.30
# Domain: super.doctorhealthy1.com
# Username: linuxuser

set -e

# Colors for output
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
RED="\033[0;31m"
NC="\033[0m" # No Color

# Configuration
DOMAIN="super.doctorhealthy1.com"
VPS_USER="linuxuser"
VPS_IP="64.176.212.30"
VPS_PASSWORD="4{QpKn}9S+52hfQE"
APP_DIR="/opt/nutrition-platform"
SERVICE_NAME="nutrition-platform"

# Build Docker image locally
echo -e "${YELLOW}Building Docker image for Vultr deployment...${NC}"
cd "/Users/khaledahmedmohamed/Desktop/trae new healthy1/nutrition-platform"
docker build -f Dockerfile.simple -t nutrition-platform:latest .

# Create deployment script for VPS
cat > deploy-to-vps.sh << 'EOF'
#!/bin/bash
set -e

# Colors for output
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
RED="\033[0;31m"
NC="\033[0m"

echo -e "${YELLOW}Starting Vultr VPS deployment...${NC}"

# Update system
echo -e "${YELLOW}Updating system packages...${NC}"
sudo apt update && sudo apt upgrade -y

# Install Docker if not present
if ! command -v docker &> /dev/null; then
    echo -e "${YELLOW}Installing Docker...${NC}"
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo systemctl enable docker
    sudo systemctl start docker
    sudo usermod -aG docker $USER
fi

# Install Docker Compose if not present
if ! command -v docker-compose &> /dev/null; then
    echo -e "${YELLOW}Installing Docker Compose...${NC}"
    sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
fi

# Create app directory
echo -e "${YELLOW}Creating application directory...${NC}"
sudo mkdir -p /opt/nutrition-platform
cd /opt/nutrition-platform

# Create production environment file
cat > .env << 'ENVEOF'
ENVIRONMENT=production
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
LOG_LEVEL=info
LOG_FORMAT=json
METRICS_ENABLED=true
HEALTH_CHECK_ENABLED=true
RELIGIOUS_FILTER_ENABLED=true
FILTER_ALCOHOL=true
FILTER_PORK=true
DEFAULT_LANGUAGE=en
SUPPORTED_LANGUAGES=en,ar
RTL_LANGUAGES=ar

# Database Configuration (SQLite for simplicity)
DB_HOST=localhost
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=nutrition_user
DB_PASSWORD=secure_password
DB_SSL_MODE=disable

# Security Configuration
JWT_SECRET=your_jwt_secret_key_here_change_in_production
API_KEY_SECRET=your_api_key_secret_here_change_in_production
ENCRYPTION_KEY=your_encryption_key_here_change_in_production
CORS_ALLOWED_ORIGINS=https://super.doctorhealthy1.com
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60s
SECURITY_HEADERS_ENABLED=true

# Performance Configuration
COMPRESSION_ENABLED=true
CACHE_TTL=3600s
DEFAULT_PAGE_SIZE=20
MAX_PAGE_SIZE=100
ENVEOF

# Create docker-compose.yml
cat > docker-compose.yml << 'COMPOSE_EOF'
version: '3.8'
services:
  app:
    image: nutrition-platform:latest
    container_name: nutrition-platform
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - ENVIRONMENT=production
      - SERVER_HOST=0.0.0.0
      - SERVER_PORT=8080
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
echo -e "${YELLOW}Creating systemd service...${NC}"
sudo tee /etc/systemd/system/nutrition-platform.service > /dev/null << 'SERVICE_EOF'
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

# Install and configure Nginx
echo -e "${YELLOW}Installing Nginx...${NC}"
sudo apt install -y nginx curl

# Create Nginx configuration
echo -e "${YELLOW}Configuring Nginx...${NC}"
sudo tee /etc/nginx/sites-available/super.doctorhealthy1.com > /dev/null << 'NGINX_EOF'
server {
    listen 80;
    listen [::]:80;
    server_name super.doctorhealthy1.com;
    
    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name super.doctorhealthy1.com;
    
    # SSL configuration will be added by certbot
    
    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    
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
sudo ln -sf /etc/nginx/sites-available/super.doctorhealthy1.com /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default

# Test Nginx configuration
sudo nginx -t
sudo systemctl reload nginx

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable nutrition-platform

# Create startup script
cat > start-production.sh << 'START_EOF'
#!/bin/bash
set -e

echo "Starting Nutrition Platform in production..."

# Load Docker image if it exists
if [ -f nutrition-platform.tar ]; then
    echo "Loading Docker image..."
    docker load < nutrition-platform.tar
    rm nutrition-platform.tar
fi

# Start services
docker-compose up -d

# Check if application is running
sleep 10
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Nutrition Platform is running successfully!${NC}"
    echo -e "${GREEN}🌐 Application available at: https://super.doctorhealthy1.com${NC}"
else
    echo -e "${RED}❌ Application health check failed${NC}"
    docker-compose logs
fi
START_EOF

chmod +x start-production.sh

echo -e "${GREEN}✅ VPS setup completed!${NC}"
echo -e "${YELLOW}Next steps:${NC}"
echo -e "1. Upload the Docker image: scp nutrition-platform.tar linuxuser@64.176.212.30:/opt/nutrition-platform/"
echo -e "2. Upload this script: scp deploy-to-vps.sh linuxuser@64.176.212.30:/opt/nutrition-platform/"
echo -e "3. SSH to server: ssh linuxuser@64.176.212.30"
echo -e "4. Run: cd /opt/nutrition-platform && chmod +x deploy-to-vps.sh && ./deploy-to-vps.sh"
echo -e "5. Run: ./start-production.sh"
EOF

chmod +x deploy-to-vps.sh

echo -e "${GREEN}✅ Deployment scripts created successfully!${NC}"
echo -e "${YELLOW}To deploy to your Vultr server:${NC}"
echo -e "1. Copy Docker image to server:"
echo -e "   scp nutrition-platform.tar linuxuser@64.176.212.30:/opt/nutrition-platform/"
echo -e "2. Copy deployment script:"
echo -e "   scp deploy-to-vps.sh linuxuser@64.176.212.30:/opt/nutrition-platform/"
echo -e "3. SSH to server:"
echo -e "   ssh linuxuser@64.176.212.30"
echo -e "4. Run deployment:"
echo -e "   cd /opt/nutrition-platform && chmod +x deploy-to-vps.sh && ./deploy-to-vps.sh"