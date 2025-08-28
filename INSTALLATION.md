# Installation Guide - Insurance Quote App with Enhanced Passport OCR

## Quick Start

```bash
# Clone and setup
git clone <repository-url>
cd insurance-quote-app

# Install dependencies
./install.sh

# Start application
./run.sh
```

## System Requirements

### Operating System
- **macOS**: 10.15+ (Catalina or later)
- **Ubuntu**: 18.04+ LTS
- **Debian**: 10+ (Buster or later)
- **Windows**: WSL2 recommended

### Hardware Requirements
- **RAM**: 4GB minimum, 8GB recommended
- **Storage**: 2GB free space
- **CPU**: Multi-core recommended for OCR processing

## Detailed Installation

### 1. Install Go

#### macOS (Homebrew)
```bash
brew install go
```

#### Ubuntu/Debian
```bash
# Download and install Go 1.21+
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### Verify Installation
```bash
go version
# Should output: go version go1.21.0 darwin/amd64 (or similar)
```

### 2. Install Python 3.8+

#### macOS
```bash
# Python should be pre-installed, but you can update via Homebrew
brew install python3
```

#### Ubuntu/Debian
```bash
sudo apt-get update
sudo apt-get install python3 python3-pip python3-dev
```

#### Verify Installation
```bash
python3 --version
pip3 --version
```

### 3. Install Tesseract OCR

#### macOS (Homebrew)
```bash
# Install Tesseract and Leptonica
brew install tesseract leptonica

# Verify installation paths
brew --prefix tesseract
brew --prefix leptonica
```

#### Ubuntu/Debian
```bash
# Install Tesseract and development libraries
sudo apt-get install tesseract-ocr tesseract-ocr-eng
sudo apt-get install libtesseract-dev libleptonica-dev

# Install additional language packs if needed
sudo apt-get install tesseract-ocr-fra tesseract-ocr-deu
```

#### Verify Tesseract Installation
```bash
tesseract --version
tesseract --list-langs
```

### 4. Install Python Dependencies

```bash
# Install required Python packages
pip3 install --upgrade pip
pip3 install passporteye pytesseract pillow numpy opencv-python

# Verify PassportEye installation
python3 -c "import passporteye; print('PassportEye installed successfully')"
```

### 5. Clone and Setup Project

```bash
# Clone the repository
git clone <your-repository-url>
cd insurance-quote-app

# Install Go dependencies
go mod tidy

# Make scripts executable
chmod +x run.sh
chmod +x *.py
```

### 6. Configure Environment

#### Automatic Configuration (Recommended)
The `run.sh` script automatically configures environment variables for macOS Homebrew installations.

#### Manual Configuration
If you have custom installation paths, edit `run.sh`:

```bash
# Find your Tesseract installation
which tesseract
brew --prefix tesseract  # macOS only

# Find tessdata directory
find /usr -name "tessdata" 2>/dev/null
find /opt -name "tessdata" 2>/dev/null

# Update run.sh with your paths
export TESSDATA_PREFIX="/your/path/to/tessdata"
export CGO_CPPFLAGS="-I/your/tesseract/include -I/your/leptonica/include"
export CGO_LDFLAGS="-L/your/tesseract/lib -ltesseract -L/your/leptonica/lib -lleptonica"
```

### 7. Build and Test

```bash
# Build the application
go build -o insurance-quote-app

# Test OCR functionality
python3 passporteye_simplified_extractor.py --help

# Start the application
./run.sh
```

## Frontend Setup (Optional)

If you want to rebuild the frontend:

```bash
cd insurance-frontend

# Install Node.js dependencies
npm install

# Build frontend
npm run build

# Copy build files to static directory
cp -r build/* ../static/
```

## Verification

### 1. Test Go Application
```bash
# Check if server starts
curl http://localhost:3000/api/health
```

### 2. Test OCR Functionality
```bash
# Test debug endpoint
curl http://localhost:3000/api/debug/ocr-test

# Test with sample image (if available)
curl -X POST -F "file=@sample_passport.png" \
     -F "documentType=passport" \
     http://localhost:3000/api/process-document
```

### 3. Test Python Scripts
```bash
# Test basic functionality
python3 -c "
import passporteye
import pytesseract
from PIL import Image
print('All Python dependencies working')
"
```

## Troubleshooting

### Common Installation Issues

#### 1. Go Build Errors
```
# Error: package not found
go mod tidy
go clean -modcache
go mod download
```

#### 2. Tesseract Not Found
```bash
# macOS: Reinstall via Homebrew
brew uninstall tesseract leptonica
brew install tesseract leptonica

# Ubuntu: Reinstall packages
sudo apt-get remove tesseract-ocr libtesseract-dev
sudo apt-get install tesseract-ocr libtesseract-dev
```

#### 3. Python Import Errors
```bash
# Upgrade pip and reinstall
pip3 install --upgrade pip
pip3 uninstall passporteye pytesseract pillow
pip3 install passporteye pytesseract pillow
```

#### 4. CGO Compilation Errors
```bash
# Check compiler installation
# macOS
xcode-select --install

# Ubuntu
sudo apt-get install build-essential
```

#### 5. Permission Errors
```bash
# Make scripts executable
chmod +x run.sh
chmod +x *.py

# Fix ownership if needed
sudo chown -R $USER:$USER .
```

### Environment-Specific Issues

#### macOS Apple Silicon (M1/M2)
```bash
# Use Homebrew paths for Apple Silicon
export CGO_CPPFLAGS="-I/opt/homebrew/include"
export CGO_LDFLAGS="-L/opt/homebrew/lib"
```

#### Ubuntu 20.04+
```bash
# Install additional dependencies
sudo apt-get install pkg-config
```

#### Windows (WSL2)
```bash
# Install Windows Subsystem for Linux
# Follow Ubuntu installation steps within WSL2
```

## Performance Optimization

### 1. Image Processing
```bash
# Install additional image processing libraries
pip3 install opencv-python-headless
```

### 2. Concurrent Processing
```bash
# Set Go runtime parameters
export GOMAXPROCS=4  # Adjust based on CPU cores
```

### 3. Memory Management
```bash
# Set memory limits if needed
export GOMEMLIMIT=2GiB
```

## Security Considerations

### 1. File Permissions
```bash
# Secure sensitive files
chmod 600 config.json
chmod 700 sessions/
```

### 2. Network Security
```bash
# Run on localhost only for development
# Use reverse proxy (nginx) for production
```

## Production Deployment

### 1. Environment Variables
```bash
# Set production environment
export GO_ENV=production
export PORT=3000
```

### 2. Process Management
```bash
# Use systemd, supervisor, or PM2 for process management
# Example systemd service file available in docs/
```

### 3. Monitoring
```bash
# Set up logging
export LOG_LEVEL=info
export LOG_FILE=/var/log/insurance-app.log
```

## Support

### Getting Help
1. **Check logs**: `tail -f app.log`
2. **Test components individually**: Use debug endpoints
3. **Verify dependencies**: Run verification commands
4. **Check environment**: Ensure all paths are correct

### Reporting Issues
Include the following information:
- Operating system and version
- Go version (`go version`)
- Python version (`python3 --version`)
- Tesseract version (`tesseract --version`)
- Error messages and logs
- Steps to reproduce

---

**Installation Guide Version**: 1.0.0  
**Last Updated**: August 29, 2025  
**Compatible with**: Go 1.19+, Python 3.8+, Tesseract 4.0+
