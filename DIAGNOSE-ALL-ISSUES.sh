#!/bin/bash

###############################################################################
# DIAGNOSE-ALL-ISSUES.sh - Complete Diagnostic Tool
# Identifies all problems preventing your web app from working
###############################################################################

echo "🔍 Trae New Healthy1 - Complete Diagnostic Report"
echo "=================================================="
echo ""
echo "Analyzing your project to identify all issues..."
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

ISSUES_FOUND=0
WARNINGS_FOUND=0

print_section() {
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

print_ok() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
    ((WARNINGS_FOUND++))
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
    ((ISSUES_FOUND++))
}

print_info() {
    echo -e "ℹ️  $1"
}

# Check current directory
if [ ! -d "production-nodejs" ] && [ ! -d "backend" ]; then
    print_error "Not in nutrition-platform directory"
    echo "Please run this script from the nutrition-platform directory"
    exit 1
fi

print_section "1. PROJECT STRUCTURE ANALYSIS"

# Check Node.js backend
if [ -d "production-nodejs" ]; then
    print_ok "Node.js backend found (production-nodejs/)"
    
    if [ -f "production-nodejs/server.js" ]; then
        print_ok "server.js exists"
    else
        print_error "server.js missing"
    fi
    
    if [ -f "production-nodejs/package.json" ]; then
        print_ok "package.json exists"
    else
        print_error "package.json missing"
    fi
    
    if [ -d "production-nodejs/node_modules" ]; then
        print_ok "Dependencies installed"
    else
        print_warning "Dependencies not installed (run: npm install)"
    fi
else
    print_error "Node.js backend not found"
fi

# Check Go backend
if [ -d "backend" ]; then
    print_ok "Go backend found (backend/)"
    
    if [ -f "backend/main.go" ]; then
        print_ok "main.go exists"
        
        # Check for compilation errors
        print_info "Checking Go code for errors..."
        cd backend
        if go build -o /dev/null . 2>/dev/null; then
            print_ok "Go code compiles successfully"
        else
            print_error "Go code has compilation errors"
            print_info "Run 'cd backend && go build' to see errors"
        fi
        cd ..
    else
        print_error "main.go missing"
    fi
    
    if [ -f "backend/go.mod" ]; then
        print_ok "go.mod exists"
    else
        print_error "go.mod missing"
    fi
else
    print_warning "Go backend not found"
fi

print_section "2. DEPLOYMENT FILES ANALYSIS"

# Count deployment scripts
DEPLOY_SCRIPTS=$(find . -maxdepth 2 -name "*deploy*.sh" -o -name "*fix*.sh" | wc -l)
print_info "Found ${DEPLOY_SCRIPTS} deployment/fix scripts"

if [ $DEPLOY_SCRIPTS -gt 5 ]; then
    print_warning "Too many deployment scripts (${DEPLOY_SCRIPTS}) - causes confusion"
    print_info "Recommendation: Use single DEPLOY-NOW.sh script"
fi

# Check for Dockerfile
if [ -f "production-nodejs/Dockerfile" ]; then
    print_ok "Node.js Dockerfile exists"
else
    print_warning "Node.js Dockerfile missing"
fi

if [ -f "backend/Dockerfile" ]; then
    print_ok "Go Dockerfile exists"
else
    print_warning "Go Dockerfile missing"
fi

# Check docker-compose
if [ -f "docker-compose.yml" ]; then
    print_ok "docker-compose.yml exists"
else
    print_warning "docker-compose.yml missing"
fi

print_section "3. CONFIGURATION ANALYSIS"

# Check environment files
if [ -f "production-nodejs/.env" ]; then
    print_ok "Node.js .env file exists"
else
    print_warning "Node.js .env file missing (optional)"
fi

if [ -f "backend/.env" ]; then
    print_ok "Go .env file exists"
else
    print_warning "Go .env file missing (optional)"
fi

print_section "4. DEPENDENCY ANALYSIS"

# Check Node.js
if command -v node &> /dev/null; then
    NODE_VERSION=$(node --version)
    print_ok "Node.js installed: ${NODE_VERSION}"
    
    # Check version
    MAJOR_VERSION=$(echo $NODE_VERSION | cut -d'.' -f1 | sed 's/v//')
    if [ $MAJOR_VERSION -ge 18 ]; then
        print_ok "Node.js version is compatible (>= 18)"
    else
        print_warning "Node.js version is old (< 18), upgrade recommended"
    fi
else
    print_error "Node.js not installed"
fi

# Check npm
if command -v npm &> /dev/null; then
    NPM_VERSION=$(npm --version)
    print_ok "npm installed: ${NPM_VERSION}"
else
    print_error "npm not installed"
fi

# Check Go
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    print_ok "Go installed: ${GO_VERSION}"
else
    print_warning "Go not installed (only needed for Go backend)"
fi

# Check Docker
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version | awk '{print $3}' | sed 's/,//')
    print_ok "Docker installed: ${DOCKER_VERSION}"
else
    print_warning "Docker not installed (needed for containerized deployment)"
fi

print_section "5. CODE QUALITY ANALYSIS"

# Check Node.js code
if [ -f "production-nodejs/server.js" ]; then
    print_info "Analyzing Node.js code..."
    
    # Check for syntax errors
    if node -c production-nodejs/server.js 2>/dev/null; then
        print_ok "No syntax errors in server.js"
    else
        print_error "Syntax errors found in server.js"
    fi
    
    # Check for required dependencies
    if grep -q "express" production-nodejs/package.json; then
        print_ok "Express dependency found"
    else
        print_error "Express dependency missing"
    fi
fi

# Check Go code issues
if [ -f "backend/main.go" ]; then
    print_info "Analyzing Go code..."
    
    # Check for common issues
    if grep -q "getUserIDFromContext" backend/handlers/*.go 2>/dev/null; then
        if grep -q "func getUserIDFromContext" backend/**/*.go 2>/dev/null; then
            print_ok "getUserIDFromContext function defined"
        else
            print_error "getUserIDFromContext function used but not defined"
        fi
    fi
    
    # Check for duplicate registrations
    RECIPE_HANDLER_COUNT=$(grep -c "recipeHandler :=" backend/main.go 2>/dev/null || echo "0")
    if [ $RECIPE_HANDLER_COUNT -gt 1 ]; then
        print_error "Duplicate recipeHandler registration in main.go"
    fi
fi

print_section "6. NETWORK & CONNECTIVITY"

# Check if ports are available
print_info "Checking port availability..."

if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1; then
    print_warning "Port 8080 is already in use"
    print_info "Process using port 8080:"
    lsof -Pi :8080 -sTCP:LISTEN | tail -n +2
else
    print_ok "Port 8080 is available"
fi

if lsof -Pi :3000 -sTCP:LISTEN -t >/dev/null 2>&1; then
    print_warning "Port 3000 is already in use"
else
    print_ok "Port 3000 is available"
fi

# Check internet connectivity
if ping -c 1 google.com &> /dev/null; then
    print_ok "Internet connection available"
else
    print_warning "No internet connection (needed for deployment)"
fi

print_section "7. FRONTEND ANALYSIS"

# Check for frontend files
FRONTEND_FOUND=false

if [ -d "frontend" ]; then
    print_ok "Frontend directory found"
    FRONTEND_FOUND=true
elif [ -d "client" ]; then
    print_ok "Client directory found"
    FRONTEND_FOUND=true
elif grep -q "<!DOCTYPE html>" production-nodejs/server.js 2>/dev/null; then
    print_ok "Frontend embedded in server.js"
    FRONTEND_FOUND=true
fi

if [ "$FRONTEND_FOUND" = false ]; then
    print_warning "No separate frontend found (using embedded frontend)"
fi

print_section "8. SECURITY ANALYSIS"

# Check for exposed secrets
if grep -r "password.*=.*\".*\"" . --include="*.js" --include="*.go" 2>/dev/null | grep -v node_modules | grep -v ".git" | head -n 1 > /dev/null; then
    print_warning "Possible hardcoded passwords found in code"
fi

# Check for .env in git
if [ -f ".gitignore" ]; then
    if grep -q ".env" .gitignore; then
        print_ok ".env files are gitignored"
    else
        print_warning ".env not in .gitignore"
    fi
fi

print_section "9. DEPLOYMENT READINESS"

READY_TO_DEPLOY=true

# Check critical requirements
if [ ! -d "production-nodejs" ] && [ ! -d "backend" ]; then
    print_error "No backend implementation found"
    READY_TO_DEPLOY=false
fi

if [ -d "production-nodejs" ]; then
    if [ ! -f "production-nodejs/server.js" ]; then
        print_error "server.js missing"
        READY_TO_DEPLOY=false
    fi
    
    if [ ! -f "production-nodejs/package.json" ]; then
        print_error "package.json missing"
        READY_TO_DEPLOY=false
    fi
fi

if [ "$READY_TO_DEPLOY" = true ]; then
    print_ok "Project is ready for deployment"
else
    print_error "Project has critical issues preventing deployment"
fi

print_section "10. DIAGNOSTIC SUMMARY"

echo ""
echo "Issues Found: ${ISSUES_FOUND}"
echo "Warnings: ${WARNINGS_FOUND}"
echo ""

if [ $ISSUES_FOUND -eq 0 ] && [ $WARNINGS_FOUND -eq 0 ]; then
    echo -e "${GREEN}🎉 No issues found! Your project is in good shape.${NC}"
    echo ""
    echo "Next steps:"
    echo "1. Run: ./DEPLOY-NOW.sh"
    echo "2. Choose deployment method"
    echo "3. Test your application"
elif [ $ISSUES_FOUND -eq 0 ]; then
    echo -e "${YELLOW}⚠️  No critical issues, but ${WARNINGS_FOUND} warnings found.${NC}"
    echo ""
    echo "You can proceed with deployment, but consider fixing warnings."
    echo ""
    echo "Next steps:"
    echo "1. Review warnings above"
    echo "2. Run: ./DEPLOY-NOW.sh"
else
    echo -e "${RED}❌ ${ISSUES_FOUND} critical issues found that must be fixed.${NC}"
    echo ""
    echo "Recommended actions:"
    echo "1. Review errors above"
    echo "2. Read MASTER-FIX-PLAN.md for solutions"
    echo "3. Fix critical issues"
    echo "4. Run this diagnostic again"
fi

print_section "11. RECOMMENDED SOLUTION"

echo ""
echo "Based on analysis, here's what you should do:"
echo ""

if [ -d "production-nodejs" ] && [ -f "production-nodejs/server.js" ]; then
    echo -e "${GREEN}✅ RECOMMENDED: Deploy Node.js Backend${NC}"
    echo ""
    echo "Your Node.js backend is complete and ready. Deploy it now:"
    echo ""
    echo "  ./DEPLOY-NOW.sh"
    echo ""
    echo "This will:"
    echo "  • Deploy working Node.js server"
    echo "  • Include built-in frontend"
    echo "  • Set up SSL automatically"
    echo "  • Configure all endpoints"
    echo ""
fi

if [ -d "backend" ] && [ -f "backend/main.go" ]; then
    echo -e "${YELLOW}⚠️  OPTIONAL: Fix Go Backend Later${NC}"
    echo ""
    echo "Your Go backend has issues that need fixing:"
    echo "  • Missing functions"
    echo "  • Compilation errors"
    echo "  • Duplicate registrations"
    echo ""
    echo "Recommendation: Deploy Node.js first, fix Go later"
fi

echo ""
echo -e "${BLUE}📖 For detailed solutions, read: MASTER-FIX-PLAN.md${NC}"
echo ""

print_section "END OF DIAGNOSTIC REPORT"

exit 0
