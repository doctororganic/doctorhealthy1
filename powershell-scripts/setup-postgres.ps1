# Setup PostgreSQL database on Fly.io for Nutrition Platform

# Set error action preference to stop on errors
$ErrorActionPreference = "Stop"

# Colors for output
$Green = "Green"
$Yellow = "Yellow"
$Red = "Red"
$NC = "White"  # No Color

Write-Host "# Setting up PostgreSQL database on Fly.io..." -ForegroundColor $Yellow

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
$POSTGRES_PASSWORD = $env:POSTGRES_PASSWORD
if (-not $POSTGRES_PASSWORD) {
    $POSTGRES_PASSWORD = & openssl rand -base64 16
    Write-Host "# Generated PostgreSQL password: $($POSTGRES_PASSWORD)" -ForegroundColor $Yellow
    Write-Host "# Please save this password securely!" -ForegroundColor $Yellow
    $env:POSTGRES_PASSWORD = $POSTGRES_PASSWORD
}

# Check if the app exists on Fly.io
try {
    & flyctl status --app nutrition-platform-db | Out-Null
    $appExists = $true
} catch {
    $appExists = $false
}

if (-not $appExists) {
    Write-Host "# Creating PostgreSQL database on Fly.io..." -ForegroundColor $Yellow

    # Create volume for persistent data
    Write-Host "# Creating volume for PostgreSQL data..." -ForegroundColor $Yellow
    & flyctl volumes create pg_data --size 10 --app nutrition-platform-db

    # Deploy PostgreSQL
    Write-Host "# Deploying PostgreSQL to Fly.io..." -ForegroundColor $Yellow
    & flyctl deploy --config fly-postgres.toml --app nutrition-platform-db
} else {
    Write-Host "# PostgreSQL database already exists on Fly.io." -ForegroundColor $Yellow
    Write-Host "# Updating configuration..." -ForegroundColor $Yellow
    & flyctl deploy --config fly-postgres.toml --app nutrition-platform-db
}

# Get connection string
Write-Host "# Getting PostgreSQL connection details..." -ForegroundColor $Yellow
$statusOutput = & flyctl status --app nutrition-platform-db
$DB_HOST = ($statusOutput | Select-String -Pattern 'v4: (\d+\.\d+\.\d+\.\d+)' | ForEach-Object { $_.Matches.Groups[1].Value })[0]

Write-Host "# PostgreSQL database setup completed!" -ForegroundColor $Green
Write-Host "Connection details:" -ForegroundColor $Green
Write-Host "  Host: $DB_HOST" -ForegroundColor $Green
Write-Host "  Port: 5432" -ForegroundColor $Green
Write-Host "  Database: nutrition_platform" -ForegroundColor $Green
Write-Host "  Username: nutrition_user" -ForegroundColor $Green
Write-Host "  Password: $POSTGRES_PASSWORD" -ForegroundColor $Green

# Update the main app configuration
Write-Host "# Updating main application configuration..." -ForegroundColor $Yellow
$envContent = @"
DB_HOST=$DB_HOST
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=nutrition_user
DB_PASSWORD=$POSTGRES_PASSWORD
DB_SSL_MODE=require
"@

Set-Content -Path ".env.production" -Value $envContent

Write-Host "# Environment configuration saved to .env.production" -ForegroundColor $Green
Write-Host "You can now deploy the main application using ./deploy-to-fly.sh" -ForegroundColor $Green