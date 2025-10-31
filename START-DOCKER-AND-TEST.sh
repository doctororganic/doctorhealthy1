#!/bin/bash

echo "🐳 Starting Docker Desktop and Running Tests"
echo "============================================="

# Check if Docker Desktop is running
if ! docker info &> /dev/null; then
    echo "⚠️  Docker Desktop is not running"
    echo ""
    echo "Please start Docker Desktop manually:"
    echo "  1. Open Docker Desktop application"
    echo "  2. Wait for it to start (whale icon in menu bar)"
    echo "  3. Run this script again"
    echo ""
    echo "Or run: open -a Docker"
    exit 1
fi

echo "✅ Docker is running"
echo ""

# Run the live test and deploy
./LIVE-TEST-DEPLOY-NOW.sh
