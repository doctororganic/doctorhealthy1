#!/bin/bash

# Coolify MCP Server Deployment Script
# This script helps deploy and manage the Coolify MCP server

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
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

# Check if Node.js is installed
check_nodejs() {
    if ! command -v node &> /dev/null; then
        print_error "Node.js is not installed. Please install Node.js 18+ to continue."
        exit 1
    fi

    local version=$(node -v | cut -d 'v' -f 2 | cut -d '.' -f 1)
    if [ "$version" -lt 18 ]; then
        print_error "Node.js version 18+ is required. Current version: $(node -v)"
        exit 1
    fi

    print_success "Node.js $(node -v) is installed"
}

# Check if npm is installed
check_npm() {
    if ! command -v npm &> /dev/null; then
        print_error "npm is not installed. Please install npm to continue."
        exit 1
    fi

    print_success "npm $(npm -v) is installed"
}

# Install dependencies
install_dependencies() {
    print_status "Installing dependencies..."

    if [ ! -d "node_modules" ]; then
        npm install
        print_success "Dependencies installed successfully"
    else
        print_status "Dependencies already installed"
    fi
}

# Build the project
build_project() {
    print_status "Building the project..."

    npm run build

    if [ $? -eq 0 ]; then
        print_success "Project built successfully"
    else
        print_error "Failed to build the project"
        exit 1
    fi
}

# Create .env file if it doesn't exist
setup_environment() {
    if [ ! -f ".env" ]; then
        print_warning ".env file not found. Creating template..."

        cat > .env << EOL
# Coolify API Configuration
COOLIFY_API_BASE_URL=https://api.doctorhealthy1.com/api/v1
COOLIFY_API_TOKEN=your_api_token_here

# Instructions:
# 1. Replace 'your_api_token_here' with your actual Coolify API token
# 2. Update the API base URL if using a different Coolify instance
# 3. Never commit this file to version control
EOL

        print_warning "Please edit the .env file and add your Coolify API token before running the server"
        exit 1
    else
        print_success "Environment file already exists"
    fi
}

# Validate environment configuration
validate_environment() {
    print_status "Validating environment configuration..."

    if ! command -v dotenv &> /dev/null; then
        npm install dotenv --save
    fi

    # Check if required environment variables are set
    if [ -z "$COOLIFY_API_BASE_URL" ] || [ -z "$COOLIFY_API_TOKEN" ]; then
        print_error "Required environment variables are not set"
        print_error "Please check your .env file and ensure COOLIFY_API_BASE_URL and COOLIFY_API_TOKEN are configured"
        exit 1
    fi

    print_success "Environment configuration is valid"
}

# Test the server
test_server() {
    print_status "Testing the server..."

    # Start server in background for testing
    node dist/index.js &
    SERVER_PID=$!

    # Wait a moment for server to start
    sleep 3

    # Check if server process is still running
    if kill -0 $SERVER_PID 2>/dev/null; then
        print_success "Server started successfully (PID: $SERVER_PID)"

        # Stop the test server
        kill $SERVER_PID
        sleep 2

        if [ $? -eq 0 ]; then
            print_success "Server test completed successfully"
        fi
    else
        print_error "Server failed to start"
        exit 1
    fi
}

# Main deployment function
deploy() {
    print_status "Starting Coolify MCP Server deployment..."

    check_nodejs
    check_npm
    setup_environment
    validate_environment
    install_dependencies
    build_project
    test_server

    print_success "🎉 Coolify MCP Server deployment completed successfully!"
    print_status ""
    print_status "To start the server:"
    print_status "  node dist/index.js"
    print_status ""
    print_status "The server will run as an MCP server and can be used with MCP-compatible clients."
}

# Help function
show_help() {
    echo "Coolify MCP Server Deployment Script"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  deploy    Deploy and configure the Coolify MCP server"
    echo "  build     Build the project only"
    echo "  test      Test the server functionality"
    echo "  help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 deploy    # Full deployment process"
    echo "  $0 build     # Build the project"
    echo "  $0 test      # Test server functionality"
}

# Main script logic
case "${1:-deploy}" in
    "deploy")
        deploy
        ;;
    "build")
        check_nodejs
        check_npm
        install_dependencies
        build_project
        ;;
    "test")
        check_nodejs
        validate_environment
        test_server
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