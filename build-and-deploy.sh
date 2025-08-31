#!/bin/bash

echo "🔧 CLIENT-UX Build and Deploy Script"
echo "======================================"

# Build the React frontend
echo "📦 Building React frontend..."
cd insurance-frontend
npm run build
if [ $? -ne 0 ]; then
    echo "❌ Frontend build failed!"
    exit 1
fi
cd ..

# Clean the static directory completely
echo "🧹 Cleaning static directory..."
rm -rf static/*

# Copy files with correct structure to prevent static/static/ nesting
echo "📁 Copying build files with correct structure..."

# Copy root files (index.html, manifest.json, etc.)
cp insurance-frontend/build/*.* static/ 2>/dev/null || true

# Copy CSS files directly to static/css (not static/static/css)
if [ -d "insurance-frontend/build/static/css" ]; then
    echo "🎨 Copying CSS files..."
    mkdir -p static/css
    cp -r insurance-frontend/build/static/css/* static/css/
fi

# Copy JS files directly to static/js (not static/static/js)
if [ -d "insurance-frontend/build/static/js" ]; then
    echo "📜 Copying JS files..."
    mkdir -p static/js
    cp -r insurance-frontend/build/static/js/* static/js/
fi

# Copy any other static assets
if [ -d "insurance-frontend/build/static/media" ]; then
    echo "🖼️  Copying media files..."
    mkdir -p static/media
    cp -r insurance-frontend/build/static/media/* static/media/
fi

echo "✅ Static files organized correctly:"
echo "   📁 static/css/ - CSS files"
echo "   📁 static/js/  - JavaScript files"
echo "   📄 static/index.html - Main HTML file"

# Verify no nested static directory exists
if [ -d "static/static" ]; then
    echo "❌ ERROR: Nested static/static/ directory detected!"
    echo "🔧 Removing nested directory..."
    rm -rf static/static
fi

# Build the Go backend
echo "🏗️  Building Go backend..."
go build -o client-ux .
if [ $? -ne 0 ]; then
    echo "❌ Backend build failed!"
    exit 1
fi

# Set execute permissions on the binary (PERMANENT FIX)
echo "🔐 Setting execute permissions..."
chmod +x client-ux

echo "🎉 Build completed successfully!"
echo "🚀 Ready to run: ./client-ux"
echo "🌐 Will serve on: http://localhost:3000"