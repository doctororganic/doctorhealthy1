package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SecurityHeaders returns a middleware that sets security headers
func SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("X-Frame-Options", "DENY")
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			c.Response().Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
			
			return next(c)
		}
	}
}

// AuthMiddleware provides basic authentication middleware (stub implementation)
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// For now, just set a dummy user context
			c.Set("user_id", "temp-user")
			c.Set("is_admin", false)
			return next(c)
		}
	}
}

// CORS middleware
func CORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		ExposeHeaders:    []string{echo.HeaderContentLength},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	})
}

// RateLimiter middleware (simple implementation using memory store)
func RateLimiter() echo.MiddlewareFunc {
	store := middleware.NewRateLimiterMemoryStore(20) // 20 requests per time unit
	return middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: store,
		Skipper: func(c echo.Context) bool {
			// Skip rate limiting for health endpoints
			return c.Request().URL.Path == "/health"
		},
	})
}

// CustomLogger middleware
func CustomLogger() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} method=${method} uri=${uri} status=${status} latency=${latency}\n",
	})
}

// CustomRecover middleware
func CustomRecover() echo.MiddlewareFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	})
}
