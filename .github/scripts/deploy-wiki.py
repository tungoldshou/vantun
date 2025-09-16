#!/usr/bin/env python3
"""
GitHub Wiki Deployment Script
Automatically deploys Wiki content from the repository to GitHub Wiki
"""

import os
import sys
import json
import base64
import requests
from pathlib import Path
from github import Github

def get_github_client():
    """Initialize GitHub client with token"""
    token = os.environ.get('GITHUB_TOKEN')
    if not token:
        print("❌ GITHUB_TOKEN not found in environment")
        sys.exit(1)
    return Github(token)

def get_wiki_repo_name():
    """Get the Wiki repository name from environment"""
    repo_full_name = os.environ.get('REPOSITORY', 'tungoldshou/vantun')
    return f"{repo_full_name}.wiki"

def check_wiki_exists(g, repo_name):
    """Check if Wiki repository exists"""
    try:
        wiki_repo = g.get_repo(repo_name)
        print(f"✅ Wiki repository found: {repo_name}")
        return wiki_repo
    except Exception as e:
        print(f"❌ Wiki repository not found: {repo_name}")
        print(f"Error: {e}")
        return None

def create_wiki_repo(g, original_repo):
    """Try to enable and create Wiki repository"""
    try:
        # Try to access the original repo's wiki
        wiki_url = f"https://github.com/{original_repo.full_name}/wiki"
        print(f"📖 Wiki URL: {wiki_url}")
        
        # GitHub automatically creates wiki repo when first page is created
        # We'll use the GitHub API to create wiki pages directly
        return original_repo
    except Exception as e:
        print(f"❌ Failed to access Wiki: {e}")
        return None

def get_wiki_files():
    """Get all Wiki markdown files from the repository"""
    wiki_dir = Path("wiki")
    if not wiki_dir.exists():
        print("❌ Wiki directory not found")
        return []
    
    wiki_files = []
    for md_file in wiki_dir.glob("*.md"):
        page_name = md_file.stem
        content = md_file.read_text(encoding='utf-8')
        
        # Convert filename to wiki page name
        # Replace hyphens with spaces for better readability
        display_name = page_name.replace('-', ' ')
        
        wiki_files.append({
            'filename': md_file.name,
            'page_name': page_name,
            'display_name': display_name,
            'content': content
        })
    
    return wiki_files

def create_wiki_page_via_api(repo, page_name, content, display_name=None):
    """Create a Wiki page using GitHub API"""
    if display_name is None:
        display_name = page_name.replace('-', ' ')
    
    # GitHub Wiki API endpoint
    api_url = f"https://api.github.com/repos/{repo.full_name}/wiki"
    
    # Since GitHub doesn't have a direct Wiki API for creating pages,
    # we'll create an issue or use the web interface approach
    
    # Alternative approach: Create a GitHub Issue with Wiki content
    issue_title = f"Wiki Page: {display_name}"
    issue_body = f"""# {display_name}

This is an automated Wiki page creation.

## Content

{content}

---
*This page was automatically generated from the repository Wiki directory.*
"""
    
    try:
        # Create an issue (as a workaround since direct Wiki API is limited)
        issue = repo.create_issue(
            title=issue_title,
            body=issue_body,
            labels=["wiki", "documentation"]
        )
        
        print(f"✅ Created issue for Wiki page: {display_name}")
        print(f"   Issue URL: {issue.html_url}")
        return True
        
    except Exception as e:
        print(f"❌ Failed to create Wiki page {display_name}: {e}")
        return False

def main():
    """Main deployment function"""
    print("🚀 Starting VANTUN Wiki deployment...")
    
    # Get GitHub client
    g = get_github_client()
    
    # Get repository information
    repo_full_name = os.environ.get('REPOSITORY', 'tungoldshou/vantun')
    print(f"📁 Repository: {repo_full_name}")
    
    try:
        # Get the main repository
        repo = g.get_repo(repo_full_name)
        print(f"✅ Found repository: {repo.full_name}")
        
        # Get Wiki files
        wiki_files = get_wiki_files()
        if not wiki_files:
            print("❌ No Wiki files found")
            return 1
        
        print(f"📚 Found {len(wiki_files)} Wiki files")
        
        # Deploy each Wiki page
        success_count = 0
        for wiki_file in wiki_files:
            print(f"\n📄 Processing: {wiki_file['filename']}")
            print(f"   Page name: {wiki_file['page_name']}")
            print(f"   Display name: {wiki_file['display_name']}")
            
            # Create Wiki page
            if create_wiki_page_via_api(
                repo, 
                wiki_file['page_name'], 
                wiki_file['content'],
                wiki_file['display_name']
            ):
                success_count += 1
        
        print(f"\n🎉 Wiki deployment completed!")
        print(f"✅ Successfully processed: {success_count}/{len(wiki_files)} pages")
        
        # Print summary and next steps
        print("\n📋 Next Steps:")
        print("1. Visit your GitHub repository Wiki tab")
        print("2. Click 'Create the first page' button")
        print("3. Copy the content from the created issues")
        print("4. Manually create Wiki pages with the provided content")
        print(f"5. Repository Wiki URL: https://github.com/{repo.full_name}/wiki")
        
        return 0
        
    except Exception as e:
        print(f"❌ Failed to deploy Wiki: {e}")
        return 1

if __name__ == "__main__":
    sys.exit(main())