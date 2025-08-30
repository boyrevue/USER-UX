#!/bin/bash

# CLIENT-UX Architecture Cleanup Script
# This script implements the cleanup plan for better AI manageability

set -e

echo "🏗️ CLIENT-UX Architecture Cleanup Starting..."

# 1. Clean up session files (keep only last 10)
echo "📁 Cleaning up session files..."
cd sessions/
ls -t *.json | tail -n +11 | xargs -r rm -f
echo "✅ Cleaned session files (kept last 10)"
cd ..

# 2. Remove duplicate static files
echo "📁 Cleaning up static directory..."
find static/ -name "*.map" -delete 2>/dev/null || true
find static/ -name "*.LICENSE.txt" -delete 2>/dev/null || true
echo "✅ Removed duplicate static files"

# 3. Clean up node_modules if needed (optional)
if [ "$1" = "--deep" ]; then
    echo "🧹 Deep cleaning node_modules..."
    rm -rf insurance-frontend/node_modules/
    rm -rf settings-frontend/node_modules/ 2>/dev/null || true
    rm -rf settings_app/node_modules/ 2>/dev/null || true
    echo "✅ Removed node_modules (run npm install to restore)"
fi

# 4. Remove old backup files
echo "🗑️ Removing backup files..."
find . -name "*.backup" -delete 2>/dev/null || true
find . -name "*.bak" -delete 2>/dev/null || true
find . -name "*~" -delete 2>/dev/null || true
echo "✅ Removed backup files"

# 5. Create .gitignore for build artifacts if not exists
if [ ! -f .gitignore ]; then
    echo "📝 Creating .gitignore..."
    cat > .gitignore << 'EOF'
# Build artifacts
client-ux
*.exe
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out

# Dependency directories
node_modules/
jspm_packages/

# Build directories
build/
dist/

# Runtime files
*.log
server.log
app.log

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# Temporary files
tmp/
temp/
*.tmp

# Session files (keep only recent ones)
sessions/*.json
!sessions/README.md

# Static build files
static/static/
EOF
    echo "✅ Created .gitignore"
fi

# 6. Check file sizes and report issues
echo "📏 Checking file sizes..."
echo "Files larger than 1000 lines:"
find . -name "*.go" -o -name "*.tsx" -o -name "*.ttl" | while read file; do
    if [ -f "$file" ]; then
        lines=$(wc -l < "$file" 2>/dev/null || echo 0)
        if [ "$lines" -gt 1000 ]; then
            echo "  ⚠️  $file: $lines lines"
        fi
    fi
done

echo "📊 Cleanup Summary:"
echo "  Sessions: $(ls sessions/*.json 2>/dev/null | wc -l) files"
echo "  Static size: $(du -sh static/ 2>/dev/null | cut -f1)"
echo "  Total Go lines: $(find . -name "*.go" -exec wc -l {} \; 2>/dev/null | awk '{sum+=$1} END {print sum}')"
echo "  Total TTL lines: $(find . -name "*.ttl" -exec wc -l {} \; 2>/dev/null | awk '{sum+=$1} END {print sum}')"

echo "✅ CLIENT-UX Architecture Cleanup Complete!"
echo "🎯 Next steps: Run 'npm run build' in insurance-frontend/ to rebuild"
