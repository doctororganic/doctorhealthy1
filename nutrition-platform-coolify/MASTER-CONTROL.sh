#!/bin/bash

################################################################################
# MASTER CONTROL
# Central control panel for all operations
################################################################################

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m'

clear

echo -e "${CYAN}"
cat << 'EOF'
╔═══════════════════════════════════════════════════════════════╗
║                                                               ║
║         NUTRITION PLATFORM - MASTER CONTROL                   ║
║                                                               ║
║         Automated Orchestration & Deployment System           ║
║                                                               ║
╚═══════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

show_menu() {
    echo ""
    echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}  MAIN MENU${NC}"
    echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    echo ""
    echo -e "  ${CYAN}1.${NC} Complete Setup (All-in-One)"
    echo -e "  ${CYAN}2.${NC} Run Tests (Parallel)"
    echo -e "  ${CYAN}3.${NC} Build & Package"
    echo -e "  ${CYAN}4.${NC} Deploy to Production"
    echo -e "  ${CYAN}5.${NC} Monitor Application"
    echo -e "  ${CYAN}6.${NC} Run Load Tests"
    echo -e "  ${CYAN}7.${NC} Security Scan"
    echo -e "  ${CYAN}8.${NC} Auto-Fix Issues"
    echo -e "  ${CYAN}9.${NC} Generate Docker Compose"
    echo -e "  ${CYAN}10.${NC} Build Frontend"
    echo ""
    echo -e "  ${RED}0.${NC} Exit"
    echo ""
    echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    echo -n "Select option: "
}

run_command() {
    local cmd=$1
    local name=$2
    
    echo ""
    echo -e "${GREEN}Running: $name${NC}"
    echo -e "${BLUE}─────────────────────────────────────────────────────────────────${NC}"
    
    if [ -f "$cmd" ]; then
        chmod +x "$cmd"
        ./"$cmd"
    else
        echo -e "${RED}Error: $cmd not found${NC}"
    fi
    
    echo ""
    echo -e "${BLUE}─────────────────────────────────────────────────────────────────${NC}"
    echo -n "Press Enter to continue..."
    read
}

while true; do
    clear
    echo -e "${CYAN}"
    cat << 'EOF'
╔═══════════════════════════════════════════════════════════════╗
║         NUTRITION PLATFORM - MASTER CONTROL                   ║
╚═══════════════════════════════════════════════════════════════╝
EOF
    echo -e "${NC}"
    
    show_menu
    read choice
    
    case $choice in
        1)
            run_command "COMPLETE-SETUP.sh" "Complete Setup"
            ;;
        2)
            run_command "PARALLEL-TEST-RUNNER.sh" "Parallel Tests"
            ;;
        3)
            run_command "AUTO-FACTORY-ORCHESTRATOR.sh" "Build & Package"
            ;;
        4)
            echo ""
            echo -n "Enter SSH_HOST: "
            read ssh_host
            echo -n "Enter SSH_USER (default: root): "
            read ssh_user
            ssh_user=${ssh_user:-root}
            
            SSH_HOST=$ssh_host SSH_USER=$ssh_user run_command "SSH-DEPLOY.sh" "Deploy to Production"
            ;;
        5)
            echo ""
            echo -n "Enter API_URL (default: http://localhost:8080): "
            read api_url
            api_url=${api_url:-http://localhost:8080}
            
            API_URL=$api_url run_command "REAL-TIME-MONITOR.sh" "Monitor Application"
            ;;
        6)
            run_command "LOAD-TEST.sh" "Load Tests"
            ;;
        7)
            run_command "SECURITY-SCAN.sh" "Security Scan"
            ;;
        8)
            run_command "AUTO-FIX-AGENT.sh" "Auto-Fix Issues"
            ;;
        9)
            run_command "DOCKER-COMPOSE-GENERATOR.sh" "Generate Docker Compose"
            ;;
        10)
            run_command "FRONTEND-BUILDER.sh" "Build Frontend"
            ;;
        0)
            echo ""
            echo -e "${GREEN}Goodbye!${NC}"
            exit 0
            ;;
        *)
            echo ""
            echo -e "${RED}Invalid option${NC}"
            sleep 2
            ;;
    esac
done
