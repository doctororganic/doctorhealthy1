#!/bin/bash

# Comprehensive Test Runner for Nutrition Platform Backend
# This script runs all test suites with proper configuration and reporting

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TEST_RESULTS_DIR="$PROJECT_ROOT/test-results"
COVERAGE_FILE="$TEST_RESULTS_DIR/coverage.out"
COVERAGE_HTML="$TEST_RESULTS_DIR/coverage.html"
TEST_REPORT="$TEST_RESULTS_DIR/test-report.json"

# Create test results directory
mkdir -p "$TEST_RESULTS_DIR"

# Logging function
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Function to check if service is running
check_service() {
    local service=$1
    local port=$2
    
    if nc -z localhost "$port" 2>/dev/null; then
        log_success "$service is running on port $port"
        return 0
    else
        log_warning "$service is not running on port $port"
        return 1
    fi
}

# Function to start services if needed
start_services() {
    log "Starting required services..."
    
    # Start Redis if not running
    if ! check_service "Redis" 6379; then
        log "Starting Redis..."
        if command -v docker >/dev/null 2>&1; then
            docker run -d --name nutrition-redis-test -p 6379:6379 redis:6-alpine
            sleep 3
        else
            log_error "Docker not found. Please start Redis manually on port 6379"
            exit 1
        fi
    fi
    
    # Start application server for E2E tests
    if ! check_service "Application" 8080; then
        log "Starting application server..."
        cd "$PROJECT_ROOT"
        go run main.go > "$TEST_RESULTS_DIR/server.log" 2>&1 &
        SERVER_PID=$!
        echo $SERVER_PID > "$TEST_RESULTS_DIR/server.pid"
        
        # Wait for server to be ready
        for i in {1..30}; do
            if check_service "Application" 8080; then
                log_success "Application server started successfully"
                break
            fi
            sleep 1
        done
        
        if ! check_service "Application" 8080; then
            log_error "Failed to start application server"
            exit 1
        fi
    fi
}

# Function to stop services
stop_services() {
    log "Stopping services..."
    
    # Stop application server
    if [ -f "$TEST_RESULTS_DIR/server.pid" ]; then
        SERVER_PID=$(cat "$TEST_RESULTS_DIR/server.pid")
        if kill -0 "$SERVER_PID" 2>/dev/null; then
            kill "$SERVER_PID"
            log_success "Application server stopped"
        fi
        rm -f "$TEST_RESULTS_DIR/server.pid"
    fi
    
    # Stop Redis Docker container
    if command -v docker >/dev/null 2>&1; then
        if docker ps -q -f name=nutrition-redis-test | grep -q .; then
            docker stop nutrition-redis-test >/dev/null
            docker rm nutrition-redis-test >/dev/null
            log_success "Redis container stopped"
        fi
    fi
}

# Function to run unit tests
run_unit_tests() {
    log "Running unit tests..."
    
    cd "$PROJECT_ROOT"
    
    # Run unit tests with coverage
    if go test ./middleware ./cache ./utils -v -race -cover -coverprofile="$COVERAGE_FILE" > "$TEST_RESULTS_DIR/unit-tests.log" 2>&1; then
        log_success "Unit tests passed"
        return 0
    else
        log_error "Unit tests failed"
        cat "$TEST_RESULTS_DIR/unit-tests.log"
        return 1
    fi
}

# Function to run integration tests
run_integration_tests() {
    log "Running integration tests..."
    
    cd "$PROJECT_ROOT"
    
    if go test ./tests/integration -v > "$TEST_RESULTS_DIR/integration-tests.log" 2>&1; then
        log_success "Integration tests passed"
        return 0
    else
        log_error "Integration tests failed"
        cat "$TEST_RESULTS_DIR/integration-tests.log"
        return 1
    fi
}

# Function to run E2E tests
run_e2e_tests() {
    log "Running E2E tests..."
    
    cd "$PROJECT_ROOT"
    
    if go test ./tests/e2e -v -timeout=5m > "$TEST_RESULTS_DIR/e2e-tests.log" 2>&1; then
        log_success "E2E tests passed"
        return 0
    else
        log_error "E2E tests failed"
        cat "$TEST_RESULTS_DIR/e2e-tests.log"
        return 1
    fi
}

# Function to run performance tests
run_performance_tests() {
    log "Running performance tests..."
    
    cd "$PROJECT_ROOT"
    
    if go test ./tests/performance -v -timeout=10m > "$TEST_RESULTS_DIR/performance-tests.log" 2>&1; then
        log_success "Performance tests passed"
        return 0
    else
        log_error "Performance tests failed"
        cat "$TEST_RESULTS_DIR/performance-tests.log"
        return 1
    fi
}

# Function to run security tests
run_security_tests() {
    log "Running security tests..."
    
    cd "$PROJECT_ROOT"
    
    if go test ./tests/security -v > "$TEST_RESULTS_DIR/security-tests.log" 2>&1; then
        log_success "Security tests passed"
        return 0
    else
        log_error "Security tests failed"
        cat "$TEST_RESULTS_DIR/security-tests.log"
        return 1
    fi
}

# Function to run database tests
run_database_tests() {
    log "Running database tests..."
    
    cd "$PROJECT_ROOT"
    
    if go test ./tests/database -v > "$TEST_RESULTS_DIR/database-tests.log" 2>&1; then
        log_success "Database tests passed"
        return 0
    else
        log_error "Database tests failed"
        cat "$TEST_RESULTS_DIR/database-tests.log"
        return 1
    fi
}

# Function to generate coverage report
generate_coverage_report() {
    log "Generating coverage report..."
    
    cd "$PROJECT_ROOT"
    
    if [ -f "$COVERAGE_FILE" ]; then
        # Generate HTML coverage report
        go tool cover -html="$COVERAGE_FILE" -o "$COVERAGE_HTML"
        
        # Generate coverage summary
        go tool cover -func="$COVERAGE_FILE" > "$TEST_RESULTS_DIR/coverage-summary.txt"
        
        log_success "Coverage report generated: $COVERAGE_HTML"
        
        # Display coverage summary
        log "Coverage Summary:"
        cat "$TEST_RESULTS_DIR/coverage-summary.txt"
    else
        log_warning "No coverage file found"
    fi
}

# Function to generate test report
generate_test_report() {
    log "Generating test report..."
    
    local report_file="$TEST_RESULTS_DIR/test-summary.md"
    
    cat > "$report_file" << EOF
# Test Execution Summary

**Date:** $(date)
**Project:** Nutrition Platform Backend

## Test Results

| Test Suite | Status | Log File |
|-------------|--------|-----------|
EOF

    # Add test results
    if [ -f "$TEST_RESULTS_DIR/unit-tests.log" ]; then
        echo "| Unit Tests | $(grep -q "PASS" "$TEST_RESULTS_DIR/unit-tests.log" && echo "✅ PASSED" || echo "❌ FAILED") | [unit-tests.log](unit-tests.log) |" >> "$report_file"
    fi
    
    if [ -f "$TEST_RESULTS_DIR/integration-tests.log" ]; then
        echo "| Integration Tests | $(grep -q "PASS" "$TEST_RESULTS_DIR/integration-tests.log" && echo "✅ PASSED" || echo "❌ FAILED") | [integration-tests.log](integration-tests.log) |" >> "$report_file"
    fi
    
    if [ -f "$TEST_RESULTS_DIR/e2e-tests.log" ]; then
        echo "| E2E Tests | $(grep -q "PASS" "$TEST_RESULTS_DIR/e2e-tests.log" && echo "✅ PASSED" || echo "❌ FAILED") | [e2e-tests.log](e2e-tests.log) |" >> "$report_file"
    fi
    
    if [ -f "$TEST_RESULTS_DIR/performance-tests.log" ]; then
        echo "| Performance Tests | $(grep -q "PASS" "$TEST_RESULTS_DIR/performance-tests.log" && echo "✅ PASSED" || echo "❌ FAILED") | [performance-tests.log](performance-tests.log) |" >> "$report_file"
    fi
    
    if [ -f "$TEST_RESULTS_DIR/security-tests.log" ]; then
        echo "| Security Tests | $(grep -q "PASS" "$TEST_RESULTS_DIR/security-tests.log" && echo "✅ PASSED" || echo "❌ FAILED") | [security-tests.log](security-tests.log) |" >> "$report_file"
    fi
    
    if [ -f "$TEST_RESULTS_DIR/database-tests.log" ]; then
        echo "| Database Tests | $(grep -q "PASS" "$TEST_RESULTS_DIR/database-tests.log" && echo "✅ PASSED" || echo "❌ FAILED") | [database-tests.log](database-tests.log) |" >> "$report_file"
    fi
    
    cat >> "$report_file" << EOF

## Coverage Report

EOF
    
    if [ -f "$TEST_RESULTS_DIR/coverage-summary.txt" ]; then
        cat "$TEST_RESULTS_DIR/coverage-summary.txt" >> "$report_file"
    fi
    
    log_success "Test report generated: $report_file"
}

# Function to cleanup test artifacts
cleanup() {
    log "Cleaning up test artifacts..."
    
    # Stop services
    stop_services
    
    # Remove temporary files
    find "$TEST_RESULTS_DIR" -name "*.tmp" -delete 2>/dev/null || true
    
    log_success "Cleanup completed"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -u, --unit        Run unit tests only"
    echo "  -i, --integration Run integration tests only"
    echo "  -e, --e2e        Run E2E tests only"
    echo "  -p, --performance Run performance tests only"
    echo "  -s, --security    Run security tests only"
    echo "  -d, --database    Run database tests only"
    echo "  -a, --all         Run all tests (default)"
    echo "  -c, --coverage    Generate coverage report"
    echo "  -r, --report      Generate test report"
    echo "  --cleanup          Cleanup test artifacts"
    echo "  --no-services      Don't start/stop services"
    echo "  -h, --help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                           # Run all tests"
    echo "  $0 --unit --coverage        # Run unit tests with coverage"
    echo "  $0 --e2e --no-services     # Run E2E tests without starting services"
    echo "  $0 --cleanup                 # Cleanup test artifacts"
}

# Main execution logic
main() {
    local run_unit=false
    local run_integration=false
    local run_e2e=false
    local run_performance=false
    local run_security=false
    local run_database=false
    local run_all=true
    local generate_coverage=false
    local generate_report=false
    local cleanup_only=false
    local no_services=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -u|--unit)
                run_unit=true
                run_all=false
                shift
                ;;
            -i|--integration)
                run_integration=true
                run_all=false
                shift
                ;;
            -e|--e2e)
                run_e2e=true
                run_all=false
                shift
                ;;
            -p|--performance)
                run_performance=true
                run_all=false
                shift
                ;;
            -s|--security)
                run_security=true
                run_all=false
                shift
                ;;
            -d|--database)
                run_database=true
                run_all=false
                shift
                ;;
            -a|--all)
                run_all=true
                shift
                ;;
            -c|--coverage)
                generate_coverage=true
                shift
                ;;
            -r|--report)
                generate_report=true
                shift
                ;;
            --cleanup)
                cleanup_only=true
                shift
                ;;
            --no-services)
                no_services=true
                shift
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    # Handle cleanup only
    if [ "$cleanup_only" = true ]; then
        cleanup
        exit 0
    fi
    
    # Set up trap for cleanup
    trap cleanup EXIT
    
    # Start services if needed
    if [ "$no_services" = false ] && { [ "$run_all" = true ] || [ "$run_e2e" = true ] || [ "$run_performance" = true ]; }; then
        start_services
    fi
    
    local test_failed=false
    
    # Run tests based on flags
    if [ "$run_all" = true ] || [ "$run_unit" = true ]; then
        if ! run_unit_tests; then
            test_failed=true
        fi
    fi
    
    if [ "$run_all" = true ] || [ "$run_integration" = true ]; then
        if ! run_integration_tests; then
            test_failed=true
        fi
    fi
    
    if [ "$run_all" = true ] || [ "$run_e2e" = true ]; then
        if ! run_e2e_tests; then
            test_failed=true
        fi
    fi
    
    if [ "$run_all" = true ] || [ "$run_performance" = true ]; then
        if ! run_performance_tests; then
            test_failed=true
        fi
    fi
    
    if [ "$run_all" = true ] || [ "$run_security" = true ]; then
        if ! run_security_tests; then
            test_failed=true
        fi
    fi
    
    if [ "$run_all" = true ] || [ "$run_database" = true ]; then
        if ! run_database_tests; then
            test_failed=true
        fi
    fi
    
    # Generate coverage report if requested
    if [ "$generate_coverage" = true ]; then
        generate_coverage_report
    fi
    
    # Generate test report if requested
    if [ "$generate_report" = true ]; then
        generate_test_report
    fi
    
    # Exit with appropriate code
    if [ "$test_failed" = true ]; then
        log_error "Some tests failed"
        exit 1
    else
        log_success "All tests passed"
        exit 0
    fi
}

# Check dependencies
check_dependencies() {
    log "Checking dependencies..."
    
    # Check Go
    if ! command -v go >/dev/null 2>&1; then
        log_error "Go is not installed"
        exit 1
    fi
    
    # Check nc for service checking
    if ! command -v nc >/dev/null 2>&1; then
        log_warning "netcat is not installed, service checking may not work"
    fi
    
    log_success "Dependencies checked"
}

# Run main function with all arguments
check_dependencies
main "$@"
