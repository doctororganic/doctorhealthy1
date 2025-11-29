# Testing Guide

This guide covers all aspects of testing for the Nutrition Platform backend, including unit tests, integration tests, security tests, and performance tests.

## Table of Contents

1. [Test Structure](#test-structure)
2. [Running Tests](#running-tests)
3. [Test Types](#test-types)
4. [Writing Tests](#writing-tests)
5. [Test Coverage](#test-coverage)
6. [CI/CD Integration](#cicd-integration)
7. [Best Practices](#best-practices)

## Test Structure

```
backend/tests/
├── performance/
│   └── load_test.go          # Performance and load testing
├── security/
│   └── security_test.go       # Security vulnerability testing
├── integration/
│   └── handlers_test.go      # Integration tests for handlers
├── contract/
│   └── api_contract_test.go   # API contract tests
├── e2e/
│   └── phase1_e2e_test.go    # End-to-end tests
└── README.md                  # This guide
```

### Test Categories

#### 1. Performance Tests (`performance/`)
- **Load Testing**: Test system under high load
- **Stress Testing**: Test system limits
- **Concurrency Testing**: Test concurrent access patterns
- **Memory Testing**: Monitor memory consumption
- **Cache Testing**: Verify cache performance

#### 2. Security Tests (`security/`)
- **Input Validation**: Test for malicious inputs
- **Authentication**: Test auth security measures
- **Rate Limiting**: Verify rate limiting works
- **CORS Testing**: Ensure proper CORS configuration
- **Security Headers**: Verify security headers are present
- **Attack Prevention**: SQL injection, XSS, path traversal

#### 3. Integration Tests (`integration/`)
- **Handler Integration**: Test handler with real dependencies
- **Database Integration**: Test with actual database
- **API Integration**: Test complete API workflows

#### 4. Contract Tests (`contract/`)
- **API Contracts**: Verify API response formats
- **Schema Validation**: Ensure response schemas match
- **Backward Compatibility**: Prevent breaking changes

#### 5. End-to-End Tests (`e2e/`)
- **User Workflows**: Test complete user journeys
- **Multi-Service**: Test interactions between services
- **Real Environment**: Test in production-like environment

## Running Tests

### Prerequisites

Ensure the test environment is set up:

```bash
# Install test dependencies
go mod download

# Set up test database
make test-db-setup

# Start required services (Redis, etc.)
make start-test-services
```

### Running All Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run tests with verbose output
go test ./tests/... -v

# Run tests with race detection
go test ./tests/... -race
```

### Running Specific Test Categories

```bash
# Performance tests
go test ./tests/performance/... -v

# Security tests
go test ./tests/security/... -v

# Integration tests
go test ./tests/integration/... -v

# Contract tests
go test ./tests/contract/... -v

# E2E tests (requires full environment)
go test ./tests/e2e/... -v -tags=e2e
```

### Running Individual Tests

```bash
# Specific test function
go test ./tests/performance -run TestAPILoad -v

# Benchmark tests
go test ./tests/performance -bench=. -v

# Tests with specific timeout
go test ./tests/... -timeout=30m
```

### Test with Different Environments

```bash
# Development environment
ENV=dev go test ./tests/...

# Staging environment
ENV=staging go test ./tests/...

# Production environment (read-only tests only)
ENV=prod go test ./tests/... -tags=readonly
```

## Test Types

### Performance Tests

#### Load Testing
```bash
# Run load tests with different configurations
go test ./tests/performance -run TestAPILoad -v

# Results include:
# - Total requests
# - Success/failure rates
# - Response times
# - Requests per second
# - Error rates
```

#### Memory Testing
```bash
# Test memory consumption under load
go test ./tests/performance -run TestMemoryUsage -v

# Monitors:
# - Initial memory usage
# - Memory increase during test
# - Memory per request
# - Memory leaks
```

#### Cache Testing
```bash
# Test cache performance and hit rates
go test ./tests/performance -run TestCachePerformance -v

# Validates:
# - Cache hit rates
# - Cache response times
# - Cache invalidation
```

### Security Tests

#### Input Validation
```bash
# Test malicious input handling
go test ./tests/security -run TestInputValidation -v

# Tests for:
# - SQL injection attempts
# - XSS attacks
# - Large payloads
# - Invalid formats
```

#### Authentication Security
```bash
# Test authentication security measures
go test ./tests/security -run TestAuthenticationSecurity -v

# Validates:
# - JWT token validation
# - Authorization header handling
# - Token expiration
```

#### Rate Limiting
```bash
# Test rate limiting functionality
go test ./tests/security -run TestRateLimiting -v

# Ensures:
# - Rate limits are enforced
# - Rate limit headers are present
# - 429 responses are returned
```

### Integration Tests

#### API Integration
```bash
# Test complete API workflows
go test ./tests/integration -run TestNutritionDataWorkflow -v

# Tests:
# - Complete request/response cycles
# - Database interactions
# - Error handling
```

## Writing Tests

### Test Structure Template

```go
package security

import (
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
)

// TestSuite groups related tests
type SecurityTestSuite struct {
    suite.Suite
    client    *http.Client
    baseURL    string
}

// SetupSuite runs before all tests in suite
func (suite *SecurityTestSuite) SetupSuite() {
    suite.client = &http.Client{Timeout: 10 * time.Second}
    suite.baseURL = "http://localhost:8080"
}

// TearDownSuite runs after all tests in suite
func (suite *SecurityTestSuite) TearDownSuite() {
    // Cleanup resources
}

// SetupTest runs before each test
func (suite *SecurityTestSuite) SetupTest() {
    // Prepare for each test
}

// TearDownTest runs after each test
func (suite *SecurityTestSuite) TearDownTest() {
    // Cleanup after each test
}

// TestExample shows test structure
func (suite *SecurityTestSuite) TestExample() {
    // Arrange
    expected := "expected value"
    input := "test input"
    
    // Act
    result := someFunction(input)
    
    // Assert
    assert.Equal(suite.T(), expected, result)
    require.NotNil(suite.T(), result)
}

// Run the test suite
func TestSecurityTestSuite(t *testing.T) {
    suite.Run(t, new(SecurityTestSuite))
}
```

### Best Practices for Writing Tests

#### 1. Test Naming
```go
// Good: Descriptive and follows conventions
func TestInputValidation_RejectsSQLInjection(t *testing.T) {}

func TestRateLimiting_Returns429AfterLimit(t *testing.T) {}

// Bad: Vague or non-descriptive
func TestStuff(t *testing.T) {}
func TestSecurity(t *testing.T) {}
```

#### 2. Test Structure (AAA Pattern)
```go
func TestUserCreation_ValidData_ReturnsSuccess(t *testing.T) {
    // Arrange
    validUser := User{
        Email:    "test@example.com",
        Password: "securepassword123",
        Name:     "Test User",
    }
    
    // Act
    result, err := CreateUser(validUser)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, validUser.Email, result.Email)
    assert.NotEmpty(t, result.ID)
}
```

#### 3. Table-Driven Tests
```go
func TestEmailValidation(t *testing.T) {
    tests := []struct {
        name        string
        email       string
        expected    bool
        description string
    }{
        {
            name:        "Valid email",
            email:       "user@example.com",
            expected:    true,
            description: "Should accept valid email format",
        },
        {
            name:        "Invalid email no domain",
            email:       "user@",
            expected:    false,
            description: "Should reject email without domain",
        },
        {
            name:        "Invalid email no @",
            email:       "userexample.com",
            expected:    false,
            description: "Should reject email without @ symbol",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ValidateEmail(tt.email)
            assert.Equal(t, tt.expected, result, tt.description)
        })
    }
}
```

#### 4. Mock Testing
```go
func TestUserService_CreateUser_RepositoryError(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{
        CreateFunc: func(user *User) error {
            return errors.New("database error")
        },
    }
    
    service := NewUserService(mockRepo)
    user := &User{Email: "test@example.com"}
    
    // Act
    err := service.CreateUser(user)
    
    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "database error")
}
```

## Test Coverage

### Generating Coverage Reports

```bash
# Generate coverage report
make test-coverage

# Generate HTML coverage report
go test ./tests/... -coverprofile=coverage.out -v
go tool cover -html=coverage.out -o coverage.html

# Coverage by package
go test ./tests/... -coverprofile=coverage.out
go tool cover -func=coverage.out

# Coverage threshold check
go test ./tests/... -coverprofile=coverage.out
coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
if [ $(echo "$coverage < 80" | bc -l) -eq 1 ]; then
    echo "Coverage below 80%: $coverage%"
    exit 1
fi
```

### Coverage Requirements

- **Overall Coverage**: Minimum 80%
- **Critical Paths**: 95% coverage required
- **Security Functions**: 100% coverage required
- **Error Handling**: 100% coverage required

### Coverage Exclusions

```go
//go:build !integration
// +build !integration

// Exclude from integration tests
func testHelper() {}
```

## CI/CD Integration

### GitHub Actions Workflow

```yaml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run tests
        run: |
          make test-coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```

### Local Pre-commit Testing

```bash
#!/bin/sh
# .husky/pre-commit

# Run tests
make test

# Check coverage
make test-coverage

# Run security tests
go test ./tests/security/... -v

echo "All tests passed!"
```

## Best Practices

### 1. Test Organization

- **Group related tests** in test suites
- **Use descriptive names** that explain what's being tested
- **Follow AAA pattern**: Arrange, Act, Assert
- **Keep tests independent** and isolated

### 2. Test Data Management

```go
// Use test fixtures
func setupTestData(t *testing.T) *TestDatabase {
    db := setupTestDB(t)
    
    // Insert test data
    _, err := db.Exec(`INSERT INTO users (email, name) VALUES (?, ?)`, 
        "test@example.com", "Test User")
    require.NoError(t, err)
    
    return db
}

// Clean up after tests
func cleanupTestData(t *testing.T, db *TestDatabase) {
    db.Exec("DELETE FROM users WHERE email LIKE 'test%'")
}
```

### 3. Test Environment

```go
// Use environment-specific configs
func getTestConfig() Config {
    return Config{
        DatabaseURL: "sqlite://test.db",
        RedisURL:   "redis://localhost:6379/1",
        LogLevel:    "debug",
    }
}
```

### 4. Error Testing

```go
func TestErrorCases(t *testing.T) {
    tests := []struct {
        name     string
        input     interface{}
        expected  string
    }{
        {"Empty input", "", "input cannot be empty"},
        {"Nil input", nil, "input is required"},
        {"Invalid type", 123, "input must be string"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := SomeFunction(tt.input)
            
            assert.Error(t, err)
            assert.Contains(t, err.Error(), tt.expected)
        })
    }
}
```

### 5. Performance Testing Guidelines

```go
// Use subtests for different scenarios
func TestPerformance(t *testing.T) {
    scenarios := []struct {
        name     string
        concurrency int
        duration time.Duration
    }{
        {"Low load", 10, 10 * time.Second},
        {"Medium load", 50, 30 * time.Second},
        {"High load", 100, 60 * time.Second},
    }
    
    for _, scenario := range scenarios {
        t.Run(scenario.name, func(t *testing.T) {
            result := runLoadTest(scenario.concurrency, scenario.duration)
            
            assert.Less(t, result.AvgResponseTime, 500*time.Millisecond)
            assert.Greater(t, result.RequestsPerSec, 100.0)
            assert.Less(t, result.ErrorRate, 1.0)
        })
    }
}
```

### 6. Security Testing Guidelines

```go
// Test common attack vectors
func TestSecurityAttacks(t *testing.T) {
    attacks := []struct {
        type     string
        payload  string
        expected int
    }{
        {"SQL Injection", "'; DROP TABLE users; --", 400},
        {"XSS", "<script>alert('xss')</script>", 400},
        {"Path Traversal", "../../../etc/passwd", 404},
    }
    
    for _, attack := range attacks {
        t.Run(attack.type, func(t *testing.T) {
            resp := makeRequest(attack.payload)
            assert.NotEqual(t, 200, resp.StatusCode)
        })
    }
}
```

## Troubleshooting

### Common Test Issues

#### 1. Race Conditions
```bash
# Detect race conditions
go test ./tests/... -race

# Common causes:
# - Shared state between tests
# - Concurrent access to resources
# - Global variables
```

#### 2. Time-based Test Failures
```go
// Use fake time for deterministic tests
func TestTimeBasedLogic(t *testing.T) {
    fakeTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
    
    // Use fake time in your logic
    result := processWithTime(someInput, fakeTime)
    
    assert.Equal(t, expected, result)
}
```

#### 3. External Dependencies
```go
// Mock external dependencies
type MockHTTPClient struct {
    DoFunc func(*http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
    return m.DoFunc(req)
}
```

#### 4. Database Test Issues
```go
// Use transactions for test isolation
func setupTestDB(t *testing.T) *sql.DB {
    db := openTestDB()
    tx, err := db.Begin()
    require.NoError(t, err)
    
    // Use transaction for test
    t.Cleanup(func() {
        tx.Rollback()
    })
    
    return db
}
```

## Test Execution Tips

### Running Tests Efficiently

```bash
# Run only changed packages
go test ./tests/... -short

# Run tests in parallel
go test ./tests/... -parallel 4

# Run failed tests only
go test ./tests/... -run TestFailed

# Run tests with specific tags
go test ./tests/... -tags=integration -v
```

### Debugging Tests

```bash
# Enable test debugging
go test ./tests/... -v -run TestSpecific

# Use test outputs
go test ./tests/... -v 2>&1 | grep -E "(PASS|FAIL|panic)"

# Generate test profile
go test ./tests/... -cpuprofile=cpu.prof -memprofile=mem.prof
```

This comprehensive testing guide should help you write effective, maintainable tests for the Nutrition Platform backend.
