#!/bin/bash

################################################################################
# DOCKER COMPOSE GENERATOR
# Generates production-ready docker-compose.yml
################################################################################

set -e

GREEN='\033[0;32m'
NC='\033[0m'

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

log() { echo -e "${GREEN}[$(date +'%H:%M:%S')]${NC} $1"; }

log "Generating docker-compose.production.yml..."

cat > "$PROJECT_ROOT/docker-compose.production.yml" << 'EOF'
version: '3.8'

services:
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_NAME=nutrition_platform
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    depends_on:
      - postgres
      - redis
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=nutrition_platform
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./frontend/build:/usr/share/nginx/html:ro
    depends_on:
      - backend
    restart: unless-stopped

volumes:
  postgres_data:
EOF

log "✓ docker-compose.production.yml created"
