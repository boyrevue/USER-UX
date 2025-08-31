#!/bin/bash

# CLIENT-UX Server Startup Script
# Fixes all the recurring issues permanently

echo "🚀 CLIENT-UX Server Startup Script"
echo "=================================="

# 1. Ensure we're in the right directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"
echo "📁 Working directory: $(pwd)"

# 2. Kill any existing processes on port 3000
echo "🔄 Cleaning up existing processes..."
pkill -f "go run.*client-ux" 2>/dev/null || true
pkill -f "client-ux" 2>/dev/null || true
lsof -ti:3000 | xargs kill -9 2>/dev/null || true
sleep 2

# 3. Check if we have the main.go file
if [ ! -f "main.go" ]; then
    echo "❌ Error: main.go not found in $(pwd)"
    echo "   Make sure you're running this from the client-ux directory"
    exit 1
fi

# 4. Check if go.mod exists
if [ ! -f "go.mod" ]; then
    echo "❌ Error: go.mod not found. Please run 'go mod init' first"
    exit 1
fi

# 5. Build the binary first (faster startup)
echo "🔨 Building client-ux binary..."
go build -o client-ux . 2>&1
if [ $? -ne 0 ]; then
    echo "❌ Build failed. Check your Go code for errors."
    exit 1
fi

# 6. Set executable permissions
chmod +x client-ux

# 7. Start the server
echo "🌐 Starting CLIENT-UX server..."
echo "   URL: http://localhost:3000"
echo "   Press Ctrl+C to stop"
echo ""

# Run in foreground so you can see logs and stop with Ctrl+C
./client-ux

# Cleanup on exit
trap 'echo "🛑 Shutting down..."; pkill -f client-ux; exit 0' INT TERM