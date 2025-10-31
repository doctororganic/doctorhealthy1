#!/bin/bash

################################################################################
# PARALLEL TEST RUNNER
# Runs all tests in parallel with real-time reporting
################################################################################

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"
LOG_DIR="$PROJECT_ROOT/logs/tests"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

mkdir -p "$LOG_DIR"

log() { echo -e "${GREEN}[$(date +'%H:%M:%S')]${NC} $1"; }
log_error() { echo -e "${RED}[$(date +'%H:%M:%S')] ERROR:${NC} $1"; }
log_success() { echo -e "${GREEN}[$(date +'%H:%M:%S')] ✓${NC} $1"; }
log_info() { echo -e "${BLUE}[$(date +'%H:%M:%S')] INFO:${NC} $1"; }

################################################################################
# Test Functions
################################################################################

test_models() {
    log_info "Testing models..."
    cd "$BACKEND_DIR"
    go test ./models/... -v -cover -coverprofile="$LOG_DIR/models_coverage.out" > "$LOG_DIR/models_$TIMESTAMP.log" 2>&1
    echo $? > "$LOG_DIR/models_exit.txt"
}

test_handlers() {
    log_info "Testing handlers..."
    cd "$BACKEND_DIR"
    go test ./handlers/... -v -cover -coverprofile="$LOG_DIR/handlers_coverage.out" > "$LOG_DIR/handlers_$TIMESTAMP.log" 2>&1
    echo $? > "$LOG_DIR/handlers_exit.txt"
}

test_services() {
    log_info "Testing services..."
    cd "$BACKEND_DIR"
    go test ./services/... -v -cover -coverprofile="$LOG_DIR/services_coverage.out" > "$LOG_DIR/services_$TIMESTAMP.log" 2>&1
    echo $? > "$LOG_DIR/services_exit.txt"
}

test_middleware() {
    log_info "Testing middleware..."
    cd "$BACKEND_DIR"
    go test ./middleware/... -v -cover -coverprofile="$LOG_DIR/middleware_coverage.out" > "$LOG_DIR/middleware_$TIMESTAMP.log" 2>&1
    echo $? > "$LOG_DIR/middleware_exit.txt"
}

test_security() {
    log_info "Testing security..."
    cd "$BACKEND_DIR"
    go test ./security/... -v -cover -coverprofile="$LOG_DIR/security_coverage.out" > "$LOG_DIR/security_$TIMESTAMP.log" 2>&1
    echo $? > "$LOG_DIR/security_exit.txt"
}

test_integration() {
    log_info "Running integration tests..."
    cd "$BACKEND_DIR"
    go test ./tests/... -v -tags=integration > "$LOG_DIR/integration_$TIMESTAMP.log" 2>&1
    echo $? > "$LOG_DIR/integration_exit.txt"
}

################################################################################
# Main
################################################################################

main() {
    log "═══════════════════════════════════════════════════════════════"
    log "PARALLEL TEST RUNNER - STARTING"
    log "═══════════════════════════════════════════════════════════════"
    
    # Run tests in parallel
    test_models &
    PID_MODELS=$!
    
    test_handlers &
    PID_HANDLERS=$!
    
    test_services &
    PID_SERVICES=$!
    
    test_middleware &
    PID_MIDDLEWARE=$!
    
    test_security &
    PID_SECURITY=$!
    
    test_integration &
    PID_INTEGRATION=$!
    
    # Wait for all tests to complete
    log_info "Waiting for tests to complete..."
    wait $PID_MODELS $PID_HANDLERS $PID_SERVICES $PID_MIDDLEWARE $PID_SECURITY $PID_INTEGRATION
    
    # Check results
    log "═══════════════════════════════════════════════════════════════"
    log "TEST RESULTS"
    log "═══════════════════════════════════════════════════════════════"
    
    local failed=0
    local total=0
    local passed=0
    
    for test in models handlers services middleware security integration; do
        total=$((total + 1))
        if [ -f "$LOG_DIR/${test}_exit.txt" ]; then
            exit_code=$(cat "$LOG_DIR/${test}_exit.txt")
            if [ "$exit_code" -eq 0 ]; then
                log_success "$test tests passed"
                passed=$((passed + 1))
                
                # Show coverage if available
                if [ -f "$LOG_DIR/${test}_coverage.out" ]; then
                    coverage=$(go tool cover -func="$LOG_DIR/${test}_coverage.out" 2>/dev/null | grep total | awk '{print $3}')
                    if [ -n "$coverage" ]; then
                        log_info "  Coverage: $coverage"
                    fi
                fi
            else
                log_error "$test tests failed (exit code: $exit_code)"
                log_info "  Log: $LOG_DIR/${test}_$TIMESTAMP.log"
                failed=1
            fi
            rm "$LOG_DIR/${test}_exit.txt"
        fi
    done
    
    echo ""
    log "═══════════════════════════════════════════════════════════════"
    log "SUMMARY"
    log "═══════════════════════════════════════════════════════════════"
    log_info "Total test suites: $total"
    log_success "Passed: $passed"
    
    if [ $failed -eq 0 ]; then
        log_success "All tests passed!"
        log_info "Logs saved to: $LOG_DIR"
        exit 0
    else
        log_error "Some tests failed"
        log_info "Check logs in: $LOG_DIR"
        exit 1
    fi
}

main "$@"
