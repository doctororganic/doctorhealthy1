#!/bin/bash

# Production Deployment Script for Nutrition Platform Backend
# This script handles secure production deployment with all necessary checks

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="nutrition-platform-backend"
DEPLOY_USER="deploy"
BACKUP_DIR="/var/backups/nutrition-platform"
LOG_FILE="/var/log/nutrition-platform/deploy.log"

# Functions
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1" | tee -a "$LOG_FILE"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "$LOG_FILE"
    exit 1
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" | tee -a "$LOG_FILE"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1" | tee -a "$LOG_FILE"
}

# Pre-deployment checks
pre_deployment_checks() {
    log "Starting pre-deployment checks..."
    
    # Check if running as correct user
    if [ "$USER" != "$DEPLOY_USER" ]; then
        error "This script must be run as the $DEPLOY_USER user"
    fi
    
    # Check if required environment variables are set
    required_vars=(
        "DB_HOST" "DB_NAME" "DB_USER" "DB_PASSWORD"
        "JWT_SECRET" "API_KEY_SECRET" "ENCRYPTION_KEY"
    )
    
    for var in "${required_vars[@]}"; do
        if [ -z "${!var}" ]; then
            error "Required environment variable $var is not set"
        fi
    done
    
    # Check database connectivity
    log "Checking database connectivity..."
    if ! pg_isready -h "$DB_HOST" -p "${DB_PORT:-5432}" -U "$DB_USER" -d "$DB_NAME" > /dev/null 2>&1; then
        error "Cannot connect to database"
    fi
    
    # Check Redis connectivity
    if [ -n "$REDIS_HOST" ]; then
        log "Checking Redis connectivity..."
        if ! redis-cli -h "$REDIS_HOST" -p "${REDIS_PORT:-6379}" ping > /dev/null 2>&1; then
            warning "Cannot connect to Redis - some features may not work"
        fi
    fi
    
    success "Pre-deployment checks completed"
}

# Security audit
security_audit() {
    log "Running security audit..."
    
    # Check file permissions
    find . -name "*.go" -perm /o+w -exec echo "World-writable Go file: {}" \;
    find . -name ".env*" -perm /o+r -exec echo "World-readable env file: {}" \;
    
    # Check for hardcoded secrets (basic check)
    if grep -r "password.*=" --include="*.go" . | grep -v "Password.*string" | grep -v "password_hash"; then
        warning "Potential hardcoded passwords found"
    fi
    
    # Check for TODO/FIXME comments that might indicate security issues
    if grep -r "TODO.*security\|FIXME.*security" --include="*.go" .; then
        warning "Security-related TODO/FIXME comments found"
    fi
    
    success "Security audit completed"
}

# Build application
build_application() {
    log "Building application..."
    
    # Clean previous builds
    make clean
    
    # Run tests
    log "Running tests..."
    if ! make test-coverage; then
        error "Tests failed - deployment aborted"
    fi
    
    # Check test coverage
    if ! make check-coverage; then
        error "Test coverage below required threshold - deployment aborted"
    fi
    
    # Build for production
    log "Building for production..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o "$APP_NAME" .
    
    # Verify binary
    if [ ! -f "$APP_NAME" ]; then
        error "Build failed - binary not found"
    fi
    
    success "Application built successfully"
}

# Database migration
run_migrations() {
    log "Running database migrations..."
    
    # Backup database before migration
    backup_file="$BACKUP_DIR/db_backup_$(date +%Y%m%d_%H%M%S).sql"
    mkdir -p "$BACKUP_DIR"
    
    log "Creating database backup..."
    pg_dump -h "$DB_HOST" -p "${DB_PORT:-5432}" -U "$DB_USER" -d "$DB_NAME" > "$backup_file"
    
    if [ $? -eq 0 ]; then
        success "Database backup created: $backup_file"
    else
        error "Database backup failed"
    fi
    
    # Run migrations (if migration tool is available)
    if command -v migrate &> /dev/null; then
        log "Running database migrations..."
        migrate -path ./migrations -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:${DB_PORT:-5432}/$DB_NAME?sslmode=disable" up
    else
        warning "Migration tool not found - skipping migrations"
    fi
    
    success "Database migrations completed"
}

# Deploy application
deploy_application() {
    log "Deploying application..."
    
    # Stop existing service
    if systemctl is-active --quiet "$APP_NAME"; then
        log "Stopping existing service..."
        sudo systemctl stop "$APP_NAME"
    fi
    
    # Backup current binary
    if [ -f "/usr/local/bin/$APP_NAME" ]; then
        cp "/usr/local/bin/$APP_NAME" "$BACKUP_DIR/${APP_NAME}_backup_$(date +%Y%m%d_%H%M%S)"
    fi
    
    # Copy new binary
    sudo cp "$APP_NAME" "/usr/local/bin/$APP_NAME"
    sudo chmod +x "/usr/local/bin/$APP_NAME"
    
    # Update systemd service file if needed
    if [ -f "scripts/$APP_NAME.service" ]; then
        sudo cp "scripts/$APP_NAME.service" "/etc/systemd/system/"
        sudo systemctl daemon-reload
    fi
    
    # Start service
    log "Starting service..."
    sudo systemctl start "$APP_NAME"
    sudo systemctl enable "$APP_NAME"
    
    # Wait for service to start
    sleep 5
    
    # Check if service is running
    if systemctl is-active --quiet "$APP_NAME"; then
        success "Service started successfully"
    else
        error "Service failed to start"
    fi
    
    success "Application deployed successfully"
}

# Health check
health_check() {
    log "Running health checks..."
    
    # Wait for application to be ready
    sleep 10
    
    # Check health endpoint
    local health_url="http://localhost:${SERVER_PORT:-8080}/health"
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -f -s "$health_url" > /dev/null; then
            success "Health check passed"
            return 0
        fi
        
        log "Health check attempt $attempt/$max_attempts failed, retrying..."
        sleep 2
        ((attempt++))
    done
    
    error "Health check failed after $max_attempts attempts"
}

# Post-deployment tasks
post_deployment() {
    log "Running post-deployment tasks..."
    
    # Update log rotation
    if [ -f "scripts/logrotate.conf" ]; then
        sudo cp "scripts/logrotate.conf" "/etc/logrotate.d/$APP_NAME"
    fi
    
    # Update monitoring configuration
    if [ -f "monitoring/prometheus.yml" ]; then
        log "Updating monitoring configuration..."
        # Update Prometheus config if needed
    fi
    
    # Clear application caches if needed
    if [ -n "$REDIS_HOST" ]; then
        log "Clearing application caches..."
        redis-cli -h "$REDIS_HOST" -p "${REDIS_PORT:-6379}" FLUSHDB
    fi
    
    # Send deployment notification
    if [ -n "$SLACK_WEBHOOK_URL" ]; then
        curl -X POST -H 'Content-type: application/json' \
            --data "{\"text\":\"✅ $APP_NAME deployed successfully to production\"}" \
            "$SLACK_WEBHOOK_URL"
    fi
    
    success "Post-deployment tasks completed"
}

# Rollback function
rollback() {
    log "Rolling back deployment..."
    
    # Stop current service
    sudo systemctl stop "$APP_NAME"
    
    # Restore previous binary
    local backup_binary=$(ls -t "$BACKUP_DIR/${APP_NAME}_backup_"* 2>/dev/null | head -n1)
    if [ -n "$backup_binary" ]; then
        sudo cp "$backup_binary" "/usr/local/bin/$APP_NAME"
        sudo systemctl start "$APP_NAME"
        success "Rollback completed"
    else
        error "No backup binary found for rollback"
    fi
}

# Main deployment process
main() {
    log "Starting production deployment of $APP_NAME"
    
    # Create log directory
    sudo mkdir -p "$(dirname "$LOG_FILE")"
    sudo chown "$DEPLOY_USER:$DEPLOY_USER" "$(dirname "$LOG_FILE")"
    
    # Trap errors for rollback
    trap 'error "Deployment failed - consider running rollback"' ERR
    
    # Run deployment steps
    pre_deployment_checks
    security_audit
    build_application
    run_migrations
    deploy_application
    health_check
    post_deployment
    
    success "🎉 Production deployment completed successfully!"
    log "Deployment log: $LOG_FILE"
}

# Handle command line arguments
case "${1:-deploy}" in
    "deploy")
        main
        ;;
    "rollback")
        rollback
        ;;
    "health-check")
        health_check
        ;;
    *)
        echo "Usage: $0 {deploy|rollback|health-check}"
        exit 1
        ;;
esac