# Raygun Production Deployment Guide

## Environment Configuration

### Development Environment
```go
raygun, err := raygun4go.New("nutrition-platform-dev", "J5KNVQg46P71JymsDyPWiQ")
```

### Staging Environment
```go
raygun, err := raygun4go.New("nutrition-platform-staging", "J5KNVQg46P71JymsDyPWiQ")
```

### Production Environment
```go
raygun, err := raygun4go.New("nutrition-platform", "J5KNVQg46P71JymsDyPWiQ")
```

## Environment Variables
Set these in your deployment environment:

```bash
# Raygun Configuration
RAYGUN_API_KEY=J5KNVQg46P71JymsDyPWiQ
RAYGUN_APP_NAME=nutrition-platform
RAYGUN_ENVIRONMENT=production

# Application Environment
APP_ENV=production
LOG_LEVEL=info
```

## Deployment Steps

### 1. Environment-Specific Configuration
Update the Raygun initialization based on environment:

```go
func getRaygunConfig() (appName, apiKey string) {
    env := os.Getenv("APP_ENV")
    switch env {
    case "production":
        return "nutrition-platform", os.Getenv("RAYGUN_API_KEY")
    case "staging":
        return "nutrition-platform-staging", os.Getenv("RAYGUN_API_KEY")
    default:
        return "nutrition-platform-dev", os.Getenv("RAYGUN_API_KEY")
    }
}

// In main()
appName, apiKey := getRaygunConfig()
raygun, err := raygun4go.New(appName, apiKey)
```

### 2. Error Severity Levels
Tag errors appropriately:

```go
// Critical errors that require immediate attention
raygun.Tags([]string{"critical", "database", "production"})

// Warning level errors
raygun.Tags([]string{"warning", "validation", "user-input"})

// Info level for tracking
raygun.Tags([]string{"info", "performance", "metrics"})
```

### 3. User Context Enhancement
Add comprehensive user context:

```go
func enhanceRaygunContext(c echo.Context) {
    // User information
    if userID := c.Get("user_id"); userID != nil {
        raygun.User(fmt.Sprintf("%v", userID))
    }
    
    // Session information
    if sessionID := c.Get("session_id"); sessionID != nil {
        raygun.CustomData(map[string]interface{}{
            "session_id": sessionID,
            "user_role":  c.Get("user_role"),
            "tenant_id":  c.Get("tenant_id"),
        })
    }
    
    // Request context
    raygun.CustomData(map[string]interface{}{
        "endpoint":     c.Path(),
        "method":       c.Request().Method,
        "ip_address":  c.RealIP(),
        "user_agent":   c.Request().UserAgent(),
        "request_size": c.Request().ContentLength,
    })
}
```

## Monitoring and Alerts

### 1. Dashboard Setup
- Login to https://app.raygun.com
- Navigate to Application Settings
- Configure notification channels (email, Slack, etc.)
- Set up alert rules for critical errors

### 2. Alert Configuration
Set up alerts for:
- Critical errors (> 5 occurrences in 5 minutes)
- Database connection failures
- Authentication failures
- Payment processing errors
- API response time > 5 seconds

### 3. Error Grouping
Configure custom grouping rules:
- Group by user ID for user-specific issues
- Group by endpoint for API problems
- Group by error type for systematic issues

## Performance Considerations

### 1. Asynchronous Error Reporting
For high-traffic production:

```go
raygun.Asynchronous(true)  // Send errors in background
```

### 2. Error Sampling
For very high traffic applications:

```go
// Sample 10% of non-critical errors
if math.Random() > 0.1 && !isCritical(err) {
    return  // Don't send to Raygun
}
```

### 3. Rate Limiting
Implement client-side rate limiting:

```go
var lastErrorTime = make(map[string]time.Time)
var errorCooldown = 5 * time.Minute

func shouldSendError(err error) bool {
    errorKey := err.Error()
    if lastSent, exists := lastErrorTime[errorKey]; exists {
        if time.Since(lastSent) < errorCooldown {
            return false
        }
    }
    lastErrorTime[errorKey] = time.Now()
    return true
}
```

## Testing in Production

### 1. Smoke Test
Deploy with a test error to verify integration:

```go
if os.Getenv("RAYGUN_TEST") == "true" {
    go func() {
        time.Sleep(30 * time.Second)  // Wait for startup
        raygun.CreateError("Production deployment test - Raygun integration verified")
    }()
}
```

### 2. Health Check Integration
Include Raygun status in health checks:

```go
func healthHandler(c echo.Context) error {
    // Test Raygun connectivity
    testErr := raygun.CreateError("Health check test")
    
    status := map[string]interface{}{
        "status": "ok",
        "raygun": testErr == nil,
        "timestamp": time.Now().Unix(),
    }
    
    return c.JSON(http.StatusOK, status)
}
```

## Security Considerations

### 1. Sensitive Data Filtering
Ensure no sensitive data is sent:

```go
func sanitizeCustomData(data map[string]interface{}) map[string]interface{} {
    sensitiveKeys := []string{"password", "token", "secret", "key"}
    
    for key := range data {
        for _, sensitive := range sensitiveKeys {
            if strings.Contains(strings.ToLower(key), sensitive) {
                delete(data, key)
                break
            }
        }
    }
    
    return data
}
```

### 2. PII Compliance
Configure Raygun to handle PII according to your requirements:
- Enable/disable user tracking
- Configure data retention policies
- Set up GDPR compliance features

## Troubleshooting

### Common Issues
1. **API Key Invalid**: Verify the API key matches the application
2. **No Errors Appearing**: Check network connectivity and firewall rules
3. **High Memory Usage**: Enable asynchronous mode
4. **Missing Context**: Ensure middleware is properly configured

### Debug Mode
Enable debug logging for troubleshooting:

```go
raygun.LogToStdOut(true)  // Log all Raygun activities
```

## Rollback Plan
If Raygun causes issues:
1. Set `RAYGUN_ENABLED=false` environment variable
2. Add conditional initialization:
   ```go
   if os.Getenv("RAYGUN_ENABLED") != "false" {
       // Initialize Raygun
   }
   ```
3. Deploy with Raygun disabled
4. Investigate and fix the issue
5. Re-enable Raygun when resolved