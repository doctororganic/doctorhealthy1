#!/bin/bash
# Script to push using a personal access token

echo "🔐 GitHub Authentication Required"
echo ""
echo "To push to doctororganic/doctorhealthy1, you need a Personal Access Token"
echo ""
echo "Steps:"
echo "1. Go to: https://github.com/settings/tokens"
echo "2. Click 'Generate new token' → 'Generate new token (classic)'"
echo "3. Name: 'doctorhealthy1-deployment'"
echo "4. Select scope: 'repo' (full control)"
echo "5. Click 'Generate token'"
echo "6. Copy the token"
echo ""
read -sp "Paste your Personal Access Token here: " TOKEN
echo ""

if [ -z "$TOKEN" ]; then
    echo "❌ No token provided. Exiting."
    exit 1
fi

# Configure git to use token
git remote set-url origin https://${TOKEN}@github.com/doctororganic/doctorhealthy1.git

echo "📤 Pushing to GitHub..."
git push -u origin main

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Successfully pushed to GitHub!"
    echo "🔗 Repository: https://github.com/doctororganic/doctorhealthy1"
    echo "🔗 Actions: https://github.com/doctororganic/doctorhealthy1/actions"
    
    # Reset remote URL to remove token
    git remote set-url origin https://github.com/doctororganic/doctorhealthy1.git
else
    echo ""
    echo "❌ Push failed. Please check:"
    echo "   - Token has 'repo' scope"
    echo "   - Token is for doctororganic account"
    echo "   - Repository exists and is accessible"
fi

