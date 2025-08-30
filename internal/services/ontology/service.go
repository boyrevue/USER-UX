package ontology

import (
	"fmt"
	"io/ioutil"
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

type Service struct {
	ontologyPath string
}

func NewService() *Service {
	return &Service{
		ontologyPath: "ontology",
	}
}

func (s *Service) GetFormDefinitions() (map[string]OntologySection, error) {
	// Read the modular auto insurance ontologies
	driverData, err := ioutil.ReadFile(s.ontologyPath + "/AI_Driver_Details.ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to read AI_Driver_Details.ttl: %v", err)
	}

	vehicleData, err := ioutil.ReadFile(s.ontologyPath + "/AI_Vehicle_Details.ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to read AI_Vehicle_Details.ttl: %v", err)
	}

	policyData, err := ioutil.ReadFile(s.ontologyPath + "/AI_Policy_Details.ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to read AI_Policy_Details.ttl: %v", err)
	}

	claimsData, err := ioutil.ReadFile(s.ontologyPath + "/AI_Claims_History.ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to read AI_Claims_History.ttl: %v", err)
	}

	paymentsData, err := ioutil.ReadFile(s.ontologyPath + "/AI_Insurance_Payments.ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to read AI_Insurance_Payments.ttl: %v", err)
	}

	complianceData, err := ioutil.ReadFile(s.ontologyPath + "/AI_Data_Compliance.ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to read AI_Data_Compliance.ttl: %v", err)
	}

	userUxData, err := ioutil.ReadFile(s.ontologyPath + "/user_ux.ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to read user_ux.ttl: %v", err)
	}

	documentsData, err := ioutil.ReadFile(s.ontologyPath + "/user_documents.ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to read user_documents.ttl: %v", err)
	}

	personalDocsData, err := ioutil.ReadFile(s.ontologyPath + "/personal_documents_ontology.ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to read personal_documents_ontology.ttl: %v", err)
	}

	// Combine all TTL content
	_ = string(driverData) + "\n" + string(vehicleData) + "\n" + string(policyData) + "\n" +
		string(claimsData) + "\n" + string(paymentsData) + "\n" + string(complianceData) + "\n" +
		string(userUxData) + "\n" + string(documentsData) + "\n" + string(personalDocsData)

	// Parse the combined TTL content
	sections := make(map[string]OntologySection)

	// Extract driver fields
	driverFields := s.extractFieldsFromTTL(string(driverData), "autoins:Driver")
	sections["drivers"] = OntologySection{
		ID:     "drivers",
		Label:  "Driver Details",
		Fields: driverFields,
	}

	// Extract vehicle fields
	vehicleFields := s.extractFieldsFromTTL(string(vehicleData), "autoins:Vehicle")
	sections["vehicles"] = OntologySection{
		ID:     "vehicles",
		Label:  "Vehicle Details",
		Fields: vehicleFields,
	}

	// Extract claims fields
	claimsFields := s.extractFieldsFromTTL(string(claimsData), "autoins:Claim")
	sections["claims"] = OntologySection{
		ID:     "claims",
		Label:  "Claims History",
		Fields: claimsFields,
	}

	// Extract settings fields
	settingsFields := s.extractFieldsFromTTL(string(userUxData), "autoins:ApplicationSettings")
	sections["settings"] = OntologySection{
		ID:     "settings",
		Label:  "Application Settings",
		Fields: settingsFields,
	}

	return sections, nil
}

func (s *Service) extractFieldsFromTTL(ttlContent, domain string) []OntologyField {
	var fields []OntologyField

	// Regular expressions to extract property definitions
	propertyPattern := regexp.MustCompile(`(?m)^autoins:(\w+)\s+a\s+owl:DatatypeProperty\s*;`)
	labelPattern := regexp.MustCompile(`rdfs:label\s+"([^"]+)"`)
	domainPattern := regexp.MustCompile(`rdfs:domain\s+(\S+)`)
	rangePattern := regexp.MustCompile(`rdfs:range\s+(\S+)`)
	requiredPattern := regexp.MustCompile(`autoins:isRequired\s+"(true|false)"\^\^xsd:boolean`)
	helpTextPattern := regexp.MustCompile(`autoins:formHelpText\s+"([^"]+)"`)

	// Find all property definitions
	propertyMatches := propertyPattern.FindAllStringSubmatch(ttlContent, -1)

	for _, match := range propertyMatches {
		if len(match) < 2 {
			continue
		}

		propertyName := match[1]

		// Extract the full property block
		propertyStart := strings.Index(ttlContent, match[0])
		if propertyStart == -1 {
			continue
		}

		// Find the end of this property block (next property or end of file)
		nextPropertyStart := len(ttlContent)
		nextMatch := propertyPattern.FindStringIndex(ttlContent[propertyStart+len(match[0]):])
		if nextMatch != nil {
			nextPropertyStart = propertyStart + len(match[0]) + nextMatch[0]
		}

		propertyBlock := ttlContent[propertyStart:nextPropertyStart]

		// Check if this property belongs to the specified domain
		domainMatches := domainPattern.FindStringSubmatch(propertyBlock)
		if len(domainMatches) < 2 || !strings.Contains(domainMatches[1], domain) {
			continue
		}

		field := OntologyField{
			Property: propertyName,
			Domain:   domain,
		}

		// Extract label
		if labelMatches := labelPattern.FindStringSubmatch(propertyBlock); len(labelMatches) >= 2 {
			field.Label = labelMatches[1]
		} else {
			field.Label = propertyName
		}

		// Extract type from range
		if rangeMatches := rangePattern.FindStringSubmatch(propertyBlock); len(rangeMatches) >= 2 {
			rangeType := rangeMatches[1]
			switch {
			case strings.Contains(rangeType, "xsd:string"):
				field.Type = "text"
			case strings.Contains(rangeType, "xsd:date"):
				field.Type = "date"
			case strings.Contains(rangeType, "xsd:int") || strings.Contains(rangeType, "xsd:decimal"):
				field.Type = "number"
			case strings.Contains(rangeType, "xsd:boolean"):
				field.Type = "checkbox"
			default:
				field.Type = "text"
			}
		} else {
			field.Type = "text"
		}

		// Extract required status
		if requiredMatches := requiredPattern.FindStringSubmatch(propertyBlock); len(requiredMatches) >= 2 {
			field.Required = requiredMatches[1] == "true"
		}

		// Extract help text
		if helpMatches := helpTextPattern.FindStringSubmatch(propertyBlock); len(helpMatches) >= 2 {
			field.HelpText = helpMatches[1]
		}

		fields = append(fields, field)
	}

	return fields
}
