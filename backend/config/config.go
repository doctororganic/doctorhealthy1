package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Server Configuration
	ServerPort        int
	ServerHost        string
	Environment       string
	Debug             bool
	ReadTimeout       int
	WriteTimeout      int
	KeepAliveTimeout  int
	MaxRequestSize    string

	// Database Configuration
	DBHost            string
	DBPort            int
	DBName            string
	DBUser            string
	DBPassword        string
	DBSSLMode         string
	DBMaxConnections  int
	DBMaxIdleConns    int
	DBConnLifetime    time.Duration

	// Redis Configuration
	RedisHost         string
	RedisPort         int
	RedisPassword     string
	RedisDB           int
	RedisMaxRetries   int
	RedisPoolSize     int

	// Security Configuration
	JWTSecret         string
	APIKeySecret      string
	EncryptionKey     string
	CORSOrigins       []string
	RateLimitRequests int
	RateLimitWindow   time.Duration
	SecurityHeaders   bool

	// Performance Configuration
	CompressionEnabled bool
	CompressionLevel   int
	CacheTTL          time.Duration

	// Monitoring Configuration
	MetricsEnabled    bool
	MetricsPort       int
	MetricsPath       string
	LogLevel          string
	LogFormat         string
	ErrorTracking     bool
	CorrelationHeader string

	// Circuit Breaker Configuration
	CircuitBreakerEnabled   bool
	CircuitBreakerThreshold int
	CircuitBreakerTimeout   time.Duration
	CircuitBreakerMaxReqs   int

	// Pagination Configuration
	DefaultPageSize int
	MaxPageSize     int

	// Content Filtering
	ReligiousFilterEnabled bool
	FilterAlcohol         bool
	FilterPork            bool
	FilterStrictMode      bool

	// File Upload Configuration
	UploadMaxSize     string
	UploadAllowedTypes []string
	UploadPath        string

	// Backup Configuration
	BackupEnabled   bool
	BackupInterval  time.Duration
	BackupRetention time.Duration
	BackupPath      string

	// Health Check Configuration
	HealthCheckEnabled  bool
	HealthCheckInterval time.Duration
	HealthCheckTimeout  time.Duration

	// Localization
	DefaultLanguage     string
	SupportedLanguages  []string
	RTLLanguages        []string

	// PWA Configuration
	PWAName            string
	PWAShortName       string
	PWADescription     string
	PWAThemeColor      string
	PWABackgroundColor string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		// Server Configuration
		ServerPort:       getEnvAsInt("SERVER_PORT", 8080),
		ServerHost:       getEnv("SERVER_HOST", "0.0.0.0"),
		Environment:      getEnv("ENVIRONMENT", "development"),
		Debug:           getEnvAsBool("DEBUG", false),
		ReadTimeout:     getEnvAsInt("READ_TIMEOUT", 30),
		WriteTimeout:    getEnvAsInt("WRITE_TIMEOUT", 30),
		KeepAliveTimeout: getEnvAsInt("KEEP_ALIVE_TIMEOUT", 65),
		MaxRequestSize:  getEnv("MAX_REQUEST_SIZE", "10MB"),

		// Database Configuration
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnvAsInt("DB_PORT", 5432),
		DBName:           getEnv("DB_NAME", "nutrition_platform"),
		DBUser:           getEnv("DB_USER", "nutrition_user"),
		DBPassword:       getEnv("DB_PASSWORD", ""),
		DBSSLMode:        getEnv("DB_SSL_MODE", "disable"),
		DBMaxConnections: getEnvAsInt("DB_MAX_CONNECTIONS", 25),
		DBMaxIdleConns:   getEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 5),
		DBConnLifetime:   getEnvAsDuration("DB_CONNECTION_LIFETIME", "300s"),

		// Redis Configuration
		RedisHost:       getEnv("REDIS_HOST", "localhost"),
		RedisPort:       getEnvAsInt("REDIS_PORT", 6379),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisDB:         getEnvAsInt("REDIS_DB", 0),
		RedisMaxRetries: getEnvAsInt("REDIS_MAX_RETRIES", 3),
		RedisPoolSize:   getEnvAsInt("REDIS_POOL_SIZE", 10),

		// Security Configuration
		JWTSecret:         getEnv("JWT_SECRET", "your_jwt_secret_key_here"),
		APIKeySecret:      getEnv("API_KEY_SECRET", "your_api_key_secret_here"),
		EncryptionKey:     getEnv("ENCRYPTION_KEY", "your_encryption_key_here"),
		CORSOrigins:       getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
		RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getEnvAsDuration("RATE_LIMIT_WINDOW", "60s"),
		SecurityHeaders:   getEnvAsBool("SECURITY_HEADERS_ENABLED", true),

		// Performance Configuration
		CompressionEnabled: getEnvAsBool("COMPRESSION_ENABLED", true),
		CompressionLevel:   getEnvAsInt("COMPRESSION_LEVEL", 6),
		CacheTTL:          getEnvAsDuration("CACHE_TTL", "3600s"),

		// Monitoring Configuration
		MetricsEnabled:    getEnvAsBool("METRICS_ENABLED", true),
		MetricsPort:       getEnvAsInt("METRICS_PORT", 9090),
		MetricsPath:       getEnv("METRICS_PATH", "/metrics"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		LogFormat:         getEnv("LOG_FORMAT", "json"),
		ErrorTracking:     getEnvAsBool("ERROR_TRACKING_ENABLED", true),
		CorrelationHeader: getEnv("CORRELATION_ID_HEADER", "X-Correlation-ID"),

		// Circuit Breaker Configuration
		CircuitBreakerEnabled:   getEnvAsBool("CIRCUIT_BREAKER_ENABLED", true),
		CircuitBreakerThreshold: getEnvAsInt("CIRCUIT_BREAKER_THRESHOLD", 5),
		CircuitBreakerTimeout:   getEnvAsDuration("CIRCUIT_BREAKER_TIMEOUT", "60s"),
		CircuitBreakerMaxReqs:   getEnvAsInt("CIRCUIT_BREAKER_MAX_REQUESTS", 3),

		// Pagination Configuration
		DefaultPageSize: getEnvAsInt("DEFAULT_PAGE_SIZE", 20),
		MaxPageSize:     getEnvAsInt("MAX_PAGE_SIZE", 100),

		// Content Filtering
		ReligiousFilterEnabled: getEnvAsBool("RELIGIOUS_FILTER_ENABLED", true),
		FilterAlcohol:         getEnvAsBool("FILTER_ALCOHOL", true),
		FilterPork:            getEnvAsBool("FILTER_PORK", true),
		FilterStrictMode:      getEnvAsBool("FILTER_STRICT_MODE", true),

		// File Upload Configuration
		UploadMaxSize:     getEnv("UPLOAD_MAX_SIZE", "5MB"),
		UploadAllowedTypes: getEnvAsSlice("UPLOAD_ALLOWED_TYPES", []string{"image/jpeg", "image/png", "image/webp"}),
		UploadPath:        getEnv("UPLOAD_PATH", "/uploads"),

		// Backup Configuration
		BackupEnabled:   getEnvAsBool("BACKUP_ENABLED", true),
		BackupInterval:  getEnvAsDuration("BACKUP_INTERVAL", "24h"),
		BackupRetention: getEnvAsDuration("BACKUP_RETENTION", "720h"), // 30 days
		BackupPath:      getEnv("BACKUP_PATH", "/backups"),

		// Health Check Configuration
		HealthCheckEnabled:  getEnvAsBool("HEALTH_CHECK_ENABLED", true),
		HealthCheckInterval: getEnvAsDuration("HEALTH_CHECK_INTERVAL", "30s"),
		HealthCheckTimeout:  getEnvAsDuration("HEALTH_CHECK_TIMEOUT", "10s"),

		// Localization
		DefaultLanguage:    getEnv("DEFAULT_LANGUAGE", "en"),
		SupportedLanguages: getEnvAsSlice("SUPPORTED_LANGUAGES", []string{"en", "ar"}),
		RTLLanguages:       getEnvAsSlice("RTL_LANGUAGES", []string{"ar"}),

		// PWA Configuration
		PWAName:            getEnv("PWA_NAME", "Nutrition & Training Platform"),
		PWAShortName:       getEnv("PWA_SHORT_NAME", "NutriTrain"),
		PWADescription:     getEnv("PWA_DESCRIPTION", "Comprehensive nutrition and training platform"),
		PWAThemeColor:      getEnv("PWA_THEME_COLOR", "#007bff"),
		PWABackgroundColor: getEnv("PWA_BACKGROUND_COLOR", "#ffffff"),
	}
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key, defaultValue string) time.Duration {
	value := getEnv(key, defaultValue)
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	log.Printf("Invalid duration format for %s: %s, using default: %s", key, value, defaultValue)
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	return time.Minute // fallback
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// GetDatabaseURL returns the database connection URL
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

// GetRedisURL returns the Redis connection URL
func (c *Config) GetRedisURL() string {
	if c.RedisPassword != "" {
		return fmt.Sprintf("redis://:%s@%s:%d/%d", c.RedisPassword, c.RedisHost, c.RedisPort, c.RedisDB)
	}
	return fmt.Sprintf("redis://%s:%d/%d", c.RedisHost, c.RedisPort, c.RedisDB)
}