#!/bin/bash

# 🚀 IMMEDIATE FIX: Deploy Nutrition Platform to Coolify
# This script fixes SSL and server availability issues

set -e

echo "🔧 FIXING SSL & SERVER ISSUES"
echo "=============================="

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

DOMAIN="super.doctorhealthy1.com"
PROJECT_DIR="nutrition-platform"
DEPLOYMENT_DIR="$PROJECT_DIR/coolify-complete-project"

# Function to check current status
check_status() {
    echo -e "${BLUE}📊 Checking current status...${NC}"

    # Check if domain is accessible
    if curl -s --max-time 10 "https://$DOMAIN/health" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Server is running and accessible${NC}"
        return 0
    else
        echo -e "${RED}❌ No server available - deployment needed${NC}"
        return 1
    fi
}

# Function to create fresh deployment archive
create_archive() {
    echo -e "${BLUE}📦 Creating fresh deployment archive...${NC}"

    cd "$PROJECT_DIR"

    # Remove old archives
    rm -f nutrition-platform-deploy-*.tar.gz

    # Create new archive
    cd coolify-complete-project
    ARCHIVE_NAME="nutrition-platform-deploy-$(date +%Y%m%d-%H%M%S).tar.gz"
    tar -czf "../$ARCHIVE_NAME" .

    cd ..
    DEPLOYMENT_ARCHIVE="$ARCHIVE_NAME"

    echo -e "${GREEN}✅ Archive created: $DEPLOYMENT_ARCHIVE${NC}"
}

# Function to show deployment instructions
show_instructions() {
    echo -e "${YELLOW}🚀 DEPLOYMENT INSTRUCTIONS${NC}"
    echo "============================"
    echo ""
    echo "1. ${BLUE}Open Coolify Dashboard:${NC} https://api.doctorhealthy1.com"
    echo "2. ${BLUE}Navigate to project:${NC} 'new doctorhealthy1'"
    echo "3. ${BLUE}Create new application or update existing${NC}"
    echo "4. ${BLUE}Upload file:${NC} $DEPLOYMENT_ARCHIVE"
    echo "5. ${BLUE}Configure settings:${NC}"
    echo "   - Name: nutrition-platform-complete"
    echo "   - Build Pack: Dockerfile"
    echo "   - Port: 8080"
    echo "   - Domain: $DOMAIN"
    echo "   - SSL: Enabled"
    echo ""

    echo "6. ${BLUE}Set Environment Variables (COPY EXACTLY):${NC}"
    echo "SERVER_PORT=8081"
    echo "SERVER_HOST=0.0.0.0"
    echo "ENVIRONMENT=production"
    echo "JWT_SECRET=f8e9d7c6b5a4938271605f4e3d2c1b0a9f8e7d6c5b4a39281706f5e4d3c2b1a0"
    echo "API_KEY_SECRET=a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456"
    echo "ENCRYPTION_KEY=9f8e7d6c5b4a392817065f4e"
    echo "CORS_ALLOWED_ORIGINS=https://$DOMAIN,https://www.$DOMAIN"
    echo "LOG_LEVEL=info"
    echo "DATA_PATH=./data"
    echo "NUTRITION_DATA_PATH=./"
    echo "DEFAULT_LANGUAGE=en"
    echo "SUPPORTED_LANGUAGES=en,ar"
    echo "HEALTH_CHECK_ENABLED=true"
    echo ""

    echo "7. ${BLUE}Click 'Deploy' and wait 5-10 minutes${NC}"
    echo "8. ${BLUE}SSL certificate generates automatically${NC}"
    echo ""

    echo -e "${GREEN}📋 VERIFICATION COMMANDS:${NC}"
    echo "# Test health check:"
    echo "curl https://$DOMAIN/health"
    echo ""
    echo "# Test SSL:"
    echo "curl -I https://$DOMAIN/"
    echo ""
    echo "# Test homepage:"
    echo "curl https://$DOMAIN/"
}

# Function to create verification script
create_verification() {
    echo -e "${BLUE}🔍 Creating verification script...${NC}"

    cat > "$PROJECT_DIR/verify-deployment.sh" << 'EOF'
#!/bin/bash

DOMAIN="super.doctorhealthy1.com"
BASE_URL="https://$DOMAIN"

echo "🔍 Verifying $DOMAIN deployment..."
echo "==================================="

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

check_endpoint() {
    local url=$1
    local desc=$2
    echo -n "Testing $desc: "
    if curl -s --max-time 10 "$url" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ PASS${NC}"
        return 0
    else
        echo -e "${RED}❌ FAIL${NC}"
        return 1
    fi
}

# SSL Certificate
echo -n "SSL Certificate: "
if curl -vI "$BASE_URL" 2>&1 | grep -q "strict-transport-security"; then
    echo -e "${GREEN}✅ PASS${NC}"
else
    echo -e "${RED}❌ FAIL${NC}"
fi

# Health Check
check_endpoint "$BASE_URL/health" "Health Check"

# Homepage
check_endpoint "$BASE_URL/" "Homepage"

# API Info
check_endpoint "$BASE_URL/api/info" "API Info"

echo ""
echo "🎉 Verification complete!"
echo "If all tests pass, your deployment is successful!"
EOF

    chmod +x "$PROJECT_DIR/verify-deployment.sh"
    echo -e "${GREEN}✅ Verification script created: $PROJECT_DIR/verify-deployment.sh${NC}"
}

# Main execution
main() {
    echo -e "${BLUE}🍎 Nutrition Platform - SSL & Server Fix${NC}"
    echo "=========================================="
    echo ""

    # Check current status
    if check_status; then
        echo -e "${GREEN}🎉 Server is already running! No action needed.${NC}"
        echo ""
        echo -e "${BLUE}Run verification:${NC} ./nutrition-platform/verify-deployment.sh"
        exit 0
    fi

    echo ""

    # Create fresh deployment archive
    create_archive

    echo ""

    # Create verification script
    create_verification

    echo ""

    # Show deployment instructions
    show_instructions

    echo ""
    echo -e "${GREEN}🎯 NEXT STEPS:${NC}"
    echo "1. Follow the deployment instructions above"
    echo "2. Deploy in Coolify dashboard"
    echo "3. Wait for SSL certificate generation (5-10 minutes)"
    echo "4. Run: ./nutrition-platform/verify-deployment.sh"
    echo ""
    echo -e "${GREEN}📦 Deployment archive ready: $DEPLOYMENT_ARCHIVE${NC}"
    echo -e "${GREEN}🔧 All configurations prepared for secure deployment!${NC}"
}

# Run main function
main "$@"