#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}📊 LIVE MONITORING DASHBOARD${NC}"
echo "=============================="

while true; do
    clear
    echo -e "${BLUE}📊 LIVE MONITORING - $(date)${NC}"
    echo "=============================="
    
    # Container status
    echo -e "\n${BLUE}🐳 Container Status:${NC}"
    docker-compose -f docker-compose.production.yml ps
    
    # Health checks
    echo -e "\n${BLUE}💚 Health Checks:${NC}"
    
    # Backend health
    if curl -sf http://localhost:8080/health >/dev/null 2>&1; then
        echo -e "Backend:  ${GREEN}✅ Healthy${NC}"
    else
        echo -e "Backend:  ${RED}❌ Down${NC}"
    fi
    
    # Frontend health
    if curl -sf http://localhost:3000 >/dev/null 2>&1; then
        echo -e "Frontend: ${GREEN}✅ Healthy${NC}"
    else
        echo -e "Frontend: ${RED}❌ Down${NC}"
    fi
    
    # Database health
    if docker-compose -f docker-compose.production.yml exec -T postgres pg_isready -U nutrition_user >/dev/null 2>&1; then
        echo -e "Database: ${GREEN}✅ Healthy${NC}"
    else
        echo -e "Database: ${RED}❌ Down${NC}"
    fi
    
    # Redis health
    if docker-compose -f docker-compose.production.yml exec -T redis redis-cli ping >/dev/null 2>&1; then
        echo -e "Redis:    ${GREEN}✅ Healthy${NC}"
    else
        echo -e "Redis:    ${RED}❌ Down${NC}"
    fi
    
    # Resource usage
    echo -e "\n${BLUE}💻 Resource Usage:${NC}"
    docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}" | head -6
    
    # Recent logs
    echo -e "\n${BLUE}📝 Recent Backend Logs:${NC}"
    docker-compose -f docker-compose.production.yml logs --tail=5 backend 2>/dev/null | tail -5
    
    echo -e "\n${YELLOW}Press Ctrl+C to exit${NC}"
    sleep 5
done
