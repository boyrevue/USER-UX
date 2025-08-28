#!/bin/bash

echo "🚀 Building and deploying Money Supermarket Infiltration System..."

# Kill any running server
pkill -f insurance-quote-app 2>/dev/null || true

# Build frontend with correct settings
cd insurance-frontend
echo "📦 Building React frontend..."
npm run build

# Clean and copy static files properly
cd ..
echo "🧹 Cleaning static directory..."
rm -rf static/*

echo "📁 Copying build files with correct structure..."
# Copy root files (index.html, manifest.json, etc.)
cp insurance-frontend/build/*.* static/ 2>/dev/null || true

# Copy static assets with proper structure (this is the key fix)
if [ -d "insurance-frontend/build/static" ]; then
    echo "🔧 Copying static assets to correct locations..."
    # Copy CSS files
    if [ -d "insurance-frontend/build/static/css" ]; then
        mkdir -p static/css
        cp -r insurance-frontend/build/static/css/* static/css/
    fi
    # Copy JS files
    if [ -d "insurance-frontend/build/static/js" ]; then
        mkdir -p static/js
        cp -r insurance-frontend/build/static/js/* static/js/
    fi
    # Copy any other static assets
    if [ -d "insurance-frontend/build/static/media" ]; then
        mkdir -p static/media
        cp -r insurance-frontend/build/static/media/* static/media/
    fi
else
    echo "⚠️  No static directory found in build, copying all files..."
    cp -r insurance-frontend/build/* static/
fi

echo "✅ Static files organized correctly"

# Build Go backend
echo "🔨 Building Go backend..."
go build -o insurance-quote-app .

# Start server
echo "🚀 Starting Money Supermarket Infiltration System..."
./insurance-quote-app &

sleep 3

# Test the deployment
echo "🧪 Testing deployment..."
if curl -s http://localhost:3000/static/css/main.*.css | head -1 | grep -q "<!doctype"; then
    echo "❌ CSS file not found or returning HTML"
else
    echo "✅ CSS file serving correctly"
fi

if curl -s http://localhost:3000/static/js/main.*.js | head -1 | grep -q "<!doctype"; then
    echo "❌ JS file not found or returning HTML"
else
    echo "✅ JS file serving correctly"
fi

echo "🎯 Money Supermarket Infiltration System deployed!"
echo "🌐 Access at: http://localhost:3000"
echo "🥷 Navigate to: Navigate & Fill → Auto Fill → Money Supermarket Infiltration"
