#!/bin/bash

################################################################################
# AUTO-FIX AGENT
# Automatically detects and fixes common issues
################################################################################

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"

log() { echo -e "${GREEN}[$(date +'%H:%M:%S')]${NC} $1"; }
log_error() { echo -e "${RED}[$(date +'%H:%M:%S')] ERROR:${NC} $1"; }
log_success() { echo -e "${GREEN}[$(date +'%H:%M:%S')] ✓${NC} $1"; }
log_info() { echo -e "${BLUE}[$(date +'%H:%M:%S')] INFO:${NC} $1"; }

fix_go_modules() {
    log_info "Fixing Go modules..."
    cd "$BACKEND_DIR"
    go mod tidy
    go mod download
    go clean -cache
    log_success "Go modules fixed"
}

fix_permissions() {
    log_info "Fixing file permissions..."
    chmod +x "$BACKEND_DIR"/*.sh 2>/dev/null || true
    chmod +x "$BACKEND_DIR"/bin/* 2>/dev/null || true
    log_success "Permissions fixed"
}

fix_database() {
    log_info "Fixing database..."
    cd "$BACKEND_DIR"
    rm -f nutrition_platform.db
    log_success "Database reset"
}

fix_docker() {
    log_info "Cleaning Docker..."
    docker system prune -f || true
    log_success "Docker cleaned"
}

main() {
    log "═══════════════════════════════════════════════════════════════"
    log "AUTO-FIX AGENT - STARTING"
    log "═══════════════════════════════════════════════════════════════"
    
    fix_go_modules
    fix_permissions
    fix_database
    fix_docker
    
    log_success "All fixes applied!"
}

main "$@"
