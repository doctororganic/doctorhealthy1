#!/bin/bash

echo "🧪 TESTING COOLIFY MCP SETUP"
echo "============================"

# 1. Check MCP installation
echo ""
echo "1️⃣  Checking MCP installation..."
if command -v npx &> /dev/null; then
    echo "✅ npx installed"
    npx --version
else
    echo "❌ npx not found"
    exit 1
fi

# 2. Check Coolify MCP package
echo ""
echo "2️⃣  Checking Coolify MCP package..."
if npx -y coolify-mcp-server --version 2>/dev/null; then
    echo "✅ Coolify MCP server available"
else
    echo "⚠️  Installing Coolify MCP server..."
    npm install -g coolify-mcp-server
fi

# 3. Load credentials
echo ""
echo "3️⃣  Loading credentials..."
if [ -f ".coolify-credentials.enc" ]; then
    source .coolify-credentials.enc
    echo "✅ Credentials loaded"
    echo "   URL: $COOLIFY_BASE_URL"
    echo "   Token: ${COOLIFY_TOKEN:0:20}..."
else
    echo "❌ Credentials file not found"
    exit 1
fi

# 4. Test API connection
echo ""
echo "4️⃣  Testing Coolify API..."
response=$(curl -s -w "\n%{http_code}" \
    -H "Authorization: Bearer $COOLIFY_TOKEN" \
    "$COOLIFY_BASE_URL/api/v1/servers" 2>&1)

http_code=$(echo "$response" | tail -1)
body=$(echo "$response" | head -n -1)

if [ "$http_code" = "200" ]; then
    echo "✅ API connection successful"
    echo "$body" | jq '.' 2>/dev/null || echo "$body"
else
    echo "❌ API connection failed (HTTP $http_code)"
    echo "Response: $body"
    
    # Debug info
    echo ""
    echo "🔍 Debug Information:"
    echo "   URL: $COOLIFY_BASE_URL/api/v1/servers"
    echo "   Token length: ${#COOLIFY_TOKEN}"
    echo "   Token format: ${COOLIFY_TOKEN:0:5}...${COOLIFY_TOKEN: -5}"
fi

# 5. Check MCP config
echo ""
echo "5️⃣  Checking MCP configuration..."
if [ -f "$HOME/.kiro/settings/mcp.json" ]; then
    echo "✅ MCP config found"
    cat "$HOME/.kiro/settings/mcp.json" | jq '.' 2>/dev/null || cat "$HOME/.kiro/settings/mcp.json"
else
    echo "⚠️  MCP config not found at $HOME/.kiro/settings/mcp.json"
fi

# 6. Test deployment readiness
echo ""
echo "6️⃣  Checking deployment readiness..."
checks=0
total=5

[ -f "docker-compose.production.yml" ] && ((checks++)) && echo "✅ docker-compose.production.yml"
[ -f "backend/Dockerfile.secure" ] && ((checks++)) && echo "✅ backend/Dockerfile.secure"
[ -f "frontend-nextjs/Dockerfile.secure" ] && ((checks++)) && echo "✅ frontend-nextjs/Dockerfile.secure"
[ -f ".env.production" ] && ((checks++)) && echo "✅ .env.production"
[ -f "nginx/production.conf" ] && ((checks++)) && echo "✅ nginx/production.conf"

echo ""
echo "📊 Deployment Readiness: $checks/$total"

if [ $checks -eq $total ]; then
    echo "✅ All deployment files present"
else
    echo "⚠️  Some deployment files missing"
fi

# Summary
echo ""
echo "================================"
echo "🎯 TEST SUMMARY"
echo "================================"
if [ "$http_code" = "200" ] && [ $checks -eq $total ]; then
    echo "✅ ALL TESTS PASSED"
    echo ""
    echo "Ready to deploy with:"
    echo "  ./COOLIFY-MCP-DEPLOY.sh"
else
    echo "❌ SOME TESTS FAILED"
    echo ""
    echo "Fix issues above before deploying"
fi
echo ""
