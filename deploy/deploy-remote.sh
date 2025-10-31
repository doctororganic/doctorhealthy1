#!/bin/bash

# Stop existing container
docker stop nutrition-platform 2>/dev/null || true
docker rm nutrition-platform 2>/dev/null || true

# Pull and run new container
docker run -d \\
    --name nutrition-platform \\
    -p 80:80 \\
    -p 443:443 \\
    -v \$(pwd)/nginx.conf:/etc/nginx/nginx.conf \\
    -v \$(pwd)/app:/app \\
    nginx:alpine

echo "✅ Deployment completed!"
