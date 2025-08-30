package bipro

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// GDVProcessor handles German insurance industry standard data format (GDV)
type GDVProcessor struct {
	version     string
	fieldSpecs  map[string]GDVFieldSpec
	recordTypes map[string]GDVRecordType
}

// NewGDVProcessor creates a new GDV format processor
func NewGDVProcessor() *GDVProcessor {
	processor := &GDVProcessor{
		version:     "GDV 2018",
		fieldSpecs:  initializeGDVFieldSpecs(),
		recordTypes: initializeGDVRecordTypes(),
	}
	return processor
}

// GDVRecord represents a single GDV data record
type GDVRecord struct {
	RecordType string            `json:"recordType"`
	Length     int               `json:"length"`
	Fields     map[string]string `json:"fields"`
	RawData    string            `json:"rawData"`
	LineNumber int               `json:"lineNumber"`
	Errors     []GDVError        `json:"errors,omitempty"`
}

// GDVFieldSpec defines the specification for a GDV field
type GDVFieldSpec struct {
	Position    int    `json:"position"`    // Starting position (1-based)
	Length      int    `json:"length"`      // Field length
	Type        string `json:"type"`        // A=Alpha, N=Numeric, D=Date
	Format      string `json:"format"`      // Date format, numeric format, etc.
	Description string `json:"description"` // Field description
	Required    bool   `json:"required"`    // Whether field is mandatory
}

// GDVRecordType defines the structure of a GDV record type
type GDVRecordType struct {
	Code        string                  `json:"code"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Length      int                     `json:"length"`
	Fields      map[string]GDVFieldSpec `json:"fields"`
}

// GDVError represents a GDV processing error
type GDVError struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Field    string `json:"field,omitempty"`
	Position int    `json:"position,omitempty"`
	Severity string `json:"severity"`
}

// GDVProcessingResult represents the result of GDV processing
type GDVProcessingResult struct {
	Records     []GDVRecord `json:"records"`
	Summary     GDVSummary  `json:"summary"`
	Errors      []GDVError  `json:"errors"`
	ProcessedAt time.Time   `json:"processedAt"`
}

// GDVSummary provides summary statistics for processed GDV data
type GDVSummary struct {
	TotalRecords     int            `json:"totalRecords"`
	RecordTypeCounts map[string]int `json:"recordTypeCounts"`
	ErrorCount       int            `json:"errorCount"`
	WarningCount     int            `json:"warningCount"`
	ProcessingTime   time.Duration  `json:"processingTime"`
}

// ProcessGDVData processes GDV format data from ZIP file or raw text
func (gdv *GDVProcessor) ProcessGDVData(data []byte, metadata BiPROMetadata) error {
	startTime := time.Now()

	// Determine if data is ZIP compressed
	var textData string
	var err error

	if isZipData(data) {
		textData, err = gdv.extractFromZip(data)
		if err != nil {
			return fmt.Errorf("failed to extract ZIP data: %w", err)
		}
	} else {
		textData = string(data)
	}

	// Process GDV records
	result, err := gdv.parseGDVText(textData)
	if err != nil {
		return fmt.Errorf("failed to parse GDV data: %w", err)
	}

	result.ProcessedAt = time.Now()
	result.Summary.ProcessingTime = time.Since(startTime)

	// Store processed data (implementation would save to database)
	return gdv.storeProcessedData(result, metadata)
}

// parseGDVText parses GDV text format into structured records
func (gdv *GDVProcessor) parseGDVText(text string) (*GDVProcessingResult, error) {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")

	result := &GDVProcessingResult{
		Records: make([]GDVRecord, 0),
		Summary: GDVSummary{
			RecordTypeCounts: make(map[string]int),
		},
		Errors: make([]GDVError, 0),
	}

	for lineNum, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue // Skip empty lines
		}

		record, err := gdv.parseGDVRecord(line, lineNum+1)
		if err != nil {
			result.Errors = append(result.Errors, GDVError{
				Code:     "RECORD_PARSE_ERROR",
				Message:  err.Error(),
				Position: lineNum + 1,
				Severity: "ERROR",
			})
			continue
		}

		result.Records = append(result.Records, *record)
		result.Summary.RecordTypeCounts[record.RecordType]++

		// Count errors and warnings
		for _, recordError := range record.Errors {
			if recordError.Severity == "ERROR" {
				result.Summary.ErrorCount++
			} else if recordError.Severity == "WARNING" {
				result.Summary.WarningCount++
			}
		}
	}

	result.Summary.TotalRecords = len(result.Records)

	return result, nil
}

// parseGDVRecord parses a single GDV record line
func (gdv *GDVProcessor) parseGDVRecord(line string, lineNum int) (*GDVRecord, error) {
	if len(line) < 4 {
		return nil, fmt.Errorf("record too short: minimum 4 characters required")
	}

	// Extract record type (first 4 characters)
	recordType := line[:4]

	recordSpec, exists := gdv.recordTypes[recordType]
	if !exists {
		return nil, fmt.Errorf("unknown record type: %s", recordType)
	}

	// Validate record length
	if len(line) != recordSpec.Length {
		return nil, fmt.Errorf("invalid record length: expected %d, got %d", recordSpec.Length, len(line))
	}

	record := &GDVRecord{
		RecordType: recordType,
		Length:     len(line),
		Fields:     make(map[string]string),
		RawData:    line,
		LineNumber: lineNum,
		Errors:     make([]GDVError, 0),
	}

	// Parse fields according to specification
	for fieldName, fieldSpec := range recordSpec.Fields {
		value, err := gdv.extractField(line, fieldSpec)
		if err != nil {
			record.Errors = append(record.Errors, GDVError{
				Code:     "FIELD_EXTRACTION_ERROR",
				Message:  err.Error(),
				Field:    fieldName,
				Position: fieldSpec.Position,
				Severity: "ERROR",
			})
			continue
		}

		// Validate field value
		if err := gdv.validateField(fieldName, value, fieldSpec); err != nil {
			record.Errors = append(record.Errors, GDVError{
				Code:     "FIELD_VALIDATION_ERROR",
				Message:  err.Error(),
				Field:    fieldName,
				Position: fieldSpec.Position,
				Severity: "WARNING",
			})
		}

		record.Fields[fieldName] = value
	}

	return record, nil
}

// extractField extracts a field value from a GDV record line
func (gdv *GDVProcessor) extractField(line string, spec GDVFieldSpec) (string, error) {
	startPos := spec.Position - 1 // Convert to 0-based indexing
	endPos := startPos + spec.Length

	if startPos < 0 || endPos > len(line) {
		return "", fmt.Errorf("field position out of bounds: %d-%d", startPos, endPos)
	}

	value := line[startPos:endPos]

	// Trim spaces for alpha fields
	if spec.Type == "A" {
		value = strings.TrimSpace(value)
	}

	return value, nil
}

// validateField validates a field value according to its specification
func (gdv *GDVProcessor) validateField(fieldName, value string, spec GDVFieldSpec) error {
	// Check required fields
	if spec.Required && strings.TrimSpace(value) == "" {
		return fmt.Errorf("required field %s is empty", fieldName)
	}

	// Type-specific validation
	switch spec.Type {
	case "N": // Numeric
		if strings.TrimSpace(value) != "" {
			if _, err := strconv.ParseFloat(strings.ReplaceAll(value, " ", ""), 64); err != nil {
				return fmt.Errorf("invalid numeric value: %s", value)
			}
		}
	case "D": // Date
		if strings.TrimSpace(value) != "" {
			if err := gdv.validateDateField(value, spec.Format); err != nil {
				return fmt.Errorf("invalid date value: %s", value)
			}
		}
	case "A": // Alpha
		// Alpha fields generally don't need special validation
		break
	}

	return nil
}

// validateDateField validates date fields according to GDV date formats
func (gdv *GDVProcessor) validateDateField(value, format string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil // Empty dates are handled by required field validation
	}

	var layout string
	switch format {
	case "DDMMYYYY":
		layout = "02012006"
	case "DDMMYY":
		layout = "020106"
	case "MMYYYY":
		layout = "012006"
	case "YYYY":
		layout = "2006"
	default:
		layout = "020106" // Default format
	}

	_, err := time.Parse(layout, value)
	return err
}

// extractFromZip extracts GDV data from ZIP file
func (gdv *GDVProcessor) extractFromZip(data []byte) (string, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}

	var content strings.Builder

	for _, file := range reader.File {
		if strings.HasSuffix(strings.ToLower(file.Name), ".gdv") ||
			strings.HasSuffix(strings.ToLower(file.Name), ".txt") {

			rc, err := file.Open()
			if err != nil {
				continue
			}

			fileContent, err := io.ReadAll(rc)
			rc.Close()

			if err != nil {
				continue
			}

			content.Write(fileContent)
			content.WriteString("\n")
		}
	}

	return content.String(), nil
}

// storeProcessedData stores processed GDV data (implementation would save to database)
func (gdv *GDVProcessor) storeProcessedData(result *GDVProcessingResult, metadata BiPROMetadata) error {
	// Implementation would store to database with proper indexing
	// For now, we'll just log the processing result

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize processing result: %w", err)
	}

	// Log processing result (in production, this would be stored in database)
	fmt.Printf("GDV Processing Result:\n%s\n", string(jsonData))

	return nil
}

// isZipData checks if the data is ZIP compressed
func isZipData(data []byte) bool {
	return len(data) >= 4 &&
		data[0] == 0x50 && data[1] == 0x4B &&
		(data[2] == 0x03 || data[2] == 0x05 || data[2] == 0x07) &&
		(data[3] == 0x04 || data[3] == 0x06 || data[3] == 0x08)
}

// ============================================================
// GDV FIELD SPECIFICATIONS
// ============================================================

// initializeGDVFieldSpecs initializes common GDV field specifications
func initializeGDVFieldSpecs() map[string]GDVFieldSpec {
	return map[string]GDVFieldSpec{
		"RECORD_TYPE": {
			Position:    1,
			Length:      4,
			Type:        "N",
			Description: "Record type identifier",
			Required:    true,
		},
		"POLICY_NUMBER": {
			Position:    5,
			Length:      17,
			Type:        "A",
			Description: "Policy number",
			Required:    true,
		},
		"CUSTOMER_NUMBER": {
			Position:    22,
			Length:      17,
			Type:        "A",
			Description: "Customer number",
			Required:    false,
		},
		"BIRTH_DATE": {
			Position:    39,
			Length:      8,
			Type:        "D",
			Format:      "DDMMYYYY",
			Description: "Date of birth",
			Required:    false,
		},
		"GENDER": {
			Position:    47,
			Length:      1,
			Type:        "A",
			Description: "Gender (M/F)",
			Required:    false,
		},
		"POSTAL_CODE": {
			Position:    48,
			Length:      5,
			Type:        "N",
			Description: "Postal code",
			Required:    false,
		},
	}
}

// initializeGDVRecordTypes initializes GDV record type specifications
func initializeGDVRecordTypes() map[string]GDVRecordType {
	return map[string]GDVRecordType{
		"0100": {
			Code:        "0100",
			Name:        "Address Record",
			Description: "Customer address information",
			Length:      256,
			Fields: map[string]GDVFieldSpec{
				"RECORD_TYPE":     {Position: 1, Length: 4, Type: "N", Required: true},
				"POLICY_NUMBER":   {Position: 5, Length: 17, Type: "A", Required: true},
				"CUSTOMER_NUMBER": {Position: 22, Length: 17, Type: "A", Required: false},
				"LAST_NAME":       {Position: 39, Length: 30, Type: "A", Required: true},
				"FIRST_NAME":      {Position: 69, Length: 30, Type: "A", Required: true},
				"BIRTH_DATE":      {Position: 99, Length: 8, Type: "D", Format: "DDMMYYYY", Required: false},
				"GENDER":          {Position: 107, Length: 1, Type: "A", Required: false},
				"STREET":          {Position: 108, Length: 30, Type: "A", Required: false},
				"HOUSE_NUMBER":    {Position: 138, Length: 7, Type: "A", Required: false},
				"POSTAL_CODE":     {Position: 145, Length: 5, Type: "N", Required: false},
				"CITY":            {Position: 150, Length: 25, Type: "A", Required: false},
				"COUNTRY":         {Position: 175, Length: 3, Type: "A", Required: false},
			},
		},
		"0200": {
			Code:        "0200",
			Name:        "Contract Record",
			Description: "Insurance contract information",
			Length:      256,
			Fields: map[string]GDVFieldSpec{
				"RECORD_TYPE":       {Position: 1, Length: 4, Type: "N", Required: true},
				"POLICY_NUMBER":     {Position: 5, Length: 17, Type: "A", Required: true},
				"CONTRACT_STATUS":   {Position: 22, Length: 1, Type: "N", Required: true},
				"START_DATE":        {Position: 23, Length: 8, Type: "D", Format: "DDMMYYYY", Required: true},
				"END_DATE":          {Position: 31, Length: 8, Type: "D", Format: "DDMMYYYY", Required: false},
				"PREMIUM_ANNUAL":    {Position: 39, Length: 12, Type: "N", Required: true},
				"PREMIUM_FREQUENCY": {Position: 51, Length: 1, Type: "N", Required: true},
				"COVERAGE_TYPE":     {Position: 52, Length: 3, Type: "A", Required: true},
				"EXCESS":            {Position: 55, Length: 8, Type: "N", Required: false},
			},
		},
		"0300": {
			Code:        "0300",
			Name:        "Vehicle Record",
			Description: "Vehicle information",
			Length:      256,
			Fields: map[string]GDVFieldSpec{
				"RECORD_TYPE":   {Position: 1, Length: 4, Type: "N", Required: true},
				"POLICY_NUMBER": {Position: 5, Length: 17, Type: "A", Required: true},
				"VEHICLE_ID":    {Position: 22, Length: 17, Type: "A", Required: true},
				"MAKE":          {Position: 39, Length: 20, Type: "A", Required: true},
				"MODEL":         {Position: 59, Length: 25, Type: "A", Required: true},
				"YEAR":          {Position: 84, Length: 4, Type: "N", Required: true},
				"VIN":           {Position: 88, Length: 17, Type: "A", Required: false},
				"REGISTRATION":  {Position: 105, Length: 12, Type: "A", Required: true},
				"ENGINE_SIZE":   {Position: 117, Length: 6, Type: "N", Required: false},
				"FUEL_TYPE":     {Position: 123, Length: 1, Type: "A", Required: false},
				"VEHICLE_VALUE": {Position: 124, Length: 10, Type: "N", Required: true},
			},
		},
		"0400": {
			Code:        "0400",
			Name:        "Claims Record",
			Description: "Claims information",
			Length:      256,
			Fields: map[string]GDVFieldSpec{
				"RECORD_TYPE":   {Position: 1, Length: 4, Type: "N", Required: true},
				"POLICY_NUMBER": {Position: 5, Length: 17, Type: "A", Required: true},
				"CLAIM_NUMBER":  {Position: 22, Length: 17, Type: "A", Required: true},
				"CLAIM_DATE":    {Position: 39, Length: 8, Type: "D", Format: "DDMMYYYY", Required: true},
				"INCIDENT_DATE": {Position: 47, Length: 8, Type: "D", Format: "DDMMYYYY", Required: true},
				"CLAIM_TYPE":    {Position: 55, Length: 3, Type: "A", Required: true},
				"CLAIM_AMOUNT":  {Position: 58, Length: 12, Type: "N", Required: true},
				"PAID_AMOUNT":   {Position: 70, Length: 12, Type: "N", Required: false},
				"FAULT_STATUS":  {Position: 82, Length: 1, Type: "N", Required: true},
				"CLAIM_STATUS":  {Position: 83, Length: 2, Type: "N", Required: true},
			},
		},
	}
}
