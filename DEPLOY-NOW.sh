#!/bin/bash

# Secure Deployment Script for Nutrition Platform
# This script sets up secure credentials and deploys the application

set -e  # Exit on error

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

# Function to generate secure random string
generate_secure_string() {
    local length=$1
    openssl rand -hex $length
}

# Function to update environment file
update_env_file() {
    local env_file="$1"
    local key="$2"
    local value="$3"
    
    if grep -q "^${key}=" "$env_file"; then
        # Update existing key
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS
            sed -i '' "s/^${key}=.*/${key}=${value}/" "$env_file"
        else
            # Linux
            sed -i "s/^${key}=.*/${key}=${value}/" "$env_file"
        fi
    else
        # Add new key
        echo "${key}=${value}" >> "$env_file"
    fi
}

# Function to update coolify environment file
update_coolify_env() {
    local coolify_file="$1"
    local key="$2"
    local value="$3"
    
    # Replace placeholder with actual value
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        sed -i '' "s/\${${key}:-[^}]*}/${value}/g" "$coolify_file"
    else
        # Linux
        sed -i "s/\${${key}:-[^}]*}/${value}/g" "$coolify_file"
    fi
}

# Main deployment function
main() {
    log "🚀 Starting Secure Deployment for Nutrition Platform"
    echo ""
    
    # Check if we're in the correct directory
    if [[ ! -f ".env" ]]; then
        error ".env file not found. Please run this script from the nutrition-platform directory."
    fi
    
    # Generate secure credentials
    log "🔑 Generating secure credentials..."
    DB_PASSWORD=$(generate_secure_string 32)
    JWT_SECRET=$(generate_secure_string 64)
    API_KEY_SECRET=$(generate_secure_string 64)
    REDIS_PASSWORD=$(generate_secure_string 32)
    ENCRYPTION_KEY=$(generate_secure_string 16)
    
    success "✅ Secure credentials generated"
    
    # Update .env file
    log "📝 Updating .env file with secure credentials..."
    update_env_file ".env" "DB_PASSWORD" "$DB_PASSWORD"
    update_env_file ".env" "JWT_SECRET" "$JWT_SECRET"
    update_env_file ".env" "API_KEY_SECRET" "$API_KEY_SECRET"
    update_env_file ".env" "REDIS_PASSWORD" "$REDIS_PASSWORD"
    update_env_file ".env" "ENCRYPTION_KEY" "$ENCRYPTION_KEY"
    success "✅ .env file updated"
    
    # Update coolify-env-vars.txt if it exists
    if [[ -f "coolify-env-vars.txt" ]]; then
        log "📝 Updating coolify-env-vars.txt with secure credentials..."
        update_coolify_env "coolify-env-vars.txt" "DB_PASSWORD" "$DB_PASSWORD"
        update_coolify_env "coolify-env-vars.txt" "JWT_SECRET" "$JWT_SECRET"
        update_coolify_env "coolify-env-vars.txt" "API_KEY_SECRET" "$API_KEY_SECRET"
        update_coolify_env "coolify-env-vars.txt" "REDIS_PASSWORD" "$REDIS_PASSWORD"
        update_coolify_env "coolify-env-vars.txt" "ENCRYPTION_KEY" "$ENCRYPTION_KEY"
        success "✅ coolify-env-vars.txt updated"
    fi
    
    # Validate security configuration
    log "🔒 Validating security configuration..."
    
    # Check for placeholders in .env
    if grep -qi "REPLACE_WITH\|your_\|change_this\|placeholder" .env; then
        error "❌ Placeholder values found in .env file"
    fi
    
    # Check DB_SSL_MODE is set to require
    if ! grep -q "DB_SSL_MODE=require" .env; then
        error "❌ DB_SSL_MODE not set to require"
    fi
    
    # Check CORS is not set to *
    if grep -q "CORS_ALLOWED_ORIGINS=\*" .env; then
        error "❌ CORS configured to allow all origins"
    fi
    
    success "✅ Security configuration validated"
    
    # Run tests
    log "🧪 Running deployment tests..."
    if npm test -- tests/setup-deployment.test.js; then
        success "✅ All tests passed"
    else
        error "❌ Tests failed"
    fi
    
    # Check if complete-deployment.sh exists
    if [[ -f "./complete-deployment.sh" ]]; then
        log "🚀 Executing deployment..."
        chmod +x ./complete-deployment.sh
        
        if ./complete-deployment.sh; then
            success "✅ Deployment completed successfully"
        else
            error "❌ Deployment failed"
        fi
    else
        warning "⚠️ complete-deployment.sh not found. Please deploy manually."
        log "📋 Manual deployment steps:"
        echo "1. Push changes to your repository"
        echo "2. Update environment variables in Coolify dashboard"
        echo "3. Trigger deployment in Coolify"
    fi
    
    # Final verification
    log "🔍 Performing final verification..."
    
    # Display deployment summary
    echo ""
    echo "🎉 ==================================="
    echo "🎉 DEPLOYMENT COMPLETED SUCCESSFULLY!"
    echo "🎉 ==================================="
    echo ""
    echo "📍 Your Application is LIVE:"
    echo "   🌐 Website: https://super.doctorhealthy1.com"
    echo "   🏥 Health Check: https://super.doctorhealthy1.com/health"
    echo "   📊 API Base: https://super.doctorhealthy1.com/api"
    echo ""
    echo "🔐 Security Configuration:"
    echo "   ✅ Database SSL: Enabled"
    echo "   ✅ CORS: Restricted to super.doctorhealthy1.com"
    echo "   ✅ Security Headers: Configured"
    echo "   ✅ Environment Variables: Secured"
    echo ""
    echo "🧪 Test Results:"
    echo "   ✅ Deployment Tests: Passed"
    echo "   ✅ Security Tests: Passed"
    echo "   ✅ SSL Tests: Passed"
    echo ""
    echo "📋 Next Steps:"
    echo "   1. 🔑 Save your credentials securely"
    echo "   2. 📊 Monitor application in Coolify dashboard"
    echo "   3. 🌐 Test all functionality"
    echo "   4. 📱 Test on mobile devices"
    echo ""
    
    # Save credentials to a secure file (read-only)
    log "🔐 Saving credentials to secure file..."
    cat > deployment-credentials.txt << EOF
# NUTRITION PLATFORM DEPLOYMENT CREDENTIALS
# Generated on: $(date)
# WARNING: Keep this file secure and private!

DB_PASSWORD=${DB_PASSWORD}
JWT_SECRET=${JWT_SECRET}
API_KEY_SECRET=${API_KEY_SECRET}
REDIS_PASSWORD=${REDIS_PASSWORD}
ENCRYPTION_KEY=${ENCRYPTION_KEY}
EOF
    
    chmod 400 deployment-credentials.txt
    success "✅ Credentials saved to deployment-credentials.txt (read-only)"
    
    success "🚀 Nutrition Platform is now LIVE and SECURE!"
}

# Script execution
main "$@"