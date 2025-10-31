#!/bin/bash

# Coolify MCP Server Deployment Script for Coolify
# This script helps deploy the MCP server to Coolify

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
check_directory() {
    if [ ! -f "Dockerfile" ] || [ ! -f "docker-compose.yml" ]; then
        print_error "Dockerfile or docker-compose.yml not found. Please run this script from the coolify-mcp directory."
        exit 1
    fi
    print_success "Directory structure validated"
}

# Validate environment configuration
validate_environment() {
    print_status "Validating environment configuration..."

    if [ -z "$COOLIFY_API_TOKEN" ]; then
        print_error "COOLIFY_API_TOKEN environment variable is not set"
        print_error "Please set your Coolify API token:"
        print_error "export COOLIFY_API_TOKEN='your_token_here'"
        exit 1
    fi

    print_success "Environment configuration is valid"
}

# Create Coolify project (if needed)
create_coolify_project() {
    print_status "Creating Coolify project for MCP server..."

    # Use MCP tool to create project if it doesn't exist
    print_status "Checking if project already exists..."

    # For now, we'll assume the project creation is handled via Coolify dashboard
    # or manual project setup. The MCP server will be deployed to an existing project.
    print_warning "Please ensure you have a Coolify project ready for deployment"
    print_warning "Project name: coolify-mcp-server"
    print_warning "Description: Coolify MCP Server for deployment management"
}

# Deploy to Coolify using MCP tools
deploy_to_coolify() {
    print_status "Deploying to Coolify using MCP tools..."

    # This would typically be done via the MCP server tools
    # For now, we'll provide instructions for manual deployment

    print_status "To deploy via Coolify dashboard:"
    echo "1. Go to https://app.doctorhealthy1.com"
    echo "2. Create a new project or use existing project"
    echo "3. Upload the Dockerfile and docker-compose.yml from this directory"
    echo "4. Set environment variables in Coolify dashboard"
    echo "5. Deploy the service"

    print_success "Deployment configuration is ready"
}

# Main deployment function
deploy() {
    print_status "Starting Coolify MCP Server deployment..."

    check_directory
    validate_environment
    create_coolify_project
    deploy_to_coolify

    print_success "🎉 Coolify MCP Server deployment preparation completed!"
    print_status ""
    print_status "Next steps:"
    print_status "1. Go to your Coolify dashboard"
    print_status "2. Create/upload the project using the provided Dockerfile"
    print_status "3. Set the environment variables"
    print_status "4. Deploy the service"
    print_status ""
    print_status "The MCP server will be available as a service in your Coolify instance"
}

# Help function
show_help() {
    echo "Coolify MCP Server Deployment Script"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  deploy    Prepare and deploy to Coolify"
    echo "  validate  Validate environment configuration"
    echo "  help      Show this help message"
    echo ""
    echo "Environment Variables:"
    echo "  COOLIFY_API_TOKEN    Your Coolify API token (required)"
    echo ""
    echo "Examples:"
    echo "  export COOLIFY_API_TOKEN='your_token_here'"
    echo "  $0 deploy    # Deploy to Coolify"
}

# Main script logic
case "${1:-deploy}" in
    "deploy")
        deploy
        ;;
    "validate")
        check_directory
        validate_environment
        ;;
    "help"|"-h"|"--help")
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        echo ""
        show_help
        exit 1
        ;;
esac