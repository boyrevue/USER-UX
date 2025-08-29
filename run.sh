#!/bin/bash

# Set environment variables for Tesseract and Leptonica
export PKG_CONFIG_PATH="/opt/homebrew/lib/pkgconfig"
export CGO_CPPFLAGS="-I/opt/homebrew/Cellar/tesseract/5.5.1/include -I/opt/homebrew/Cellar/leptonica/1.85.0/include/leptonica"
export CGO_LDFLAGS="-L/opt/homebrew/Cellar/tesseract/5.5.1/lib -L/opt/homebrew/Cellar/leptonica/1.85.0/lib -ltesseract -lleptonica"
export CGO_ENABLED=1
export TESSDATA_PREFIX="/opt/homebrew/share/tessdata"

echo "Starting Insurance Quote App..."
echo "Environment variables set:"
echo "CGO_CPPFLAGS: $CGO_CPPFLAGS"
echo "CGO_LDFLAGS: $CGO_LDFLAGS"
echo "TESSDATA_PREFIX: $TESSDATA_PREFIX"
echo ""

# Run the application
go run .