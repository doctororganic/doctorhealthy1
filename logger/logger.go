package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"nutrition-platform/config"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger wraps logrus logger with additional functionality
type Logger struct {
	*logrus.Logger
	config config.LoggingConfig
}

// NewLogger creates a new logger instance
func NewLogger(cfg config.LoggingConfig) *Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set formatter
	if cfg.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "function",
				logrus.FieldKeyFile:  "file",
			},
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}

	// Set output
	if cfg.File != "" {
		// Create directory if it doesn't exist
		dir := filepath.Dir(cfg.File)
		if err := os.MkdirAll(dir, 0755); err != nil {
			logger.SetOutput(os.Stdout)
			logger.WithError(err).Error("Failed to create log directory, using stdout")
		} else {
			// Use lumberjack for log rotation
			logger.SetOutput(&lumberjack.Logger{
				Filename:   cfg.File,
				MaxSize:    cfg.MaxSize,    // megabytes
				MaxBackups: cfg.MaxBackups, // number of backups
				MaxAge:     cfg.MaxAge,     // days
				Compress:   true,
			})
		}
	} else {
		logger.SetOutput(os.Stdout)
	}

	appLogger := &Logger{
		Logger: logger,
		config: cfg,
	}

	return appLogger
}

// Info logs an info message with fields
func (l *Logger) Info(message string, fields map[string]interface{}) {
	l.WithFields(fields).Info(message)
}

// Error logs an error message with fields
func (l *Logger) Error(message string, fields map[string]interface{}) {
	l.WithFields(fields).Error(message)
}

// Warn logs a warning message with fields
func (l *Logger) Warn(message string, fields map[string]interface{}) {
	l.WithFields(fields).Warn(message)
}

// Debug logs a debug message with fields
func (l *Logger) Debug(message string, fields map[string]interface{}) {
	l.WithFields(fields).Debug(message)
}

// Fatal logs a fatal message with fields and exits
func (l *Logger) Fatal(message string, fields map[string]interface{}) {
	l.WithFields(fields).Fatal(message)
}

// Panic logs a panic message with fields and panics
func (l *Logger) Panic(message string, fields map[string]interface{}) {
	l.WithFields(fields).Panic(message)
}

// WithRequest adds request-related fields to the logger
func (l *Logger) WithRequest(c echo.Context) *logrus.Entry {
	fields := logrus.Fields{
		"method":           c.Request().Method,
		"path":             c.Request().URL.Path,
		"query":            c.Request().URL.RawQuery,
		"remote_ip":        c.RealIP(),
		"user_agent":       c.Request().UserAgent(),
		"request_id":       c.Response().Header().Get("X-Request-ID"),
		"correlation_id":   c.Get("correlation_id"),
	}

	// Add user information if available
	if userID, ok := c.Get("user_id").(uint); ok {
		fields["user_id"] = userID
	}
	if userRole, ok := c.Get("user_role").(string); ok {
		fields["user_role"] = userRole
	}

	return l.WithFields(fields)
}

// WithError adds error information to the logger
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// WithField adds a single field to the logger
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

// GetEchoLogger returns a logger compatible with Echo framework
func (l *Logger) GetEchoLogger() echo.Logger {
	return echo.NewEchoLogger(l.Logger.Writer())
}

// GetLogLevel returns the current log level
func (l *Logger) GetLogLevel() logrus.Level {
	return l.Logger.GetLevel()
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level logrus.Level) {
	l.Logger.SetLevel(level)
}

// SetFormat sets the log format
func (l *Logger) SetFormat(format string) {
	if format == "json" {
		l.Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		l.Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}
}

// HTTPLogger middleware for Echo
func HTTPLogger(logger *Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Process request
			err := next(c)

			// Calculate request duration
			duration := time.Since(start)

			// Log request
			fields := map[string]interface{}{
				"method":          c.Request().Method,
				"path":            c.Request().URL.Path,
				"status":          c.Response().Status,
				"duration":        duration.String(),
				"duration_ms":     duration.Milliseconds(),
				"remote_ip":       c.RealIP(),
				"user_agent":      c.Request().UserAgent(),
				"request_id":      c.Response().Header().Get("X-Request-ID"),
				"correlation_id":  c.Get("correlation_id"),
				"response_size":   c.Response().Size(),
			}

			// Add user information if available
			if userID, ok := c.Get("user_id").(uint); ok {
				fields["user_id"] = userID
			}
			if userRole, ok := c.Get("user_role").(string); ok {
				fields["user_role"] = userRole
			}

			// Add error information if available
			if err != nil {
				fields["error"] = err.Error()
				logger.Error("HTTP Request Error", fields)
			} else {
				// Use different log levels based on status code
				switch {
				case c.Response().Status >= 500:
					logger.Error("HTTP Request", fields)
				case c.Response().Status >= 400:
					logger.Warn("HTTP Request", fields)
				default:
					logger.Info("HTTP Request", fields)
				}
			}

			return err
		}
	}
}

// RequestLogger middleware for detailed request logging
func RequestLogger(logger *Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Add response time header
			c.Response().Header().Set("X-Response-Time", "0ms")

			start := time.Now()
			err := next(c)
			duration := time.Since(start)

			// Update response time header
			c.Response().Header().Set("X-Response-Time", duration.String())

			// Log detailed request information
			fields := logger.WithRequest(c).WithFields(map[string]interface{}{
				"status":       c.Response().Status,
				"duration":     duration.String(),
				"duration_ms":  duration.Milliseconds(),
				"size":         c.Response().Size(),
				"latency":      duration.Seconds(),
			})

			if err != nil {
				fields.WithError(err).Error("Request completed with error")
			} else {
				fields.Info("Request completed")
			}

			return err
		}
	}
}

// AuditLogger for security and compliance logging
type AuditLogger struct {
	*Logger
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(cfg config.LoggingConfig) *AuditLogger {
	// Audit logs should always be in JSON format for easy parsing
	auditCfg := cfg
	auditCfg.Format = "json"
	auditCfg.Level = "info" // Always log audit events

	return &AuditLogger{
		Logger: NewLogger(auditCfg),
	}
}

// LogAuthentication logs authentication events
func (a *AuditLogger) LogAuthentication(userID uint, email, action, ip, userAgent string, success bool) {
	a.WithFields(map[string]interface{}{
		"event_type":    "authentication",
		"user_id":       userID,
		"email":         email,
		"action":        action, // login, logout, register, password_reset
		"ip_address":    ip,
		"user_agent":    userAgent,
		"success":       success,
		"timestamp":     time.Now().UTC(),
	}).Info("Authentication event")
}

// LogDataAccess logs data access events
func (a *AuditLogger) LogDataAccess(userID uint, resource, action, ip string) {
	a.WithFields(map[string]interface{}{
		"event_type": "data_access",
		"user_id":    userID,
		"resource":   resource, // users, foods, recipes, etc.
		"action":     action,   // read, write, delete
		"ip_address": ip,
		"timestamp":  time.Now().UTC(),
	}).Info("Data access event")
}

// LogConfigurationChange logs configuration changes
func (a *AuditLogger) LogConfigurationChange(userID uint, section, field, oldValue, newValue string) {
	a.WithFields(map[string]interface{}{
		"event_type":  "configuration_change",
		"user_id":     userID,
		"section":     section,
		"field":       field,
		"old_value":   oldValue,
		"new_value":   newValue,
		"timestamp":   time.Now().UTC(),
	}).Info("Configuration change event")
}

// LogSecurityEvent logs security-related events
func (a *AuditLogger) LogSecurityEvent(eventType, description, severity, ip string, details map[string]interface{}) {
	fields := map[string]interface{}{
		"event_type":  "security",
		"type":        eventType,
		"description": description,
		"severity":    severity, // low, medium, high, critical
		"ip_address":  ip,
		"timestamp":   time.Now().UTC(),
	}

	// Add additional details
	for k, v := range details {
		fields[k] = v
	}

	a.WithFields(fields).Info("Security event")
}

// PerformanceLogger for performance monitoring
type PerformanceLogger struct {
	*Logger
}

// NewPerformanceLogger creates a new performance logger
func NewPerformanceLogger(cfg config.LoggingConfig) *PerformanceLogger {
	return &PerformanceLogger{
		Logger: NewLogger(cfg),
	}
}

// LogDatabaseQuery logs database query performance
func (p *PerformanceLogger) LogDatabaseQuery(query string, duration time.Duration, rowsAffected int, err error) {
	fields := map[string]interface{}{
		"event_type":    "database_query",
		"query":         query,
		"duration_ms":   duration.Milliseconds(),
		"rows_affected": rowsAffected,
		"timestamp":     time.Now().UTC(),
	}

	if err != nil {
		fields["error"] = err.Error()
		p.WithFields(fields).Error("Database query failed")
	} else {
		p.WithFields(fields).Debug("Database query executed")
	}
}

// LogCacheOperation logs cache operations
func (p *PerformanceLogger) LogCacheOperation(operation, key string, hit bool, duration time.Duration) {
	p.WithFields(map[string]interface{}{
		"event_type":  "cache_operation",
		"operation":   operation, // get, set, delete
		"key":         key,
		"hit":         hit,
		"duration_ms": duration.Milliseconds(),
		"timestamp":   time.Now().UTC(),
	}).Debug("Cache operation")
}

// LogAPIRequest logs external API request performance
func (p *PerformanceLogger) LogAPIRequest(url, method string, statusCode int, duration time.Duration, err error) {
	fields := map[string]interface{}{
		"event_type":  "external_api_request",
		"url":         url,
		"method":      method,
		"status_code": statusCode,
		"duration_ms": duration.Milliseconds(),
		"timestamp":   time.Now().UTC(),
	}

	if err != nil {
		fields["error"] = err.Error()
		p.WithFields(fields).Error("External API request failed")
	} else {
		p.WithFields(fields).Debug("External API request completed")
	}
}

// SlowQueryDetector detects and logs slow queries
type SlowQueryDetector struct {
	logger    *Logger
	threshold time.Duration
}

// NewSlowQueryDetector creates a new slow query detector
func NewSlowQueryDetector(logger *Logger, threshold time.Duration) *SlowQueryDetector {
	return &SlowQueryDetector{
		logger:    logger,
		threshold: threshold,
	}
}

// CheckQuery checks if a query is slow and logs it
func (s *SlowQueryDetector) CheckQuery(query string, duration time.Duration, err error) {
	if duration > s.threshold {
		fields := map[string]interface{}{
			"event_type":  "slow_query",
			"query":       query,
			"duration_ms": duration.Milliseconds(),
			"threshold_ms": s.threshold.Milliseconds(),
			"timestamp":   time.Now().UTC(),
		}

		if err != nil {
			fields["error"] = err.Error()
		}

		s.logger.WithFields(fields).Warn("Slow query detected")
	}
}
