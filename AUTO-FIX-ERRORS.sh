#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔧 AUTO-FIX DEPLOYMENT ERRORS${NC}"
echo "=============================="

# Function to fix common errors
fix_port_conflict() {
    echo -e "${YELLOW}Checking for port conflicts...${NC}"
    for port in 8080 3000 5432 6379; do
        if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
            echo -e "${RED}Port $port is in use${NC}"
            echo "Killing process on port $port..."
            lsof -ti:$port | xargs kill -9 2>/dev/null || true
            echo -e "${GREEN}✅ Port $port freed${NC}"
        fi
    done
}

fix_docker_issues() {
    echo -e "${YELLOW}Checking Docker...${NC}"
    
    # Remove dangling images
    echo "Removing dangling images..."
    docker image prune -f
    
    # Remove stopped containers
    echo "Removing stopped containers..."
    docker container prune -f
    
    # Remove unused volumes
    echo "Removing unused volumes..."
    docker volume prune -f
    
    echo -e "${GREEN}✅ Docker cleaned${NC}"
}

fix_permissions() {
    echo -e "${YELLOW}Fixing file permissions...${NC}"
    chmod +x *.sh 2>/dev/null || true
    chmod +x scripts/*.sh 2>/dev/null || true
    echo -e "${GREEN}✅ Permissions fixed${NC}"
}

fix_env_file() {
    echo -e "${YELLOW}Checking .env.production...${NC}"
    if [ ! -f ".env.production" ]; then
        echo -e "${RED}.env.production missing${NC}"
        echo "Run ./LIVE-TEST-DEPLOY.sh to generate it"
        return 1
    fi
    echo -e "${GREEN}✅ .env.production exists${NC}"
}

fix_missing_dirs() {
    echo -e "${YELLOW}Creating missing directories...${NC}"
    mkdir -p logs/{backend,frontend,postgres,redis}
    mkdir -p ssl
    mkdir -p backups
    echo -e "${GREEN}✅ Directories created${NC}"
}

restart_failed_services() {
    echo -e "${YELLOW}Checking for failed services...${NC}"
    
    failed=$(docker-compose -f docker-compose.production.yml ps --format json 2>/dev/null | jq -r 'select(.State != "running") | .Service' 2>/dev/null || echo "")
    
    if [ -n "$failed" ]; then
        echo -e "${RED}Failed services found: $failed${NC}"
        echo "Restarting..."
        docker-compose -f docker-compose.production.yml restart $failed
        echo -e "${GREEN}✅ Services restarted${NC}"
    else
        echo -e "${GREEN}✅ All services running${NC}"
    fi
}

# Run all fixes
echo -e "\n${BLUE}Running auto-fixes...${NC}"
fix_port_conflict
fix_docker_issues
fix_permissions
fix_env_file
fix_missing_dirs
restart_failed_services

echo -e "\n${GREEN}🎉 Auto-fix complete!${NC}"
echo "Run ./LIVE-TEST-DEPLOY.sh to deploy"
