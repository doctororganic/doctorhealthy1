#!/bin/bash

# Pre-deployment checks
set -e

echo "Running pre-deployment checks..."

# Check Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed"
    exit 1
fi
echo "✅ Docker installed"

# Check Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed"
    exit 1
fi
echo "✅ Docker Compose installed"

# Check .env file exists
if [ ! -f ".env" ]; then
    echo "❌ .env file not found"
    exit 1
fi
echo "✅ .env file exists"

# Check required environment variables
required_vars=("DB_PASSWORD" "JWT_SECRET" "API_KEY_SECRET")
for var in "${required_vars[@]}"; do
    if ! grep -q "^${var}=" .env; then
        echo "❌ Missing required variable: $var"
        exit 1
    fi
done
echo "✅ All required environment variables set"

# Check ports are available
ports=(80 443 8080 3000 5432 6379)
for port in "${ports[@]}"; do
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
        echo "⚠️  Port $port is already in use"
    fi
done

echo "✅ Pre-deployment checks passed"
