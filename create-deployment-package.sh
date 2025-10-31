#!/bin/bash

# Create Deployment Package Script
# Creates a ZIP file for Coolify deployment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Function to print colored output
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Create deployment directory
log "📦 Creating deployment package..."
rm -rf nutrition-platform-coolify
mkdir -p nutrition-platform-coolify

# Copy backend files
log "📁 Copying backend files..."
cp -r backend/* nutrition-platform-coolify/

# Copy frontend files
log "📁 Copying frontend files..."
cp -r frontend/* nutrition-platform-coolify/

# Copy configuration files
log "⚙️ Copying configuration files..."
cp docker-compose.yml nutrition-platform-coolify/
cp Dockerfile* nutrition-platform-coolify/
cp .env.production nutrition-platform-coolify/
cp nginx.conf nutrition-platform-coolify/

# Copy documentation
log "📚 Copying documentation..."
cp README.md nutrition-platform-coolify/
cp -r docs nutrition-platform-coolify/ 2>/dev/null || true

# Create deployment info
cat > nutrition-platform-coolify/DEPLOYMENT-INFO.md << 'EOF'
# Nutrition Platform Deployment Information
# Generated: 2025-10-13

## Application Details
- Name: nutrition-platform-secure
- Description: AI-powered nutrition platform with enterprise security
- Version: 1.0.0

## Deployment Configuration
- Build Pack: Dockerfile
- Dockerfile Location: backend/Dockerfile
- Build Context: ./
- Start Command: (will be auto-detected)
- Port: 8080

## Environment Variables
All environment variables are in .env.production file

## Services Required
- PostgreSQL 15
- Redis 7-alpine

## Domain
- Primary: super.doctorhealthy1.com
- Secondary: my.doctorhealthy1.com

## Health Check
- Path: /health
- Interval: 30s
EOF

# Create ZIP file
log "📦 Creating ZIP package..."
cd nutrition-platform-coolify
zip -r ../nutrition-platform-coolify-$(date +%Y%m%d-%H%M%S).zip .
cd ..

success "✅ Deployment package created successfully!"

# List created files
echo ""
echo "📋 Created Files:"
ls -la nutrition-platform-coolify*.zip

echo ""
echo "🚀 Ready for Coolify deployment!"
echo "📦 Upload the ZIP file to Coolify following the deployment steps."