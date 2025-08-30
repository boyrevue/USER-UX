package claims_portal

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"time"
)

// ============================================================
// UK CLAIMS PORTAL A2A INTEGRATION SERVICE
// ============================================================
// Implements integration with:
// - Claims Portal (RTA, EL/PL) - SOAP + REST (dual running until June 2026)
// - Official Injury Claim (OIC) - SOAP + REST
// - Supports FNOL and full claims lifecycle
// ============================================================

// ClaimsPortalService handles UK Claims Portal A2A integration
type ClaimsPortalService struct {
	soapEndpoint string
	restEndpoint string
	oicEndpoint  string
	apiKey       string
	version      string
	dualMode     bool // SOAP + REST dual running until June 2026
}

// NewClaimsPortalService creates a new Claims Portal service
func NewClaimsPortalService() *ClaimsPortalService {
	return &ClaimsPortalService{
		soapEndpoint: "https://et.claimsportal.crif.com/soap",
		restEndpoint: "https://et.claimsportal.crif.com/rest",
		oicEndpoint:  "https://www.officialinjuryclaim.org.uk/api",
		version:      "2024.1",
		dualMode:     true, // Until June 2026
	}
}

// ============================================================
// A2A MESSAGE STRUCTURES
// ============================================================

// A2AMessageHeader represents standard A2A message header
type A2AMessageHeader struct {
	MessageID     string    `json:"messageId" xml:"MessageId"`
	MessageType   string    `json:"messageType" xml:"MessageType"`
	Version       string    `json:"version" xml:"Version"`
	Timestamp     time.Time `json:"timestamp" xml:"Timestamp"`
	SenderID      string    `json:"senderId" xml:"SenderId"`
	ReceiverID    string    `json:"receiverId" xml:"ReceiverId"`
	WorkflowID    string    `json:"workflowId" xml:"WorkflowId"`
	TestIndicator bool      `json:"testIndicator" xml:"TestIndicator"`
	SecurityToken string    `json:"securityToken" xml:"SecurityToken"`
}

// ============================================================
// FNOL (FIRST NOTIFICATION OF LOSS) STRUCTURES
// ============================================================

// FNOLRequest represents First Notification of Loss
type FNOLRequest struct {
	MessageHeader             A2AMessageHeader      `json:"messageHeader" xml:"MessageHeader"`
	ClaimantRepresentativeRef string                `json:"claimantRepresentativeRef" xml:"ClaimantRepresentativeRef"`
	AccidentDetails           AccidentDetails       `json:"accidentDetails" xml:"AccidentDetails"`
	ClaimantDetails           ClaimantDetails       `json:"claimantDetails" xml:"ClaimantDetails"`
	DefendantDetails          DefendantDetails      `json:"defendantDetails" xml:"DefendantDetails"`
	VehicleDetails            A2AVehicleDetails     `json:"vehicleDetails" xml:"VehicleDetails"`
	InjuryDetails             InjuryDetails         `json:"injuryDetails" xml:"InjuryDetails"`
	PropertyDamage            PropertyDamage        `json:"propertyDamage" xml:"PropertyDamage"`
	LiabilityDetails          LiabilityDetails      `json:"liabilityDetails" xml:"LiabilityDetails"`
	RepresentativeDetails     RepresentativeDetails `json:"representativeDetails" xml:"RepresentativeDetails"`
	ECallData                 *ECallData            `json:"eCallData,omitempty" xml:"ECallData,omitempty"`
	TelematicsData            *TelematicsData       `json:"telematicsData,omitempty" xml:"TelematicsData,omitempty"`
}

// AccidentDetails represents accident information
type AccidentDetails struct {
	AccidentDate          time.Time `json:"accidentDate" xml:"AccidentDate"`
	AccidentTime          string    `json:"accidentTime" xml:"AccidentTime"`
	AccidentLocation      string    `json:"accidentLocation" xml:"AccidentLocation"`
	AccidentPostcode      string    `json:"accidentPostcode" xml:"AccidentPostcode"`
	AccidentCountry       string    `json:"accidentCountry" xml:"AccidentCountry"`
	AccidentDescription   string    `json:"accidentDescription" xml:"AccidentDescription"`
	WeatherConditions     string    `json:"weatherConditions" xml:"WeatherConditions"`
	RoadConditions        string    `json:"roadConditions" xml:"RoadConditions"`
	LightConditions       string    `json:"lightConditions" xml:"LightConditions"`
	SpeedLimit            int       `json:"speedLimit" xml:"SpeedLimit"`
	PoliceAttended        bool      `json:"policeAttended" xml:"PoliceAttended"`
	PoliceForce           string    `json:"policeForce" xml:"PoliceForce"`
	PoliceReference       string    `json:"policeReference" xml:"PoliceReference"`
	WitnessesPresent      bool      `json:"witnessesPresent" xml:"WitnessesPresent"`
	NumberOfVehicles      int       `json:"numberOfVehicles" xml:"NumberOfVehicles"`
	AccidentCircumstances string    `json:"accidentCircumstances" xml:"AccidentCircumstances"`
	GeoCoordinates        string    `json:"geoCoordinates,omitempty" xml:"GeoCoordinates,omitempty"`
}

// ClaimantDetails represents claimant information
type ClaimantDetails struct {
	Title                string    `json:"title" xml:"Title"`
	FirstName            string    `json:"firstName" xml:"FirstName"`
	LastName             string    `json:"lastName" xml:"LastName"`
	DateOfBirth          time.Time `json:"dateOfBirth" xml:"DateOfBirth"`
	Gender               string    `json:"gender" xml:"Gender"`
	Address              Address   `json:"address" xml:"Address"`
	ContactDetails       Contact   `json:"contactDetails" xml:"ContactDetails"`
	Occupation           string    `json:"occupation" xml:"Occupation"`
	NationalInsurance    string    `json:"nationalInsurance" xml:"NationalInsurance"`
	DrivingLicence       string    `json:"drivingLicence" xml:"DrivingLicence"`
	RelationshipToDriver string    `json:"relationshipToDriver" xml:"RelationshipToDriver"`
	SeatPosition         string    `json:"seatPosition" xml:"SeatPosition"`
	SeatbeltWorn         bool      `json:"seatbeltWorn" xml:"SeatbeltWorn"`
	HelmetWorn           bool      `json:"helmetWorn" xml:"HelmetWorn"`
}

// DefendantDetails represents defendant/insured information
type DefendantDetails struct {
	Title             string    `json:"title" xml:"Title"`
	FirstName         string    `json:"firstName" xml:"FirstName"`
	LastName          string    `json:"lastName" xml:"LastName"`
	DateOfBirth       time.Time `json:"dateOfBirth" xml:"DateOfBirth"`
	Address           Address   `json:"address" xml:"Address"`
	ContactDetails    Contact   `json:"contactDetails" xml:"ContactDetails"`
	DrivingLicence    string    `json:"drivingLicence" xml:"DrivingLicence"`
	PolicyNumber      string    `json:"policyNumber" xml:"PolicyNumber"`
	InsurerReference  string    `json:"insurerReference" xml:"InsurerReference"`
	VehicleOwner      bool      `json:"vehicleOwner" xml:"VehicleOwner"`
	PermissionToDrive bool      `json:"permissionToDrive" xml:"PermissionToDrive"`
}

// A2AVehicleDetails represents vehicle information for A2A
type A2AVehicleDetails struct {
	VRM                string  `json:"vrm" xml:"VRM"`
	VIN                string  `json:"vin" xml:"VIN"`
	Make               string  `json:"make" xml:"Make"`
	Model              string  `json:"model" xml:"Model"`
	YearOfManufacture  int     `json:"yearOfManufacture" xml:"YearOfManufacture"`
	EngineSize         string  `json:"engineSize" xml:"EngineSize"`
	FuelType           string  `json:"fuelType" xml:"FuelType"`
	Colour             string  `json:"colour" xml:"Colour"`
	VehicleType        string  `json:"vehicleType" xml:"VehicleType"`
	MOTExpiryDate      string  `json:"motExpiryDate" xml:"MOTExpiryDate"`
	TaxExpiryDate      string  `json:"taxExpiryDate" xml:"TaxExpiryDate"`
	InsuranceCompany   string  `json:"insuranceCompany" xml:"InsuranceCompany"`
	PolicyNumber       string  `json:"policyNumber" xml:"PolicyNumber"`
	PolicyStartDate    string  `json:"policyStartDate" xml:"PolicyStartDate"`
	PolicyEndDate      string  `json:"policyEndDate" xml:"PolicyEndDate"`
	LocationOfVehicle  string  `json:"locationOfVehicle" xml:"LocationOfVehicle"`
	VehicleRecoverable bool    `json:"vehicleRecoverable" xml:"VehicleRecoverable"`
	EstimatedDamage    float64 `json:"estimatedDamage" xml:"EstimatedDamage"`
}

// InjuryDetails represents injury information
type InjuryDetails struct {
	InjurySustained       bool     `json:"injurySustained" xml:"InjurySustained"`
	InjuryType            string   `json:"injuryType" xml:"InjuryType"`
	InjurySeverity        string   `json:"injurySeverity" xml:"InjurySeverity"`
	BodyPartsAffected     []string `json:"bodyPartsAffected" xml:"BodyPartsAffected>BodyPart"`
	HospitalAttendance    bool     `json:"hospitalAttendance" xml:"HospitalAttendance"`
	HospitalName          string   `json:"hospitalName" xml:"HospitalName"`
	TreatmentReceived     string   `json:"treatmentReceived" xml:"TreatmentReceived"`
	OngoingTreatment      bool     `json:"ongoingTreatment" xml:"OngoingTreatment"`
	TimeOffWork           bool     `json:"timeOffWork" xml:"TimeOffWork"`
	TimeOffWorkDays       int      `json:"timeOffWorkDays" xml:"TimeOffWorkDays"`
	FullRecoveryExpected  bool     `json:"fullRecoveryExpected" xml:"FullRecoveryExpected"`
	PreExistingConditions bool     `json:"preExistingConditions" xml:"PreExistingConditions"`
}

// PropertyDamage represents property damage details
type PropertyDamage struct {
	PropertyDamaged     bool    `json:"propertyDamaged" xml:"PropertyDamaged"`
	PropertyType        string  `json:"propertyType" xml:"PropertyType"`
	PropertyDescription string  `json:"propertyDescription" xml:"PropertyDescription"`
	EstimatedValue      float64 `json:"estimatedValue" xml:"EstimatedValue"`
	RepairEstimate      float64 `json:"repairEstimate" xml:"RepairEstimate"`
	TotalLoss           bool    `json:"totalLoss" xml:"TotalLoss"`
}

// LiabilityDetails represents liability assessment
type LiabilityDetails struct {
	LiabilityAdmitted   bool     `json:"liabilityAdmitted" xml:"LiabilityAdmitted"`
	LiabilityDenied     bool     `json:"liabilityDenied" xml:"LiabilityDenied"`
	LiabilityDisputed   bool     `json:"liabilityDisputed" xml:"LiabilityDisputed"`
	LiabilityPercentage int      `json:"liabilityPercentage" xml:"LiabilityPercentage"`
	LiabilityReason     string   `json:"liabilityReason" xml:"LiabilityReason"`
	ContributoryFactors []string `json:"contributoryFactors" xml:"ContributoryFactors>Factor"`
}

// RepresentativeDetails represents legal representative information
type RepresentativeDetails struct {
	CompanyName     string  `json:"companyName" xml:"CompanyName"`
	SolicitorName   string  `json:"solicitorName" xml:"SolicitorName"`
	SRANumber       string  `json:"sraNumber" xml:"SRANumber"`
	Address         Address `json:"address" xml:"Address"`
	ContactDetails  Contact `json:"contactDetails" xml:"ContactDetails"`
	ReferenceNumber string  `json:"referenceNumber" xml:"ReferenceNumber"`
}

// Address represents address information
type Address struct {
	AddressLine1 string `json:"addressLine1" xml:"AddressLine1"`
	AddressLine2 string `json:"addressLine2" xml:"AddressLine2"`
	AddressLine3 string `json:"addressLine3" xml:"AddressLine3"`
	Town         string `json:"town" xml:"Town"`
	County       string `json:"county" xml:"County"`
	Postcode     string `json:"postcode" xml:"Postcode"`
	Country      string `json:"country" xml:"Country"`
}

// Contact represents contact information
type Contact struct {
	TelephoneNumber string `json:"telephoneNumber" xml:"TelephoneNumber"`
	MobileNumber    string `json:"mobileNumber" xml:"MobileNumber"`
	EmailAddress    string `json:"emailAddress" xml:"EmailAddress"`
	FaxNumber       string `json:"faxNumber" xml:"FaxNumber"`
}

// ============================================================
// eCALL AND TELEMATICS DATA INTEGRATION
// ============================================================

// ECallData represents eCall MSD (Minimum Set of Data) from EN 15722
type ECallData struct {
	MessageID             string    `json:"messageId" xml:"MessageId"`
	Timestamp             time.Time `json:"timestamp" xml:"Timestamp"`
	Position              Position  `json:"position" xml:"Position"`
	VehicleDirection      float64   `json:"vehicleDirection" xml:"VehicleDirection"`
	CrashSeverity         string    `json:"crashSeverity" xml:"CrashSeverity"`
	NumberOfPassengers    int       `json:"numberOfPassengers" xml:"NumberOfPassengers"`
	VehicleIdentification VehicleID `json:"vehicleIdentification" xml:"VehicleIdentification"`
	AdditionalData        string    `json:"additionalData,omitempty" xml:"AdditionalData,omitempty"`
}

// Position represents GPS coordinates from eCall
type Position struct {
	Latitude  float64 `json:"latitude" xml:"Latitude"`
	Longitude float64 `json:"longitude" xml:"Longitude"`
	Accuracy  float64 `json:"accuracy" xml:"Accuracy"`
}

// VehicleID represents vehicle identification from eCall
type VehicleID struct {
	VIN            string `json:"vin" xml:"VIN"`
	VehicleType    string `json:"vehicleType" xml:"VehicleType"`
	PropulsionType string `json:"propulsionType" xml:"PropulsionType"`
}

// TelematicsData represents additional telematics information
type TelematicsData struct {
	DeviceID         string           `json:"deviceId" xml:"DeviceId"`
	Timestamp        time.Time        `json:"timestamp" xml:"Timestamp"`
	Speed            float64          `json:"speed" xml:"Speed"`
	Acceleration     float64          `json:"acceleration" xml:"Acceleration"`
	Deceleration     float64          `json:"deceleration" xml:"Deceleration"`
	ImpactForce      float64          `json:"impactForce" xml:"ImpactForce"`
	SeatbeltStatus   string           `json:"seatbeltStatus" xml:"SeatbeltStatus"`
	AirbagDeployment bool             `json:"airbagDeployment" xml:"AirbagDeployment"`
	VehicleRollover  bool             `json:"vehicleRollover" xml:"VehicleRollover"`
	EmergencyBraking bool             `json:"emergencyBraking" xml:"EmergencyBraking"`
	DrivingBehaviour DrivingBehaviour `json:"drivingBehaviour" xml:"DrivingBehaviour"`
}

// DrivingBehaviour represents driving behaviour analysis
type DrivingBehaviour struct {
	HarshAcceleration int     `json:"harshAcceleration" xml:"HarshAcceleration"`
	HarshBraking      int     `json:"harshBraking" xml:"HarshBraking"`
	HarshCornering    int     `json:"harshCornering" xml:"HarshCornering"`
	Speeding          int     `json:"speeding" xml:"Speeding"`
	NightDriving      float64 `json:"nightDriving" xml:"NightDriving"`
	MotorwayDriving   float64 `json:"motorwayDriving" xml:"MotorwayDriving"`
	UrbanDriving      float64 `json:"urbanDriving" xml:"UrbanDriving"`
}

// ============================================================
// A2A RESPONSE STRUCTURES
// ============================================================

// FNOLResponse represents FNOL acknowledgment
type FNOLResponse struct {
	MessageHeader      A2AMessageHeader `json:"messageHeader" xml:"MessageHeader"`
	ClaimReference     string           `json:"claimReference" xml:"ClaimReference"`
	Status             string           `json:"status" xml:"Status"`
	AcknowledgmentDate time.Time        `json:"acknowledgmentDate" xml:"AcknowledgmentDate"`
	NextSteps          []string         `json:"nextSteps" xml:"NextSteps>Step"`
	RequiredDocuments  []string         `json:"requiredDocuments" xml:"RequiredDocuments>Document"`
	Errors             []A2AError       `json:"errors,omitempty" xml:"Errors>Error,omitempty"`
	Warnings           []A2AWarning     `json:"warnings,omitempty" xml:"Warnings>Warning,omitempty"`
}

// A2AError represents processing errors
type A2AError struct {
	Code     string `json:"code" xml:"Code"`
	Message  string `json:"message" xml:"Message"`
	Field    string `json:"field" xml:"Field"`
	Severity string `json:"severity" xml:"Severity"`
}

// A2AWarning represents processing warnings
type A2AWarning struct {
	Code    string `json:"code" xml:"Code"`
	Message string `json:"message" xml:"Message"`
	Field   string `json:"field" xml:"Field"`
}

// ============================================================
// CLAIMS LIFECYCLE STRUCTURES
// ============================================================

// LiabilityDecision represents liability decision
type LiabilityDecision struct {
	MessageHeader       A2AMessageHeader `json:"messageHeader" xml:"MessageHeader"`
	ClaimReference      string           `json:"claimReference" xml:"ClaimReference"`
	DecisionDate        time.Time        `json:"decisionDate" xml:"DecisionDate"`
	LiabilityAdmitted   bool             `json:"liabilityAdmitted" xml:"LiabilityAdmitted"`
	LiabilityPercentage int              `json:"liabilityPercentage" xml:"LiabilityPercentage"`
	Reasoning           string           `json:"reasoning" xml:"Reasoning"`
	SupportingEvidence  []string         `json:"supportingEvidence" xml:"SupportingEvidence>Evidence"`
}

// SettlementOffer represents settlement offer
type SettlementOffer struct {
	MessageHeader  A2AMessageHeader `json:"messageHeader" xml:"MessageHeader"`
	ClaimReference string           `json:"claimReference" xml:"ClaimReference"`
	OfferDate      time.Time        `json:"offerDate" xml:"OfferDate"`
	OfferAmount    float64          `json:"offerAmount" xml:"OfferAmount"`
	OfferBreakdown OfferBreakdown   `json:"offerBreakdown" xml:"OfferBreakdown"`
	ValidUntil     time.Time        `json:"validUntil" xml:"ValidUntil"`
	PaymentMethod  string           `json:"paymentMethod" xml:"PaymentMethod"`
	Conditions     []string         `json:"conditions" xml:"Conditions>Condition"`
}

// OfferBreakdown represents settlement offer breakdown
type OfferBreakdown struct {
	GeneralDamages  float64 `json:"generalDamages" xml:"GeneralDamages"`
	SpecialDamages  float64 `json:"specialDamages" xml:"SpecialDamages"`
	LossOfEarnings  float64 `json:"lossOfEarnings" xml:"LossOfEarnings"`
	MedicalExpenses float64 `json:"medicalExpenses" xml:"MedicalExpenses"`
	PropertyDamage  float64 `json:"propertyDamage" xml:"PropertyDamage"`
	LegalCosts      float64 `json:"legalCosts" xml:"LegalCosts"`
	Interest        float64 `json:"interest" xml:"Interest"`
	TotalAmount     float64 `json:"totalAmount" xml:"TotalAmount"`
}

// ============================================================
// SERVICE METHODS
// ============================================================

// SubmitFNOL submits First Notification of Loss to Claims Portal
func (cps *ClaimsPortalService) SubmitFNOL(request FNOLRequest) (*FNOLResponse, error) {
	// Validate FNOL request
	if err := cps.validateFNOLRequest(request); err != nil {
		return nil, fmt.Errorf("FNOL validation failed: %w", err)
	}

	// Parse eCall data if present
	if request.ECallData != nil {
		if err := cps.processECallData(request.ECallData); err != nil {
			// Log warning but continue - eCall data is supplementary
			fmt.Printf("eCall data processing warning: %v\n", err)
		}
	}

	// Submit via both SOAP and REST if dual mode enabled
	var response *FNOLResponse
	var err error

	if cps.dualMode {
		// Try REST first, fallback to SOAP
		response, err = cps.submitFNOLREST(request)
		if err != nil {
			fmt.Printf("REST submission failed, trying SOAP: %v\n", err)
			response, err = cps.submitFNOLSOAP(request)
		}
	} else {
		// Use REST only (post June 2026)
		response, err = cps.submitFNOLREST(request)
	}

	if err != nil {
		return nil, fmt.Errorf("FNOL submission failed: %w", err)
	}

	return response, nil
}

// SubmitLiabilityDecision submits liability decision
func (cps *ClaimsPortalService) SubmitLiabilityDecision(decision LiabilityDecision) error {
	// Validate liability decision
	if decision.ClaimReference == "" {
		return fmt.Errorf("claim reference required")
	}

	// Submit decision
	return cps.submitLiabilityDecision(decision)
}

// SubmitSettlementOffer submits settlement offer
func (cps *ClaimsPortalService) SubmitSettlementOffer(offer SettlementOffer) error {
	// Validate settlement offer
	if offer.ClaimReference == "" {
		return fmt.Errorf("claim reference required")
	}

	if offer.OfferAmount <= 0 {
		return fmt.Errorf("offer amount must be positive")
	}

	// Submit offer
	return cps.submitSettlementOffer(offer)
}

// ============================================================
// VALIDATION METHODS
// ============================================================

// validateFNOLRequest validates FNOL request against A2A standards
func (cps *ClaimsPortalService) validateFNOLRequest(request FNOLRequest) error {
	// Required fields validation
	if request.ClaimantRepresentativeRef == "" {
		return fmt.Errorf("claimant representative reference required")
	}

	if request.AccidentDetails.AccidentDate.IsZero() {
		return fmt.Errorf("accident date required")
	}

	if request.AccidentDetails.AccidentLocation == "" {
		return fmt.Errorf("accident location required")
	}

	// Validate claimant details
	if request.ClaimantDetails.FirstName == "" || request.ClaimantDetails.LastName == "" {
		return fmt.Errorf("claimant name required")
	}

	if request.ClaimantDetails.DateOfBirth.IsZero() {
		return fmt.Errorf("claimant date of birth required")
	}

	// Validate defendant details
	if request.DefendantDetails.PolicyNumber == "" {
		return fmt.Errorf("defendant policy number required")
	}

	// Validate vehicle details
	if request.VehicleDetails.VRM == "" {
		return fmt.Errorf("vehicle registration required")
	}

	// Validate injury details if injury claimed
	if request.InjuryDetails.InjurySustained {
		if request.InjuryDetails.InjuryType == "" {
			return fmt.Errorf("injury type required when injury sustained")
		}
		if request.InjuryDetails.InjurySeverity == "" {
			return fmt.Errorf("injury severity required when injury sustained")
		}
	}

	return nil
}

// ============================================================
// eCALL DATA PROCESSING
// ============================================================

// processECallData processes eCall MSD data according to EN 15722
func (cps *ClaimsPortalService) processECallData(eCallData *ECallData) error {
	// Validate eCall data format
	if eCallData.MessageID == "" {
		return fmt.Errorf("eCall message ID required")
	}

	// Validate position data
	if eCallData.Position.Latitude == 0 && eCallData.Position.Longitude == 0 {
		return fmt.Errorf("eCall position data invalid")
	}

	// Validate crash severity
	validSeverities := []string{"Minor", "Moderate", "Severe", "Fatal"}
	severityValid := false
	for _, severity := range validSeverities {
		if eCallData.CrashSeverity == severity {
			severityValid = true
			break
		}
	}
	if !severityValid {
		return fmt.Errorf("invalid crash severity: %s", eCallData.CrashSeverity)
	}

	// Process vehicle identification
	if eCallData.VehicleIdentification.VIN == "" {
		return fmt.Errorf("eCall vehicle VIN required")
	}

	return nil
}

// ============================================================
// SUBMISSION METHODS
// ============================================================

// submitFNOLREST submits FNOL via REST API
func (cps *ClaimsPortalService) submitFNOLREST(request FNOLRequest) (*FNOLResponse, error) {
	// Mock implementation - in production would call actual Claims Portal REST API
	response := &FNOLResponse{
		MessageHeader: A2AMessageHeader{
			MessageID:   generateA2AMessageID(),
			MessageType: "FNOLResponse",
			Version:     cps.version,
			Timestamp:   time.Now(),
			SenderID:    "CLAIMS-PORTAL",
			ReceiverID:  request.MessageHeader.SenderID,
		},
		ClaimReference:     generateClaimReference(),
		Status:             "Acknowledged",
		AcknowledgmentDate: time.Now(),
		NextSteps: []string{
			"Liability investigation will commence",
			"Medical records will be requested",
			"Vehicle inspection will be arranged",
		},
		RequiredDocuments: []string{
			"Medical records",
			"Repair estimates",
			"Witness statements",
			"Police report",
		},
	}

	return response, nil
}

// submitFNOLSOAP submits FNOL via SOAP API
func (cps *ClaimsPortalService) submitFNOLSOAP(request FNOLRequest) (*FNOLResponse, error) {
	// Mock implementation - in production would call actual Claims Portal SOAP API
	return cps.submitFNOLREST(request) // Same response structure
}

// submitLiabilityDecision submits liability decision
func (cps *ClaimsPortalService) submitLiabilityDecision(decision LiabilityDecision) error {
	// Mock implementation
	fmt.Printf("Liability decision submitted for claim %s: %d%% liability\n",
		decision.ClaimReference, decision.LiabilityPercentage)
	return nil
}

// submitSettlementOffer submits settlement offer
func (cps *ClaimsPortalService) submitSettlementOffer(offer SettlementOffer) error {
	// Mock implementation
	fmt.Printf("Settlement offer submitted for claim %s: £%.2f\n",
		offer.ClaimReference, offer.OfferAmount)
	return nil
}

// ============================================================
// OIC (OFFICIAL INJURY CLAIM) INTEGRATION
// ============================================================

// OICClaim represents Official Injury Claim structure
type OICClaim struct {
	MessageHeader     A2AMessageHeader  `json:"messageHeader" xml:"MessageHeader"`
	ClaimReference    string            `json:"claimReference" xml:"ClaimReference"`
	InjuryDetails     InjuryDetails     `json:"injuryDetails" xml:"InjuryDetails"`
	MedicalEvidence   []MedicalEvidence `json:"medicalEvidence" xml:"MedicalEvidence>Evidence"`
	LossOfEarnings    LossOfEarnings    `json:"lossOfEarnings" xml:"LossOfEarnings"`
	CareAndAssistance CareAndAssistance `json:"careAndAssistance" xml:"CareAndAssistance"`
	TravelExpenses    TravelExpenses    `json:"travelExpenses" xml:"TravelExpenses"`
}

// MedicalEvidence represents medical evidence
type MedicalEvidence struct {
	EvidenceType        string    `json:"evidenceType" xml:"EvidenceType"`
	MedicalProfessional string    `json:"medicalProfessional" xml:"MedicalProfessional"`
	Date                time.Time `json:"date" xml:"Date"`
	Diagnosis           string    `json:"diagnosis" xml:"Diagnosis"`
	Prognosis           string    `json:"prognosis" xml:"Prognosis"`
	TreatmentRequired   string    `json:"treatmentRequired" xml:"TreatmentRequired"`
	Cost                float64   `json:"cost" xml:"Cost"`
}

// LossOfEarnings represents loss of earnings claim
type LossOfEarnings struct {
	GrossWeeklyEarnings float64 `json:"grossWeeklyEarnings" xml:"GrossWeeklyEarnings"`
	NetWeeklyEarnings   float64 `json:"netWeeklyEarnings" xml:"NetWeeklyEarnings"`
	TimeOffWork         int     `json:"timeOffWork" xml:"TimeOffWork"` // Days
	TotalLoss           float64 `json:"totalLoss" xml:"TotalLoss"`
	FutureEarningsLoss  float64 `json:"futureEarningsLoss" xml:"FutureEarningsLoss"`
}

// CareAndAssistance represents care and assistance claim
type CareAndAssistance struct {
	CareRequired  bool    `json:"careRequired" xml:"CareRequired"`
	CareProvider  string  `json:"careProvider" xml:"CareProvider"`
	HoursPerWeek  int     `json:"hoursPerWeek" xml:"HoursPerWeek"`
	WeeksRequired int     `json:"weeksRequired" xml:"WeeksRequired"`
	HourlyRate    float64 `json:"hourlyRate" xml:"HourlyRate"`
	TotalCost     float64 `json:"totalCost" xml:"TotalCost"`
}

// TravelExpenses represents travel expenses
type TravelExpenses struct {
	MedicalAppointments float64 `json:"medicalAppointments" xml:"MedicalAppointments"`
	Physiotherapy       float64 `json:"physiotherapy" xml:"Physiotherapy"`
	Other               float64 `json:"other" xml:"Other"`
	TotalExpenses       float64 `json:"totalExpenses" xml:"TotalExpenses"`
}

// SubmitOICClaim submits claim to Official Injury Claim portal
func (cps *ClaimsPortalService) SubmitOICClaim(claim OICClaim) error {
	// Validate OIC claim
	if claim.ClaimReference == "" {
		return fmt.Errorf("claim reference required for OIC")
	}

	if !claim.InjuryDetails.InjurySustained {
		return fmt.Errorf("injury must be sustained for OIC claim")
	}

	// Submit to OIC portal
	fmt.Printf("OIC claim submitted for reference %s\n", claim.ClaimReference)
	return nil
}

// ============================================================
// UTILITY FUNCTIONS
// ============================================================

// generateA2AMessageID generates unique A2A message ID
func generateA2AMessageID() string {
	return fmt.Sprintf("A2A-%d", time.Now().UnixNano())
}

// generateClaimReference generates unique claim reference
func generateClaimReference() string {
	return fmt.Sprintf("CLM-%d", time.Now().UnixNano())
}

// ToJSON converts A2A data to JSON format
func (cps *ClaimsPortalService) ToJSON(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

// ToXML converts A2A data to XML format
func (cps *ClaimsPortalService) ToXML(data interface{}) ([]byte, error) {
	return xml.MarshalIndent(data, "", "  ")
}
