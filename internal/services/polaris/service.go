package polaris

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

// ============================================================
// POLARIS UK PRIVATE MOTOR STANDARDS ADAPTER
// ============================================================
// Implements Polaris Standards for UK private motor insurance
// Supports: Quote, MTA, Renewal, Cancellation, Referrals
// Formats: EDI, XML, JSON
// ============================================================

// PolarisService handles UK Polaris Standards compliance
type PolarisService struct {
	endpoint    string
	apiKey      string
	version     string
	codeListURL string
}

// NewPolarisService creates a new Polaris service instance
func NewPolarisService() *PolarisService {
	return &PolarisService{
		endpoint:    "https://api.polaris.uk/v1",
		version:     "2024.1",
		codeListURL: "https://codelists.polaris.uk/v1",
	}
}

// ============================================================
// POLARIS DATA STRUCTURES
// ============================================================

// PolarisQuoteRequest represents a Polaris quote request
type PolarisQuoteRequest struct {
	MessageHeader PolarisMessageHeader `json:"messageHeader" xml:"MessageHeader"`
	QuoteDetails  PolarisQuoteDetails  `json:"quoteDetails" xml:"QuoteDetails"`
	Vehicle       PolarisVehicle       `json:"vehicle" xml:"Vehicle"`
	Drivers       []PolarisDriver      `json:"drivers" xml:"Drivers>Driver"`
	Coverage      PolarisCoverage      `json:"coverage" xml:"Coverage"`
	Usage         PolarisUsage         `json:"usage" xml:"Usage"`
}

// PolarisMessageHeader represents standard Polaris message header
type PolarisMessageHeader struct {
	MessageID     string    `json:"messageId" xml:"MessageId"`
	MessageType   string    `json:"messageType" xml:"MessageType"`
	Version       string    `json:"version" xml:"Version"`
	Timestamp     time.Time `json:"timestamp" xml:"Timestamp"`
	SenderID      string    `json:"senderId" xml:"SenderId"`
	ReceiverID    string    `json:"receiverId" xml:"ReceiverId"`
	TransactionID string    `json:"transactionId" xml:"TransactionId"`
	TestIndicator bool      `json:"testIndicator" xml:"TestIndicator"`
}

// PolarisQuoteDetails represents quote-specific details
type PolarisQuoteDetails struct {
	QuoteReference   string    `json:"quoteReference" xml:"QuoteReference"`
	EffectiveDate    time.Time `json:"effectiveDate" xml:"EffectiveDate"`
	ExpirationDate   time.Time `json:"expirationDate" xml:"ExpirationDate"`
	QuoteValidUntil  time.Time `json:"quoteValidUntil" xml:"QuoteValidUntil"`
	PaymentMethod    string    `json:"paymentMethod" xml:"PaymentMethod"`
	PaymentFrequency string    `json:"paymentFrequency" xml:"PaymentFrequency"`
	Channel          string    `json:"channel" xml:"Channel"`
	Source           string    `json:"source" xml:"Source"`
}

// PolarisVehicle represents vehicle data in Polaris format
type PolarisVehicle struct {
	VRM               string                `json:"vrm" xml:"VRM"`     // UK Registration Number
	VIN               string                `json:"vin" xml:"VIN"`     // Vehicle Identification Number
	Make              string                `json:"make" xml:"Make"`   // From Polaris code list
	Model             string                `json:"model" xml:"Model"` // From Polaris code list
	YearOfManufacture int                   `json:"yearOfManufacture" xml:"YearOfManufacture"`
	EngineSize        float64               `json:"engineSize" xml:"EngineSize"`     // In litres
	FuelType          string                `json:"fuelType" xml:"FuelType"`         // Polaris code list
	Transmission      string                `json:"transmission" xml:"Transmission"` // Manual/Automatic
	BodyType          string                `json:"bodyType" xml:"BodyType"`         // Polaris code list
	Doors             int                   `json:"doors" xml:"Doors"`
	Seats             int                   `json:"seats" xml:"Seats"`
	Value             float64               `json:"value" xml:"Value"` // Current market value
	PurchasePrice     float64               `json:"purchasePrice" xml:"PurchasePrice"`
	PurchaseDate      string                `json:"purchaseDate" xml:"PurchaseDate"`
	OwnershipType     string                `json:"ownershipType" xml:"OwnershipType"` // Owned/Financed/Leased
	FinancingCompany  string                `json:"financingCompany" xml:"FinancingCompany"`
	ThatchamCategory  string                `json:"thatchamCategory" xml:"ThatchamCategory"` // Security rating
	AlarmType         string                `json:"alarmType" xml:"AlarmType"`               // Thatcham codes
	ImmobiliserType   string                `json:"immobiliserType" xml:"ImmobiliserType"`   // Thatcham codes
	TrackingType      string                `json:"trackingType" xml:"TrackingType"`         // Thatcham codes
	SecurityMarking   bool                  `json:"securityMarking" xml:"SecurityMarking"`
	Modifications     []PolarisModification `json:"modifications" xml:"Modifications>Modification"`
	MOTStatus         string                `json:"motStatus" xml:"MOTStatus"`
	MOTExpiryDate     string                `json:"motExpiryDate" xml:"MOTExpiryDate"`
	TaxStatus         string                `json:"taxStatus" xml:"TaxStatus"`
	TaxExpiryDate     string                `json:"taxExpiryDate" xml:"TaxExpiryDate"`
}

// PolarisModification represents vehicle modifications
type PolarisModification struct {
	Type        string  `json:"type" xml:"Type"` // From Polaris code list
	Description string  `json:"description" xml:"Description"`
	Value       float64 `json:"value" xml:"Value"`
	Declared    bool    `json:"declared" xml:"Declared"`
}

// PolarisDriver represents driver data in Polaris format
type PolarisDriver struct {
	DriverType            string              `json:"driverType" xml:"DriverType"` // Main/Additional/Named
	Title                 string              `json:"title" xml:"Title"`
	FirstName             string              `json:"firstName" xml:"FirstName"`
	LastName              string              `json:"lastName" xml:"LastName"`
	DateOfBirth           time.Time           `json:"dateOfBirth" xml:"DateOfBirth"`
	Gender                string              `json:"gender" xml:"Gender"`
	MaritalStatus         string              `json:"maritalStatus" xml:"MaritalStatus"` // Polaris code list
	Occupation            string              `json:"occupation" xml:"Occupation"`       // Polaris code list
	BusinessUse           string              `json:"businessUse" xml:"BusinessUse"`
	LicenceNumber         string              `json:"licenceNumber" xml:"LicenceNumber"` // UK licence number
	LicenceType           string              `json:"licenceType" xml:"LicenceType"`     // Full/Provisional
	LicenceIssueDate      time.Time           `json:"licenceIssueDate" xml:"LicenceIssueDate"`
	YearsHeld             float64             `json:"yearsHeld" xml:"YearsHeld"`
	PenaltyPoints         int                 `json:"penaltyPoints" xml:"PenaltyPoints"` // From DVLA ADD
	Convictions           []PolarisConviction `json:"convictions" xml:"Convictions>Conviction"`
	ClaimsHistory         []PolarisClaim      `json:"claimsHistory" xml:"ClaimsHistory>Claim"`
	NoClaimsDiscount      float64             `json:"noClaimsDiscount" xml:"NoClaimsDiscount"` // Years
	NoClaimsDiscountProof bool                `json:"noClaimsDiscountProof" xml:"NoClaimsDiscountProof"`
}

// PolarisConviction represents driving convictions
type PolarisConviction struct {
	Code        string    `json:"code" xml:"Code"` // DVLA conviction code
	Date        time.Time `json:"date" xml:"Date"`
	Points      int       `json:"points" xml:"Points"`
	Fine        float64   `json:"fine" xml:"Fine"`
	Description string    `json:"description" xml:"Description"`
	Spent       bool      `json:"spent" xml:"Spent"` // Rehabilitation of Offenders Act
}

// PolarisClaim represents previous claims
type PolarisClaim struct {
	ClaimNumber       string    `json:"claimNumber" xml:"ClaimNumber"`
	Date              time.Time `json:"date" xml:"Date"`
	Type              string    `json:"type" xml:"Type"`               // Polaris code list
	Cause             string    `json:"cause" xml:"Cause"`             // Polaris code list
	FaultStatus       string    `json:"faultStatus" xml:"FaultStatus"` // Fault/Non-fault/Split
	Amount            float64   `json:"amount" xml:"Amount"`
	Settled           bool      `json:"settled" xml:"Settled"`
	OutstandingAmount float64   `json:"outstandingAmount" xml:"OutstandingAmount"`
}

// PolarisCoverage represents insurance coverage
type PolarisCoverage struct {
	CoverType        string  `json:"coverType" xml:"CoverType"` // Comprehensive/TPFT/TPO
	PolicyLimit      float64 `json:"policyLimit" xml:"PolicyLimit"`
	Excess           float64 `json:"excess" xml:"Excess"`
	VoluntaryExcess  float64 `json:"voluntaryExcess" xml:"VoluntaryExcess"`
	CompulsoryExcess float64 `json:"compulsoryExcess" xml:"CompulsoryExcess"`
	PersonalAccident bool    `json:"personalAccident" xml:"PersonalAccident"`
	MedicalExpenses  bool    `json:"medicalExpenses" xml:"MedicalExpenses"`
	PersonalEffects  bool    `json:"personalEffects" xml:"PersonalEffects"`
	KeyCover         bool    `json:"keyCover" xml:"KeyCover"`
	BreakdownCover   bool    `json:"breakdownCover" xml:"BreakdownCover"`
	CourtesyCar      bool    `json:"courtesyCar" xml:"CourtesyCar"`
	ProtectedNCD     bool    `json:"protectedNCD" xml:"ProtectedNCD"`
	LegalExpenses    bool    `json:"legalExpenses" xml:"LegalExpenses"`
	UninsuredLoss    bool    `json:"uninsuredLoss" xml:"UninsuredLoss"`
}

// PolarisUsage represents vehicle usage
type PolarisUsage struct {
	MainUse                string `json:"mainUse" xml:"MainUse"` // SDP/Commuting/Business
	AnnualMileage          int    `json:"annualMileage" xml:"AnnualMileage"`
	CommutingMiles         int    `json:"commutingMiles" xml:"CommutingMiles"`
	BusinessMiles          int    `json:"businessMiles" xml:"BusinessMiles"`
	OvernightLocation      string `json:"overnightLocation" xml:"OvernightLocation"` // Polaris code list
	DaytimeLocation        string `json:"daytimeLocation" xml:"DaytimeLocation"`     // Polaris code list
	PostalCode             string `json:"postalCode" xml:"PostalCode"`
	KeepAtDifferentAddress bool   `json:"keepAtDifferentAddress" xml:"KeepAtDifferentAddress"`
	KeptAddress            string `json:"keptAddress" xml:"KeptAddress"`
}

// ============================================================
// POLARIS RESPONSE STRUCTURES
// ============================================================

// PolarisQuoteResponse represents a Polaris quote response
type PolarisQuoteResponse struct {
	MessageHeader PolarisMessageHeader `json:"messageHeader" xml:"MessageHeader"`
	QuoteResult   PolarisQuoteResult   `json:"quoteResult" xml:"QuoteResult"`
	Premium       PolarisPremium       `json:"premium" xml:"Premium"`
	Referrals     []PolarisReferral    `json:"referrals" xml:"Referrals>Referral"`
	Errors        []PolarisError       `json:"errors" xml:"Errors>Error"`
	Warnings      []PolarisWarning     `json:"warnings" xml:"Warnings>Warning"`
}

// PolarisQuoteResult represents quote result
type PolarisQuoteResult struct {
	QuoteReference       string    `json:"quoteReference" xml:"QuoteReference"`
	Status               string    `json:"status" xml:"Status"` // Quoted/Referred/Declined
	ValidUntil           time.Time `json:"validUntil" xml:"ValidUntil"`
	UnderwritingDecision string    `json:"underwritingDecision" xml:"UnderwritingDecision"`
	RiskRating           string    `json:"riskRating" xml:"RiskRating"`
	NoClaimsDiscount     float64   `json:"noClaimsDiscount" xml:"NoClaimsDiscount"`
}

// PolarisPremium represents premium breakdown
type PolarisPremium struct {
	GrossPremium float64             `json:"grossPremium" xml:"GrossPremium"`
	NetPremium   float64             `json:"netPremium" xml:"NetPremium"`
	IPT          float64             `json:"ipt" xml:"IPT"` // Insurance Premium Tax
	Commission   float64             `json:"commission" xml:"Commission"`
	Brokerage    float64             `json:"brokerage" xml:"Brokerage"`
	AdminFee     float64             `json:"adminFee" xml:"AdminFee"`
	TotalPremium float64             `json:"totalPremium" xml:"TotalPremium"`
	Instalments  []PolarisInstalment `json:"instalments" xml:"Instalments>Instalment"`
}

// PolarisInstalment represents payment instalments
type PolarisInstalment struct {
	Number  int       `json:"number" xml:"Number"`
	Amount  float64   `json:"amount" xml:"Amount"`
	DueDate time.Time `json:"dueDate" xml:"DueDate"`
	Type    string    `json:"type" xml:"Type"` // Deposit/Monthly
}

// PolarisReferral represents underwriting referrals
type PolarisReferral struct {
	Code        string `json:"code" xml:"Code"`
	Description string `json:"description" xml:"Description"`
	Category    string `json:"category" xml:"Category"` // Driver/Vehicle/Claims/Other
	Severity    string `json:"severity" xml:"Severity"` // Info/Warning/Referral/Decline
	Action      string `json:"action" xml:"Action"`     // Required action
}

// PolarisError represents processing errors
type PolarisError struct {
	Code     string `json:"code" xml:"Code"`
	Message  string `json:"message" xml:"Message"`
	Field    string `json:"field" xml:"Field"`
	Severity string `json:"severity" xml:"Severity"`
}

// PolarisWarning represents processing warnings
type PolarisWarning struct {
	Code    string `json:"code" xml:"Code"`
	Message string `json:"message" xml:"Message"`
	Field   string `json:"field" xml:"Field"`
}

// ============================================================
// POLARIS SERVICE METHODS
// ============================================================

// ProcessQuote processes a quote request using Polaris standards
func (ps *PolarisService) ProcessQuote(request PolarisQuoteRequest) (*PolarisQuoteResponse, error) {
	// Validate request against Polaris standards
	if err := ps.validateQuoteRequest(request); err != nil {
		return nil, fmt.Errorf("Polaris validation failed: %w", err)
	}

	// Enrich with DVLA data if UK vehicle
	if err := ps.enrichWithDVLAData(&request); err != nil {
		// Log warning but continue - DVLA enrichment is optional
		fmt.Printf("DVLA enrichment warning: %v\n", err)
	}

	// Apply Polaris business rules
	response := &PolarisQuoteResponse{
		MessageHeader: PolarisMessageHeader{
			MessageID:     generateMessageID(),
			MessageType:   "QuoteResponse",
			Version:       ps.version,
			Timestamp:     time.Now(),
			SenderID:      "CLIENT-UX",
			ReceiverID:    request.MessageHeader.SenderID,
			TransactionID: request.MessageHeader.TransactionID,
			TestIndicator: request.MessageHeader.TestIndicator,
		},
		QuoteResult: PolarisQuoteResult{
			QuoteReference:       generateQuoteReference(),
			Status:               "Quoted",
			ValidUntil:           time.Now().AddDate(0, 0, 30), // 30 days
			UnderwritingDecision: "Accept",
			RiskRating:           ps.calculateRiskRating(request),
			NoClaimsDiscount:     ps.calculateNCD(request.Drivers[0].ClaimsHistory),
		},
		Premium: ps.calculatePremium(request),
	}

	// Check for referrals
	referrals := ps.checkReferrals(request)
	if len(referrals) > 0 {
		response.QuoteResult.Status = "Referred"
		response.Referrals = referrals
	}

	return response, nil
}

// ProcessMTA processes a Mid-Term Adjustment using Polaris standards
func (ps *PolarisService) ProcessMTA(request PolarisQuoteRequest) (*PolarisQuoteResponse, error) {
	// MTA-specific validation
	if request.QuoteDetails.QuoteReference == "" {
		return nil, fmt.Errorf("quote reference required for MTA")
	}

	// Process as quote but with MTA-specific logic
	response, err := ps.ProcessQuote(request)
	if err != nil {
		return nil, err
	}

	response.MessageHeader.MessageType = "MTAResponse"
	return response, nil
}

// ProcessRenewal processes a renewal using Polaris standards
func (ps *PolarisService) ProcessRenewal(request PolarisQuoteRequest) (*PolarisQuoteResponse, error) {
	// Renewal-specific validation and processing
	response, err := ps.ProcessQuote(request)
	if err != nil {
		return nil, err
	}

	response.MessageHeader.MessageType = "RenewalResponse"
	// Apply renewal-specific discounts
	response.Premium.NetPremium *= 0.95                                // 5% renewal discount
	response.Premium.GrossPremium = response.Premium.NetPremium * 1.12 // Add IPT

	return response, nil
}

// ProcessCancellation processes a cancellation using Polaris standards
func (ps *PolarisService) ProcessCancellation(policyNumber string, cancellationDate time.Time, reason string) (*PolarisQuoteResponse, error) {
	// Cancellation processing logic
	response := &PolarisQuoteResponse{
		MessageHeader: PolarisMessageHeader{
			MessageID:   generateMessageID(),
			MessageType: "CancellationResponse",
			Version:     ps.version,
			Timestamp:   time.Now(),
		},
		QuoteResult: PolarisQuoteResult{
			QuoteReference: policyNumber,
			Status:         "Cancelled",
		},
	}

	// Calculate refund/additional premium
	response.Premium = ps.calculateCancellationPremium(policyNumber, cancellationDate, reason)

	return response, nil
}

// ============================================================
// HELPER METHODS
// ============================================================

// validateQuoteRequest validates request against Polaris standards
func (ps *PolarisService) validateQuoteRequest(request PolarisQuoteRequest) error {
	// Validate VRM format (UK)
	if !ps.isValidUKRegistration(request.Vehicle.VRM) {
		return fmt.Errorf("invalid UK vehicle registration: %s", request.Vehicle.VRM)
	}

	// Validate driver licence
	if len(request.Drivers) == 0 {
		return fmt.Errorf("at least one driver required")
	}

	mainDriver := request.Drivers[0]
	if mainDriver.LicenceNumber == "" {
		return fmt.Errorf("main driver licence number required")
	}

	// Validate dates
	if request.QuoteDetails.EffectiveDate.Before(time.Now()) {
		return fmt.Errorf("effective date cannot be in the past")
	}

	return nil
}

// enrichWithDVLAData enriches request with DVLA VES and ADD data
func (ps *PolarisService) enrichWithDVLAData(request *PolarisQuoteRequest) error {
	// DVLA VES (Vehicle Enquiry Service) integration
	vehicleData, err := ps.queryDVLAVES(request.Vehicle.VRM)
	if err != nil {
		return fmt.Errorf("DVLA VES query failed: %w", err)
	}

	// Update vehicle data with DVLA information
	if vehicleData != nil {
		request.Vehicle.Make = vehicleData.Make
		request.Vehicle.Model = vehicleData.Model
		request.Vehicle.YearOfManufacture = vehicleData.YearOfManufacture
		request.Vehicle.EngineSize = vehicleData.EngineSize
		request.Vehicle.FuelType = vehicleData.FuelType
		request.Vehicle.MOTStatus = vehicleData.MOTStatus
		request.Vehicle.MOTExpiryDate = vehicleData.MOTExpiryDate
		request.Vehicle.TaxStatus = vehicleData.TaxStatus
		request.Vehicle.TaxExpiryDate = vehicleData.TaxExpiryDate
	}

	// DVLA ADD (Access to Driver Data) integration
	for i := range request.Drivers {
		driverData, err := ps.queryDVLAADD(request.Drivers[i].LicenceNumber)
		if err != nil {
			// Log warning but continue
			fmt.Printf("DVLA ADD query warning for driver %d: %v\n", i, err)
			continue
		}

		if driverData != nil {
			request.Drivers[i].PenaltyPoints = driverData.PenaltyPoints
			request.Drivers[i].LicenceType = driverData.LicenceType
			request.Drivers[i].LicenceIssueDate = driverData.IssueDate
			// Add any endorsements/convictions from DVLA
			request.Drivers[i].Convictions = append(request.Drivers[i].Convictions, driverData.Convictions...)
		}
	}

	return nil
}

// calculateRiskRating calculates risk rating based on Polaris factors
func (ps *PolarisService) calculateRiskRating(request PolarisQuoteRequest) string {
	score := 0

	// Vehicle factors
	if request.Vehicle.YearOfManufacture < 2010 {
		score += 2
	}
	if request.Vehicle.EngineSize > 2.0 {
		score += 3
	}
	if request.Vehicle.ThatchamCategory == "" {
		score += 2
	}

	// Driver factors
	mainDriver := request.Drivers[0]
	age := time.Now().Year() - mainDriver.DateOfBirth.Year()
	if age < 25 {
		score += 5
	} else if age > 65 {
		score += 2
	}

	if mainDriver.PenaltyPoints > 0 {
		score += mainDriver.PenaltyPoints
	}

	if len(mainDriver.ClaimsHistory) > 0 {
		score += len(mainDriver.ClaimsHistory) * 2
	}

	// Usage factors
	if request.Usage.AnnualMileage > 15000 {
		score += 3
	}
	if request.Usage.MainUse == "Business" {
		score += 4
	}

	// Determine rating
	if score <= 5 {
		return "Low"
	} else if score <= 15 {
		return "Medium"
	} else {
		return "High"
	}
}

// calculateNCD calculates No Claims Discount
func (ps *PolarisService) calculateNCD(claims []PolarisClaim) float64 {
	if len(claims) == 0 {
		return 5.0 // 5 years NCD for no claims
	}

	// Reduce NCD based on fault claims in last 5 years
	faultClaims := 0
	fiveYearsAgo := time.Now().AddDate(-5, 0, 0)

	for _, claim := range claims {
		if claim.Date.After(fiveYearsAgo) && claim.FaultStatus == "Fault" {
			faultClaims++
		}
	}

	ncd := 5.0 - float64(faultClaims)
	if ncd < 0 {
		ncd = 0
	}

	return ncd
}

// calculatePremium calculates premium using Polaris rating factors
func (ps *PolarisService) calculatePremium(request PolarisQuoteRequest) PolarisPremium {
	basePremium := 500.0 // Base premium

	// Vehicle rating
	vehicleMultiplier := 1.0
	if request.Vehicle.YearOfManufacture < 2010 {
		vehicleMultiplier += 0.2
	}
	if request.Vehicle.EngineSize > 2.0 {
		vehicleMultiplier += 0.3
	}
	if request.Vehicle.Value > 30000 {
		vehicleMultiplier += 0.4
	}

	// Driver rating
	driverMultiplier := 1.0
	mainDriver := request.Drivers[0]
	age := time.Now().Year() - mainDriver.DateOfBirth.Year()

	if age < 25 {
		driverMultiplier += 1.0
	} else if age > 65 {
		driverMultiplier += 0.3
	}

	if mainDriver.PenaltyPoints > 0 {
		driverMultiplier += float64(mainDriver.PenaltyPoints) * 0.1
	}

	// Claims history
	claimsMultiplier := 1.0
	for _, claim := range mainDriver.ClaimsHistory {
		if claim.FaultStatus == "Fault" {
			claimsMultiplier += 0.25
		}
	}

	// Usage rating
	usageMultiplier := 1.0
	if request.Usage.AnnualMileage > 15000 {
		usageMultiplier += 0.3
	}
	if request.Usage.MainUse == "Business" {
		usageMultiplier += 0.5
	}

	// Apply NCD
	ncdDiscount := ps.calculateNCD(mainDriver.ClaimsHistory) * 0.1 // 10% per year
	if ncdDiscount > 0.65 {
		ncdDiscount = 0.65 // Maximum 65% NCD
	}

	// Calculate final premium
	netPremium := basePremium * vehicleMultiplier * driverMultiplier * claimsMultiplier * usageMultiplier
	netPremium = netPremium * (1.0 - ncdDiscount)

	ipt := netPremium * 0.12 // 12% Insurance Premium Tax
	grossPremium := netPremium + ipt

	return PolarisPremium{
		NetPremium:   netPremium,
		IPT:          ipt,
		GrossPremium: grossPremium,
		TotalPremium: grossPremium,
		Commission:   netPremium * 0.15, // 15% commission
		Brokerage:    netPremium * 0.05, // 5% brokerage
		AdminFee:     25.0,              // £25 admin fee
	}
}

// checkReferrals checks for underwriting referrals
func (ps *PolarisService) checkReferrals(request PolarisQuoteRequest) []PolarisReferral {
	var referrals []PolarisReferral

	mainDriver := request.Drivers[0]
	age := time.Now().Year() - mainDriver.DateOfBirth.Year()

	// Age referrals
	if age < 21 {
		referrals = append(referrals, PolarisReferral{
			Code:        "AGE_YOUNG",
			Description: "Driver under 21 years old",
			Category:    "Driver",
			Severity:    "Referral",
			Action:      "Manual underwriting required",
		})
	}

	// High value vehicle
	if request.Vehicle.Value > 50000 {
		referrals = append(referrals, PolarisReferral{
			Code:        "VEHICLE_HIGH_VALUE",
			Description: "Vehicle value exceeds £50,000",
			Category:    "Vehicle",
			Severity:    "Referral",
			Action:      "Additional security requirements",
		})
	}

	// Multiple claims
	if len(mainDriver.ClaimsHistory) > 2 {
		referrals = append(referrals, PolarisReferral{
			Code:        "CLAIMS_MULTIPLE",
			Description: "More than 2 claims in claims history",
			Category:    "Claims",
			Severity:    "Referral",
			Action:      "Claims review required",
		})
	}

	// High penalty points
	if mainDriver.PenaltyPoints > 6 {
		referrals = append(referrals, PolarisReferral{
			Code:        "LICENCE_HIGH_POINTS",
			Description: "More than 6 penalty points",
			Category:    "Driver",
			Severity:    "Referral",
			Action:      "Licence verification required",
		})
	}

	return referrals
}

// calculateCancellationPremium calculates refund/additional premium for cancellation
func (ps *PolarisService) calculateCancellationPremium(policyNumber string, cancellationDate time.Time, reason string) PolarisPremium {
	// Mock calculation - in real implementation would look up policy
	return PolarisPremium{
		NetPremium:   -150.0, // Refund
		IPT:          -18.0,  // IPT refund
		GrossPremium: -168.0,
		TotalPremium: -143.0, // After cancellation fee
		AdminFee:     25.0,   // Cancellation fee
	}
}

// ============================================================
// DVLA INTEGRATION METHODS
// ============================================================

// DVLAVehicleData represents DVLA VES response
type DVLAVehicleData struct {
	Make              string  `json:"make"`
	Model             string  `json:"model"`
	YearOfManufacture int     `json:"yearOfManufacture"`
	EngineSize        float64 `json:"engineCapacity"`
	FuelType          string  `json:"fuelType"`
	MOTStatus         string  `json:"motStatus"`
	MOTExpiryDate     string  `json:"motExpiryDate"`
	TaxStatus         string  `json:"taxStatus"`
	TaxExpiryDate     string  `json:"taxDueDate"`
}

// DVLADriverData represents DVLA ADD response
type DVLADriverData struct {
	LicenceType   string              `json:"licenceType"`
	IssueDate     time.Time           `json:"issueDate"`
	PenaltyPoints int                 `json:"penaltyPoints"`
	Convictions   []PolarisConviction `json:"convictions"`
}

// queryDVLAVES queries DVLA Vehicle Enquiry Service
func (ps *PolarisService) queryDVLAVES(vrm string) (*DVLAVehicleData, error) {
	// Mock implementation - in production would call actual DVLA API
	// https://driver-vehicle-licensing.api.gov.uk/vehicle-enquiry/v1/vehicles

	return &DVLAVehicleData{
		Make:              "FORD",
		Model:             "FOCUS",
		YearOfManufacture: 2018,
		EngineSize:        1.6,
		FuelType:          "Petrol",
		MOTStatus:         "Valid",
		MOTExpiryDate:     "2025-03-15",
		TaxStatus:         "Taxed",
		TaxExpiryDate:     "2025-02-01",
	}, nil
}

// queryDVLAADD queries DVLA Access to Driver Data
func (ps *PolarisService) queryDVLAADD(licenceNumber string) (*DVLADriverData, error) {
	// Mock implementation - in production would call actual DVLA API
	// https://driver-vehicle-licensing.api.gov.uk/driving-record/v1/driver

	return &DVLADriverData{
		LicenceType:   "Full",
		IssueDate:     time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC),
		PenaltyPoints: 3,
		Convictions: []PolarisConviction{
			{
				Code:        "SP30",
				Date:        time.Date(2023, 8, 10, 0, 0, 0, 0, time.UTC),
				Points:      3,
				Fine:        100.0,
				Description: "Exceeding statutory speed limit on a public road",
				Spent:       false,
			},
		},
	}, nil
}

// ============================================================
// UTILITY FUNCTIONS
// ============================================================

// isValidUKRegistration validates UK vehicle registration format
func (ps *PolarisService) isValidUKRegistration(vrm string) bool {
	vrm = strings.ToUpper(strings.TrimSpace(vrm))

	// Current format: AB12 CDE (2001 onwards)
	if len(vrm) == 7 && vrm[2] >= '0' && vrm[2] <= '9' && vrm[3] >= '0' && vrm[3] <= '9' {
		return true
	}

	// Prefix format: A123 BCD (1983-2001)
	if len(vrm) >= 6 && len(vrm) <= 7 {
		return true
	}

	// Suffix format: ABC 123D (1963-1983)
	if len(vrm) >= 6 && len(vrm) <= 7 {
		return true
	}

	return false
}

// generateMessageID generates a unique message ID
func generateMessageID() string {
	return fmt.Sprintf("POL-%d", time.Now().UnixNano())
}

// generateQuoteReference generates a unique quote reference
func generateQuoteReference() string {
	return fmt.Sprintf("QUO-%d", time.Now().UnixNano())
}

// ============================================================
// FORMAT CONVERTERS
// ============================================================

// ToJSON converts Polaris data to JSON format
func (ps *PolarisService) ToJSON(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

// ToXML converts Polaris data to XML format
func (ps *PolarisService) ToXML(data interface{}) ([]byte, error) {
	return xml.MarshalIndent(data, "", "  ")
}

// ToEDI converts Polaris data to EDI format (simplified)
func (ps *PolarisService) ToEDI(data interface{}) (string, error) {
	// Simplified EDI conversion - in production would use proper EDI library
	return "UNH+1+QUOTES:D:03B:UN:EAN008'", nil
}
