package custom

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// StructuredLogger wraps zap logger for Echo
type StructuredLogger struct {
	logger *zap.Logger
	cfg    *Config
}

// Config holds logger configuration
type Config struct {
	Level      string `json:"level" yaml:"level"`
	Format     string `json:"format" yaml:"format"`
	OutputPath string `json:"output_path" yaml:"output_path"`
}

// NewStructuredLogger creates a new structured logger
func NewStructuredLogger(cfg *Config) (*StructuredLogger, error) {
	var zapCfg zap.Config

	switch cfg.Format {
	case "json":
		zapCfg = zap.NewProductionConfig()
		zapCfg.OutputPaths = []string{cfg.OutputPath}
		if cfg.OutputPath == "" {
			zapCfg.OutputPaths = []string{"stdout"}
		}
	case "console":
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapCfg.OutputPaths = []string{"stdout"}
	default:
		zapCfg = zap.NewProductionConfig()
	}

	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}
	zapCfg.Level = zap.NewAtomicLevelAt(level)

	logger, err := zapCfg.Build()
	if err != nil {
		return nil, err
	}

	return &StructuredLogger{
		logger: logger,
		cfg:    cfg,
	}, nil
}

// Middleware returns Echo middleware for structured logging
func (l *StructuredLogger) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Process request
			err := next(c)

			// Calculate duration
			duration := time.Since(start)

			// Get request info
			req := c.Request()
			res := c.Response()

			// Get correlation ID
			correlationID := c.Get("correlation_id")
			if correlationID == nil {
				correlationID = "unknown"
			}

			// Build log fields
			fields := []zap.Field{
				zap.String("correlation_id", correlationID.(string)),
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.String("path", c.Path()),
				zap.Int("status", res.Status),
				zap.Duration("duration", duration),
				zap.String("ip", c.RealIP()),
				zap.String("user_agent", req.UserAgent()),
				zap.String("referer", req.Referer()),
				zap.Int64("size", res.Size),
			}

			// Add error field if there's an error
			if err != nil {
				fields = append(fields, zap.Error(err))
			}

			// Log based on status code
			switch {
			case res.Status >= 500:
				l.logger.Error("HTTP request completed with server error", fields...)
			case res.Status >= 400:
				l.logger.Warn("HTTP request completed with client error", fields...)
			case res.Status >= 300:
				l.logger.Info("HTTP request completed with redirection", fields...)
			default:
				l.logger.Info("HTTP request completed successfully", fields...)
			}

			return err
		}
	}
}

// AccessLog creates a middleware that logs HTTP requests in a traditional format
func (l *StructuredLogger) AccessLog() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			duration := time.Since(start)
			req := c.Request()
			res := c.Response()

			correlationID := c.Get("correlation_id")
			if correlationID == nil {
				correlationID = "unknown"
			}

			l.logger.Info("HTTP Access",
				zap.String("correlation_id", correlationID.(string)),
				zap.String("remote_ip", c.RealIP()),
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.Int("status", res.Status),
				zap.Duration("duration", duration),
				zap.Int64("size", res.Size),
				zap.String("user_agent", req.UserAgent()),
			)

			return err
		}
	}
}

// RequestLogger logs incoming requests with detailed information
func (l *StructuredLogger) RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()

			correlationID := c.Get("correlation_id")
			if correlationID == nil {
				correlationID = "unknown"
			}

			fields := []zap.Field{
				zap.String("correlation_id", correlationID.(string)),
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.String("path", c.Path()),
				zap.String("ip", c.RealIP()),
				zap.String("user_agent", req.UserAgent()),
				zap.String("referer", req.Referer()),
				zap.String("host", req.Host),
				zap.String("protocol", req.Proto),
			}

			// Add query parameters if present
			if req.URL.RawQuery != "" {
				fields = append(fields, zap.String("query", req.URL.RawQuery))
			}

			// Add content type if present
			if contentType := req.Header.Get("Content-Type"); contentType != "" {
				fields = append(fields, zap.String("content_type", contentType))
			}

			// Add content length if present
			if req.ContentLength > 0 {
				fields = append(fields, zap.Int64("content_length", req.ContentLength))
			}

			l.logger.Info("Incoming HTTP request", fields...)

			return next(c)
		}
	}
}

// ErrorLogger logs errors with context
func (l *StructuredLogger) ErrorLogger() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		correlationID := c.Get("correlation_id")
		if correlationID == nil {
			correlationID = "unknown"
		}

		req := c.Request()
		res := c.Response()

		fields := []zap.Field{
			zap.String("correlation_id", correlationID.(string)),
			zap.Error(err),
			zap.String("method", req.Method),
			zap.String("uri", req.RequestURI),
			zap.String("path", c.Path()),
			zap.String("ip", c.RealIP()),
			zap.String("user_agent", req.UserAgent()),
			zap.Int("status", res.Status),
		}

		// Add stack trace for server errors
		if res.Status >= 500 {
			fields = append(fields, zap.String("stack_trace", "stack trace would be here"))
		}

		l.logger.Error("HTTP error occurred", fields...)

		// Call the original error handler or provide default response
		if !c.Response().Committed {
			if c.Request().Method == "HEAD" {
				c.NoContent(res.Status)
			} else {
				c.JSON(res.Status, map[string]interface{}{
					"error":          true,
					"message":        "Internal server error",
					"correlation_id": correlationID,
					"timestamp":      time.Now().UTC().Format(time.RFC3339),
				})
			}
		}
	}
}

// Sync flushes any buffered log entries
func (l *StructuredLogger) Sync() {
	l.logger.Sync()
}

// Close closes the logger
func (l *StructuredLogger) Close() error {
	return l.logger.Sync()
}
