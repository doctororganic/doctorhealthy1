#!/bin/bash
# Vultr VPS Deployment Script - Server-side

set -e

# Colors for output
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
RED="\033[0;31m"
NC="\033[0m"

DOMAIN="super.doctorhealthy1.com"
APP_DIR="/opt/nutrition-platform"

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
    echo -e "${YELLOW}Please logout and login again for Docker permissions to take effect${NC}"
fi

# Install Docker Compose if not present
if ! command -v docker-compose &> /dev/null; then
    echo -e "${YELLOW}Installing Docker Compose...${NC}"
    sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
fi

# Create app directory
echo -e "${YELLOW}Creating application directory...${NC}"
sudo mkdir -p $APP_DIR
cd $APP_DIR

# Create production environment file
echo -e "${YELLOW}Creating production environment...${NC}"
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

# Database Configuration (SQLite)
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
echo -e "${YELLOW}Creating Docker Compose configuration...${NC}"
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
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
COMPOSE_EOF

# Load Docker image
echo -e "${YELLOW}Loading Docker image...${NC}"
docker load < nutrition-platform.tar

# Start application
echo -e "${YELLOW}Starting application...${NC}"
docker-compose up -d

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
sudo ln -sf /etc/nginx/sites-available/super.doctorhealthy1.com /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default

# Test Nginx configuration
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx

# Check if application is running
echo -e "${YELLOW}Checking application status...${NC}"
sleep 10
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Nutrition Platform is running successfully!${NC}"
    echo -e "${GREEN}🌐 Application available at: http://super.doctorhealthy1.com${NC}"
else
    echo -e "${RED}❌ Application health check failed${NC}"
    docker-compose logs
fi

echo -e "${YELLOW}To set up SSL (HTTPS):${NC}"
echo -e "1. Install Certbot: sudo apt install certbot python3-certbot-nginx"
echo -e "2. Run: sudo certbot --nginx -d super.doctorhealthy1.com"
echo -e "3. Test auto-renewal: sudo certbot renew --dry-run"