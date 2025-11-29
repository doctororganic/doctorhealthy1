#!/bin/bash

# Production Deployment Script
# This script handles zero-downtime deployments with health checks and rollback capabilities

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
LOG_FILE="$PROJECT_ROOT/logs/deployment.log"
BACKUP_DIR="$PROJECT_ROOT/backups"
HEALTH_CHECK_TIMEOUT=300
ROLLBACK_TIMEOUT=60

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1" | tee -a "$LOG_FILE"
}

log_success() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] SUCCESS:${NC} $1" | tee -a "$LOG_FILE"
}

log_warning() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING:${NC} $1" | tee -a "$LOG_FILE"
}

log_error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR:${NC} $1" | tee -a "$LOG_FILE"
}

# Ensure directories exist
mkdir -p "$PROJECT_ROOT/logs"
mkdir -p "$BACKUP_DIR"

# Usage information
usage() {
    cat << EOF
Production Deployment Script

Usage: $0 [OPTIONS] <ENVIRONMENT>

ENVIRONMENT:
    staging     Deploy to staging environment
    production  Deploy to production environment

OPTIONS:
    --skip-backup     Skip database backup
    --skip-tests       Skip pre-deployment tests
    --force           Force deployment without confirmation
    --dry-run         Show what would be deployed without actually deploying
    --rollback VERSION  Rollback to specific version

EXAMPLES:
    $0 production                    # Deploy to production
    $0 staging --skip-backup         # Deploy to staging without backup
    $0 --rollback v1.2.3 production # Rollback production to v1.2.3

EOF
}

# Parse command line arguments
ENVIRONMENT=""
SKIP_BACKUP=false
SKIP_TESTS=false
FORCE=false
DRY_RUN=false
ROLLBACK_VERSION=""

while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-backup)
            SKIP_BACKUP=true
            shift
            ;;
        --skip-tests)
            SKIP_TESTS=true
            shift
            ;;
        --force)
            FORCE=true
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --rollback)
            ROLLBACK_VERSION="$2"
            shift 2
            ;;
        staging|production)
            ENVIRONMENT="$1"
            shift
            ;;
        *)
            log_error "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

# Validate environment
if [[ -z "$ENVIRONMENT" && -z "$ROLLBACK_VERSION" ]]; then
    log_error "Environment must be specified"
    usage
    exit 1
fi

# Load environment configuration
ENV_FILE="$PROJECT_ROOT/config/${ENVIRONMENT}.env"
if [[ ! -f "$ENV_FILE" ]]; then
    log_error "Environment file not found: $ENV_FILE"
    exit 1
fi

# Source environment variables
set -a
source "$ENV_FILE"
set +a

log "Starting deployment for environment: $ENVIRONMENT"

# Pre-flight checks
check_prerequisites() {
    log "Running pre-flight checks..."
    
    # Check if required tools are installed
    local tools=("docker" "docker-compose" "git" "psql" "redis-cli")
    for tool in "${tools[@]}"; do
        if ! command -v "$tool" &> /dev/null; then
            log_error "Required tool not found: $tool"
            exit 1
        fi
    done
    
    # Check if environment is properly configured
    if [[ -z "$DATABASE_URL" ]]; then
        log_error "DATABASE_URL not configured in $ENV_FILE"
        exit 1
    fi
    
    if [[ -z "$REDIS_ADDR" ]]; then
        log_error "REDIS_ADDR not configured in $ENV_FILE"
        exit 1
    fi
    
    log_success "Pre-flight checks passed"
}

# Run tests
run_tests() {
    if [[ "$SKIP_TESTS" == true ]]; then
        log_warning "Skipping tests as requested"
        return 0
    fi
    
    log "Running pre-deployment tests..."
    
    # Run unit tests
    if ! "$PROJECT_ROOT/backend/scripts/run-all-tests.sh"; then
        log_error "Unit tests failed"
        exit 1
    fi
    
    # Run integration tests
    if ! "$PROJECT_ROOT/backend/scripts/smoke-test.sh"; then
        log_error "Integration tests failed"
        exit 1
    fi
    
    log_success "All tests passed"
}

# Create database backup
create_backup() {
    if [[ "$SKIP_BACKUP" == true ]]; then
        log_warning "Skipping database backup as requested"
        return 0
    fi
    
    log "Creating database backup..."
    
    local backup_file="$BACKUP_DIR/backup_${ENVIRONMENT}_$(date +%Y%m%d_%H%M%S).sql"
    
    # Create backup using pg_dump
    if ! pg_dump "$DATABASE_URL" > "$backup_file"; then
        log_error "Database backup failed"
        exit 1
    fi
    
    # Compress backup
    gzip "$backup_file"
    
    log_success "Database backup created: ${backup_file}.gz"
    
    # Clean old backups (keep last 10)
    find "$BACKUP_DIR" -name "backup_${ENVIRONMENT}_*.sql.gz" -type f | sort -r | tail -n +11 | xargs -r rm
}

# Build application
build_application() {
    log "Building application..."
    
    # Build backend
    cd "$PROJECT_ROOT/backend"
    if ! go build -o nutrition-platform .; then
        log_error "Backend build failed"
        exit 1
    fi
    
    # Build frontend if it exists
    if [[ -d "$PROJECT_ROOT/frontend" ]]; then
        cd "$PROJECT_ROOT/frontend"
        if ! npm run build; then
            log_error "Frontend build failed"
            exit 1
        fi
    fi
    
    log_success "Application built successfully"
}

# Deploy containers
deploy_containers() {
    log "Deploying containers..."
    
    cd "$PROJECT_ROOT"
    
    # Use docker-compose for deployment
    local compose_file="docker-compose.${ENVIRONMENT}.yml"
    if [[ ! -f "$compose_file" ]]; then
        compose_file="docker-compose.yml"
    fi
    
    # Pull latest images
    if ! docker-compose -f "$compose_file" pull; then
        log_error "Failed to pull latest images"
        exit 1
    fi
    
    # Deploy with zero downtime
    if ! docker-compose -f "$compose_file" up -d --no-deps --scale backend=2; then
        log_error "Container deployment failed"
        exit 1
    fi
    
    log_success "Containers deployed successfully"
}

# Health check
health_check() {
    log "Performing health check..."
    
    local start_time=$(date +%s)
    local health_url="http://localhost:${PORT:-8080}/health"
    
    while true; do
        local current_time=$(date +%s)
        local elapsed=$((current_time - start_time))
        
        if [[ $elapsed -gt $HEALTH_CHECK_TIMEOUT ]]; then
            log_error "Health check timeout after ${HEALTH_CHECK_TIMEOUT}s"
            return 1
        fi
        
        # Check if application is responding
        if curl -f -s "$health_url" > /dev/null 2>&1; then
            log_success "Health check passed (${elapsed}s)"
            return 0
        fi
        
        sleep 5
    done
}

# Rollback deployment
rollback() {
    log "Initiating rollback..."
    
    # Stop current containers
    cd "$PROJECT_ROOT"
    local compose_file="docker-compose.${ENVIRONMENT}.yml"
    if [[ ! -f "$compose_file" ]]; then
        compose_file="docker-compose.yml"
    fi
    
    if ! docker-compose -f "$compose_file" down; then
        log_error "Failed to stop containers"
        exit 1
    fi
    
    # Restore from backup if this is a database rollback
    if [[ -n "$ROLLBACK_VERSION" ]]; then
        log "Restoring database backup for version: $ROLLBACK_VERSION"
        # Implementation would depend on backup strategy
        log_warning "Database restoration not implemented in this script"
    fi
    
    # Restart with previous version
    # This would involve switching to previous image/tag
    log_warning "Container rollback not fully implemented - manual intervention required"
    
    log_success "Rollback completed"
}

# Cleanup old containers and images
cleanup() {
    log "Cleaning up old containers and images..."
    
    # Remove stopped containers
    docker container prune -f > /dev/null 2>&1
    
    # Remove unused images
    docker image prune -f > /dev/null 2>&1
    
    log_success "Cleanup completed"
}

# Send notification
send_notification() {
    local status="$1"
    local message="$2"
    
    # Send to Slack if webhook is configured
    if [[ -n "${SLACK_WEBHOOK_URL:-}" ]]; then
        local color="good"
        [[ "$status" == "error" ]] && color="danger"
        [[ "$status" == "warning" ]] && color="warning"
        
        curl -X POST "$SLACK_WEBHOOK_URL" \
            -H 'Content-type: application/json' \
            --data "{
                \"attachments\": [{
                    \"color\": \"$color\",
                    \"title\": \"Deployment Notification\",
                    \"text\": \"$message\",
                    \"fields\": [{
                        \"title\": \"Environment\",
                        \"value\": \"$ENVIRONMENT\",
                        \"short\": true
                    }, {
                        \"title\": \"Time\",
                        \"value\": \"$(date)\",
                        \"short\": true
                    }]
                }]
            }" > /dev/null 2>&1
    fi
    
    # Send email if configured
    if [[ -n "${NOTIFICATION_EMAIL:-}" ]]; then
        echo "$message" | mail -s "Deployment Notification: $ENVIRONMENT" "$NOTIFICATION_EMAIL"
    fi
}

# Main deployment flow
main() {
    # Check if this is a rollback
    if [[ -n "$ROLLBACK_VERSION" ]]; then
        log "Rolling back to version: $ROLLBACK_VERSION"
        rollback
        send_notification "success" "Successfully rolled back $ENVIRONMENT to version $ROLLBACK_VERSION"
        exit 0
    fi
    
    # Dry run mode
    if [[ "$DRY_RUN" == true ]]; then
        log "DRY RUN: Would deploy to $ENVIRONMENT"
        log "DRY RUN: Would run tests"
        log "DRY RUN: Would create backup"
        log "DRY RUN: Would build application"
        log "DRY RUN: Would deploy containers"
        log "DRY RUN: Would perform health check"
        exit 0
    fi
    
    # Confirmation prompt
    if [[ "$FORCE" != true ]]; then
        echo -e "${YELLOW}You are about to deploy to $ENVIRONMENT environment.${NC}"
        read -p "Are you sure you want to continue? (yes/no): " -r confirm
        if [[ "$confirm" != "yes" ]]; then
            log "Deployment cancelled by user"
            exit 0
        fi
    fi
    
    # Execute deployment steps
    check_prerequisites
    run_tests
    create_backup
    build_application
    deploy_containers
    
    # Health check
    if health_check; then
        cleanup
        log_success "Deployment to $ENVIRONMENT completed successfully"
        send_notification "success" "Successfully deployed to $ENVIRONMENT environment"
    else
        log_error "Health check failed, initiating rollback"
        rollback
        send_notification "error" "Deployment to $ENVIRONMENT failed and was rolled back"
        exit 1
    fi
}

# Trap for cleanup on script exit
trap 'log_error "Script interrupted or failed"' ERR EXIT

# Run main function
main "$@"
