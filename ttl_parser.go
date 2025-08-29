package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type OntologyField struct {
	Property string        `json:"property"`
	Label    string        `json:"label"`
	Type     string        `json:"type"`
	Required bool          `json:"required"`
	HelpText string        `json:"helpText,omitempty"`
	Options  []FieldOption `json:"options,omitempty"`
	Domain   string        `json:"domain,omitempty"`
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

	// Extract properties by parsing TTL patterns
	driverFields := []OntologyField{}
	vehicleFields := []OntologyField{}
	claimsFields := []OntologyField{}
	settingsFields := []OntologyField{}
	documentsFields := []OntologyField{}

	// Regex patterns for TTL parsing (handles autoins:, settings:, docs:, banking:, comms: prefixes)
	propertyPattern := regexp.MustCompile(`(autoins|settings|banking|comms|docs):(\w+)\s+a\s+owl:DatatypeProperty\s*;`)
	labelPattern := regexp.MustCompile(`rdfs:label\s+"([^"]+)"\s*;`)
	domainPattern := regexp.MustCompile(`rdfs:domain\s+(autoins|settings|docs|foaf):(\w+)\s*;`)
	rangePattern := regexp.MustCompile(`rdfs:range\s+xsd:(\w+)\s*;`)
	requiredPattern := regexp.MustCompile(`(autoins|docs):isRequired\s+"(true|false)"\^\^xsd:boolean\s*;`)
	helpTextPattern := regexp.MustCompile(`(autoins|docs):formHelpText\s+"([^"]+)"\s*;`)
	enumPattern := regexp.MustCompile(`(autoins|docs):enumerationValues\s+\(([^)]+)\)\s*;`)

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
		var label, domain, range_, helpText string
		var required bool
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
			enumStr := strings.ReplaceAll(enumMatch[2], `"`, "")
			enumValues = strings.Fields(enumStr)
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
					options = append(options, FieldOption{
						Value: strings.TrimSpace(value),
						Label: strings.TrimSpace(value),
					})
				}
			}
			if len(options) > 0 {
				if len(options) <= 3 {
					fieldType = "radio"
				} else {
					fieldType = "select"
				}
			}
		}

		// Create field
		field := OntologyField{
			Property: propName,
			Label:    label,
			Type:     fieldType,
			Required: required,
			HelpText: helpText,
			Options:  options,
			Domain:   domain,
		}

		// Add to appropriate section based on domain
		switch domain {
		case "Driver":
			driverFields = append(driverFields, field)
		case "Vehicle":
			vehicleFields = append(vehicleFields, field)
		case "Claims", "ClaimsHistory":
			claimsFields = append(claimsFields, field)
		case "Settings", "BankAccount", "CreditCard", "InsurancePayments", "CommunicationChannel":
			settingsFields = append(settingsFields, field)
		case "PassportDocument", "DrivingLicenceDocument", "IdentityCardDocument", "UtilityBillDocument", "BankStatementDocument", "MedicalCertificateDocument", "InsuranceDocument", "PersonalDocument":
			documentsFields = append(documentsFields, field)
		// Comprehensive Insurance Entity Classes
		case "InsuranceEntity", "MotorInsuranceDocument", "DriverDocument", "VehicleDocument", "ClaimsDocument", "PolicyDocument", "FinancialDocument":
			// Assign to appropriate section based on document type
			if strings.Contains(strings.ToLower(domain), "driver") {
				driverFields = append(driverFields, field)
			} else if strings.Contains(strings.ToLower(domain), "vehicle") {
				vehicleFields = append(vehicleFields, field)
			} else if strings.Contains(strings.ToLower(domain), "claim") {
				claimsFields = append(claimsFields, field)
			} else {
				// Default insurance documents to documents section
				documentsFields = append(documentsFields, field)
			}
		// Specific Insurance Document Sub-classes
		case "InsuranceCertificate", "InsuranceSchedule", "ProofOfNoClaims", "QuoteDocument", "RenewalNotice":
			documentsFields = append(documentsFields, field)
		case "DrivingLicence", "ConvictionCertificate", "MedicalCertificate", "PassPlusCertificate":
			driverFields = append(driverFields, field)
		case "VehicleRegistrationDocument", "MOTCertificate", "VehicleValuation", "ModificationCertificate", "SecurityDeviceCertificate":
			vehicleFields = append(vehicleFields, field)
		case "ClaimForm", "AccidentReport", "PoliceReport", "RepairEstimate", "SettlementLetter":
			claimsFields = append(claimsFields, field)
		case "PolicyWording", "EndorsementDocument", "CancellationNotice":
			documentsFields = append(documentsFields, field)
		case "PremiumInvoice", "PaymentReceipt", "RefundNotice", "DirectDebitMandate":
			documentsFields = append(documentsFields, field)
		// FOAF classes (Person, Organization)
		case "Person", "Organization":
			// These are typically used as ranges for object properties, not domains for datatype properties
			// But if they appear as domains, add to settings
			settingsFields = append(settingsFields, field)
		}
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
