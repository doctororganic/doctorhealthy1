#!/bin/bash

# 🚀 Quick Deploy Script for Nutrition Platform
# This script deploys the nutrition platform to production using Docker

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Function to check if running as root
check_root() {
    if [[ $EUID -eq 0 ]]; then
        print_error "Please do not run this script as root"
        exit 1
    fi
}

# Function to check if Docker is installed
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    if ! docker info &> /dev/null; then
        print_error "Docker is not running. Please start Docker first."
        exit 1
    fi
}

# Function to check if Docker Compose is installed
check_docker_compose() {
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
}

# Function to check if .env.local file exists
check_env_file() {
    if [[ ! -f .env.local ]]; then
        print_error ".env.local file not found. Please create it first."
        exit 1
    fi
}

# Function to stop existing containers
stop_containers() {
    print_status "Stopping existing containers..."
    docker-compose -f docker-compose.production.yml down 2>/dev/null || true
}

# Function to remove old containers and images
cleanup() {
    print_status "Cleaning up old containers and images..."
    docker system prune -f
}

# Function to build Docker images
build_images() {
    print_status "Building Docker images..."
    docker-compose -f docker-compose.production.yml build
}

# Function to start containers
start_containers() {
    print_status "Starting containers..."
    docker-compose -f docker-compose.production.yml up -d
}

# Function to wait for services to be ready
wait_for_services() {
    print_status "Waiting for services to be ready..."
    
    # Wait for backend to be ready
    local max_attempts=30
    local attempt=1
    
    while [[ $attempt -le $max_attempts ]]; do
        if curl -f http://localhost:8080/api/health &>/dev/null; then
            print_status "Backend is ready!"
            break
        fi
        
        print_warning "Backend not ready yet (attempt $attempt/$max_attempts)..."
        sleep 5
        ((attempt++))
    done
    
    if [[ $attempt -gt $max_attempts ]]; then
        print_error "Backend failed to start within expected time"
        return 1
    fi
    
    # Wait for frontend to be ready
    attempt=1
    while [[ $attempt -le $max_attempts ]]; do
        if curl -f http://localhost:3000 &>/dev/null; then
            print_status "Frontend is ready!"
            break
        fi
        
        print_warning "Frontend not ready yet (attempt $attempt/$max_attempts)..."
        sleep 5
        ((attempt++))
    done
    
    if [[ $attempt -gt $max_attempts ]]; then
        print_error "Frontend failed to start within expected time"
        return 1
    fi
}

# Function to run health checks
run_health_checks() {
    print_status "Running health checks..."
    
    # Check backend health
    echo "Backend health check:"
    curl -f http://localhost:8080/api/health || {
        print_error "Backend health check failed"
        return 1
    }
    
    # Check frontend health
    echo "Frontend health check:"
    curl -f http://localhost:3000 || {
        print_error "Frontend health check failed"
        return 1
    }
    
    print_status "All health checks passed!"
}

# Function to show deployment information
show_deployment_info() {
    print_status "Deployment completed successfully!"
    echo ""
    echo "🌐 Application URLs:"
    echo "  Frontend: http://localhost:3000"
    echo "  Backend:  http://localhost:8080"
    echo "  API Documentation: http://localhost:8080/docs"
    echo ""
    echo "📊 Monitoring:"
    echo "  Health Check: http://localhost:8080/api/health"
    echo "  Monitoring Dashboard: http://localhost:3000/monitoring"
    echo ""
    echo "🛠 Management Commands:"
    echo "  View logs: docker-compose -f docker-compose.production.yml logs -f"
    echo "  Stop services: docker-compose -f docker-compose.production.yml down"
    echo "  Restart services: docker-compose -f docker-compose.production.yml restart"
    echo "  Update application: git pull && docker-compose -f docker-compose.production.yml up -d --build"
    echo ""
}

# Main deployment function
deploy() {
    print_status "🚀 Starting deployment of Nutrition Platform..."
    
    # Check prerequisites
    check_root
    check_docker
    check_docker_compose
    check_env_file
    
    # Stop existing containers
    stop_containers
    
    # Clean up
    cleanup
    
    # Build images
    build_images
    
    # Start containers
    start_containers
    
    # Wait for services
    wait_for_services
    
    # Run health checks
    run_health_checks
    
    # Show deployment information
    show_deployment_info
}

# Function to stop the application
stop() {
    print_status "🛑 Stopping Nutrition Platform..."
    stop_containers
    print_status "Application stopped successfully!"
}

# Function to restart the application
restart() {
    print_status "🔄 Restarting Nutrition Platform..."
    stop_containers
    start_containers
    wait_for_services
    run_health_checks
    show_deployment_info
}

# Function to show logs
logs() {
    print_status "📋 Showing logs for Nutrition Platform..."
    docker-compose -f docker-compose.production.yml logs -f
}

# Function to update the application
update() {
    print_status "🔄 Updating Nutrition Platform..."
    git pull
    build_images
    stop_containers
    start_containers
    wait_for_services
    run_health_checks
    show_deployment_info
}

# Function to backup data
backup() {
    print_status "💾 Backing up data..."
    
    # Create backup directory
    mkdir -p backups/$(date +%Y-%m-%d)
    
    # Backup database
    docker exec nutrition_postgres pg_dump -U nutrition_user nutrition_platform > backups/$(date +%Y-%m-%d)/database.sql
    
    # Backup Redis
    docker exec nutrition_redis redis-cli --rdb > backups/$(date +%Y-%m-%d)/redis.rdb
    
    print_status "Backup completed successfully!"
}

# Function to show status
status() {
    print_status "📊 Nutrition Platform Status:"
    echo ""
    
    # Show container status
    echo "Container Status:"
    docker-compose -f docker-compose.production.yml ps
    
    echo ""
    
    # Show resource usage
    echo "Resource Usage:"
    docker stats --no-stream
    
    echo ""
    
    # Show health status
    echo "Health Status:"
    curl -f http://localhost:8080/api/health 2>/dev/null && echo "✅ Healthy" || echo "❌ Unhealthy"
    curl -f http://localhost:3000 2>/dev/null && echo "✅ Healthy" || echo "❌ Unhealthy"
}

# Function to show help
help() {
    echo "Nutrition Platform Deployment Script"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  deploy     Deploy the application"
    echo "  stop       Stop the application"
    echo "  restart    Restart the application"
    echo "  logs       Show logs"
    echo "  update     Update the application"
    echo "  backup     Backup data"
    echo "  status     Show status"
    echo "  help       Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 deploy     # Deploy the application"
    echo "  $0 stop       # Stop the application"
    echo "  $0 restart    # Restart the application"
    echo "  $0 logs       # Show logs"
    echo "  $0 update     # Update the application"
    echo "  $0 backup     # Backup data"
    echo "  $0 status     # Show status"
    echo "  $0 help       # Show this help message"
}

# Parse command line arguments
case "${1:-}" in
    deploy)
        deploy
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    logs)
        logs
        ;;
    update)
        update
        ;;
    backup)
        backup
        ;;
    status)
        status
        ;;
    help|--help|-h)
        help
        ;;
    *)
        print_error "Unknown command: $1"
        help
        exit 1
        ;;
esac

exit 0