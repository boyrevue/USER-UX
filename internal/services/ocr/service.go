package ocr

import (
	"encoding/json"
	"fmt"
	"io"
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

// ProcessResult represents the final OCR processing result
type ProcessResult struct {
	DocumentType    string                 `json:"documentType"`
	UploadType      string                 `json:"uploadType"`
	ExtractedFields map[string]interface{} `json:"extractedFields"`
	Confidence      float64                `json:"confidence"`
	ProcessedAt     string                 `json:"processedAt"`
	Images          map[string]string      `json:"images,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

type Service struct {
	staticDir string
}

func NewService() *Service {
	return &Service{
		staticDir: "static",
	}
}

func (s *Service) ProcessUpload(r *http.Request) (*ProcessResult, error) {
	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		return nil, fmt.Errorf("failed to parse multipart form: %v", err)
	}

	uploadType := r.FormValue("uploadType")
	selectedDocumentType := r.FormValue("selectedDocumentType")

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		return nil, fmt.Errorf("no files uploaded")
	}

	file := files[0]

	// Save uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Create unique filename
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	filename := fmt.Sprintf("%s_%s_%s", uploadType, selectedDocumentType, timestamp)

	// Ensure static directory exists
	os.MkdirAll(filepath.Join(s.staticDir, "mrz"), 0755)

	// Save file
	filePath := filepath.Join(s.staticDir, "mrz", filename+".png")
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	// Process based on document type
	switch selectedDocumentType {
	case "passport":
		return s.processPassport(filePath, uploadType, timestamp)
	case "driving_licence":
		return s.processDrivingLicence(filePath, uploadType, timestamp)
	default:
		return s.processGenericDocument(filePath, uploadType, selectedDocumentType, timestamp)
	}
}

func (s *Service) processPassport(filePath, uploadType, timestamp string) (*ProcessResult, error) {
	// Use PassportEye for MRZ extraction
	result, err := s.extractWithPassportEye(filePath, timestamp)
	if err != nil {
		return nil, fmt.Errorf("passport processing failed: %v", err)
	}

	return &ProcessResult{
		DocumentType:    "passport",
		UploadType:      uploadType,
		ExtractedFields: result.ExtractedData,
		Confidence:      result.Confidence,
		ProcessedAt:     time.Now().Format(time.RFC3339),
		Images: map[string]string{
			"original": filePath,
		},
		Metadata: map[string]interface{}{
			"engine":          "passporteye",
			"ocrConfidence":   result.Confidence,
			"fieldConfidence": result.Confidence,
			"mrzExtracted":    result.Success,
		},
	}, nil
}

func (s *Service) processDrivingLicence(filePath, uploadType, timestamp string) (*ProcessResult, error) {
	// Use Tesseract for driving licence OCR
	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(filePath)
	text, err := client.Text()
	if err != nil {
		return nil, fmt.Errorf("tesseract OCR failed: %v", err)
	}

	// Extract driving licence data from OCR text
	extractedData := s.extractDrivingLicenceData(text)

	return &ProcessResult{
		DocumentType:    "driving_licence",
		UploadType:      uploadType,
		ExtractedFields: extractedData,
		Confidence:      0.8, // Default confidence for Tesseract
		ProcessedAt:     time.Now().Format(time.RFC3339),
		Images: map[string]string{
			"original": filePath,
		},
		Metadata: map[string]interface{}{
			"engine":          "tesseract",
			"ocrConfidence":   0.8,
			"fieldConfidence": 0.8,
		},
	}, nil
}

func (s *Service) processGenericDocument(filePath, uploadType, documentType, timestamp string) (*ProcessResult, error) {
	// Use Tesseract for generic document OCR
	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(filePath)
	text, err := client.Text()
	if err != nil {
		return nil, fmt.Errorf("tesseract OCR failed: %v", err)
	}

	return &ProcessResult{
		DocumentType: documentType,
		UploadType:   uploadType,
		ExtractedFields: map[string]interface{}{
			"rawText": text,
		},
		Confidence:  0.7,
		ProcessedAt: time.Now().Format(time.RFC3339),
		Images: map[string]string{
			"original": filePath,
		},
		Metadata: map[string]interface{}{
			"engine":          "tesseract",
			"ocrConfidence":   0.7,
			"fieldConfidence": 0.7,
		},
	}, nil
}

func (s *Service) extractWithPassportEye(imagePath, timestamp string) (*PassportEyeResult, error) {
	// Call PassportEye Python script
	cmd := exec.Command("python3", "passporteye_extractor.py", imagePath, timestamp)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("passporteye execution failed: %v", err)
	}

	var result PassportEyeResult
	err = json.Unmarshal(output, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse passporteye result: %v", err)
	}

	return &result, nil
}

func (s *Service) extractDrivingLicenceData(text string) map[string]interface{} {
	data := make(map[string]interface{})

	// Extract licence number (pattern: ABCDE123456FG789)
	licencePattern := `[A-Z]{5}\d{6}[A-Z]{2}\d{3}`
	if matches := regexp.MustCompile(licencePattern).FindStringSubmatch(text); len(matches) > 0 {
		data["licenceNumber"] = matches[0]
	}

	// Extract dates (DD/MM/YYYY or DD.MM.YYYY)
	datePattern := `\b(\d{2})[\/\.](\d{2})[\/\.](\d{4})\b`
	dates := regexp.MustCompile(datePattern).FindAllStringSubmatch(text, -1)
	if len(dates) >= 2 {
		data["issueDate"] = fmt.Sprintf("%s-%s-%s", dates[0][3], dates[0][2], dates[0][1])
		data["expiryDate"] = fmt.Sprintf("%s-%s-%s", dates[1][3], dates[1][2], dates[1][1])
	}

	// Extract name (look for common patterns)
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 5 && strings.Contains(strings.ToUpper(line), "MR ") || strings.Contains(strings.ToUpper(line), "MRS ") || strings.Contains(strings.ToUpper(line), "MISS ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				data["firstName"] = parts[1]
				if len(parts) >= 3 {
					data["lastName"] = parts[2]
				}
			}
		}
	}

	data["rawText"] = text
	return data
}
