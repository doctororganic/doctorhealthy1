#!/bin/bash

# Docker deployment script for Nutrition Platform

set -e

echo "🚀 Starting Nutrition Platform Docker Deployment"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

print_status "Docker is running ✓"

# Create network if it doesn't exist
if ! docker network ls | grep -q "nutrition_network"; then
    print_status "Creating Docker network..."
    docker network create nutrition_network
else
    print_status "Docker network already exists ✓"
fi

# Build frontend container
print_status "Building frontend container..."
cd frontend
docker build -t nutrition-frontend:latest .
cd ..

# Run frontend container
print_status "Starting frontend container..."
docker run -d \
    --name nutrition-frontend \
    --network nutrition_network \
    -p 8080:80 \
    --restart unless-stopped \
    nutrition-frontend:latest

print_status "Frontend container started successfully!"

# Check container status
if docker ps | grep -q "nutrition-frontend"; then
    print_status "✅ Frontend container is running"
    print_status "🌐 Access your application at: http://localhost:8080"
else
    print_error "❌ Frontend container failed to start"
    docker logs nutrition-frontend
    exit 1
fi

print_status "🎉 Deployment completed successfully!"
print_status "📊 Container status:"
docker ps --filter "name=nutrition-frontend" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo ""
print_status "Useful commands:"
echo "  • View logs: docker logs nutrition-frontend"
echo "  • Stop container: docker stop nutrition-frontend"
echo "  • Remove container: docker rm nutrition-frontend"
echo "  • Restart container: docker restart nutrition-frontend"