#!/bin/bash

# Fix Static Directory Structure
# This script properly copies React build files and fixes the double static/static/ nesting

echo "🔧 Fixing static directory structure..."

# Remove old static content (except mrz directory which contains uploaded images)
find static -type f ! -path "static/mrz/*" -delete 2>/dev/null || true
find static -type d -empty -delete 2>/dev/null || true

# Ensure mrz directory exists
mkdir -p static/mrz

# Copy React build files correctly
echo "📁 Copying React build files..."
cp -r insurance-frontend/build/* static/

# Fix the double static/static/ nesting if it exists
if [ -d "static/static" ]; then
    echo "⚠️  Found nested static/static/ - fixing..."
    
    # Move JS files
    if [ -d "static/static/js" ]; then
        mkdir -p static/js
        mv static/static/js/* static/js/
        rmdir static/static/js
    fi
    
    # Move CSS files  
    if [ -d "static/static/css" ]; then
        mkdir -p static/css
        mv static/static/css/* static/css/
        rmdir static/static/css
    fi
    
    # Remove empty nested static directory
    rmdir static/static 2>/dev/null || true
    
    echo "✅ Fixed nested static directory structure"
else
    echo "✅ No nested static directory found"
fi

# Verify critical files exist
echo "🔍 Verifying critical files..."
critical_files=(
    "static/index.html"
    "static/manifest.json"
    "static/js/main.*.js"
    "static/css/main.*.css"
)

for pattern in "${critical_files[@]}"; do
    if ls $pattern 1> /dev/null 2>&1; then
        echo "✅ Found: $pattern"
    else
        echo "❌ Missing: $pattern"
    fi
done

echo "🎉 Static directory structure fixed!"
