package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// SHACL Transformation Engine
type SHACLTransformer struct {
	Rules      []TransformationRule `json:"rules"`
	Ontologies map[string]Ontology  `json:"ontologies"`
	Namespaces map[string]string    `json:"namespaces"`
}

type TransformationRule struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	SourceType  string         `json:"sourceType"`
	TargetType  string         `json:"targetType"`
	Pattern     string         `json:"pattern"`
	Transform   string         `json:"transform"`
	Validation  ValidationRule `json:"validation"`
	Priority    int            `json:"priority"`
}

type ValidationRule struct {
	Required      bool     `json:"required"`
	MinLength     int      `json:"minLength"`
	MaxLength     int      `json:"maxLength"`
	Pattern       string   `json:"pattern"`
	AllowedValues []string `json:"allowedValues"`
	DataType      string   `json:"dataType"`
}

type Ontology struct {
	Classes    map[string]OntologyClass    `json:"classes"`
	Properties map[string]OntologyProperty `json:"properties"`
	Shapes     map[string]SHACLShape       `json:"shapes"`
}

type OntologyClass struct {
	URI        string            `json:"uri"`
	Label      map[string]string `json:"label"` // language -> label
	Comment    map[string]string `json:"comment"`
	SubClassOf []string          `json:"subClassOf"`
	Properties []string          `json:"properties"`
}

type OntologyProperty struct {
	URI      string            `json:"uri"`
	Label    map[string]string `json:"label"`
	Comment  map[string]string `json:"comment"`
	Domain   []string          `json:"domain"`
	Range    []string          `json:"range"`
	DataType string            `json:"dataType"`
}

type SHACLShape struct {
	URI         string               `json:"uri"`
	TargetClass string               `json:"targetClass"`
	Properties  []SHACLPropertyShape `json:"properties"`
	Closed      bool                 `json:"closed"`
}

type SHACLPropertyShape struct {
	Path          string   `json:"path"`
	DataType      string   `json:"dataType"`
	MinCount      int      `json:"minCount"`
	MaxCount      int      `json:"maxCount"`
	Pattern       string   `json:"pattern"`
	MinLength     int      `json:"minLength"`
	MaxLength     int      `json:"maxLength"`
	AllowedValues []string `json:"allowedValues"`
	Message       string   `json:"message"`
}

type TransformationRequest struct {
	SourceData  map[string]interface{} `json:"sourceData"`
	TargetShape string                 `json:"targetShape"`
	Mappings    []FieldMapping         `json:"mappings"`
	Options     TransformOptions       `json:"options"`
}

type TransformOptions struct {
	StrictValidation bool                 `json:"strictValidation"`
	Language         string               `json:"language"`
	CustomRules      []TransformationRule `json:"customRules"`
	PreserveOriginal bool                 `json:"preserveOriginal"`
}

type TransformationResult struct {
	Success          bool                   `json:"success"`
	TransformedData  map[string]interface{} `json:"transformedData"`
	ValidationErrors []SHACLValidationError `json:"validationErrors"`
	Warnings         []string               `json:"warnings"`
	AppliedRules     []string               `json:"appliedRules"`
	ProcessedAt      time.Time              `json:"processedAt"`
}

type SHACLValidationError struct {
	Field    string `json:"field"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
	Rule     string `json:"rule"`
	Value    string `json:"value"`
}

// SHACL Transform Handler
func SHACLTransformHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request TransformationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Initialize transformer
	transformer, err := NewSHACLTransformer()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to initialize transformer: %v", err), http.StatusInternalServerError)
		return
	}

	// Perform transformation
	result, err := transformer.Transform(request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Transformation failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Initialize SHACL Transformer
func NewSHACLTransformer() (*SHACLTransformer, error) {
	transformer := &SHACLTransformer{
		Rules:      []TransformationRule{},
		Ontologies: make(map[string]Ontology),
		Namespaces: map[string]string{
			"autoins":  "http://example.org/autoins#",
			"settings": "http://example.org/settings#",
			"xsd":      "http://www.w3.org/2001/XMLSchema#",
			"sh":       "http://www.w3.org/ns/shacl#",
		},
	}

	// Load ontologies and rules
	err := transformer.LoadOntologies()
	if err != nil {
		return nil, err
	}

	err = transformer.LoadTransformationRules()
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

// Load ontologies from TTL files
func (st *SHACLTransformer) LoadOntologies() error {
	// Load autoins.ttl
	autoinsOntology, err := st.parseTTLFile("ontology/autoins.ttl")
	if err != nil {
		return fmt.Errorf("failed to load autoins ontology: %v", err)
	}
	st.Ontologies["autoins"] = autoinsOntology

	// Load settings.ttl
	settingsOntology, err := st.parseTTLFile("ontology/settings.ttl")
	if err != nil {
		return fmt.Errorf("failed to load settings ontology: %v", err)
	}
	st.Ontologies["settings"] = settingsOntology

	return nil
}

// Load transformation rules
func (st *SHACLTransformer) LoadTransformationRules() error {
	// Define common transformation rules
	st.Rules = []TransformationRule{
		{
			ID:          "date_format_uk_to_iso",
			Name:        "UK Date to ISO Format",
			Description: "Transform UK date format (DD/MM/YYYY) to ISO format (YYYY-MM-DD)",
			SourceType:  "string",
			TargetType:  "xsd:date",
			Pattern:     `^(\d{2})/(\d{2})/(\d{4})$`,
			Transform:   "$3-$2-$1",
			Priority:    10,
			Validation: ValidationRule{
				Required: true,
				Pattern:  `^\d{4}-\d{2}-\d{2}$`,
				DataType: "xsd:date",
			},
		},
		{
			ID:          "phone_format_uk",
			Name:        "UK Phone Number Formatting",
			Description: "Standardize UK phone numbers to +44 format",
			SourceType:  "string",
			TargetType:  "xsd:string",
			Pattern:     `^(?:0|\+44\s?)?(\d{4})\s?(\d{3})\s?(\d{3})$`,
			Transform:   "+44 $1 $2 $3",
			Priority:    5,
			Validation: ValidationRule{
				Required: false,
				Pattern:  `^\+44\s\d{4}\s\d{3}\s\d{3}$`,
				DataType: "xsd:string",
			},
		},
		{
			ID:          "postcode_format_uk",
			Name:        "UK Postcode Formatting",
			Description: "Standardize UK postcodes to uppercase with proper spacing",
			SourceType:  "string",
			TargetType:  "xsd:string",
			Pattern:     `^([A-Za-z]{1,2}\d{1,2}[A-Za-z]?)\s*(\d[A-Za-z]{2})$`,
			Transform:   "UPPER($1 $2)",
			Priority:    3,
			Validation: ValidationRule{
				Required: false,
				Pattern:  `^[A-Z]{1,2}\d{1,2}[A-Z]?\s\d[A-Z]{2}$`,
				DataType: "xsd:string",
			},
		},
		{
			ID:          "licence_number_format",
			Name:        "Driving Licence Number Formatting",
			Description: "Format UK driving licence numbers",
			SourceType:  "string",
			TargetType:  "xsd:string",
			Pattern:     `^([A-Za-z]+)(\d+)([A-Za-z]+\d+[A-Za-z]+)$`,
			Transform:   "UPPER($1$2$3)",
			Priority:    8,
			Validation: ValidationRule{
				Required: true,
				Pattern:  `^[A-Z]{5}\d{6}[A-Z]{2}\d{2}[A-Z]$`,
				DataType: "xsd:string",
			},
		},
		{
			ID:          "name_capitalization",
			Name:        "Name Proper Capitalization",
			Description: "Ensure proper capitalization for names",
			SourceType:  "string",
			TargetType:  "xsd:string",
			Pattern:     `^(.+)$`,
			Transform:   "TITLE_CASE($1)",
			Priority:    1,
			Validation: ValidationRule{
				Required:  true,
				MinLength: 1,
				MaxLength: 50,
				DataType:  "xsd:string",
			},
		},
	}

	return nil
}

// Transform data according to SHACL rules
func (st *SHACLTransformer) Transform(request TransformationRequest) (*TransformationResult, error) {
	result := &TransformationResult{
		Success:          true,
		TransformedData:  make(map[string]interface{}),
		ValidationErrors: []SHACLValidationError{},
		Warnings:         []string{},
		AppliedRules:     []string{},
		ProcessedAt:      time.Now(),
	}

	// Get target shape definition
	targetShape, exists := st.getShape(request.TargetShape)
	if !exists {
		return nil, fmt.Errorf("target shape not found: %s", request.TargetShape)
	}

	// Process each field mapping
	for _, mapping := range request.Mappings {
		sourceValue, exists := request.SourceData[mapping.SourceField.Name]
		if !exists && mapping.SourceField.Required {
			result.ValidationErrors = append(result.ValidationErrors, SHACLValidationError{
				Field:    mapping.SourceField.Name,
				Message:  fmt.Sprintf("Required field '%s' is missing", mapping.SourceField.Name),
				Severity: "error",
				Rule:     "required_field",
			})
			result.Success = false
			continue
		}

		if !exists {
			continue // Skip optional missing fields
		}

		// Apply transformations
		transformedValue, err := st.applyTransformations(sourceValue, mapping, request.Options)
		if err != nil {
			result.ValidationErrors = append(result.ValidationErrors, SHACLValidationError{
				Field:    mapping.SourceField.Name,
				Message:  fmt.Sprintf("Transformation failed: %v", err),
				Severity: "error",
				Rule:     "transformation_error",
				Value:    fmt.Sprintf("%v", sourceValue),
			})
			result.Success = false
			continue
		}

		// Validate against SHACL shape
		validationErrors := st.validateAgainstShape(mapping.OntologyProp, transformedValue, targetShape)
		if len(validationErrors) > 0 {
			result.ValidationErrors = append(result.ValidationErrors, validationErrors...)
			if request.Options.StrictValidation {
				result.Success = false
				continue
			}
		}

		// Store transformed value
		result.TransformedData[mapping.OntologyProp] = transformedValue
	}

	return result, nil
}

// Apply transformations to a value
func (st *SHACLTransformer) applyTransformations(value interface{}, mapping FieldMapping, options TransformOptions) (interface{}, error) {
	stringValue := fmt.Sprintf("%v", value)

	// Find applicable transformation rules
	applicableRules := st.findApplicableRules(mapping.SourceField.Type, mapping.OntologyProp)

	// Add custom rules
	applicableRules = append(applicableRules, options.CustomRules...)

	// Sort by priority
	for i := 0; i < len(applicableRules)-1; i++ {
		for j := i + 1; j < len(applicableRules); j++ {
			if applicableRules[i].Priority < applicableRules[j].Priority {
				applicableRules[i], applicableRules[j] = applicableRules[j], applicableRules[i]
			}
		}
	}

	transformedValue := stringValue

	// Apply transformations in priority order
	for _, rule := range applicableRules {
		if rule.Pattern != "" {
			matched, err := regexp.MatchString(rule.Pattern, transformedValue)
			if err != nil {
				continue
			}
			if matched {
				transformedValue = st.applyTransformationRule(transformedValue, rule)
			}
		}
	}

	// Convert to appropriate data type
	return st.convertToTargetType(transformedValue, mapping.OntologyProp)
}

// Apply a specific transformation rule
func (st *SHACLTransformer) applyTransformationRule(value string, rule TransformationRule) string {
	if rule.Pattern == "" {
		return value
	}

	re := regexp.MustCompile(rule.Pattern)

	// Handle special transformation functions
	transform := rule.Transform
	if strings.Contains(transform, "UPPER(") {
		// Extract content within UPPER()
		upperRe := regexp.MustCompile(`UPPER\(([^)]+)\)`)
		transform = upperRe.ReplaceAllStringFunc(transform, func(match string) string {
			content := upperRe.FindStringSubmatch(match)[1]
			expanded := re.ReplaceAllString(value, content)
			return strings.ToUpper(expanded)
		})
		return transform
	}

	if strings.Contains(transform, "TITLE_CASE(") {
		titleRe := regexp.MustCompile(`TITLE_CASE\(([^)]+)\)`)
		transform = titleRe.ReplaceAllStringFunc(transform, func(match string) string {
			content := titleRe.FindStringSubmatch(match)[1]
			expanded := re.ReplaceAllString(value, content)
			return strings.Title(strings.ToLower(expanded))
		})
		return transform
	}

	// Standard regex replacement
	return re.ReplaceAllString(value, transform)
}

// Find applicable transformation rules
func (st *SHACLTransformer) findApplicableRules(sourceType, targetProperty string) []TransformationRule {
	var applicable []TransformationRule

	for _, rule := range st.Rules {
		// Check if rule applies to this source type or target property
		if rule.SourceType == sourceType || strings.Contains(targetProperty, rule.Name) {
			applicable = append(applicable, rule)
		}
	}

	return applicable
}

// Convert value to target data type
func (st *SHACLTransformer) convertToTargetType(value string, targetProperty string) (interface{}, error) {
	// Determine target type from ontology
	targetType := st.getPropertyDataType(targetProperty)

	switch targetType {
	case "xsd:integer":
		// Convert to integer if possible
		if value == "" {
			return nil, nil
		}
		// Implementation for integer conversion
		return value, nil
	case "xsd:date":
		// Validate date format
		if matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, value); matched {
			return value, nil
		}
		return nil, fmt.Errorf("invalid date format: %s", value)
	case "xsd:boolean":
		// Convert to boolean
		lower := strings.ToLower(value)
		if lower == "true" || lower == "yes" || lower == "1" {
			return true, nil
		}
		if lower == "false" || lower == "no" || lower == "0" {
			return false, nil
		}
		return nil, fmt.Errorf("invalid boolean value: %s", value)
	default:
		return value, nil
	}
}

// Validate value against SHACL shape
func (st *SHACLTransformer) validateAgainstShape(property string, value interface{}, shape SHACLShape) []SHACLValidationError {
	var errors []SHACLValidationError

	// Find property shape
	var propShape *SHACLPropertyShape
	for _, ps := range shape.Properties {
		if ps.Path == property {
			propShape = &ps
			break
		}
	}

	if propShape == nil {
		return errors // No validation rules for this property
	}

	stringValue := fmt.Sprintf("%v", value)

	// Check required (minCount)
	if propShape.MinCount > 0 && (value == nil || stringValue == "") {
		errors = append(errors, SHACLValidationError{
			Field:    property,
			Message:  propShape.Message,
			Severity: "error",
			Rule:     "minCount",
			Value:    stringValue,
		})
	}

	// Check pattern
	if propShape.Pattern != "" && stringValue != "" {
		matched, err := regexp.MatchString(propShape.Pattern, stringValue)
		if err == nil && !matched {
			errors = append(errors, SHACLValidationError{
				Field:    property,
				Message:  fmt.Sprintf("Value does not match required pattern: %s", propShape.Pattern),
				Severity: "error",
				Rule:     "pattern",
				Value:    stringValue,
			})
		}
	}

	// Check length constraints
	if propShape.MinLength > 0 && len(stringValue) < propShape.MinLength {
		errors = append(errors, SHACLValidationError{
			Field:    property,
			Message:  fmt.Sprintf("Value too short (minimum %d characters)", propShape.MinLength),
			Severity: "error",
			Rule:     "minLength",
			Value:    stringValue,
		})
	}

	if propShape.MaxLength > 0 && len(stringValue) > propShape.MaxLength {
		errors = append(errors, SHACLValidationError{
			Field:    property,
			Message:  fmt.Sprintf("Value too long (maximum %d characters)", propShape.MaxLength),
			Severity: "error",
			Rule:     "maxLength",
			Value:    stringValue,
		})
	}

	return errors
}

// Helper functions
func (st *SHACLTransformer) parseTTLFile(filename string) (Ontology, error) {
	// Mock implementation - would parse actual TTL files
	ontology := Ontology{
		Classes:    make(map[string]OntologyClass),
		Properties: make(map[string]OntologyProperty),
		Shapes:     make(map[string]SHACLShape),
	}

	// Add some mock definitions
	ontology.Classes["autoins:Person"] = OntologyClass{
		URI: "http://example.org/autoins#Person",
		Label: map[string]string{
			"en": "Person",
			"de": "Person",
		},
		Properties: []string{"autoins:firstName", "autoins:lastName"},
	}

	return ontology, nil
}

func (st *SHACLTransformer) getShape(shapeName string) (SHACLShape, bool) {
	for _, ontology := range st.Ontologies {
		if shape, exists := ontology.Shapes[shapeName]; exists {
			return shape, true
		}
	}
	return SHACLShape{}, false
}

func (st *SHACLTransformer) getPropertyDataType(property string) string {
	for _, ontology := range st.Ontologies {
		if prop, exists := ontology.Properties[property]; exists {
			return prop.DataType
		}
	}
	return "xsd:string" // Default
}
