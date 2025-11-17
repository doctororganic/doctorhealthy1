#!/bin/bash

# Production Deployment Script for Nutrition Platform
# This script handles the complete deployment process

set -e  # Exit on any error

# Configuration
PROJECT_NAME="nutrition-platform"
BACKUP_DIR="/opt/backups/${PROJECT_NAME}"
LOG_FILE="/var/log/${PROJECT_NAME}-deploy.log"
ENV_FILE=".env.production"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_FILE"
}

# Print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root
check_root() {
    if [[ $EUID -eq 0 ]]; then
        print_error "This script should not be run as root"
        exit 1
    fi
}

# Check if Docker is installed
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose is not installed"
        exit 1
    fi
}

# Create backup
create_backup() {
    print_status "Creating backup..."
    BACKUP_DATE=$(date +%Y%m%d_%H%M%S)
    BACKUP_PATH="${BACKUP_DIR}/${BACKUP_DATE}"
    
    mkdir -p "$BACKUP_PATH"
    
    # Backup database
    if [ -f "data/nutrition.db" ]; then
        cp data/nutrition.db "$BACKUP_PATH/"
        log "Database backed up to $BACKUP_PATH/nutrition.db"
    fi
    
    # Backup configuration
    if [ -f "$ENV_FILE" ]; then
        cp "$ENV_FILE" "$BACKUP_PATH/"
        log "Configuration backed up to $BACKUP_PATH/$ENV_FILE"
    fi
    
    # Backup logs
    if [ -d "logs" ]; then
        cp -r logs "$BACKUP_PATH/"
        log "Logs backed up to $BACKUP_PATH/logs"
    fi
    
    print_status "Backup created at $BACKUP_PATH"
}

# Validate environment
validate_environment() {
    print_status "Validating environment..."
    
    if [ ! -f "$ENV_FILE" ]; then
        print_error "Environment file $ENV_FILE not found"
        exit 1
    fi
    
    # Check required environment variables
    source "$ENV_FILE"
    
    if [ -z "$JWT_SECRET" ] || [ ${#JWT_SECRET} -lt 32 ]; then
        print_error "JWT_SECRET must be at least 32 characters"
        exit 1
    fi
    
    print_status "Environment validation passed"
}

# Build Docker images
build_images() {
    print_status "Building Docker images..."
    
    # Build backend
    print_status "Building backend image..."
    docker build -t nutrition-backend:latest backend/
    
    # Build frontend
    print_status "Building frontend image..."
    docker build -t nutrition-frontend:latest frontend/
    
    print_status "Docker images built successfully"
}

# Deploy services
deploy_services() {
    print_status "Deploying services..."
    
    # Stop existing services
    docker-compose down 2>/dev/null || true
    
    # Pull latest images (if any)
    docker-compose pull
    
    # Start services
    docker-compose up -d
    
    # Wait for services to be healthy
    print_status "Waiting for services to be healthy..."
    sleep 30
    
    # Check service health
    if docker-compose ps | grep -q "Up (healthy)"; then
        print_status "Services are healthy"
    else
        print_warning "Some services may not be healthy yet"
        docker-compose ps
    fi
}

# Run health checks
run_health_checks() {
    print_status "Running health checks..."
    
    # Check backend health
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        print_status "Backend health check passed"
    else
        print_error "Backend health check failed"
        exit 1
    fi
    
    # Check frontend
    if curl -f http://localhost:3000 > /dev/null 2>&1; then
        print_status "Frontend health check passed"
    else
        print_warning "Frontend may still be starting"
    fi
    
    # Check database
    if [ -f "data/nutrition.db" ]; then
        print_status "Database file exists"
    else
        print_warning "Database file not found - will be created on first run"
    fi
}

# Cleanup old backups
cleanup_backups() {
    print_status "Cleaning up old backups..."
    
    # Keep last 7 days of backups
    find "$BACKUP_DIR" -type d -mtime +7 -exec rm -rf {} + 2>/dev/null || true
    
    print_status "Backup cleanup completed"
}

# Main deployment function
main() {
    print_status "Starting deployment of Nutrition Platform..."
    
    # Run deployment steps
    check_root
    check_docker
    create_backup
    validate_environment
    build_images
    deploy_services
    run_health_checks
    cleanup_backups
    
    print_status "Deployment completed successfully!"
    print_status "Application is running at: http://localhost"
    print_status "API is available at: http://localhost/api"
    print_status "Health check: http://localhost/health"
    
    # Show running services
    echo
    print_status "Running services:"
    docker-compose ps
}

# Handle script interruption
trap 'print_error "Deployment interrupted"; exit 1' INT

# Run main function
main "$@"
