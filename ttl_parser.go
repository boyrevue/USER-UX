package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type OntologyField struct {
	Property               string        `json:"property"`
	Label                  string        `json:"label"`
	Type                   string        `json:"type"`
	Required               bool          `json:"required"`
	HelpText               string        `json:"helpText,omitempty"`
	Options                []FieldOption `json:"options,omitempty"`
	Domain                 string        `json:"domain,omitempty"`
	ConditionalDisplay     string        `json:"conditionalDisplay,omitempty"`
	ConditionalRequirement string        `json:"conditionalRequirement,omitempty"`
	IsMultiSelect          bool          `json:"isMultiSelect"`
	FormType               string        `json:"formType"`
	EnumerationValues      []string      `json:"enumerationValues"`
	ArrayItemStructure     string        `json:"arrayItemStructure,omitempty"`
	FormSection            string        `json:"formSection,omitempty"`
	FormInfoText           string        `json:"formInfoText,omitempty"`
	DefaultValue           string        `json:"defaultValue,omitempty"`
	RequiresAIValidation   bool          `json:"requiresAIValidation"`
	AIValidationPrompt     string        `json:"aiValidationPrompt,omitempty"`
	MinInclusive           *int          `json:"minInclusive,omitempty"`
	MaxInclusive           *int          `json:"maxInclusive,omitempty"`
}

type FieldOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type OntologySection struct {
	ID     string          `json:"id"`
	Label  string          `json:"label"`
	Fields []OntologyField `json:"fields"`
}

// ParseTTLOntology - parses all TTL ontology files (insurance, app config, documents)
func ParseTTLOntology() (map[string]OntologySection, error) {
	// Read the modular auto insurance ontologies
	driverData, err := ioutil.ReadFile("ontology/AI_Driver_Details.ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to read AI_Driver_Details.ttl: %v", err)
	}

	vehicleData, err := ioutil.ReadFile("ontology/AI_Vehicle_Details.ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to read AI_Vehicle_Details.ttl: %v", err)
	}

	policyData, err := ioutil.ReadFile("ontology/AI_Policy_Details.ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to read AI_Policy_Details.ttl: %v", err)
	}

	claimsData, err := ioutil.ReadFile("ontology/AI_Claims_History.ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to read AI_Claims_History.ttl: %v", err)
	}

	paymentsData, err := ioutil.ReadFile("ontology/AI_Insurance_Payments.ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to read AI_Insurance_Payments.ttl: %v", err)
	}

	// Read the GDPR compliance ontology
	complianceData, err := ioutil.ReadFile("ontology/AI_Data_Compliance.ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to read AI_Data_Compliance.ttl: %v", err)
	}

	// Read the USER-UX app configuration TTL file
	userUxData, err := ioutil.ReadFile("ontology/user_ux.ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to read user_ux.ttl: %v", err)
	}

	// Read the user documents TTL file
	documentsData, err := ioutil.ReadFile("ontology/user_documents.ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to read user_documents.ttl: %v", err)
	}

	// Read the comprehensive personal documents ontology
	personalDocsData, err := ioutil.ReadFile("ontology/personal_documents_ontology.ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to read personal_documents_ontology.ttl: %v", err)
	}

	// Combine all TTL files for parsing
	content := string(driverData) + "\n" + string(vehicleData) + "\n" + string(policyData) + "\n" + string(claimsData) + "\n" + string(paymentsData) + "\n" + string(complianceData) + "\n" + string(userUxData) + "\n" + string(documentsData) + "\n" + string(personalDocsData)

	// Extract properties by parsing TTL patterns with deduplication maps
	driverFieldsMap := make(map[string]OntologyField)
	vehicleFieldsMap := make(map[string]OntologyField)
	claimsFieldsMap := make(map[string]OntologyField)
	settingsFieldsMap := make(map[string]OntologyField)
	documentsFieldsMap := make(map[string]OntologyField)

	// Regex patterns for TTL parsing (handles autoins:, settings:, docs:, banking:, comms:, : prefixes)
	propertyPattern := regexp.MustCompile(`(autoins|settings|banking|comms|docs|):(\w+)\s+a\s+owl:DatatypeProperty\s*;`)
	labelPattern := regexp.MustCompile(`rdfs:label\s+"([^"]+)"(?:@\w+)?\s*;`)
	domainPattern := regexp.MustCompile(`rdfs:domain\s+(autoins|settings|docs|foaf|):(\w+)\s*;`)
	rangePattern := regexp.MustCompile(`rdfs:range\s+xsd:(\w+)\s*;`)
	requiredPattern := regexp.MustCompile(`(autoins|docs|):isRequired\s+"(true|false)"\^\^xsd:boolean\s*;`)
	helpTextPattern := regexp.MustCompile(`(autoins|docs|):formHelpText\s+"([^"]+)"\s*;`)
	enumPattern := regexp.MustCompile(`(autoins|docs|):enumerationValues\s+\(([^)]+)\)\s*;`)
	conditionalDisplayPattern := regexp.MustCompile(`(autoins|docs|):conditionalDisplay\s+"([^"]+)"\s*;`)
	conditionalRequirementPattern := regexp.MustCompile(`(autoins|docs|):conditionalRequirement\s+"([^"]+)"\s*;`)
	isMultiSelectPattern := regexp.MustCompile(`(autoins|docs|):isMultiSelect\s+"(true|false)"\^\^xsd:boolean\s*;`)
	formTypePattern := regexp.MustCompile(`(autoins|docs|):formType\s+"([^"]+)"\s*;`)
	arrayItemStructurePattern := regexp.MustCompile(`(autoins|docs|):arrayItemStructure\s+"([^"]+)"\s*;`)
	formSectionPattern := regexp.MustCompile(`(autoins|docs|):formSection\s+"([^"]+)"\s*[;.]`)
	formInfoTextPattern := regexp.MustCompile(`(autoins|docs|):formInfoText\s+"([^"]+)"\s*[;.]`)
	defaultValuePattern := regexp.MustCompile(`(autoins|docs|):defaultValue\s+"([^"]+)"\s*[;.\s]`)
	requiresAIValidationPattern := regexp.MustCompile(`(autoins|docs|):requiresAIValidation\s+"(true|false)"\^\^xsd:boolean\s*[;.]`)
	aiValidationPromptPattern := regexp.MustCompile(`(autoins|docs|):aiValidationPrompt\s+"([^"]+)"\s*[;.]`)
	minInclusivePattern := regexp.MustCompile(`(autoins|docs|):minInclusive\s+(\d+)\s*[;.]`)
	maxInclusivePattern := regexp.MustCompile(`(autoins|docs|):maxInclusive\s+(\d+)\s*[;.]`)

	// Find all properties
	propertyMatches := propertyPattern.FindAllStringSubmatch(content, -1)

	for _, match := range propertyMatches {
		if len(match) < 3 {
			continue
		}

		// propPrefix := match[1]  // autoins, settings, banking, comms (not used currently)
		propName := match[2]

		// Find the property block (from property declaration to next property or end)
		propStart := strings.Index(content, match[0])
		if propStart == -1 {
			continue
		}

		// Find the end of this property block (ends with a period)
		var propBlock string
		blockEnd := strings.Index(content[propStart:], " .")
		if blockEnd == -1 {
			// Try to find next property as fallback
			nextPropIndex := strings.Index(content[propStart+1:], "autoins:")
			if nextPropIndex == -1 {
				propBlock = content[propStart:]
			} else {
				propBlock = content[propStart : propStart+1+nextPropIndex]
			}
		} else {
			propBlock = content[propStart : propStart+blockEnd+2]
		}

		// Extract information from this property block
		var label, domain, range_, helpText, conditionalDisplay, conditionalRequirement, formType string
		var required, isMultiSelect bool
		var enumValues []string

		if labelMatch := labelPattern.FindStringSubmatch(propBlock); len(labelMatch) > 1 {
			label = labelMatch[1]
		}

		if domainMatch := domainPattern.FindStringSubmatch(propBlock); len(domainMatch) > 2 {
			domain = domainMatch[2] // Extract the class name (Driver, Vehicle, etc.)
		}

		if rangeMatch := rangePattern.FindStringSubmatch(propBlock); len(rangeMatch) > 1 {
			range_ = rangeMatch[1]
		}

		if requiredMatch := requiredPattern.FindStringSubmatch(propBlock); len(requiredMatch) > 2 {
			required = requiredMatch[2] == "true"
		}

		if helpMatch := helpTextPattern.FindStringSubmatch(propBlock); len(helpMatch) > 2 {
			helpText = helpMatch[2]
		}

		if enumMatch := enumPattern.FindStringSubmatch(propBlock); len(enumMatch) > 2 {
			enumStr := enumMatch[2]
			// Extract quoted values from enumeration
			quotedValuePattern := regexp.MustCompile(`"([^"]+)"`)
			quotedMatches := quotedValuePattern.FindAllStringSubmatch(enumStr, -1)
			for _, match := range quotedMatches {
				if len(match) > 1 {
					enumValues = append(enumValues, match[1])
				}
			}
		}

		if conditionalDisplayMatch := conditionalDisplayPattern.FindStringSubmatch(propBlock); len(conditionalDisplayMatch) > 2 {
			conditionalDisplay = conditionalDisplayMatch[2]
		}

		if conditionalRequirementMatch := conditionalRequirementPattern.FindStringSubmatch(propBlock); len(conditionalRequirementMatch) > 2 {
			conditionalRequirement = conditionalRequirementMatch[2]
		}

		if isMultiSelectMatch := isMultiSelectPattern.FindStringSubmatch(propBlock); len(isMultiSelectMatch) > 2 {
			isMultiSelect = isMultiSelectMatch[2] == "true"
		}

		if formTypeMatch := formTypePattern.FindStringSubmatch(propBlock); len(formTypeMatch) > 2 {
			formType = formTypeMatch[2]
		}

		var arrayItemStructure string
		if arrayItemStructureMatch := arrayItemStructurePattern.FindStringSubmatch(propBlock); len(arrayItemStructureMatch) > 2 {
			arrayItemStructure = arrayItemStructureMatch[2]
		}

		var formSection string
		if formSectionMatch := formSectionPattern.FindStringSubmatch(propBlock); len(formSectionMatch) > 2 {
			formSection = formSectionMatch[2]
		}

		var formInfoText string
		if formInfoTextMatch := formInfoTextPattern.FindStringSubmatch(propBlock); len(formInfoTextMatch) > 2 {
			formInfoText = formInfoTextMatch[2]
		}

		var defaultValue string
		if defaultValueMatch := defaultValuePattern.FindStringSubmatch(propBlock); len(defaultValueMatch) > 2 {
			defaultValue = defaultValueMatch[2]
			if propName == "isMainDriver" || propName == "manualOrAuto" {
				fmt.Printf("DEBUG: Found defaultValue for %s: '%s'\n", propName, defaultValue)
			}
		} else if propName == "isMainDriver" || propName == "manualOrAuto" {
			fmt.Printf("DEBUG: No defaultValue match found for %s\n", propName)
			startPos := len(propBlock) - 300
			if startPos < 0 {
				startPos = 0
			}
			fmt.Printf("DEBUG: Property block snippet: %s\n", propBlock[startPos:])
		}

		var requiresAIValidation bool
		if requiresAIValidationMatch := requiresAIValidationPattern.FindStringSubmatch(propBlock); len(requiresAIValidationMatch) > 2 {
			requiresAIValidation = requiresAIValidationMatch[2] == "true"
		}

		var aiValidationPrompt string
		if aiValidationPromptMatch := aiValidationPromptPattern.FindStringSubmatch(propBlock); len(aiValidationPromptMatch) > 2 {
			aiValidationPrompt = aiValidationPromptMatch[2]
		}

		// Extract minInclusive and maxInclusive values
		var minInclusive, maxInclusive *int
		if minInclusiveMatch := minInclusivePattern.FindStringSubmatch(propBlock); len(minInclusiveMatch) > 2 {
			if val, err := strconv.Atoi(minInclusiveMatch[2]); err == nil {
				minInclusive = &val
			}
		}
		if maxInclusiveMatch := maxInclusivePattern.FindStringSubmatch(propBlock); len(maxInclusiveMatch) > 2 {
			if val, err := strconv.Atoi(maxInclusiveMatch[2]); err == nil {
				maxInclusive = &val
			}
		}

		// Debug logging for accident fields
		if propName == "accidentDate" || propName == "accidentFault" {
			fmt.Printf("DEBUG: Parsing %s - formSection: '%s'\n", propName, formSection)
			fmt.Printf("DEBUG: Property block snippet: %s\n", propBlock[len(propBlock)-200:])
		}

		// Skip if no domain or label
		if domain == "" || label == "" {
			fmt.Printf("DEBUG: Skipping field %s - domain: '%s', label: '%s'\n", propName, domain, label)
			continue
		}

		fmt.Printf("DEBUG: Processing field %s - domain: '%s', label: '%s'\n", propName, domain, label)

		// Determine field type
		fieldType := "text"
		switch range_ {
		case "boolean":
			fieldType = "radio"
		case "date":
			fieldType = "date"
		default:
			if strings.Contains(propName, "email") {
				fieldType = "email"
			} else if strings.Contains(propName, "phone") {
				fieldType = "tel"
			}
		}

		// Handle enumeration values
		var options []FieldOption
		if len(enumValues) > 0 {
			for _, value := range enumValues {
				if strings.TrimSpace(value) != "" {
					// Clean up the value and create a nice label
					cleanValue := strings.TrimSpace(value)
					// Convert underscores to spaces and capitalize first letter
					displayLabel := strings.ReplaceAll(cleanValue, "_", " ")
					displayLabel = strings.Title(strings.ToLower(displayLabel))

					options = append(options, FieldOption{
						Value: cleanValue,
						Label: displayLabel,
					})
				}
			}
			if len(options) > 0 {
				if isMultiSelect {
					fieldType = "checkbox"
				} else if len(options) <= 3 {
					fieldType = "radio"
				} else {
					fieldType = "select"
				}
			}
		}

		// Create field
		field := OntologyField{
			Property:               propName,
			Label:                  label,
			Type:                   fieldType,
			Required:               required,
			HelpText:               helpText,
			Options:                options,
			Domain:                 domain,
			ConditionalDisplay:     conditionalDisplay,
			ConditionalRequirement: conditionalRequirement,
			IsMultiSelect:          isMultiSelect,
			FormType:               formType,
			EnumerationValues:      enumValues,
			ArrayItemStructure:     arrayItemStructure,
			FormSection:            formSection,
			FormInfoText:           formInfoText,
			DefaultValue:           defaultValue,
			RequiresAIValidation:   requiresAIValidation,
			AIValidationPrompt:     aiValidationPrompt,
			MinInclusive:           minInclusive,
			MaxInclusive:           maxInclusive,
		}

		// Debug output for specific fields
		if propName == "disabilityTypes" || propName == "automaticOnly" || propName == "adaptationTypes" {
			fmt.Printf("DEBUG: Creating field %s - IsMultiSelect: %v, FormType: '%s', EnumValues: %v\n", propName, isMultiSelect, formType, enumValues)
			fmt.Printf("DEBUG: Field struct: %+v\n", field)
		}

		// Add to appropriate section based on domain (with deduplication)
		switch domain {
		case "Driver", "Conviction":
			// Always add/overwrite (keep the latest definition)
			driverFieldsMap[field.Property] = field
		case "Vehicle":
			// Only add if not already present (deduplication by property name)
			if _, exists := vehicleFieldsMap[field.Property]; !exists {
				vehicleFieldsMap[field.Property] = field
			}
		case "Claims", "ClaimsHistory":
			// Only add if not already present (deduplication by property name)
			if _, exists := claimsFieldsMap[field.Property]; !exists {
				claimsFieldsMap[field.Property] = field
			}
		case "Settings", "BankAccount", "CreditCard", "InsurancePayments", "CommunicationChannel":
			// Only add if not already present (deduplication by property name)
			if _, exists := settingsFieldsMap[field.Property]; !exists {
				settingsFieldsMap[field.Property] = field
			}
		case "PassportDocument", "DrivingLicenceDocument", "IdentityCardDocument", "UtilityBillDocument", "BankStatementDocument", "MedicalCertificateDocument", "InsuranceDocument", "PersonalDocument", "IdentityDocument", "FinancialDocument", "EducationalDocument", "EmploymentDocument", "PropertyDocument", "MedicalDocument", "TravelDocument", "VehicleDocument":
			// Only add if not already present (deduplication by property name)
			if _, exists := documentsFieldsMap[field.Property]; !exists {
				documentsFieldsMap[field.Property] = field
			}
		// Comprehensive Insurance Entity Classes
		case "InsuranceEntity", "MotorInsuranceDocument", "DriverDocument", "ClaimsDocument", "PolicyDocument":
			// Assign to appropriate section based on document type (with deduplication)
			if strings.Contains(strings.ToLower(domain), "driver") {
				if _, exists := driverFieldsMap[field.Property]; !exists {
					driverFieldsMap[field.Property] = field
				}
			} else if strings.Contains(strings.ToLower(domain), "vehicle") {
				if _, exists := vehicleFieldsMap[field.Property]; !exists {
					vehicleFieldsMap[field.Property] = field
				}
			} else if strings.Contains(strings.ToLower(domain), "claim") {
				if _, exists := claimsFieldsMap[field.Property]; !exists {
					claimsFieldsMap[field.Property] = field
				}
			} else {
				// Default insurance documents to documents section
				if _, exists := documentsFieldsMap[field.Property]; !exists {
					documentsFieldsMap[field.Property] = field
				}
			}
		// Specific Insurance Document Sub-classes
		case "InsuranceCertificate", "InsuranceSchedule", "ProofOfNoClaims", "QuoteDocument", "RenewalNotice":
			if _, exists := documentsFieldsMap[field.Property]; !exists {
				documentsFieldsMap[field.Property] = field
			}
		case "DrivingLicence", "ConvictionCertificate", "MedicalCertificate", "PassPlusCertificate":
			if _, exists := driverFieldsMap[field.Property]; !exists {
				driverFieldsMap[field.Property] = field
			}
		case "VehicleRegistrationDocument", "MOTCertificate", "VehicleValuation", "ModificationCertificate", "SecurityDeviceCertificate":
			if _, exists := vehicleFieldsMap[field.Property]; !exists {
				vehicleFieldsMap[field.Property] = field
			}
		case "ClaimForm", "AccidentReport", "PoliceReport", "RepairEstimate", "SettlementLetter":
			if _, exists := claimsFieldsMap[field.Property]; !exists {
				claimsFieldsMap[field.Property] = field
			}
		case "PolicyWording", "EndorsementDocument", "CancellationNotice":
			if _, exists := documentsFieldsMap[field.Property]; !exists {
				documentsFieldsMap[field.Property] = field
			}
		case "PremiumInvoice", "PaymentReceipt", "RefundNotice", "DirectDebitMandate":
			if _, exists := documentsFieldsMap[field.Property]; !exists {
				documentsFieldsMap[field.Property] = field
			}
		// FOAF classes (Person, Organization)
		case "Person", "Organization":
			// These are typically used as ranges for object properties, not domains for datatype properties
			// But if they appear as domains, add to settings
			if _, exists := settingsFieldsMap[field.Property]; !exists {
				settingsFieldsMap[field.Property] = field
			}
		}
	}

	// Convert maps to slices for final output
	driverFields := make([]OntologyField, 0, len(driverFieldsMap))
	for _, field := range driverFieldsMap {
		driverFields = append(driverFields, field)
	}

	vehicleFields := make([]OntologyField, 0, len(vehicleFieldsMap))
	for _, field := range vehicleFieldsMap {
		vehicleFields = append(vehicleFields, field)
	}

	claimsFields := make([]OntologyField, 0, len(claimsFieldsMap))
	for _, field := range claimsFieldsMap {
		claimsFields = append(claimsFields, field)
	}

	settingsFields := make([]OntologyField, 0, len(settingsFieldsMap))
	for _, field := range settingsFieldsMap {
		settingsFields = append(settingsFields, field)
	}

	documentsFields := make([]OntologyField, 0, len(documentsFieldsMap))
	for _, field := range documentsFieldsMap {
		documentsFields = append(documentsFields, field)
	}

	// Build sections
	sections := map[string]OntologySection{
		"drivers": {
			ID:     "drivers",
			Label:  "Driver Details",
			Fields: driverFields,
		},
		"vehicles": {
			ID:     "vehicles",
			Label:  "Vehicle Details",
			Fields: vehicleFields,
		},
		"claims": {
			ID:     "claims",
			Label:  "Claims History",
			Fields: claimsFields,
		},
		"settings": {
			ID:     "settings",
			Label:  "Application Settings",
			Fields: settingsFields,
		},
		"documents": {
			ID:     "documents",
			Label:  "Personal Documents",
			Fields: documentsFields,
		},
	}

	return sections, nil
}

func HandleOntologyAPI(w http.ResponseWriter, r *http.Request) {
	sections, err := ParseTTLOntology()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse ontology: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sections)
}
