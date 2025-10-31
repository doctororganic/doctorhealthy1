# Error Handling and Security Implementation

This document describes the comprehensive error handling and security measures implemented in the Nutrition Platform backend.

## Overview

The platform implements a multi-layered security approach with robust error handling, input validation, rate limiting, API key management, and comprehensive logging.

## Components

### 1. Error Handling System

#### Core Components
- **`errors/errors.go`**: Defines structured error types and codes
- **`middleware/error_handler.go`**: Centralized error handling middleware
- **`validation/input_validator.go`**: Input validation and security checks

#### Error Types
```go
// Authentication errors
ErrInvalidAPIKey     // Invalid API key provided
ErrExpiredAPIKey     // API key has expired
ErrRevokedAPIKey     // API key has been revoked
ErrMissingAPIKey     // API key is required
ErrInsufficientScope // Insufficient permissions

// Rate limiting errors
ErrRateLimitExceeded // Rate limit exceeded
ErrQuotaExceeded     // Quota exceeded

// Validation errors
ErrInvalidInput      // Invalid input provided
ErrMissingParameter  // Missing required parameter
ErrInvalidFormat     // Invalid format
ErrInvalidRange      // Invalid range

// Resource errors
ErrResourceNotFound  // Resource not found
ErrResourceExists    // Resource already exists
ErrResourceLocked    // Resource is locked

// Security errors
ErrSecurityViolation // Security violation detected
ErrSuspiciousActivity // Suspicious activity detected
ErrIPBlocked         // IP address is blocked
```

#### Error Response Format
```json
{
  "error": {
    "code": "INVALID_API_KEY",
    "message": "Invalid API key provided",
    "details": "API key format is invalid",
    "timestamp": "2024-01-15T10:30:00Z",
    "request_id": "req_123456789",
    "path": "/api/v1/meals",
    "method": "GET"
  },
  "success": false,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### 2. Security Logger

#### Features
- **Structured Logging**: JSON-formatted security events
- **Real-time Monitoring**: Immediate detection of security threats
- **Alert System**: Configurable alerts for security violations
- **Metrics Tracking**: Security metrics and statistics
- **IP Blocking**: Automatic blocking of suspicious IPs

#### Security Events
```go
type SecurityEvent struct {
    Type        string            // Event type (authentication_failure, rate_limit_exceeded, etc.)
    Severity    string            // Severity level (low, medium, high, critical)
    Message     string            // Human-readable message
    Details     map[string]string // Additional details
    Timestamp   time.Time         // Event timestamp
    IPAddress   string            // Source IP address
    UserAgent   string            // User agent string
    APIKeyID    string            // API key ID (if applicable)
    Endpoint    string            // Requested endpoint
    Method      string            // HTTP method
    StatusCode  int               // Response status code
    ResponseTime time.Duration    // Response time
}
```

#### Alert Rules
- **High Authentication Failures**: 5+ failures in 5 minutes
- **Critical Security Violation**: Immediate blocking
- **Rate Limit Abuse**: 10+ violations in 1 minute

### 3. Input Validation

#### Security Checks
- **SQL Injection Detection**: Pattern-based detection of SQL injection attempts
- **XSS Prevention**: Cross-site scripting attack detection
- **Command Injection**: Prevention of command injection attacks
- **Path Traversal**: Detection of directory traversal attempts
- **Suspicious Patterns**: Detection of unusual character patterns

#### Validation Rules
```go
type ValidationRule struct {
    Field           string
    Required        bool
    MinLength       int
    MaxLength       int
    Pattern         string
    AllowedValues   []string
    CustomValidator func(value string) error
}
```

#### Built-in Validators
- Email format validation
- URL format validation
- Phone number validation
- Date format validation
- Integer/Float range validation
- UUID format validation

### 4. Security Configuration

#### Configuration Categories
- **API Key Management**: Key generation, expiration, rotation
- **Rate Limiting**: Global, per-API-key, and per-IP limits
- **Authentication**: JWT settings, password requirements
- **Encryption**: Algorithms, key sizes, TLS configuration
- **Logging**: Log levels, file rotation, retention
- **CORS**: Cross-origin resource sharing settings
- **Security Headers**: HTTP security headers configuration
- **Input Validation**: Validation rules and thresholds
- **Monitoring**: Metrics, alerts, and notifications

#### Environment Variables
```bash
# API Key Configuration
API_KEY_ENABLED=true
API_KEY_LENGTH=32
API_KEY_PREFIX=np_
API_KEY_DEFAULT_EXPIRATION=8760h  # 1 year

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_GLOBAL_REQUESTS=10000
RATE_LIMIT_GLOBAL_WINDOW=1h

# Security
SECURITY_HEADERS_ENABLED=true
VALIDATION_ENABLED=true
VALIDATION_BLOCK_SQL_INJECTION=true
VALIDATION_BLOCK_XSS=true

# Logging
LOG_LEVEL=info
SECURITY_LOG_PATH=./logs/security.log
LOG_STACK_TRACE=true

# Environment
ENVIRONMENT=production
DEBUG_MODE=false
SANITIZE_ERRORS=true
```

## Implementation Details

### 1. Error Handler Middleware

The error handler middleware provides:
- **Centralized Error Processing**: All errors are processed through a single handler
- **Context-Aware Logging**: Errors include request context (IP, user agent, endpoint)
- **Security Event Generation**: Security-related errors generate security events
- **Error Sanitization**: Sensitive information is removed in production
- **Panic Recovery**: Graceful handling of application panics

### 2. Security Middleware Stack

```go
// Security middleware order (important!)
e.Use(middlewareError.PanicRecovery(errorConfig))     // Panic recovery
e.Use(middlewareError.RequestValidator())             // Request validation
e.Use(middleware.Logger())                            // Request logging
e.Use(middleware.CORS())                              // CORS handling
e.Use(middleware.Gzip())                              // Response compression
e.Use(middleware.RateLimiter(...))                    // Rate limiting
e.Use(middlewareCustom.RequestID())                   // Request ID generation
e.Use(middlewareCustom.SecurityHeaders())             // Security headers
```

### 3. API Key Authentication

```go
// Optional API key authentication
e.Use(middlewareCustom.OptionalAPIKeyAuth())

// Required API key with specific scope
e.Use(middlewareCustom.RequireAPIKeyScope("read"))
e.Use(middlewareCustom.RequireAPIKeyScope("write"))
e.Use(middlewareCustom.RequireAPIKeyScope("admin"))
```

### 4. Input Validation Usage

```go
// Validate individual input
validationContext := &validation.ValidationContext{
    UserRole:  "user",
    IPAddress: c.RealIP(),
    UserAgent: c.Request().UserAgent(),
    Endpoint:  c.Request().URL.Path,
    Method:    c.Request().Method,
}

if err := inputValidator.ValidateInput(userInput, validationContext); err != nil {
    return err // Returns structured API error
}

// Validate multiple fields with rules
rules := []validation.ValidationRule{
    {
        Field:     "email",
        Required:  true,
        MaxLength: 255,
        CustomValidator: validation.ValidateEmail,
    },
    {
        Field:         "age",
        Required:      true,
        CustomValidator: func(value string) error {
            return validation.ValidateInteger(value, 1, 120)
        },
    },
}

validationErrors := inputValidator.ValidateFields(requestData, rules, validationContext)
if validationErrors.HasErrors() {
    return validationErrors // Returns validation error response
}
```

## Security Features

### 1. Automatic Threat Detection
- **SQL Injection**: Pattern-based detection with immediate blocking
- **XSS Attacks**: Script injection detection and prevention
- **Command Injection**: System command detection
- **Path Traversal**: Directory traversal attempt detection
- **Brute Force**: Failed authentication attempt tracking
- **Rate Limit Abuse**: Excessive request detection

### 2. IP-based Security
- **Automatic IP Blocking**: IPs with multiple violations are automatically blocked
- **Whitelist/Blacklist**: Configurable IP allow/deny lists
- **Geolocation Blocking**: Optional geographic restrictions
- **Suspicious Activity Tracking**: Behavioral analysis and scoring

### 3. API Key Security
- **Secure Generation**: Cryptographically secure random generation
- **Scope-based Access**: Fine-grained permission control
- **Automatic Expiration**: Configurable key lifetimes
- **Usage Tracking**: Detailed usage statistics and monitoring
- **Rate Limiting**: Per-key rate limits and quotas

### 4. Security Headers
```http
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Referrer-Policy: strict-origin-when-cross-origin
Content-Security-Policy: default-src 'self'
Strict-Transport-Security: max-age=31536000; includeSubDomains
Permissions-Policy: geolocation=(), microphone=(), camera=()
Cache-Control: no-cache, no-store, must-revalidate
```

## Monitoring and Alerting

### 1. Security Metrics
- Total security events
- Events by type and severity
- Failed authentication attempts
- Rate limit violations
- Blocked IPs and suspicious activity
- Response times and error rates

### 2. Alert Thresholds
- **Failed Authentication**: 10+ attempts in 5 minutes
- **Rate Limit Violations**: 50+ violations in 5 minutes
- **Security Violations**: 5+ violations in 5 minutes
- **Error Rate**: >5% error rate in 5 minutes
- **Response Time**: >5 seconds average in 5 minutes

### 3. Notification Channels
- **Email Alerts**: Critical security events
- **Webhook Notifications**: Real-time event streaming
- **Log Aggregation**: Centralized log collection
- **Metrics Dashboard**: Real-time monitoring

## Best Practices

### 1. Error Handling
- Always use structured errors with appropriate codes
- Include sufficient context for debugging
- Sanitize errors in production environments
- Log all errors with appropriate severity levels
- Provide user-friendly error messages

### 2. Security
- Validate all input at the application boundary
- Use parameterized queries to prevent SQL injection
- Implement proper authentication and authorization
- Apply the principle of least privilege
- Regularly rotate API keys and secrets
- Monitor and alert on security events

### 3. Performance
- Use efficient validation algorithms
- Implement proper caching for validation rules
- Monitor performance impact of security measures
- Optimize logging for high-throughput scenarios

### 4. Maintenance
- Regularly update security patterns and rules
- Review and update alert thresholds
- Perform security audits and penetration testing
- Keep dependencies updated
- Monitor security advisories

## Testing

The implementation includes comprehensive tests:
- **Unit Tests**: Individual component testing
- **Integration Tests**: End-to-end security testing
- **Security Tests**: Penetration testing and vulnerability assessment
- **Performance Tests**: Load testing with security measures enabled
- **Compliance Tests**: Verification against security standards

## Compliance

The implementation addresses requirements from:
- **OWASP API Security Top 10**
- **NIST Cybersecurity Framework**
- **ISO 27001 Information Security Management**
- **PCI DSS** (where applicable)
- **GDPR** (data protection aspects)

## Conclusion

This comprehensive error handling and security implementation provides:
- **Robust Protection**: Multi-layered security against common threats
- **Comprehensive Monitoring**: Real-time threat detection and alerting
- **Flexible Configuration**: Environment-specific security settings
- **Developer-Friendly**: Easy-to-use APIs and clear documentation
- **Production-Ready**: Scalable and performant security measures

The system is designed to be both secure and maintainable, providing strong protection while remaining easy to operate and extend.