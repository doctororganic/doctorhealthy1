# Nutrition Platform - Vercel Deployment Script
# Email: ieltspass111@gmail.com

# Set error action preference to stop on errors
$ErrorActionPreference = "Stop"

Write-Host "🚀 Starting Vercel Deployment for Nutrition Platform" -ForegroundColor Cyan
Write-Host "=================================================" -ForegroundColor Cyan

# Check if we're in the right directory
if (-not (Test-Path "vercel.json")) {
    Write-Host "❌ Error: vercel.json not found. Please run this script from the nutrition-platform directory." -ForegroundColor Red
    exit 1
}

# Check if Node.js/npm is installed
try {
    & npx --version | Out-Null
} catch {
    Write-Host "❌ Error: Node.js/npm not found. Please install Node.js first." -ForegroundColor Red
    exit 1
}

Write-Host "✅ Environment check passed" -ForegroundColor Green

# Validate project files
Write-Host "🔍 Validating project files..." -ForegroundColor Yellow

# Check JavaScript files
Write-Host "  - Checking JavaScript files..." -ForegroundColor Yellow
try {
    Get-ChildItem -Path "./frontend" -Filter "*.js" -Recurse | ForEach-Object {
        & node -c $_.FullName
    }
} catch {
    Write-Host "❌ JavaScript validation failed" -ForegroundColor Red
    exit 1
}

# Check JSON files
Write-Host "  - Checking JSON files..." -ForegroundColor Yellow
try {
    Get-ChildItem -Filter "*.json" -Recurse | ForEach-Object {
        $jsonContent = Get-Content $_.FullName -Raw
        $jsonObject = ConvertFrom-Json $jsonContent
    }
} catch {
    Write-Host "❌ JSON validation failed" -ForegroundColor Red
    exit 1
}

Write-Host "✅ File validation passed" -ForegroundColor Green

# Check Vercel authentication
Write-Host "🔐 Checking Vercel authentication..." -ForegroundColor Magenta
try {
    $user = & npx vercel whoami
    Write-Host "✅ Already logged in to Vercel" -ForegroundColor Green
    Write-Host "   Logged in as: $user" -ForegroundColor Green
} catch {
    Write-Host "⚠️  Not logged in to Vercel" -ForegroundColor Yellow
    Write-Host "📧 Please login with: ieltspass111@gmail.com" -ForegroundColor Cyan
    Write-Host "" -ForegroundColor Cyan
    Write-Host "🔑 Running Vercel login..." -ForegroundColor Cyan
    Write-Host "   1. Select 'Continue with Email'" -ForegroundColor Cyan
    Write-Host "   2. Enter: ieltspass111@gmail.com" -ForegroundColor Cyan
    Write-Host "   3. Check your email for verification link" -ForegroundColor Cyan
    Write-Host "   4. Click the link to complete login" -ForegroundColor Cyan
    Write-Host "" -ForegroundColor Cyan

    & npx vercel login

    # Verify login was successful
    try {
        & npx vercel whoami | Out-Null
        Write-Host "✅ Login successful!" -ForegroundColor Green
    } catch {
        Write-Host "❌ Login failed. Please try again." -ForegroundColor Red
        exit 1
    }
}

# Deploy to production
Write-Host "🚀 Deploying to Vercel production..." -ForegroundColor Cyan
Write-Host "   This may take a few minutes..." -ForegroundColor Cyan
Write-Host "" -ForegroundColor Cyan

# Run deployment
& npx vercel --prod --yes

if ($LASTEXITCODE -eq 0) {
    Write-Host "" -ForegroundColor Green
    Write-Host "🎉 Deployment successful!" -ForegroundColor Green
    Write-Host "=================================================" -ForegroundColor Green
    Write-Host "✅ Your Nutrition Platform is now live!" -ForegroundColor Green
    Write-Host "" -ForegroundColor Green
    Write-Host "📱 Available Features:" -ForegroundColor Green
    Write-Host "   • Personalized Nutrition Planning" -ForegroundColor Green
    Write-Host "   • Diet Plan Generation" -ForegroundColor Green
    Write-Host "   • Workout Recommendations" -ForegroundColor Green
    Write-Host "   • Medical Condition Support" -ForegroundColor Green
    Write-Host "   • System Validation Dashboard" -ForegroundColor Green
    Write-Host "" -ForegroundColor Green
    Write-Host "🔗 Access your app at the URL provided above" -ForegroundColor Green
    Write-Host "📊 View deployment details in Vercel dashboard" -ForegroundColor Green
    Write-Host "" -ForegroundColor Green
} else {
    Write-Host "❌ Deployment failed" -ForegroundColor Red
    Write-Host "💡 Troubleshooting:" -ForegroundColor Yellow
    Write-Host "   1. Check your internet connection" -ForegroundColor Yellow
    Write-Host "   2. Verify Vercel authentication: npx vercel whoami" -ForegroundColor Yellow
    Write-Host "   3. Review the error messages above" -ForegroundColor Yellow
    Write-Host "   4. Try running: npx vercel --prod manually" -ForegroundColor Yellow
    exit 1
}