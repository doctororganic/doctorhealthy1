#!/bin/bash
# Quick Setup Script - Runs everything in one go
# Usage: ./quick-setup.sh

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🚀 Quick Setup - Nutrition Platform${NC}"
echo "=================================="

# 1. Check Go installation
echo -e "\n${YELLOW}1. Checking Go installation...${NC}"
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi
echo "✅ Go version: $(go version)"

# 2. Install dependencies
echo -e "\n${YELLOW}2. Installing Go dependencies...${NC}"
go mod download
go mod tidy
echo "✅ Dependencies installed"

# 3. Run migrations
echo -e "\n${YELLOW}3. Running database migrations...${NC}"
if [ -f "run_migrations.sh" ]; then
    chmod +x run_migrations.sh
    ./run_migrations.sh
else
    echo "⚠️  Migration script not found, skipping..."
fi

# 4. Build application
echo -e "\n${YELLOW}4. Building application...${NC}"
go build -o bin/nutrition-platform ./main.go
echo "✅ Build successful"

# 5. Create .env file if it doesn't exist
echo -e "\n${YELLOW}5. Setting up environment...${NC}"
if [ ! -f ".env" ]; then
    if [ -f "env.example" ]; then
        cp env.example .env
        echo "✅ Created .env from env.example"
        echo "⚠️  Please edit .env with your configuration"
    else
        echo "⚠️  env.example not found, creating basic .env..."
        cat > .env << EOF
PORT=8080
ENVIRONMENT=development
JWT_SECRET=$(openssl rand -base64 32)
DB_PATH=./nutrition-platform.db
EOF
        echo "✅ Created basic .env"
    fi
else
    echo "✅ .env already exists"
fi

# 6. Create necessary directories
echo -e "\n${YELLOW}6. Creating directories...${NC}"
mkdir -p uploads logs backups
echo "✅ Directories created"

echo -e "\n${GREEN}✅ Setup complete!${NC}"
echo -e "\nNext steps:"
echo "  1. Edit .env file with your configuration"
echo "  2. Run: go run main.go"
echo "  3. Test: curl http://localhost:8080/health"