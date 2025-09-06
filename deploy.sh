#!/bin/bash

# Nutrition Platform - Vercel Deployment Script
# Email: ieltspass111@gmail.com

set -e

echo "🚀 Starting Vercel Deployment for Nutrition Platform"
echo "================================================="

# Check if we're in the right directory
if [ ! -f "vercel.json" ]; then
    echo "❌ Error: vercel.json not found. Please run this script from the nutrition-platform directory."
    exit 1
fi

# Check if Vercel CLI is installed
if ! command -v npx &> /dev/null; then
    echo "❌ Error: Node.js/npm not found. Please install Node.js first."
    exit 1
fi

echo "✅ Environment check passed"

# Validate project files
echo "🔍 Validating project files..."

# Check JavaScript files
echo "  - Checking JavaScript files..."
find ./frontend -name '*.js' -exec node -c {} \; 2>/dev/null || {
    echo "❌ JavaScript validation failed"
    exit 1
}

# Check JSON files
echo "  - Checking JSON files..."
find . -name '*.json' -exec node -e "JSON.parse(require('fs').readFileSync('{}', 'utf8'))" \; 2>/dev/null || {
    echo "❌ JSON validation failed"
    exit 1
}

echo "✅ File validation passed"

# Check Vercel authentication
echo "🔐 Checking Vercel authentication..."
if npx vercel whoami &>/dev/null; then
    echo "✅ Already logged in to Vercel"
    USER=$(npx vercel whoami)
    echo "   Logged in as: $USER"
else
    echo "⚠️  Not logged in to Vercel"
    echo "📧 Please login with: ieltspass111@gmail.com"
    echo ""
    echo "🔑 Running Vercel login..."
    echo "   1. Select 'Continue with Email'"
    echo "   2. Enter: ieltspass111@gmail.com"
    echo "   3. Check your email for verification link"
    echo "   4. Click the link to complete login"
    echo ""
    
    npx vercel login
    
    # Verify login was successful
    if npx vercel whoami &>/dev/null; then
        echo "✅ Login successful!"
    else
        echo "❌ Login failed. Please try again."
        exit 1
    fi
fi

# Deploy to production
echo "🚀 Deploying to Vercel production..."
echo "   This may take a few minutes..."
echo ""

# Run deployment
npx vercel --prod --yes

if [ $? -eq 0 ]; then
    echo ""
    echo "🎉 Deployment successful!"
    echo "================================================="
    echo "✅ Your Nutrition Platform is now live!"
    echo ""
    echo "📱 Available Features:"
    echo "   • Personalized Nutrition Planning"
    echo "   • Diet Plan Generation"
    echo "   • Workout Recommendations"
    echo "   • Medical Condition Support"
    echo "   • System Validation Dashboard"
    echo ""
    echo "🔗 Access your app at the URL provided above"
    echo "📊 View deployment details in Vercel dashboard"
    echo ""
else
    echo "❌ Deployment failed"
    echo "💡 Troubleshooting:"
    echo "   1. Check your internet connection"
    echo "   2. Verify Vercel authentication: npx vercel whoami"
    echo "   3. Review the error messages above"
    echo "   4. Try running: npx vercel --prod manually"
    exit 1
fi