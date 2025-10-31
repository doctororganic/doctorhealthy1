#!/bin/bash
set -e

echo "🔥 LIVE TESTING & ERROR DETECTION"
echo "=================================="
echo "Starting comprehensive test suite..."
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

ERRORS=0
WARNINGS=0
PASSED=0

log_error() {
    echo -e "${RED}❌ ERROR: $1${NC}"
    ((ERRORS++))
}

log_warning() {
    echo -e "${YELLOW}⚠️  WARNING: $1${NC}"
    ((WARNINGS++))
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
    ((PASSED++))
}

log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Test 1: Check Go syntax
echo "1️⃣  Testing Go Backend Syntax..."
if cd backend && go build -o /tmp/test-build ./... 2>&1 | tee /tmp/go-build.log; then
    log_success "Go backend compiles"
    rm -f /tmp/test-build
else
    log_error "Go compilation failed"
    cat /tmp/go-build.log
fi
cd ..

# Test 2: Check Go dependencies
echo ""
echo "2️⃣  Checking Go Dependencies..."
cd backend
if go mod verify; then
    log_success "Go modules verified"
else
    log_error "Go module issues detected"
fi
if go mod tidy -v 2>&1 | grep -q "error"; then
    log_error "Go mod tidy found issues"
else
    log_success "Go dependencies clean"
fi
cd ..

# Test 3: Check Docker files
echo ""
echo "3️⃣  Validating Docker Configuration..."
for dockerfile in backend/Dockerfile backend/Dockerfile.secure frontend-nextjs/Dockerfile frontend-nextjs/Dockerfile.secure; do
    if [ -f "$dockerfile" ]; then
        if docker build -f "$dockerfile" -t test-image --no-cache . > /dev/null 2>&1; then
            log_success "Dockerfile valid: $dockerfile"
        else
            log_error "Dockerfile invalid: $dockerfile"
        fi
    else
        log_warning "Dockerfile missing: $dockerfile"
    fi
done

# Test 4: Check docker-compose
echo ""
echo "4️⃣  Validating Docker Compose..."
if docker-compose -f docker-compose.production.yml config > /dev/null 2>&1; then
    log_success "docker-compose.production.yml valid"
else
    log_error "docker-compose.production.yml has errors"
    docker-compose -f docker-compose.production.yml config
fi

# Test 5: Check environment variables
echo ""
echo "5️⃣  Checking Environment Configuration..."
if [ -f ".env.production" ]; then
    log_success ".env.production exists"
    
    # Check for placeholder values
    if grep -q "CHANGE_ME\|your_\|example\|password123" .env.production; then
        log_error "Found placeholder values in .env.production"
    else
        log_success "No placeholder values found"
    fi
    
    # Check required vars
    required_vars=("DB_PASSWORD" "JWT_SECRET" "REDIS_PASSWORD")
    for var in "${required_vars[@]}"; do
        if grep -q "^${var}=" .env.production; then
            log_success "Required var present: $var"
        else
            log_error "Missing required var: $var"
        fi
    done
else
    log_error ".env.production missing"
fi

# Test 6: Check CORS configuration
echo ""
echo "6️⃣  Validating CORS Configuration..."
if [ -f "nginx/production.conf" ]; then
    if grep -q 'Access-Control-Allow-Origin "\*"' nginx/production.conf; then
        log_error "CORS allows all origins (*) - SECURITY RISK"
    else
        log_success "CORS properly restricted"
    fi
else
    log_warning "nginx/production.conf not found"
fi

# Test 7: Check SSL configuration
echo ""
echo "7️⃣  Checking SSL/TLS Configuration..."
if grep -q "DB_SSL_MODE=disable" backend/.env.example backend/.env 2>/dev/null; then
    log_error "Database SSL disabled - SECURITY RISK"
else
    log_success "Database SSL configuration secure"
fi

# Test 8: Check for exposed secrets
echo ""
echo "8️⃣  Scanning for Exposed Secrets..."
secret_patterns=("password.*=.*['\"][^$]" "secret.*=.*['\"][^$]" "key.*=.*['\"][^$]")
found_secrets=0
for pattern in "${secret_patterns[@]}"; do
    if grep -r -i -E "$pattern" backend/ --include="*.go" --include="*.env" 2>/dev/null | grep -v "placeholder\|example\|\${"; then
        ((found_secrets++))
    fi
done
if [ $found_secrets -eq 0 ]; then
    log_success "No hardcoded secrets found"
else
    log_error "Found $found_secrets potential hardcoded secrets"
fi

# Test 9: Check port conflicts
echo ""
echo "9️⃣  Checking Port Availability..."
ports=(3000 8080 5432 6379 80 443)
for port in "${ports[@]}"; do
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
        log_warning "Port $port already in use"
    else
        log_success "Port $port available"
    fi
done

# Test 10: Check disk space
echo ""
echo "🔟 Checking System Resources..."
disk_usage=$(df . | awk 'NR==2 {print $5}' | sed 's/%//')
if [ "$disk_usage" -gt 90 ]; then
    log_error "Disk usage critical: ${disk_usage}%"
elif [ "$disk_usage" -gt 80 ]; then
    log_warning "Disk usage high: ${disk_usage}%"
else
    log_success "Disk space OK: ${disk_usage}% used"
fi

# Test 11: Check memory
available_mem=$(free -m 2>/dev/null | awk 'NR==2{print $7}' || sysctl -n hw.memsize 2>/dev/null | awk '{print $0/1024/1024}')
if [ -n "$available_mem" ]; then
    if [ "$available_mem" -lt 1000 ]; then
        log_warning "Low memory: ${available_mem}MB available"
    else
        log_success "Memory OK: ${available_mem}MB available"
    fi
fi

# Test 12: Run Go tests
echo ""
echo "1️⃣2️⃣  Running Go Unit Tests..."
cd backend
if go test ./... -v 2>&1 | tee /tmp/go-test.log; then
    log_success "Go tests passed"
else
    log_error "Go tests failed"
    tail -20 /tmp/go-test.log
fi
cd ..

# Summary
echo ""
echo "=================================="
echo "📊 TEST SUMMARY"
echo "=================================="
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${YELLOW}Warnings: $WARNINGS${NC}"
echo -e "${RED}Errors: $ERRORS${NC}"
echo ""

if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}✅ ALL CRITICAL TESTS PASSED!${NC}"
    echo "Ready for deployment"
    exit 0
else
    echo -e "${RED}❌ FOUND $ERRORS CRITICAL ERRORS${NC}"
    echo "Fix errors before deployment"
    exit 1
fi
