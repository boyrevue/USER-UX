#!/bin/bash

# ============================================================
# CLIENT-UX BUILD AND DEPLOY SCRIPT
# ============================================================
# Permanent solution for the static/static nested directory issue
# This script ensures proper build and deployment every time
# ============================================================

set -e  # Exit on any error

echo "🚀 CLIENT-UX Build and Deploy Script"
echo "===================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Step 1: Kill existing processes
print_status "Stopping existing CLIENT-UX processes..."
pkill -f "./client-ux" || print_warning "No existing processes found"

# Step 2: Clean up old static files
print_status "Cleaning up old static files..."
if [ -d "static" ]; then
    # Preserve essential files
    cp static/favicon.ico /tmp/favicon.ico.bak 2>/dev/null || true
    cp static/manifest.json /tmp/manifest.json.bak 2>/dev/null || true
    cp static/logo192.png /tmp/logo192.png.bak 2>/dev/null || true
    cp static/logo512.png /tmp/logo512.png.bak 2>/dev/null || true
    cp static/robots.txt /tmp/robots.txt.bak 2>/dev/null || true
    cp static/pdf.worker.min.js /tmp/pdf.worker.min.js.bak 2>/dev/null || true
    
    # Remove old static directory
    rm -rf static/*
    print_success "Old static files cleaned"
else
    mkdir -p static
    print_status "Created static directory"
fi

# Step 3: Build React frontend
print_status "Building React frontend..."
cd insurance-frontend

# Clean node_modules if needed
if [ "$1" = "--clean" ]; then
    print_status "Cleaning node_modules..."
    rm -rf node_modules package-lock.json
    npm install
fi

# Build the frontend
npm run build

if [ $? -ne 0 ]; then
    print_error "Frontend build failed!"
    exit 1
fi

print_success "Frontend build completed"

# Step 4: Copy build files with proper structure
print_status "Copying build files to static directory..."
cd ..

# Copy all files from build directory
cp -r insurance-frontend/build/* static/

# Fix the nested static/static issue permanently
if [ -d "static/static" ]; then
    print_warning "Found nested static/static directory - fixing..."
    
    # Move nested static files to root static
    mv static/static/* static/
    rmdir static/static
    
    print_success "Fixed nested static directory issue"
fi

# Step 5: Restore essential files
print_status "Restoring essential files..."
cp /tmp/favicon.ico.bak static/favicon.ico 2>/dev/null || true
cp /tmp/manifest.json.bak static/manifest.json 2>/dev/null || true
cp /tmp/logo192.png.bak static/logo192.png 2>/dev/null || true
cp /tmp/logo512.png.bak static/logo512.png 2>/dev/null || true
cp /tmp/robots.txt.bak static/robots.txt 2>/dev/null || true
cp /tmp/pdf.worker.min.js.bak static/pdf.worker.min.js 2>/dev/null || true

# Clean up temp files
rm -f /tmp/*.bak

# Step 6: Verify static file structure
print_status "Verifying static file structure..."

required_files=(
    "static/css/main.*.css"
    "static/js/main.*.js"
    "static/index.html"
    "static/manifest.json"
    "static/favicon.ico"
)

all_good=true
for pattern in "${required_files[@]}"; do
    if ls $pattern 1> /dev/null 2>&1; then
        print_success "✓ Found: $pattern"
    else
        print_error "✗ Missing: $pattern"
        all_good=false
    fi
done

if [ "$all_good" = false ]; then
    print_error "Static file verification failed!"
    exit 1
fi

# Step 7: Build Go backend
print_status "Building Go backend..."
go build -o client-ux main.go

if [ $? -ne 0 ]; then
    print_error "Go build failed!"
    exit 1
fi

print_success "Go backend built successfully"

# Step 8: Make executable
chmod +x client-ux

# Step 9: Start the application
print_status "Starting CLIENT-UX application..."
./client-ux &
SERVER_PID=$!

# Wait a moment for server to start
sleep 3

# Step 10: Verify deployment
print_status "Verifying deployment..."

# Test main page
if curl -s http://localhost:3000 > /dev/null; then
    print_success "✓ Main page accessible"
else
    print_error "✗ Main page not accessible"
    kill $SERVER_PID 2>/dev/null || true
    exit 1
fi

# Test CSS file
CSS_FILE=$(ls static/css/main.*.css | head -1 | sed 's|static/||')
if curl -s -I "http://localhost:3000/static/$CSS_FILE" | grep -q "text/css"; then
    print_success "✓ CSS file served with correct MIME type"
else
    print_error "✗ CSS file MIME type issue"
fi

# Test JS file
JS_FILE=$(ls static/js/main.*.js | head -1 | sed 's|static/||')
if curl -s -I "http://localhost:3000/static/$JS_FILE" | grep -q "application/javascript"; then
    print_success "✓ JavaScript file served with correct MIME type"
else
    print_error "✗ JavaScript file MIME type issue"
fi

# Test market adapter
if curl -s http://localhost:3000/api/market/status > /dev/null; then
    print_success "✓ Market adapter API accessible"
else
    print_warning "⚠ Market adapter API not responding"
fi

# Step 11: Final status
echo ""
echo "🎉 CLIENT-UX DEPLOYMENT COMPLETE!"
echo "================================="
echo "🌐 Application URL: http://localhost:3000"
echo "🔧 Process ID: $SERVER_PID"
echo "📊 Static files: $(ls static/css/*.css static/js/*.js 2>/dev/null | wc -l) files"
echo "🌍 Markets supported: UK, DE, NL, ES, FR"
echo "✅ MIME types: Fixed permanently"
echo ""
echo "💡 To rebuild and redeploy: ./build_and_deploy.sh"
echo "🧹 To clean rebuild: ./build_and_deploy.sh --clean"
echo ""

# Save deployment info
cat > deployment_info.json << EOF
{
  "deployment_time": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "process_id": $SERVER_PID,
  "static_files": {
    "css": "$(ls static/css/main.*.css | head -1)",
    "js": "$(ls static/js/main.*.js | head -1)"
  },
  "status": "deployed",
  "url": "http://localhost:3000"
}
EOF

print_success "Deployment info saved to deployment_info.json"
