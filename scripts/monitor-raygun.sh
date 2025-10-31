#!/bin/bash

# Raygun Monitoring Script
# This script helps monitor and test Raygun error reporting

set -e

# Configuration
API_KEY="J5KNVQg46P71JymsDyPWiQ"
APP_NAME="nutrition-platform"
RAYGUN_API="https://api.raygun.io/v3/entries"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Test Raygun connectivity
test_connectivity() {
    log_info "Testing Raygun API connectivity..."
    
    response=$(curl -s -o /dev/null -w "%{http_code}" \
        -X POST \
        -H "X-ApiKey: $API_KEY" \
        -H "Content-Type: application/json" \
        -d '{
            "occurredOn": "'$(date -u +%Y-%m-%dT%H:%M:%S.%3NZ)'",
            "details": {
                "error": {
                    "message": "Connectivity test from monitoring script",
                    "stackTrace": []
                },
                "machineName": "monitoring-script",
                "version": "1.0.0"
            }
        }' \
        "$RAYGUN_API" 2>/dev/null)
    
    if [ "$response" = "202" ]; then
        log_success "Raygun API connectivity test passed"
        return 0
    else
        log_error "Raygun API connectivity test failed (HTTP $response)"
        return 1
    fi
}

# Check recent errors
check_recent_errors() {
    log_info "Checking for recent errors in Raygun..."
    
    echo "To check recent errors, visit:"
    echo "🔗 https://app.raygun.com"
    echo ""
    echo "Manual steps to check recent errors:"
    echo "1. Login to your Raygun dashboard"
    echo "2. Select the '$APP_NAME' application"
    echo "3. View the 'Errors' tab"
    echo "4. Filter by time range (last hour, last 24 hours, etc.)"
    echo ""
}

# Generate test error
generate_test_error() {
    log_info "Generating test error to verify monitoring..."
    
    cd backend
    go run -c "
package main

import (
    \"fmt\"
    \"log\"
    \"time\"

    \"github.com/MindscapeHQ/raygun4go\"
)

func main() {
    raygun, err := raygun4go.New(\"nutrition-platform-monitor\", \"$API_KEY\")
    if err != nil {
        log.Fatalf(\"Failed to create Raygun client: %v\", err)
    }
    defer raygun.HandleError()

    raygun.Tags([]string{\"monitoring\", \"test\", \"$(date +%s)\"})
    raygun.CustomData(map[string]interface{}{
        \"script\": \"monitor-raygun.sh\",
        \"timestamp\": time.Now().Unix(),
        \"test_type\": \"manual_monitoring_test\",
    })

    err = raygun.CreateError(\"Manual test error from monitoring script\")
    if err != nil {
        log.Printf(\"Failed to send error: %v\", err)
    } else {
        fmt.Println(\"Test error sent successfully to Raygun\")
    }
}
" 2>/dev/null || {
        log_warning "Could not generate test error via Go. Try running the application manually."
    }
    
    cd ..
}

# Check application health
check_app_health() {
    log_info "Checking application health..."
    
    if pgrep -f "nutrition-platform" > /dev/null; then
        log_success "Application process is running"
    else
        log_warning "Application process not found"
    fi
    
    # Check if application is responding on health endpoint
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        log_success "Application health endpoint responding"
    else
        log_warning "Application health endpoint not responding"
    fi
}

# Show dashboard instructions
show_dashboard_info() {
    log_info "Raygun Dashboard Information:"
    echo ""
    echo "📊 Dashboard: https://app.raygun.com"
    echo "🔑 API Key: $API_KEY"
    echo "📱 App Name: $APP_NAME"
    echo ""
    echo "Quick Links:"
    echo "• Error Dashboard: https://app.raygun.com/errors"
    echo "• Performance: https://app.raygun.com/performance"
    echo "• Users: https://app.raygun.com/users"
    echo "• Settings: https://app.raygun.com/settings"
    echo ""
}

# Main menu
show_menu() {
    echo ""
    echo "🔍 Raygun Monitoring Menu"
    echo "========================="
    echo "1. Test Raygun connectivity"
    echo "2. Generate test error"
    echo "3. Check application health"
    echo "4. Show dashboard information"
    echo "5. Check recent errors"
    echo "6. Exit"
    echo ""
    read -p "Select an option (1-6): " choice
}

# Main execution
main() {
    echo "🚀 Raygun Monitoring Script"
    echo "=========================="
    echo ""
    
    while true; do
        show_menu
        
        case $choice in
            1)
                test_connectivity
                ;;
            2)
                generate_test_error
                ;;
            3)
                check_app_health
                ;;
            4)
                show_dashboard_info
                ;;
            5)
                check_recent_errors
                ;;
            6)
                log_info "Exiting monitoring script"
                exit 0
                ;;
            *)
                log_error "Invalid option. Please select 1-6."
                ;;
        esac
        
        echo ""
        read -p "Press Enter to continue..."
    done
}

# Check if running with arguments
if [ $# -eq 1 ]; then
    case $1 in
        "test")
            test_connectivity
            ;;
        "error")
            generate_test_error
            ;;
        "health")
            check_app_health
            ;;
        "dashboard")
            show_dashboard_info
            ;;
        *)
            echo "Usage: $0 [test|error|health|dashboard]"
            exit 1
            ;;
    esac
else
    main
fi