#!/bin/bash

# Migration Runner Script for Nutrition Platform
# This script runs all database migrations in order

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
DB_FILE="nutrition-platform.db"
MIGRATIONS_DIR="migrations"

echo -e "${BLUE}🚀 Nutrition Platform Migration Runner${NC}"
echo "=================================="

# Check if we're in the right directory
if [ ! -f "$DB_FILE" ] && [ ! -d "$MIGRATIONS_DIR" ]; then
    echo -e "${RED}❌ Error: Database file or migrations directory not found${NC}"
    echo "Please run this script from the backend directory"
    exit 1
fi

# Create database file if it doesn't exist
if [ ! -f "$DB_FILE" ]; then
    echo -e "${YELLOW}⚠️  Database file not found. Creating new database...${NC}"
    touch "$DB_FILE"
    echo -e "${GREEN}✅ Database file created: $DB_FILE${NC}"
fi

# Check if migrations directory exists
if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo -e "${RED}❌ Error: Migrations directory '$MIGRATIONS_DIR' not found${NC}"
    exit 1
fi

# Get list of migration files
MIGRATION_FILES=$(find "$MIGRATIONS_DIR" -name "*.sql" -type f | sort -V)

if [ -z "$MIGRATION_FILES" ]; then
    echo -e "${YELLOW}⚠️  No migration files found in $MIGRATIONS_DIR${NC}"
    exit 0
fi

# Create migrations table if it doesn't exist
echo -e "${BLUE}📋 Setting up migrations tracking...${NC}"
sqlite3 "$DB_FILE" "
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
"

# Function to check if migration has been applied
is_migration_applied() {
    local version=$1
    local count=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM schema_migrations WHERE version='$version';")
    [ "$count" -eq 1 ]
}

# Count total migrations
TOTAL_MIGRATIONS=$(echo "$MIGRATION_FILES" | wc -l)
CURRENT_MIGRATION=0

echo -e "${BLUE}📊 Found $TOTAL_MIGRATIONS migration files${NC}"
echo ""

# Run each migration
for migration_file in $MIGRATION_FILES; do
    CURRENT_MIGRATION=$((CURRENT_MIGRATION + 1))
    migration_name=$(basename "$migration_file")
    version="${migration_name%.*}"
    
    echo -e "${BLUE}[${CURRENT_MIGRATION}/${TOTAL_MIGRATIONS}]${NC} Processing: $migration_name"
    
    if is_migration_applied "$version"; then
        echo -e "${YELLOW}   ⏭️  Already applied - skipping${NC}"
        continue
    fi
    
    echo -e "${YELLOW}   🔄 Applying migration...${NC}"
    
    # Run the migration with error handling
    if sqlite3 "$DB_FILE" < "$migration_file"; then
        # Mark migration as applied
        sqlite3 "$DB_FILE" "INSERT INTO schema_migrations (version) VALUES ('$version');"
        echo -e "${GREEN}   ✅ Migration applied successfully${NC}"
    else
        echo -e "${RED}   ❌ Migration failed!${NC}"
        echo -e "${RED}   Error in file: $migration_file${NC}"
        echo -e "${RED}   Please fix the migration and try again${NC}"
        exit 1
    fi
    
    echo ""
done

# Verify database structure
echo -e "${BLUE}🔍 Verifying database structure...${NC}"

# Get list of tables
TABLES=$(sqlite3 "$DB_FILE" ".tables")
echo -e "${GREEN}📋 Tables in database:${NC}"
echo "$TABLES" | tr ' ' '\n' | while read -r table; do
    if [ -n "$table" ]; then
        row_count=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM '$table';" 2>/dev/null || echo "N/A")
        echo -e "   ${GREEN}✓${NC} $table ($row_count rows)"
    fi
done

echo ""
echo -e "${GREEN}🎉 All migrations completed successfully!${NC}"

# Show applied migrations
echo ""
echo -e "${BLUE}📜 Applied migrations:${NC}"
sqlite3 "$DB_FILE" "SELECT version, applied_at FROM schema_migrations ORDER BY applied_at;" | while IFS='|' read -r version applied_at; do
    echo -e "   ${GREEN}✓${NC} $version (applied at $applied_at)"
done

echo ""
echo -e "${GREEN}✅ Database is ready to use!${NC}"
echo ""
echo -e "${BLUE}💡 Next steps:${NC}"
echo "   1. Start the backend server: go run main.go"
echo "   2. Run the frontend: npm run dev"
echo "   3. Test the API endpoints"
echo ""
echo -e "${BLUE}🔧 Useful commands:${NC}"
echo "   View database: sqlite3 $DB_FILE"
echo "   List tables: sqlite3 $DB_FILE '.tables'"
echo "   View schema: sqlite3 $DB_FILE '.schema'"
