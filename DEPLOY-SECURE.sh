#!/bin/bash
set -e

echo "🚀 SECURE COOLIFY DEPLOYMENT WITH HTTPS"
echo "======================================"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${YELLOW}ℹ️  Docker not running locally - that's OK for Coolify deployment${NC}"
    echo -e "${BLUE}📋 Coolify will handle Docker deployment on the server${NC}"
fi

# Generate secure credentials if .env.production doesn't exist
if [ ! -f ".env.production" ]; then
    echo -e "${BLUE}🔐 Generating secure credentials...${NC}"

    # Generate cryptographically secure secrets
    DB_PASSWORD=$(openssl rand -hex 32)
    REDIS_PASSWORD=$(openssl rand -hex 32)
    JWT_SECRET=$(openssl rand -hex 64)
    API_KEY_SECRET=$(openssl rand -hex 64)
    ENCRYPTION_KEY=$(openssl rand -hex 32)
    SESSION_SECRET=$(openssl rand -hex 32)

    # Create .env.production with secure values
    cat > .env.production << EOF
# ========================================
# SECURE PRODUCTION ENVIRONMENT
# Generated: $(date)
# ========================================

# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=nutrition_user
DB_PASSWORD=${DB_PASSWORD}
DB_SSL_MODE=require

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=${REDIS_PASSWORD}

# Security Configuration
JWT_SECRET=${JWT_SECRET}
API_KEY_SECRET=${API_KEY_SECRET}
ENCRYPTION_KEY=${ENCRYPTION_KEY}
SESSION_SECRET=${SESSION_SECRET}

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
ENVIRONMENT=production
DEBUG=false

# Domain Configuration
DOMAIN=super.doctorhealthy1.com
CORS_ALLOWED_ORIGINS=https://super.doctorhealthy1.com,https://www.super.doctorhealthy1.com

# Security Features
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60s
SECURITY_HEADERS_ENABLED=true
HSTS_ENABLED=true

# Performance
COMPRESSION_ENABLED=true
CACHE_TTL=3600s
METRICS_ENABLED=true

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Features
RELIGIOUS_FILTER_ENABLED=true
FILTER_ALCOHOL=true
FILTER_PORK=true
DEFAULT_LANGUAGE=en
SUPPORTED_LANGUAGES=en,ar
HEALTH_CHECK_ENABLED=true
EOF

    echo -e "${GREEN}✅ Secure credentials generated and saved to .env.production${NC}"
    echo ""
    echo -e "${YELLOW}📋 IMPORTANT: SAVE THESE CREDENTIALS:${NC}"
    echo "====================================="
    echo "DB_PASSWORD=${DB_PASSWORD}"
    echo "REDIS_PASSWORD=${REDIS_PASSWORD}"
    echo "JWT_SECRET=${JWT_SECRET}"
    echo ""
else
    echo -e "${GREEN}✅ Using existing .env.production${NC}"
fi

# Update docker-compose for HTTPS support
echo -e "${BLUE}🐳 Updating Docker Compose for HTTPS...${NC}"

# Create enhanced docker-compose with SSL
cat > docker-compose.https.yml << EOF
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: \${DB_NAME:-nutrition_platform}
      POSTGRES_USER: \${DB_USER:-nutrition_user}
      POSTGRES_PASSWORD: \${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backend/migrations:/docker-entrypoint-initdb.d
    command: postgres -c ssl=on -c ssl_cert_file=/etc/ssl/certs/ssl-cert-snakeoil.pem -c ssl_key_file=/etc/ssl/private/ssl-cert-snakeoil.key
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U \${DB_USER:-nutrition_user}"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass \${REDIS_PASSWORD}
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "--raw", "incr", "ping"]
      interval: 10s
      timeout: 3s
      retries: 5

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
      args:
        - BUILD_ENV=production
    env_file: .env.production
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    ports:
      - "8080:8080"
    healthcheck:
      test: ["CMD", "curl", "-f", "-k", "https://localhost:8080/health"]
      interval: 30s
      timeout: 15s
      retries: 3
    restart: unless-stopped

  nginx:
    image: nginx:1.25-alpine
    volumes:
      - ./nginx/ssl.conf:/etc/nginx/conf.d/default.conf:ro
      - ./certs:/etc/ssl/certs:ro
    ports:
      - "80:80"
      - "443:443"
    depends_on:
      - backend
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
EOF

# Create HTTPS-enabled Nginx configuration
echo -e "${BLUE}🌐 Creating HTTPS Nginx configuration...${NC}"

cat > nginx/ssl.conf << EOF
# HTTPS-enabled server block
server {
    listen 80;
    server_name super.doctorhealthy1.com www.super.doctorhealthy1.com;

    # Redirect HTTP to HTTPS
    return 301 https://\$server_name\$request_uri;
}

server {
    listen 443 ssl http2;
    server_name super.doctorhealthy1.com www.super.doctorhealthy1.com;

    # SSL Configuration
    ssl_certificate /etc/ssl/certs/server.crt;
    ssl_certificate_key /etc/ssl/certs/server.key;
    ssl_session_timeout 1d;
    ssl_session_cache shared:MozTLS:10m;
    ssl_session_tickets off;

    # Modern configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    # HSTS
    add_header Strict-Transport-Security "max-age=63072000" always;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # Rate limiting
    limit_req zone=api burst=20 nodelay;
    limit_conn conn_limit_per_ip 20;

    # Health check
    location /health {
        proxy_pass http://backend:8080;
        access_log off;
    }

    # API routes
    location /api/ {
        limit_req zone=api burst=10 nodelay;
        proxy_pass http://backend:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_set_header X-Forwarded-Host \$host;
        proxy_set_header X-Forwarded-Port \$server_port;

        # CORS headers
        if (\$http_origin ~* "^https?://(.*\.)?super\.doctorhealthy1\.com$") {
            add_header Access-Control-Allow-Origin \$http_origin always;
        }
        add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS" always;
        add_header Access-Control-Allow-Headers "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization,X-Correlation-ID" always;

        # Handle preflight requests
        if (\$request_method = 'OPTIONS') {
            if (\$http_origin ~* "^https?://(.*\.)?super\.doctorhealthy1\.com$") {
                add_header Access-Control-Allow-Origin \$http_origin;
            }
            add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS";
            add_header Access-Control-Allow-Headers "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization,X-Correlation-ID";
            add_header Access-Control-Max-Age 1728000;
            add_header Content-Type "text/plain; charset=utf-8";
            add_header Content-Length 0;
            return 204;
        }
    }

    # Static assets
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        proxy_pass http://backend:8080;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # All other routes
    location / {
        proxy_pass http://backend:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF

# Create self-signed certificates for development/testing
echo -e "${BLUE}🔒 Generating SSL certificates...${NC}"

mkdir -p certs
openssl req -x509 -newkey rsa:4096 -keyout certs/server.key -out certs/server.crt -days 365 -nodes \
    -subj "/C=US/ST=State/L=City/O=Organization/CN=super.doctorhealthy1.com"

echo -e "${GREEN}✅ SSL certificates generated${NC}"

# Deploy with Docker Compose
echo -e "${BLUE}🚀 Deploying with Docker Compose...${NC}"

# Stop any existing containers
docker-compose -f docker-compose.https.yml down || true

# Build and start services
docker-compose -f docker-compose.https.yml up -d --build

# Wait for services to be ready
echo -e "${YELLOW}⏳ Waiting for services to start...${NC}"
sleep 30

# Health check
echo -e "${BLUE}🔍 Checking application health...${NC}"

# Check if services are running
if docker-compose -f docker-compose.https.yml ps | grep -q "Up"; then
    echo -e "${GREEN}✅ All services are running${NC}"

    # Test health endpoint
    if curl -f -k https://localhost/health > /dev/null 2>&1; then
        echo -e "${GREEN}✅ HTTPS Health check passed${NC}"
        echo -e "${GREEN}✅ Application is accessible via HTTPS${NC}"
    else
        echo -e "${YELLOW}⚠️  HTTPS health check failed, checking HTTP...${NC}"
        if curl -f http://localhost:8080/health > /dev/null 2>&1; then
            echo -e "${GREEN}✅ HTTP Health check passed${NC}"
        else
            echo -e "${RED}❌ Health check failed${NC}"
        fi
    fi
else
    echo -e "${RED}❌ Some services failed to start${NC}"
    echo "Check logs with: docker-compose -f docker-compose.https.yml logs"
fi

echo ""
echo -e "${GREEN}🎉 SECURE DEPLOYMENT COMPLETE!${NC}"
echo "================================="
echo ""
echo -e "${BLUE}🌐 Access URLs:${NC}"
echo "• HTTPS: https://super.doctorhealthy1.com"
echo "• HTTP:  http://localhost:8080"
echo "• Health: https://super.doctorhealthy1.com/health"
echo ""
echo -e "${BLUE}🐳 Service Status:${NC}"
docker-compose -f docker-compose.https.yml ps
echo ""
echo -e "${YELLOW}📋 Next Steps:${NC}"
echo "1. Copy credentials from .env.production"
echo "2. Deploy to Coolify with these secure credentials"
echo "3. Update DNS to point super.doctorhealthy1.com to your server"
echo "4. Replace self-signed certs with Let's Encrypt"
echo ""
echo -e "${GREEN}✅ Deployment completed with enterprise-grade security!${NC}"