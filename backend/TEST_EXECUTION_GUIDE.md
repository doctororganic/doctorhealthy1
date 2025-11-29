# Test Execution Guide

## Overview

This guide provides comprehensive instructions for running all test suites in the nutrition platform backend project. It covers local development, CI/CD integration, and troubleshooting.

## Prerequisites

### Required Software
- **Go 1.21+**: Primary development language
- **Redis 6+**: For cache testing (optional, tests will skip if not available)
- **SQLite**: Default database (included with Go)
- **Git**: Version control

### Development Environment
```bash
# Verify Go installation
go version

# Verify Redis (optional)
redis-cli ping

# Set up workspace
cd nutrition-platform/backend
```

## Test Structure

```
backend/
├── middleware/
│   ├── cache_test.go              # Cache middleware unit tests
│   └── enhanced_rate_limiter_test.go  # Rate limiter tests
├── cache/
│   └── redis_cache_test.go        # Redis cache unit tests
├── utils/
│   └── json_loader_test.go        # JSON loader tests
├── tests/
│   ├── e2e_test.go               # Basic E2E tests
│   ├── e2e/
│   │   └── phase1_e2e_test.go    # Phase 1 E2E test suite
│   └── integration/
│       └── handlers_test.go      # Integration tests
```

## Quick Start

### Run All Tests
```bash
# From backend directory
go test ./... -v
```

### Run with Coverage
```bash
# Generate coverage report
go test ./... -cover -coverprofile=coverage.out

# View HTML coverage report
go tool cover -html=coverage.out

# View coverage in terminal
go tool cover -func=coverage.out
```

### Run Specific Test Types
```bash
# Unit tests only
go test ./middleware ./cache ./utils -v

# Integration tests only
go test ./tests/integration -v

# E2E tests only (requires running server)
go test ./tests/e2e -v
```

## Detailed Test Execution

### Unit Tests

#### Cache Middleware Tests
```bash
# Run cache middleware tests
go test ./middleware -run TestCache -v

# Run specific cache test
go test ./middleware -run TestCacheMiddleware_HitMiss -v

# Run with race detection
go test ./middleware -race -v

# Run benchmarks
go test ./middleware -bench=BenchmarkMemoryCache -v
```

#### Redis Cache Tests
```bash
# Run Redis cache tests (requires Redis)
go test ./cache -v

# Run Redis tests with race detection
go test ./cache -race -v

# Skip Redis tests if Redis not available
go test ./cache -v -tags=!redis
```

#### Rate Limiter Tests
```bash
# Run rate limiter tests
go test ./middleware -run TestRateLimiter -v

# Run concurrency tests
go test ./middleware -run TestRateLimiterConcurrency -v

# Run benchmarks
go test ./middleware -bench=BenchmarkMemoryStore -v
```

### Integration Tests

#### Handler Integration Tests
```bash
# Run all integration tests
go test ./tests/integration -v

# Run specific endpoint tests
go test ./tests/integration -run TestNutritionDataHandler_GetRecipes -v

# Run error handling tests
go test ./tests/integration -run TestNutritionDataHandler_ErrorHandling -v
```

### End-to-End Tests

#### Basic E2E Tests
```bash
# Start server first (in separate terminal)
go run main.go

# Run E2E tests
go test ./tests -run TestE2E -v

# Run specific E2E test
go test ./tests -run TestE2E_HealthCheck -v
```

#### Phase 1 E2E Tests
```bash
# Start server first (in separate terminal)
go run main.go

# Run Phase 1 E2E test suite
go test ./tests/e2e -v

# Run specific Phase 1 test
go test ./tests/e2e -run TestPhase1E2ETestSuite/TestCacheHitMiss -v

# Skip E2E tests in short mode
go test ./tests/e2e -v -short
```

## Test Configuration

### Environment Variables
```bash
# Redis configuration (for Redis tests)
export REDIS_HOST=localhost:6379
export REDIS_PASSWORD=
export REDIS_DB=0

# Test configuration
export TEST_TIMEOUT=30s
export TEST_DB_PATH=:memory:
export LOG_LEVEL=debug
```

### Test Tags
```bash
# Run tests with specific tags
go test ./... -tags=integration -v

# Skip tests with specific tags
go test ./... -tags=!slow -v
```

### Test Parallelism
```bash
# Run tests in parallel (default is number of CPUs)
go test ./... -parallel 4 -v

# Run tests sequentially
go test ./... -parallel 1 -v
```

## Performance Testing

### Benchmark Tests
```bash
# Run all benchmarks
go test ./... -bench=. -benchmem

# Run specific benchmark
go test ./middleware -bench=BenchmarkMemoryStore_Allow -benchmem

# Run benchmarks with custom count
go test ./... -bench=. -benchtime=5s -count=3
```

### Load Testing
```bash
# Install load testing tool
go install github.com/tsenart/vegeta@latest

# Run load test against API
echo "GET http://localhost:8080/api/v1/nutrition-data/recipes" | \
  vegeta attack -duration=30s -rate=50 | vegeta report

# Run load test with custom headers
echo "GET http://localhost:8080/health" | \
  vegeta attack -header="Authorization: Bearer token" -duration=30s -rate=100
```

## CI/CD Integration

### GitHub Actions
```yaml
# .github/workflows/test.yml
name: Test Suite
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      redis:
        image: redis:6
        ports:
          - 6379:6379
    
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      
      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
      
      - name: Download dependencies
        run: go mod download
      
      - name: Run unit tests
        run: go test ./middleware ./cache ./utils -v -race -cover
      
      - name: Run integration tests
        run: go test ./tests/integration -v
      
      - name: Start server for E2E tests
        run: |
          go run main.go &
          sleep 10
      
      - name: Run E2E tests
        run: go test ./tests/e2e -v
```

### Docker Testing
```bash
# Build test container
docker build -t nutrition-platform-test -f Dockerfile.test .

# Run tests in container
docker run --rm nutrition-platform-test

# Run tests with Redis
docker run --network host nutrition-platform-test
```

## Troubleshooting

### Common Issues

#### Redis Connection Errors
```bash
# Error: "Redis not available, skipping test"
# Solution: Start Redis or use mock

# Start Redis
docker run -d -p 6379:6379 redis:6

# Or use mock for testing
go test ./cache -v -tags=mock_redis
```

#### Port Already in Use
```bash
# Error: "bind: address already in use"
# Solution: Kill existing process or use different port

# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>

# Or use different port
export PORT=8081
go run main.go
```

#### Test Timeouts
```bash
# Error: "test timed out"
# Solution: Increase timeout or run tests sequentially

# Increase timeout
go test ./tests/e2e -v -timeout=60s

# Run sequentially
go test ./tests/e2e -v -parallel 1
```

#### Race Conditions
```bash
# Error: "race detected"
# Solution: Run with race detection to identify issues

# Run with race detection
go test ./... -race -v

# Run specific test with race detection
go test ./middleware -run TestRateLimiterConcurrency -race -v
```

### Debug Mode

### Verbose Output
```bash
# Maximum verbosity
go test ./... -v -test.v

# Show test execution time
go test ./... -v -test.run

# Show coverage details
go test ./... -cover -covermode=atomic -v
```

### Test Debugging
```bash
# Run single test with debugging
go test ./middleware -run TestCacheMiddleware_HitMiss -v -test.v

# Use delve debugger
dlv test ./middleware -- -test.run TestCacheMiddleware_HitMiss

# Print test output immediately
go test ./... -v -test.run
```

## Test Data Management

### Test Fixtures
```bash
# Load test fixtures
go test ./tests/integration -v -fixtures=../fixtures

# Generate test data
go test ./tests/integration -v -generate-data

# Clean test data
go test ./tests/integration -v -cleanup
```

### Database Testing
```bash
# Use in-memory database for testing
export TEST_DB=:memory:
go test ./tests/integration -v

# Use test database file
export TEST_DB=./test.db
go test ./tests/integration -v

# Reset test database
go test ./tests/integration -v -reset-db
```

## Performance Monitoring

### Test Performance Metrics
```bash
# Generate performance profile
go test ./... -cpuprofile=cpu.prof -memprofile=mem.prof -v

# Analyze CPU profile
go tool pprof cpu.prof

# Analyze memory profile
go tool pprof mem.prof
```

### Continuous Monitoring
```bash
# Run tests with performance monitoring
go test ./... -v -metrics

# Generate performance report
go test ./... -v -report=performance.html
```

## Best Practices

### Before Running Tests
1. **Clean workspace**: Remove old test artifacts
2. **Check dependencies**: Ensure all modules are up to date
3. **Start services**: Start Redis and other required services
4. **Set environment**: Configure test environment variables

### During Test Execution
1. **Monitor resources**: Watch CPU and memory usage
2. **Check logs**: Monitor application logs for errors
3. **Verify coverage**: Ensure adequate test coverage
4. **Performance**: Monitor response times and resource usage

### After Test Execution
1. **Review results**: Check for failed tests and warnings
2. **Analyze coverage**: Identify untested code paths
3. **Performance analysis**: Review performance metrics
4. **Cleanup**: Remove test artifacts and temporary data

## Test Reports

### Coverage Reports
```bash
# Generate HTML coverage report
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Generate coverage summary
go test ./... -cover -coverprofile=coverage.out
go tool cover -func=coverage.out > coverage.txt

# Generate coverage badge
go test ./... -cover -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total | awk '{print $3}'
```

### Test Reports
```bash
# Generate test report
go test ./... -v -json > test-report.json

# Generate JUnit XML report
go test ./... -v -json | go-junit-report > test-report.xml

# Generate HTML test report
gotestsum --format testname --junitfile report.xml ./...
```

## Automation Scripts

### Run All Tests Script
```bash
#!/bin/bash
# run-all-tests.sh

echo "Running unit tests..."
go test ./middleware ./cache ./utils -v -race -cover

echo "Running integration tests..."
go test ./tests/integration -v

echo "Starting server for E2E tests..."
go run main.go &
SERVER_PID=$!
sleep 10

echo "Running E2E tests..."
go test ./tests/e2e -v

echo "Cleaning up..."
kill $SERVER_PID

echo "All tests completed!"
```

### Quick Test Script
```bash
#!/bin/bash
# quick-test.sh

# Run critical tests only
go test ./middleware -run TestCache -v
go test ./cache -run TestRedisCache_GetSet -v
go test ./tests/e2e -run TestPhase1E2ETestSuite/TestCacheHitMiss -v
```

---

## Conclusion

This guide provides comprehensive instructions for executing all test suites in the nutrition platform backend. Regular test execution ensures code quality, identifies issues early, and maintains system reliability.

For additional information, refer to:
- [Testing Analysis and Recommendations](./TESTING_ANALYSIS_AND_RECOMMENDATIONS.md)
- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)