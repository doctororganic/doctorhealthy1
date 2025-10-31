package main

import (
	"log"
	"nutrition-platform/migrations"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	// Initialize database and run migrations on startup
	db, err := gorm.Open(sqlite.Open("./nutrition_platform.db"), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}

	// Run migrations
	mm, err := migrations.NewMigrationManager(db, "./migrations", "./backups")
	if err != nil {
		log.Printf("Failed to create migration manager: %v", err)
		return
	}

	if err := mm.Migrate(); err != nil {
		log.Printf("Migration failed: %v", err)
		return
	}

	log.Println("Database initialized successfully")
}