package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"nutrition-platform/config"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// DB holds the database connection
var DB *sql.DB

// InitDatabase initializes the database connection
func InitDatabase(cfg *config.Config) error {
	var err error

	// Use SQLite for development, PostgreSQL for production
	if cfg.Environment == "development" || cfg.Environment == "test" {
		// SQLite connection for development and testing
		DB, err = sql.Open("sqlite3", "./nutrition_platform.db")
	} else {
		// PostgreSQL connection for production
		connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
		DB, err = sql.Open("postgres", connStr)
	}

	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	DB.SetMaxOpenConns(cfg.DBMaxConnections)
	DB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	DB.SetConnMaxLifetime(cfg.DBConnLifetime)

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")

	// Create tables if they don't exist
	if err = createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

// createTables creates the necessary database tables
func createTables() error {
	queries := []string{
		// Users table
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			first_name VARCHAR(50),
			last_name VARCHAR(50),
			date_of_birth DATE,
			gender VARCHAR(10),
			height DECIMAL(5,2),
			weight DECIMAL(5,2),
			activity_level VARCHAR(20),
			goals TEXT,
			dietary_restrictions TEXT[],
			religious_filter_enabled BOOLEAN DEFAULT true,
			filter_alcohol BOOLEAN DEFAULT true,
			filter_pork BOOLEAN DEFAULT true,
			preferred_language VARCHAR(5) DEFAULT 'en',
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// Foods table
		`CREATE TABLE IF NOT EXISTS foods (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			name_ar VARCHAR(100),
			description TEXT,
			description_ar TEXT,
			category VARCHAR(50),
			subcategory VARCHAR(50),
			barcode VARCHAR(50),
			brand VARCHAR(100),
			serving_size DECIMAL(8,2),
			serving_unit VARCHAR(20),
			calories_per_100g DECIMAL(8,2),
			protein_per_100g DECIMAL(8,2),
			carbs_per_100g DECIMAL(8,2),
			fat_per_100g DECIMAL(8,2),
			fiber_per_100g DECIMAL(8,2),
			sugar_per_100g DECIMAL(8,2),
			sodium_per_100g DECIMAL(8,2),
			ingredients TEXT[],
			allergens TEXT[],
			contains_alcohol BOOLEAN DEFAULT false,
			contains_pork BOOLEAN DEFAULT false,
			is_halal BOOLEAN DEFAULT true,
			is_kosher BOOLEAN DEFAULT false,
			is_vegetarian BOOLEAN DEFAULT false,
			is_vegan BOOLEAN DEFAULT false,
			image_url VARCHAR(255),
			verified BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// Exercises table
		`CREATE TABLE IF NOT EXISTS exercises (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			name_ar VARCHAR(100),
			description TEXT,
			description_ar TEXT,
			category VARCHAR(50),
			muscle_groups TEXT[],
			equipment VARCHAR(50),
			difficulty_level VARCHAR(20),
			instructions TEXT,
			instructions_ar TEXT,
			calories_per_minute DECIMAL(5,2),
			met_value DECIMAL(4,2),
			image_url VARCHAR(255),
			video_url VARCHAR(255),
			verified BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// Meal plans table
		`CREATE TABLE IF NOT EXISTS meal_plans (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			start_date DATE,
			end_date DATE,
			total_calories INTEGER,
			total_protein DECIMAL(8,2),
			total_carbs DECIMAL(8,2),
			total_fat DECIMAL(8,2),
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// Workout plans table
		`CREATE TABLE IF NOT EXISTS workout_plans (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			start_date DATE,
			end_date DATE,
			difficulty_level VARCHAR(20),
			goal VARCHAR(50),
			days_per_week INTEGER,
			duration_weeks INTEGER,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// User food logs table
		`CREATE TABLE IF NOT EXISTS user_food_logs (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			food_id INTEGER REFERENCES foods(id),
			quantity DECIMAL(8,2) NOT NULL,
			unit VARCHAR(20),
			meal_type VARCHAR(20),
			consumed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			calories DECIMAL(8,2),
			protein DECIMAL(8,2),
			carbs DECIMAL(8,2),
			fat DECIMAL(8,2),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// User exercise logs table
		`CREATE TABLE IF NOT EXISTS user_exercise_logs (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			exercise_id INTEGER REFERENCES exercises(id),
			duration_minutes INTEGER,
			sets INTEGER,
			reps INTEGER,
			weight DECIMAL(6,2),
			calories_burned DECIMAL(8,2),
			notes TEXT,
			performed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// API keys table
		`CREATE TABLE IF NOT EXISTS api_keys (
			id VARCHAR(50) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			key_hash VARCHAR(255) UNIQUE NOT NULL,
			prefix VARCHAR(10) NOT NULL,
			user_id VARCHAR(50),
			status VARCHAR(20) DEFAULT 'active',
			scopes JSONB NOT NULL,
			rate_limit INTEGER DEFAULT 100,
			expires_at TIMESTAMP,
			last_used_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			metadata JSONB DEFAULT '{}'
		);`,

		// API key usage table
		`CREATE TABLE IF NOT EXISTS api_key_usage (
			id SERIAL PRIMARY KEY,
			api_key_id VARCHAR(50) REFERENCES api_keys(id) ON DELETE CASCADE,
			endpoint VARCHAR(255) NOT NULL,
			method VARCHAR(10) NOT NULL,
			status_code INTEGER NOT NULL,
			response_time BIGINT NOT NULL,
			ip_address INET,
			user_agent TEXT,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// API metrics snapshots table
		`CREATE TABLE IF NOT EXISTS api_metrics_snapshots (
			id SERIAL PRIMARY KEY,
			api_key_id VARCHAR(50) REFERENCES api_keys(id) ON DELETE CASCADE,
			metrics_data JSONB NOT NULL,
			timestamp TIMESTAMP NOT NULL,
			UNIQUE(api_key_id, timestamp)
		);`,

		// Usage alerts table
		`CREATE TABLE IF NOT EXISTS usage_alerts (
			id VARCHAR(50) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			condition VARCHAR(100) NOT NULL,
			threshold DECIMAL(10,2) NOT NULL,
			enabled BOOLEAN DEFAULT true,
			last_triggered TIMESTAMP,
			metadata JSONB DEFAULT '{}',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// System metrics table
		`CREATE TABLE IF NOT EXISTS system_metrics (
			id SERIAL PRIMARY KEY,
			metric_name VARCHAR(100) NOT NULL,
			metric_value DECIMAL(15,6),
			metric_type VARCHAR(50),
			labels JSONB,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	log.Println("Database tables created successfully")
	return nil
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// HealthCheck performs a database health check
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return DB.PingContext(ctx)
}
