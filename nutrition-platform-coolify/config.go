// config.go
package main

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port             string
	DBPath           string
	JWTSecret        string
	APIRateLimit     int
	EnableHTTPS      bool
	Domain           string
	LogLevel         string
	DBDriver         string
	MaxRequestSize   int64
	ReadTimeout      int
	WriteTimeout     int
	KeepAliveTimeout int
}

func LoadConfig() Config {
	config := Config{
		Port:             getEnv("PORT", "8080"),
		DBPath:           getEnv("DB_PATH", "./nutrition_platform.db"),
		JWTSecret:        getEnv("JWT_SECRET", "default-secret-change-in-production"),
		APIRateLimit:     getEnvAsInt("API_RATE_LIMIT", 100),
		EnableHTTPS:      getEnvAsBool("ENABLE_HTTPS", false),
		Domain:           getEnv("DOMAIN", "localhost"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		DBDriver:         getEnv("DB_DRIVER", "sqlite3"),
		MaxRequestSize:   getEnvAsInt64("MAX_REQUEST_SIZE", 10*1024*1024), // 10MB default
		ReadTimeout:      getEnvAsInt("READ_TIMEOUT", 30),
		WriteTimeout:     getEnvAsInt("WRITE_TIMEOUT", 30),
		KeepAliveTimeout: getEnvAsInt("KEEP_ALIVE_TIMEOUT", 300),
	}

	// Validate required environment variables
	if config.JWTSecret == "default-secret-change-in-production" {
		log.Println("WARNING: Using default JWT secret. Please set JWT_SECRET in production.")
	}

	validateConfig(config)

	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func validateConfig(config Config) {
	// Validate port
	if port, err := strconv.Atoi(config.Port); err != nil || port < 1 || port > 65535 {
		log.Printf("WARNING: Invalid port number '%s', using default 8080", config.Port)
		config.Port = "8080"
	}

	// Validate rate limit
	if config.APIRateLimit < 1 {
		log.Printf("WARNING: Invalid API rate limit '%d', using default 100", config.APIRateLimit)
		config.APIRateLimit = 100
	}

	// Validate max request size
	if config.MaxRequestSize < 1024 { // Minimum 1KB
		log.Printf("WARNING: Invalid max request size '%d', using default 10MB", config.MaxRequestSize)
		config.MaxRequestSize = 10 * 1024 * 1024
	}

	log.Printf("Configuration loaded successfully: Port=%s, DB=%s, RateLimit=%d",
		config.Port, config.DBDriver, config.APIRateLimit)
}
