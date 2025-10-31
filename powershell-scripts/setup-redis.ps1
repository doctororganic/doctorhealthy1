# Setup Redis caching on Fly.io for Nutrition Platform

# Set error action preference to stop on errors
$ErrorActionPreference = "Stop"

# Colors for output
$Green = "Green"
$Yellow = "Yellow"
$Red = "Red"
$NC = "White"  # No Color

Write-Host "# Setting up Redis caching on Fly.io..." -ForegroundColor $Yellow

# Check if flyctl is installed
try {
    & flyctl version | Out-Null
} catch {
    Write-Host "# Error: flyctl is not installed. Please install it first." -ForegroundColor $Red
    Write-Host "Visit https://fly.io/docs/hands-on/install-flyctl/ for installation instructions." -ForegroundColor $Red
    exit 1
}

# Check if user is authenticated
try {
    & flyctl auth whoami | Out-Null
} catch {
    Write-Host "# You are not logged in to Fly.io. Please log in first." -ForegroundColor $Yellow
    & flyctl auth login
}

# Generate a secure password if not provided
$REDIS_PASSWORD = $env:REDIS_PASSWORD
if (-not $REDIS_PASSWORD) {
    $REDIS_PASSWORD = & openssl rand -base64 16
    Write-Host "# Generated Redis password: $($REDIS_PASSWORD)" -ForegroundColor $Yellow
    Write-Host "# Please save this password securely!" -ForegroundColor $Yellow
    $env:REDIS_PASSWORD = $REDIS_PASSWORD
}

# Check if the app exists on Fly.io
try {
    & flyctl status --app nutrition-platform-redis | Out-Null
    $appExists = $true
} catch {
    $appExists = $false
}

if (-not $appExists) {
    Write-Host "# Creating Redis instance on Fly.io..." -ForegroundColor $Yellow

    # Create volume for persistent data
    Write-Host "# Creating volume for Redis data..." -ForegroundColor $Yellow
    & flyctl volumes create redis_data --size 1 --app nutrition-platform-redis

    # Deploy Redis
    Write-Host "# Deploying Redis to Fly.io..." -ForegroundColor $Yellow
    & flyctl deploy --config fly-redis.toml --app nutrition-platform-redis
} else {
    Write-Host "# Redis instance already exists on Fly.io." -ForegroundColor $Yellow
    Write-Host "# Updating configuration..." -ForegroundColor $Yellow
    & flyctl deploy --config fly-redis.toml --app nutrition-platform-redis
}

# Get connection string
Write-Host "# Getting Redis connection details..." -ForegroundColor $Yellow
$REDIS_HOST = (& flyctl status --app nutrition-platform-redis | Select-String -Pattern 'v4: (\d+\.\d+\.\d+\.\d+)' | ForEach-Object { $_.Matches.Groups[1].Value })[0]

Write-Host "# Redis setup completed!" -ForegroundColor $Green
Write-Host "Connection details:" -ForegroundColor $Green
Write-Host "  Host: $REDIS_HOST" -ForegroundColor $Green
Write-Host "  Port: 6379" -ForegroundColor $Green
Write-Host "  Password: $REDIS_PASSWORD" -ForegroundColor $Green

# Update the main app configuration
Write-Host "# Updating main application configuration..." -ForegroundColor $Yellow

# Check if .env.production exists
$envExists = Test-Path ".env.production"

if ($envExists) {
    # Add Redis configuration to existing file
    Write-Host "# Adding Redis configuration to .env.production..." -ForegroundColor $Yellow
    $redisConfig = @"

# Redis Configuration
REDIS_HOST=$REDIS_HOST
REDIS_PORT=6379
REDIS_PASSWORD=$REDIS_PASSWORD
REDIS_DB=0
REDIS_ENABLED=true
"@
    Add-Content -Path ".env.production" -Value $redisConfig
} else {
    # Create new file with Redis configuration
    Write-Host "# Creating .env.production with Redis configuration..." -ForegroundColor $Yellow
    $envContent = @"
# Redis Configuration
REDIS_HOST=$REDIS_HOST
REDIS_PORT=6379
REDIS_PASSWORD=$REDIS_PASSWORD
REDIS_DB=0
REDIS_ENABLED=true
"@
    Set-Content -Path ".env.production" -Value $envContent
}

Write-Host "# Redis configuration saved to .env.production" -ForegroundColor $Green
Write-Host "You can now deploy the main application using ./deploy-to-fly.sh" -ForegroundColor $Green