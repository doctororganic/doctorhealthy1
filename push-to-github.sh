#!/bin/bash
# Script to push code to GitHub repository

echo "🚀 Preparing to push to GitHub..."
echo "Repository: doctororganic/doctorhealthy1"
echo ""

# Check if we're in the right directory
if [ ! -f ".github/workflows/ci.yml" ]; then
    echo "❌ Error: Please run this script from the nutrition-platform directory"
    exit 1
fi

# Set remote URL
git remote set-url origin https://github.com/doctororganic/doctorhealthy1.git

echo "📋 Current status:"
git status --short | head -10
echo ""

# Check if there are uncommitted changes
if [ -n "$(git status --porcelain)" ]; then
    echo "⚠️  Warning: You have uncommitted changes"
    read -p "Do you want to commit them? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        git add .
        git commit --no-verify -m "Update project files"
    fi
fi

echo ""
echo "🔐 Authentication required"
echo "You'll need to authenticate with GitHub."
echo ""
echo "Option 1: Use Personal Access Token"
echo "  - Get token from: https://github.com/settings/tokens"
echo "  - Username: doctororganic"
echo "  - Password: <your personal access token>"
echo ""
echo "Option 2: Use GitHub CLI"
echo "  - Run: gh auth login"
echo "  - Then run this script again"
echo ""

read -p "Ready to push? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "📤 Pushing to GitHub..."
    git push -u origin main
    
    if [ $? -eq 0 ]; then
        echo ""
        echo "✅ Successfully pushed to GitHub!"
        echo "🔗 View repository: https://github.com/doctororganic/doctorhealthy1"
        echo "🔗 View Actions: https://github.com/doctororganic/doctorhealthy1/actions"
    else
        echo ""
        echo "❌ Push failed. Please check authentication."
        echo "See DEPLOYMENT_INSTRUCTIONS.md for help"
    fi
else
    echo "Push cancelled."
fi
