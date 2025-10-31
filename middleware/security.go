package middleware

import (
	"net/http"
	"strings"

	"nutrition-platform/config"

	"github.com/labstack/echo/v4"
)

// SecurityHeaders adds comprehensive security headers
func SecurityHeaders(cfg config.SecurityConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Content Security Policy
			if cfg.EnableContentSecurity {
				csp := "default-src 'self'; " +
					"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
					"style-src 'self' 'unsafe-inline'; " +
					"img-src 'self' data: https:; " +
					"font-src 'self'; " +
					"connect-src 'self'; " +
					"frame-ancestors 'none'; " +
					"base-uri 'self'; " +
					"form-action 'self'"
				
				c.Response().Header().Set("Content-Security-Policy", csp)
			}

			// XSS Protection
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")

			// Content Type Options
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")

			// Frame Options
			c.Response().Header().Set("X-Frame-Options", "DENY")

			// Referrer Policy
			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Permissions Policy
			permissionsPolicy := "geolocation=(), " +
				"microphone=(), " +
				"camera=(), " +
				"payment=(), " +
				"usb=(), " +
				"magnetometer=(), " +
				"gyroscope=(), " +
				"accelerometer=()"
			c.Response().Header().Set("Permissions-Policy", permissionsPolicy)

			// Strict Transport Security (only in HTTPS)
			if c.Request().TLS != nil {
				hsts := "max-age=31536000; includeSubDomains; preload"
				c.Response().Header().Set("Strict-Transport-Security", hsts)
			}

			// Remove server information
			c.Response().Header().Set("Server", "")

			return next(c)
		}
	}
}

// CSRFProtection adds CSRF protection
func CSRFProtection() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip CSRF for GET, HEAD, OPTIONS requests
			if c.Request().Method == "GET" || 
			   c.Request().Method == "HEAD" || 
			   c.Request().Method == "OPTIONS" {
				return next(c)
			}

			// Skip CSRF for API authentication endpoints
			if strings.HasPrefix(c.Request().URL.Path, "/api/v1/auth/") {
				return next(c)
			}

			// Get CSRF token from header
			csrfToken := c.Request().Header.Get("X-CSRF-Token")
			if csrfToken == "" {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"error": "CSRF token required",
					"code":  "CSRF_REQUIRED",
				})
			}

			// Get expected CSRF token from session/cookie
			expectedToken := c.Request().Header.Get("X-Expected-CSRF")
			if expectedToken == "" || csrfToken != expectedToken {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"error": "Invalid CSRF token",
					"code":  "INVALID_CSRF",
				})
			}

			return next(c)
		}
	}
}

// RequestValidation validates incoming requests
func RequestValidation() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Validate content length
			if c.Request().ContentLength > 10*1024*1024 { // 10MB limit
				return c.JSON(http.StatusRequestEntityTooLarge, map[string]interface{}{
					"error": "Request entity too large",
					"code":  "PAYLOAD_TOO_LARGE",
				})
			}

			// Validate content type for POST/PUT requests
			if c.Request().Method == "POST" || c.Request().Method == "PUT" {
				contentType := c.Request().Header.Get("Content-Type")
				if !strings.Contains(contentType, "application/json") &&
				   !strings.Contains(contentType, "multipart/form-data") &&
				   !strings.Contains(contentType, "application/x-www-form-urlencoded") {
					return c.JSON(http.StatusUnsupportedMediaType, map[string]interface{}{
						"error": "Unsupported media type",
						"code":  "UNSUPPORTED_MEDIA_TYPE",
					})
				}
			}

			// Add request ID if not present
			requestID := c.Response().Header().Get("X-Request-ID")
			if requestID == "" {
				requestID = generateRequestID()
				c.Response().Header().Set("X-Request-ID", requestID)
			}

			return next(c)
		}
	}
}

// IPWhitelist restricts access to specific IP addresses
func IPWhitelist(allowedIPs []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if len(allowedIPs) == 0 {
				return next(c) // No whitelist configured
			}

			clientIP := c.RealIP()
			if clientIP == "" {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"error": "Unable to determine client IP",
					"code":  "IP_DETECTION_FAILED",
				})
			}

			// Check if client IP is in whitelist
			allowed := false
			for _, ip := range allowedIPs {
				if ip == clientIP {
					allowed = true
					break
				}
			}

			if !allowed {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"error": "Access denied from this IP",
					"code":  "IP_NOT_ALLOWED",
				})
			}

			return next(c)
		}
	}
}

// APIKeyMiddleware validates API keys
func APIKeyMiddleware(validKeys []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip API key validation for auth endpoints
			if strings.HasPrefix(c.Request().URL.Path, "/api/v1/auth/") {
				return next(c)
			}

			// Get API key from header
			apiKey := c.Request().Header.Get("X-API-Key")
			if apiKey == "" {
				// Try to get from query parameter
				apiKey = c.QueryParam("api_key")
			}

			if apiKey == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "API key required",
					"code":  "API_KEY_REQUIRED",
				})
			}

			// Validate API key
			valid := false
			for _, key := range validKeys {
				if key == apiKey {
					valid = true
					break
				}
			}

			if !valid {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "Invalid API key",
					"code":  "INVALID_API_KEY",
				})
			}

			// Store API key in context for rate limiting
			c.Set("api_key", apiKey)

			return next(c)
		}
	}
}

// RateLimitByUser applies rate limiting based on user ID
func RateLimitByUser(requestsPerMinute int) echo.MiddlewareFunc {
	// This would integrate with Redis for distributed rate limiting
	// For now, it's a placeholder implementation
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user ID from context (set by JWT middleware)
			userID, ok := c.Get("user_id").(uint)
			if !ok {
				// If no user ID, use IP address for rate limiting
				userID = hashIP(c.RealIP())
			}

			// Store user ID for rate limiting
			c.Set("rate_limit_key", userID)

			return next(c)
		}
	}
}

// CorrelationID adds or forwards correlation ID for request tracing
func CorrelationID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check for existing correlation ID
			correlationID := c.Request().Header.Get("X-Correlation-ID")
			if correlationID == "" {
				correlationID = generateCorrelationID()
			}

			// Set correlation ID in response header
			c.Response().Header().Set("X-Correlation-ID", correlationID)

			// Store in context for logging
			c.Set("correlation_id", correlationID)

			return next(c)
		}
	}
}

// ResponseLogger logs response information
func ResponseLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Process request
			err := next(c)

			// Log response
			status := c.Response().Status
			method := c.Request().Method
			path := c.Request().URL.Path
			userAgent := c.Request().UserAgent()
			clientIP := c.RealIP()

			// Get correlation ID if available
			correlationID, _ := c.Get("correlation_id").(string)

			// Log response details (this would integrate with your logging system)
			logResponse(map[string]interface{}{
				"method":         method,
				"path":           path,
				"status":         status,
				"user_agent":     userAgent,
				"client_ip":      clientIP,
				"correlation_id": correlationID,
				"response_time":  c.Response().Header().Get("X-Response-Time"),
			})

			return err
		}
	}
}

// Helper functions

func generateRequestID() string {
	// Generate a unique request ID
	return "req_" + generateRandomString(16)
}

func generateCorrelationID() string {
	// Generate a unique correlation ID
	return "corr_" + generateRandomString(32)
}

func generateRandomString(length int) string {
	// Simple random string generator (in production, use crypto/rand)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		// This is a simplified version - use crypto/rand in production
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}

func hashIP(ip string) uint {
	// Simple hash function for IP addresses
	// In production, use a proper hashing algorithm
	hash := uint(0)
	for _, c := range ip {
		hash = hash*31 + uint(c)
	}
	return hash
}

func logResponse(data map[string]interface{}) {
	// This would integrate with your logging system
	// For now, just a placeholder
	// log.Info("HTTP Response", data)
}
