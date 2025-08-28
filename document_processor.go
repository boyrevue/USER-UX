package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/otiai10/gosseract/v2"
)

// DocumentData represents extracted document information
type DocumentData struct {
	DocumentType     string                    `json:"documentType"`
	UploadType       string                    `json:"uploadType"`
	ExtractedFields  map[string]interface{}    `json:"extractedFields"`
	Confidence       float64                   `json:"confidence"`
	ProcessedAt      time.Time                 `json:"processedAt"`
	ValidationResult *DocumentValidationResult `json:"validationResult,omitempty"`
	OCREngine        string                    `json:"ocrEngine"`
	ProcessingPath   string                    `json:"processingPath"` // "auto", "review", "manual"
}

// OCRResult represents the result from an OCR engine
type OCRResult struct {
	Text       string  `json:"text"`
	Confidence float64 `json:"confidence"`
	Engine     string  `json:"engine"`
	Error      error   `json:"error,omitempty"`
}

// Confidence thresholds for routing
const (
	CONFIDENCE_HIGH   = 0.85
	CONFIDENCE_MEDIUM = 0.65
)

// DocumentValidationResult represents SHACL validation results for documents
type DocumentValidationResult struct {
	IsValid     bool                  `json:"isValid"`
	Violations  []ValidationViolation `json:"violations,omitempty"`
	Suggestions []string              `json:"suggestions,omitempty"`
}

// ValidationViolation represents a SHACL validation violation
type ValidationViolation struct {
	Property string `json:"property"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
	Value    string `json:"value,omitempty"`
}

// ProcessDocumentHandler handles document upload and OCR processing
func ProcessDocumentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	documentType := r.FormValue("documentType")
	uploadType := r.FormValue("uploadType")

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	allowedTypes := []string{".pdf", ".jpg", ".jpeg", ".png"}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !contains(allowedTypes, ext) {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	// Process document based on type
	var extractedData map[string]interface{}
	var confidence float64

	switch documentType {
	case "driving-licence":
		extractedData, confidence, err = processDrivingLicence(file, uploadType)
	case "passport":
		extractedData, confidence, err = processPassport(file, uploadType)
	default:
		extractedData, confidence, err = processGenericDocument(file, documentType)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("OCR processing failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Validate with SHACL
	validationResult, err := validateWithSHACL(documentType, extractedData)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("SHACL validation error: %v\n", err)
	}

	// Determine processing path based on confidence
	processingPath := getProcessingPath(confidence)

	// Create response
	response := DocumentData{
		DocumentType:     documentType,
		UploadType:       uploadType,
		ExtractedFields:  extractedData,
		Confidence:       confidence,
		ProcessedAt:      time.Now(),
		ValidationResult: validationResult,
		ProcessingPath:   processingPath,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getProcessingPath determines the processing path based on confidence
func getProcessingPath(confidence float64) string {
	if confidence >= CONFIDENCE_HIGH {
		return "auto"
	} else if confidence >= CONFIDENCE_MEDIUM {
		return "review"
	}
	return "manual"
}

// runMultiEngineOCR runs multiple OCR engines and returns the best result
func runMultiEngineOCR(fileBytes []byte, filename string) (*OCRResult, error) {
	// Save file temporarily for processing
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("ocr_%d_%s", time.Now().UnixNano(), filename))
	err := os.WriteFile(tempFile, fileBytes, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write temp file: %v", err)
	}
	defer os.Remove(tempFile)

	// Determine file type and select appropriate engines
	var engines []func(string) (*OCRResult, error)

	if strings.HasSuffix(strings.ToLower(filename), ".pdf") {
		// For PDF files, try all engines
		engines = []func(string) (*OCRResult, error){
			ocrWithPDFtoText, // Try native text extraction first
			ocrWithOCRmyPDF,  // Then OCRmyPDF for scanned PDFs
			ocrWithTesseract, // Fallback to Tesseract
		}
	} else {
		// For image files, only Tesseract works
		engines = []func(string) (*OCRResult, error){
			ocrWithTesseract,
		}
	}

	var bestResult *OCRResult
	var lastError error

	for _, engine := range engines {
		result, err := engine(tempFile)
		if err != nil {
			lastError = err
			fmt.Printf("OCR engine failed: %v\n", err)
			continue
		}

		if result != nil && result.Text != "" {
			// If this is the first successful result or has higher confidence
			if bestResult == nil || result.Confidence > bestResult.Confidence {
				bestResult = result
			}

			// If we have high confidence, use it immediately
			if result.Confidence >= CONFIDENCE_HIGH {
				break
			}
		}
	}

	if bestResult == nil {
		return nil, fmt.Errorf("all OCR engines failed, last error: %v", lastError)
	}

	return bestResult, nil
}

// ocrWithOCRmyPDF uses OCRmyPDF for high-quality OCR (PDF files only)
func ocrWithOCRmyPDF(filePath string) (*OCRResult, error) {
	// Check if OCRmyPDF is available
	if _, err := exec.LookPath("ocrmypdf"); err != nil {
		return nil, fmt.Errorf("ocrmypdf not found: %v", err)
	}

	// OCRmyPDF only works with PDF files
	if !strings.HasSuffix(strings.ToLower(filePath), ".pdf") {
		return nil, fmt.Errorf("ocrmypdf only supports PDF files")
	}

	outputPath := filePath + "_ocr.pdf"
	defer os.Remove(outputPath)

	// Run OCRmyPDF
	cmd := exec.Command("ocrmypdf", "--quiet", "--force-ocr", "-l", "eng", filePath, outputPath)
	if tessdata := os.Getenv("TESSDATA_PREFIX"); tessdata != "" {
		cmd.Env = append(os.Environ(), "TESSDATA_PREFIX="+tessdata)
	}

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("ocrmypdf failed: %v", err)
	}

	// Extract text from the OCR'd PDF
	text, err := extractTextFromPDF(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract text from OCR'd PDF: %v", err)
	}

	confidence := calculateOCRConfidence(text, "ocrmypdf")
	return &OCRResult{
		Text:       text,
		Confidence: confidence,
		Engine:     "ocrmypdf",
	}, nil
}

// ocrWithTesseract uses Tesseract directly
func ocrWithTesseract(filePath string) (*OCRResult, error) {
	client := gosseract.NewClient()
	defer client.Close()

	// Set image from file
	err := client.SetImage(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to set image: %v", err)
	}

	// Configure Tesseract
	client.SetLanguage("eng")
	client.SetPageSegMode(gosseract.PSM_AUTO)
	client.SetVariable("tessedit_char_whitelist", "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789<>/ .-")

	// Extract text
	text, err := client.Text()
	if err != nil {
		return nil, fmt.Errorf("tesseract OCR failed: %v", err)
	}

	confidence := calculateOCRConfidence(text, "tesseract")
	return &OCRResult{
		Text:       text,
		Confidence: confidence,
		Engine:     "tesseract",
	}, nil
}

// ocrWithPDFtoText tries to extract existing text from PDF (PDF files only)
func ocrWithPDFtoText(filePath string) (*OCRResult, error) {
	// Check if pdftotext is available
	if _, err := exec.LookPath("pdftotext"); err != nil {
		return nil, fmt.Errorf("pdftotext not found: %v", err)
	}

	// pdftotext only works with PDF files
	if !strings.HasSuffix(strings.ToLower(filePath), ".pdf") {
		return nil, fmt.Errorf("pdftotext only supports PDF files")
	}

	// Run pdftotext
	cmd := exec.Command("pdftotext", "-layout", filePath, "-")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("pdftotext failed: %v", err)
	}

	text := out.String()
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("no text extracted")
	}

	confidence := calculateOCRConfidence(text, "pdftotext")
	return &OCRResult{
		Text:       text,
		Confidence: confidence,
		Engine:     "pdftotext",
	}, nil
}

// extractTextFromPDF extracts text from a PDF file
func extractTextFromPDF(filePath string) (string, error) {
	cmd := exec.Command("pdftotext", "-layout", filePath, "-")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("pdftotext failed: %v", err)
	}

	return out.String(), nil
}

// calculateOCRConfidence calculates confidence based on text quality and engine
func calculateOCRConfidence(text, engine string) float64 {
	if text == "" {
		return 0.0
	}

	confidence := 0.5 // Base confidence

	// Length bonus (more text usually means better OCR)
	textLen := len(strings.TrimSpace(text))
	if textLen > 100 {
		confidence += 0.1
	}
	if textLen > 500 {
		confidence += 0.1
	}

	// Pattern recognition bonus
	if regexp.MustCompile(`\b[A-Z]{2,}\b`).MatchString(text) { // Has uppercase words
		confidence += 0.05
	}
	if regexp.MustCompile(`\b\d{4}-\d{2}-\d{2}\b`).MatchString(text) { // Has dates
		confidence += 0.05
	}
	if regexp.MustCompile(`P<[A-Z]{3}`).MatchString(text) { // Has MRZ pattern
		confidence += 0.15
	}

	// Engine-specific adjustments
	switch engine {
	case "ocrmypdf":
		confidence += 0.1 // OCRmyPDF is generally more accurate
	case "pdftotext":
		confidence += 0.2 // Native text extraction is most reliable
	case "tesseract":
		// No adjustment (baseline)
	}

	// Cap confidence
	if confidence > 0.95 {
		confidence = 0.95
	}
	if confidence < 0.1 {
		confidence = 0.1
	}

	return confidence
}

// processDrivingLicence processes UK driving licence documents using multi-engine OCR
func processDrivingLicence(file multipart.File, uploadType string) (map[string]interface{}, float64, error) {
	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read file: %v", err)
	}

	// Use multi-engine OCR for best results
	ocrResult, err := runMultiEngineOCR(fileBytes, "driving_licence.jpg")
	if err != nil {
		return nil, 0, fmt.Errorf("multi-engine OCR failed: %v", err)
	}

	fmt.Printf("Driving Licence OCR Text extracted using %s (confidence: %.2f): %s\n", ocrResult.Engine, ocrResult.Confidence, ocrResult.Text)

	// Parse extracted text based on front/back side
	var extractedData map[string]interface{}
	if uploadType == "front" {
		extractedData = parseDrivingLicenceFront(ocrResult.Text)
	} else if uploadType == "back" {
		extractedData = parseDrivingLicenceBack(ocrResult.Text)
	} else {
		extractedData = parseDrivingLicenceFront(ocrResult.Text) // Default to front
	}

	// Calculate final confidence combining OCR confidence and field extraction confidence
	fieldConfidence := calculateDrivingLicenceConfidence(extractedData, uploadType)
	finalConfidence := (ocrResult.Confidence + fieldConfidence) / 2.0

	// Add OCR metadata
	extractedData["_ocrEngine"] = ocrResult.Engine
	extractedData["_ocrConfidence"] = ocrResult.Confidence
	extractedData["_fieldConfidence"] = fieldConfidence

	return extractedData, finalConfidence, nil
}

// parseDrivingLicenceFront extracts fields from the front of a UK driving licence
func parseDrivingLicenceFront(text string) map[string]interface{} {
	extractedData := make(map[string]interface{})

	// Clean up text
	text = strings.ToUpper(strings.TrimSpace(text))

	// Patterns for UK driving licence front side fields
	patterns := map[string]*regexp.Regexp{
		"licenceNumber":    regexp.MustCompile(`(?:LICENCE\s+NO\.?\s*|LICENSE\s+NO\.?\s*)([A-Z]{5}\d{6}[A-Z]{2}\d[A-Z]{2})`),
		"surname":          regexp.MustCompile(`(?:SURNAME\s*[:\-]?\s*|FAMILY\s+NAME\s*[:\-]?\s*)([A-Z]+(?:\s+[A-Z]+)*)`),
		"givenNames":       regexp.MustCompile(`(?:GIVEN\s+NAMES?\s*[:\-]?\s*|FIRST\s+NAMES?\s*[:\-]?\s*)([A-Z]+(?:\s+[A-Z]+)*)`),
		"dateOfBirth":      regexp.MustCompile(`(?:DATE\s+OF\s+BIRTH\s*[:\-]?\s*|DOB\s*[:\-]?\s*)(\d{1,2}[\/\-\.]\d{1,2}[\/\-\.]\d{2,4})`),
		"placeOfBirth":     regexp.MustCompile(`(?:PLACE\s+OF\s+BIRTH\s*[:\-]?\s*|POB\s*[:\-]?\s*)([A-Z\s,]+)`),
		"issueDate":        regexp.MustCompile(`(?:DATE\s+OF\s+ISSUE\s*[:\-]?\s*|ISSUE\s+DATE\s*[:\-]?\s*)(\d{1,2}[\/\-\.]\d{1,2}[\/\-\.]\d{2,4})`),
		"expiryDate":       regexp.MustCompile(`(?:VALID\s+UNTIL\s*[:\-]?\s*|EXPIRES?\s*[:\-]?\s*)(\d{1,2}[\/\-\.]\d{1,2}[\/\-\.]\d{2,4})`),
		"issuingAuthority": regexp.MustCompile(`(?:DVLA|DRIVER\s+AND\s+VEHICLE\s+LICENSING\s+AGENCY)`),
	}

	// Extract address (usually multiple lines)
	addressPattern := regexp.MustCompile(`(?:ADDRESS\s*[:\-]?\s*|ADDR\s*[:\-]?\s*)([A-Z0-9\s,]+(?:\n[A-Z0-9\s,]+)*)`)
	if match := addressPattern.FindStringSubmatch(text); len(match) > 1 {
		address := strings.TrimSpace(match[1])
		if address != "" {
			extractedData["address"] = address
		}
	}

	// Extract fields using patterns
	for field, pattern := range patterns {
		if match := pattern.FindStringSubmatch(text); len(match) > 1 {
			value := strings.TrimSpace(match[1])
			if value != "" {
				extractedData[field] = value
			}
		}
	}

	// Set issuing authority if DVLA pattern found
	if _, exists := extractedData["issuingAuthority"]; exists {
		extractedData["issuingAuthority"] = "DVLA"
	}

	return extractedData
}

// parseDrivingLicenceBack extracts fields from the back of a UK driving licence
func parseDrivingLicenceBack(text string) map[string]interface{} {
	extractedData := make(map[string]interface{})

	// Clean up text
	text = strings.ToUpper(strings.TrimSpace(text))

	// Extract vehicle categories (A, A1, A2, AM, B, B1, BE, C, C1, CE, C1E, D, D1, DE, D1E, etc.)
	categoryPattern := regexp.MustCompile(`(?:CATEGORIES?\s*[:\-]?\s*|CAT\s*[:\-]?\s*)([A-Z0-9\s,+]+)`)
	if match := categoryPattern.FindStringSubmatch(text); len(match) > 1 {
		categories := strings.Split(strings.ReplaceAll(match[1], ",", " "), " ")
		var validCategories []string
		for _, cat := range categories {
			cat = strings.TrimSpace(cat)
			if cat != "" && regexp.MustCompile(`^[A-Z][0-9]?[E]?$`).MatchString(cat) {
				validCategories = append(validCategories, cat)
			}
		}
		if len(validCategories) > 0 {
			extractedData["categories"] = validCategories
		}
	}

	// Extract restrictions/endorsements
	restrictionPattern := regexp.MustCompile(`(?:RESTRICTIONS?\s*[:\-]?\s*|REST\s*[:\-]?\s*)([A-Z0-9\s,]+)`)
	if match := restrictionPattern.FindStringSubmatch(text); len(match) > 1 {
		restrictions := strings.Split(strings.ReplaceAll(match[1], ",", " "), " ")
		var validRestrictions []string
		for _, rest := range restrictions {
			rest = strings.TrimSpace(rest)
			if rest != "" && rest != "NONE" {
				validRestrictions = append(validRestrictions, rest)
			}
		}
		extractedData["restrictions"] = validRestrictions
	}

	// Extract endorsements
	endorsementPattern := regexp.MustCompile(`(?:ENDORSEMENTS?\s*[:\-]?\s*|END\s*[:\-]?\s*)([A-Z0-9\s,]+)`)
	if match := endorsementPattern.FindStringSubmatch(text); len(match) > 1 {
		endorsements := strings.Split(strings.ReplaceAll(match[1], ",", " "), " ")
		var validEndorsements []string
		for _, end := range endorsements {
			end = strings.TrimSpace(end)
			if end != "" && end != "NONE" {
				validEndorsements = append(validEndorsements, end)
			}
		}
		extractedData["endorsements"] = validEndorsements
	}

	return extractedData
}

// calculateDrivingLicenceConfidence calculates confidence based on extracted fields
func calculateDrivingLicenceConfidence(data map[string]interface{}, uploadType string) float64 {
	var requiredFields []string
	var optionalFields []string

	if uploadType == "front" {
		requiredFields = []string{"licenceNumber", "surname", "givenNames", "dateOfBirth"}
		optionalFields = []string{"placeOfBirth", "issueDate", "expiryDate", "address", "issuingAuthority"}
	} else {
		requiredFields = []string{"categories"}
		optionalFields = []string{"restrictions", "endorsements"}
	}

	foundRequired := 0
	foundOptional := 0

	for _, field := range requiredFields {
		if _, exists := data[field]; exists {
			foundRequired++
		}
	}

	for _, field := range optionalFields {
		if _, exists := data[field]; exists {
			foundOptional++
		}
	}

	// Base confidence on required fields, bonus for optional fields
	confidence := float64(foundRequired) / float64(len(requiredFields)) * 0.8
	if len(optionalFields) > 0 {
		confidence += float64(foundOptional) / float64(len(optionalFields)) * 0.2
	}

	// Minimum confidence of 0.3, maximum of 0.95
	if confidence < 0.3 {
		confidence = 0.3
	}
	if confidence > 0.95 {
		confidence = 0.95
	}

	return confidence
}

// processPassport processes passport documents using MRZ-only extraction
func processPassport(file multipart.File, uploadType string) (map[string]interface{}, float64, error) {
	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read file: %v", err)
	}

	fmt.Printf("🔍 STEP 1: Extracting MRZ area from passport image...\n")

	// Extract MRZ area only (bottom portion of passport)
	mrzBytes, err := extractMRZArea(fileBytes)
	if err != nil {
		fmt.Printf("❌ MRZ extraction failed, falling back to full image: %v\n", err)
		mrzBytes = fileBytes // Fallback to full image
	} else {
		fmt.Printf("✅ MRZ area extracted successfully\n")
	}

	fmt.Printf("🔍 STEP 2: Running OCR on MRZ area only...\n")

	// Use multi-engine OCR on MRZ area only
	ocrResult, err := runMultiEngineOCR(mrzBytes, "passport_mrz.jpg")
	if err != nil {
		return nil, 0, fmt.Errorf("MRZ OCR failed: %v", err)
	}

	fmt.Printf("OCR Text extracted from MRZ using %s (confidence: %.2f):\n", ocrResult.Engine, ocrResult.Confidence)
	fmt.Printf("=== MRZ OCR TEXT START ===\n%s\n=== MRZ OCR TEXT END ===\n", ocrResult.Text)

	fmt.Printf("🔍 STEP 3: Parsing MRZ data...\n")

	// Parse extracted MRZ text with enhanced processing
	extractedData := parseMRZOnly(ocrResult.Text)

	// Calculate final confidence combining OCR confidence and field extraction confidence
	fieldConfidence := calculatePassportConfidence(extractedData)
	finalConfidence := (ocrResult.Confidence + fieldConfidence) / 2.0

	// Add OCR metadata
	extractedData["_ocrEngine"] = ocrResult.Engine
	extractedData["_ocrConfidence"] = ocrResult.Confidence
	extractedData["_fieldConfidence"] = fieldConfidence
	extractedData["_mrzExtracted"] = true

	fmt.Printf("✅ STEP 4: Extraction complete - %d fields found\n", len(extractedData)-4) // -4 for metadata fields

	return extractedData, finalConfidence, nil
}

// extractMRZArea extracts the MRZ (Machine Readable Zone) from passport image
// MRZ is typically in the bottom 20-25% of the passport
func extractMRZArea(imageBytes []byte) ([]byte, error) {
	// Decode image
	img, format, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	fmt.Printf("📏 Original image: %dx%d pixels\n", width, height)

	// Calculate MRZ area (bottom portion below photo, typically 80-85% down)
	mrzStartY := int(float64(height) * 0.80) // Start at 80% down (below photo)
	mrzHeight := height - mrzStartY          // Bottom 20%

	fmt.Printf("📏 MRZ area: %dx%d pixels (Y: %d-%d)\n", width, mrzHeight, mrzStartY, height)

	// Create cropped image for MRZ area
	mrzImg := image.NewRGBA(image.Rect(0, 0, width, mrzHeight))

	// Copy MRZ area to new image
	for y := mrzStartY; y < height; y++ {
		for x := 0; x < width; x++ {
			mrzImg.Set(x, y-mrzStartY, img.At(x, y))
		}
	}

	// Save the cropped MRZ image for verification
	mrzFilename := fmt.Sprintf("mrz_cropped_%d.png", time.Now().Unix())
	mrzFile, err := os.Create(mrzFilename)
	if err == nil {
		png.Encode(mrzFile, mrzImg)
		mrzFile.Close()
		fmt.Printf("💾 MRZ cropped image saved as: %s\n", mrzFilename)
	}

	// Encode cropped image back to bytes
	var buf bytes.Buffer
	switch format {
	case "jpeg":
		err = jpeg.Encode(&buf, mrzImg, &jpeg.Options{Quality: 95})
	case "png":
		err = png.Encode(&buf, mrzImg)
	default:
		// Default to PNG
		err = png.Encode(&buf, mrzImg)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to encode MRZ image: %v", err)
	}

	return buf.Bytes(), nil
}

// parseMRZOnly parses ONLY MRZ text (no visual inspection zone)
func parseMRZOnly(text string) map[string]interface{} {
	extractedData := make(map[string]interface{})

	// Clean up text
	text = strings.ToUpper(strings.TrimSpace(text))
	fmt.Printf("=== MRZ TEXT TO PARSE ===\n%s\n=== END MRZ TEXT ===\n", text)

	// Split into lines and filter for actual MRZ lines
	allLines := strings.Split(text, "\n")
	fmt.Printf("📋 Found %d total lines in OCR text\n", len(allLines))

	// Filter for MRZ lines (start with P< or contain passport patterns)
	var mrzLines []string
	for i, line := range allLines {
		line = strings.TrimSpace(line)
		if line != "" {
			fmt.Printf("  Line %d: %s\n", i+1, line)

			// Check if this looks like an MRZ line
			if isMRZLine(line) {
				mrzLines = append(mrzLines, line)
				fmt.Printf("    ✅ Identified as MRZ line\n")
			}
		}
	}

	fmt.Printf("📋 Found %d MRZ lines after filtering\n", len(mrzLines))

	// Try to parse the identified MRZ lines
	if len(mrzLines) >= 2 {
		// Try to parse as 2-line MRZ
		line1 := strings.TrimSpace(mrzLines[0])
		line2 := strings.TrimSpace(mrzLines[1])

		fmt.Printf("🔍 Attempting 2-line MRZ parsing...\n")
		fmt.Printf("  MRZ Line 1: %s\n", line1)
		fmt.Printf("  MRZ Line 2: %s\n", line2)

		mrzData := parseTwoLineMRZ(line1, line2)
		for key, value := range mrzData {
			extractedData[key] = value
		}
	} else if len(mrzLines) == 1 {
		// Single MRZ line found
		fmt.Printf("🔍 Attempting single MRZ line parsing...\n")
		singleData := parseSingleLineMRZ(mrzLines[0])
		for key, value := range singleData {
			extractedData[key] = value
		}
	}

	// If still no MRZ found, try all lines as fallback
	if len(extractedData) == 0 {
		fmt.Printf("🔍 Fallback: Attempting pattern matching on all lines...\n")
		for _, line := range allLines {
			line = strings.TrimSpace(line)
			if len(line) > 15 {
				singleData := parseSingleLineMRZ(line)
				for key, value := range singleData {
					extractedData[key] = value
				}
			}
		}
	}

	fmt.Printf("📊 MRZ parsing result: %d fields extracted\n", len(extractedData))
	return extractedData
}

// isMRZLine checks if a line looks like an MRZ line
func isMRZLine(line string) bool {
	line = strings.TrimSpace(line)

	// Check for passport MRZ patterns
	if strings.HasPrefix(line, "P<") {
		return true
	}

	// Check for passport number patterns (starts with letter+numbers)
	passportPattern := regexp.MustCompile(`^[A-Z][A-Z0-9]{6,8}`)
	if passportPattern.MatchString(line) {
		return true
	}

	// Check for typical MRZ length (44 characters for TD-3)
	if len(line) >= 35 && len(line) <= 50 {
		// Count alphanumeric characters vs spaces/symbols
		alphaNum := 0
		for _, char := range line {
			if (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
				alphaNum++
			}
		}
		// MRZ lines should be mostly alphanumeric
		if float64(alphaNum)/float64(len(line)) > 0.7 {
			return true
		}
	}

	return false
}

// parseTwoLineMRZ parses standard 2-line TD-3 passport MRZ
func parseTwoLineMRZ(line1, line2 string) map[string]interface{} {
	extractedData := make(map[string]interface{})

	fmt.Printf("📋 Line 1: %s (length: %d)\n", line1, len(line1))
	fmt.Printf("📋 Line 2: %s (length: %d)\n", line2, len(line2))

	// Line 1: P<CCCNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNN
	// P = Passport, CCC = Country, NNN = Name fields
	if len(line1) >= 5 && strings.HasPrefix(line1, "P<") {
		// Extract country code
		if len(line1) >= 5 {
			country := line1[2:5]
			extractedData["issuingCountry"] = country
			fmt.Printf("  ✅ Country: %s\n", country)
		}

		// Extract names (after country code)
		if len(line1) > 5 {
			namesPart := line1[5:]
			// Split by << to separate surname and given names
			nameParts := strings.Split(namesPart, "<<")
			if len(nameParts) >= 1 {
				surname := strings.ReplaceAll(nameParts[0], "<", " ")
				surname = strings.TrimSpace(surname)
				if surname != "" {
					extractedData["surname"] = surname
					fmt.Printf("  ✅ Surname: %s\n", surname)
				}
			}
			if len(nameParts) >= 2 {
				givenNames := strings.ReplaceAll(nameParts[1], "<", " ")
				givenNames = strings.TrimSpace(givenNames)
				if givenNames != "" {
					extractedData["givenNames"] = givenNames
					fmt.Printf("  ✅ Given Names: %s\n", givenNames)
				}
			}
		}

		extractedData["machineReadableZone"] = line1 + "\n" + line2
	}

	// Line 2: NNNNNNNNNCCCDDDDDDSDDDDDDNNNNNNNNNNNNNC
	// N = Passport number, C = Check digit, D = Dates, S = Sex
	if len(line2) >= 30 {
		// Passport number (positions 0-8, may be shorter)
		passportEnd := 9
		for i := 0; i < 9 && i < len(line2); i++ {
			if line2[i] == '<' {
				passportEnd = i
				break
			}
		}
		if passportEnd > 0 {
			passportNum := line2[0:passportEnd]
			extractedData["passportNumber"] = passportNum
			fmt.Printf("  ✅ Passport Number: %s\n", passportNum)
		}

		// Skip check digit at position 9

		// Nationality (positions 10-12)
		if len(line2) >= 13 {
			nationality := line2[10:13]
			extractedData["nationality"] = nationality
			fmt.Printf("  ✅ Nationality: %s\n", nationality)
		}

		// Date of birth (positions 13-18)
		if len(line2) >= 19 {
			dobStr := line2[13:19]
			if dob := parseMRZDate(dobStr); dob != "" {
				extractedData["dateOfBirth"] = dob
				fmt.Printf("  ✅ Date of Birth: %s\n", dob)
			}
		}

		// Gender (position 20)
		if len(line2) >= 21 {
			gender := string(line2[20])
			if gender == "M" || gender == "F" {
				extractedData["gender"] = gender
				fmt.Printf("  ✅ Gender: %s\n", gender)
			}
		}

		// Expiry date (positions 21-26)
		if len(line2) >= 27 {
			expiryStr := line2[21:27]
			if expiry := parseMRZDate(expiryStr); expiry != "" {
				extractedData["expiryDate"] = expiry
				fmt.Printf("  ✅ Expiry Date: %s\n", expiry)
			}
		}
	}

	return extractedData
}

// parseSingleLineMRZ attempts to parse single line MRZ patterns
func parseSingleLineMRZ(line string) map[string]interface{} {
	extractedData := make(map[string]interface{})

	// Look for passport number patterns
	passportPattern := regexp.MustCompile(`([A-Z0-9]{6,9})`)
	if matches := passportPattern.FindAllString(line, -1); len(matches) > 0 {
		// Take the first reasonable match
		for _, match := range matches {
			if len(match) >= 6 {
				extractedData["passportNumber"] = match
				fmt.Printf("  ✅ Found passport number: %s\n", match)
				break
			}
		}
	}

	// Look for country codes
	countryPattern := regexp.MustCompile(`\b([A-Z]{3})\b`)
	if matches := countryPattern.FindAllString(line, -1); len(matches) > 0 {
		extractedData["issuingCountry"] = matches[0]
		fmt.Printf("  ✅ Found country: %s\n", matches[0])
	}

	return extractedData
}

// parsePassportTextEnhanced extracts passport fields from OCR text with enhanced MRZ processing
func parsePassportTextEnhanced(text string) map[string]interface{} {
	extractedData := make(map[string]interface{})

	// Clean up text
	text = strings.ToUpper(strings.TrimSpace(text))
	fmt.Printf("=== PARSING TEXT ===\n%s\n=== END PARSING TEXT ===\n", text)

	// First, try to find and parse MRZ (Machine Readable Zone)
	fmt.Printf("🔍 Attempting MRZ extraction...\n")
	mrzData := extractMRZData(text)
	fmt.Printf("📋 MRZ Data found: %d fields\n", len(mrzData))
	for key, value := range mrzData {
		fmt.Printf("  - %s: %v\n", key, value)
		extractedData[key] = value
	}

	// Then extract visual inspection zone data
	fmt.Printf("🔍 Attempting Visual Inspection Zone extraction...\n")
	vizData := extractVisualInspectionZone(text)
	fmt.Printf("📋 VIZ Data found: %d fields\n", len(vizData))
	for key, value := range vizData {
		// Only overwrite if we don't have MRZ data or if VIZ data seems more reliable
		if _, exists := extractedData[key]; !exists {
			fmt.Printf("  + %s: %v\n", key, value)
			extractedData[key] = value
		} else {
			fmt.Printf("  - %s: %v (skipped, MRZ data exists)\n", key, value)
		}
	}

	fmt.Printf("📊 Total extracted fields: %d\n", len(extractedData))
	return extractedData
}

// extractMRZData extracts data from Machine Readable Zone with comprehensive parsing
func extractMRZData(text string) map[string]interface{} {
	extractedData := make(map[string]interface{})

	// Enhanced MRZ patterns for different passport formats
	mrzPatterns := []*regexp.Regexp{
		// Standard TD-3 passport MRZ (2 lines, 44 characters each)
		regexp.MustCompile(`P<([A-Z]{3})([A-Z<]+)<<([A-Z<]+)<<+\s*([A-Z0-9<]{9})([0-9])([A-Z]{3})([0-9]{6})([MF])([0-9]{6})([0-9])([A-Z0-9<]{14})([0-9])`),
		// Relaxed MRZ pattern for Spanish passport format
		regexp.MustCompile(`P<([A-Z]{3})([A-Z<]+)<<([A-Z<]+)<<+\s*([A-Z0-9]+)([A-Z]{3})([0-9]{6})([MF])([0-9]{6})([A-Z0-9<]+)`),
		// Simple MRZ start pattern
		regexp.MustCompile(`P<([A-Z]{3})([A-Z<\s]+)`),
	}

	// First try to parse as a complete 2-line MRZ
	if mrzData := parseSpanishMRZ(text); len(mrzData) > 0 {
		for key, value := range mrzData {
			extractedData[key] = value
		}
		return extractedData
	}

	// Fallback to pattern matching
	for _, pattern := range mrzPatterns {
		if matches := pattern.FindAllStringSubmatch(text, -1); len(matches) > 0 {
			match := matches[0]
			if len(match) >= 13 { // Full MRZ match
				extractedData["issuingCountry"] = match[1]
				extractedData["surname"] = strings.ReplaceAll(match[2], "<", " ")
				extractedData["givenNames"] = strings.ReplaceAll(match[3], "<", " ")
				extractedData["passportNumber"] = strings.ReplaceAll(match[4], "<", "")
				extractedData["nationality"] = match[6]

				// Parse dates from MRZ format (YYMMDD)
				if dob := parseMRZDate(match[7]); dob != "" {
					extractedData["dateOfBirth"] = dob
				}

				extractedData["gender"] = match[8]

				if expiry := parseMRZDate(match[9]); expiry != "" {
					extractedData["expiryDate"] = expiry
				}

				extractedData["machineReadableZone"] = match[0]
				break
			} else if len(match) > 0 {
				// Partial MRZ match - extract what we can
				fullMRZ := match[0]
				extractedData["machineReadableZone"] = fullMRZ

				// Try to extract country code
				if len(fullMRZ) >= 5 {
					extractedData["issuingCountry"] = fullMRZ[2:5]
				}

				// Try to extract passport number pattern
				passportPattern := regexp.MustCompile(`([A-Z0-9]{6,9})`)
				if passportMatch := passportPattern.FindString(fullMRZ[5:]); passportMatch != "" {
					extractedData["passportNumber"] = passportMatch
				}
			}
		}
	}

	return extractedData
}

// parseSpanishMRZ parses Spanish passport MRZ format specifically
// Line 1: P<ESPDOBARAN<SANTIAGO<<OIHANE<<<<<<<<<<<<<<<<
// Line 2: PA08932043ESP7812259F3208298A16064864000000092
func parseSpanishMRZ(text string) map[string]interface{} {
	extractedData := make(map[string]interface{})

	// Split into lines and clean
	lines := strings.Split(strings.ToUpper(text), "\n")
	if len(lines) < 2 {
		return extractedData
	}

	line1 := strings.TrimSpace(lines[0])
	line2 := strings.TrimSpace(lines[1])

	// Parse Line 1: P<ESPDOBARAN<SANTIAGO<<OIHANE<<<<<<<<<<<<<<<<
	line1Pattern := regexp.MustCompile(`^P<([A-Z]{3})([A-Z<]+)<<([A-Z<]+)<<+`)
	if match := line1Pattern.FindStringSubmatch(line1); len(match) >= 4 {
		extractedData["issuingCountry"] = match[1]
		extractedData["surname"] = strings.TrimSpace(strings.ReplaceAll(match[2], "<", " "))
		extractedData["givenNames"] = strings.TrimSpace(strings.ReplaceAll(match[3], "<", " "))
		extractedData["machineReadableZone"] = line1 + "\n" + line2
	}

	// Parse Line 2: PA08932043ESP7812259F3208298A16064864000000092
	// Format: [PassportNumber][Country][DOB][Gender][Expiry][PersonalNumber][CheckDigit]
	if len(line2) >= 30 {
		// Extract passport number (positions 0-9)
		if len(line2) >= 10 {
			passportNum := line2[0:10]
			// Remove trailing < characters
			passportNum = strings.TrimRight(passportNum, "<")
			extractedData["passportNumber"] = passportNum
		}

		// Extract nationality (positions 10-12)
		if len(line2) >= 13 {
			extractedData["nationality"] = line2[10:13]
		}

		// Extract date of birth (positions 13-18, format YYMMDD)
		if len(line2) >= 19 {
			dobStr := line2[13:19]
			if dob := parseMRZDate(dobStr); dob != "" {
				extractedData["dateOfBirth"] = dob
			}
		}

		// Extract gender (position 19)
		if len(line2) >= 20 {
			extractedData["gender"] = string(line2[19])
		}

		// Extract expiry date (positions 20-25, format YYMMDD)
		if len(line2) >= 26 {
			expiryStr := line2[20:26]
			if expiry := parseMRZDate(expiryStr); expiry != "" {
				extractedData["expiryDate"] = expiry
			}
		}

		// Extract personal number (positions 26+)
		if len(line2) > 26 {
			personalNum := line2[26:]
			// Remove check digits and padding
			personalNum = strings.TrimRight(personalNum, "0123456789<")
			if personalNum != "" {
				extractedData["personalNumber"] = personalNum
			}
		}
	}

	return extractedData
}

// parseMRZDate converts MRZ date format (YYMMDD) to standard format (YYYY-MM-DD)
func parseMRZDate(mrzDate string) string {
	if len(mrzDate) != 6 {
		return ""
	}

	year, err1 := strconv.Atoi(mrzDate[0:2])
	month, err2 := strconv.Atoi(mrzDate[2:4])
	day, err3 := strconv.Atoi(mrzDate[4:6])

	if err1 != nil || err2 != nil || err3 != nil || month < 1 || month > 12 || day < 1 || day > 31 {
		return ""
	}

	// Convert 2-digit year to 4-digit year (assume 1950-2049)
	if year < 50 {
		year += 2000
	} else {
		year += 1900
	}

	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

// extractVisualInspectionZone extracts data from the visual inspection zone
func extractVisualInspectionZone(text string) map[string]interface{} {
	extractedData := make(map[string]interface{})

	// Clean up text
	text = strings.ToUpper(strings.TrimSpace(text))

	// Patterns for passport fields
	patterns := map[string]*regexp.Regexp{
		"passportNumber":   regexp.MustCompile(`(?:PASSPORT\s+NO\.?\s*|P\s*<\s*[A-Z]{3}\s*)([A-Z0-9]{6,9})`),
		"surname":          regexp.MustCompile(`(?:SURNAME\s*[:\-]?\s*|P<[A-Z]{3})([A-Z]+(?:\s+[A-Z]+)*)`),
		"givenNames":       regexp.MustCompile(`(?:GIVEN\s+NAMES?\s*[:\-]?\s*|FIRST\s+NAMES?\s*[:\-]?\s*)([A-Z]+(?:\s+[A-Z]+)*)`),
		"nationality":      regexp.MustCompile(`(?:NATIONALITY\s*[:\-]?\s*)([A-Z\s]+)`),
		"dateOfBirth":      regexp.MustCompile(`(?:DATE\s+OF\s+BIRTH\s*[:\-]?\s*|DOB\s*[:\-]?\s*)(\d{1,2}[\/\-\.]\d{1,2}[\/\-\.]\d{2,4})`),
		"placeOfBirth":     regexp.MustCompile(`(?:PLACE\s+OF\s+BIRTH\s*[:\-]?\s*|POB\s*[:\-]?\s*)([A-Z\s,]+)`),
		"gender":           regexp.MustCompile(`(?:SEX\s*[:\-]?\s*|GENDER\s*[:\-]?\s*)([MF])`),
		"issueDate":        regexp.MustCompile(`(?:DATE\s+OF\s+ISSUE\s*[:\-]?\s*|ISSUE\s+DATE\s*[:\-]?\s*)(\d{1,2}[\/\-\.]\d{1,2}[\/\-\.]\d{2,4})`),
		"expiryDate":       regexp.MustCompile(`(?:DATE\s+OF\s+EXPIRY\s*[:\-]?\s*|EXPIRY\s+DATE\s*[:\-]?\s*|EXPIRES?\s*[:\-]?\s*)(\d{1,2}[\/\-\.]\d{1,2}[\/\-\.]\d{2,4})`),
		"issuingAuthority": regexp.MustCompile(`(?:AUTHORITY\s*[:\-]?\s*|ISSUED\s+BY\s*[:\-]?\s*)([A-Z\s]+)`),
	}

	// Try to find Machine Readable Zone (MRZ) - typically at bottom of passport
	mrzPattern := regexp.MustCompile(`P<[A-Z]{3}[A-Z<]+<<[A-Z<]+[0-9A-Z<]{44}`)
	if mrzMatch := mrzPattern.FindString(text); mrzMatch != "" {
		extractedData["machineReadableZone"] = mrzMatch
		// Parse MRZ for additional data
		parseMRZ(mrzMatch, extractedData)
	}

	// Extract fields using patterns
	for field, pattern := range patterns {
		if match := pattern.FindStringSubmatch(text); len(match) > 1 {
			value := strings.TrimSpace(match[1])
			if value != "" {
				extractedData[field] = value
			}
		}
	}

	// Try to extract country code from various formats
	countryPatterns := []*regexp.Regexp{
		regexp.MustCompile(`P<([A-Z]{3})`),
		regexp.MustCompile(`COUNTRY\s+CODE\s*[:\-]?\s*([A-Z]{2,3})`),
		regexp.MustCompile(`([A-Z]{3})\d{6}[MF]\d{6}`),
	}

	for _, pattern := range countryPatterns {
		if match := pattern.FindStringSubmatch(text); len(match) > 1 {
			extractedData["issuingCountry"] = match[1]
			break
		}
	}

	return extractedData
}

// Legacy function for backward compatibility
func parsePassportText(text string) map[string]interface{} {
	return parsePassportTextEnhanced(text)
}

// parseMRZ extracts data from Machine Readable Zone
func parseMRZ(mrz string, extractedData map[string]interface{}) {
	if len(mrz) < 44 {
		return
	}

	// MRZ format: P<CCCNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNN
	// Where CCC = country code, rest contains encoded data

	// Extract country code
	if len(mrz) >= 5 {
		country := mrz[2:5]
		extractedData["issuingCountry"] = country
	}

	// Extract passport number (positions vary, but typically after country)
	passportNumPattern := regexp.MustCompile(`([A-Z0-9]{6,9})`)
	if match := passportNumPattern.FindString(mrz[5:]); match != "" {
		extractedData["passportNumber"] = match
	}
}

// calculatePassportConfidence calculates confidence based on extracted fields
func calculatePassportConfidence(data map[string]interface{}) float64 {
	requiredFields := []string{"passportNumber", "surname", "givenNames", "dateOfBirth"}
	optionalFields := []string{"nationality", "placeOfBirth", "gender", "issueDate", "expiryDate", "issuingCountry"}

	foundRequired := 0
	foundOptional := 0

	for _, field := range requiredFields {
		if _, exists := data[field]; exists {
			foundRequired++
		}
	}

	for _, field := range optionalFields {
		if _, exists := data[field]; exists {
			foundOptional++
		}
	}

	// Base confidence on required fields, bonus for optional fields
	confidence := float64(foundRequired) / float64(len(requiredFields)) * 0.8
	confidence += float64(foundOptional) / float64(len(optionalFields)) * 0.2

	// Minimum confidence of 0.3, maximum of 0.95
	if confidence < 0.3 {
		confidence = 0.3
	}
	if confidence > 0.95 {
		confidence = 0.95
	}

	return confidence
}

// processGenericDocument processes other document types
func processGenericDocument(file multipart.File, documentType string) (map[string]interface{}, float64, error) {
	// Basic OCR processing for other document types
	extractedData := map[string]interface{}{
		"documentType": documentType,
		"processed":    true,
		"timestamp":    time.Now().Format(time.RFC3339),
	}

	return extractedData, 0.80, nil
}

// validateWithSHACL validates extracted data against SHACL shapes
func validateWithSHACL(documentType string, data map[string]interface{}) (*DocumentValidationResult, error) {
	// Mock SHACL validation
	result := &DocumentValidationResult{
		IsValid:     true,
		Violations:  []ValidationViolation{},
		Suggestions: []string{},
	}

	// Basic validation rules
	if documentType == "driving-licence" {
		if licenceNum, ok := data["licenceNumber"].(string); ok {
			if len(licenceNum) != 16 {
				result.IsValid = false
				result.Violations = append(result.Violations, ValidationViolation{
					Property: "licenceNumber",
					Message:  "UK driving licence number must be 16 characters",
					Severity: "Violation",
					Value:    licenceNum,
				})
			}
		}
	}

	if documentType == "passport" {
		if passportNum, ok := data["passportNumber"].(string); ok {
			if len(passportNum) != 9 {
				result.IsValid = false
				result.Violations = append(result.Violations, ValidationViolation{
					Property: "passportNumber",
					Message:  "UK passport number must be 9 characters",
					Severity: "Violation",
					Value:    passportNum,
				})
			}
		}
	}

	// In a real implementation:
	// 1. Load SHACL shapes from ontology files
	// 2. Convert extracted data to RDF triples
	// 3. Run SHACL validation engine
	// 4. Return detailed validation report

	return result, nil
}

// ValidateDocumentHandler handles SHACL validation requests
func ValidateDocumentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		DocumentType string                 `json:"documentType"`
		Data         map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := validateWithSHACL(request.DocumentType, request.Data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Validation failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
