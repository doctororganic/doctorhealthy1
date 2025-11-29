#!/bin/bash

# Database Test Setup Script
# Sets up test database and runs migrations for testing

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
TEST_DB_PATH="$PROJECT_ROOT/test_nutrition.db"
MIGRATIONS_DIR="$PROJECT_ROOT/migrations"

# Logging functions
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if SQLite is available
check_sqlite() {
    if ! command -v sqlite3 >/dev/null 2>&1; then
        log_error "SQLite3 is not installed. Please install it first."
        exit 1
    fi
    log_success "SQLite3 is available"
}

# Clean up existing test database
cleanup_test_db() {
    if [ -f "$TEST_DB_PATH" ]; then
        log "Cleaning up existing test database..."
        rm -f "$TEST_DB_PATH"
        log_success "Test database cleaned up"
    fi
}

# Create test database
create_test_db() {
    log "Creating test database..."
    
    # Create database file
    touch "$TEST_DB_PATH"
    chmod 666 "$TEST_DB_PATH"
    
    # Verify database creation
    if [ -f "$TEST_DB_PATH" ]; then
        log_success "Test database created at: $TEST_DB_PATH"
    else
        log_error "Failed to create test database"
        exit 1
    fi
}

# Run migrations on test database
run_migrations() {
    log "Running migrations on test database..."
    
    cd "$PROJECT_ROOT"
    
    # Set environment for test database
    export DATABASE_URL="sqlite://$TEST_DB_PATH"
    
    # Run migration script
    if [ -f "$SCRIPT_DIR/run_migrations.sh" ]; then
        "$SCRIPT_DIR/run_migrations.sh"
        log_success "Migrations completed"
    else
        log_warning "Migration script not found, creating basic schema..."
        
        # Create basic schema if migration script doesn't exist
        sqlite3 "$TEST_DB_PATH" << EOF
-- Users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Nutrition goals table
CREATE TABLE IF NOT EXISTS nutrition_goals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    daily_calories INTEGER,
    daily_protein REAL,
    daily_carbs REAL,
    daily_fats REAL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Weight logs table
CREATE TABLE IF NOT EXISTS weight_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    weight REAL NOT NULL,
    body_fat_percentage REAL,
    notes TEXT,
    log_date DATE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Body measurements table
CREATE TABLE IF NOT EXISTS body_measurements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    measurement_type TEXT NOT NULL,
    value REAL NOT NULL,
    unit TEXT NOT NULL,
    measurement_date DATE NOT NULL,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Exercises table
CREATE TABLE IF NOT EXISTS exercises (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    description TEXT,
    instructions TEXT,
    muscle_groups TEXT,
    equipment TEXT,
    difficulty TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Workouts table
CREATE TABLE IF NOT EXISTS workouts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    duration_minutes INTEGER,
    difficulty TEXT,
    muscle_groups TEXT,
    exercises TEXT, -- JSON array of exercise IDs and reps/sets
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Nutrition data tables
CREATE TABLE IF NOT EXISTS recipes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    ingredients TEXT, -- JSON
    instructions TEXT,
    prep_time_minutes INTEGER,
    cook_time_minutes INTEGER,
    servings INTEGER,
    calories_per_serving INTEGER,
    protein_per_serving REAL,
    carbs_per_serving REAL,
    fats_per_serving REAL,
    fiber_per_serving REAL,
    sugar_per_serving REAL,
    sodium_per_serving REAL,
    cuisine_type TEXT,
    meal_type TEXT,
    dietary_restrictions TEXT, -- JSON array
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS meals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    foods TEXT, -- JSON
    total_calories INTEGER,
    total_protein REAL,
    total_carbs REAL,
    total_fats REAL,
    meal_type TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Health conditions table
CREATE TABLE IF NOT EXISTS health_conditions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    symptoms TEXT, -- JSON
    dietary_recommendations TEXT, -- JSON
    foods_to_avoid TEXT, -- JSON
    recommended_foods TEXT, -- JSON
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Vitamins and minerals table
CREATE TABLE IF NOT EXISTS vitamins_minerals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL, -- 'vitamin' or 'mineral'
    description TEXT,
    daily_recommended_amount REAL,
    unit TEXT,
    food_sources TEXT, -- JSON
    deficiency_symptoms TEXT, -- JSON
    benefits TEXT, -- JSON
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Injuries table
CREATE TABLE IF NOT EXISTS injuries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    body_part TEXT NOT NULL,
    injury_type TEXT NOT NULL,
    severity TEXT,
    recovery_time_weeks INTEGER,
    exercises_to_avoid TEXT, -- JSON
    recommended_exercises TEXT, -- JSON
    nutritional_recommendations TEXT, -- JSON
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample test data
INSERT OR IGNORE INTO users (id, email, password_hash, name) VALUES 
    (1, 'test@example.com', 'hashed_password', 'Test User'),
    (2, 'admin@example.com', 'hashed_admin_password', 'Admin User');

INSERT OR IGNORE INTO nutrition_goals (user_id, daily_calories, daily_protein, daily_carbs, daily_fats) VALUES
    (1, 2000, 150, 250, 65),
    (2, 2500, 180, 300, 80);

INSERT OR IGNORE INTO weight_logs (user_id, weight, log_date) VALUES
    (1, 70.5, '2023-01-01'),
    (1, 70.2, '2023-01-08'),
    (2, 85.0, '2023-01-01');

INSERT OR IGNORE INTO exercises (id, name, category, description, muscle_groups, difficulty) VALUES
    (1, 'Push-ups', 'strength', 'Classic upper body exercise', 'chest, shoulders, triceps', 'beginner'),
    (2, 'Squats', 'strength', 'Lower body compound exercise', 'quads, glutes, hamstrings', 'beginner'),
    (3, 'Running', 'cardio', 'Cardiovascular exercise', 'legs, core', 'intermediate');

INSERT OR IGNORE INTO workouts (id, name, description, duration_minutes, difficulty, muscle_groups, exercises) VALUES
    (1, 'Full Body Workout', 'Complete workout for all major muscle groups', 45, 'intermediate', 'full_body', '[{"exercise_id": 1, "sets": 3, "reps": 15}, {"exercise_id": 2, "sets": 3, "reps": 20}]');

INSERT OR IGNORE INTO recipes (id, name, description, calories_per_serving, protein_per_serving, carbs_per_serving, fats_per_serving, cuisine_type, meal_type) VALUES
    (1, 'Grilled Chicken Salad', 'Healthy salad with grilled chicken', 350, 30, 15, 20, 'american', 'lunch'),
    (2, 'Vegetable Stir Fry', 'Mixed vegetables with rice', 400, 15, 50, 15, 'asian', 'dinner');

EOF
        
        log_success "Basic test schema created"
    fi
}

# Verify database setup
verify_setup() {
    log "Verifying database setup..."
    
    # Check if tables exist
    tables=$(sqlite3 "$TEST_DB_PATH" "SELECT name FROM sqlite_master WHERE type='table';")
    
    if [ -z "$tables" ]; then
        log_error "No tables found in test database"
        exit 1
    fi
    
    table_count=$(echo "$tables" | wc -l)
    log_success "Found $table_count tables in test database"
    
    # Check if sample data exists
    user_count=$(sqlite3 "$TEST_DB_PATH" "SELECT COUNT(*) FROM users;")
    if [ "$user_count" -gt 0 ]; then
        log_success "Sample data inserted successfully ($user_count users)"
    else
        log_warning "No sample data found"
    fi
}

# Create test environment file
create_test_env() {
    log "Creating test environment file..."
    
    cat > "$PROJECT_ROOT/.env.test" << EOF
# Test Environment Configuration
DATABASE_URL=sqlite://$TEST_DB_PATH
REDIS_URL=redis://localhost:6379/1
JWT_SECRET=test-secret-key-for-development-only
CORS_ORIGINS=http://localhost:3000,http://localhost:3001
LOG_LEVEL=debug
CACHE_TTL=300
RATE_LIMIT_REQUESTS=1000
ENVIRONMENT=test
PORT=8080
EOF
    
    log_success "Test environment file created: .env.test"
}

# Show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -c, --cleanup     Clean up existing test database only"
    echo "  -s, --setup       Setup test database only"
    echo "  -v, --verify      Verify setup only"
    echo "  -a, --all         Clean, setup, and verify (default)"
    echo "  -h, --help       Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                # Full setup"
    echo "  $0 --cleanup      # Clean up only"
    echo "  $0 --verify       # Verify existing setup"
}

# Main function
main() {
    local cleanup_only=false
    local setup_only=false
    local verify_only=false
    local run_all=true
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -c|--cleanup)
                cleanup_only=true
                run_all=false
                shift
                ;;
            -s|--setup)
                setup_only=true
                run_all=false
                shift
                ;;
            -v|--verify)
                verify_only=true
                run_all=false
                shift
                ;;
            -a|--all)
                run_all=true
                shift
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    log "Database Test Setup for Nutrition Platform"
    
    # Check dependencies
    check_sqlite
    
    # Execute based on flags
    if [ "$run_all" = true ] || [ "$cleanup_only" = true ]; then
        cleanup_test_db
    fi
    
    if [ "$run_all" = true ] || [ "$setup_only" = true ]; then
        create_test_db
        run_migrations
        create_test_env
    fi
    
    if [ "$run_all" = true ] || [ "$verify_only" = true ]; then
        verify_setup
    fi
    
    log_success "Database test setup completed!"
    echo ""
    echo "Test database location: $TEST_DB_PATH"
    echo "To use for testing: export DATABASE_URL=sqlite://$TEST_DB_PATH"
}

# Run main function
main "$@"
