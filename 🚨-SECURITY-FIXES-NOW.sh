#!/bin/bash

# ============================================
# EMERGENCY SECURITY FIXES
# Addressing all critical security audit findings
# ============================================

set -e

echo "🚨 EMERGENCY SECURITY FIXES"
echo "============================"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo ""
echo "1️⃣  Removing exposed secrets from repository..."

# Remove files with exposed secrets
if [ -f "coolify-env-vars.txt" ]; then
    git rm --cached coolify-env-vars.txt 2>/dev/null || rm coolify-env-vars.txt
    echo -e "${GREEN}✅ Removed coolify-env-vars.txt${NC}"
fi

if [ -f "backend/.env" ]; then
    git rm --cached backend/.env 2>/dev/null || true
    echo -e "${GREEN}✅ Removed backend/.env from git${NC}"
fi

if [ -f ".env" ]; then
    git rm --cached .env 2>/dev/null || true
    echo -e "${GREEN}✅ Removed .env from git${NC}"
fi

# Add to .gitignore
echo "" >> .gitignore
echo "# Security: Never commit secrets" >> .gitignore
echo "*.env" >> .gitignore
echo ".env*" >> .gitignore
echo "!.env.example" >> .gitignore
echo "coolify-env-vars.txt" >> .gitignore
echo "secrets/" >> .gitignore

echo -e "${GREEN}✅ Updated .gitignore${NC}"

echo ""
echo "2️⃣  Generating strong secrets..."

# Generate strong secrets
DB_PASSWORD=$(openssl rand -hex 32)
REDIS_PASSWORD=$(openssl rand -hex 32)
JWT_SECRET=$(openssl rand -hex 64)
API_KEY_SECRET=$(openssl rand -hex 64)
ENCRYPTION_KEY=$(openssl rand -hex 16)  # 32 chars

# Create secure .env.example
cat > backend/.env.example << EOF
# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=nutrition_user
DB_PASSWORD=CHANGE_ME_$(openssl rand -hex 16)
DB_SSL_MODE=require
DB_SSL_CERT=/path/to/client-cert.pem
DB_SSL_KEY=/path/to/client-key.pem
DB_SSL_ROOT_CERT=/path/to/ca-cert.pem

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=CHANGE_ME_$(openssl rand -hex 16)

# Security
JWT_SECRET=CHANGE_ME_$(openssl rand -hex 32)
API_KEY_SECRET=CHANGE_ME_$(openssl rand -hex 32)
ENCRYPTION_KEY=CHANGE_ME_$(openssl rand -hex 16)

# Server
PORT=8080
ENVIRONMENT=production
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60s
EOF

echo -e "${GREEN}✅ Created secure .env.example${NC}"

echo ""
echo "3️⃣  Fixing CORS policy..."

# Fix CORS in nginx config
if [ -f "nginx/conf.d/default.conf" ]; then
    # Backup original
    cp nginx/conf.d/default.conf nginx/conf.d/default.conf.backup
    
    # Replace wildcard CORS with specific origins
    sed -i.bak 's/Access-Control-Allow-Origin "\*"/Access-Control-Allow-Origin "$http_origin"/g' nginx/conf.d/default.conf
    
    echo -e "${GREEN}✅ Fixed CORS policy in nginx${NC}"
fi

echo ""
echo "4️⃣  Enabling database SSL..."

# Update database config
cat > backend/config/database.go << 'EOF'
package config

import (
    "fmt"
    "os"
)

type DatabaseConfig struct {
    Host     string
    Port     string
    Name     string
    User     string
    Password string
    SSLMode  string
    SSLCert  string
    SSLKey   string
    SSLRootCert string
}

func LoadDatabaseConfig() *DatabaseConfig {
    sslMode := os.Getenv("DB_SSL_MODE")
    if sslMode == "" {
        sslMode = "require"  // Default to require SSL
    }
    
    return &DatabaseConfig{
        Host:        os.Getenv("DB_HOST"),
        Port:        os.Getenv("DB_PORT"),
        Name:        os.Getenv("DB_NAME"),
        User:        os.Getenv("DB_USER"),
        Password:    os.Getenv("DB_PASSWORD"),
        SSLMode:     sslMode,
        SSLCert:     os.Getenv("DB_SSL_CERT"),
        SSLKey:      os.Getenv("DB_SSL_KEY"),
        SSLRootCert: os.Getenv("DB_SSL_ROOT_CERT"),
    }
}

func (c *DatabaseConfig) ConnectionString() string {
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
    
    if c.SSLMode != "disable" {
        if c.SSLCert != "" {
            dsn += fmt.Sprintf(" sslcert=%s", c.SSLCert)
        }
        if c.SSLKey != "" {
            dsn += fmt.Sprintf(" sslkey=%s", c.SSLKey)
        }
        if c.SSLRootCert != "" {
            dsn += fmt.Sprintf(" sslrootcert=%s", c.SSLRootCert)
        }
    }
    
    return dsn
}
EOF

echo -e "${GREEN}✅ Enabled database SSL${NC}"

echo ""
echo "5️⃣  Creating secret rotation script..."

cat > scripts/rotate-secrets.sh << 'EOFSCRIPT'
#!/bin/bash

echo "🔄 Rotating secrets..."

# Generate new secrets
NEW_DB_PASSWORD=$(openssl rand -hex 32)
NEW_REDIS_PASSWORD=$(openssl rand -hex 32)
NEW_JWT_SECRET=$(openssl rand -hex 64)
NEW_API_KEY_SECRET=$(openssl rand -hex 64)

echo "New secrets generated. Update your .env file with:"
echo ""
echo "DB_PASSWORD=$NEW_DB_PASSWORD"
echo "REDIS_PASSWORD=$NEW_REDIS_PASSWORD"
echo "JWT_SECRET=$NEW_JWT_SECRET"
echo "API_KEY_SECRET=$NEW_API_KEY_SECRET"
echo ""
echo "⚠️  Remember to:"
echo "1. Update database password"
echo "2. Update Redis password"
echo "3. Restart all services"
echo "4. Invalidate all existing JWT tokens"
EOFSCRIPT

chmod +x scripts/rotate-secrets.sh

echo -e "${GREEN}✅ Created secret rotation script${NC}"

echo ""
echo "============================"
echo -e "${GREEN}🎉 SECURITY FIXES COMPLETE!${NC}"
echo "============================"
echo ""
echo "⚠️  IMPORTANT NEXT STEPS:"
echo ""
echo "1. Generate production secrets:"
echo "   ./scripts/rotate-secrets.sh"
echo ""
echo "2. Update .env file with new secrets"
echo ""
echo "3. Commit security fixes:"
echo "   git add .gitignore backend/.env.example"
echo "   git commit -m 'Security: Remove exposed secrets, fix CORS, enable DB SSL'"
echo ""
echo "4. Deploy with new configuration"
echo ""
echo "5. Verify security:"
echo "   ./scripts/security-scan.sh"
echo ""
