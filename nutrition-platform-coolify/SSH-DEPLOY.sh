#!/bin/bash

################################################################################
# SSH DEPLOYMENT SCRIPT
# Automated deployment to remote server via SSH
################################################################################

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

log() { echo -e "${GREEN}[$(date +'%H:%M:%S')]${NC} $1"; }
log_error() { echo -e "${RED}[$(date +'%H:%M:%S')] ERROR:${NC} $1"; }
log_success() { echo -e "${GREEN}[$(date +'%H:%M:%S')] ✓${NC} $1"; }
log_info() { echo -e "${BLUE}[$(date +'%H:%M:%S')] INFO:${NC} $1"; }

# Configuration
SSH_HOST="${SSH_HOST:-}"
SSH_USER="${SSH_USER:-root}"
SSH_PORT="${SSH_PORT:-22}"
DEPLOY_PATH="${DEPLOY_PATH:-/opt/nutrition-platform}"

if [ -z "$SSH_HOST" ]; then
    log_error "SSH_HOST environment variable is required"
    echo ""
    echo "Usage:"
    echo "  SSH_HOST=your-server.com SSH_USER=user ./SSH-DEPLOY.sh"
    echo ""
    echo "Optional variables:"
    echo "  SSH_PORT=22 (default)"
    echo "  DEPLOY_PATH=/opt/nutrition-platform (default)"
    exit 1
fi

log "═══════════════════════════════════════════════════════════════"
log "SSH DEPLOYMENT TO $SSH_HOST"
log "═══════════════════════════════════════════════════════════════"

# Build deployment package
log "Building deployment package..."
./AUTO-FACTORY-ORCHESTRATOR.sh

# Find latest deployment package
DEPLOY_PACKAGE=$(ls -t deploy_*.tar.gz | head -1)

if [ -z "$DEPLOY_PACKAGE" ]; then
    log_error "No deployment package found"
    exit 1
fi

log_info "Deployment package: $DEPLOY_PACKAGE"
log_info "Package size: $(du -h "$DEPLOY_PACKAGE" | cut -f1)"

# Test SSH connection
log "Testing SSH connection..."
if ! ssh -p "$SSH_PORT" -o ConnectTimeout=10 "$SSH_USER@$SSH_HOST" "echo 'SSH connection successful'" > /dev/null 2>&1; then
    log_error "Failed to connect to $SSH_HOST"
    exit 1
fi
log_success "SSH connection successful"

# Upload to server
log "Uploading to server..."
scp -P "$SSH_PORT" "$DEPLOY_PACKAGE" "$SSH_USER@$SSH_HOST:/tmp/" || {
    log_error "Failed to upload deployment package"
    exit 1
}
log_success "Upload completed"

# Deploy on server
log "Deploying on server..."
ssh -p "$SSH_PORT" "$SSH_USER@$SSH_HOST" << ENDSSH
    set -e
    
    echo "Creating deployment directory..."
    mkdir -p $DEPLOY_PATH
    cd $DEPLOY_PATH
    
    # Backup current deployment
    if [ -d "current" ]; then
        echo "Backing up current deployment..."
        mv current "backup_\$(date +%Y%m%d_%H%M%S)"
        
        # Keep only last 3 backups
        ls -dt backup_* | tail -n +4 | xargs rm -rf 2>/dev/null || true
    fi
    
    # Extract new deployment
    echo "Extracting new deployment..."
    mkdir -p current
    tar -xzf /tmp/$DEPLOY_PACKAGE -C current/
    
    # Stop current service
    echo "Stopping current service..."
    systemctl stop nutrition-platform 2>/dev/null || true
    
    # Start new service
    cd current
    chmod +x bin/server
    
    # Create systemd service if it doesn't exist
    if [ ! -f /etc/systemd/system/nutrition-platform.service ]; then
        echo "Creating systemd service..."
        cat > /etc/systemd/system/nutrition-platform.service << 'EOF'
[Unit]
Description=Nutrition Platform API
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$DEPLOY_PATH/current
ExecStart=$DEPLOY_PATH/current/bin/server
Restart=always
RestartSec=10
Environment="PORT=8080"

[Install]
WantedBy=multi-user.target
EOF
    fi
    
    # Reload and start
    echo "Starting service..."
    systemctl daemon-reload
    systemctl enable nutrition-platform
    systemctl start nutrition-platform
    
    # Cleanup
    rm /tmp/$DEPLOY_PACKAGE
    
    echo "Deployment completed successfully"
ENDSSH

log_success "Deployment completed!"

# Verify deployment
log "Verifying deployment..."
sleep 10

if curl -f "http://$SSH_HOST:8080/health" > /dev/null 2>&1; then
    log_success "Application is running!"
else
    log_error "Application health check failed"
    log_info "Checking service status..."
    ssh -p "$SSH_PORT" "$SSH_USER@$SSH_HOST" "systemctl status nutrition-platform"
    exit 1
fi

log "═══════════════════════════════════════════════════════════════"
log "DEPLOYMENT SUCCESSFUL"
log "═══════════════════════════════════════════════════════════════"
log_info "Application URL: http://$SSH_HOST:8080"
log_info "Health Check: http://$SSH_HOST:8080/health"
log_info "API Docs: http://$SSH_HOST:8080/api/v1"
