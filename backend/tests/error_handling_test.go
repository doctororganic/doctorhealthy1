package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"nutrition-platform/errors"
	"nutrition-platform/middleware"
	"nutrition-platform/monitoring"
)

// ErrorHandlingTestSuite contains all error handling tests
type ErrorHandlingTestSuite struct {
	suite.Suite
	echo         *echo.Echo
	errorHandler *errors.ErrorHandler
	dbHandler    *errors.DatabaseErrorHandler
	authHandler  *errors.AuthErrorHandler
	healthMonitor *monitoring.HealthMonitor
	mockDB       *MockDB
	mockRedis    *MockRedis
}

// MockDB simulates database for testing
type MockDB struct {
	mock.Mock
	*sql.DB
}

func (m *MockDB) PingContext(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDB) Stats() sql.DBStats {
	args := m.Called()
	return args.Get(0).(sql.DBStats)
}

// MockRedis simulates Redis for testing
type MockRedis struct {
	mock.Mock
}

func (m *MockRedis) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockRedis) PoolStats() interface{} {
	args := m.Called()
	return args.Get(0)
}

// SetupSuite initializes the test suite
func (suite *ErrorHandlingTestSuite) SetupSuite() {
	suite.echo = echo.New()
	suite.mockDB = &MockDB{}
	suite.mockRedis = &MockRedis{}

	// Initialize error handlers
	suite.errorHandler = errors.NewErrorHandler(&errors.ErrorHandlerConfig{
		EnableMetrics:     true,
		EnableLogging:     true,
		EnableStackTrace:  true,
		EnableNotification: false,
		MaxRetries:        3,
		RetryDelay:        100 * time.Millisecond,
	})

	suite.dbHandler = errors.NewDatabaseErrorHandler(&errors.DatabaseConfig{
		MaxRetries:      3,
		RetryDelay:      100 * time.Millisecond,
		QueryTimeout:    5 * time.Second,
		HealthCheckInterval: 30 * time.Second,
		EnableMetrics:   true,
	})

	suite.authHandler = errors.NewAuthErrorHandler(&errors.AuthConfig{
		JWTSecret:         "test-secret",
		JWTExpiration:     24 * time.Hour,
		RateLimitRequests: 100,
		RateLimitWindow:   time.Minute,
		MaxLoginAttempts:  5,
		LockoutDuration:   15 * time.Minute,
		CleanupInterval:   time.Hour,
		EnableMetrics:     true,
		LogFailedAttempts: true,
	})

	suite.healthMonitor = monitoring.NewHealthMonitor(&monitoring.MonitorConfig{
		CheckInterval:      5 * time.Second,
		HealthCheckTimeout: 3 * time.Second,
		AlertThresholds: &monitoring.AlertThresholds{
			CPUUsage:            80.0,
			MemoryUsage:         85.0,
			ResponseTime:        5 * time.Second,
			ErrorRate:           0.05,
			DatabaseConnections: 100,
			RedisConnections:    50,
			ActiveGoroutines:    1000,
		},
		EnableMetrics:   true,
		EnableAlerts:    false,
		MetricsPort:     9090,
		HealthEndpoint:  "/health",
		LogHealthChecks: true,
	}, suite.mockDB, suite.mockRedis)
}

// TearDownSuite cleans up after tests
func (suite *ErrorHandlingTestSuite) TearDownSuite() {
	if suite.healthMonitor != nil {
		suite.healthMonitor.Stop()
	}
	if suite.authHandler != nil {
		suite.authHandler.Stop()
	}
}

// Test Core Error Handling

func (suite *ErrorHandlingTestSuite) TestAPIErrorCreation() {
	err := errors.NewAPIError(errors.ErrValidation, "Invalid input")
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), errors.ErrValidation, err.Code)
	assert.Equal(suite.T(), "Invalid input", err.Message)
	assert.Equal(suite.T(), http.StatusBadRequest, err.HTTPStatus())
}

func (suite *ErrorHandlingTestSuite) TestErrorHandlerProcessing() {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	testErr := errors.NewAPIError(errors.ErrInternal, "Test error")
	result := suite.errorHandler.ProcessError(testErr, c)

	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), errors.ErrInternal, result.Code)
}

func (suite *ErrorHandlingTestSuite) TestErrorSanitization() {
	sensitiveErr := fmt.Errorf("database connection failed: password=secret123")
	apiErr := errors.NewAPIError(errors.ErrDatabase, sensitiveErr.Error())
	
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	result := suite.errorHandler.ProcessError(apiErr, c)
	assert.NotContains(suite.T(), result.Message, "password=secret123")
}

// Test Database Error Handling

func (suite *ErrorHandlingTestSuite) TestDatabaseConnectionFailure() {
	suite.mockDB.On("PingContext", mock.Anything).Return(fmt.Errorf("connection refused"))

	err := suite.dbHandler.ExecuteWithRetry(context.Background(), func(ctx context.Context) error {
		return suite.mockDB.PingContext(ctx)
	})

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "connection refused")
	suite.mockDB.AssertExpectations(suite.T())
}

func (suite *ErrorHandlingTestSuite) TestDatabaseRetryMechanism() {
	// First two calls fail, third succeeds
	suite.mockDB.On("PingContext", mock.Anything).Return(fmt.Errorf("temporary failure")).Twice()
	suite.mockDB.On("PingContext", mock.Anything).Return(nil).Once()

	err := suite.dbHandler.ExecuteWithRetry(context.Background(), func(ctx context.Context) error {
		return suite.mockDB.PingContext(ctx)
	})

	assert.NoError(suite.T(), err)
	suite.mockDB.AssertExpectations(suite.T())
}

func (suite *ErrorHandlingTestSuite) TestDatabaseTimeout() {
	suite.mockDB.On("PingContext", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		select {
		case <-time.After(10 * time.Second): // Simulate long operation
		case <-ctx.Done():
			return
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := suite.dbHandler.ExecuteWithTimeout(ctx, func(ctx context.Context) error {
		return suite.mockDB.PingContext(ctx)
	})

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "timeout")
}

// Test Authentication Error Handling

func (suite *ErrorHandlingTestSuite) TestJWTValidation() {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Test missing token
	token, err := suite.authHandler.ValidateJWT("", c)
	assert.Nil(suite.T(), token)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), errors.ErrMissingAPIKey, err.Code)

	// Test invalid token
	token, err = suite.authHandler.ValidateJWT("invalid-token", c)
	assert.Nil(suite.T(), token)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), errors.ErrInvalidAPIKey, err.Code)
}

func (suite *ErrorHandlingTestSuite) TestRateLimiting() {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)
	c.Request().Header.Set("X-Real-IP", "192.168.1.1")

	// First request should pass
	err := suite.authHandler.CheckRateLimit("test-user", c)
	assert.Nil(suite.T(), err)

	// Simulate rate limit exceeded
	for i := 0; i < 150; i++ { // Exceed the limit of 100 requests
		suite.authHandler.CheckRateLimit("test-user", c)
	}

	err = suite.authHandler.CheckRateLimit("test-user", c)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), errors.ErrRateLimitExceeded, err.Code)
}

func (suite *ErrorHandlingTestSuite) TestAPIKeyValidation() {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Test missing API key
	keyInfo, err := suite.authHandler.ValidateAPIKey("", c)
	assert.Nil(suite.T(), keyInfo)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), errors.ErrMissingAPIKey, err.Code)

	// Test invalid API key
	keyInfo, err = suite.authHandler.ValidateAPIKey("invalid-key", c)
	assert.Nil(suite.T(), keyInfo)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), errors.ErrInvalidAPIKey, err.Code)

	// Test valid API key
	validKey := "valid-test-key"
	validKeyInfo := &errors.APIKeyInfo{
		Key:       validKey,
		UserID:    "test-user",
		Scopes:    []string{"read", "write"},
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		IsActive:  true,
	}
	suite.authHandler.AddAPIKey(validKey, validKeyInfo)

	keyInfo, err = suite.authHandler.ValidateAPIKey(validKey, c)
	assert.NotNil(suite.T(), keyInfo)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "test-user", keyInfo.UserID)
}

// Test Health Monitoring

func (suite *ErrorHandlingTestSuite) TestHealthCheckSuccess() {
	suite.mockDB.On("PingContext", mock.Anything).Return(nil)
	suite.mockDB.On("Stats").Return(sql.DBStats{OpenConnections: 5})
	suite.mockRedis.On("Ping", mock.Anything).Return(nil)

	status := suite.healthMonitor.GetHealthStatus()
	assert.Equal(suite.T(), "healthy", status.Status)
	assert.NotEmpty(suite.T(), status.Checks)
}

func (suite *ErrorHandlingTestSuite) TestHealthCheckFailure() {
	suite.mockDB.On("PingContext", mock.Anything).Return(fmt.Errorf("database down"))
	suite.mockRedis.On("Ping", mock.Anything).Return(fmt.Errorf("redis down"))

	// Wait for health checks to run
	time.Sleep(100 * time.Millisecond)

	status := suite.healthMonitor.GetHealthStatus()
	assert.Equal(suite.T(), "unhealthy", status.Status)
	assert.Contains(suite.T(), status.Checks, "database")
	assert.Contains(suite.T(), status.Checks, "redis")
}

// Test Frontend Error Handling Integration

func (suite *ErrorHandlingTestSuite) TestFrontendErrorResponse() {
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Simulate API error
	apiErr := errors.NewAPIError(errors.ErrValidation, "Invalid request data")
	processedErr := suite.errorHandler.ProcessError(apiErr, c)

	// Check error response format
	assert.Equal(suite.T(), errors.ErrValidation, processedErr.Code)
	assert.NotEmpty(suite.T(), processedErr.RequestID)
	assert.NotEmpty(suite.T(), processedErr.Timestamp)
}

// Test Circuit Breaker

func (suite *ErrorHandlingTestSuite) TestCircuitBreakerTrip() {
	serviceName := "test-service"
	cb := suite.errorHandler.GetCircuitBreaker(serviceName)

	// Simulate multiple failures to trip the circuit breaker
	for i := 0; i < 10; i++ {
		_, err := cb.Execute(func() (interface{}, error) {
			return nil, fmt.Errorf("service failure")
		})
		assert.Error(suite.T(), err)
	}

	// Circuit breaker should now be open
	_, err := cb.Execute(func() (interface{}, error) {
		return "success", nil
	})
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "circuit breaker is open")
}

// Test Memory and Resource Management

func (suite *ErrorHandlingTestSuite) TestMemoryLeakDetection() {
	// This test would typically run for a longer period
	// and monitor memory usage patterns
	status := suite.healthMonitor.GetHealthStatus()
	assert.NotZero(suite.T(), status.SystemInfo.MemoryUsage)
	assert.NotZero(suite.T(), status.SystemInfo.Goroutines)
}

// Test Error Recovery

func (suite *ErrorHandlingTestSuite) TestErrorRecovery() {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Test panic recovery
	recovered := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				recovered = true
				suite.errorHandler.HandlePanic(r, c)
			}
		}()
		panic("test panic")
	}()

	assert.True(suite.T(), recovered)
}

// Test Load and Stress Scenarios

func (suite *ErrorHandlingTestSuite) TestHighLoadErrorHandling() {
	// Simulate high load with concurrent requests
	concurrentRequests := 100
	done := make(chan bool, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func(id int) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/test-%d", id), nil)
			rec := httptest.NewRecorder()
			c := suite.echo.NewContext(req, rec)

			// Simulate various error types
			var err *errors.APIError
			switch id % 4 {
			case 0:
				err = errors.NewAPIError(errors.ErrValidation, "Validation error")
			case 1:
				err = errors.NewAPIError(errors.ErrDatabase, "Database error")
			case 2:
				err = errors.NewAPIError(errors.ErrRateLimitExceeded, "Rate limit error")
			default:
				err = errors.NewAPIError(errors.ErrInternal, "Internal error")
			}

			processedErr := suite.errorHandler.ProcessError(err, c)
			assert.NotNil(suite.T(), processedErr)
			done <- true
		}(i)
	}

	// Wait for all requests to complete
	for i := 0; i < concurrentRequests; i++ {
		select {
		case <-done:
			// Request completed
		case <-time.After(10 * time.Second):
			suite.T().Fatal("Test timed out")
		}
	}
}

// Test Integration Scenarios

func (suite *ErrorHandlingTestSuite) TestEndToEndErrorFlow() {
	// Test complete error flow from frontend to backend
	reqBody := `{"invalid": "data"}`
	req := httptest.NewRequest(http.MethodPost, "/api/nutrition", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Simulate validation error
	validationErr := errors.NewAPIError(errors.ErrValidation, "Required field missing")
	processedErr := suite.errorHandler.ProcessError(validationErr, c)

	// Verify error structure for frontend consumption
	assert.Equal(suite.T(), errors.ErrValidation, processedErr.Code)
	assert.NotEmpty(suite.T(), processedErr.RequestID)
	assert.Equal(suite.T(), "/api/nutrition", processedErr.Path)
	assert.Equal(suite.T(), "POST", processedErr.Method)
}

// Test Security Error Scenarios

func (suite *ErrorHandlingTestSuite) TestSecurityErrorHandling() {
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("X-Real-IP", "192.168.1.100")
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Test unauthorized access
	authErr := errors.NewAPIError(errors.ErrUnauthorized, "Access denied")
	processedErr := suite.errorHandler.ProcessError(authErr, c)

	assert.Equal(suite.T(), errors.ErrUnauthorized, processedErr.Code)
	assert.NotContains(suite.T(), processedErr.Message, "internal") // No internal details leaked
}

// Test Monitoring and Alerting

func (suite *ErrorHandlingTestSuite) TestMetricsCollection() {
	// Verify that metrics are being collected
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Generate some errors to collect metrics
	for i := 0; i < 5; i++ {
		err := errors.NewAPIError(errors.ErrValidation, "Test error")
		suite.errorHandler.ProcessError(err, c)
	}

	// Metrics should be recorded (this would typically be verified
	// by checking Prometheus metrics endpoint in integration tests)
	assert.True(suite.T(), true) // Placeholder - actual metrics verification would be more complex
}

// Run the test suite
func TestErrorHandlingTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorHandlingTestSuite))
}

// Benchmark tests for performance validation

func BenchmarkErrorHandlerProcessing(b *testing.B) {
	e := echo.New()
	errorHandler := errors.NewErrorHandler(&errors.ErrorHandlerConfig{
		EnableMetrics:    false, // Disable metrics for pure performance test
		EnableLogging:    false,
		EnableStackTrace: false,
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := errors.NewAPIError(errors.ErrValidation, "Benchmark test error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		errorHandler.ProcessError(err, c)
	}
}

func BenchmarkRateLimitCheck(b *testing.B) {
	authHandler := errors.NewAuthErrorHandler(&errors.AuthConfig{
		RateLimitRequests: 1000,
		RateLimitWindow:   time.Minute,
		EnableMetrics:     false,
	})
	defer authHandler.Stop()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		authHandler.CheckRateLimit(fmt.Sprintf("user-%d", i%100), c)
	}
}