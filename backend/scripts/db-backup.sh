#!/bin/bash
# Database Backup Script
# Usage: ./scripts/db-backup.sh [restore backup_file.db]

DB_FILE="nutrition-platform.db"
BACKUP_DIR="backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

if [ "$1" = "restore" ] && [ -n "$2" ]; then
    echo "🔄 Restoring database from $2..."
    cp "$2" "$DB_FILE"
    echo "✅ Database restored!"
    exit 0
fi

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Create backup
BACKUP_FILE="$BACKUP_DIR/backup_$TIMESTAMP.db"
cp "$DB_FILE" "$BACKUP_FILE"

# Compress backup
gzip "$BACKUP_FILE"
BACKUP_FILE="${BACKUP_FILE}.gz"

echo "✅ Backup created: $BACKUP_FILE"

# Keep only last 10 backups
ls -t "$BACKUP_DIR"/backup_*.db.gz | tail -n +11 | xargs rm -f

echo "✅ Old backups cleaned (keeping last 10)"