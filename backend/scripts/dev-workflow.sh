#!/bin/bash
# Development Workflow Helper
# Usage: ./scripts/dev-workflow.sh [command]

case "$1" in
    "start")
        echo "🚀 Starting development server..."
        go run main.go
        ;;
    "test-all")
        echo "🧪 Running all tests..."
        go test -v ./...
        go test -v ./tests/...
        ;;
    "check")
        echo "🔍 Running checks..."
        go fmt ./...
        go vet ./...
        go build -o /dev/null ./main.go
        echo "✅ All checks passed!"
        ;;
    "reset-db")
        echo "🗑️  Resetting database..."
        rm -f nutrition-platform.db
        ./run_migrations.sh
        echo "✅ Database reset!"
        ;;
    "logs")
        echo "📋 Showing logs..."
        tail -f logs/app.log 2>/dev/null || echo "No log file found"
        ;;
    *)
        echo "Available commands:"
        echo "  start      - Start development server"
        echo "  test-all   - Run all tests"
        echo "  check      - Format, vet, and build"
        echo "  reset-db   - Reset database"
        echo "  logs       - Show application logs"
        ;;
esac