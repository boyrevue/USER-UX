#!/bin/bash

# Insurance Quote App - Enhanced Passport OCR Installation Script
# This script automates the installation of all dependencies

set -e  # Exit on any error

echo "🚀 Installing Insurance Quote App with Enhanced Passport OCR..."
echo "================================================================"

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

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Detect operating system
detect_os() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if command_exists apt-get; then
            OS="ubuntu"
        elif command_exists yum; then
            OS="centos"
        else
            OS="linux"
        fi
    else
        OS="unknown"
    fi
    print_status "Detected OS: $OS"
}

# Install dependencies based on OS
install_system_deps() {
    print_status "Installing system dependencies..."
    
    case $OS in
        "macos")
            if ! command_exists brew; then
                print_error "Homebrew not found. Please install Homebrew first:"
                print_error "https://brew.sh"
                exit 1
            fi
            
            print_status "Installing Tesseract and Leptonica via Homebrew..."
            brew install tesseract leptonica
            
            print_status "Installing Python3 if needed..."
            brew install python3
            ;;
            
        "ubuntu")
            print_status "Updating package list..."
            sudo apt-get update
            
            print_status "Installing Tesseract OCR and development libraries..."
            sudo apt-get install -y tesseract-ocr tesseract-ocr-eng
            sudo apt-get install -y libtesseract-dev libleptonica-dev
            
            print_status "Installing Python3 and development tools..."
            sudo apt-get install -y python3 python3-pip python3-dev
            sudo apt-get install -y build-essential pkg-config
            ;;
            
        *)
            print_error "Unsupported operating system: $OS"
            print_error "Please install dependencies manually:"
            print_error "- Go 1.19+"
            print_error "- Python 3.8+"
            print_error "- Tesseract OCR"
            print_error "- Leptonica"
            exit 1
            ;;
    esac
}

# Check Go installation
check_go() {
    print_status "Checking Go installation..."
    
    if ! command_exists go; then
        print_error "Go not found. Please install Go 1.19+ first:"
        print_error "https://golang.org/dl/"
        exit 1
    fi
    
    GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | sed 's/go//')
    MAJOR_VERSION=$(echo $GO_VERSION | cut -d. -f1)
    MINOR_VERSION=$(echo $GO_VERSION | cut -d. -f2)
    
    if [ "$MAJOR_VERSION" -lt 1 ] || ([ "$MAJOR_VERSION" -eq 1 ] && [ "$MINOR_VERSION" -lt 19 ]); then
        print_error "Go version $GO_VERSION found, but Go 1.19+ is required"
        exit 1
    fi
    
    print_success "Go $GO_VERSION found"
}

# Check Python installation
check_python() {
    print_status "Checking Python installation..."
    
    if ! command_exists python3; then
        print_error "Python3 not found"
        exit 1
    fi
    
    PYTHON_VERSION=$(python3 --version | grep -o '[0-9]\+\.[0-9]\+')
    MAJOR_VERSION=$(echo $PYTHON_VERSION | cut -d. -f1)
    MINOR_VERSION=$(echo $PYTHON_VERSION | cut -d. -f2)
    
    if [ "$MAJOR_VERSION" -lt 3 ] || ([ "$MAJOR_VERSION" -eq 3 ] && [ "$MINOR_VERSION" -lt 8 ]); then
        print_error "Python $PYTHON_VERSION found, but Python 3.8+ is required"
        exit 1
    fi
    
    print_success "Python $PYTHON_VERSION found"
}

# Install Python dependencies
install_python_deps() {
    print_status "Installing Python dependencies..."
    
    # Upgrade pip first
    python3 -m pip install --upgrade pip
    
    # Install required packages
    python3 -m pip install passporteye pytesseract pillow numpy opencv-python
    
    # Verify installations
    python3 -c "import passporteye; print('✓ PassportEye installed')"
    python3 -c "import pytesseract; print('✓ PyTesseract installed')"
    python3 -c "import PIL; print('✓ Pillow installed')"
    
    print_success "Python dependencies installed"
}

# Install Go dependencies
install_go_deps() {
    print_status "Installing Go dependencies..."
    
    if [ ! -f "go.mod" ]; then
        print_error "go.mod not found. Are you in the correct directory?"
        exit 1
    fi
    
    go mod tidy
    go mod download
    
    print_success "Go dependencies installed"
}

# Verify Tesseract installation
verify_tesseract() {
    print_status "Verifying Tesseract installation..."
    
    if ! command_exists tesseract; then
        print_error "Tesseract not found in PATH"
        exit 1
    fi
    
    TESSERACT_VERSION=$(tesseract --version 2>&1 | head -n1)
    print_success "Tesseract found: $TESSERACT_VERSION"
    
    # Check for English language data
    if ! tesseract --list-langs 2>/dev/null | grep -q "eng"; then
        print_warning "English language data not found for Tesseract"
        print_warning "You may need to install tesseract-ocr-eng package"
    else
        print_success "English language data found"
    fi
}

# Configure environment
configure_environment() {
    print_status "Configuring environment..."
    
    # Make scripts executable
    chmod +x run.sh 2>/dev/null || true
    chmod +x *.py 2>/dev/null || true
    
    # Create necessary directories
    mkdir -p static/mrz
    mkdir -p sessions
    
    print_success "Environment configured"
}

# Build application
build_app() {
    print_status "Building application..."
    
    # Set CGO environment for building
    export CGO_ENABLED=1
    
    if [ "$OS" = "macos" ]; then
        # macOS Homebrew paths
        export CGO_CPPFLAGS="-I$(brew --prefix tesseract)/include -I$(brew --prefix leptonica)/include/leptonica"
        export CGO_LDFLAGS="-L$(brew --prefix tesseract)/lib -ltesseract -L$(brew --prefix leptonica)/lib -lleptonica"
    fi
    
    go build -o insurance-quote-app
    
    if [ -f "insurance-quote-app" ]; then
        print_success "Application built successfully"
    else
        print_error "Failed to build application"
        exit 1
    fi
}

# Test installation
test_installation() {
    print_status "Testing installation..."
    
    # Test Python scripts
    if [ -f "passporteye_simplified_extractor.py" ]; then
        python3 passporteye_simplified_extractor.py --help >/dev/null 2>&1 || true
        print_success "Python scripts are executable"
    fi
    
    # Test Go application
    if [ -f "insurance-quote-app" ]; then
        print_success "Go application built successfully"
    fi
    
    print_success "Installation test completed"
}

# Main installation process
main() {
    echo
    print_status "Starting installation process..."
    echo
    
    detect_os
    check_go
    check_python
    install_system_deps
    verify_tesseract
    install_python_deps
    install_go_deps
    configure_environment
    build_app
    test_installation
    
    echo
    echo "================================================================"
    print_success "🎉 Installation completed successfully!"
    echo
    print_status "Next steps:"
    echo "  1. Start the application: ./run.sh"
    echo "  2. Open browser: http://localhost:3000"
    echo "  3. Upload a passport image to test OCR"
    echo
    print_status "Documentation:"
    echo "  - README: PASSPORT_OCR_README.md"
    echo "  - Installation: INSTALLATION.md"
    echo
    print_status "Troubleshooting:"
    echo "  - Check logs: tail -f app.log"
    echo "  - Debug endpoint: curl http://localhost:3000/api/debug/ocr-test"
    echo "================================================================"
}

# Run main function
main "$@"
