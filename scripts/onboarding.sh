#!/bin/bash

# Developer onboarding script for nutrition platform
# This script helps new developers set up their environment and understand the coordination system

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to display colored headers
header() {
    echo -e "${PURPLE}=====================================${NC}"
    echo -e "${PURPLE}$1${NC}"
    echo -e "${PURPLE}=====================================${NC}"
}

# Function to display success messages
success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# Function to display warning messages
warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Function to display error messages
error() {
    echo -e "${RED}❌ $1${NC}"
}

# Function to display info messages
info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Function to prompt for input
prompt() {
    echo -e "${CYAN}📝 $1${NC}"
    read -r response
    echo "$response"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to validate developer ID
validate_dev_id() {
    local dev_id="$1"
    
    if [[ ! "$dev_id" =~ ^DEV[1-5]$ ]]; then
        error "Invalid developer ID. Must be DEV1, DEV2, DEV3, DEV4, or DEV5"
        return 1
    fi
    
    return 0
}

# Main onboarding process
main() {
    header "NUTRITION PLATFORM DEVELOPER ONBOARDING"
    
    echo ""
    echo "Welcome to the Nutrition Platform development team!"
    echo "This script will help you set up your development environment"
    echo "and understand our multi-developer coordination system."
    echo ""
    
    # Step 1: Developer identification
    header "STEP 1: DEVELOPER IDENTIFICATION"
    
    echo "Please identify your role in the team:"
    echo ""
    echo "DEV1 - Testing & Quality Assurance"
    echo "  - Owns: backend/tests/**, middleware/*_test.go, e2e-tests/"
    echo "  - Responsible: Testing infrastructure, CI/CD testing pipeline"
    echo ""
    echo "DEV2 - Frontend Integration"
    echo "  - Owns: frontend-nextjs/src/app/**, components/ui/**, hooks/"
    echo "  - Responsible: Frontend build processes, UI components"
    echo ""
    echo "DEV3 - Backend API & Services"
    echo "  - Owns: handlers/, services/, repositories/, models/"
    echo "  - Responsible: API endpoints, database schema, business logic"
    echo ""
    echo "DEV4 - DevOps & Infrastructure"
    echo "  - Owns: .github/workflows/, Makefile, docker-compose.yml"
    echo "  - Responsible: Deployment pipelines, infrastructure"
    echo ""
    echo "DEV5 - Documentation & API Reference"
    echo "  - Owns: docs/, README.md, CHANGELOG.md"
    echo "  - Responsible: Documentation, API reference"
    echo ""
    
    while true; do
        dev_id=$(prompt "Enter your developer ID (DEV1-DEV5):")
        
        if validate_dev_id "$dev_id"; then
            success "Developer ID set to: $dev_id"
            break
        fi
    done
    
    # Step 2: Environment setup
    header "STEP 2: ENVIRONMENT SETUP"
    
    echo "Checking required tools..."
    
    # Check for Git
    if command_exists git; then
        success "Git is installed"
        git_version=$(git --version)
        info "Version: $git_version"
    else
        error "Git is not installed. Please install Git first."
        exit 1
    fi
    
    # Check for Go (for backend developers)
    if [[ "$dev_id" =~ ^DEV[13]$ ]]; then
        if command_exists go; then
            success "Go is installed"
            go_version=$(go version)
            info "Version: $go_version"
        else
            error "Go is not installed. Please install Go first."
            info "Visit: https://golang.org/dl/"
            exit 1
        fi
    fi
    
    # Check for Node.js (for frontend developers)
    if [[ "$dev_id" =~ ^DEV[25]$ ]]; then
        if command_exists node; then
            success "Node.js is installed"
            node_version=$(node --version)
            info "Version: $node_version"
        else
            error "Node.js is not installed. Please install Node.js first."
            info "Visit: https://nodejs.org/"
            exit 1
        fi
        
        if command_exists npm; then
            success "npm is installed"
            npm_version=$(npm --version)
            info "Version: $npm_version"
        else
            error "npm is not installed. Please install npm first."
            exit 1
        fi
    fi
    
    # Check for Docker (for DevOps)
    if [[ "$dev_id" == "DEV4" ]]; then
        if command_exists docker; then
            success "Docker is installed"
            docker_version=$(docker --version)
            info "Version: $docker_version"
        else
            error "Docker is not installed. Please install Docker first."
            info "Visit: https://docs.docker.com/get-docker/"
            exit 1
        fi
    fi
    
    # Step 3: Git configuration
    header "STEP 3: GIT CONFIGURATION"
    
    echo "Configuring Git for multi-developer coordination..."
    
    # Get current Git configuration
    current_name=$(git config user.name)
    current_email=$(git config user.email)
    
    echo "Current Git configuration:"
    echo "  Name: $current_name"
    echo "  Email: $current_email"
    echo ""
    
    if prompt "Do you want to update your Git configuration? (y/n):" | grep -iq "^y"; then
        new_name=$(prompt "Enter your full name:")
        new_email=$(prompt "Enter your email address:")
        
        git config user.name "$new_name"
        git config user.email "$new_email"
        
        success "Git configuration updated"
        info "Name: $new_name"
        info "Email: $new_email"
    fi
    
    # Step 4: Lock system setup
    header "STEP 4: LOCK SYSTEM SETUP"
    
    echo "Setting up file locking system..."
    
    # Create locks directory
    mkdir -p .locks
    success "Locks directory created"
    
    # Test lock management script
    if [ -f "scripts/manage-locks.sh" ]; then
        success "Lock management script found"
        
        # Test the script
        if ./scripts/manage-locks.sh help >/dev/null 2>&1; then
            success "Lock management script is working"
        else
            error "Lock management script is not working"
            chmod +x scripts/manage-locks.sh
            success "Fixed script permissions"
        fi
    else
        error "Lock management script not found"
        exit 1
    fi
    
    # Step 5: Pre-commit hooks
    header "STEP 5: PRE-COMMIT HOOKS"
    
    echo "Setting up pre-commit hooks..."
    
    if [ -d ".husky" ]; then
        success "Husky directory found"
        
        if [ -f ".husky/pre-commit" ]; then
            success "Pre-commit hook found"
            
            if [ -x ".husky/pre-commit" ]; then
                success "Pre-commit hook is executable"
            else
                chmod +x .husky/pre-commit
                success "Fixed pre-commit hook permissions"
            fi
        else
            error "Pre-commit hook not found"
            exit 1
        fi
    else
        error "Husky directory not found"
        exit 1
    fi
    
    # Step 6: Development environment
    header "STEP 6: DEVELOPMENT ENVIRONMENT"
    
    echo "Setting up development environment..."
    
    # Backend setup for DEV1 and DEV3
    if [[ "$dev_id" =~ ^DEV[13]$ ]]; then
        echo "Setting up backend environment..."
        
        if [ -d "backend" ]; then
            cd backend
            
            # Copy environment file
            if [ -f ".env.example" ] && [ ! -f ".env" ]; then
                cp .env.example .env
                success "Backend environment file created"
            fi
            
            # Download Go modules
            if go mod download; then
                success "Go modules downloaded"
            else
                error "Failed to download Go modules"
                exit 1
            fi
            
            cd ..
        else
            error "Backend directory not found"
            exit 1
        fi
    fi
    
    # Frontend setup for DEV2 and DEV5
    if [[ "$dev_id" =~ ^DEV[25]$ ]]; then
        echo "Setting up frontend environment..."
        
        if [ -d "frontend-nextjs" ]; then
            cd frontend-nextjs
            
            # Copy environment file
            if [ -f ".env.local.example" ] && [ ! -f ".env.local" ]; then
                cp .env.local.example .env.local
                success "Frontend environment file created"
            fi
            
            # Install npm dependencies
            if npm install; then
                success "npm dependencies installed"
            else
                error "Failed to install npm dependencies"
                exit 1
            fi
            
            cd ..
        else
            error "Frontend directory not found"
            exit 1
        fi
    fi
    
    # Step 7: Feature branch creation
    header "STEP 7: FEATURE BRANCH CREATION"
    
    echo "Creating your feature branch..."
    
    branch_name="${dev_id,,}-initial-setup"
    
    if git checkout -b "$branch_name"; then
        success "Feature branch created: $branch_name"
    else
        warning "Failed to create feature branch (might already exist)"
    fi
    
    # Step 8: Coordination guidelines
    header "STEP 8: COORDINATION GUIDELINES"
    
    echo "Key coordination guidelines for $dev_id:"
    echo ""
    
    case "$dev_id" in
        "DEV1")
            echo "- You own all testing infrastructure"
            echo "- Coordinate before modifying shared middleware"
            echo "- Do not modify handler implementations"
            echo "- Create new test files only"
            echo "- Run tests before committing"
            ;;
        "DEV2")
            echo "- You own frontend components and pages"
            echo "- Do not modify backend code"
            echo "- Coordinate before modifying shared components"
            echo "- Use feature branches for new pages"
            echo "- Test frontend before committing"
            ;;
        "DEV3")
            echo "- You own backend API and services"
            echo "- Do not modify frontend code"
            echo "- Coordinate before changing API endpoints"
            echo "- Maintain backward compatibility"
            echo "- Document all API changes"
            ;;
        "DEV4")
            echo "- You own DevOps and infrastructure"
            echo "- Do not modify application code"
            echo "- Coordinate before changing build processes"
            echo "- Test scripts in isolation"
            echo "- Own deployment pipelines"
            ;;
        "DEV5")
            echo "- You own documentation"
            echo "- Do not modify code files"
            echo "- Coordinate before updating architecture docs"
            echo "- Update CHANGELOG.md last"
            echo "- Maintain API reference accuracy"
            ;;
    esac
    
    echo ""
    echo "General rules for all developers:"
    echo "- Always check for locks before modifying shared files"
    echo "- Create feature branches for all work"
    echo "- Make atomic commits"
    echo "- Test before merging"
    echo "- Communicate before modifying high-risk files"
    echo ""
    
    # Step 9: Quick reference
    header "STEP 9: QUICK REFERENCE"
    
    echo "Useful commands:"
    echo ""
    echo "Lock management:"
    echo "  ./scripts/manage-locks.sh lock <file> <dev_id> <description> <eta>"
    echo "  ./scripts/manage-locks.sh unlock <file>"
    echo "  ./scripts/manage-locks.sh check"
    echo ""
    echo "Development:"
    echo "  make dev                    # Start development servers"
    echo "  make test                   # Run all tests"
    echo "  make build                  # Build project"
    echo ""
    echo "Git workflow:"
    echo "  git checkout -b <branch>    # Create feature branch"
    echo "  git add .                   # Stage changes"
    echo "  git commit -m \"message\"     # Commit changes"
    echo "  git push origin <branch>    # Push to remote"
    echo ""
    
    # Step 10: Completion
    header "ONBOARDING COMPLETE"
    
    success "You're all set up for development!"
    success "Developer ID: $dev_id"
    success "Feature branch: $branch_name"
    
    echo ""
    echo "Next steps:"
    echo "1. Read the MULTI_DEVELOPER_COORDINATION_GUIDELINES.md file"
    echo "2. Check for any active locks with: ./scripts/manage-locks.sh check"
    echo "3. Start developing on your feature branch"
    echo "4. Communicate with the team before modifying shared files"
    echo ""
    
    info "Welcome to the team! 🎉"
    echo ""
    
    # Test the lock system
    echo "Testing lock system..."
    if ./scripts/manage-locks.sh lock "README.md" "$dev_id" "Testing lock system" "1min"; then
        success "Lock system test passed"
        ./scripts/manage-locks.sh unlock "README.md" >/dev/null 2>&1
    else
        error "Lock system test failed"
        exit 1
    fi
    
    success "Onboarding completed successfully!"
}

# Run the main function
main "$@"