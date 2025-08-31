package main

import "time"

// Ontology structs
type OntologyData struct {
	Categories map[string]Category      `json:"-"`
	Fields     map[string][]Field       `json:"-"`
	Subforms   map[string]Subform       `json:"-"`
	Schemes    map[string]ConceptScheme `json:"-"`
}

type Category struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Icon        string `json:"icon"`
	Section     string `json:"section"`
	Order       int    `json:"order"`
	Description string `json:"description"`
}

type Field struct {
	Property               string   `json:"property"`
	Label                  string   `json:"label"`
	Type                   string   `json:"type"`
	Required               bool     `json:"required"`
	Pattern                string   `json:"pattern,omitempty"`
	ErrorMessage           string   `json:"errorMessage,omitempty"`
	HelpText               string   `json:"helpText,omitempty"`
	Options                []Option `json:"options,omitempty"`
	ConditionalRequirement string   `json:"conditionalRequirement,omitempty"`
	ConditionalDisplay     string   `json:"conditionalDisplay,omitempty"`
	Min                    int      `json:"min,omitempty"`
	Max                    int      `json:"max,omitempty"`
	Multiple               bool     `json:"multiple,omitempty"`
	TriggerSubform         string   `json:"triggerSubform,omitempty"`
	IsMultiSelect          bool     `json:"isMultiSelect"`
	FormType               string   `json:"formType"`
	EnumerationValues      []string `json:"enumerationValues"`
	ArrayItemStructure     string   `json:"arrayItemStructure,omitempty"`
	FormSection            string   `json:"formSection,omitempty"`
	FormInfoText           string   `json:"formInfoText,omitempty"`
	DefaultValue           string   `json:"defaultValue"`
	RequiresAIValidation   bool     `json:"requiresAIValidation"`
	AIValidationPrompt     string   `json:"aiValidationPrompt,omitempty"`
}

type Option struct {
	Value   string `json:"value"`
	Label   string `json:"label"`
	Trigger string `json:"trigger,omitempty"`
}

type Subform struct {
	ID           string             `json:"id"`
	Title        string             `json:"title"`
	TriggerValue string             `json:"triggerValue"`
	WarningText  string             `json:"warningText,omitempty"`
	Fields       []Field            `json:"fields"`
	Subforms     map[string]Subform `json:"subforms,omitempty"`
}

type ConceptScheme struct{}

// Enhanced Session structs
type Driver struct {
	ID             string                   `json:"id"`
	Classification string                   `json:"classification"` // MAIN or NAMED
	FirstName      string                   `json:"firstName"`
	LastName       string                   `json:"lastName"`
	DateOfBirth    time.Time                `json:"dateOfBirth"`
	Email          string                   `json:"email"`
	Phone          string                   `json:"phone"`
	Address        string                   `json:"address"`
	Postcode       string                   `json:"postcode"`
	LicenceType    string                   `json:"licenceType"`
	LicenceNumber  string                   `json:"licenceNumber"`
	YearsHeld      int                      `json:"yearsHeld"`
	Relationship   string                   `json:"relationship,omitempty"`
	SameAddress    bool                     `json:"sameAddress"`
	HasConvictions bool                     `json:"hasConvictions"`
	Convictions    []map[string]interface{} `json:"convictions,omitempty"`
}

type Vehicle struct {
	Registration      string                 `json:"registration"`
	Make              string                 `json:"make"`
	Model             string                 `json:"model"`
	Year              int                    `json:"year"`
	Mileage           int                    `json:"mileage"`
	Value             float64                `json:"value"`
	OvernightLocation string                 `json:"overnightLocation"`
	HasModifications  bool                   `json:"hasModifications"`
	Modifications     map[string]interface{} `json:"modifications,omitempty"`
}

type Policy struct {
	CoverType       string    `json:"coverType"`
	StartDate       time.Time `json:"startDate"`
	VoluntaryExcess int       `json:"voluntaryExcess"`
	NCDYears        int       `json:"ncdYears"`
	ProtectNCD      bool      `json:"protectNCD"`
}

type ClaimsHistory struct {
	HasClaims bool                     `json:"hasClaims"`
	Claims    []map[string]interface{} `json:"claims,omitempty"`
}

type Payment struct {
	Frequency string `json:"frequency"`
	Method    string `json:"method"`
}

type Extras struct {
	BreakdownCover string `json:"breakdownCover"`
	LegalExpenses  bool   `json:"legalExpenses"`
	CourtesyCar    bool   `json:"courtesyCar"`
}

type Marketing struct {
	EmailMarketing bool `json:"emailMarketing"`
	SMSMarketing   bool `json:"smsMarketing"`
	PostMarketing  bool `json:"postMarketing"`
}

type QuoteSession struct {
	ID           string                            `json:"id"`
	Language     string                            `json:"language"`
	Drivers      []Driver                          `json:"drivers"`
	Vehicle      Vehicle                           `json:"vehicle"`
	Policy       Policy                            `json:"policy"`
	Claims       ClaimsHistory                     `json:"claims"`
	Payment      Payment                           `json:"payment"`
	Extras       Extras                            `json:"extras"`
	Marketing    Marketing                         `json:"marketing"`
	Progress     map[string]bool                   `json:"progress"`
	FormData     map[string]map[string]interface{} `json:"formData"`
	CreatedAt    time.Time                         `json:"createdAt"`
	LastAccessed time.Time                         `json:"lastAccessed"`
	CompletedAt  *time.Time                        `json:"completedAt,omitempty"`
}

// Validation
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors"`
}
