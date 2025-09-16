#!/bin/bash

# VANTUN GitHub Actions History Cleanup Script
# This script cleans up GitHub Actions workflow history while preserving documentation

echo "🧹 Starting GitHub Actions history cleanup..."

# Store current branch
CURRENT_BRANCH=$(git branch --show-current)

# Create a backup branch
echo "📋 Creating backup branch..."
git branch backup-$(date +%Y%m%d-%H%M%S)

# Get the latest commit hash
LATEST_COMMIT=$(git rev-parse HEAD)

echo "🔄 Rewriting Git history to remove Actions workflow history..."

# Create a new orphan branch for clean history
git checkout --orphan clean-main

# Remove all files from the new branch
git rm -rf .

# Checkout only the files we want to keep from the latest commit
git checkout $LATEST_COMMIT -- .

# Remove GitHub Actions workflows (they'll be added back in the next commit)
rm -rf .github/workflows/
rm -rf .github/scripts/

# Create new commit with clean history
echo "📝 Creating new clean commit..."
git add -A
git commit -m "Clean repository with documentation only

- Removed GitHub Actions workflow history
- Preserved all documentation and source code
- Cleaned up repository for fresh start

This commit represents a clean state with only essential files."

# Force push to main (be careful with this!)
echo "⚠️  Force pushing to main branch..."
echo "This will overwrite the remote main branch history!"
read -p "Are you sure you want to continue? (yes/no): " CONFIRM

if [[ "$CONFIRM" == "yes" ]]; then
    echo "🚀 Force pushing clean history..."
    git push origin clean-main:main --force
    
    echo "✅ History cleanup completed!"
    echo "📍 Your repository now has a clean history without Actions workflows."
    
    # Switch back to main
    git checkout main
    git branch -D clean-main
    
else
    echo "❌ Cleanup cancelled. No changes were made."
    git checkout $CURRENT_BRANCH
    git branch -D clean-main
fi

echo "🎉 Cleanup process completed!"
echo "📚 Your documentation is available at: https://tungoldshou.github.io/vantun/"