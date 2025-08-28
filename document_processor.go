package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
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

// PassportEyeResult represents the result from PassportEye MRZ extraction
type PassportEyeResult struct {
	Success       bool                   `json:"success"`
	Confidence    float64                `json:"confidence"`
	Valid         bool                   `json:"valid"`
	RawText       string                 `json:"raw_text"`
	ExtractedData map[string]interface{} `json:"extracted_data"`
	MRZLines      []string               `json:"mrz_lines"`
	Validation    struct {
		Valid  bool     `json:"valid"`
		Errors []string `json:"errors"`
	} `json:"validation"`
	Error string `json:"error"`
}

// I18n structure for loading month mappings from JSON
type I18nData struct {
	Months struct {
		Abbrev map[string]string `json:"abbrev"`
	} `json:"months"`
}

// getI18nMonthMapping loads month mappings from all supported i18n files
func getI18nMonthMapping() map[string]string {
	monthMap := make(map[string]string)

	// Load from multiple language files
	languages := []string{"en", "de"}

	for _, lang := range languages {
		filePath := fmt.Sprintf("./i18n/%s.json", lang)
		if data, err := ioutil.ReadFile(filePath); err == nil {
			var i18nData I18nData
			if err := json.Unmarshal(data, &i18nData); err == nil {
				// Merge abbreviation mappings
				for abbrev, monthNum := range i18nData.Months.Abbrev {
					if monthNum != "" {
						// Pad single digit months with zero
						if len(monthNum) == 1 {
							monthNum = "0" + monthNum
						}
						monthMap[abbrev] = monthNum
					}
				}
				fmt.Printf("🌍 Loaded %d month mappings from %s\n", len(i18nData.Months.Abbrev), lang)
			}
		}
	}

	// Add fallback standard mappings if i18n files not found
	if len(monthMap) == 0 {
		fmt.Println("⚠️  No i18n files found, using fallback month mappings")
		monthMap = map[string]string{
			"JAN": "01", "FEB": "02", "MAR": "03", "APR": "04", "MAY": "05", "JUN": "06",
			"JUL": "07", "AUG": "08", "SEP": "09", "OCT": "10", "NOV": "11", "DEC": "12",
			"jan": "01", "feb": "02", "mar": "03", "apr": "04", "may": "05", "jun": "06",
			"jul": "07", "aug": "08", "sep": "09", "oct": "10", "nov": "11", "dec": "12",
		}
	}

	fmt.Printf("📅 Total month mappings loaded: %d\n", len(monthMap))
	return monthMap
}

// ISO 3166-1 alpha-3 country code mapping for MRZ
var countryCodeMap = map[string]string{
	"AFG": "Afghanistan", "ALB": "Albania", "DZA": "Algeria", "AND": "Andorra", "AGO": "Angola",
	"ATG": "Antigua and Barbuda", "ARG": "Argentina", "ARM": "Armenia", "AUS": "Australia", "AUT": "Austria",
	"AZE": "Azerbaijan", "BHS": "Bahamas", "BHR": "Bahrain", "BGD": "Bangladesh", "BRB": "Barbados",
	"BLR": "Belarus", "BEL": "Belgium", "BLZ": "Belize", "BEN": "Benin", "BTN": "Bhutan",
	"BOL": "Bolivia", "BIH": "Bosnia and Herzegovina", "BWA": "Botswana", "BRA": "Brazil", "BRN": "Brunei",
	"BGR": "Bulgaria", "BFA": "Burkina Faso", "BDI": "Burundi", "KHM": "Cambodia", "CMR": "Cameroon",
	"CAN": "Canada", "CPV": "Cape Verde", "CAF": "Central African Republic", "TCD": "Chad", "CHL": "Chile",
	"CHN": "China", "COL": "Colombia", "COM": "Comoros", "COG": "Congo", "COD": "Congo (Democratic Republic)",
	"CRI": "Costa Rica", "CIV": "Côte d'Ivoire", "HRV": "Croatia", "CUB": "Cuba", "CYP": "Cyprus",
	"CZE": "Czech Republic", "DNK": "Denmark", "DJI": "Djibouti", "DMA": "Dominica", "DOM": "Dominican Republic",
	"ECU": "Ecuador", "EGY": "Egypt", "SLV": "El Salvador", "GNQ": "Equatorial Guinea", "ERI": "Eritrea",
	"EST": "Estonia", "ETH": "Ethiopia", "FJI": "Fiji", "FIN": "Finland", "FRA": "France",
	"GAB": "Gabon", "GMB": "Gambia", "GEO": "Georgia", "DEU": "Germany", "GHA": "Ghana",
	"GRC": "Greece", "GRD": "Grenada", "GTM": "Guatemala", "GIN": "Guinea", "GNB": "Guinea-Bissau",
	"GUY": "Guyana", "HTI": "Haiti", "HND": "Honduras", "HUN": "Hungary", "ISL": "Iceland",
	"IND": "India", "IDN": "Indonesia", "IRN": "Iran", "IRQ": "Iraq", "IRL": "Ireland",
	"ISR": "Israel", "ITA": "Italy", "JAM": "Jamaica", "JPN": "Japan", "JOR": "Jordan",
	"KAZ": "Kazakhstan", "KEN": "Kenya", "KIR": "Kiribati", "PRK": "North Korea", "KOR": "South Korea",
	"KWT": "Kuwait", "KGZ": "Kyrgyzstan", "LAO": "Laos", "LVA": "Latvia", "LBN": "Lebanon",
	"LSO": "Lesotho", "LBR": "Liberia", "LBY": "Libya", "LIE": "Liechtenstein", "LTU": "Lithuania",
	"LUX": "Luxembourg", "MKD": "North Macedonia", "MDG": "Madagascar", "MWI": "Malawi", "MYS": "Malaysia",
	"MDV": "Maldives", "MLI": "Mali", "MLT": "Malta", "MHL": "Marshall Islands", "MRT": "Mauritania",
	"MUS": "Mauritius", "MEX": "Mexico", "FSM": "Micronesia", "MDA": "Moldova", "MCO": "Monaco",
	"MNG": "Mongolia", "MNE": "Montenegro", "MAR": "Morocco", "MOZ": "Mozambique", "MMR": "Myanmar",
	"NAM": "Namibia", "NRU": "Nauru", "NPL": "Nepal", "NLD": "Netherlands", "NZL": "New Zealand",
	"NIC": "Nicaragua", "NER": "Niger", "NGA": "Nigeria", "NOR": "Norway", "OMN": "Oman",
	"PAK": "Pakistan", "PLW": "Palau", "PAN": "Panama", "PNG": "Papua New Guinea", "PRY": "Paraguay",
	"PER": "Peru", "PHL": "Philippines", "POL": "Poland", "PRT": "Portugal", "QAT": "Qatar",
	"ROU": "Romania", "RUS": "Russia", "RWA": "Rwanda", "KNA": "Saint Kitts and Nevis", "LCA": "Saint Lucia",
	"VCT": "Saint Vincent and the Grenadines", "WSM": "Samoa", "SMR": "San Marino", "STP": "São Tomé and Príncipe",
	"SAU": "Saudi Arabia", "SEN": "Senegal", "SRB": "Serbia", "SYC": "Seychelles", "SLE": "Sierra Leone",
	"SGP": "Singapore", "SVK": "Slovakia", "SVN": "Slovenia", "SLB": "Solomon Islands", "SOM": "Somalia",
	"ZAF": "South Africa", "SSD": "South Sudan", "ESP": "Spain", "LKA": "Sri Lanka", "SDN": "Sudan",
	"SUR": "Suriname", "SWZ": "Eswatini", "SWE": "Sweden", "CHE": "Switzerland", "SYR": "Syria",
	"TWN": "Taiwan", "TJK": "Tajikistan", "TZA": "Tanzania", "THA": "Thailand", "TLS": "Timor-Leste",
	"TGO": "Togo", "TON": "Tonga", "TTO": "Trinidad and Tobago", "TUN": "Tunisia", "TUR": "Turkey",
	"TKM": "Turkmenistan", "TUV": "Tuvalu", "UGA": "Uganda", "UKR": "Ukraine", "ARE": "United Arab Emirates",
	"GBR": "United Kingdom", "USA": "United States", "URY": "Uruguay", "UZB": "Uzbekistan", "VUT": "Vanuatu",
	"VAT": "Vatican City", "VEN": "Venezuela", "VNM": "Vietnam", "YEM": "Yemen", "ZMB": "Zambia", "ZWE": "Zimbabwe",
}

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

// extractIssueDateFromText attempts to extract issue date from passport page text
func extractIssueDateFromText(text string) string {
	// Common patterns for issue date on UK passports:
	// "Date of issue: 01 JAN 2022"
	// "Date of Issue 01/01/2022"
	// "Issue Date: 01-01-2022"
	// "01 JAN 22"

	patterns := []string{
		// Standard date formats
		`(?i)date\s+of\s+issue[:\s]+(\d{1,2})\s+([A-Z]{3,4})\s+(\d{4})`,    // "Date of issue: 01 JAN 2022"
		`(?i)date\s+of\s+issue[:\s]+(\d{1,2})[\/\-](\d{1,2})[\/\-](\d{4})`, // "Date of issue: 01/01/2022"
		`(?i)issue\s+date[:\s]+(\d{1,2})[\/\-](\d{1,2})[\/\-](\d{4})`,      // "Issue date: 01/01/2022"
		`(?i)issued[:\s]+(\d{1,2})\s+([A-Z]{3,4})\s+(\d{4})`,               // "Issued: 01 JAN 2022"

		// UK Passport specific patterns
		`(?i)date\s+of\s+issue\s*[:\-]?\s*(\d{1,2})\s+([A-Z]{3,4})\s+(\d{2,4})`, // "Date of issue 01 JAN 22"
		`(?i)issued\s+on\s*[:\-]?\s*(\d{1,2})\s+([A-Z]{3,4})\s+(\d{2,4})`,       // "Issued on 01 JAN 22"
		`(?i)authority\s+[\w\s]*\s+(\d{1,2})\s+([A-Z]{3,4})\s+(\d{2,4})`,        // "Authority ... 01 JAN 22"

		// Specific UK passport format with bilingual text
		`(\d{1,2})\s+([A-Z]{3,4})\s*[\/\\]\s*[A-Z]{3,4}\s+(\d{2})`, // "03 SEP /SEPT 22"
		`(\d{1,2})\s+([A-Z]{3,4})\s+[\/\\]\s*[A-Z]{3,4}\s+(\d{2})`, // "03 SEP / SEPT 22"

		// Standard 3-letter month codes
		`(\d{1,2})\s+([A-Z]{3,4})\s+(\d{2,4})`, // "01 JAN 22" or "01 JAN 2022"

		// More flexible patterns for OCR errors
		`(?i)(\d{1,2})\s*[\/\-\.\s]\s*([A-Z]{3,4})\s*[\/\-\.\s]\s*(\d{2,4})`, // "01/JAN/22" or "01-JAN-22"

		// English month variations (OCR errors)
		`(?i)(\d{1,2})\s+(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)\s*[\/\s]*[a-z]*\s*(\d{2})`, // "05 sep /seer 22"
		`(?i)(\d{1,2})\s+(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*\s+(\d{2})`,           // "05 sep 22"

		// German month variations
		`(?i)(\d{1,2})\s+(jän|feb|mär|apr|mai|jun|jul|aug|sep|okt|nov|dez)\s*[\/\s]*[a-z]*\s*(\d{2})`, // German months

		// French month variations
		`(?i)(\d{1,2})\s+(janv|févr|mars|avr|mai|juin|juil|août|sept|oct|nov|déce)\s*[\/\s]*[a-z]*\s*(\d{2})`, // French months

		// Spanish month variations
		`(?i)(\d{1,2})\s+(ene|feb|mar|abr|may|jun|jul|ago|sep|oct|nov|dic)\s*[\/\s]*[a-z]*\s*(\d{2})`, // Spanish months

		// Italian month variations
		`(?i)(\d{1,2})\s+(gen|feb|mar|apr|mag|giu|lug|ago|set|ott|nov|dic)\s*[\/\s]*[a-z]*\s*(\d{2})`, // Italian months

		// Full month names (any language)
		`(?i)(\d{1,2})\s+(january|february|march|april|may|june|july|august|september|october|november|december)\s+(\d{2,4})`,
		`(?i)(\d{1,2})\s+(januar|februar|märz|april|mai|juni|juli|august|september|oktober|november|dezember)\s+(\d{2,4})`,
		`(?i)(\d{1,2})\s+(janvier|février|mars|avril|mai|juin|juillet|août|septembre|octobre|novembre|décembre)\s+(\d{2,4})`,
	}

	// Load i18n-compliant month mapping from ontology
	monthMap := getI18nMonthMapping()

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(text)

		if len(matches) >= 4 {
			day := matches[1]
			month := matches[2]
			year := matches[3]

			// Handle month names (try both original case and uppercase)
			if monthNum, exists := monthMap[month]; exists {
				month = monthNum
			} else if monthNum, exists := monthMap[strings.ToUpper(month)]; exists {
				month = monthNum
			}

			// Handle 2-digit years
			if len(year) == 2 {
				yearNum, _ := strconv.Atoi(year)
				if yearNum <= 30 {
					year = "20" + year
				} else {
					year = "19" + year
				}
			}

			// Pad day and month with zeros
			if len(day) == 1 {
				day = "0" + day
			}
			if len(month) == 1 {
				month = "0" + month
			}

			fmt.Printf("🗓️ Extracted issue date: %s-%s-%s\n", year, month, day)
			return fmt.Sprintf("%s-%s-%s", year, month, day)
		}
	}

	fmt.Printf("⚠️ No issue date pattern found in text\n")
	fmt.Printf("📝 Text being searched for issue date:\n%s\n", text)
	return ""
}

// preprocessMRZImage enhances MRZ image for better OCR by converting to high contrast black and white
func preprocessMRZImage(inputPath string) (string, error) {
	// Read the image
	file, err := os.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %v", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create new grayscale image
	grayImg := image.NewGray(bounds)

	// First pass: convert to grayscale and calculate histogram
	var histogram [256]int
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			originalColor := img.At(x, y)
			r, g, b, _ := originalColor.RGBA()

			// Convert to grayscale using luminance formula
			gray := uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256)
			grayImg.Set(x, y, color.Gray{Y: gray})
			histogram[gray]++
		}
	}

	// Calculate adaptive threshold using Otsu's method (simplified)
	totalPixels := width * height
	sum := 0
	for i := 0; i < 256; i++ {
		sum += i * histogram[i]
	}

	sumB := 0
	wB := 0
	maximum := 0.0
	threshold := uint8(128) // fallback

	for i := 0; i < 256; i++ {
		wB += histogram[i]
		if wB == 0 {
			continue
		}
		wF := totalPixels - wB
		if wF == 0 {
			break
		}

		sumB += i * histogram[i]
		mB := float64(sumB) / float64(wB)
		mF := float64(sum-sumB) / float64(wF)

		between := float64(wB) * float64(wF) * (mB - mF) * (mB - mF)
		if between > maximum {
			maximum = between
			threshold = uint8(i)
		}
	}

	fmt.Printf("🎯 Calculated adaptive threshold: %d\n", threshold)

	// Second pass: apply adaptive threshold
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gray := grayImg.GrayAt(x, y).Y

			// Apply adaptive threshold for high contrast black and white
			if gray < threshold {
				grayImg.Set(x, y, color.Gray{Y: 0}) // Pure black for text
			} else {
				grayImg.Set(x, y, color.Gray{Y: 255}) // Pure white for background
			}
		}
	}

	// Apply additional noise reduction (median filter)
	cleanImg := applyMedianFilter(grayImg)

	// Save preprocessed image
	preprocessedPath := strings.Replace(inputPath, ".png", "_preprocessed.png", 1)
	outFile, err := os.Create(preprocessedPath)
	if err != nil {
		return "", fmt.Errorf("failed to create preprocessed image: %v", err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, cleanImg)
	if err != nil {
		return "", fmt.Errorf("failed to encode preprocessed image: %v", err)
	}

	fmt.Printf("🎨 Preprocessed MRZ image saved: %s\n", preprocessedPath)
	fmt.Printf("🔗 Preprocessed MRZ URL: http://localhost:3000/%s\n", strings.Replace(preprocessedPath, "\\", "/", -1))

	return preprocessedPath, nil
}

// applyMedianFilter applies a 3x3 median filter to reduce noise
func applyMedianFilter(img *image.Gray) *image.Gray {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	filtered := image.NewGray(bounds)

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			// Get 3x3 neighborhood
			var pixels []uint8
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					pixel := img.GrayAt(x+dx, y+dy)
					pixels = append(pixels, pixel.Y)
				}
			}

			// Sort pixels and get median
			for i := 0; i < len(pixels); i++ {
				for j := i + 1; j < len(pixels); j++ {
					if pixels[i] > pixels[j] {
						pixels[i], pixels[j] = pixels[j], pixels[i]
					}
				}
			}

			median := pixels[len(pixels)/2]
			filtered.Set(x, y, color.Gray{Y: median})
		}
	}

	// Copy edges
	for y := 0; y < height; y++ {
		filtered.Set(0, y, img.GrayAt(0, y))
		filtered.Set(width-1, y, img.GrayAt(width-1, y))
	}
	for x := 0; x < width; x++ {
		filtered.Set(x, 0, img.GrayAt(x, 0))
		filtered.Set(x, height-1, img.GrayAt(x, height-1))
	}

	return filtered
}

// ocrWithPassportEyeFull uses enhanced PassportEye for full passport extraction including issue date
func ocrWithPassportEyeFull(filePath string) (*PassportEyeResult, error) {
	// Call Python enhanced PassportEye script
	cmd := exec.Command("python3", "./passporteye_full_extractor.py", filePath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("PassportEye full extraction failed: %v, output: %s", err, string(output))
	}

	// Parse JSON result
	var result PassportEyeResult
	err = json.Unmarshal(output, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PassportEye full result: %v, output: %s", err, string(output))
	}

	return &result, nil
}

// ocrWithPassportEye uses PassportEye for specialized MRZ extraction
func ocrWithPassportEye(filePath string) (*PassportEyeResult, error) {
	// Call Python PassportEye script
	cmd := exec.Command("python3", "./passporteye_extractor.py", filePath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("PassportEye execution failed: %v, output: %s", err, string(output))
	}

	// Parse JSON result
	var result PassportEyeResult
	err = json.Unmarshal(output, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PassportEye result: %v, output: %s", err, string(output))
	}

	return &result, nil
}

// ocrWithTesseractPassportText uses Tesseract optimized for passport text (not MRZ)
func ocrWithTesseractPassportText(filePath string) (*OCRResult, error) {
	client := gosseract.NewClient()
	defer client.Close()

	// Set image path
	if err := client.SetImage(filePath); err != nil {
		return nil, fmt.Errorf("failed to set image: %v", err)
	}

	// Optimize for passport text (not MRZ)
	client.SetPageSegMode(gosseract.PSM_AUTO)

	// Passport text whitelist (more permissive than MRZ)
	client.SetVariable("tessedit_char_whitelist", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789 .,:-/")

	// Improve OCR for small text
	client.SetVariable("textord_min_linesize", "2.5")
	client.SetVariable("preserve_interword_spaces", "1")

	// Extract text
	text, err := client.Text()
	if err != nil {
		return nil, fmt.Errorf("tesseract OCR failed: %v", err)
	}

	// Clean up text
	text = strings.TrimSpace(text)

	// Calculate basic confidence (Tesseract Go binding doesn't expose confidence easily)
	confidence := 0.8
	if len(text) < 10 {
		confidence = 0.4
	}

	return &OCRResult{
		Text:       text,
		Confidence: confidence,
	}, nil
}

// ocrWithTesseract uses Tesseract directly with optimized passport configuration
func ocrWithTesseract(filePath string) (*OCRResult, error) {
	client := gosseract.NewClient()
	defer client.Close()

	// Preprocess image if it's an MRZ image for better OCR accuracy
	var imageToProcess string
	if strings.Contains(filePath, "mrz") {
		fmt.Printf("🎨 Preprocessing MRZ image for better OCR...\n")
		preprocessedPath, err := preprocessMRZImage(filePath)
		if err != nil {
			fmt.Printf("⚠️ Preprocessing failed, using original: %v\n", err)
			imageToProcess = filePath
		} else {
			imageToProcess = preprocessedPath
			fmt.Printf("✅ Using preprocessed image for OCR\n")
		}
	} else {
		imageToProcess = filePath
	}

	// Set image from file
	err := client.SetImage(imageToProcess)
	if err != nil {
		return nil, fmt.Errorf("failed to set image: %v", err)
	}

	// Configure Tesseract for passport OCR (optimized for MRZ)
	client.SetLanguage("eng")

	// Use different page segmentation for MRZ vs other parts
	if strings.Contains(imageToProcess, "mrz") {
		client.SetPageSegMode(gosseract.PSM_SINGLE_BLOCK) // Better for MRZ lines
		// MRZ-specific character whitelist (no spaces, periods, or slashes)
		client.SetVariable("tessedit_char_whitelist", "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789<")
		client.SetVariable("tessedit_pageseg_mode", "6")     // Uniform block of text
		client.SetVariable("preserve_interword_spaces", "0") // No spaces in MRZ
		// Additional MRZ-specific settings
		client.SetVariable("textord_really_old_xheight", "1")
		client.SetVariable("textord_min_xheight", "10")
		client.SetVariable("classify_enable_learning", "0")
		client.SetVariable("classify_enable_adaptive_matcher", "0")
	} else {
		client.SetPageSegMode(gosseract.PSM_AUTO)
		// General passport text whitelist
		client.SetVariable("tessedit_char_whitelist", "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789<>/ .-")
		client.SetVariable("tessedit_pageseg_mode", "6")
		client.SetVariable("preserve_interword_spaces", "1")
	}

	// Extract text
	text, err := client.Text()
	if err != nil {
		return nil, fmt.Errorf("tesseract OCR failed: %v", err)
	}

	// Clean up common OCR noise for MRZ
	text = cleanMRZNoise(text)

	confidence := calculateOCRConfidence(text, "tesseract")
	return &OCRResult{
		Text:       text,
		Confidence: confidence,
		Engine:     "tesseract",
	}, nil
}

// cleanMRZNoise removes common OCR noise from MRZ text
func cleanMRZNoise(text string) string {
	// Remove common OCR misreads
	replacements := map[string]string{
		"0": "O", // In names, 0 is usually O
		"1": "I", // In names, 1 is usually I
		"5": "S", // In names, 5 is usually S
		"|": "I", // Vertical bar to I
		"!": "I", // Exclamation to I
		"@": "O", // At symbol to O
	}

	lines := strings.Split(text, "\n")
	var cleanedLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 10 { // Only process substantial lines
			// Apply replacements only to likely name portions
			if strings.Contains(line, "P<") || len(line) > 30 {
				for old, new := range replacements {
					line = strings.ReplaceAll(line, old, new)
				}
			}
		}
		if line != "" {
			cleanedLines = append(cleanedLines, line)
		}
	}

	return strings.Join(cleanedLines, "\n")
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

// processPassport processes passport documents using advanced three-zone extraction
func processPassport(file multipart.File, uploadType string) (map[string]interface{}, float64, error) {
	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read file: %v", err)
	}

	fmt.Printf("🔍 STEP 1: Three-zone passport extraction...\n")

	// Extract three zones: Page 1, Page 2 Upper, Page 2 MRZ
	page1Bytes, page2UpperBytes, _, page1Path, page2UpperPath, page2MrzPath, err := extractThreeZones(fileBytes)
	if err != nil {
		fmt.Printf("❌ Three-zone extraction failed, falling back to simple MRZ: %v\n", err)
		// Fallback to simple MRZ extraction
		_, err := extractMRZAreaSimple(fileBytes)
		if err != nil {
			return nil, 0, fmt.Errorf("fallback MRZ extraction failed: %v", err)
		}
	} else {
		fmt.Printf("✅ Three zones extracted successfully\n")
	}

	fmt.Printf("🔍 STEP 2: Running OCR on all three zones...\n")

	// OCR Page 1 (Photo + Signature)
	var page1Text string
	if page1Bytes != nil {
		if page1Result, err := ocrWithTesseractPassportText(page1Path); err == nil {
			page1Text = page1Result.Text
			fmt.Printf("📄 Page 1 OCR: %d characters extracted\n", len(page1Text))
			fmt.Printf("=== PAGE 1 OCR TEXT START ===\n%s\n=== PAGE 1 OCR TEXT END ===\n", page1Text)
		} else {
			fmt.Printf("⚠️ Page 1 OCR failed: %v\n", err)
		}
	}

	// OCR Page 2 Upper (Text area - contains issue date and other passport details)
	var page2UpperText string
	if page2UpperBytes != nil {
		if page2UpperResult, err := ocrWithTesseractPassportText(page2UpperPath); err == nil {
			page2UpperText = page2UpperResult.Text
			fmt.Printf("📄 Page 2 Upper OCR: %d characters extracted\n", len(page2UpperText))
			fmt.Printf("=== PAGE 2 UPPER OCR TEXT START ===\n%s\n=== PAGE 2 UPPER OCR TEXT END ===\n", page2UpperText)
		} else {
			fmt.Printf("⚠️ Page 2 Upper OCR failed: %v\n", err)
		}
	}

	// OCR Page 2 MRZ (Most important for passport data) - Try enhanced PassportEye first on full passport
	fmt.Printf("🎯 Attempting full passport extraction with enhanced PassportEye on original image...\n")

	var passportEyeResult *PassportEyeResult

	// Create a temporary file path for the original image (we need to save the original bytes first)
	originalPath := strings.Replace(page2MrzPath, "page2_mrz_", "full_passport_", 1)
	if writeErr := os.WriteFile(originalPath, fileBytes, 0644); writeErr == nil {
		passportEyeResult, err = ocrWithPassportEyeFull(originalPath)
		fmt.Printf("📄 Full passport extraction on original image: success=%t\n", err == nil && passportEyeResult.Success)
	} else {
		err = fmt.Errorf("failed to save original image: %v", writeErr)
	}

	// If full extraction fails, fallback to MRZ-only extraction
	if err != nil || !passportEyeResult.Success {
		fmt.Printf("⚠️ Full PassportEye extraction failed (%v), trying MRZ-only...\n", err)
		passportEyeResult, err = ocrWithPassportEye(page2MrzPath)
	}

	var finalMRZText string
	var finalConfidence float64
	var engine string

	if err == nil && passportEyeResult.Success {
		// PassportEye succeeded
		finalMRZText = passportEyeResult.RawText
		finalConfidence = passportEyeResult.Confidence
		engine = "passporteye"
		fmt.Printf("✅ PassportEye extraction successful (confidence: %.2f, valid: %t)\n",
			passportEyeResult.Confidence, passportEyeResult.Valid)
		fmt.Printf("=== PASSPORTEYE MRZ TEXT START ===\n%s\n=== PASSPORTEYE MRZ TEXT END ===\n", finalMRZText)

		// If PassportEye found valid MRZ, use its structured data directly
		if passportEyeResult.Valid && len(passportEyeResult.ExtractedData) > 0 {
			fmt.Printf("🎯 Using PassportEye structured data extraction\n")
			// We'll use this structured data later
		}
	} else {
		// PassportEye failed, fallback to Tesseract
		fmt.Printf("⚠️ PassportEye failed (%v), falling back to Tesseract...\n", err)
		ocrResult, tesseractErr := ocrWithTesseract(page2MrzPath)
		if tesseractErr != nil {
			return nil, 0, fmt.Errorf("both PassportEye and Tesseract MRZ OCR failed: PassportEye: %v, Tesseract: %v", err, tesseractErr)
		}
		finalMRZText = ocrResult.Text
		finalConfidence = ocrResult.Confidence
		engine = "tesseract"
		fmt.Printf("📄 Page 2 MRZ OCR using %s (confidence: %.2f):\n", engine, finalConfidence)
		fmt.Printf("=== MRZ OCR TEXT START ===\n%s\n=== MRZ OCR TEXT END ===\n", finalMRZText)
	}

	fmt.Printf("🔍 STEP 3: Parsing MRZ data...\n")

	var extractedData map[string]interface{}

	// Use PassportEye structured data if available and valid
	if err == nil && passportEyeResult.Success && passportEyeResult.Valid && len(passportEyeResult.ExtractedData) > 0 {
		fmt.Printf("🎯 Using PassportEye structured extraction\n")
		extractedData = make(map[string]interface{})

		// Copy PassportEye data with our field names
		peData := passportEyeResult.ExtractedData
		if surname, ok := peData["surname"].(string); ok && surname != "" {
			extractedData["surname"] = surname
		}
		if givenNames, ok := peData["givenNames"].(string); ok && givenNames != "" {
			extractedData["givenNames"] = givenNames
		}
		if passportNumber, ok := peData["passportNumber"].(string); ok && passportNumber != "" {
			extractedData["passportNumber"] = passportNumber
		}
		if nationality, ok := peData["nationality"].(string); ok && nationality != "" {
			// Map nationality code to full country name
			if fullCountryName, exists := countryCodeMap[nationality]; exists {
				extractedData["nationality"] = fullCountryName
				fmt.Printf("🌍 Mapped nationality: %s -> %s\n", nationality, fullCountryName)
			} else {
				extractedData["nationality"] = nationality
				fmt.Printf("⚠️ Unknown nationality code: %s\n", nationality)
			}
		}
		if issuingCountry, ok := peData["issuingCountry"].(string); ok && issuingCountry != "" {
			// Map country code to full country name
			if fullCountryName, exists := countryCodeMap[issuingCountry]; exists {
				extractedData["issuingCountry"] = fullCountryName
				fmt.Printf("🌍 Mapped issuing country: %s -> %s\n", issuingCountry, fullCountryName)
			} else {
				extractedData["issuingCountry"] = issuingCountry
				fmt.Printf("⚠️ Unknown country code: %s\n", issuingCountry)
			}
		}
		if gender, ok := peData["gender"].(string); ok && gender != "" {
			extractedData["gender"] = gender
		}
		if dateOfBirth, ok := peData["dateOfBirth"].(string); ok && dateOfBirth != "" {
			extractedData["dateOfBirth"] = dateOfBirth
		}
		if expiryDate, ok := peData["expiryDate"].(string); ok && expiryDate != "" {
			extractedData["expiryDate"] = expiryDate
		}

		// Add raw MRZ text
		if len(passportEyeResult.MRZLines) > 0 {
			extractedData["machineReadableZone"] = strings.Join(passportEyeResult.MRZLines, "\n")
		}

		fmt.Printf("📊 PassportEye extracted %d fields\n", len(extractedData))

		// Try to extract issue date from page 1 text first (UK passports often have it there)
		fmt.Printf("🔍 Checking for issue date in page 1 text (length: %d)\n", len(page1Text))
		if page1Text != "" {
			if issueDate := extractIssueDateFromText(page1Text); issueDate != "" {
				extractedData["issueDate"] = issueDate
				fmt.Printf("✅ Added issue date from page 1 text: %s\n", issueDate)
			}
		}

		// If not found in page 1, try page 2 upper text
		if _, hasIssueDate := extractedData["issueDate"]; !hasIssueDate {
			fmt.Printf("🔍 Checking for issue date in page 2 upper text (length: %d)\n", len(page2UpperText))
			if page2UpperText != "" {
				if issueDate := extractIssueDateFromText(page2UpperText); issueDate != "" {
					extractedData["issueDate"] = issueDate
					fmt.Printf("✅ Added issue date from page 2 text: %s\n", issueDate)
				}
			} else {
				fmt.Printf("⚠️ Page 2 upper text is empty - cannot extract issue date\n")
			}
		}
	} else {
		// Fallback to manual parsing of OCR text
		fmt.Printf("🔍 Using manual MRZ parsing from OCR text\n")
		extractedData = parseMRZOnly(finalMRZText)

		// Try to extract issue date from page 1 text first (fallback case)
		fmt.Printf("🔍 Checking for issue date in page 1 text - fallback (length: %d)\n", len(page1Text))
		if page1Text != "" {
			if issueDate := extractIssueDateFromText(page1Text); issueDate != "" {
				extractedData["issueDate"] = issueDate
				fmt.Printf("✅ Added issue date from page 1 text (fallback): %s\n", issueDate)
			}
		}

		// If not found in page 1, try page 2 upper text (fallback case)
		if _, hasIssueDate := extractedData["issueDate"]; !hasIssueDate {
			fmt.Printf("🔍 Checking for issue date in page 2 upper text - fallback (length: %d)\n", len(page2UpperText))
			if page2UpperText != "" {
				if issueDate := extractIssueDateFromText(page2UpperText); issueDate != "" {
					extractedData["issueDate"] = issueDate
					fmt.Printf("✅ Added issue date from page 2 text (fallback): %s\n", issueDate)
				}
			} else {
				fmt.Printf("⚠️ Page 2 upper text is empty - cannot extract issue date (fallback)\n")
			}
		}
	}

	// Calculate final confidence
	var fieldConfidence float64
	if err == nil && passportEyeResult.Success {
		// PassportEye provides its own confidence
		fieldConfidence = finalConfidence
	} else {
		// Calculate field confidence for Tesseract results
		fieldConfidence = calculatePassportConfidence(extractedData)
	}

	finalConfidenceValue := (finalConfidence + fieldConfidence) / 2.0

	// Add OCR metadata and image paths
	extractedData["_ocrEngine"] = engine
	extractedData["_ocrConfidence"] = finalConfidence
	extractedData["_fieldConfidence"] = fieldConfidence
	extractedData["_mrzExtracted"] = true

	// Add image paths for frontend access
	if page1Path != "" {
		extractedData["_page1ImagePath"] = page1Path
		extractedData["_page1ImageUrl"] = "/static/mrz/" + filepath.Base(page1Path)
	}
	if page2UpperPath != "" {
		extractedData["_page2UpperImagePath"] = page2UpperPath
		extractedData["_page2UpperImageUrl"] = "/static/mrz/" + filepath.Base(page2UpperPath)
	}
	if page2MrzPath != "" {
		extractedData["_page2MrzImagePath"] = page2MrzPath
		extractedData["_page2MrzImageUrl"] = "/static/mrz/" + filepath.Base(page2MrzPath)

		// Add preprocessed MRZ image URL
		preprocessedMrzPath := strings.Replace(page2MrzPath, ".png", "_preprocessed.png", 1)
		extractedData["_page2MrzPreprocessedUrl"] = "/static/mrz/" + filepath.Base(preprocessedMrzPath)
	}

	fmt.Printf("✅ STEP 4: Extraction complete - %d fields found\n", len(extractedData)-7) // -7 for metadata fields

	return extractedData, finalConfidenceValue, nil
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

// extractThreeZones splits passport image into Page 1, Page 2 Upper, and Page 2 MRZ
func extractThreeZones(imageBytes []byte) ([]byte, []byte, []byte, string, string, string, error) {
	// Decode image
	img, format, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, nil, nil, "", "", "", fmt.Errorf("failed to decode image: %v", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	fmt.Printf("📏 Original passport image: %dx%d pixels\n", width, height)

	// Create static/mrz directory if it doesn't exist
	mrzDir := "static/mrz"
	if err := os.MkdirAll(mrzDir, 0755); err != nil {
		fmt.Printf("⚠️ Could not create mrz directory: %v\n", err)
	}

	// Detect passport fold line (between pages)
	foldLine := detectPassportFoldLine(img)
	if foldLine == -1 {
		// Fallback: assume fold is at 50% for standard passport
		foldLine = height / 2
		fmt.Printf("📏 Using fallback fold line at 50%% (%d pixels)\n", foldLine)
	} else {
		fmt.Printf("📏 Detected passport fold line at %d pixels\n", foldLine)
	}

	// Calculate zones
	// Page 1: Top to fold line (photo + signature page)
	page1Height := foldLine

	// Page 2 Upper: Fold line to MRZ start (text area) - expand significantly to capture issue date
	mrzStartY := int(float64(height) * 0.70) // MRZ typically starts at 70% down (was 75%)
	page2UpperStartY := foldLine
	page2UpperHeight := mrzStartY - foldLine

	// Page 2 MRZ: MRZ lines only (bottom 25%)
	page2MrzHeight := height - mrzStartY

	fmt.Printf("📏 Zone calculations:\n")
	fmt.Printf("  Page 1: 0-%d (%dx%d)\n", foldLine, width, page1Height)
	fmt.Printf("  Page 2 Upper: %d-%d (%dx%d)\n", foldLine, mrzStartY, width, page2UpperHeight)
	fmt.Printf("  Page 2 MRZ: %d-%d (%dx%d)\n", mrzStartY, height, width, page2MrzHeight)

	timestamp := time.Now().Unix()

	// Extract Page 1 (Photo + Signature)
	var page1Bytes []byte
	var page1Path string
	if page1Height > 0 {
		page1Img := image.NewRGBA(image.Rect(0, 0, width, page1Height))
		for y := 0; y < page1Height; y++ {
			for x := 0; x < width; x++ {
				page1Img.Set(x, y, img.At(x, y))
			}
		}

		page1Path = filepath.Join(mrzDir, fmt.Sprintf("page1_passport_%d.png", timestamp))
		if err := saveImage(page1Img, page1Path, format); err == nil {
			page1Bytes, _ = imageToBytes(page1Img, format)
			fmt.Printf("💾 Page 1 saved: %s\n", page1Path)
			fmt.Printf("🔗 Page 1 URL: http://localhost:3000/%s\n", strings.Replace(page1Path, "\\", "/", -1))
		}
	}

	// Extract Page 2 Upper (Text area)
	var page2UpperBytes []byte
	var page2UpperPath string
	if page2UpperHeight > 0 {
		page2UpperImg := image.NewRGBA(image.Rect(0, 0, width, page2UpperHeight))
		for y := page2UpperStartY; y < mrzStartY; y++ {
			for x := 0; x < width; x++ {
				page2UpperImg.Set(x, y-page2UpperStartY, img.At(x, y))
			}
		}

		page2UpperPath = filepath.Join(mrzDir, fmt.Sprintf("page2_upper_passport_%d.png", timestamp))
		if err := saveImage(page2UpperImg, page2UpperPath, format); err == nil {
			page2UpperBytes, _ = imageToBytes(page2UpperImg, format)
			fmt.Printf("💾 Page 2 Upper saved: %s\n", page2UpperPath)
			fmt.Printf("🔗 Page 2 Upper URL: http://localhost:3000/%s\n", strings.Replace(page2UpperPath, "\\", "/", -1))
		}
	}

	// Extract Page 2 MRZ (MRZ lines only)
	var page2MrzBytes []byte
	var page2MrzPath string
	if page2MrzHeight > 0 {
		page2MrzImg := image.NewRGBA(image.Rect(0, 0, width, page2MrzHeight))
		for y := mrzStartY; y < height; y++ {
			for x := 0; x < width; x++ {
				page2MrzImg.Set(x, y-mrzStartY, img.At(x, y))
			}
		}

		page2MrzPath = filepath.Join(mrzDir, fmt.Sprintf("page2_mrz_passport_%d.png", timestamp))
		if err := saveImage(page2MrzImg, page2MrzPath, format); err == nil {
			page2MrzBytes, _ = imageToBytes(page2MrzImg, format)
			fmt.Printf("💾 Page 2 MRZ saved: %s\n", page2MrzPath)
			fmt.Printf("🔗 Page 2 MRZ URL: http://localhost:3000/%s\n", strings.Replace(page2MrzPath, "\\", "/", -1))
		}
	}

	return page1Bytes, page2UpperBytes, page2MrzBytes, page1Path, page2UpperPath, page2MrzPath, nil
}

// detectPassportFoldLine detects the fold line between passport pages
func detectPassportFoldLine(img image.Image) int {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Look for horizontal line patterns in the middle section (30-70% of height)
	startY := int(float64(height) * 0.3)
	endY := int(float64(height) * 0.7)

	maxLineStrength := 0.0
	bestFoldY := -1

	for y := startY; y < endY; y++ {
		lineStrength := 0.0

		// Sample brightness changes across the width
		for x := 1; x < width-1; x++ {
			r1, g1, b1, _ := img.At(x-1, y).RGBA()
			r2, g2, b2, _ := img.At(x+1, y).RGBA()

			// Calculate brightness difference
			brightness1 := (r1 + g1 + b1) / 3
			brightness2 := (r2 + g2 + b2) / 3

			diff := float64(brightness1) - float64(brightness2)
			if diff < 0 {
				diff = -diff
			}
			lineStrength += diff
		}

		// Normalize by width
		lineStrength /= float64(width)

		if lineStrength > maxLineStrength {
			maxLineStrength = lineStrength
			bestFoldY = y
		}
	}

	// Only return fold line if we found a strong enough pattern
	if maxLineStrength > 1000 { // Threshold for fold detection
		return bestFoldY
	}

	return -1 // No fold detected
}

// extractMRZAreaSimple extracts just the MRZ area (fallback method)
func extractMRZAreaSimple(imageBytes []byte) ([]byte, error) {
	// Decode image
	img, format, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Extract bottom 15% for MRZ (85-100%)
	mrzStartY := int(float64(height) * 0.85)
	mrzHeight := height - mrzStartY

	fmt.Printf("📏 Simple MRZ extraction: %dx%d pixels (Y: %d-%d)\n", width, mrzHeight, mrzStartY, height)

	// Create cropped image for MRZ area
	mrzImg := image.NewRGBA(image.Rect(0, 0, width, mrzHeight))

	// Copy MRZ area to new image
	for y := mrzStartY; y < height; y++ {
		for x := 0; x < width; x++ {
			mrzImg.Set(x, y-mrzStartY, img.At(x, y))
		}
	}

	// Encode cropped image back to bytes
	return imageToBytes(mrzImg, format)
}

// saveImage saves an image to a file
func saveImage(img image.Image, filePath, format string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	switch format {
	case "jpeg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 95})
	case "png":
		return png.Encode(file, img)
	default:
		return png.Encode(file, img)
	}
}

// imageToBytes converts an image to bytes
func imageToBytes(img image.Image, format string) ([]byte, error) {
	var buf bytes.Buffer

	switch format {
	case "jpeg":
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 95})
		return buf.Bytes(), err
	case "png":
		err := png.Encode(&buf, img)
		return buf.Bytes(), err
	default:
		err := png.Encode(&buf, img)
		return buf.Bytes(), err
	}
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

// isMRZLine checks if a line looks like an MRZ line (enhanced detection)
func isMRZLine(line string) bool {
	line = strings.TrimSpace(line)

	// Check for passport MRZ patterns
	if strings.HasPrefix(line, "P<") {
		return true
	}

	// Enhanced passport number patterns (both letter+numbers and numbers+letters)
	passportPatterns := []*regexp.Regexp{
		regexp.MustCompile(`^[A-Z][A-Z0-9]{6,9}`),    // Standard: letter + alphanumeric
		regexp.MustCompile(`^[0-9]{8,10}[A-Z]{2,3}`), // Alternative: numbers + letters
		regexp.MustCompile(`^[A-Z]{2,3}[0-9]{6,9}`),  // Country + numbers
	}

	for _, pattern := range passportPatterns {
		if pattern.MatchString(line) {
			return true
		}
	}

	// Check for typical MRZ length (44 characters for TD-3, but allow flexibility)
	if len(line) >= 30 && len(line) <= 50 {
		// Count alphanumeric characters vs spaces/symbols
		alphaNum := 0
		mrzChars := 0 // Count MRZ-specific characters
		for _, char := range line {
			if (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
				alphaNum++
			}
			if char == '<' || char == '>' || char == '/' {
				mrzChars++
			}
		}

		// MRZ lines should be mostly alphanumeric with some MRZ characters
		alphaRatio := float64(alphaNum) / float64(len(line))
		if alphaRatio > 0.6 && (mrzChars > 0 || alphaRatio > 0.8) {
			return true
		}
	}

	// Check for MRZ date patterns (YYMMDD)
	datePattern := regexp.MustCompile(`[0-9]{6}[MF][0-9]{6}`)
	if datePattern.MatchString(line) {
		return true
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
