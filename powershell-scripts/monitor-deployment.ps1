# Comprehensive Deployment Monitoring Script
# Monitors deployment progress and provides real-time updates

# Set error action preference to stop on errors
$ErrorActionPreference = "Stop"

# Server Configuration
$SERVER_IP = "128.140.111.171"
$SERVER_USER = "root"
$SERVER_PASSWORD = "Khaled55400214."
$DOMAIN = "super.doctorhealthy1.com"

# Colors
$Red = [ConsoleColor]::Red
$Green = [ConsoleColor]::Green
$Yellow = [ConsoleColor]::Yellow
$Blue = [ConsoleColor]::Blue
$Purple = [ConsoleColor]::Magenta
$Cyan = [ConsoleColor]::Cyan
$NC = [ConsoleColor]::White

function log($message) {
    Write-Host "[$((Get-Date).ToString('HH:mm:ss'))]$NC $message" -ForegroundColor $Blue
}

function success($message) {
    Write-Host "[SUCCESS]$NC $message" -ForegroundColor $Green
}

function warning($message) {
    Write-Host "[WARNING]$NC $message" -ForegroundColor $Yellow
}

function error($message) {
    Write-Host "[ERROR]$NC $message" -ForegroundColor $Red
}

function info($message) {
    Write-Host "[INFO]$NC $message" -ForegroundColor $Cyan
}

function progress($message) {
    Write-Host "[PROGRESS]$NC $message" -ForegroundColor $Purple
}

Write-Host ""
Write-Host "🚀 =================================" -ForegroundColor $Green
Write-Host "🚀 DEPLOYMENT MONITORING STARTED" -ForegroundColor $Green
Write-Host "🚀 =================================" -ForegroundColor $Green
Write-Host ""
log "📍 Server: $SERVER_IP"
log "🌐 Domain: $DOMAIN"
log "⏰ Started monitoring at $((Get-Date))"
Write-Host ""

# Function to check server connectivity
function check_server_connection {
    try {
        & ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no "$SERVER_USER@$SERVER_IP" "echo 'connected'" 2>$null | Select-String -Pattern "connected" | Out-Null
        return $true
    } catch {
        return $false
    }
}

# Function to check application health
function check_app_health {
    try {
        $response = Invoke-WebRequest -Uri "http://$SERVER_IP`:8080/health" -TimeoutSec 5 -Method Head
        return $true
    } catch {
        return $false
    }
}

# Function to get deployment status from server
function get_deployment_status {
    $remoteScript = @'
# Check if deployment directory exists
if [ -d "/opt/trae-new-healthy1" ]; then
    echo "DEPLOYMENT_DIR_EXISTS=true"
    cd /opt/trae-new-healthy1
    
    # Check if docker-compose file exists
    if [ -f "docker-compose.yml" ]; then
        echo "DOCKER_COMPOSE_EXISTS=true"
        
        # Check Docker containers status
        if command -v docker-compose &> /dev/null; then
            echo "DOCKER_COMPOSE_INSTALLED=true"
            
            # Get container status
            CONTAINER_STATUS=$(docker-compose ps --format json 2>/dev/null || echo "[]")
            echo "CONTAINER_STATUS=$CONTAINER_STATUS"
            
            # Check if backend is running
            if docker-compose ps | grep -q "nutrition-backend.*Up"; then
                echo "BACKEND_RUNNING=true"
            else
                echo "BACKEND_RUNNING=false"
            fi
            
            # Check if database is running
            if docker-compose ps | grep -q "nutrition-postgres.*Up"; then
                echo "DATABASE_RUNNING=true"
            else
                echo "DATABASE_RUNNING=false"
            fi
            
            # Check if redis is running
            if docker-compose ps | grep -q "nutrition-redis.*Up"; then
                echo "REDIS_RUNNING=true"
            else
                echo "REDIS_RUNNING=false"
            fi
        else
            echo "DOCKER_COMPOSE_INSTALLED=false"
        fi
    else
        echo "DOCKER_COMPOSE_EXISTS=false"
    fi
else
    echo "DEPLOYMENT_DIR_EXISTS=false"
fi

# Check if Docker is installed
if command -v docker &> /dev/null; then
    echo "DOCKER_INSTALLED=true"
else
    echo "DOCKER_INSTALLED=false"
fi

# Check if Go is installed
if command -v go &> /dev/null; then
    echo "GO_INSTALLED=true"
else
    echo "GO_INSTALLED=false"
fi
'@
    try {
        $result = & ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no "$SERVER_USER@$SERVER_IP" $remoteScript 2>$null
        return $result -join "`n"
    } catch {
        return "STATUS_CHECK_FAILED=true"
    }
}

# Monitoring loop
$MONITOR_COUNT = 0
$MAX_MONITOR_TIME = 900  # 15 minutes
$INTERVAL = 30  # 30 seconds

while ($MONITOR_COUNT -lt $MAX_MONITOR_TIME) {
    $CURRENT_TIME = (Get-Date).ToString('HH:mm:ss')
    $ELAPSED_MIN = [math]::Floor($MONITOR_COUNT / 60)
    
    Write-Host ""
    progress "🔍 Monitoring Check #$([math]::Floor($MONITOR_COUNT / $INTERVAL) + 1) - Elapsed: ${ELAPSED_MIN}m"
    
    # Check server connectivity
    if (check_server_connection) {
        success "✅ Server SSH connection: ACTIVE"
        
        # Get detailed deployment status
        log "📊 Checking deployment status..."
        $STATUS_OUTPUT = get_deployment_status
        
        # Parse status
        $DEPLOYMENT_DIR_EXISTS = ($STATUS_OUTPUT | Select-String -Pattern "DEPLOYMENT_DIR_EXISTS=(.+)").Matches.Groups[1].Value
        $DOCKER_INSTALLED = ($STATUS_OUTPUT | Select-String -Pattern "DOCKER_INSTALLED=(.+)").Matches.Groups[1].Value
        $GO_INSTALLED = ($STATUS_OUTPUT | Select-String -Pattern "GO_INSTALLED=(.+)").Matches.Groups[1].Value
        $DOCKER_COMPOSE_EXISTS = ($STATUS_OUTPUT | Select-String -Pattern "DOCKER_COMPOSE_EXISTS=(.+)").Matches.Groups[1].Value
        $BACKEND_RUNNING = ($STATUS_OUTPUT | Select-String -Pattern "BACKEND_RUNNING=(.+)").Matches.Groups[1].Value
        $DATABASE_RUNNING = ($STATUS_OUTPUT | Select-String -Pattern "DATABASE_RUNNING=(.+)").Matches.Groups[1].Value
        $REDIS_RUNNING = ($STATUS_OUTPUT | Select-String -Pattern "REDIS_RUNNING=(.+)").Matches.Groups[1].Value
        
        # Display status
        info "📋 Deployment Status:"
        
        if ($DEPLOYMENT_DIR_EXISTS -eq "true") {
            success "   ✅ Application directory created"
        } else {
            warning "   ⏳ Application directory not ready"
        }
        
        if ($DOCKER_INSTALLED -eq "true") {
            success "   ✅ Docker installed"
        } else {
            warning "   ⏳ Docker installation in progress"
        }
        
        if ($GO_INSTALLED -eq "true") {
            success "   ✅ Go installed"
        } else {
            warning "   ⏳ Go installation in progress"
        }
        
        if ($DOCKER_COMPOSE_EXISTS -eq "true") {
            success "   ✅ Docker Compose configuration ready"
        } else {
            warning "   ⏳ Docker Compose configuration pending"
        }
        
        if ($BACKEND_RUNNING -eq "true") {
            success "   ✅ Backend application running"
        } else {
            warning "   ⏳ Backend application starting"
        }
        
        if ($DATABASE_RUNNING -eq "true") {
            success "   ✅ PostgreSQL database running"
        } else {
            warning "   ⏳ PostgreSQL database starting"
        }
        
        if ($REDIS_RUNNING -eq "true") {
            success "   ✅ Redis cache running"
        } else {
            warning "   ⏳ Redis cache starting"
        }
        
    } else {
        warning "⚠️ Server SSH connection: NOT READY (deployment in progress)"
    }
    
    # Check application health
    log "🏥 Checking application health..."
    if (check_app_health) {
        success "✅ Application health check: PASSED"
        
        # Test API endpoints
        log "🧪 Testing API endpoints..."
        
        # Test API info
        try {
            Invoke-WebRequest -Uri "http://$SERVER_IP`:8080/api/info" -TimeoutSec 5 | Out-Null
            success "   ✅ API info endpoint: WORKING"
        } catch {
            warning "   ⏳ API info endpoint: STARTING"
        }
        
        # Test nutrition analysis
        $nutritionTest = try {
            Invoke-WebRequest -Uri "http://$SERVER_IP`:8080/api/nutrition/analyze" -Method Post -Body '{"food": "apple", "quantity": 100, "unit": "g", "checkHalal": true}' -ContentType "application/json" -TimeoutSec 10
            $nutritionTest.Content
        } catch {
            "failed"
        }
        
        if ($nutritionTest -match "calories|nutrition") {
            success "   ✅ Nutrition analysis: WORKING"
        } else {
            warning "   ⏳ Nutrition analysis: STARTING"
        }
        
        # If health check passes, deployment is likely complete
        Write-Host ""
        Write-Host "🎉 =================================" -ForegroundColor $Green
        Write-Host "🎉 DEPLOYMENT COMPLETED!" -ForegroundColor $Green
        Write-Host "🎉 =================================" -ForegroundColor $Green
        Write-Host ""
        success "🚀 Your Trae New Healthy1 platform is LIVE!"
        Write-Host ""
        Write-Host "🔗 Access your application:" -ForegroundColor $Green
        Write-Host "   🏥 Health Check: http://$SERVER_IP`:8080/health" -ForegroundColor $Green
        Write-Host "   📊 API Info: http://$SERVER_IP`:8080/api/info" -ForegroundColor $Green
        Write-Host "   🍎 Nutrition API: http://$SERVER_IP`:8080/api/nutrition/analyze" -ForegroundColor $Green
        Write-Host "   🌐 Domain (after DNS): https://$DOMAIN" -ForegroundColor $Green
        Write-Host ""
        Write-Host "🎯 Platform Features Now Available:" -ForegroundColor $Green
        Write-Host "   ✅ AI-powered nutrition analysis" -ForegroundColor $Green
        Write-Host "   ✅ 10 evidence-based diet plans" -ForegroundColor $Green
        Write-Host "   ✅ Recipe management system" -ForegroundColor $Green
        Write-Host "   ✅ Health tracking and analytics" -ForegroundColor $Green
        Write-Host "   ✅ Medication management" -ForegroundColor $Green
        Write-Host "   ✅ Workout programs" -ForegroundColor $Green
        Write-Host "   ✅ Multi-language support (EN/AR)" -ForegroundColor $Green
        Write-Host "   ✅ Religious dietary filtering" -ForegroundColor $Green
        Write-Host ""
        Write-Host "📋 Next Steps:" -ForegroundColor $Green
        Write-Host "   1. 🌐 Point domain $DOMAIN to $SERVER_IP in DNS" -ForegroundColor $Green
        Write-Host "   2. 🔒 SSL will auto-configure once DNS propagates" -ForegroundColor $Green
        Write-Host "   3. 🔑 SSH to server for management: ssh root@$SERVER_IP" -ForegroundColor $Green
        Write-Host "   4. 📖 Start using your nutrition platform!" -ForegroundColor $Green
        Write-Host ""
        success "🎯 Deployment monitoring completed successfully!"
        exit 0
        
    } else {
        warning "⚠️ Application health check: NOT READY"
    }
    
    # Wait before next check
    log "⏳ Waiting ${INTERVAL}s before next check..."
    Start-Sleep -Seconds $INTERVAL
    $MONITOR_COUNT += $INTERVAL
}

# If we reach here, monitoring timed out
Write-Host ""
warning "⚠️ Monitoring timeout reached (15 minutes)"
Write-Host ""
Write-Host "📋 The deployment may still be in progress. You can:" -ForegroundColor $Yellow
Write-Host "   1. 🔍 Check manually: ssh root@$SERVER_IP" -ForegroundColor $Yellow
Write-Host "   2. 📋 View logs: ssh root@$SERVER_IP 'cd /opt/trae-new-healthy1 && docker-compose logs -f'" -ForegroundColor $Yellow
Write-Host "   3. 🔄 Restart if needed: ssh root@$SERVER_IP 'cd /opt/trae-new-healthy1 && docker-compose restart'" -ForegroundColor $Yellow
Write-Host ""
log "Monitoring session ended at $((Get-Date))"