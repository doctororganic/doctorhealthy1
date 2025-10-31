#!/bin/bash

# 🚀 Complete Coolify Deployment Script for Nutrition Platform
# This script handles the full deployment process including Docker build, security setup, and Coolify configuration

set -e

echo "🍎 Starting Nutrition Platform Deployment to Coolify"
echo "=================================================="

# Configuration
PROJECT_NAME="nutrition-platform-complete"
DOMAIN="super.doctorhealthy1.com"
COOLIFY_PROJECT_ID="j0w00gog0c84owww80csk0c4"
DEPLOYMENT_DIR="nutrition-platform/coolify-complete-project"

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

# Step 1: Validate project structure
validate_project() {
    print_status "Step 1: Validating project structure..."

    if [ ! -d "$DEPLOYMENT_DIR" ]; then
        print_error "Deployment directory not found: $DEPLOYMENT_DIR"
        exit 1
    fi

    # Check required files
    required_files=("Dockerfile" "docker-compose.yml" "main.go" "go.mod" "nginx/nginx.conf" "nginx/conf.d/default.conf")
    for file in "${required_files[@]}"; do
        if [ ! -f "$DEPLOYMENT_DIR/$file" ]; then
            print_error "Required file missing: $file"
            exit 1
        fi
    done

    print_success "Project structure validated"
}

# Step 2: Test Docker build locally
test_docker_build() {
    print_status "Step 2: Testing Docker build locally..."

    cd "$DEPLOYMENT_DIR"

    # Build Docker image
    if docker build -t nutrition-platform-test .; then
        print_success "Docker build successful"
    else
        print_error "Docker build failed"
        exit 1
    fi

    # Test container startup
    if docker run -d --name test-container -p 8080:8080 nutrition-platform-test; then
        print_success "Container started successfully"

        # Wait for health check
        sleep 10

        # Test health endpoint
        if curl -f http://localhost:8080/health; then
            print_success "Health check passed"
        else
            print_error "Health check failed"
            docker logs test-container
        fi

        # Cleanup
        docker stop test-container
        docker rm test-container
    else
        print_error "Container failed to start"
        exit 1
    fi

    cd - > /dev/null
}

# Step 3: Create deployment archive
create_deployment_archive() {
    print_status "Step 3: Creating deployment archive..."

    cd "$DEPLOYMENT_DIR"

    # Create archive with timestamp
    TIMESTAMP=$(date +"%Y%m%d-%H%M%S")
    ARCHIVE_NAME="nutrition-platform-deploy-$TIMESTAMP.tar.gz"

    # Create archive excluding unnecessary files
    tar -czf "../$ARCHIVE_NAME" \
        --exclude='.git' \
        --exclude='node_modules' \
        --exclude='.DS_Store' \
        --exclude='*.log' \
        --exclude='logs/' \
        --exclude='tmp/' \
        --exclude='temp/' \
        .

    cd - > /dev/null

    DEPLOYMENT_ARCHIVE="$DEPLOYMENT_DIR/../$ARCHIVE_NAME"
    print_success "Deployment archive created: $DEPLOYMENT_ARCHIVE"
}

# Step 4: Generate secure environment variables
generate_secure_env_vars() {
    print_status "Step 4: Generating secure environment variables..."

    # Generate secure random values
    JWT_SECRET=$(openssl rand -hex 32)
    API_KEY_SECRET=$(openssl rand -hex 32)
    ENCRYPTION_KEY=$(openssl rand -hex 16)
    DB_PASSWORD=$(openssl rand -hex 16)
    REDIS_PASSWORD=$(openssl rand -hex 16)

    # Create environment variables file for Coolify
    cat > nutrition-platform/coolify-env-vars-final.txt << EOF
# Production Environment Variables for Coolify Deployment
# Generated on $(date)

# Server Configuration
SERVER_PORT=8081
SERVER_HOST=0.0.0.0
ENVIRONMENT=production
DEBUG=false

# Security - SECURE GENERATED VALUES
JWT_SECRET=$JWT_SECRET
API_KEY_SECRET=$API_KEY_SECRET
ENCRYPTION_KEY=$ENCRYPTION_KEY

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://$DOMAIN,https://www.$DOMAIN

# Database
DB_HOST=postgres
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=nutrition_user
DB_PASSWORD=$DB_PASSWORD
DB_SSL_MODE=disable

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=$REDIS_PASSWORD
REDIS_DB=0

# Application Settings
LOG_LEVEL=info
DATA_PATH=./data
NUTRITION_DATA_PATH=./
DEFAULT_LANGUAGE=en
SUPPORTED_LANGUAGES=en,ar
HEALTH_CHECK_ENABLED=true

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60s

# Security Headers
SECURITY_HEADERS_ENABLED=true
COMPRESSION_ENABLED=true
COMPRESSION_LEVEL=6

# Monitoring
METRICS_ENABLED=true
METRICS_PORT=9090
METRICS_PATH=/metrics

# Upload and Storage
UPLOAD_PATH=./uploads
BACKUP_PATH=./backups
DEFAULT_PAGE_SIZE=20
MAX_PAGE_SIZE=100

# Religious Compliance
RELIGIOUS_FILTER_ENABLED=true
FILTER_ALCOHOL=true
FILTER_PORK=true
FILTER_STRICT_MODE=true

# Localization
RTL_LANGUAGES=ar
EOF

    print_success "Secure environment variables generated"
    print_warning "IMPORTANT: Save the generated values securely!"
    print_warning "JWT_SECRET: $JWT_SECRET"
    print_warning "API_KEY_SECRET: $API_KEY_SECRET"
    print_warning "ENCRYPTION_KEY: $ENCRYPTION_KEY"
}

# Step 5: Provide Coolify deployment instructions
provide_deployment_instructions() {
    print_status "Step 5: Providing Coolify deployment instructions..."

    cat << 'EOF'

🔧 MANUAL COOLIFY DEPLOYMENT STEPS
===================================

1. **Access Coolify Dashboard:**
   - URL: https://api.doctorhealthy1.com
   - Navigate to project: "new doctorhealthy1"

2. **Create New Application:**
   - Click "Create Application"
   - Choose "Upload ZIP" as source type
   - Upload the generated archive: nutrition-platform-deploy-*.tar.gz

3. **Configure Application:**
   - Name: nutrition-platform-complete
   - Build Pack: Dockerfile
   - Dockerfile Location: ./Dockerfile
   - Port: 8080

4. **Set Environment Variables:**
   - Copy all variables from: nutrition-platform/coolify-env-vars-final.txt
   - Verify each value is set correctly

5. **Configure Domain:**
   - Primary domain: super.doctorhealthy1.com
   - Enable SSL: Yes (automatic)
   - Force HTTPS: Yes

6. **Deploy:**
   - Click "Deploy" button
   - Monitor build logs
   - Wait for completion (5-10 minutes)

7. **Verify Deployment:**
   - Test health endpoint: https://super.doctorhealthy1.com/health
   - Check SSL certificate
   - Test API endpoints

EOF
}

# Step 6: Post-deployment verification script
create_verification_script() {
    print_status "Step 6: Creating post-deployment verification script..."

    cat > nutrition-platform/verify-deployment.sh << 'EOF'
#!/bin/bash

# Post-Deployment Verification Script
DOMAIN="super.doctorhealthy1.com"
BASE_URL="https://$DOMAIN"

echo "🔍 Verifying Nutrition Platform Deployment"
echo "=========================================="

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

check_endpoint() {
    local url=$1
    local expected_code=${2:-200}
    local description=$3

    echo -n "Testing $description: "
    if curl -s -o /dev/null -w "%{http_code}" "$url" | grep -q "^$expected_code$"; then
        echo -e "${GREEN}✓ PASS${NC}"
        return 0
    else
        echo -e "${RED}✗ FAIL${NC}"
        return 1
    fi
}

# Health Check
check_endpoint "$BASE_URL/health" 200 "Health Check"

# API Info
check_endpoint "$BASE_URL/api/info" 200 "API Info"

# Nutrition Analysis (will fail without API key, but should return 401/403)
curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/nutrition/analyze" | grep -q "^40[13]$" && \
    echo -e "Testing Nutrition Analysis: ${GREEN}✓ PASS${NC} (Protected)" || \
    echo -e "Testing Nutrition Analysis: ${YELLOW}⚠ CHECK${NC} (Unexpected response)"

# Frontend
check_endpoint "$BASE_URL/" 200 "Frontend"

# SSL Certificate Check
echo -n "Testing SSL Certificate: "
if echo | openssl s_client -servername "$DOMAIN" -connect "$DOMAIN:443" 2>/dev/null | openssl x509 -noout -dates > /dev/null 2>&1; then
    echo -e "${GREEN}✓ PASS${NC}"
else
    echo -e "${RED}✗ FAIL${NC}"
fi

# CORS Headers Check
echo -n "Testing CORS Headers: "
cors_header=$(curl -s -I "$BASE_URL/api/info" | grep -i "access-control-allow-origin" | head -1)
if [[ $cors_header == *"super.doctorhealthy1.com"* ]]; then
    echo -e "${GREEN}✓ PASS${NC}"
else
    echo -e "${RED}✗ FAIL${NC}"
fi

echo ""
echo "🎉 Verification Complete!"
echo "If all tests pass, your deployment is successful!"
EOF

    chmod +x nutrition-platform/verify-deployment.sh
    print_success "Verification script created: nutrition-platform/verify-deployment.sh"
}

# Step 7: Security audit
perform_security_audit() {
    print_status "Step 7: Performing security audit..."

    echo "🔒 Security Audit Results:"
    echo "=========================="

    # Check for hardcoded secrets
    if grep -r "password\|secret\|key" "$DEPLOYMENT_DIR" --include="*.go" --include="*.js" --include="*.env*" | grep -v "secure_" | grep -v "generated" | grep -v "example"; then
        print_warning "Potential hardcoded secrets found - review before deployment"
    else
        print_success "No hardcoded secrets detected"
    fi

    # Check file permissions
    if find "$DEPLOYMENT_DIR" -name "*.key" -o -name "*.pem" -o -name "*secret*" | grep -q .; then
        print_warning "Sensitive files detected - ensure proper permissions"
    else
        print_success "No sensitive files with incorrect extensions"
    fi

    # Check for debug mode
    if grep -r "DEBUG.*true\|debug.*true" "$DEPLOYMENT_DIR" --include="*.go" --include="*.js"; then
        print_warning "Debug mode enabled in production code"
    else
        print_success "Debug mode properly disabled for production"
    fi

    print_success "Security audit completed"
}

# Main execution
main() {
    echo "🍎 Nutrition Platform - Complete Coolify Deployment"
    echo "=================================================="
    echo ""

    validate_project
    test_docker_build
    create_deployment_archive
    generate_secure_env_vars
    provide_deployment_instructions
    create_verification_script
    perform_security_audit

    echo ""
    print_success "🎉 Pre-deployment preparation complete!"
    echo ""
    echo "📦 Deployment archive: $DEPLOYMENT_ARCHIVE"
    echo "📋 Environment variables: nutrition-platform/coolify-env-vars-final.txt"
    echo "🔍 Verification script: nutrition-platform/verify-deployment.sh"
    echo ""
    echo "🚀 Ready for Coolify deployment!"
    echo "Follow the manual steps provided above to complete the deployment."
}

# Run main function
main "$@"