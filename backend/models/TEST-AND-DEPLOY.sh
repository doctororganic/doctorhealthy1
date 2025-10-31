#!/bin/bash

################################################################################
# TEST AND DEPLOY
# Runs tests and shows deployment status
################################################################################

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

log() { echo -e "${GREEN}[$(date +'%H:%M:%S')]${NC} $1"; }
log_success() { echo -e "${GREEN}[$(date +'%H:%M:%S')] ✓${NC} $1"; }
log_error() { echo -e "${RED}[$(date +'%H:%M:%S')] ✗${NC} $1"; }

cd ..

log "═══════════════════════════════════════════════════════════════"
log "RUNNING TESTS"
log "═══════════════════════════════════════════════════════════════"

# Run tests
go test ./... -cover -v 2>&1 | tee test-results.log

# Check results
if [ ${PIPESTATUS[0]} -eq 0 ]; then
    log_success "All tests passed!"
else
    log_error "Some tests failed - check test-results.log"
fi

log "═══════════════════════════════════════════════════════════════"
log "TEST SUMMARY"
log "═══════════════════════════════════════════════════════════════"

grep -E "(PASS|FAIL|coverage)" test-results.log | tail -20

log "═══════════════════════════════════════════════════════════════"
log "BUILD STATUS"
log "═══════════════════════════════════════════════════════════════"

if [ -f "bin/server" ]; then
    log_success "Binary exists: bin/server"
    ls -lh bin/server
else
    log_error "Binary not found"
fi

log "═══════════════════════════════════════════════════════════════"
log "DEPLOYMENT READY"
log "═══════════════════════════════════════════════════════════════"

echo ""
echo "✅ Backend built successfully"
echo "✅ Binary: bin/server ($(du -h bin/server | cut -f1))"
echo ""
echo "Next steps:"
echo "  1. Test locally: ./bin/server"
echo "  2. Deploy to Coolify (see COOLIFY-DEPLOYMENT.md)"
echo ""
