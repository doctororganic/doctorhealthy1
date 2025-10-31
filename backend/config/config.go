package config

import (
	"os"
	"strconv"
	"strings"
)

// FileStorageConfig holds file storage configuration
type FileStorageConfig struct {
	StorageType string `json:"storage_type"`
	BasePath    string `json:"base_path"`
	BaseURL     string `json:"base_url"`
	S3Bucket    string `json:"s3_bucket"`
	S3Region    string `json:"s3_region"`
	S3URL       string `json:"s3_url"`
}

// AppConfig holds application configuration
type AppConfig struct {
	Port              string             `json:"port"`
	DatabaseURL       string             `json:"database_url"`
	JWTSecret         string             `json:"jwt_secret"`
	Environment       string             `json:"environment"`
	ReadTimeout       int                `json:"read_timeout"`
	WriteTimeout      int                `json:"write_timeout"`
	KeepAliveTimeout  int                `json:"keep_alive_timeout"`
	FileStorage       FileStorageConfig   `json:"file_storage"`
	EmailConfig       EmailConfig        `json:"email"`
	PushConfig        PushConfig         `json:"push"`
}

// EmailConfig holds email service configuration
type EmailConfig struct {
	Provider   string `json:"provider"`
	SMTPHost   string `json:"smtp_host"`
	SMTPPort   int    `json:"smtp_port"`
	SMTPUser   string `json:"smtp_user"`
	SMTPPass   string `json:"smtp_pass"`
	FromEmail  string `json:"from_email"`
	FromName   string `json:"from_name"`
}

// PushConfig holds push notification configuration
type PushConfig struct {
	FCMServerKey string `json:"fcm_server_key"`
	APNSKeyPath  string `json:"apns_key_path"`
	APNSKeyID    string `json:"apns_key_id"`
	APNSTeamID   string `json:"apns_team_id"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *AppConfig {
	config := &AppConfig{
		Port:             getEnv("PORT", "8080"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://localhost/nutrition_platform?sslmode=disable"),
		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		Environment:      getEnv("ENVIRONMENT", "development"),
		ReadTimeout:      getEnvAsInt("READ_TIMEOUT", 30),
		WriteTimeout:     getEnvAsInt("WRITE_TIMEOUT", 30),
		KeepAliveTimeout: getEnvAsInt("KEEP_ALIVE_TIMEOUT", 60),
		FileStorage: FileStorageConfig{
			StorageType: getEnv("STORAGE_TYPE", "local"),
			BasePath:    getEnv("FILE_STORAGE_PATH", "./uploads"),
			BaseURL:     getEnv("FILE_STORAGE_URL", "http://localhost:8080/uploads"),
			S3Bucket:    getEnv("S3_BUCKET", ""),
			S3Region:    getEnv("S3_REGION", "us-east-1"),
			S3URL:       getEnv("S3_URL", ""),
		},
		EmailConfig: EmailConfig{
			Provider:  getEnv("EMAIL_PROVIDER", "smtp"),
			SMTPHost:  getEnv("SMTP_HOST", "localhost"),
			SMTPPort:  getEnvAsInt("SMTP_PORT", 587),
			SMTPUser:  getEnv("SMTP_USER", ""),
			SMTPPass:  getEnv("SMTP_PASS", ""),
			FromEmail: getEnv("FROM_EMAIL", "noreply@nutrition-platform.com"),
			FromName:  getEnv("FROM_NAME", "Nutrition Platform"),
		},
		PushConfig: PushConfig{
			FCMServerKey: getEnv("FCM_SERVER_KEY", ""),
			APNSKeyPath:  getEnv("APNS_KEY_PATH", ""),
			APNSKeyID:    getEnv("APNS_KEY_ID", ""),
			APNSTeamID:   getEnv("APNS_TEAM_ID", ""),
		},
	}

	// Validate required configuration
	if config.JWTSecret == "your-secret-key-change-in-production" && config.Environment == "production" {
		panic("JWT_SECRET must be set in production")
	}

	return config
}

// IsDevelopment checks if the environment is development
func (c *AppConfig) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction checks if the environment is production
func (c *AppConfig) IsProduction() bool {
	return c.Environment == "production"
}

// IsLocalStorage checks if using local file storage
func (c *AppConfig) IsLocalStorage() bool {
	return c.FileStorage.StorageType == "local"
}

// IsS3Storage checks if using S3 file storage
func (c *AppConfig) IsS3Storage() bool {
	return c.FileStorage.StorageType == "s3"
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

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
