#!/bin/bash

# ============================================
# FINAL DEPLOYMENT SCRIPT
# Trae New Healthy1 - Nutrition Platform
# Version: 2.0
# Last Updated: $(date +%Y-%m-%d)
# ============================================

# Enable strict error handling
set -euo pipefail

# Script configuration
readonly SCRIPT_NAME="$(basename "$0")"
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly LOG_FILE="/tmp/${SCRIPT_NAME%.*}.log"
readonly DOCKER_BUILD_LOG="/tmp/docker-build.log"
readonly TEST_CONTAINER_NAME="trae-test"
readonly TEST_IMAGE_NAME="trae-healthy1-test"
readonly TEST_PORT="8080"
readonly APP_PORT="3000"

# Deployment configuration
readonly REQUIRED_FILES=(
    "Dockerfile"
    ".dockerignore"
    "production-nodejs/server.js"
    "production-nodejs/package.json"
)

# Environment variables
export NODE_ENV="${NODE_ENV:-production}"
export PORT="${PORT:-$APP_PORT}"
export HOST="${HOST:-0.0.0.0}"

# ============================================
# COLOR DEFINITIONS
# ============================================
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly PURPLE='\033[0;35m'
readonly CYAN='\033[0;36m'
readonly NC='\033[0m' # No Color

# ============================================
# LOGGING FUNCTIONS
# ============================================
log() {
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo -e "${BLUE}[INFO]${NC} $1" | tee -a "$LOG_FILE"
    echo "[$timestamp] [INFO] $1" >> "$LOG_FILE"
}

success() {
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo -e "${GREEN}[SUCCESS]${NC} $1" | tee -a "$LOG_FILE"
    echo "[$timestamp] [SUCCESS] $1" >> "$LOG_FILE"
}

error() {
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "$LOG_FILE"
    echo "[$timestamp] [ERROR] $1" >> "$LOG_FILE" >&2
}

warning() {
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo -e "${YELLOW}[WARNING]${NC} $1" | tee -a "$LOG_FILE"
    echo "[$timestamp] [WARNING] $1" >> "$LOG_FILE"
}

step() {
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo -e "${PURPLE}[STEP]${NC} $1" | tee -a "$LOG_FILE"
    echo "[$timestamp] [STEP] $1" >> "$LOG_FILE"
}

# ============================================
# UTILITY FUNCTIONS
# ============================================
cleanup() {
    local exit_code=$?
    log "Cleaning up resources..."
    
    # Stop and remove test container if it exists
    if docker ps -q -f name="$TEST_CONTAINER_NAME" | grep -q .; then
        log "Stopping test container..."
        docker stop "$TEST_CONTAINER_NAME" >/dev/null 2>&1 || true
    fi
    
    # Remove test container if it exists
    if docker ps -a -q -f name="$TEST_CONTAINER_NAME" | grep -q .; then
        log "Removing test container..."
        docker rm "$TEST_CONTAINER_NAME" >/dev/null 2>&1 || true
    fi
    
    # Remove test image if it exists
    if docker images -q "$TEST_IMAGE_NAME" | grep -q .; then
        log "Removing test image..."
        docker rmi "$TEST_IMAGE_NAME" >/dev/null 2>&1 || true
    fi
    
    if [ $exit_code -ne 0 ]; then
        error "Script failed with exit code $exit_code"
        error "Check logs at: $LOG_FILE"
    fi
    
    exit $exit_code
}

check_command() {
    local cmd="$1"
    local name="${2:-$cmd}"
    
    if command -v "$cmd" &> /dev/null; then
        log "$name is available"
        return 0
    else
        error "$name is not installed or not in PATH"
        return 1
    fi
}

wait_for_container() {
    local container_name="$1"
    local max_attempts="${2:-30}"
    local attempt=1
    
    log "Waiting for container to start..."
    
    while [ $attempt -le $max_attempts ]; do
        if docker ps -f name="$container_name" --format "table {{.Status}}" | grep -q "Up"; then
            success "Container is running"
            return 0
        fi
        
        log "Attempt $attempt/$max_attempts: Container not ready yet..."
        sleep 2
        ((attempt++))
    done
    
    error "Container failed to start within ${max_attempts} attempts"
    return 1
}

# ============================================
# VALIDATION FUNCTIONS
# ============================================
validate_prerequisites() {
    log "Validating prerequisites..."
    
    # Check required commands
    local required_commands=("docker" "curl")
    for cmd in "${required_commands[@]}"; do
        check_command "$cmd" || exit 1
    done
    
    # Check if Docker daemon is running
    if ! docker info >/dev/null 2>&1; then
        error "Docker daemon is not running"
        exit 1
    fi
    
    success "All prerequisites validated"
}

validate_files() {
    log "Checking required files..."
    local all_exist=true
    
    for file in "${REQUIRED_FILES[@]}"; do
        if [ -f "$file" ]; then
            success "✓ $file"
        else
            error "✗ $file missing"
            all_exist=false
        fi
    done
    
    if [ "$all_exist" = false ]; then
        error "Missing required files. Cannot proceed."
        exit 1
    fi
    
    success "All required files present"
}

validate_code() {
    log "Validating code..."
    
    if check_command "node" "Node.js"; then
        log "Checking Node.js syntax..."
        if node --check production-nodejs/server.js 2>>"$LOG_FILE"; then
            success "✓ Node.js syntax valid"
        else
            error "✗ Node.js syntax error. Check $LOG_FILE for details."
            exit 1
        fi
    else
        warning "Node.js not found, skipping syntax check"
    fi
    
    success "Code validation complete"
}

# ============================================
# DOCKER FUNCTIONS
# ============================================
test_docker_build() {
    log "Building Docker image (this may take 2-3 minutes)..."
    
    if docker build -t "$TEST_IMAGE_NAME" -f Dockerfile . > "$DOCKER_BUILD_LOG" 2>&1; then
        success "✓ Docker build successful"
    else
        error "✗ Docker build failed"
        error "Check $DOCKER_BUILD_LOG for details"
        exit 1
    fi
}

test_container() {
    log "Starting test container..."
    
    if docker run -d -p "${TEST_PORT}:${APP_PORT}" --name "$TEST_CONTAINER_NAME" "$TEST_IMAGE_NAME" >>"$LOG_FILE" 2>&1; then
        success "Test container started"
    else
        error "Failed to start test container"
        exit 1
    fi
    
    # Wait for container to be ready
    wait_for_container "$TEST_CONTAINER_NAME" || exit 1
    
    log "Testing health endpoint..."
    local max_attempts=10
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -f --max-time 10 "http://localhost:${TEST_PORT}/health" >/dev/null 2>&1; then
            success "✓ Health check passed"
            return 0
        fi
        
        log "Health check attempt $attempt/$max_attempts failed, retrying..."
        sleep 3
        ((attempt++))
    done
    
    error "✗ Health check failed after $max_attempts attempts"
    error "Container logs:"
    docker logs "$TEST_CONTAINER_NAME" 2>&1 | tee -a "$LOG_FILE"
    exit 1
}

# ============================================
# DISPLAY FUNCTIONS
# ============================================
display_header() {
    clear
    echo ""
    echo -e "${CYAN}╔════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║                                            ║${NC}"
    echo -e "${CYAN}║     🚀 FINAL DEPLOYMENT SCRIPT v2.0 🚀    ║${NC}"
    echo -e "${CYAN}║     Trae New Healthy1 Platform            ║${NC}"
    echo -e "${CYAN}║                                            ║${NC}"
    echo -e "${CYAN}╚════════════════════════════════════════════╝${NC}"
    echo ""
}

display_summary() {
    echo ""
    echo -e "${CYAN}╔════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║         DEPLOYMENT READY SUMMARY           ║${NC}"
    echo -e "${CYAN}╚════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${GREEN}✓${NC} Code validation: PASSED"
    echo -e "${GREEN}✓${NC} Docker build: PASSED"
    echo -e "${GREEN}✓${NC} Container test: PASSED"
    echo -e "${GREEN}✓${NC} Health check: PASSED"
    echo -e "${GREEN}✓${NC} All systems: GO"
    echo ""
    echo -e "${BLUE}Platform:${NC} Coolify"
    echo -e "${BLUE}Domain:${NC} super.doctorhealthy1.com"
    echo -e "${BLUE}Port:${NC} $APP_PORT (internal)"
    echo -e "${BLUE}SSL:${NC} Auto-configured"
    echo ""
}

display_deployment_options() {
    echo ""
    echo -e "${CYAN}╔════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║           DEPLOYMENT OPTIONS               ║${NC}"
    echo -e "${CYAN}╚════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${YELLOW}Option 1: Coolify Dashboard (Recommended)${NC}"
    echo "  1. Login to https://api.doctorhealthy1.com"
    echo "  2. Navigate to 'new doctorhealthy1' project"
    echo "  3. Create/update application"
    echo "  4. Set domain: super.doctorhealthy1.com"
    echo "  5. Set port: $APP_PORT"
    echo "  6. Add environment variables:"
    echo "     - NODE_ENV=$NODE_ENV"
    echo "     - PORT=$PORT"
    echo "     - HOST=$HOST"
    echo "  7. Click 'Deploy'"
    echo ""
    echo -e "${YELLOW}Option 2: Manual Docker Deploy${NC}"
    echo "  docker build -t trae-healthy1 ."
    echo "  docker run -d -p $APP_PORT:$APP_PORT trae-healthy1"
    echo ""
}

display_success() {
    echo ""
    echo -e "${CYAN}╔════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║              🎉 SUCCESS! 🎉                ║${NC}"
    echo -e "${CYAN}╚════════════════════════════════════════════╝${NC}"
    echo ""
    success "All pre-deployment checks passed!"
    success "Your platform is ready to deploy!"
    echo ""
    echo -e "${GREEN}Next:${NC} Choose a deployment option above"
    echo -e "${GREEN}Docs:${NC} See FINAL-DEPLOYMENT-GUIDE.md"
    echo -e "${GREEN}Logs:${NC} $LOG_FILE"
    echo -e "${GREEN}Help:${NC} Check troubleshooting guides"
    echo ""
    echo -e "${PURPLE}Your months of work are about to pay off!${NC}"
    echo -e "${PURPLE}Deploy with confidence! 🚀${NC}"
    echo ""
}

# ============================================
# MAIN EXECUTION
# ============================================
main() {
    # Set up cleanup trap
    trap cleanup EXIT
    
    # Initialize log file
    echo "Deployment log started at $(date)" > "$LOG_FILE"
    
    display_header
    
    # Step 1: Pre-deployment checks
    step "1/6: Running pre-deployment checks..."
    validate_prerequisites
    validate_files
    echo ""
    
    # Step 2: Validate code
    step "2/6: Validating code..."
    validate_code
    echo ""
    
    # Step 3: Test Docker build
    step "3/6: Testing Docker build..."
    test_docker_build
    echo ""
    
    # Step 4: Test container
    step "4/6: Testing container..."
    test_container
    echo ""
    
    # Step 5: Deployment summary
    step "5/6: Deployment summary..."
    display_summary
    
    # Step 6: Deployment instructions
    step "6/6: Next steps..."
    display_deployment_options
    
    display_success
}

# Execute main function
main "$@"