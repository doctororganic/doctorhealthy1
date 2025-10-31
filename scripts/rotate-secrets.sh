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
