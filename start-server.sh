#!/bin/bash

echo "🚀 CLIENT-UX Server Startup Script"
echo "=================================="

# Ensure the binary exists
if [ ! -f "client-ux" ]; then
    echo "❌ client-ux binary not found. Running build first..."
    ./build-and-deploy.sh
fi

# Always ensure execute permissions (PERMANENT FIX)
echo "🔐 Ensuring execute permissions..."
chmod +x client-ux

# Kill any existing processes
echo "🧹 Cleaning up any existing processes..."
pkill -f client-ux 2>/dev/null || true

# Start the server
echo "🌐 Starting server on http://localhost:3000..."
./client-ux
