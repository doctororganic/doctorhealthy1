package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	// Server Configuration
	ServerHost string
	ServerPort string
	Port       string

	// Database Configuration
	DBHost         string
	DBPort         int
	DBName         string
	DBUser         string
	DBPassword     string
	DBSSLMode      string
	DBMaxConns     int
	DBMaxIdleConns int
	DBConnLifetime time.Duration

	// Redis Configuration
	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int
	RedisMaxRetries int
	RedisPoolSize int

	// Security Configuration
	JWTSecret      string
	APIKeySecret   string
	EncryptionKey  string

	// Server timeouts
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	KeepAliveTimeout time.Duration

	// CORS Configuration
	CORSOrigins []string

	// Rate Limiting
	RateLimitRequests int
	RateLimitWindow   time.Duration

	// Logging
	LogLevel  string
	LogFormat string

	// Features
	ReligiousFilterEnabled bool
	FilterAlcohol         bool
	FilterPork           bool
	FilterStrictMode     bool

	// Health Check
	HealthCheckEnabled    bool
	HealthCheckInterval   time.Duration
	HealthCheckTimeout    time.Duration

	// Internationalization
	DefaultLanguage   string
	SupportedLanguages []string
	RTLLanguages      []string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		ServerHost: getEnv("SERVER_HOST", "localhost"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		Port:       getEnv("SERVER_PORT", "8080"),

		// Database Configuration
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnvAsInt("DB_PORT", 5432),
		DBName:         getEnv("DB_NAME", "nutrition_platform"),
		DBUser:         getEnv("DB_USER", "nutrition_user"),
		DBPassword:     getEnv("DB_PASSWORD", "REPLACE_WITH_STRONG_32_CHAR_PASSWORD"),
		DBSSLMode:      getEnv("DB_SSL_MODE", "require"),
		DBMaxConns:     getEnvAsInt("DB_MAX_CONNECTIONS", 25),
		DBMaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 5),
		DBConnLifetime: getEnvAsDuration("DB_CONNECTION_LIFETIME", 300*time.Second),

		// Redis Configuration
		RedisHost:       getEnv("REDIS_HOST", "redis"),
		RedisPort:       getEnvAsInt("REDIS_PORT", 6379),
		RedisPassword:   getEnv("REDIS_PASSWORD", "REPLACE_WITH_STRONG_32_CHAR_PASSWORD"),
		RedisDB:         getEnvAsInt("REDIS_DB", 0),
		RedisMaxRetries: getEnvAsInt("REDIS_MAX_RETRIES", 3),
		RedisPoolSize:   getEnvAsInt("REDIS_POOL_SIZE", 10),

		// Security Configuration
		JWTSecret:     getEnv("JWT_SECRET", "REPLACE_WITH_STRONG_64_CHAR_SECRET_KEY"),
		APIKeySecret:  getEnv("API_KEY_SECRET", "REPLACE_WITH_STRONG_64_CHAR_API_KEY"),
		EncryptionKey: getEnv("ENCRYPTION_KEY", "REPLACE_WITH_STRONG_32_CHAR_ENCRYPTION_KEY"),

		// Server timeouts
		ReadTimeout:  time.Duration(getEnvAsInt("READ_TIMEOUT", 30)) * time.Second,
		WriteTimeout: time.Duration(getEnvAsInt("WRITE_TIMEOUT", 30)) * time.Second,
		KeepAliveTimeout: time.Duration(getEnvAsInt("KEEP_ALIVE_TIMEOUT", 65)) * time.Second,

		// CORS Configuration
		CORSOrigins: getEnvAsSlice("CORS_ORIGINS", []string{"http://localhost:3000", "http://localhost:8080"}),

		// Rate Limiting
		RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getEnvAsDuration("RATE_LIMIT_WINDOW", 60*time.Second),

		// Logging
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "json"),

		// Features
		ReligiousFilterEnabled: getEnvAsBool("RELIGIOUS_FILTER_ENABLED", true),
		FilterAlcohol:         getEnvAsBool("FILTER_ALCOHOL", true),
		FilterPork:           getEnvAsBool("FILTER_PORK", true),
		FilterStrictMode:     getEnvAsBool("FILTER_STRICT_MODE", false),

		// Health Check
		HealthCheckEnabled:  getEnvAsBool("HEALTH_CHECK_ENABLED", true),
		HealthCheckInterval: getEnvAsDuration("HEALTH_CHECK_INTERVAL", 30*time.Second),
		HealthCheckTimeout:  getEnvAsDuration("HEALTH_CHECK_TIMEOUT", 5*time.Second),

		// Internationalization
		DefaultLanguage:    getEnv("DEFAULT_LANGUAGE", "en"),
		SupportedLanguages: getEnvAsSlice("SUPPORTED_LANGUAGES", []string{"en", "ar"}),
		RTLLanguages:       getEnvAsSlice("RTL_LANGUAGES", []string{"ar"}),
	}

	return config
}

// Helper functions
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

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if durationValue, err := time.ParseDuration(value); err == nil {
			return durationValue
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Simple comma-separated parsing
		if value != "" {
			return []string{value} // For now, return as single element
		}
	}
	return defaultValue
}