#!/bin/bash

################################################################################
# SECURITY SCAN
# Security vulnerability scanning
################################################################################

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

log() { echo -e "${GREEN}[$(date +'%H:%M:%S')]${NC} $1"; }
log_error() { echo -e "${RED}[$(date +'%H:%M:%S')] ERROR:${NC} $1"; }
log_warning() { echo -e "${YELLOW}[$(date +'%H:%M:%S')] WARNING:${NC} $1"; }
log_success() { echo -e "${GREEN}[$(date +'%H:%M:%S')] ✓${NC} $1"; }

log "═══════════════════════════════════════════════════════════════"
log "SECURITY SCAN - STARTING"
log "═══════════════════════════════════════════════════════════════"

cd "$PROJECT_ROOT"

# Go security scan
log "Scanning Go dependencies..."
if command -v gosec &> /dev/null; then
    gosec ./... || log_warning "Security issues found"
else
    log_warning "gosec not installed, skipping Go security scan"
    log "Install with: go install github.com/securego/gosec/v2/cmd/gosec@latest"
fi

# Check for common vulnerabilities
log "Checking for common vulnerabilities..."

# Check for hardcoded secrets
log "Checking for hardcoded secrets..."
if grep -r "password.*=.*\"" . --exclude-dir={.git,node_modules,vendor} | grep -v ".sh" | grep -v "example"; then
    log_warning "Potential hardcoded passwords found"
else
    log_success "No hardcoded passwords found"
fi

# Check for SQL injection vulnerabilities
log "Checking for SQL injection vulnerabilities..."
if grep -r "fmt.Sprintf.*SELECT" . --exclude-dir={.git,node_modules,vendor}; then
    log_warning "Potential SQL injection vulnerabilities found"
else
    log_success "No SQL injection vulnerabilities found"
fi

# Check file permissions
log "Checking file permissions..."
find . -type f -perm -o+w ! -path "./.git/*" ! -path "./node_modules/*" | while read file; do
    log_warning "World-writable file: $file"
done

log_success "Security scan completed"
