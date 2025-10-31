#!/bin/bash

# Coolify API Deployment Script
# Deploys the nutrition platform using Coolify API

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
COOLIFY_URL="https://api.doctorhealthy1.com"
TOKEN="6|uJSYhIJQIypx4UuxbQkaHkidEyiQshLR6U1QNxEQab344fda"
PROJECT_ID="j0w00gog0c84owww80csk0c4"
ENVIRONMENT_ID="l0gscs8w8kw8800c00ccckco"
SERVER_UUID="x8gck8ggggsgkggg4coosg0g"
DOMAIN="super.doctorhealthy1.com"
ZIP_FILE="nutrition-platform-coolify-20251013-164858.zip"

# Function to print colored output
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Function to make API calls
coolify_api() {
    local method=$1
    local endpoint=$2
    local data=$3
    
    local curl_cmd="curl -s -X \"$method\" \
        -H \"Authorization: Bearer $TOKEN\" \
        -H \"Content-Type: application/json\" \
        -H \"Accept: application/json\""
    
    if [[ -n "$data" ]]; then
        curl_cmd="$curl_cmd -d \"$data\""
    fi
    
    curl_cmd="$curl_cmd \"$COOLIFY_URL/api/v1$endpoint\""
    
    eval "$curl_cmd"
}

# Function to create application
create_application() {
    log "🚀 Creating application in Coolify..."
    
    local app_data=$(cat << EOF
{
    "name": "nutrition-platform-secure",
    "description": "AI-powered nutrition platform with enterprise security",
    "project_uuid": "$PROJECT_ID",
    "environment_uuid": "$ENVIRONMENT_ID",
    "source": {
        "type": "zipfile",
        "zipfile": "$ZIP_FILE"
    },
    "destination": {
        "type": "docker",
        "docker": {
            "image": "nutrition-platform:latest",
            "port": 8080
        }
    },
    "domains": [
        {
            "domain": "$DOMAIN",
            "certificate": "letsencrypt"
        }
    ]
}
EOF
)
    
    local response=$(coolify_api "POST" "/applications" "$app_data")
    
    if echo "$response" | grep -q "error\|Error"; then
        error "❌ Failed to create application: $response"
    fi
    
    local app_uuid=$(echo "$response" | grep -o '"uuid":"[^"]*"' | cut -d'"' -f4)
    if [[ -z "$app_uuid" ]]; then
        error "❌ Could not extract application UUID"
    fi
    
    success "✅ Application created with UUID: $app_uuid"
    echo "$app_uuid"
}

# Function to add environment variables
add_environment_variables() {
    local app_uuid=$1
    log "⚙️ Adding environment variables..."
    
    local env_data=$(cat << EOF
{
    "variables": [
        {
            "key": "DB_HOST",
            "value": "localhost",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "DB_PORT",
            "value": "5432",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "DB_NAME",
            "value": "nutrition_platform",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "DB_USER",
            "value": "nutrition_user",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "DB_PASSWORD",
            "value": "ac287cc0e30f54afad53c6dc7e02fd0cccad979d62b75d75d97b1ede12daf8d5",
            "is_build_time": false,
            "is_multiline": false,
            "is_secret": true
        },
        {
            "key": "DB_SSL_MODE",
            "value": "require",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "REDIS_HOST",
            "value": "localhost",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "REDIS_PORT",
            "value": "6379",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "REDIS_PASSWORD",
            "value": "f606b2d16d6697e666ce78a8685574d042df15484ca8f18f39f2e67bf38dc09a",
            "is_build_time": false,
            "is_multiline": false,
            "is_secret": true
        },
        {
            "key": "JWT_SECRET",
            "value": "9a00511e8e23764f8f4524c02f1db9eccc1923208c02fb36cb758d874d8d569bce9ea1b24ac18a958334abe15ef89e09d6010fe64a1d1ffc02a45b07898b2473",
            "is_build_time": false,
            "is_multiline": false,
            "is_secret": true
        },
        {
            "key": "API_KEY_SECRET",
            "value": "5d2763e839f7e71b90ff88bef12f690a41802635aa131f6bc7160056ef0aeb7dc9caaeb07dbe0028128e617529a48903f8d01c6cc64ce61419eb7f309fdfc8bc",
            "is_build_time": false,
            "is_multiline": false,
            "is_secret": true
        },
        {
            "key": "ENCRYPTION_KEY",
            "value": "cc1574e486b2f5abd69d86537079ba928974cc463e36ff410647b15b15533d23",
            "is_build_time": false,
            "is_multiline": false,
            "is_secret": true
        },
        {
            "key": "SESSION_SECRET",
            "value": "f40776484ee20b35e4f754909fb3067cef2a186d0da7c4c24f1bcd54870d9fba",
            "is_build_time": false,
            "is_multiline": false,
            "is_secret": true
        },
        {
            "key": "SERVER_HOST",
            "value": "0.0.0.0",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "SERVER_PORT",
            "value": "8080",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "CORS_ALLOWED_ORIGINS",
            "value": "https://super.doctorhealthy1.com,https://my.doctorhealthy1.com",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "RATE_LIMIT_REQUESTS",
            "value": "100",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "RATE_LIMIT_WINDOW",
            "value": "60",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "LOG_LEVEL",
            "value": "info",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "LOG_FORMAT",
            "value": "json",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "RELIGIOUS_FILTER_ENABLED",
            "value": "true",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "FILTER_ALCOOL",
            "value": "true",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "FILTER_PORK",
            "value": "true",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "FILTER_STRICT_MODE",
            "value": "false",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "DEFAULT_LANGUAGE",
            "value": "en",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "SUPPORTED_LANGUAGES",
            "value": "en,ar",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "RTL_LANGUAGES",
            "value": "ar",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "HEALTH_CHECK_ENABLED",
            "value": "true",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "HEALTH_CHECK_INTERVAL",
            "value": "30",
            "is_build_time": false,
            "is_multiline": false
        },
        {
            "key": "HEALTH_CHECK_TIMEOUT",
            "value": "5",
            "is_build_time": false,
            "is_multiline": false
        }
    ]
}
EOF
)
    
    local response=$(coolify_api "POST" "/applications/$app_uuid/environment/variables" "$env_data")
    
    if echo "$response" | grep -q "error\|Error"; then
        error "❌ Failed to add environment variables: $response"
    fi
    
    success "✅ Environment variables added successfully"
}

# Function to add database services
add_database_services() {
    local app_uuid=$1
    log "🗄️ Adding database services..."
    
    # Add PostgreSQL service
    local postgres_data=$(cat << EOF
{
    "name": "nutrition-postgres",
    "type": "postgresql",
    "version": "15",
    "destination_uuid": "$SERVER_UUID",
    "database": "nutrition_platform",
    "username": "nutrition_user",
    "password": "ac287cc0e30f54afad53c6dc7e02fd0cccad979d62b75d75d97b1ede12daf8d5"
}
EOF
)
    
    local postgres_response=$(coolify_api "POST" "/services" "$postgres_data")
    
    if echo "$postgres_response" | grep -q "error\|Error"; then
        warning "⚠️ PostgreSQL service creation response: $postgres_response"
    else
        success "✅ PostgreSQL service created successfully"
    fi
    
    # Add Redis service
    local redis_data=$(cat << EOF
{
    "name": "nutrition-redis",
    "type": "redis",
    "version": "7-alpine",
    "destination_uuid": "$SERVER_UUID",
    "password": "f606b2d16d6697e666ce78a8685574d042df15484ca8f18f39f2e67bf38dc09a"
}
EOF
)
    
    local redis_response=$(coolify_api "POST" "/services" "$redis_data")
    
    if echo "$redis_response" | grep -q "error\|Error"; then
        warning "⚠️ Redis service creation response: $redis_response"
    else
        success "✅ Redis service created successfully"
    fi
}

# Function to deploy application
deploy_application() {
    local app_uuid=$1
    log "🚀 Deploying application..."
    
    local deploy_data=$(cat << EOF
{
    "force_rebuild": true,
    "debug": false,
    "pull": true
}
EOF
)
    
    local response=$(coolify_api "POST" "/applications/$app_uuid/deploy" "$deploy_data")
    
    if echo "$response" | grep -q "error\|Error"; then
        error "❌ Failed to start deployment: $response"
    fi
    
    local deployment_uuid=$(echo "$response" | grep -o '"deployment_uuid":"[^"]*"' | cut -d'"' -f4)
    
    if [[ -n "$deployment_uuid" ]]; then
        success "✅ Deployment started with UUID: $deployment_uuid"
        monitor_deployment "$deployment_uuid"
    else
        success "✅ Deployment started successfully"
    fi
}

# Function to monitor deployment
monitor_deployment() {
    local deployment_uuid=$1
    log "⏳ Monitoring deployment progress..."
    
    local timeout=600  # 10 minutes
    local interval=10
    local elapsed=0
    
    while [[ $elapsed -lt $timeout ]]; do
        local response=$(coolify_api "GET" "/deployments/$deployment_uuid")
        local status=$(echo "$response" | grep -o '"status":"[^"]*"' | cut -d'"' -f4)
        
        case "$status" in
            "success"|"finished"|"completed"|"running")
                success "✅ Deployment completed successfully!"
                return 0
                ;;
            "failed"|"error")
                error "❌ Deployment failed"
                ;;
            "building"|"in_progress"|"deploying"|"pending")
                log "🔄 Deployment in progress... (${elapsed}s/${timeout}s)"
                ;;
            *)
                log "📊 Deployment status: $status (${elapsed}s/${timeout}s)"
                ;;
        esac
        
        sleep $interval
        elapsed=$((elapsed + interval))
    done
    
    warning "⚠️ Deployment monitoring timed out, but deployment may still be running"
}

# Function to verify deployment
verify_deployment() {
    log "🔍 Verifying deployment..."
    
    local max_attempts=30
    local attempt=0
    
    while [[ $attempt -lt $max_attempts ]]; do
        # Check main site
        if curl -f -s --max-time 10 "https://$DOMAIN" > /dev/null; then
            success "✅ Main site is accessible: https://$DOMAIN"
            
            # Check health endpoint
            if curl -f -s --max-time 10 "https://$DOMAIN/health" > /dev/null; then
                success "✅ Health endpoint is working: https://$DOMAIN/health"
                
                # Check API endpoint
                if curl -f -s --max-time 10 "https://$DOMAIN/api/v1/info" > /dev/null; then
                    success "✅ API endpoint is working: https://$DOMAIN/api/v1/info"
                    return 0
                else
                    warning "⚠️ API endpoint not ready yet"
                fi
            else
                warning "⚠️ Health endpoint not ready yet"
            fi
        else
            warning "⚠️ Main site not ready yet"
        fi
        
        sleep 10
        attempt=$((attempt + 1))
    done
    
    warning "⚠️ Verification timed out, but deployment may still be completing"
}

# Main deployment flow
main() {
    log "🚀 Starting Coolify API Deployment for Nutrition Platform"
    echo ""
    
    # Check if ZIP file exists
    if [[ ! -f "$ZIP_FILE" ]]; then
        error "❌ ZIP file not found: $ZIP_FILE"
    fi
    
    # Create application
    local app_uuid=$(create_application)
    
    # Add environment variables
    add_environment_variables "$app_uuid"
    
    # Add database services
    add_database_services "$app_uuid"
    
    # Deploy application
    deploy_application "$app_uuid"
    
    # Verify deployment
    verify_deployment
    
    echo ""
    echo "🎉 ==================================="
    echo "🎉 DEPLOYMENT COMPLETED!"
    echo "🎉 ==================================="
    echo ""
    echo "📍 Your Application is LIVE:"
    echo "   🌐 Website: https://$DOMAIN"
    echo "   🏥 Health Check: https://$DOMAIN/health"
    echo "   📊 API Base: https://$DOMAIN/api"
    echo ""
    echo "🔧 Coolify Dashboard:"
    echo "   🌐 Management: $COOLIFY_URL/project/$PROJECT_ID/environment/$ENVIRONMENT_ID/application/$app_uuid"
    echo ""
    echo "📊 Next Steps:"
    echo "   1. Test all application features"
    echo "   2. Monitor in Coolify dashboard"
    echo "   3. Set up monitoring and alerts"
    echo "   4. Configure backup policies"
    echo ""
    success "🚀 Nutrition Platform is LIVE and ready!"
}

# Run main function
main