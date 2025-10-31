#!/bin/bash

# GitHub Upload Script for Nutrition Platform
# Account: Khaledalzayat278@gmail.com

set -e

echo "🚀 GitHub Upload Script for Nutrition Platform"
echo "============================================="
echo ""
echo "📧 GitHub Account: Khaledalzayat278@gmail.com"
echo "📦 Repository: nutrition-platform"
echo ""

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo "❌ Error: Not in a git repository. Please run 'git init' first."
    exit 1
fi

# Check if we have commits
if ! git log --oneline -1 &>/dev/null; then
    echo "❌ Error: No commits found. Please commit your changes first."
    exit 1
fi

echo "✅ Git repository check passed"

# Check if remote origin exists
if git remote get-url origin &>/dev/null; then
    echo "⚠️  Remote 'origin' already exists:"
    git remote get-url origin
    echo ""
    read -p "Do you want to continue? (y/N): " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ Aborted by user"
        exit 1
    fi
else
    echo "📝 Setting up GitHub remote..."
    
    # Add the remote repository
    git remote add origin https://github.com/Khaledalzayat278/nutrition-platform.git
    
    if [ $? -eq 0 ]; then
        echo "✅ Remote repository added successfully"
    else
        echo "❌ Failed to add remote repository"
        echo "💡 Make sure you've created the repository on GitHub first:"
        echo "   https://github.com/new"
        exit 1
    fi
fi

echo ""
echo "🔄 Preparing to push to GitHub..."
echo "   Repository: https://github.com/Khaledalzayat278/nutrition-platform"
echo "   Branch: main"
echo ""

# Set main branch
echo "📋 Setting up main branch..."
git branch -M main

if [ $? -eq 0 ]; then
    echo "✅ Main branch configured"
else
    echo "❌ Failed to configure main branch"
    exit 1
fi

# Push to GitHub
echo "🚀 Pushing to GitHub..."
echo "   This may take a few minutes for the first push..."
echo ""

git push -u origin main

if [ $? -eq 0 ]; then
    echo ""
    echo "🎉 Successfully uploaded to GitHub!"
    echo "============================================="
    echo "✅ Your Nutrition Platform is now on GitHub!"
    echo ""
    echo "🔗 Repository URL:"
    echo "   https://github.com/Khaledalzayat278/nutrition-platform"
    echo ""
    echo "📱 Features uploaded:"
    echo "   • Complete nutrition planning system"
    echo "   • 50+ medical condition support"
    echo "   • Workout generator"
    echo "   • Diet planning tools"
    echo "   • System validation dashboard"
    echo "   • Production-ready deployment configs"
    echo ""
    echo "🌟 Next steps:"
    echo "   1. Visit your repository on GitHub"
    echo "   2. Add repository description and topics"
    echo "   3. Enable GitHub Pages (optional)"
    echo "   4. Share your project with the world!"
    echo ""
    echo "🚀 Deploy to production:"
    echo "   • Vercel: Run './deploy.sh'"
    echo "   • GitHub Pages: Enable in repository settings"
    echo ""
else
    echo "❌ Failed to push to GitHub"
    echo ""
    echo "💡 Troubleshooting:"
    echo "   1. Make sure you've created the repository on GitHub:"
    echo "      https://github.com/new"
    echo "   2. Repository name should be: nutrition-platform"
    echo "   3. Make sure you're logged into GitHub"
    echo "   4. Check your internet connection"
    echo ""
    echo "🔑 If authentication fails:"
    echo "   1. Go to GitHub.com → Settings → Developer settings"
    echo "   2. Generate a Personal Access Token"
    echo "   3. Use token as password when prompted"
    echo ""
    exit 1
fi

echo "✅ GitHub upload completed successfully!"