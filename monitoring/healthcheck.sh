#!/bin/bash

# ============================================
# AUTOMATED HEALTH CHECK & MONITORING
# Runs every 5 minutes via cron
# ============================================

DOMAIN="${1:-super.doctorhealthy1.com}"
ALERT_EMAIL="${2:-admin@doctorhealthy1.com}"
LOG_FILE="/var/log/nutrition-platform/healthcheck.log"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

alert() {
  log "🚨 ALERT: $1"
  # Send email alert
  echo "$1" | mail -s "Health Check Alert - $DOMAIN" "$ALERT_EMAIL" 2>/dev/null || true
}

# Test 1: Health Endpoint
log "Testing health endpoint..."
HEALTH_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "https://$DOMAIN/health")
if [ "$HEALTH_STATUS" != "200" ]; then
  alert "Health endpoint returned $HEALTH_STATUS"
  exit 1
fi
log "✅ Health endpoint OK"

# Test 2: Response Time
log "Testing response time..."
START=$(date +%s%N)
curl -s "https://$DOMAIN/health" > /dev/null
END=$(date +%s%N)
RESPONSE_TIME=$(( (END - START) / 1000000 ))

if [ $RESPONSE_TIME -gt 2000 ]; then
  alert "Slow response time: ${RESPONSE_TIME}ms"
fi
log "✅ Response time: ${RESPONSE_TIME}ms"

# Test 3: API Endpoints
log "Testing API endpoints..."
API_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "https://$DOMAIN/api/info")
if [ "$API_STATUS" != "200" ]; then
  alert "API endpoint returned $API_STATUS"
  exit 1
fi
log "✅ API endpoints OK"

# Test 4: SSL Certificate
log "Checking SSL certificate..."
EXPIRY=$(echo | openssl s_client -servername "$DOMAIN" -connect "$DOMAIN:443" 2>/dev/null | \
  openssl x509 -noout -enddate | cut -d= -f2)
EXPIRY_EPOCH=$(date -d "$EXPIRY" +%s)
NOW_EPOCH=$(date +%s)
DAYS_LEFT=$(( (EXPIRY_EPOCH - NOW_EPOCH) / 86400 ))

if [ $DAYS_LEFT -lt 30 ]; then
  alert "SSL certificate expires in $DAYS_LEFT days"
fi
log "✅ SSL certificate valid ($DAYS_LEFT days remaining)"

# Test 5: Memory Usage (if running locally)
if command -v docker &> /dev/null; then
  CONTAINER_ID=$(docker ps -q -f name=nutrition-platform)
  if [ -n "$CONTAINER_ID" ]; then
    MEMORY=$(docker stats --no-stream --format "{{.MemUsage}}" "$CONTAINER_ID")
    log "📊 Memory usage: $MEMORY"
  fi
fi

log "✅ All health checks passed"
exit 0
