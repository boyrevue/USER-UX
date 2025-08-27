package main

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// DocumentData represents extracted document information
type DocumentData struct {
	DocumentType     string                    `json:"documentType"`
	UploadType       string                    `json:"uploadType"`
	ExtractedFields  map[string]interface{}    `json:"extractedFields"`
	Confidence       float64                   `json:"confidence"`
	ProcessedAt      time.Time                 `json:"processedAt"`
	ValidationResult *DocumentValidationResult `json:"validationResult,omitempty"`
}

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

	// Create response
	response := DocumentData{
		DocumentType:     documentType,
		UploadType:       uploadType,
		ExtractedFields:  extractedData,
		Confidence:       confidence,
		ProcessedAt:      time.Now(),
		ValidationResult: validationResult,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// processDrivingLicence processes UK driving licence documents
func processDrivingLicence(file multipart.File, uploadType string) (map[string]interface{}, float64, error) {
	// Read file content (for future OCR processing)
	_, err := io.ReadAll(file)
	if err != nil {
		return nil, 0, err
	}

	// Call OCR service (mock implementation for now)
	extractedData := make(map[string]interface{})

	if uploadType == "front" {
		// Mock extracted data for front side
		extractedData = map[string]interface{}{
			"licenceNumber":    "MORGA657054SM9IJ",
			"surname":          "MORGAN",
			"givenNames":       "SARAH MEREDYTH",
			"dateOfBirth":      "1976-03-11",
			"placeOfBirth":     "UNITED KINGDOM",
			"issueDate":        "2019-01-19",
			"expiryDate":       "2029-01-18",
			"issuingAuthority": "DVLA",
			"licenceType":      "FULL",
			"address":          "122 BURNS CRESCENT, EDINBURGH EH1 9GP",
		}
	} else if uploadType == "back" {
		// Mock extracted data for back side
		extractedData = map[string]interface{}{
			"categories":    []string{"A", "B", "BE"},
			"restrictions":  []string{},
			"endorsements":  []string{},
			"licenceNumber": "MORGA657054SM9IJ", // Should match front
		}
	}

	// In a real implementation, you would:
	// 1. Send to OCR service (Tesseract, AWS Textract, Google Vision, etc.)
	// 2. Apply document-specific parsing rules
	// 3. Validate extracted fields against expected patterns

	return extractedData, 0.95, nil
}

// processPassport processes passport documents
func processPassport(file multipart.File, uploadType string) (map[string]interface{}, float64, error) {
	// Read file content (for future OCR processing)
	_, err := io.ReadAll(file)
	if err != nil {
		return nil, 0, err
	}

	// Mock extracted passport data
	extractedData := map[string]interface{}{
		"passportNumber":      "502135326",
		"surname":             "MORGAN",
		"givenNames":          "SARAH MEREDYTH",
		"nationality":         "BRITISH CITIZEN",
		"dateOfBirth":         "1976-03-11",
		"placeOfBirth":        "CROYDON, UNITED KINGDOM",
		"gender":              "F",
		"issueDate":           "2019-02-05",
		"expiryDate":          "2029-02-04",
		"issuingCountry":      "GBR",
		"issuingAuthority":    "HM PASSPORT OFFICE",
		"machineReadableZone": "P<GBRMORGAN<<SARAH<MEREDYTH<<<<<<<<<<<<<<<<502135326<5GBR7603115F2902045<<<<<<<<<<<<<<04",
	}

	// In a real implementation:
	// 1. Parse MRZ (Machine Readable Zone) using specialized libraries
	// 2. Extract biographical data from visual inspection zone
	// 3. Validate check digits and format
	// 4. Cross-reference MRZ and VIZ data for consistency

	return extractedData, 0.92, nil
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
