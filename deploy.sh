#!/bin/bash

echo "🚀 Deploying Nutrition Platform..."

# Build and start services
docker-compose down
docker-compose build
docker-compose up -d

echo ""
echo "✅ Deployment complete!"
echo ""
echo "Services:"
echo "  Backend:  http://localhost:8080"
echo "  Frontend: http://localhost:3000"
echo "  Health:   http://localhost:8080/health"
echo ""
echo "Check status: docker-compose ps"
echo "View logs:    docker-compose logs -f"
