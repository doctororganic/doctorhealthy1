# Error Handling and Logging Best Practices Implementation

This document outlines the comprehensive error handling and logging system implemented in the Nutrition Platform backend, following industry best practices.

## Table of Contents

1. [Error Handling Strategy](#error-handling-strategy)
2. [Custom Error Types](#custom-error-types)
3. [Structured Logging](#structured-logging)
4. [Global Error Handler](#global-error-handler)
5. [Correlation ID Tracking](#correlation-id-tracking)
6. [Input Validation](#input-validation)
7. [Circuit Breaker Pattern](#circuit-breaker-pattern)
8. [Retry Mechanisms](#retry-mechanisms)
9. [Security Logging](#security-logging)
10. [Performance Monitoring](#performance-monitoring)
11. [Configuration](#configuration)
12. [Testing](#testing)

## Error Handling Strategy

### Consistent Error Handling Approach

The application uses a centralized error handling strategy with the following components:

- **Custom Error Types**: Structured error types with error codes and context
- **Global Error Handler**: Centralized middleware for processing all errors
- **Structured Responses**: Consistent JSON error response format
- **Error Recovery**: Circuit breakers and retry mechanisms

### Error Hierarchy

```
APIError (base error type)
├── ValidationError
├── SecurityError
├── DatabaseError
├── TimeoutError
└── InternalServerError
```

## Custom Error Types

### APIError Structure

```go
type APIError struct {
    Code       ErrorCode   `json:"code"`
    Message    string      `json:"message"`
    Details    string      `json:"details,omitempty"`
    Timestamp  time.Time   `json:"timestamp"`
    RequestID  string      `json:"request_id,omitempty"`
    Path       string      `json:"path,omitempty"`
    Method     string      `json:"method,omitempty"`
}
```

### Error Codes

- `INVALID_API_KEY` - Authentication errors
- `RATE_LIMIT_EXCEEDED` - Rate limiting
- `INVALID_INPUT` - Validation errors
- `RESOURCE_NOT_FOUND` - Resource access errors
- `DATABASE_CONNECTION` - Database connectivity issues
- `INTERNAL_SERVER_ERROR` - Unexpected errors

### Error Response Format

```json
{
  "error": {
    "code": "INVALID_API_KEY",
    "message": "Invalid API key provided",
    "details": "Key format is invalid",
    "timestamp": "2023-05-20T14:30:45Z",
    "request_id": "abc-123-def",
    "path": "/api/v1/nutrition/analyze",
    "method": "POST"
  },
  "success": false,
  "timestamp": "2023-05-20T14:30:45Z"
}
```

## Structured Logging

### JSON Log Format

All logs are output in structured JSON format for easy parsing and analysis:

```json
{
  "timestamp": "2023-05-20T14:30:45Z",
  "level": "ERROR",
  "service": "nutrition-platform",
  "environment": "development",
  "message": "Database connection failed",
  "fields": {
    "correlation_id": "abc-123-def",
    "user_id": "user123",
    "error": "connection timeout",
    "method": "POST",
    "path": "/api/v1/nutrition/analyze",
    "status_code": 500,
    "latency_ms": 2500
  }
}
```

### Log Levels

- **DEBUG**: Detailed information for troubleshooting
- **INFO**: General information about application flow
- **WARN**: Potentially harmful situations
- **ERROR**: Error conditions that don't stop execution
- **FATAL**: Severe errors that cause application termination

### Context Enrichment

Logs include rich context information:

- **Correlation ID**: For request tracing across services
- **User ID**: For user-specific operations
- **Request Info**: Method, path, status code, latency
- **API Key**: Masked for security
- **Error Details**: Stack traces and error context

## Global Error Handler

### Middleware Integration

The global error handler is implemented as Echo middleware that:

1. **Captures All Errors**: Intercepts any error returned by handlers
2. **Converts Error Types**: Standardizes different error types to APIError
3. **Records Metrics**: Updates Prometheus metrics for monitoring
4. **Logs Errors**: Structured logging with full context
5. **Sends Notifications**: For critical errors (when configured)
6. **Sanitizes Output**: Removes sensitive information in production

### Error Processing Flow

```
Request → Handler → Error Returned → Global Handler
                                      ↓
1. Convert to APIError
2. Add request context
3. Record metrics
4. Log error
5. Send notification (if critical)
6. Sanitize for production
7. Return JSON response
```

## Correlation ID Tracking

### Automatic ID Generation

- **Incoming Requests**: Checks for `X-Correlation-ID` header
- **ID Generation**: Creates unique ID if not provided
- **Response Header**: Sets correlation ID in response
- **Context Storage**: Stores ID in Echo context for logging

### Cross-Service Tracing

```go
// Example usage in handlers
logger := Logger.WithContext(c)
logger.Info("Processing request", map[string]interface{}{
    "operation": "analyze_nutrition",
})
```

## Input Validation

### Validation Middleware

- **Early Validation**: Input validated before processing
- **Structured Errors**: Specific validation error messages
- **Security**: Prevents malformed data from causing issues

### Supported Validation Rules

```go
type NutritionRequest struct {
    Food     string  `json:"food" validate:"required"`
    Quantity float64 `json:"quantity" validate:"required,min=0"`
    Unit     string  `json:"unit" validate:"required"`
}
```

## Circuit Breaker Pattern

### Configuration

```go
errorHandler := errors.NewErrorHandler(&errors.ErrorHandlerConfig{
    CircuitBreakerEnabled: true,
    MaxRetries:           3,
    RetryDelay:           time.Second * 2,
})
```

### Circuit States

- **Closed**: Normal operation (failures < threshold)
- **Open**: Service unavailable (failures > threshold)
- **Half-Open**: Testing recovery (limited requests allowed)

### Automatic Recovery

- Monitors service health
- Automatically transitions between states
- Logs state changes for monitoring

## Retry Mechanisms

### WithRetry Function

```go
err := errorHandler.WithRetry(ctx, "database_operation", func() error {
    return performDatabaseOperation()
})
```

### Retry Conditions

- **Transient Errors**: Database timeouts, network issues
- **Configurable Limits**: Max attempts and delay between retries
- **Context Cancellation**: Respects request context timeouts

## Security Logging

### Sensitive Data Masking

```go
// API keys are masked
"api_key_prefix": "abcd1234****"

// User agents with tokens are sanitized
"user_agent": "Masked User Agent"
```

### Security Events

- **Authentication Failures**: Invalid API keys, expired tokens
- **Rate Limit Violations**: Suspicious activity detection
- **Access Violations**: Unauthorized resource access

## Performance Monitoring

### Metrics Collection

Prometheus metrics for:

- **Error Rates**: By type, severity, and component
- **Recovery Time**: Time to handle errors
- **Circuit Breaker State**: Service availability status
- **Retry Attempts**: Success/failure rates

### Performance Logging

- **Request Latency**: Response time tracking
- **Slow Operations**: Threshold-based alerting
- **Resource Usage**: Memory and CPU monitoring

## Configuration

### Environment Variables

```bash
# Logging
LOG_LEVEL=info                    # debug, info, warn, error
LOG_FORMAT=json                   # json or text
LOG_FILE=/var/log/app.log         # Optional file output
LOG_CALLER=true                   # Include caller info

# Error Handling
ENVIRONMENT=production            # development or production
ERROR_SANITIZE=true              # Remove details in production

# Circuit Breaker
CIRCUIT_BREAKER_ENABLED=true
MAX_RETRIES=3
RETRY_DELAY=2s
HEALTH_CHECK_INTERVAL=5m
ALERT_THRESHOLD=10
```

### Development vs Production

**Development:**
- Detailed error messages
- Stack traces in logs
- Debug logging enabled
- Full error context

**Production:**
- Sanitized error messages
- Masked sensitive data
- Appropriate log levels
- Error aggregation

## Testing

### Error Scenario Testing

```go
func TestErrorHandling(t *testing.T) {
    // Test invalid input validation
    req := NutritionRequest{Food: "", Quantity: -1}
    err := validateRequest(req)

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "INVALID_INPUT")
}
```

### Logging Verification

```go
func TestStructuredLogging(t *testing.T) {
    // Capture log output
    var buf bytes.Buffer
    logger := NewStructuredLogger(&LogConfig{
        EnableJSON: true,
        Writer:     &buf,
    })

    logger.Error("Test error", errors.New("test error"))

    var logEntry LogEntry
    err := json.Unmarshal(buf.Bytes(), &logEntry)

    assert.NoError(t, err)
    assert.Equal(t, "ERROR", logEntry.Level)
    assert.Equal(t, "Test error", logEntry.Message)
}
```

## Best Practices Implemented

✅ **Consistent Error Strategy**: Centralized error handling across the application
✅ **Custom Error Types**: Hierarchical error types with context
✅ **Clear Error Messages**: User-friendly messages with actionable guidance
✅ **Comprehensive Logging**: Structured JSON logs with rich context
✅ **Context Enrichment**: Correlation IDs, user info, request details
✅ **Sensitive Data Handling**: Masking and sanitization
✅ **Centralized Logging**: Single logging system with consistent format
✅ **Performance Considerations**: Async logging and sampling
✅ **Error Recovery**: Retry mechanisms and circuit breakers
✅ **User Experience**: Appropriate error responses and guidance
✅ **API Standards**: HTTP status codes and structured responses
✅ **Input Validation**: Early validation with specific error messages
✅ **Global Handlers**: Centralized exception handling
✅ **Monitoring Integration**: Metrics and alerting
✅ **Documentation**: Comprehensive error and troubleshooting guides
✅ **Testing**: Error scenario and logging verification
✅ **Environment-Specific**: Different behaviors for dev/prod

## Usage Examples

### Logging in Handlers

```go
func analyzeNutrition(c echo.Context) error {
    logger := Logger.WithContext(c)

    logger.Info("Starting nutrition analysis", map[string]interface{}{
        "food": req.Food,
        "quantity": req.Quantity,
    })

    // ... processing logic ...

    if err != nil {
        logger.Error("Analysis failed", err, map[string]interface{}{
            "food": req.Food,
        })
        return errors.ErrInternalServerError("Analysis failed")
    }

    logger.Info("Analysis completed successfully")
    return c.JSON(http.StatusOK, response)
}
```

### Error Handling

```go
func validateAPIKey(c echo.Context) error {
    apiKey := c.Request().Header.Get("X-API-Key")

    if apiKey == "" {
        return errors.ErrMissingAPIKeyError()
    }

    if !isValidAPIKey(apiKey) {
        Logger.WithContext(c).Warn("Invalid API key attempt", map[string]interface{}{
            "api_key_prefix": maskAPIKey(apiKey),
        })
        return errors.ErrInvalidAPIKeyError("API key format is invalid")
    }

    return nil
}
```

This implementation provides a robust, production-ready error handling and logging system that follows industry best practices and ensures reliable operation of the Nutrition Platform.