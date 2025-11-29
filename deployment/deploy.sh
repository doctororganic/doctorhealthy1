#!/bin/bash

# Nutrition Platform Deployment Script for Hostinger VPS
# Zero-interaction deployment with automated SSL, monitoring, and security

set -euo pipefail

# Configuration
APP_NAME="nutrition-platform"
APP_USER="nutrition"
APP_DIR="/var/www/${APP_NAME}"
BACKEND_PORT="8080"
FRONTEND_PORT="3000"
DOMAIN="nutrition-platform.com"
EMAIL="admin@nutrition-platform.com"
DB_NAME="nutrition_db"
DB_USER="nutrition_user"
REDIS_PORT="6379"
PROMETHEUS_PORT="9090"
GRAFANA_PORT="3001"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"
    exit 1
}

# Check if running as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        error "This script must be run as root"
    fi
}

# Update system packages
update_system() {
    log "Updating system packages..."
    apt-get update -y
    apt-get upgrade -y
    apt-get install -y curl wget git unzip software-properties-common apt-transport-https ca-certificates gnupg lsb-release
}

# Install Docker and Docker Compose
install_docker() {
    log "Installing Docker..."
    
    # Remove old versions
    apt-get remove -y docker docker-engine docker.io containerd runc || true
    
    # Add Docker's official GPG key
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    
    # Add Docker repository
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
    
    # Install Docker
    apt-get update -y
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
    
    # Start and enable Docker
    systemctl start docker
    systemctl enable docker
    
    # Install Docker Compose
    curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
    
    log "Docker installed successfully"
}

# Install Nginx
install_nginx() {
    log "Installing Nginx..."
    apt-get install -y nginx
    
    # Create cache directories
    mkdir -p /var/cache/nginx/api
    mkdir -p /var/cache/nginx/static
    chown -R www-data:www-data /var/cache/nginx
    
    # Enable and start Nginx
    systemctl enable nginx
    systemctl start nginx
    
    log "Nginx installed successfully"
}

# Install Certbot for SSL
install_certbot() {
    log "Installing Certbot for SSL certificates..."
    apt-get install -y certbot python3-certbot-nginx
    log "Certbot installed successfully"
}

# Install PostgreSQL
install_postgresql() {
    log "Installing PostgreSQL..."
    apt-get install -y postgresql postgresql-contrib
    
    # Start and enable PostgreSQL
    systemctl start postgresql
    systemctl enable postgresql
    
    log "PostgreSQL installed successfully"
}

# Install Redis
install_redis() {
    log "Installing Redis..."
    apt-get install -y redis-server
    
    # Configure Redis
    sed -i 's/^# maxmemory <bytes>/maxmemory 256mb/' /etc/redis/redis.conf
    sed -i 's/^# maxmemory-policy noeviction/maxmemory-policy allkeys-lru/' /etc/redis/redis.conf
    
    # Start and enable Redis
    systemctl start redis-server
    systemctl enable redis-server
    
    log "Redis installed successfully"
}

# Install Node.js
install_nodejs() {
    log "Installing Node.js..."
    curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
    apt-get install -y nodejs
    
    # Install PM2 for process management
    npm install -g pm2
    
    log "Node.js and PM2 installed successfully"
}

# Install Go
install_go() {
    log "Installing Go..."
    GO_VERSION="1.21.0"
    wget https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
    rm go${GO_VERSION}.linux-amd64.tar.gz
    
    # Add Go to PATH
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    export PATH=$PATH:/usr/local/go/bin
    
    log "Go installed successfully"
}

# Create application user
create_app_user() {
    log "Creating application user..."
    if ! id "$APP_USER" &>/dev/null; then
        useradd -r -s /bin/bash -d $APP_DIR $APP_USER
        log "User $APP_USER created"
    else
        log "User $APP_USER already exists"
    fi
}

# Setup application directory
setup_app_directory() {
    log "Setting up application directory..."
    mkdir -p $APP_DIR
    mkdir -p $APP_DIR/logs
    mkdir -p $APP_DIR/backups
    mkdir -p $APP_DIR/uploads
    mkdir -p /etc/$APP_NAME
    
    chown -R $APP_USER:$APP_USER $APP_DIR
    chmod 755 $APP_DIR
}

# Setup database
setup_database() {
    log "Setting up PostgreSQL database..."
    
    # Generate random password
    DB_PASSWORD=$(openssl rand -base64 32)
    
    # Create database and user
    sudo -u postgres psql << EOF
CREATE DATABASE $DB_NAME;
CREATE USER $DB_USER WITH ENCRYPTED PASSWORD '$DB_PASSWORD';
GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;
ALTER USER $DB_USER CREATEDB;
\q
EOF
    
    # Save database credentials
    cat > /etc/$APP_NAME/database.env << EOF
DB_HOST=localhost
DB_PORT=5432
DB_NAME=$DB_NAME
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_SSLMODE=disable
EOF
    
    chmod 600 /etc/$APP_NAME/database.env
    chown $APP_USER:$APP_USER /etc/$APP_NAME/database.env
    
    log "Database setup completed"
}

# Setup environment variables
setup_environment() {
    log "Setting up environment variables..."
    
    # Generate secrets
    JWT_SECRET=$(openssl rand -base64 64)
    ENCRYPTION_KEY=$(openssl rand -base64 32)
    API_KEY=$(openssl rand -base64 32)
    
    # Create main environment file
    cat > /etc/$APP_NAME/app.env << EOF
# Application Configuration
APP_ENV=production
APP_NAME=$APP_NAME
APP_PORT=$BACKEND_PORT
APP_HOST=0.0.0.0
APP_DOMAIN=$DOMAIN

# Security
JWT_SECRET=$JWT_SECRET
ENCRYPTION_KEY=$ENCRYPTION_KEY
API_KEY=$API_KEY

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=$REDIS_PORT
REDIS_PASSWORD=
REDIS_DB=0

# File Upload
UPLOAD_DIR=$APP_DIR/uploads
MAX_FILE_SIZE=10485760
ALLOWED_FILE_TYPES=jpg,jpeg,png,gif,pdf,doc,docx

# Email Configuration (configure with your SMTP settings)
SMTP_HOST=
SMTP_PORT=587
SMTP_USER=
SMTP_PASSWORD=
SMTP_FROM=$EMAIL

# Monitoring
PROMETHEUS_PORT=$PROMETHEUS_PORT
GRAFANA_PORT=$GRAFANA_PORT
METRICS_ENABLED=true
HEALTH_CHECK_INTERVAL=30s

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60s

# Backup Configuration
BACKUP_ENABLED=true
BACKUP_SCHEDULE="0 2 * * *"
BACKUP_RETENTION_DAYS=30
BACKUP_DIR=$APP_DIR/backups
EOF
    
    chmod 600 /etc/$APP_NAME/app.env
    chown $APP_USER:$APP_USER /etc/$APP_NAME/app.env
    
    log "Environment variables configured"
}

# Deploy application code
deploy_application() {
    log "Deploying application code..."
    
    # Clone or copy application code
    if [ -d "/tmp/nutrition-platform" ]; then
        cp -r /tmp/nutrition-platform/* $APP_DIR/
    else
        # If code is not in /tmp, create placeholder structure
        mkdir -p $APP_DIR/backend
        mkdir -p $APP_DIR/frontend
        mkdir -p $APP_DIR/deployment
        
        # Copy deployment files
        cp nginx.conf /etc/nginx/sites-available/$APP_NAME
        ln -sf /etc/nginx/sites-available/$APP_NAME /etc/nginx/sites-enabled/
        rm -f /etc/nginx/sites-enabled/default
    fi
    
    chown -R $APP_USER:$APP_USER $APP_DIR
    
    log "Application code deployed"
}

# Build and setup backend
setup_backend() {
    log "Setting up backend..."
    
    if [ -d "$APP_DIR/backend" ]; then
        cd $APP_DIR/backend
        
        # Build Go application
        sudo -u $APP_USER /usr/local/go/bin/go mod tidy
        sudo -u $APP_USER /usr/local/go/bin/go build -o bin/nutrition-platform ./cmd/server
        
        # Create systemd service
        cat > /etc/systemd/system/$APP_NAME-backend.service << EOF
[Unit]
Description=Nutrition Platform Backend
After=network.target postgresql.service redis.service
Wants=postgresql.service redis.service

[Service]
Type=simple
User=$APP_USER
Group=$APP_USER
WorkingDirectory=$APP_DIR/backend
EnvironmentFile=/etc/$APP_NAME/app.env
EnvironmentFile=/etc/$APP_NAME/database.env
ExecStart=$APP_DIR/backend/bin/nutrition-platform
Restart=always
RestartSec=5
LimitNOFILE=65536

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$APP_DIR

[Install]
WantedBy=multi-user.target
EOF
        
        # Enable and start backend service
        systemctl daemon-reload
        systemctl enable $APP_NAME-backend
        systemctl start $APP_NAME-backend
        
        log "Backend service configured and started"
    else
        warn "Backend directory not found, skipping backend setup"
    fi
}

# Build and setup frontend
setup_frontend() {
    log "Setting up frontend..."
    
    if [ -d "$APP_DIR/frontend" ]; then
        cd $APP_DIR/frontend
        
        # Install dependencies and build
        sudo -u $APP_USER npm install
        sudo -u $APP_USER npm run build
        
        # Copy built files to web directory
        mkdir -p /var/www/$APP_NAME/frontend
        cp -r dist/* /var/www/$APP_NAME/frontend/
        chown -R www-data:www-data /var/www/$APP_NAME/frontend
        
        log "Frontend built and deployed"
    else
        warn "Frontend directory not found, skipping frontend setup"
    fi
}

# Setup SSL certificates
setup_ssl() {
    log "Setting up SSL certificates..."
    
    # Stop Nginx temporarily
    systemctl stop nginx
    
    # Obtain SSL certificate
    certbot certonly --standalone --non-interactive --agree-tos --email $EMAIL -d $DOMAIN -d www.$DOMAIN
    
    # Setup auto-renewal
    echo "0 12 * * * /usr/bin/certbot renew --quiet" | crontab -
    
    # Start Nginx
    systemctl start nginx
    
    log "SSL certificates configured"
}

# Setup monitoring
setup_monitoring() {
    log "Setting up monitoring with Prometheus and Grafana..."
    
    # Create monitoring directory
    mkdir -p $APP_DIR/monitoring
    
    # Create Prometheus configuration
    cat > $APP_DIR/monitoring/prometheus.yml << EOF
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "alert_rules.yml"

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093

scrape_configs:
  - job_name: 'nutrition-platform'
    static_configs:
      - targets: ['localhost:$BACKEND_PORT']
    metrics_path: '/metrics'
    scrape_interval: 5s
    
  - job_name: 'nginx'
    static_configs:
      - targets: ['localhost:80']
    metrics_path: '/nginx_status'
    
  - job_name: 'node-exporter'
    static_configs:
      - targets: ['localhost:9100']
      
  - job_name: 'postgres-exporter'
    static_configs:
      - targets: ['localhost:9187']
EOF
    
    # Create alert rules
    cat > $APP_DIR/monitoring/alert_rules.yml << EOF
groups:
  - name: nutrition-platform-alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          
      - alert: HighResponseTime
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High response time detected"
          
      - alert: DatabaseDown
        expr: up{job="postgres-exporter"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Database is down"
EOF
    
    # Create Docker Compose for monitoring stack
    cat > $APP_DIR/monitoring/docker-compose.yml << EOF
version: '3.8'

services:
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    ports:
      - "$PROMETHEUS_PORT:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - ./alert_rules.yml:/etc/prometheus/alert_rules.yml
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--storage.tsdb.retention.time=200h'
      - '--web.enable-lifecycle'
    restart: unless-stopped
    
  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "$GRAFANA_PORT:3000"
    volumes:
      - grafana_data:/var/lib/grafana
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin123
      - GF_USERS_ALLOW_SIGN_UP=false
    restart: unless-stopped
    
  node-exporter:
    image: prom/node-exporter:latest
    container_name: node-exporter
    ports:
      - "9100:9100"
    volumes:
      - /proc:/host/proc:ro
      - /sys:/host/sys:ro
      - /:/rootfs:ro
    command:
      - '--path.procfs=/host/proc'
      - '--path.rootfs=/rootfs'
      - '--path.sysfs=/host/sys'
      - '--collector.filesystem.mount-points-exclude=^/(sys|proc|dev|host|etc)($$|/)'
    restart: unless-stopped
    
  postgres-exporter:
    image: prometheuscommunity/postgres-exporter:latest
    container_name: postgres-exporter
    ports:
      - "9187:9187"
    environment:
      - DATA_SOURCE_NAME=postgresql://$DB_USER:$(cat /etc/$APP_NAME/database.env | grep DB_PASSWORD | cut -d'=' -f2)@localhost:5432/$DB_NAME?sslmode=disable
    restart: unless-stopped

volumes:
  prometheus_data:
  grafana_data:
EOF
    
    # Start monitoring stack
    cd $APP_DIR/monitoring
    docker-compose up -d
    
    log "Monitoring stack deployed"
}

# Setup backup system
setup_backup() {
    log "Setting up backup system..."
    
    # Create backup script
    cat > $APP_DIR/scripts/backup.sh << EOF
#!/bin/bash

# Load environment variables
source /etc/$APP_NAME/app.env
source /etc/$APP_NAME/database.env

# Create backup directory with timestamp
BACKUP_DATE=\$(date +"%Y%m%d_%H%M%S")
BACKUP_PATH="\$BACKUP_DIR/\$BACKUP_DATE"
mkdir -p "\$BACKUP_PATH"

# Database backup
pg_dump -h \$DB_HOST -p \$DB_PORT -U \$DB_USER -d \$DB_NAME > "\$BACKUP_PATH/database.sql"

# Application files backup
tar -czf "\$BACKUP_PATH/uploads.tar.gz" -C \$APP_DIR uploads
tar -czf "\$BACKUP_PATH/logs.tar.gz" -C \$APP_DIR logs

# Configuration backup
cp -r /etc/$APP_NAME "\$BACKUP_PATH/config"

# Compress entire backup
tar -czf "\$BACKUP_DIR/backup_\$BACKUP_DATE.tar.gz" -C "\$BACKUP_DIR" "\$BACKUP_DATE"
rm -rf "\$BACKUP_PATH"

# Clean old backups (keep last 30 days)
find \$BACKUP_DIR -name "backup_*.tar.gz" -mtime +30 -delete

echo "Backup completed: backup_\$BACKUP_DATE.tar.gz"
EOF
    
    chmod +x $APP_DIR/scripts/backup.sh
    chown $APP_USER:$APP_USER $APP_DIR/scripts/backup.sh
    
    # Add to crontab
    echo "0 2 * * * $APP_DIR/scripts/backup.sh" | crontab -u $APP_USER -
    
    log "Backup system configured"
}

# Setup firewall
setup_firewall() {
    log "Setting up firewall..."
    
    # Install and configure UFW
    apt-get install -y ufw
    
    # Default policies
    ufw default deny incoming
    ufw default allow outgoing
    
    # Allow SSH
    ufw allow ssh
    
    # Allow HTTP and HTTPS
    ufw allow 80/tcp
    ufw allow 443/tcp
    
    # Allow monitoring ports (only from localhost)
    ufw allow from 127.0.0.1 to any port $PROMETHEUS_PORT
    ufw allow from 127.0.0.1 to any port $GRAFANA_PORT
    
    # Enable firewall
    ufw --force enable
    
    log "Firewall configured"
}

# Setup log rotation
setup_log_rotation() {
    log "Setting up log rotation..."
    
    cat > /etc/logrotate.d/$APP_NAME << EOF
$APP_DIR/logs/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 $APP_USER $APP_USER
    postrotate
        systemctl reload $APP_NAME-backend
    endscript
}

/var/log/nginx/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 www-data adm
    postrotate
        systemctl reload nginx
    endscript
}
EOF
    
    log "Log rotation configured"
}

# Health check function
health_check() {
    log "Performing health checks..."
    
    # Check services
    services=("nginx" "postgresql" "redis-server" "$APP_NAME-backend")
    for service in "${services[@]}"; do
        if systemctl is-active --quiet $service; then
            log "✓ $service is running"
        else
            error "✗ $service is not running"
        fi
    done
    
    # Check ports
    ports=("80" "443" "$BACKEND_PORT" "5432" "$REDIS_PORT")
    for port in "${ports[@]}"; do
        if netstat -tuln | grep -q ":$port "; then
            log "✓ Port $port is listening"
        else
            warn "✗ Port $port is not listening"
        fi
    done
    
    # Check SSL certificate
    if [ -f "/etc/letsencrypt/live/$DOMAIN/fullchain.pem" ]; then
        log "✓ SSL certificate exists"
    else
        warn "✗ SSL certificate not found"
    fi
    
    log "Health check completed"
}

# Main deployment function
main() {
    log "Starting Nutrition Platform deployment..."
    
    check_root
    update_system
    install_docker
    install_nginx
    install_certbot
    install_postgresql
    install_redis
    install_nodejs
    install_go
    create_app_user
    setup_app_directory
    setup_database
    setup_environment
    deploy_application
    setup_backend
    setup_frontend
    setup_ssl
    setup_monitoring
    setup_backup
    setup_firewall
    setup_log_rotation
    health_check
    
    log "Deployment completed successfully!"
    log "Application URL: https://$DOMAIN"
    log "Grafana URL: http://$DOMAIN:$GRAFANA_PORT (admin/admin123)"
    log "Prometheus URL: http://$DOMAIN:$PROMETHEUS_PORT"
    
    echo -e "\n${GREEN}Next steps:${NC}"
    echo "1. Configure your DNS to point $DOMAIN to this server"
    echo "2. Update SMTP settings in /etc/$APP_NAME/app.env"
    echo "3. Change default Grafana password"
    echo "4. Review and customize monitoring alerts"
    echo "5. Test the application and monitoring"
}

# Run main function
main "$@"