#!/bin/bash

################################################################################
# FRONTEND BUILDER
# Builds powerful React frontend with backend integration
################################################################################

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
FRONTEND_DIR="$PROJECT_ROOT/frontend"

log() { echo -e "${GREEN}[$(date +'%H:%M:%S')]${NC} $1"; }
log_error() { echo -e "${RED}[$(date +'%H:%M:%S')] ERROR:${NC} $1"; }
log_success() { echo -e "${GREEN}[$(date +'%H:%M:%S')] ✓${NC} $1"; }

log "═══════════════════════════════════════════════════════════════"
log "FRONTEND BUILDER - STARTING"
log "═══════════════════════════════════════════════════════════════"

# Create frontend directory
mkdir -p "$FRONTEND_DIR"
cd "$FRONTEND_DIR"

# Initialize if needed
if [ ! -f "package.json" ]; then
    log "Creating new React app..."
    npx create-react-app . --template typescript || {
        log_error "Failed to create React app"
        exit 1
    }
fi

# Install dependencies
log "Installing dependencies..."
npm install axios react-router-dom @tanstack/react-query recharts --legacy-peer-deps

# Build
log "Building frontend..."
npm run build

log_success "Frontend build completed!"
log "Build output: $FRONTEND_DIR/build"
