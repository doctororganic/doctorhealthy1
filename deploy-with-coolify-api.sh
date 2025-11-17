#!/bin/bash

# Deploy Nutrition Platform using Coolify API directly
# This script uses the Coolify REST API to deploy the project

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${BLUE}[DEPLOY]${NC} $1"; }
success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; }
warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }

echo ""
echo "🚀 =================================="
echo "🚀 DEPLOYING NUTRITION PLATFORM"
echo "🚀 =================================="
echo ""

# Load environment variables
if [ -f ".env.secure" ]; then
    export $(grep -v '^#' .env.secure | xargs)
    success "Loaded secure environment variables"
else
    error ".env.secure file not found"
    exit 1
fi

# Create deployment package
log "Creating deployment package..."
tar -czf nutrition-platform-deploy.tar.gz \
    --exclude='.git' \
    --exclude='node_modules' \
    --exclude='*.log' \
    --exclude='archive/' \
    --exclude='coolify-deployment-*/' \
    --exclude='nutrition-platform-coolify-*/' \
    --exclude='.env*' \
    .

success "Deployment package created"

# Get existing applications to find project ID
log "Checking for existing applications..."
APPS_RESPONSE=$(curl -s -H "Authorization: Bearer $COOLIFY_API_TOKEN" \
    -H "Accept: application/json" \
    "$COOLIFY_API_URL/api/v1/applications" 2>/dev/null)

# Extract the first application ID (assuming we want to update existing)
APP_ID=$(echo "$APPS_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*' || echo "")

if [ -z "$APP_ID" ]; then
    warning "No existing application found, will create new one"
    CREATE_NEW=true
else
    success "Found existing application ID: $APP_ID"
    CREATE_NEW=false
fi

# Deploy or update application
if [ "$CREATE_NEW" = true ]; then
    log "Creating new application..."
    
    # Create application
    CREATE_RESPONSE=$(curl -s -X POST \
        -H "Authorization: Bearer $COOLIFY_API_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"name\": \"nutrition-platform\",
            \"description\": \"AI-powered nutrition and health management platform\",
            \"server_uuid\": \"x8gck8ggggsgkggg4coosg0g\",
            \"domains\": [\"super.doctorhealthy1.com\"],
            \"build_pack\": \"dockercompose\",
            \"port\": 3000,
            \"environment_variables\": {
                \"NODE_ENV\": \"production\",
                \"PORT\": \"8080\",
                \"DB_HOST\": \"postgres\",
                \"DB_PORT\": \"5432\",
                \"DB_NAME\": \"nutrition_platform\",
                \"DB_USER\": \"nutrition_user\",
                \"DB_PASSWORD\": \"nutrition_pass\",
                \"REDIS_HOST\": \"redis\",
                \"REDIS_PORT\": \"6379\",
                \"ENVIRONMENT\": \"production\",
                \"NEXT_PUBLIC_API_URL\": \"http://backend:8080\"
            },
            \"source\": {
                \"type\": \"git\",
                \"repository\": \"nutrition-platform\",
                \"branch\": \"main\"
            },
            \"health_check\": {
                \"path\": \"/health\",
                \"port\": 8080,
                \"interval\": 30,
                \"timeout\": 10,
                \"retries\": 3
            },
            \"force_https\": true,
            \"is_static\": false
        }" \
        "$COOLIFY_API_URL/api/v1/applications" 2>/dev/null)
    
    if echo "$CREATE_RESPONSE" | grep -q "error\|Error"; then
        error "Application creation failed: $CREATE_RESPONSE"
        exit 1
    else
        APP_ID=$(echo "$CREATE_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')
        success "Application created with ID: $APP_ID"
    fi
else
    log "Updating existing application..."
    
    # Update application
    UPDATE_RESPONSE=$(curl -s -X PATCH \
        -H "Authorization: Bearer $COOLIFY_API_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"server_uuid\": \"x8gck8ggggsgkggg4coosg0g\",
            \"domains\": [\"super.doctorhealthy1.com\"],
            \"environment_variables\": {
                \"NODE_ENV\": \"production\",
                \"PORT\": \"8080\",
                \"DB_HOST\": \"postgres\",
                \"DB_PORT\": \"5432\",
                \"DB_NAME\": \"nutrition_platform\",
                \"DB_USER\": \"nutrition_user\",
                \"DB_PASSWORD\": \"nutrition_pass\",
                \"REDIS_HOST\": \"redis\",
                \"REDIS_PORT\": \"6379\",
                \"ENVIRONMENT\": \"production\",
                \"NEXT_PUBLIC_API_URL\": \"http://backend:8080\"
            },
            \"force_https\": true
        }" \
        "$COOLIFY_API_URL/api/v1/applications/$APP_ID" 2>/dev/null)
    
    if echo "$UPDATE_RESPONSE" | grep -q "error\|Error"; then
        error "Application update failed: $UPDATE_RESPONSE"
        exit 1
    else
        success "Application updated successfully"
    fi
fi

# Trigger deployment
log "Triggering deployment..."
DEPLOY_RESPONSE=$(curl -s -X POST \
    -H "Authorization: Bearer $COOLIFY_API_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"force_rebuild\": true,
        \"debug\": false,
        \"pull\": true
    }" \
    "$COOLIFY_API_URL/api/v1/applications/$APP_ID/deploy" 2>/dev/null)

if echo "$DEPLOY_RESPONSE" | grep -q "error\|Error"; then
    error "Deployment trigger failed: $DEPLOY_RESPONSE"
    exit 1
else
    success "Deployment triggered successfully"
fi

# Extract deployment ID for monitoring
DEPLOYMENT_ID=$(echo "$DEPLOY_RESPONSE" | grep -o '"deployment_id":"[^"]*"' | cut -d'"' -f4 || echo "")

# Monitor deployment
log "Monitoring deployment progress..."
TIMEOUT=600  # 10 minutes
ELAPSED=0
INTERVAL=30

while [ $ELAPSED -lt $TIMEOUT ]; do
    if [ -n "$DEPLOYMENT_ID" ] && [ "$DEPLOYMENT_ID" != "null" ]; then
        STATUS_RESPONSE=$(curl -s -H "Authorization: Bearer $COOLIFY_API_TOKEN" \
            "$COOLIFY_API_URL/api/v1/deployments/$DEPLOYMENT_ID" 2>/dev/null)
        STATUS=$(echo "$STATUS_RESPONSE" | grep -o '"status":"[^"]*"' | cut -d'"' -f4 || echo "unknown")
    else
        # Fallback: check application status
        APP_STATUS_RESPONSE=$(curl -s -H "Authorization: Bearer $COOLIFY_API_TOKEN" \
            "$COOLIFY_API_URL/api/v1/applications/$APP_ID" 2>/dev/null)
        STATUS=$(echo "$APP_STATUS_RESPONSE" | grep -o '"status":"[^"]*"' | cut -d'"' -f4 || echo "unknown")
    fi

    case "$STATUS" in
        "success"|"finished"|"completed"|"running"|"active")
            success "✅ Deployment completed successfully!"
            break
            ;;
        "failed"|"error")
            error "❌ Deployment failed"
            exit 1
            ;;
        "building"|"in_progress"|"deploying"|"pending"|"queued")
            log "🔄 Deployment in progress... ($ELAPSED/$TIMEOUT seconds)"
            ;;
        *)
            log "📊 Deployment status: $STATUS ($ELAPSED/$TIMEOUT seconds)"
            ;;
    esac

    sleep $INTERVAL
    ELAPSED=$((ELAPSED + INTERVAL))
done

if [ $ELAPSED -ge $TIMEOUT ]; then
    warning "⚠️ Deployment monitoring timed out, but deployment may still be running"
fi

# Cleanup
log "Cleaning up temporary files..."
rm -f nutrition-platform-deploy.tar.gz

echo ""
echo "🎉 ===================================="
echo "🎉 DEPLOYMENT COMPLETED!"
echo "🎉 ===================================="
echo ""
echo "📍 Your Application is LIVE:"
echo "   🌐 Website: https://super.doctorhealthy1.com"
echo "   🏥 Health Check: https://super.doctorhealthy1.com/health"
echo "   📊 API Base: https://super.doctorhealthy1.com/api"
echo ""
echo "🖥️ Coolify Dashboard:"
echo "   🌐 Management: $COOLIFY_API_URL"
echo "   📱 Application ID: $APP_ID"
echo ""
echo "📋 Next Steps:"
echo "   1. 🔑 Test all API endpoints"
echo "   2. 📊 Monitor in Coolify dashboard"
echo "   3. 🌐 Share your application URL"
echo "   4. 📱 Test on mobile devices"
echo ""
success "🚀 Nutrition Platform deployed successfully!"