package custom

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"

	"nutrition-platform/config"
)

// Metrics holds Prometheus metrics
type Metrics struct {
	RequestsTotal       *prometheus.CounterVec
	RequestDuration     *prometheus.HistogramVec
	ActiveConnections   prometheus.Gauge
	CircuitBreakerState *prometheus.GaugeVec
}

// NewMetrics creates new Prometheus metrics
func NewMetrics() *Metrics {
	m := &Metrics{
		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		ActiveConnections: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_active_connections",
				Help: "Number of active HTTP connections",
			},
		),
		CircuitBreakerState: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "circuit_breaker_state",
				Help: "Circuit breaker state (0=closed, 1=half-open, 2=open)",
			},
			[]string{"name"},
		),
	}

	// Register metrics
	prometheus.MustRegister(m.RequestsTotal)
	prometheus.MustRegister(m.RequestDuration)
	prometheus.MustRegister(m.ActiveConnections)
	prometheus.MustRegister(m.CircuitBreakerState)

	return m
}

// CorrelationID middleware adds correlation ID to requests
func CorrelationID(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			correlationID := c.Request().Header.Get(cfg.CorrelationHeader)
			if correlationID == "" {
				correlationID = uuid.New().String()
			}
			c.Response().Header().Set(cfg.CorrelationHeader, correlationID)
			c.Set("correlation_id", correlationID)
			return next(c)
		}
	}
}

// MetricsMiddleware collects HTTP metrics
func MetricsMiddleware(metrics *Metrics) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			metrics.ActiveConnections.Inc()
			defer metrics.ActiveConnections.Dec()

			err := next(c)

			duration := time.Since(start).Seconds()
			method := c.Request().Method
			path := c.Path()
			status := strconv.Itoa(c.Response().Status)

			metrics.RequestsTotal.WithLabelValues(method, path, status).Inc()
			metrics.RequestDuration.WithLabelValues(method, path).Observe(duration)

			return err
		}
	}
}

// RateLimiter middleware implements rate limiting
func RateLimiter(cfg *config.Config) echo.MiddlewareFunc {
	limiter := rate.NewLimiter(rate.Every(cfg.RateLimitWindow/time.Duration(cfg.RateLimitRequests)), cfg.RateLimitRequests)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !limiter.Allow() {
				return echo.NewHTTPError(429, "Rate limit exceeded")
			}
			return next(c)
		}
	}
}

// SecurityHeaders middleware adds security headers
func SecurityHeaders(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.SecurityHeaders {
				c.Response().Header().Set("X-Content-Type-Options", "nosniff")
				c.Response().Header().Set("X-Frame-Options", "DENY")
				c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
				c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
				c.Response().Header().Set("Content-Security-Policy", "default-src 'self'")
				c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			}
			return next(c)
		}
	}
}

// CircuitBreakerMiddleware implements circuit breaker pattern
func CircuitBreakerMiddleware(cfg *config.Config, metrics *Metrics) echo.MiddlewareFunc {
	if !cfg.CircuitBreakerEnabled {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}

	settings := gobreaker.Settings{
		Name:        "http-requests",
		MaxRequests: uint32(cfg.CircuitBreakerMaxReqs),
		Interval:    cfg.CircuitBreakerTimeout,
		Timeout:     cfg.CircuitBreakerTimeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= uint32(cfg.CircuitBreakerThreshold)
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			var state float64
			switch to {
			case gobreaker.StateClosed:
				state = 0
			case gobreaker.StateHalfOpen:
				state = 1
			case gobreaker.StateOpen:
				state = 2
			}
			metrics.CircuitBreakerState.WithLabelValues(name).Set(state)
		},
	}

	cb := gobreaker.NewCircuitBreaker(settings)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			result, err := cb.Execute(func() (interface{}, error) {
				err := next(c)
				if err != nil {
					// Consider HTTP errors as failures
					if he, ok := err.(*echo.HTTPError); ok && he.Code >= 500 {
						return nil, err
					}
				}
				return nil, err
			})
			_ = result // Ignore result
			return err
		}
	}
}

// RequestSizeLimit middleware limits request body size
func RequestSizeLimit(maxSize string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Parse max size
			size, err := parseSize(maxSize)
			if err != nil {
				return echo.NewHTTPError(500, "Invalid max request size configuration")
			}

			// Check content length
			if c.Request().ContentLength > size {
				return echo.NewHTTPError(413, "Request entity too large")
			}

			return next(c)
		}
	}
}

// parseSize parses size strings like "10MB", "5KB", etc.
func parseSize(sizeStr string) (int64, error) {
	sizeStr = strings.ToUpper(strings.TrimSpace(sizeStr))

	if len(sizeStr) < 2 {
		return 0, fmt.Errorf("invalid size format")
	}

	// Extract number and unit
	var numStr string
	var unit string

	for i, r := range sizeStr {
		if r >= '0' && r <= '9' || r == '.' {
			numStr += string(r)
		} else {
			unit = sizeStr[i:]
			break
		}
	}

	if numStr == "" {
		return 0, fmt.Errorf("no number found in size")
	}

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %w", err)
	}

	switch unit {
	case "B", "":
		return int64(num), nil
	case "KB":
		return int64(num * 1024), nil
	case "MB":
		return int64(num * 1024 * 1024), nil
	case "GB":
		return int64(num * 1024 * 1024 * 1024), nil
	default:
		return 0, fmt.Errorf("unknown unit: %s", unit)
	}
}

// ErrorHandler provides custom error handling with correlation IDs
func ErrorHandler(cfg *config.Config) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		code := 500
		message := "Internal Server Error"

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			if he.Message != nil {
				message = fmt.Sprintf("%v", he.Message)
			}
		}

		correlationID := c.Get("correlation_id")
		if correlationID == nil {
			correlationID = "unknown"
		}

		// Log error with correlation ID
		c.Logger().Errorf("Error [%v]: %v", correlationID, err)

		// Return user-friendly error response
		response := map[string]interface{}{
			"error":          true,
			"message":        message,
			"correlation_id": correlationID,
			"timestamp":      time.Now().UTC().Format(time.RFC3339),
		}

		// Add debug info in development
		if cfg.IsDevelopment() {
			response["debug"] = err.Error()
		}

		if !c.Response().Committed {
			if c.Request().Method == "HEAD" {
				err = c.NoContent(code)
			} else {
				err = c.JSON(code, response)
			}
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}
}

// HealthCheck middleware for health check endpoints
func HealthCheck() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"uptime":    time.Since(startTime).String(),
		})
	}
}

var startTime = time.Now()
