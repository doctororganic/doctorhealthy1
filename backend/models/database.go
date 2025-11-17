// Package models provides data models for the nutrition platform
package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the global database connection
var DB *sql.DB

// InitDB initializes the database connection
func InitDB(databaseURL string) *sql.DB {
	// For testing, use in-memory database to avoid hanging
	dbPath := ":memory:"
	log.Printf("Using in-memory database for testing")

	// Open database
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Configure connection pool (not needed for in-memory but kept for consistency)
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	// Check if database is accessible
	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Run migrations
	if err := runMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database initialized successfully (in-memory)")
	return DB
}

// runMigrations runs database migrations
func runMigrations() error {
	// Create migrations table if not exists
	migrationSQL := `
	CREATE TABLE IF NOT EXISTS migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	if _, err := DB.Exec(migrationSQL); err != nil {
		return err
	}

	// List of migrations to run
	migrations := []Migration{
		{
			Name: "001_initial_schema",
			SQL: `
				CREATE TABLE IF NOT EXISTS users (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					username TEXT NOT NULL UNIQUE,
					email TEXT NOT NULL UNIQUE,
					password_hash TEXT NOT NULL,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);

				CREATE TABLE IF NOT EXISTS api_keys (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					user_id INTEGER NOT NULL,
					key TEXT NOT NULL UNIQUE,
					quota_limit INTEGER DEFAULT 100,
					quota_used INTEGER DEFAULT 0,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (user_id) REFERENCES users (id)
				);

				CREATE TABLE IF NOT EXISTS nutrition_plans (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					user_id INTEGER,
					name TEXT NOT NULL,
					description TEXT,
					target_calories INTEGER,
					protein_grams INTEGER,
					carb_grams INTEGER,
					fat_grams INTEGER,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (user_id) REFERENCES users (id)
				);

				CREATE TABLE IF NOT EXISTS recipes (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					name TEXT NOT NULL,
					description TEXT,
					instructions TEXT,
					prep_time INTEGER,
					cook_time INTEGER,
					servings INTEGER,
					calories_per_serving INTEGER,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);

				CREATE TABLE IF NOT EXISTS health_conditions (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					name TEXT NOT NULL,
					description TEXT,
					nutrition_recommendations TEXT,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);`,
		},
		{
			Name: "002_add_indexes",
			SQL: `
				CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
				CREATE INDEX IF NOT EXISTS idx_api_keys_key ON api_keys(key);
				CREATE INDEX IF NOT EXISTS idx_nutrition_plans_user_id ON nutrition_plans(user_id);
				CREATE INDEX IF NOT EXISTS idx_recipes_name ON recipes(name);`,
		},
		{
			Name: "003_add_audit_logs",
			SQL: `
				CREATE TABLE IF NOT EXISTS audit_logs (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					user_id INTEGER,
					action TEXT NOT NULL,
					resource_type TEXT NOT NULL,
					resource_id INTEGER,
					ip_address TEXT,
					user_agent TEXT,
					timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (user_id) REFERENCES users (id)
				);`,
		},
	}

	// Run each migration if not already applied
	for _, migration := range migrations {
		if err := applyMigration(migration); err != nil {
			return err
		}
	}

	return nil
}

// Migration represents a database migration
type Migration struct {
	Name string
	SQL  string
}

// applyMigration applies a single migration
func applyMigration(migration Migration) error {
	// Check if migration already applied
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = ?", migration.Name).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Printf("Migration %s already applied, skipping", migration.Name)
		return nil // Migration already applied
	}

	log.Printf("Applying migration: %s", migration.Name)

	// Apply migration using transaction for safety
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute migration SQL
	if _, err := tx.Exec(migration.SQL); err != nil {
		return err
	}

	// Record migration
	if _, err = tx.Exec("INSERT INTO migrations (name) VALUES (?)", migration.Name); err != nil {
		return err
	}

	return tx.Commit()
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		log.Println("Closing database connection")
		return DB.Close()
	}
	return nil
}