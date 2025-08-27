#!/bin/bash

# 🚀 Insurance Quote App - Complete Startup Script
# This script ensures the application starts fresh every time without any forgotten discoveries

set -e  # Exit on any error

echo "=========================================="
echo "🏗️  Insurance Quote App - Fresh Start"
echo "=========================================="

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

# Check if we're in the right directory
if [ ! -f "main.go" ]; then
    print_error "This script must be run from the insurance-quote-app directory"
    exit 1
fi

print_status "Starting fresh application setup..."

# Step 1: Kill any existing processes
print_status "Killing any existing Go processes..."
pkill -f "insurance-quote-app" || true
pkill -f "go run" || true
sleep 2

# Step 2: Clean and rebuild React frontend
print_status "Building React frontend..."
cd insurance-frontend
npm install --silent
npm run build
cd ..

# Step 3: Copy React build to static directory
print_status "Copying React build to static directory..."
rm -rf static/*
cp -r insurance-frontend/build/* static/

# Step 4: Fix nested static directory issue (React creates static/static/)
print_status "Fixing nested static directory structure..."
if [ -d "static/static" ]; then
    cp -r static/static/js/* static/js/ 2>/dev/null || true
    cp -r static/static/css/* static/css/ 2>/dev/null || true
    rm -rf static/static
    print_success "Fixed nested static directory"
fi

# Step 5: Verify static files exist
print_status "Verifying static files..."
if [ ! -f "static/index.html" ]; then
    print_error "static/index.html not found!"
    exit 1
fi

if [ ! -f "static/manifest.json" ]; then
    print_error "static/manifest.json not found!"
    exit 1
fi

print_success "Static files verified"

# Step 6: Build Go backend
print_status "Building Go backend..."
go mod tidy
go build -o insurance-quote-app .

# Step 7: Verify Go binary
if [ ! -f "insurance-quote-app" ]; then
    print_error "Go binary not built!"
    exit 1
fi

print_success "Go backend built successfully"

# Step 8: Start the server
print_status "Starting server..."
print_success "Application starting on http://localhost:3000"
print_status "Press Ctrl+C to stop"

# Start the server in foreground
./insurance-quote-app
