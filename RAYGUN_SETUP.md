# Raygun Error Monitoring Setup

## Overview
This document outlines the Raygun error monitoring integration in the Nutrition Platform backend.

## Configuration
- **Application Name**: nutrition-platform
- **API Key**: J5KNVQg46P71JymsDyPWiQ
- **Library**: github.com/MindscapeHQ/raygun4go
- **Integration Location**: `backend/main.go`

## How It Works
The Raygun client is initialized at application startup and automatically:
- Captures all panics and unhandled errors
- Sends error reports to Raygun dashboard
- Includes stack traces and context information
- Provides automatic error recovery with `defer raygun.HandleError()`

## Manual Error Reporting
You can manually send errors to Raygun using:

```go
// Send an existing error
raygun.SendError(fmt.Errorf("Something went wrong"))

// Create a custom error message
raygun.CreateError("Custom error message")

// Add context information
raygun.User("user-id-123")
raygun.Tags([]string{"payment", "critical"})
raygun.CustomData(map[string]interface{}{
    "order_id": "12345",
    "amount": 99.99,
})
```

## Dashboard Access
- **URL**: https://app.raygun.com
- Login with your Raygun account credentials to view error reports

## Best Practices
1. **Add User Context**: Always set user information when available
2. **Use Tags**: Categorize errors with relevant tags
3. **Custom Data**: Include relevant business context
4. **Environment**: Set different application names for dev/staging/production

## Monitoring
- Check the Raygun dashboard regularly for new errors
- Set up alerts for critical error patterns
- Review error trends and performance impact