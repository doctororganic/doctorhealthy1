#!/bin/bash

# Coolify Deployment Helper Script
# This script provides all the exact steps and configurations for Coolify deployment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Function to print colored output
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log "🚀 NUTRITION PLATFORM COOLIFY DEPLOYMENT HELPER"
echo ""

# Display exact deployment steps
echo "📋 EXACT COOLIFY DEPLOYMENT STEPS"
echo "===================================="
echo ""
echo "🌐 STEP 1: Access Coolify Dashboard"
echo "   Go to: https://api.doctorhealthy1.com"
echo "   Login with your Coolify credentials"
echo ""
echo "📍 STEP 2: Navigate to Applications"
echo "   Click 'Applications' in left sidebar"
echo "   Click 'Add Application' button"
echo ""
echo "📦 STEP 3: Upload Deployment Package"
echo "   Select 'Upload ZIP file'"
echo "   Choose file: nutrition-platform-coolify-20251013-164858.zip"
echo "   Application Name: nutrition-platform-secure"
echo "   Description: AI-powered nutrition platform with enterprise security"
echo ""
echo "🔧 STEP 4: Configure Source Settings"
echo "   Source Type: Archive"
echo "   Archive Type: ZIP file"
echo "   Root Directory: / (root of archive)"
echo ""
echo "🏗️ STEP 5: Configure Build Settings"
echo "   Build Pack: Dockerfile"
echo "   Dockerfile Location: backend/Dockerfile"
echo "   Build Context: ./"
echo "   Install Command: (blank)"
echo "   Build Command: (blank)"
echo "   Start Command: (blank)"
echo ""
echo "🌐 STEP 6: Configure Deployment Settings"
echo "   Domain: super.doctorhealthy1.com"
echo "   Port: 8080"
echo "   Health Check Path: /health"
echo "   Health Check Interval: 30s"
echo "   Auto Deploy: Enabled"
echo ""
echo "⚙️ STEP 7: Add Environment Variables"
echo "   Click 'Environment Variables' tab"
echo "   Click 'Bulk Import'"
echo "   Copy and paste ALL variables from .env.production:"
echo ""

# Display environment variables
cat << 'EOF'
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=nutrition_user
DB_PASSWORD=ac287cc0e30f54afad53c6dc7e02fd0cccad979d62b75d75d97b1ede12daf8d5
DB_SSL_MODE=require
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0
REDIS_PASSWORD=f606b2d16d6697e666ce78a8685574d042df15484ca8f18f39f2e67bf38dc09a

# Security Configuration
JWT_SECRET=9a00511e8e23764f8f4524c02f1db9eccc1923208c02fb36cb758d874d8d569bce9ea1b24ac18a958334abe15ef89e09d6010fe64a1d1ffc02a45b07898b2473
API_KEY_SECRET=5d2763e839f7e71b90ff88bef12f690a41802635aa131f6bc7160056ef0aeb7dc9caaeb07dbe0028128e617529a48903f8d01c6cc64ce61419eb7f309fdfc8bc
ENCRYPTION_KEY=cc1574e486b2f5abd69d86537079ba928974cc463e36ff410647b15b15533d23
SESSION_SECRET=f40776484ee20b35e4f754909fb3067cef2a186d0da7c4c24f1bcd54870d9fba

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://super.doctorhealthy1.com,https://my.doctorhealthy1.com

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Features
RELIGIOUS_FILTER_ENABLED=true
FILTER_ALCOOL=true
FILTER_PORK=true
FILTER_STRICT_MODE=false

# Internationalization
DEFAULT_LANGUAGE=en
SUPPORTED_LANGUAGES=en,ar
RTL_LANGUAGES=ar

# Health Check
HEALTH_CHECK_ENABLED=true
HEALTH_CHECK_INTERVAL=30
HEALTH_CHECK_TIMEOUT=5
EOF

echo ""
echo "🗄️ STEP 8: Add Database Services"
echo "   Click 'Services' tab"
echo "   Click 'Add Service'"
echo "   Select 'PostgreSQL'"
echo "   Name: nutrition-postgres"
echo "   Version: 15"
echo "   Database: nutrition_platform"
echo "   Username: nutrition_user"
echo "   Password: ac287cc0e30f54afad53c6dc7e02fd0cccad979d62b75d75d97b1ede12daf8d5"
echo ""
echo "   Click 'Add Another Service'"
echo "   Select 'Redis'"
echo "   Name: nutrition-redis"
echo "   Version: 7-alpine"
echo "   Password: f606b2d16d6697e666ce78a8685574d042df15484ca8f18f39f2e67bf38dc09a"
echo ""
echo "🚀 STEP 9: Deploy Application"
echo "   Click 'Deploy' button (top right)"
echo "   Wait 5-10 minutes for deployment"
echo "   Monitor the deployment in the 'Deployments' tab"
echo ""

# Display verification steps
echo "🔍 POST-DEPLOYMENT VERIFICATION"
echo "==================================="
echo ""
echo "Test these URLs after deployment:"
echo ""
echo "🌐 Main Site: https://super.doctorhealthy1.com"
echo "🔍 Health Check: https://super.doctorhealthy1.com/health"
echo "📚 API Info: https://super.doctorhealthy1.com/api/v1/info"
echo "🧪 Nutrition Test: POST https://super.doctorhealthy1.com/api/v1/nutrition/analyze"
echo ""
echo "Test Payload for Nutrition API:"
echo "{"
echo "  \"food\": \"chicken breast\","
echo "  \"quantity\": 100,"
echo "  \"unit\": \"grams\","
echo "  \"checkHalal\": true,"
echo "  \"language\": \"en\""
echo "}"
echo ""
echo "Expected Health Response:"
echo "{"
echo "  \"status\": \"healthy\","
echo "  \"timestamp\": \"2025-10-13T...\","
echo "  \"version\": \"1.0.0\""
echo "}"
echo ""

# Create a monitoring script
cat > monitor-deployment.sh << 'EOF'
#!/bin/bash

# Monitor deployment script
echo "🔍 Monitoring deployment..."

DOMAIN="https://super.doctorhealthy1.com"
MAX_ATTEMPTS=30
ATTEMPT=0

while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
    echo "📊 Check attempt $((ATTEMPT + 1))/$MAX_ATTEMPTS"
    
    # Check main site
    if curl -f -s --max-time 10 "$DOMAIN" > /dev/null; then
        echo "✅ Main site is accessible: $DOMAIN"
        
        # Check health endpoint
        if curl -f -s --max-time 10 "$DOMAIN/health" > /dev/null; then
            echo "✅ Health endpoint is working: $DOMAIN/health"
            
            # Check API endpoint
            if curl -f -s --max-time 10 "$DOMAIN/api/v1/info" > /dev/null; then
                echo "✅ API endpoint is working: $DOMAIN/api/v1/info"
                echo ""
                echo "🎉 DEPLOYMENT SUCCESSFUL!"
                echo "📍 Your nutrition platform is live at: $DOMAIN"
                exit 0
            else
                echo "⚠️ API endpoint not ready yet"
            fi
        else
            echo "⚠️ Health endpoint not ready yet"
        fi
    else
        echo "⚠️ Main site not ready yet"
    fi
    
    sleep 30
    ATTEMPT=$((ATTEMPT + 1))
done

echo ""
echo "⚠️ Deployment monitoring timed out"
echo "Please check the Coolify dashboard for deployment status"
EOF

chmod +x monitor-deployment.sh

success "✅ Deployment helper script created!"
echo ""
echo "📋 NEXT ACTIONS:"
echo "1. Follow the steps above to deploy in Coolify"
echo "2. Run './monitor-deployment.sh' to monitor deployment progress"
echo "3. Check the URLs once deployment is complete"
echo ""
echo "🎊 Your nutrition platform will be live at: https://super.doctorhealthy1.com"