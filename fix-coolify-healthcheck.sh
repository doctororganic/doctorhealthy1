#!/bin/bash

echo "🔧 Quick Coolify Health Check Fix"
echo "================================="

# Option 1: Disable health check for static nginx service
echo "Option 1: Disable health check (for static sites)"
echo "Go to: Coolify Dashboard → Services → doctorhealthy1-api → Settings → Health Check → Disable"

echo ""
echo "Option 2: Fix for Go application"
echo "If this should be a Go application, update the service configuration:"
echo "1. Change 'static_image' from 'nginx:alpine' to null/empty"
echo "2. Ensure Dockerfile exists in repository root"
echo "3. Redeploy the service"

echo ""
echo "Option 3: Quick manual restart"
echo "SSH into server and restart containers:"
echo "docker restart <container_name>"

echo ""
echo "Current Status from MCP Server:"
echo "✅ Server: Running (128.140.111.171)"
echo "✅ Traefik: Running"
echo "❌ Main API: running:unhealthy (nginx static server)"
echo "❌ Test Service: exited:unhealthy"

echo ""
echo "🔧 Recommended Actions:"
echo "1. Access Coolify dashboard"
echo "2. Go to doctorhealthy1-api service"
echo "3. Edit service configuration"
echo "4. Disable health check OR fix Dockerfile configuration"
echo "5. Redeploy service"
echo ""
echo "📝 Service URLs:"
echo "Main API: https://my.doctorhealthy1.com"
echo "Test Service: http://vc40okgosockc8c8sgoko0k0.128.140.111.171.sslip.io"

