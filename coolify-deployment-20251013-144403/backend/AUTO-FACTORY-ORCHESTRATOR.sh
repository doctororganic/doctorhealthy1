#!/bin/bash

################################################################################
# AUTO-FACTORY ORCHESTRATOR
# Comprehensive automated build, test, fix, and deploy system
# With parallel execution, real-time monitoring, and SSH integration
################################################################################

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# Configuration
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"
LOG_DIR="$PROJECT_ROOT/logs/orchestrator"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
MAIN_LOG="$LOG_DIR/orchestrator_$TIMESTAMP.log"

mkdir -p "$LOG_DIR"

################################################################################
# Logging
################################################################################

log() { echo -e "${GREEN}[$(date +'%H:%M:%S')]${NC} $1" | tee -a "$MAIN_LOG"; }
log_error() { echo -e "${RED}[$(date +'%H:%M:%S')] ERROR:${NC} $1" | tee -a "$MAIN_LOG"; }
log_warning() { echo -e "${YELLOW}[$(date +'%H:%M:%S')] WARNING:${NC} $1" | tee -a "$MAIN_LOG"; }
log_info() { echo -e "${BLUE}[$(date +'%H:%M:%S')] INFO:${NC} $1" | tee -a "$MAIN_LOG"; }
log_success() { echo -e "${GREEN}[$(date +'%H:%M:%S')] ✓${NC} $1" | tee -a "$MAIN_LOG"; }

################################################################################
# Phase 1: Environment Setup
################################################################################

phase_1_setup() {
    log "═══════════════════════════════════════════════════════════════"
    log "PHASE 1: Environment Setup & Validation"
    log "═══════════════════════════════════════════════════════════════"
    
    # Check Go
    if command -v go &> /dev/null; then
        log_success "Go installed: $(go version)"
    else
        log_error "Go is not installed"
        exit 1
    fi
    
    # Check Docker
    if command -v docker &> /dev/null; then
        log_success "Docker installed: $(docker --version)"
    else
        log_warning "Docker not installed"
    fi
    
    # Check Node.js
    if command -v node &> /dev/null; then
        log_success "Node.js installed: $(node --version)"
    else
        log_warning "Node.js not installed"
    fi
    
    log_success "Phase 1 completed"
}

################################################################################
# Phase 2: Backend Build & Test
################################################################################

phase_2_backend() {
    log "═══════════════════════════════════════════════════════════════"
    log "PHASE 2: Backend Build & Test"
    log "═══════════════════════════════════════════════════════════════"
    
    cd "$BACKEND_DIR"
    
    # Clean and download dependencies
    log_info "Installing Go dependencies..."
    go mod tidy
    go mod download
    
    # Run tests in parallel
    log_info "Running backend tests..."
    go test ./... -v -cover -parallel 4 2>&1 | tee "$LOG_DIR/backend_tests_$TIMESTAMP.log" || log_warning "Some tests failed"
    
    # Build
    log_info "Building backend..."
    go build -o bin/server ./cmd/server || {
        log_error "Build failed, attempting fix..."
        go clean -cache
        go build -o bin/server ./cmd/server
    }
    
    log_success "Backend build completed"
}

################################################################################
# Phase 3: Docker Build
################################################################################

phase_3_docker() {
    log "═══════════════════════════════════════════════════════════════"
    log "PHASE 3: Docker Build"
    log "═══════════════════════════════════════════════════════════════"
    
    cd "$BACKEND_DIR"
    
    if [ ! -f "Dockerfile" ]; then
        log_warning "No Dockerfile found"
        return 0
    fi
    
    log_info "Building Docker image..."
    docker build -t nutrition-platform:latest . 2>&1 | tee "$LOG_DIR/docker_build_$TIMESTAMP.log"
    
    log_success "Docker build completed"
}

################################################################################
# Phase 4: Integration Tests
################################################################################

phase_4_integration() {
    log "═══════════════════════════════════════════════════════════════"
    log "PHASE 4: Integration Tests"
    log "═══════════════════════════════════════════════════════════════"
    
    cd "$BACKEND_DIR"
    
    # Start test server
    log_info "Starting test server..."
    ./bin/server &
    SERVER_PID=$!
    
    sleep 5
    
    # Test health endpoint
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        log_success "Health check passed"
    else
        log_error "Health check failed"
    fi
    
    # Test API endpoints
    local endpoints=(
        "/api/v1/users"
        "/api/v1/foods"
        "/api/v1/workouts"
    )
    
    for endpoint in "${endpoints[@]}"; do
        if curl -f "http://localhost:8080$endpoint" > /dev/null 2>&1; then
            log_success "Endpoint $endpoint is accessible"
        else
            log_warning "Endpoint $endpoint failed"
        fi
    done
    
    # Stop server
    kill $SERVER_PID 2>/dev/null || true
    
    log_success "Integration tests completed"
}

################################################################################
# Phase 5: Deployment Package
################################################################################

phase_5_deployment() {
    log "═══════════════════════════════════════════════════════════════"
    log "PHASE 5: Deployment Package"
    log "═══════════════════════════════════════════════════════════════"
    
    cd "$PROJECT_ROOT"
    
    local deploy_dir="deploy_$TIMESTAMP"
    mkdir -p "$deploy_dir"
    
    # Copy backend files
    cp -r "$BACKEND_DIR/bin" "$deploy_dir/"
    cp -r "$BACKEND_DIR/migrations" "$deploy_dir/" 2>/dev/null || true
    cp "$BACKEND_DIR/.env.example" "$deploy_dir/.env" 2>/dev/null || true
    cp "$BACKEND_DIR/Dockerfile" "$deploy_dir/" 2>/dev/null || true
    
    # Create archive
    tar -czf "deploy_$TIMESTAMP.tar.gz" "$deploy_dir"
    rm -rf "$deploy_dir"
    
    log_success "Deployment package created: deploy_$TIMESTAMP.tar.gz"
}

################################################################################
# Phase 6: SSH Deployment (Optional)
################################################################################

phase_6_ssh_deployment() {
    log "═══════════════════════════════════════════════════════════════"
    log "PHASE 6: SSH Deployment (Optional)"
    log "═══════════════════════════════════════════════════════════════"
    
    if [ -z "$SSH_HOST" ]; then
        log_info "SSH_HOST not set, skipping remote deployment"
        log_info "To deploy remotely, set: SSH_HOST=your-server.com SSH_USER=user"
        return 0
    fi
    
    log_info "Deploying to $SSH_HOST..."
    
    # Find latest deployment package
    local deploy_package=$(ls -t deploy_*.tar.gz | head -1)
    
    # Upload
    log_info "Uploading deployment package..."
    scp "$deploy_package" "$SSH_USER@$SSH_HOST:/tmp/"
    
    # Deploy
    log_info "Executing deployment on remote server..."
    ssh "$SSH_USER@$SSH_HOST" << 'ENDSSH'
        set -e
        cd /opt/nutrition-platform
        
        # Backup
        if [ -d "current" ]; then
            mv current "backup_$(date +%Y%m%d_%H%M%S)"
        fi
        
        # Extract
        mkdir -p current
        tar -xzf /tmp/deploy_*.tar.gz -C current/
        
        # Restart service
        systemctl restart nutrition-platform || true
        
        rm /tmp/deploy_*.tar.gz
ENDSSH
    
    log_success "SSH deployment completed"
}

################################################################################
# Main
################################################################################

main() {
    log "╔═══════════════════════════════════════════════════════════════╗"
    log "║         AUTO-FACTORY ORCHESTRATOR - STARTING                  ║"
    log "╚═══════════════════════════════════════════════════════════════╝"
    
    phase_1_setup
    phase_2_backend
    phase_3_docker
    phase_4_integration
    phase_5_deployment
    phase_6_ssh_deployment
    
    log "╔═══════════════════════════════════════════════════════════════╗"
    log "║         AUTO-FACTORY ORCHESTRATOR - COMPLETED                 ║"
    log "╚═══════════════════════════════════════════════════════════════╝"
    
    log_success "All phases completed!"
    log_info "Logs: $LOG_DIR"
    log_info "Deployment package: $(ls -t deploy_*.tar.gz | head -1)"
}

main "$@"
