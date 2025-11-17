#!/bin/bash
set -e

echo "🧹 WORKSPACE CLEANUP & OPTIMIZATION"
echo "===================================="

# Create archive directory
mkdir -p .archive
mkdir -p .archive/old-docs
mkdir -p .archive/old-scripts
mkdir -p .archive/old-deployments

# Track what we're doing
echo "" > cleanup-report.txt
echo "CLEANUP REPORT - $(date)" >> cleanup-report.txt
echo "================================" >> cleanup-report.txt

# 1. Remove duplicate/old deployment packages
echo ""
echo "1️⃣  Removing old deployment packages..."
find . -maxdepth 1 -type f \( -name "*.tar" -o -name "*.tar.gz" -o -name "*.zip" \) -exec du -h {} \; >> cleanup-report.txt
find . -maxdepth 1 -type f \( -name "*.tar" -o -name "*.tar.gz" -o -name "*.zip" \) -exec mv {} .archive/old-deployments/ \;
echo "✅ Archived deployment packages"

# 2. Consolidate duplicate documentation
echo ""
echo "2️⃣  Consolidating duplicate documentation..."

# Keep only essential docs, archive duplicates
DUPLICATE_DOCS=(
    "DEPLOY-NOW-GUIDE.md"
    "DEPLOY-NOW-INSTRUCTIONS.md"
    "DEPLOY-SECURE.md"
    "DEPLOYMENT-GUIDE.md"
    "DEPLOYMENT-INSTRUCTIONS.md"
    "DEPLOYMENT-SUCCESS-SUMMARY.md"
    "DEPLOYMENT.md"
    "MANUAL_DEPLOYMENT.md"
    "MANUAL-DEPLOYMENT-STEPS.md"
    "MANUAL-COOLIFY-DEPLOYMENT-GUIDE.md"
    "coolify-deployment-guide.md"
    "coolify-manual-final-steps.md"
    "coolify-step-by-step.md"
    "deployment-summary.md"
    "manual-fix-instructions.md"
    "ssh-commands.md"
    "ERRORS-FIXED-SUMMARY.md"
    "ERROR-MANAGEMENT-REVIEW.md"
    "CRITICAL-REVIEW-VERIFICATION.md"
    "CRITICAL-BUGS-FIXED.md"
    "TEST-FIXES-COMPLETE.md"
    "KILLER-BUGS-ANALYSIS.md"
)

for doc in "${DUPLICATE_DOCS[@]}"; do
    if [ -f "$doc" ]; then
        mv "$doc" .archive/old-docs/
        echo "Archived: $doc" >> cleanup-report.txt
    fi
done

echo "✅ Archived duplicate docs"

# 3. Consolidate duplicate scripts
echo ""
echo "3️⃣  Consolidating duplicate scripts..."

DUPLICATE_SCRIPTS=(
    "deploy-simple.sh"
    "deploy.sh"
    "DEPLOY-SECURE.sh"
    "DEPLOY-TO-COOLIFY.sh"
    "deploy-with-coolify.sh"
    "deploy-to-server.sh"
    "DEPLOY_NOW_FIX.sh"
    "DIAGNOSE-ALL-ISSUES.sh"
    "FIX-ALL-ERRORS.sh"
    "FIX-ERRORS-AUTO.sh"
    "AUTO-FIX-ERRORS.sh"
    "MONITOR-LIVE.sh"
    "MONITOR-LIVE-ERRORS.sh"
    "STRESS-TEST.sh"
    "STRESS-TEST-LIVE.sh"
    "LIVE-TEST-AND-DEPLOY.sh"
    "LIVE-TEST-DEPLOY.sh"
    "LIVE-DEPLOY-WITH-MONITORING.sh"
    "LIVE-TEST-AND-FIX.sh"
    "CONSOLIDATE-NOW.sh"
    "START-NOW.sh"
    "VERIFY-ENTERPRISE-STANDARDS.sh"
    "TEST-EVERYTHING.sh"
    "EXECUTE-IMPLEMENTATION.sh"
)

for script in "${DUPLICATE_SCRIPTS[@]}"; do
    if [ -f "$script" ]; then
        mv "$script" .archive/old-scripts/
        echo "Archived: $script" >> cleanup-report.txt
    fi
done

echo "✅ Archived duplicate scripts"

# 4. Remove old deployment directories
echo ""
echo "4️⃣  Cleaning old deployment directories..."

if [ -d "coolify-deployment-20251013-144403" ]; then
    tar -czf .archive/coolify-deployment-20251013-144403.tar.gz coolify-deployment-20251013-144403
    rm -rf coolify-deployment-20251013-144403
    echo "Archived: coolify-deployment-20251013-144403" >> cleanup-report.txt
fi

if [ -d "nutrition-platform-coolify" ]; then
    tar -czf .archive/nutrition-platform-coolify.tar.gz nutrition-platform-coolify
    rm -rf nutrition-platform-coolify
    echo "Archived: nutrition-platform-coolify" >> cleanup-report.txt
fi

echo "✅ Archived old deployment directories"

# 5. Clean build artifacts
echo ""
echo "5️⃣  Cleaning build artifacts..."

# Backend
if [ -f "backend/main" ]; then rm backend/main; fi
if [ -f "backend/nutrition-platform" ]; then rm backend/nutrition-platform; fi
if [ -f "backend/nutrition_platform.db" ]; then mv backend/nutrition_platform.db .archive/; fi

# Node modules (will be reinstalled)
if [ -d "node_modules" ]; then
    echo "Removing node_modules (will reinstall from package.json)"
    rm -rf node_modules
fi

echo "✅ Cleaned build artifacts"

# 6. Compress nutrition data (maintain quality)
echo ""
echo "6️⃣  Optimizing nutrition data..."

# Compress JSON files with gzip (lossless)
if [ -d "disease nutrition easy json files" ]; then
    cd "disease nutrition easy json files"
    for file in *.json; do
        if [ -f "$file" ]; then
            gzip -k -9 "$file"  # -k keeps original, -9 best compression
        fi
    done
    cd ..
    echo "✅ Compressed nutrition JSON files (originals kept)"
fi

# 7. Clean logs
echo ""
echo "7️⃣  Cleaning old logs..."

if [ -d "logs" ]; then
    find logs -type f -name "*.log" -mtime +7 -delete
    echo "✅ Removed logs older than 7 days"
fi

# 8. Clean cache
echo ""
echo "8️⃣  Cleaning cache..."

rm -rf cache/* 2>/dev/null || true
rm -rf .next 2>/dev/null || true
rm -rf frontend-nextjs/.next 2>/dev/null || true
rm -rf frontend/.vercel 2>/dev/null || true

echo "✅ Cleaned cache directories"

# 9. Remove duplicate Dockerfiles
echo ""
echo "9️⃣  Consolidating Dockerfiles..."

# Keep only .secure versions
if [ -f "Dockerfile.simple" ]; then mv Dockerfile.simple .archive/; fi
if [ -f "Dockerfile.web" ]; then mv Dockerfile.web .archive/; fi
if [ -f "backend/Dockerfile" ] && [ -f "backend/Dockerfile.secure" ]; then
    mv backend/Dockerfile .archive/
fi
if [ -f "frontend-nextjs/Dockerfile" ] && [ -f "frontend-nextjs/Dockerfile.secure" ]; then
    mv frontend-nextjs/Dockerfile .archive/
fi

echo "✅ Consolidated Dockerfiles"

# 10. Summary
echo ""
echo "================================"
echo "✅ CLEANUP COMPLETE"
echo "================================"
echo ""

# Calculate space saved
ARCHIVE_SIZE=$(du -sh .archive | awk '{print $1}')
CURRENT_SIZE=$(du -sh . | awk '{print $1}')

echo "📊 Results:"
echo "  Archive size: $ARCHIVE_SIZE"
echo "  Current size: $CURRENT_SIZE"
echo ""
echo "📁 Archived files in: .archive/"
echo "📋 Full report: cleanup-report.txt"
echo ""
echo "✅ Essential files preserved:"
echo "  - All source code (backend, frontend)"
echo "  - Nutrition data (compressed)"
echo "  - Active deployment scripts"
echo "  - Current documentation"
echo "  - Configuration files"
echo ""
