# VPS Deployment Script for Hostinger
# Domain: ieltspass1.com
# Author: Nutrition Platform Team

# Set error action preference to stop on errors
$ErrorActionPreference = "Stop"

# Colors for output
$Red = [ConsoleColor]::Red
$Green = [ConsoleColor]::Green
$Yellow = [ConsoleColor]::Yellow
$Blue = [ConsoleColor]::Blue
$NC = [ConsoleColor]::White

# Configuration
$DOMAIN = "ieltspass1.com"
$WWW_DOMAIN = "www.ieltspass1.com"
$EMAIL = "Khaledalzayat278@gmail.com"
$APP_DIR = "/opt/nutrition-platform"
$BACKUP_DIR = "/opt/backups"

# Functions
function log_info($message) {
    Write-Host "[INFO]$NC $message" -ForegroundColor $Blue
}

function log_success($message) {
    Write-Host "[SUCCESS]$NC $message" -ForegroundColor $Green
}

function log_warning($message) {
    Write-Host "[WARNING]$NC $message" -ForegroundColor $Yellow
}

function log_error($message) {
    Write-Host "[ERROR]$NC $message" -ForegroundColor $Red
}

function check_root {
    if ((id -u) -eq 0) {
        log_error "This script should not be run as root. Please run as a regular user with sudo privileges."
        exit 1
    }
}

function check_prerequisites {
    log_info "Checking prerequisites..."
    
    # Check if running on Ubuntu/Debian
    if (-not (Get-Command apt -ErrorAction SilentlyContinue)) {
        log_error "This script is designed for Ubuntu/Debian systems."
        exit 1
    }
    
    # Check sudo privileges
    try {
        & sudo -n true 2>$null
    } catch {
        log_error "This script requires sudo privileges."
        exit 1
    }
    
    log_success "Prerequisites check passed."
}

function install_dependencies {
    log_info "Installing system dependencies..."
    
    # Update system
    & sudo apt update
    & sudo apt upgrade -y
    
    # Install required packages
    & sudo apt install -y nginx docker.io docker-compose git certbot python3-certbot-nginx ufw curl htop unzip
    
    # Start and enable services
    & sudo systemctl start nginx
    & sudo systemctl enable nginx
    & sudo systemctl start docker
    & sudo systemctl enable docker
    
    # Add user to docker group
    & sudo usermod -aG docker $env:USER
    
    log_success "Dependencies installed successfully."
}

function configure_firewall {
    log_info "Configuring firewall..."
    
    # Reset UFW to defaults
    & sudo ufw --force reset
    
    # Set default policies
    & sudo ufw default deny incoming
    & sudo ufw default allow outgoing
    
    # Allow SSH, HTTP, HTTPS, and backend API
    & sudo ufw allow OpenSSH
    & sudo ufw allow 'Nginx Full'
    & sudo ufw allow 80/tcp
    & sudo ufw allow 443/tcp
    & sudo ufw allow 8080/tcp
    & sudo ufw allow 8081/tcp
    
    # Enable firewall
    & sudo ufw --force enable
    
    log_success "Firewall configured successfully."
}

function setup_application {
    log_info "Setting up application..."
    
    # Create application directory
    & sudo mkdir -p $APP_DIR
    & sudo chown -R "$($env:USER):$($env:USER)" $APP_DIR
    
    # Copy application files
    if (Test-Path "./nutrition-platform") {
        Copy-Item -Path "./nutrition-platform/*" -Destination $APP_DIR -Recurse
    } else {
        Copy-Item -Path "./*" -Destination $APP_DIR -Recurse
    }
    
    # Set up environment file
    if (-not (Test-Path "$APP_DIR/.env")) {
        Copy-Item -Path "$APP_DIR/.env.vps" -Destination "$APP_DIR/.env"
        log_info "Using VPS-specific environment configuration"
    } else {
        log_warning "Environment file already exists. Backing up and updating..."
        Copy-Item -Path "$APP_DIR/.env" -Destination "$APP_DIR/.env.backup.$((Get-Date).ToString('yyyyMMdd_HHmmss'))"
        Copy-Item -Path "$APP_DIR/.env.vps" -Destination "$APP_DIR/.env"
    }
    
    # Update domain-specific configurations
    $envContent = Get-Content "$APP_DIR/.env" -Raw
    $envContent = $envContent -replace "DOMAIN=.*", "DOMAIN=$DOMAIN"
    $envContent = $envContent -replace "WWW_DOMAIN=.*", "WWW_DOMAIN=$WWW_DOMAIN"
    Set-Content -Path "$APP_DIR/.env" -Value $envContent
    
    log_success "Application setup completed."
}

function build_containers {
    log_info "Building Docker containers..."
    
    Set-Location $APP_DIR
    
    # Build and start containers using VPS configuration
    & docker-compose -f docker-compose.vps.yml down 2>$null
    & docker-compose -f docker-compose.vps.yml build --no-cache
    & docker-compose -f docker-compose.vps.yml up -d
    
    # Wait for containers to be ready
    log_info "Waiting for containers to start..."
    Start-Sleep -Seconds 30
    
    # Check container status
    $containerStatus = & docker-compose -f docker-compose.vps.yml ps
    if ($containerStatus | Select-String -Pattern "Up") {
        log_success "Containers are running successfully."
    } else {
        log_error "Some containers failed to start. Check logs with: docker-compose -f docker-compose.vps.yml logs"
        exit 1
    }
}

function configure_nginx {
    log_info "Configuring Nginx..."
    
    # Remove default site
    & sudo rm -f /etc/nginx/sites-enabled/default
    
    # Create site configuration
    $nginxConfig = @"
# HTTP server - redirect to HTTPS
server {
    listen 80;
    server_name $DOMAIN $WWW_DOMAIN;
    return 301 https://`$server_name`$request_uri;
}

# HTTPS server
server {
    listen 443 ssl http2;
    server_name $DOMAIN $WWW_DOMAIN;
    
    # SSL configuration (will be configured by certbot)
    ssl_certificate /etc/letsencrypt/live/$DOMAIN/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/$DOMAIN/privkey.pem;
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
        proxy_set_header Upgrade `$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host `$host;
        proxy_set_header X-Real-IP `$remote_addr;
        proxy_set_header X-Forwarded-For `$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto `$scheme;
        proxy_cache_bypass `$http_upgrade;
    }
    
    # Backend API
    location /api/ {
        proxy_pass http://localhost:8080/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade `$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host `$host;
        proxy_set_header X-Real-IP `$remote_addr;
        proxy_set_header X-Forwarded-For `$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto `$scheme;
        proxy_cache_bypass `$http_upgrade;
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
"@
    
    & sudo bash -c "echo '$nginxConfig' > /etc/nginx/sites-available/$DOMAIN"
    
    # Enable site
    & sudo ln -sf "/etc/nginx/sites-available/$DOMAIN" "/etc/nginx/sites-enabled/"
    
    # Test nginx configuration
    try {
        & sudo nginx -t
        & sudo systemctl reload nginx
        log_success "Nginx configured successfully."
    } catch {
        log_error "Nginx configuration test failed."
        exit 1
    }
}

function setup_ssl {
    log_info "Setting up SSL certificate..."
    
    # Get SSL certificate
    try {
        & sudo certbot --nginx -d $DOMAIN -d $WWW_DOMAIN --non-interactive --agree-tos --email $EMAIL
        log_success "SSL certificate obtained successfully."
        
        # Test auto-renewal
        try {
            & sudo certbot renew --dry-run
            log_success "SSL auto-renewal test passed."
        } catch {
            log_warning "SSL auto-renewal test failed. Please check manually."
        }
        
        # Add cron job for auto-renewal
        $cronJob = "0 12 * * * /usr/bin/certbot renew --quiet"
        (Get-Content /etc/crontab -Raw; $cronJob) | Set-Content /etc/crontab
        
    } catch {
        log_error "Failed to obtain SSL certificate. Please check domain DNS settings."
        log_info "You can run SSL setup manually later with: sudo certbot --nginx -d $DOMAIN -d $WWW_DOMAIN"
    }
}

function setup_monitoring {
    log_info "Setting up monitoring and backup..."
    
    # Create backup directory
    & sudo mkdir -p $BACKUP_DIR
    & sudo chown "$($env:USER):$($env:USER)" $BACKUP_DIR
    
    # Create backup script
    $backupScript = @'
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
'@
    
    & sudo bash -c "echo '$backupScript' > /opt/backup.sh"
    
    # Make backup script executable
    & sudo chmod +x /opt/backup.sh
    
    # Add backup to cron (daily at 2 AM)
    $cronBackup = "0 2 * * * /opt/backup.sh"
    (crontab -l 2>/dev/null; $cronBackup) | crontab
    
    log_success "Monitoring and backup setup completed."
}

function show_status {
    log_info "Deployment Status:"
    Write-Host "==========================================" -ForegroundColor $Cyan
    Write-Host "Domain: https://$DOMAIN" -ForegroundColor $Cyan
    Write-Host "WWW Domain: https://$WWW_DOMAIN" -ForegroundColor $Cyan
    Write-Host "Application Directory: $APP_DIR" -ForegroundColor $Cyan
    Write-Host "Backup Directory: $BACKUP_DIR" -ForegroundColor $Cyan
    Write-Host "" -ForegroundColor $Cyan
    Write-Host "Services Status:" -ForegroundColor $Cyan
    Write-Host "- Nginx: $(& systemctl is-active nginx)" -ForegroundColor $Cyan
    Write-Host "- Docker: $(& systemctl is-active docker)" -ForegroundColor $Cyan
    Write-Host "" -ForegroundColor $Cyan
    Write-Host "Container Status:" -ForegroundColor $Cyan
    Set-Location $APP_DIR
    & docker-compose -f docker-compose.vps.yml ps
    Write-Host "" -ForegroundColor $Cyan
    Write-Host "SSL Certificate Status:" -ForegroundColor $Cyan
    try {
        & sudo certbot certificates 2>$null | Select-String -Pattern $DOMAIN -Context 0,5
    } catch {
        Write-Host "No SSL certificate found" -ForegroundColor $Cyan
    }
    Write-Host "==========================================" -ForegroundColor $Cyan
}

function show_next_steps {
    log_success "Deployment completed successfully!"
    Write-Host "" -ForegroundColor $Green
    Write-Host "Next Steps:" -ForegroundColor $Green
    Write-Host "1. Verify your domain DNS points to this server's IP address" -ForegroundColor $Green
    Write-Host "2. Visit https://$DOMAIN to test the application" -ForegroundColor $Green
    Write-Host "3. Check logs if needed: docker-compose logs -f" -ForegroundColor $Green
    Write-Host "4. Monitor system resources: htop" -ForegroundColor $Green
    Write-Host "5. Check SSL certificate: sudo certbot certificates" -ForegroundColor $Green
    Write-Host "" -ForegroundColor $Green
    Write-Host "Useful Commands:" -ForegroundColor $Green
    Write-Host "- Restart containers: cd $APP_DIR && docker-compose restart" -ForegroundColor $Green
    Write-Host "- View logs: cd $APP_DIR && docker-compose logs -f" -ForegroundColor $Green
    Write-Host "- Update SSL: sudo certbot renew" -ForegroundColor $Green
    Write-Host "- Backup manually: sudo /opt/backup.sh" -ForegroundColor $Green
    Write-Host "" -ForegroundColor $Green
    Write-Host "Support: Check VPS_HOSTINGER_DEPLOYMENT.md for detailed documentation" -ForegroundColor $Green
}

# Main execution
function main {
    log_info "Starting VPS deployment for $DOMAIN..."
    
    check_root
    check_prerequisites
    install_dependencies
    configure_firewall
    setup_application
    build_containers
    configure_nginx
    setup_ssl
    setup_monitoring
    
    show_status
    show_next_steps
}

# Run main function
main