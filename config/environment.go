package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all application configuration
type Config struct {
	// Application
	App        AppConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	JWT        JWTConfig
	Security   SecurityConfig
	CORS       CORSConfig
	Email      EmailConfig
	External   ExternalConfig
	Logging    LoggingConfig
	Monitoring MonitoringConfig
	Cache      CacheConfig
	Session    SessionConfig
	Upload     UploadConfig
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name        string
	Version     string
	Environment string
	Port        string
	URL         string
	APIURL      string
	Debug       bool
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host                string
	Port                string
	Name                string
	User                string
	Password            string
	SSLMode             string
	MaxConnections      int
	MaxIdleConnections  int
	ConnectionMaxLifetime time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret           string
	RefreshSecret    string
	AccessExpiry     time.Duration
	RefreshExpiry    time.Duration
	Issuer           string
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	APISecret                string
	BCryptRounds            int
	RateLimitRPM            int
	RateLimitBurst          int
	EnableCSRF              bool
	EnableContentSecurity   bool
	EnableHTTPSRedirect     bool
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPSecure   bool
	SMTPUser     string
	SMTPPassword string
	From         string
}

// ExternalConfig holds external API configuration
type ExternalConfig struct {
	NutritionAPIKey     string
	NutritionAPIBaseURL string
	GoogleClientID      string
	GoogleClientSecret  string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string
	Format     string
	File       string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

// MonitoringConfig holds monitoring configuration
type MonitoringConfig struct {
	Enabled               bool
	MetricsPort           string
	HealthCheckInterval   time.Duration
	SentryDSN             string
	SentryEnvironment     string
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Default      time.Duration
	UserProfile  time.Duration
	NutritionData time.Duration
	SearchResults time.Duration
}

// SessionConfig holds session configuration
type SessionConfig struct {
	Secret           string
	CookieSecure     bool
	CookieHTTPOnly   bool
	CookieMaxAge     int
}

// UploadConfig holds file upload configuration
type UploadConfig struct {
	MaxFileSize    int64
	AllowedTypes   []string
	UploadPath     string
	TempUploadPath string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	cfg := &Config{}

	// Application config
	cfg.App = AppConfig{
		Name:        getEnv("APP_NAME", "Nutrition Platform"),
		Version:     getEnv("APP_VERSION", "1.0.0"),
		Environment: getEnv("NODE_ENV", "development"),
		Port:        getEnv("PORT", "8080"),
		URL:         getEnv("APP_URL", "http://localhost:3000"),
		APIURL:      getEnv("API_URL", "http://localhost:8080"),
		Debug:       getBoolEnv("DEBUG", false),
	}

	// Database config
	cfg.Database = DatabaseConfig{
		Host:                getEnv("DB_HOST", "localhost"),
		Port:                getEnv("DB_PORT", "5432"),
		Name:                getEnv("DB_NAME", "nutrition_platform"),
		User:                getEnv("DB_USER", "nutrition_user"),
		Password:            getEnv("DB_PASSWORD", ""),
		SSLMode:             getEnv("DB_SSLMODE", "prefer"),
		MaxConnections:      getIntEnv("DB_MAX_CONNECTIONS", 25),
		MaxIdleConnections:  getIntEnv("DB_MAX_IDLE_CONNECTIONS", 5),
		ConnectionMaxLifetime: getDurationEnv("DB_CONNECTION_MAX_LIFETIME", 5*time.Minute),
	}

	// Redis config
	cfg.Redis = RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       getIntEnv("REDIS_DB", 0),
		PoolSize: getIntEnv("REDIS_POOL_SIZE", 10),
	}

	// JWT config
	cfg.JWT = JWTConfig{
		Secret:        getEnv("JWT_SECRET", ""),
		RefreshSecret: getEnv("JWT_REFRESH_SECRET", ""),
		AccessExpiry:  getDurationEnv("JWT_ACCESS_TOKEN_EXPIRY", 15*time.Minute),
		RefreshExpiry: getDurationEnv("JWT_REFRESH_TOKEN_EXPIRY", 7*24*time.Hour),
		Issuer:        getEnv("JWT_ISSUER", "nutrition-platform"),
	}

	// Security config
	cfg.Security = SecurityConfig{
		APISecret:              getEnv("API_SECRET", ""),
		BCryptRounds:          getIntEnv("BCRYPT_ROUNDS", 12),
		RateLimitRPM:          getIntEnv("RATE_LIMIT_REQUESTS_PER_MINUTE", 100),
		RateLimitBurst:        getIntEnv("RATE_LIMIT_BURST", 200),
		EnableCSRF:            getBoolEnv("ENABLE_CSRF", true),
		EnableContentSecurity: getBoolEnv("ENABLE_CONTENT_SECURITY", true),
		EnableHTTPSRedirect:   getBoolEnv("ENABLE_HTTPS_REDIRECT", false),
	}

	// CORS config
	cfg.CORS = CORSConfig{
		AllowedOrigins: getSliceEnv("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		AllowedMethods: getSliceEnv("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		AllowedHeaders: getSliceEnv("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization", "X-Requested-With"}),
	}

	// Email config
	cfg.Email = EmailConfig{
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getIntEnv("SMTP_PORT", 587),
		SMTPSecure:   getBoolEnv("SMTP_SECURE", false),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		From:         getEnv("EMAIL_FROM", "Nutrition Platform <noreply@nutritionplatform.com>"),
	}

	// External API config
	cfg.External = ExternalConfig{
		NutritionAPIKey:     getEnv("NUTRITION_API_KEY", ""),
		NutritionAPIBaseURL: getEnv("NUTRITION_API_BASE_URL", "https://api.nutrition.gov"),
		GoogleClientID:      getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:  getEnv("GOOGLE_CLIENT_SECRET", ""),
	}

	// Logging config
	cfg.Logging = LoggingConfig{
		Level:      getEnv("LOG_LEVEL", "info"),
		Format:     getEnv("LOG_FORMAT", "json"),
		File:       getEnv("LOG_FILE", "./logs/app.log"),
		MaxSize:    getIntEnv("LOG_MAX_SIZE", 100),
		MaxBackups: getIntEnv("LOG_MAX_BACKUPS", 10),
		MaxAge:     getIntEnv("LOG_MAX_AGE", 30),
	}

	// Monitoring config
	cfg.Monitoring = MonitoringConfig{
		Enabled:             getBoolEnv("METRICS_ENABLED", true),
		MetricsPort:         getEnv("METRICS_PORT", "9090"),
		HealthCheckInterval: getDurationEnv("HEALTH_CHECK_INTERVAL", 30*time.Second),
		SentryDSN:           getEnv("SENTRY_DSN", ""),
		SentryEnvironment:   getEnv("SENTRY_ENVIRONMENT", "development"),
	}

	// Cache config
	cfg.Cache = CacheConfig{
		Default:        getDurationEnv("CACHE_TTL_DEFAULT", 1*time.Hour),
		UserProfile:    getDurationEnv("CACHE_TTL_USER_PROFILE", 30*time.Minute),
		NutritionData:  getDurationEnv("CACHE_TTL_NUTRITION_DATA", 24*time.Hour),
		SearchResults:  getDurationEnv("CACHE_TTL_SEARCH_RESULTS", 15*time.Minute),
	}

	// Session config
	cfg.Session = SessionConfig{
		Secret:         getEnv("SESSION_SECRET", ""),
		CookieSecure:   getBoolEnv("SESSION_COOKIE_SECURE", false),
		CookieHTTPOnly: getBoolEnv("SESSION_COOKIE_HTTP_ONLY", true),
		CookieMaxAge:   getIntEnv("SESSION_COOKIE_MAX_AGE", 86400),
	}

	// Upload config
	cfg.Upload = UploadConfig{
		MaxFileSize:    getInt64Env("MAX_FILE_SIZE", 10*1024*1024), // 10MB
		AllowedTypes:   getSliceEnv("ALLOWED_FILE_TYPES", []string{"jpg", "jpeg", "png", "gif", "pdf", "doc", "docx"}),
		UploadPath:     getEnv("UPLOAD_PATH", "./uploads"),
		TempUploadPath: getEnv("TEMP_UPLOAD_PATH", "./temp"),
	}

	// Validate critical configuration
	if err := cfg.validate(); err != nil {
		panic(fmt.Sprintf("Configuration validation failed: %v", err))
	}

	return cfg
}

// validate validates critical configuration values
func (c *Config) validate() error {
	// Validate JWT secrets
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if len(c.JWT.Secret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}

	if c.JWT.RefreshSecret == "" {
		return fmt.Errorf("JWT_REFRESH_SECRET is required")
	}
	if len(c.JWT.RefreshSecret) < 32 {
		return fmt.Errorf("JWT_REFRESH_SECRET must be at least 32 characters long")
	}

	// Validate database configuration
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}

	// Validate API secret
	if c.Security.APISecret == "" {
		return fmt.Errorf("API_SECRET is required")
	}

	// Validate session secret
	if c.Session.Secret == "" {
		return fmt.Errorf("SESSION_SECRET is required")
	}

	return nil
}

// IsProduction returns true if the application is running in production
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// IsDevelopment returns true if the application is running in development
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getInt64Env(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
