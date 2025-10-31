#!/bin/bash

# Advanced Security Scanner for Nutrition Platform
# Comprehensive security analysis and vulnerability detection

set -e

echo "🛡️  NUTRITION PLATFORM - COMPREHENSIVE SECURITY SCAN"
echo "=================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Scan configuration
SCAN_DIR="."
EXCLUDE_DIRS=".git,node_modules,dist,build,.next"
LOG_FILE="/tmp/security-scan-$(date +%Y%m%d_%H%M%S).log"

# Initialize counters
TOTAL_ISSUES=0
CRITICAL_ISSUES=0
HIGH_ISSUES=0
MEDIUM_ISSUES=0
LOW_ISSUES=0

# Logging function
log() {
    echo -e "${2:-$GREEN}[$(date '+%Y-%m-%d %H:%M:%S')] $1${NC}" | tee -a "$LOG_FILE"
}

warn() {
    echo -e "${YELLOW}[$(date '+%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}" | tee -a "$LOG_FILE"
    ((TOTAL_ISSUES++))
    ((MEDIUM_ISSUES++))
}

error() {
    echo -e "${RED}[$(date '+%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}" | tee -a "$LOG_FILE"
    ((TOTAL_ISSUES++))
    ((HIGH_ISSUES++))
}

critical() {
    echo -e "${RED}[$(date '+%Y-%m-%d %H:%M:%S')] CRITICAL: $1${NC}" | tee -a "$LOG_FILE"
    ((TOTAL_ISSUES++))
    ((CRITICAL_ISSUES++))
}

info() {
    echo -e "${BLUE}[$(date '+%Y-%m-%d %H:%M:%S')] INFO: $1${NC}" | tee -a "$LOG_FILE"
}

# 1. SECRET DETECTION SCAN
scan_secrets() {
    log "🔍 Scanning for exposed secrets..." "$BLUE"

    # Patterns for different types of secrets
    local secret_patterns=(
        "password.*=.*['\"][^'\"]*['\"]"
        "api_key.*=.*['\"][^'\"]*['\"]"
        "secret.*=.*['\"][^'\"]*['\"]"
        "token.*=.*['\"][^'\"]*['\"]"
        "key.*=.*['\"][^'\"]*['\"]"
        "auth.*=.*['\"][^'\"]*['\"]"
        "bearer.*['\"][^'\"]*['\"]"
        "private_key.*=.*['\"][^'\"]*['\"]"
    )

    local found_secrets=()
    for pattern in "${secret_patterns[@]}"; do
        while IFS= read -r -d '' file; do
            found_secrets+=("$file")
        done < <(find "$SCAN_DIR" -type f \( -name "*.go" -o -name "*.js" -o -name "*.py" -o -name "*.json" -o -name "*.env*" -o -name "*.yml" -o -name "*.yaml" \) ! -path "*/node_modules/*" ! -path "*/.git/*" -exec grep -l "$pattern" {} \; 2>/dev/null)
    done

    if [ ${#found_secrets[@]} -gt 0 ]; then
        error "Found potential secrets in ${#found_secrets[@]} files:"
        printf '%s\n' "${found_secrets[@]}" | head -10 | while read -r file; do
            error "  - $file"
        done
        if [ ${#found_secrets[@]} -gt 10 ]; then
            error "  ... and $(( ${#found_secrets[@]} - 10 )) more files"
        fi
    else
        log "✅ No hardcoded secrets found" "$GREEN"
    fi
}

# 2. DEPENDENCY VULNERABILITY SCAN
scan_dependencies() {
    log "📦 Scanning dependencies for vulnerabilities..." "$BLUE"

    # Check if safety is installed
    if command -v safety &> /dev/null; then
        if [ -f "requirements.txt" ]; then
            info "Running Python dependency security scan..."
            if safety check --output text 2>/dev/null; then
                log "✅ Python dependencies are secure" "$GREEN"
            else
                warn "Python dependencies have known vulnerabilities"
            fi
        fi
    else
        warn "Safety tool not installed. Install with: pip install safety"
    fi

    # Check Go modules
    if [ -f "go.mod" ]; then
        info "Checking Go module dependencies..."
        if command -v govulncheck &> /dev/null; then
            if govulncheck ./... 2>/dev/null; then
                log "✅ Go dependencies are secure" "$GREEN"
            else
                warn "Go dependencies have known vulnerabilities"
            fi
        else
            warn "govulncheck not installed. Install with: go install golang.org/x/vuln/cmd/govulncheck@latest"
        fi
    fi
}

# 3. PERMISSION AND OWNERSHIP SCAN
scan_permissions() {
    log "🔒 Scanning file permissions..." "$BLUE"

    # Check for world-writable files
    while IFS= read -r -d '' file; do
        error "World-writable file found: $file"
    done < <(find "$SCAN_DIR" -type f ! -path "*/node_modules/*" ! -path "*/.git/*" -perm /002 2>/dev/null)

    # Check for files with no permissions
    while IFS= read -r -d '' file; do
        warn "File with no permissions: $file"
    done < <(find "$SCAN_DIR" -type f ! -path "*/node_modules/*" ! -path "*/.git/*" -perm /000 2>/dev/null)

    if ! find "$SCAN_DIR" -type f ! -path "*/node_modules/*" ! -path "*/.git/*" -perm /002 2>/dev/null | grep -q .; then
        log "✅ No world-writable files found" "$GREEN"
    fi
}

# 4. DOCKER SECURITY SCAN
scan_docker() {
    log "🐳 Scanning Docker configuration..." "$BLUE"

    if [ -f "Dockerfile" ]; then
        # Check for running as root
        if grep -q "USER root" Dockerfile; then
            error "Dockerfile runs as root user"
        fi

        # Check for latest tag usage
        if grep -q "FROM.*:latest" Dockerfile; then
            warn "Using 'latest' tag in Dockerfile - consider using specific versions"
        fi

        # Check for apt update without upgrade
        if grep -q "apt.*update" Dockerfile && ! grep -A5 -B5 "apt.*update" Dockerfile | grep -q "apt.*upgrade"; then
            warn "apt update without upgrade may leave system vulnerable"
        fi

        log "✅ Dockerfile security scan completed" "$GREEN"
    else
        warn "No Dockerfile found"
    fi
}

# 5. CONFIGURATION SECURITY SCAN
scan_configurations() {
    log "⚙️  Scanning configuration files..." "$BLUE"

    # Check for debug modes in production
    while IFS= read -r -d '' file; do
        if grep -q -i "debug.*true\|production.*false\|development.*true" "$file" 2>/dev/null; then
            warn "Potential debug/production configuration issue in: $(basename "$file")"
        fi
    done < <(find "$SCAN_DIR" -name "*.json" -o -name "*.yaml" -o -name "*.yml" -o -name "*.env*" ! -path "*/node_modules/*" ! -path "*/.git/*")

    # Check for exposed ports
    if [ -f "docker-compose.yml" ]; then
        if grep -q "ports:" docker-compose.yml && grep -q "0.0.0.0:" docker-compose.yml; then
            warn "Services may be exposing ports to all interfaces (0.0.0.0)"
        fi
    fi
}

# 6. CODE SECURITY SCAN
scan_code_security() {
    log "🔍 Scanning code for security issues..." "$BLUE"

    # Check for SQL injection patterns
    while IFS= read -r -d '' file; do
        if grep -q -i "SELECT.*+.*\|DROP.*\|DELETE.*\|UPDATE.*SET.*\|INSERT.*INTO" "$file" 2>/dev/null; then
            warn "Potential SQL injection pattern in: $(basename "$file")"
        fi
    done < <(find "$SCAN_DIR" -name "*.go" -o -name "*.py" -o -name "*.js" ! -path "*/node_modules/*" ! -path "*/.git/*")

    # Check for XSS patterns
    while IFS= read -r -d '' file; do
        if grep -q -i "innerHTML\|outerHTML\|document\.write" "$file" 2>/dev/null; then
            warn "Potential XSS vulnerability in: $(basename "$file")"
        fi
    done < <(find "$SCAN_DIR" -name "*.js" -o -name "*.ts" ! -path "*/node_modules/*" ! -path "*/.git/*")

    # Check for insecure random number generation
    while IFS= read -r -d '' file; do
        if grep -q "math\.random\|rand\.rand" "$file" 2>/dev/null; then
            warn "Insecure random number generation in: $(basename "$file")"
        fi
    done < <(find "$SCAN_DIR" -name "*.go" -o -name "*.js" -o -name "*.py" ! -path "*/node_modules/*" ! -path "*/.git/*")
}

# 7. NETWORK SECURITY SCAN
scan_network() {
    log "🌐 Scanning network configuration..." "$BLUE"

    # Check for hardcoded IPs
    while IFS= read -r -d '' file; do
        if grep -q -E '\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b' "$file" 2>/dev/null; then
            if ! grep -q "127\.0\.0\.1\|0\.0\.0\.0\|localhost" "$file"; then
                warn "Hardcoded IP address found in: $(basename "$file")"
            fi
        fi
    done < <(find "$SCAN_DIR" -type f ! -path "*/node_modules/*" ! -path "*/.git/*")
}

# 8. TLS/SSL CONFIGURATION SCAN
scan_tls() {
    log "🔒 Scanning TLS/SSL configuration..." "$BLUE"

    if [ -f "nginx.conf" ]; then
        # Check for SSL configuration
        if ! grep -q "ssl_certificate\|ssl_certificate_key" nginx.conf; then
            warn "No SSL certificate configuration found in nginx.conf"
        fi

        # Check for secure SSL protocols
        if grep -q "ssl_protocols.*SSLv" nginx.conf; then
            warn "Insecure SSL protocols detected in nginx.conf"
        fi

        # Check for HSTS
        if ! grep -q "add_header.*Strict-Transport-Security" nginx.conf; then
            warn "HSTS header not configured in nginx.conf"
        fi
    fi
}

# 9. BACKUP AND RECOVERY SCAN
scan_backup_recovery() {
    log "💾 Scanning backup and recovery configuration..." "$BLUE"

    # Check for backup scripts
    if ! find "$SCAN_DIR" -name "*backup*" -type f ! -path "*/node_modules/*" ! -path "*/.git/*" | grep -q .; then
        warn "No backup scripts found"
    else
        log "✅ Backup scripts found" "$GREEN"
    fi

    # Check database backup configuration
    if [ -f "docker-compose.yml" ]; then
        if ! grep -q "volumes" docker-compose.yml; then
            warn "No persistent volumes configured - data may be lost on restart"
        fi
    fi
}

# 10. COMPLIANCE SCAN
scan_compliance() {
    log "📋 Scanning for compliance issues..." "$BLUE"

    # Check for required security headers
    if [ -f "main.go" ]; then
        local security_headers_found=false
        if grep -q "X-Content-Type-Options\|X-Frame-Options\|X-XSS-Protection\|Strict-Transport-Security" main.go; then
            security_headers_found=true
        fi

        if ! $security_headers_found; then
            warn "Security headers may not be properly configured"
        else
            log "✅ Security headers configured" "$GREEN"
        fi
    fi
}

# MAIN EXECUTION
main() {
    info "Starting comprehensive security scan of nutrition platform..."
    echo "Log file: $LOG_FILE"
    echo ""

    # Run all scans
    scan_secrets
    echo ""

    scan_dependencies
    echo ""

    scan_permissions
    echo ""

    scan_docker
    echo ""

    scan_configurations
    echo ""

    scan_code_security
    echo ""

    scan_network
    echo ""

    scan_tls
    echo ""

    scan_backup_recovery
    echo ""

    scan_compliance
    echo ""

    # SUMMARY
    echo ""
    echo "📊 SECURITY SCAN SUMMARY"
    echo "======================="
    echo "Total Issues Found: $TOTAL_ISSUES"
    echo "Critical Issues: $CRITICAL_ISSUES"
    echo "High Issues: $HIGH_ISSUES"
    echo "Medium Issues: $MEDIUM_ISSUES"
    echo "Low Issues: $LOW_ISSUES"
    echo ""

    # Recommendations
    if [ $CRITICAL_ISSUES -gt 0 ]; then
        echo -e "${RED}🚨 CRITICAL ISSUES FOUND - Fix immediately before deployment!${NC}"
        echo "1. Review all hardcoded secrets and move to environment variables"
        echo "2. Update all dependencies to latest secure versions"
        echo "3. Fix file permissions"
        echo "4. Implement proper SSL/TLS configuration"
    elif [ $HIGH_ISSUES -gt 0 ]; then
        echo -e "${YELLOW}⚠️  HIGH PRIORITY ISSUES - Address before production deployment${NC}"
        echo "1. Fix security headers"
        echo "2. Update vulnerable dependencies"
        echo "3. Review and fix code security issues"
    else
        echo -e "${GREEN}✅ SCAN COMPLETED - No critical security issues found${NC}"
        echo "Regular security scans recommended for ongoing protection."
    fi

    echo ""
    info "Security scan completed. Review the log file for detailed information: $LOG_FILE"

    # Exit with error code if critical issues found
    if [ $CRITICAL_ISSUES -gt 0 ]; then
        exit 1
    fi
}

# Handle script interruption
trap 'echo -e "\n${YELLOW}Scan interrupted by user${NC}"; exit 130' INT

# Execute main function
main "$@"
