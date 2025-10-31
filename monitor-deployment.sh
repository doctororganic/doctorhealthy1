#!/bin/bash

echo "🔄 LIVE MONITORING STARTED"
echo "=========================="
echo "Press Ctrl+C to stop"
echo ""

while true; do
    clear
    echo "🔄 LIVE DEPLOYMENT MONITOR - $(date)"
    echo "===================================="
    echo ""
    
    # Container status
    echo "📦 CONTAINER STATUS:"
    docker-compose -f docker-compose.production.yml ps
    echo ""
    
    # Health checks
    echo "🏥 HEALTH CHECKS:"
    
    # Backend
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        echo "✅ Backend: HEALTHY"
    else
        echo "❌ Backend: DOWN"
    fi
    
    # Frontend
    if curl -f http://localhost:3000 > /dev/null 2>&1; then
        echo "✅ Frontend: HEALTHY"
    else
        echo "❌ Frontend: DOWN"
    fi
    
    # PostgreSQL
    if docker-compose -f docker-compose.production.yml exec -T postgres pg_isready -U nutrition_user > /dev/null 2>&1; then
        echo "✅ PostgreSQL: HEALTHY"
    else
        echo "❌ PostgreSQL: DOWN"
    fi
    
    # Redis
    if docker-compose -f docker-compose.production.yml exec -T redis redis-cli ping > /dev/null 2>&1; then
        echo "✅ Redis: HEALTHY"
    else
        echo "❌ Redis: DOWN"
    fi
    
    echo ""
    
    # Resource usage
    echo "💻 RESOURCE USAGE:"
    docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}" | head -6
    echo ""
    
    # Recent errors
    echo "🚨 RECENT ERRORS (last 5):"
    docker-compose -f docker-compose.production.yml logs --tail=100 2>&1 | grep -i "error\|fatal\|panic" | tail -5 || echo "No errors"
    echo ""
    
    echo "Refreshing in 10 seconds..."
    sleep 10
done
