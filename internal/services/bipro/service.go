package bipro

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"client-ux/internal/services/grounded"
	"client-ux/internal/services/reserve"
)

// BiPROService provides German insurance industry standard compliance
type BiPROService struct {
	classicAdapter *BiPROClassicAdapter
	nextAdapter    *BiPRONextAdapter
	gdvProcessor   *GDVProcessor
	validator      *BiPROValidator
	ontologyMapper *BiPROOntologyMapper
	reserveCalc    *reserve.ReserveCalculator
	groundedAI     *grounded.GroundedPromptEngine
}

// NewBiPROService creates a new BiPRO compliance service
func NewBiPROService() *BiPROService {
	return &BiPROService{
		classicAdapter: NewBiPROClassicAdapter(),
		nextAdapter:    NewBiPRONextAdapter(),
		gdvProcessor:   NewGDVProcessor(),
		validator:      NewBiPROValidator(),
		ontologyMapper: NewBiPROOntologyMapper(),
		reserveCalc:    reserve.NewReserveCalculator(),
		groundedAI:     grounded.NewGroundedPromptEngine(),
	}
}

// ============================================================
// NORM 420 - TARIFICATION, OFFER, APPLICATION (TAA)
// ============================================================

// Norm420TariffRequest represents BiPRO Norm 420 tariff calculation request
type Norm420TariffRequest struct {
	MessageHeader BiPROMessageHeader `json:"messageHeader"`
	RiskData      BiPRORiskData      `json:"riskData"`
	CoverageData  BiPROCoverageData  `json:"coverageData"`
	CustomerData  BiPROCustomerData  `json:"customerData"`
	RequestID     string             `json:"requestId"`
	Timestamp     time.Time          `json:"timestamp"`
}

// Norm420TariffResponse represents BiPRO Norm 420 tariff calculation response
type Norm420TariffResponse struct {
	MessageHeader BiPROMessageHeader `json:"messageHeader"`
	Premium       BiPROPremium       `json:"premium"`
	Conditions    BiPROConditions    `json:"conditions"`
	Validity      BiPROValidity      `json:"validity"`
	Calculations  []BiPROCalculation `json:"calculations"`
	ResponseID    string             `json:"responseId"`
	Timestamp     time.Time          `json:"timestamp"`
}

// BiPRORiskData represents standardized risk assessment data
type BiPRORiskData struct {
	VehicleData  BiPROVehicleData  `json:"vehicleData"`
	DriverData   BiPRODriverData   `json:"driverData"`
	UsageData    BiPROUsageData    `json:"usageData"`
	LocationData BiPROLocationData `json:"locationData"`
	HistoryData  BiPROHistoryData  `json:"historyData"`
}

// BiPROVehicleData represents vehicle information for risk assessment
type BiPROVehicleData struct {
	Make             string   `json:"make"`
	Model            string   `json:"model"`
	Year             int      `json:"year"`
	VIN              string   `json:"vin"`
	Registration     string   `json:"registration"`
	EngineSize       float64  `json:"engineSize"`
	FuelType         string   `json:"fuelType"`
	VehicleValue     float64  `json:"vehicleValue"`
	SecurityFeatures []string `json:"securityFeatures"`
	Modifications    []string `json:"modifications"`
}

// BiPROCoverageData represents insurance coverage information
type BiPROCoverageData struct {
	CoverageType     string  `json:"coverageType"`
	PolicyLimit      float64 `json:"policyLimit"`
	Excess           float64 `json:"excess"`
	NoClaimsDiscount float64 `json:"noClaimsDiscount"`
}

// BiPROCustomerData represents customer information
type BiPROCustomerData struct {
	CustomerID string `json:"customerId"`
	Title      string `json:"title"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	PostalCode string `json:"postalCode"`
}

// BiPROUsageData represents vehicle usage information
type BiPROUsageData struct {
	AnnualMileage    int    `json:"annualMileage"`
	MainUse          string `json:"mainUse"`
	OvernightParking string `json:"overnightParking"`
	DaytimeParking   string `json:"daytimeParking"`
}

// BiPROLocationData represents location-based risk factors
type BiPROLocationData struct {
	PostalCode string `json:"postalCode"`
	RiskArea   string `json:"riskArea"`
	CrimeRate  string `json:"crimeRate"`
	FloodRisk  string `json:"floodRisk"`
}

// BiPROHistoryData represents historical claims and incidents
type BiPROHistoryData struct {
	PreviousClaims   []BiPROClaim      `json:"previousClaims"`
	Convictions      []BiPROConviction `json:"convictions"`
	PreviousPolicies []string          `json:"previousPolicies"`
}

// BiPRODriverData represents driver information
type BiPRODriverData struct {
	DateOfBirth      time.Time         `json:"dateOfBirth"`
	LicenseIssueDate time.Time         `json:"licenseIssueDate"`
	LicenseType      string            `json:"licenseType"`
	Occupation       string            `json:"occupation"`
	MaritalStatus    string            `json:"maritalStatus"`
	Convictions      []BiPROConviction `json:"convictions"`
	ClaimsHistory    []BiPROClaim      `json:"claimsHistory"`
}

// BiPROConviction represents driving conviction data
type BiPROConviction struct {
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`
	Points      int       `json:"points"`
	Fine        float64   `json:"fine"`
	Description string    `json:"description"`
}

// BiPROClaim represents claims history data
type BiPROClaim struct {
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	FaultStatus string    `json:"faultStatus"`
	Settled     bool      `json:"settled"`
}

// ProcessTariffRequest processes BiPRO Norm 420 tariff calculation
func (bs *BiPROService) ProcessTariffRequest(req Norm420TariffRequest) (*Norm420TariffResponse, error) {
	// Validate request against BiPRO schemas
	if err := bs.validator.ValidateNorm420Request(req); err != nil {
		return nil, fmt.Errorf("BiPRO 420 validation failed: %w", err)
	}

	// Map BiPRO data to internal ontology format
	_ = bs.ontologyMapper.MapTariffRequestToInternal(req)

	// Calculate premium using grounded AI and reserve calculator
	reserveData := reserve.ClaimData{
		LossType:           "Comprehensive", // Default for new policy
		VehicleACV:         req.RiskData.VehicleData.VehicleValue,
		HasFraudSignals:    false, // New application
		LiabilityUncertain: false,
		PartsBackorder:     false,
	}

	reserveResult := bs.reserveCalc.CalculateReserve(reserveData)

	// Generate BiPRO compliant response
	response := &Norm420TariffResponse{
		MessageHeader: BiPROMessageHeader{
			MessageID:   generateMessageID(),
			Sender:      "CLIENT-UX",
			Receiver:    req.MessageHeader.Sender,
			NormVersion: "420.2024.1",
		},
		Premium: BiPROPremium{
			AnnualPremium:    reserveResult.BaseReserve * 1.2, // Premium markup
			MonthlyPremium:   (reserveResult.BaseReserve * 1.2) / 12,
			Currency:         "EUR",
			TaxIncluded:      true,
			PaymentFrequency: "MONTHLY",
		},
		Conditions: BiPROConditions{
			Excess:           500.0,
			CoverageType:     "COMPREHENSIVE",
			PolicyTerm:       12,
			NoClaimsDiscount: calculateNCD(req.RiskData.DriverData.ClaimsHistory),
		},
		Validity: BiPROValidity{
			ValidFrom: time.Now(),
			ValidTo:   time.Now().AddDate(0, 0, 30), // 30 days validity
		},
		Calculations: []BiPROCalculation{
			{
				Type:        "BASE_PREMIUM",
				Amount:      reserveResult.BaseReserve,
				Description: "Base premium calculation",
				Factors:     []string{"vehicle_value", "driver_age", "location"},
			},
			{
				Type:        "RISK_ADJUSTMENT",
				Amount:      reserveResult.FinalReserve - reserveResult.BaseReserve,
				Description: "Risk-based adjustments",
				Factors:     []string{"claims_history", "convictions", "vehicle_security"},
			},
		},
		ResponseID: generateResponseID(req.RequestID),
		Timestamp:  time.Now(),
	}

	return response, nil
}

// ============================================================
// NORM 430 - TRANSFER SERVICES
// ============================================================

// Norm430TransferRequest represents document transfer request
type Norm430TransferRequest struct {
	MessageHeader BiPROMessageHeader `json:"messageHeader"`
	TransferType  string             `json:"transferType"` // 430.1, 430.2, 430.4, 430.5, 430.7
	DocumentType  string             `json:"documentType"`
	Format        string             `json:"format"`      // GDV, XML, JSON, PDF
	Compression   string             `json:"compression"` // ZIP, GZIP, NONE
	Data          []byte             `json:"data"`
	Metadata      BiPROMetadata      `json:"metadata"`
}

// Norm430TransferResponse represents document transfer response
type Norm430TransferResponse struct {
	MessageHeader BiPROMessageHeader `json:"messageHeader"`
	Status        string             `json:"status"`
	TransferID    string             `json:"transferId"`
	ProcessedAt   time.Time          `json:"processedAt"`
	Errors        []BiPROError       `json:"errors,omitempty"`
}

// ProcessDocumentTransfer processes BiPRO Norm 430 document transfers
func (bs *BiPROService) ProcessDocumentTransfer(req Norm430TransferRequest) (*Norm430TransferResponse, error) {
	// Validate transfer request
	if err := bs.validator.ValidateNorm430Request(req); err != nil {
		return nil, fmt.Errorf("BiPRO 430 validation failed: %w", err)
	}

	var errors []BiPROError

	// Process based on transfer type
	switch req.TransferType {
	case "430.1": // GDV data transfer
		if err := bs.processGDVTransfer(req); err != nil {
			errors = append(errors, BiPROError{
				Code:     "GDV_PROCESSING_ERROR",
				Message:  err.Error(),
				Severity: "ERROR",
			})
		}
	case "430.2": // Payment irregularities
		if err := bs.processPaymentIrregularities(req); err != nil {
			errors = append(errors, BiPROError{
				Code:     "PAYMENT_PROCESSING_ERROR",
				Message:  err.Error(),
				Severity: "ERROR",
			})
		}
	case "430.4": // Contract transactions
		if err := bs.processContractTransactions(req); err != nil {
			errors = append(errors, BiPROError{
				Code:     "CONTRACT_PROCESSING_ERROR",
				Message:  err.Error(),
				Severity: "ERROR",
			})
		}
	case "430.5": // Claims data
		if err := bs.processClaimsData(req); err != nil {
			errors = append(errors, BiPROError{
				Code:     "CLAIMS_PROCESSING_ERROR",
				Message:  err.Error(),
				Severity: "ERROR",
			})
		}
	default:
		errors = append(errors, BiPROError{
			Code:     "UNSUPPORTED_TRANSFER_TYPE",
			Message:  fmt.Sprintf("Transfer type %s not supported", req.TransferType),
			Severity: "ERROR",
		})
	}

	status := "SUCCESS"
	if len(errors) > 0 {
		status = "ERROR"
	}

	response := &Norm430TransferResponse{
		MessageHeader: BiPROMessageHeader{
			MessageID:   generateMessageID(),
			Sender:      "CLIENT-UX",
			Receiver:    req.MessageHeader.Sender,
			NormVersion: "430.2024.1",
		},
		Status:      status,
		TransferID:  generateTransferID(),
		ProcessedAt: time.Now(),
		Errors:      errors,
	}

	return response, nil
}

// ============================================================
// NORM 440 - DIRECT ACCESS (DEEP LINK)
// ============================================================

// Norm440DeepLinkRequest represents deep link access request
type Norm440DeepLinkRequest struct {
	MessageHeader  BiPROMessageHeader `json:"messageHeader"`
	TargetSystem   string             `json:"targetSystem"`
	TargetFunction string             `json:"targetFunction"`
	Parameters     map[string]string  `json:"parameters"`
	SessionToken   string             `json:"sessionToken"`
	UserID         string             `json:"userId"`
}

// Norm440DeepLinkResponse represents deep link access response
type Norm440DeepLinkResponse struct {
	MessageHeader BiPROMessageHeader `json:"messageHeader"`
	AccessURL     string             `json:"accessUrl"`
	SessionID     string             `json:"sessionId"`
	ExpiresAt     time.Time          `json:"expiresAt"`
	Status        string             `json:"status"`
}

// ProcessDeepLinkRequest processes BiPRO Norm 440 deep link requests
func (bs *BiPROService) ProcessDeepLinkRequest(req Norm440DeepLinkRequest) (*Norm440DeepLinkResponse, error) {
	// Validate deep link request
	if err := bs.validator.ValidateNorm440Request(req); err != nil {
		return nil, fmt.Errorf("BiPRO 440 validation failed: %w", err)
	}

	// Generate secure access URL with embedded parameters
	accessURL := bs.generateSecureAccessURL(req.TargetSystem, req.TargetFunction, req.Parameters)

	// Create session for SSO
	sessionID := generateSessionID()
	expiresAt := time.Now().Add(30 * time.Minute) // 30 minute session

	response := &Norm440DeepLinkResponse{
		MessageHeader: BiPROMessageHeader{
			MessageID:   generateMessageID(),
			Sender:      "CLIENT-UX",
			Receiver:    req.MessageHeader.Sender,
			NormVersion: "440.2024.1",
		},
		AccessURL: accessURL,
		SessionID: sessionID,
		ExpiresAt: expiresAt,
		Status:    "SUCCESS",
	}

	return response, nil
}

// ============================================================
// SUPPORTING DATA STRUCTURES
// ============================================================

// BiPROMessageHeader represents standard BiPRO message header
type BiPROMessageHeader struct {
	MessageID   string    `json:"messageId"`
	Sender      string    `json:"sender"`
	Receiver    string    `json:"receiver"`
	NormVersion string    `json:"normVersion"`
	Timestamp   time.Time `json:"timestamp"`
}

// BiPROPremium represents premium calculation results
type BiPROPremium struct {
	AnnualPremium    float64 `json:"annualPremium"`
	MonthlyPremium   float64 `json:"monthlyPremium"`
	Currency         string  `json:"currency"`
	TaxIncluded      bool    `json:"taxIncluded"`
	PaymentFrequency string  `json:"paymentFrequency"`
}

// BiPROConditions represents policy conditions
type BiPROConditions struct {
	Excess           float64 `json:"excess"`
	CoverageType     string  `json:"coverageType"`
	PolicyTerm       int     `json:"policyTerm"`
	NoClaimsDiscount float64 `json:"noClaimsDiscount"`
}

// BiPROValidity represents offer validity period
type BiPROValidity struct {
	ValidFrom time.Time `json:"validFrom"`
	ValidTo   time.Time `json:"validTo"`
}

// BiPROCalculation represents individual calculation components
type BiPROCalculation struct {
	Type        string   `json:"type"`
	Amount      float64  `json:"amount"`
	Description string   `json:"description"`
	Factors     []string `json:"factors"`
}

// BiPROMetadata represents transfer metadata
type BiPROMetadata struct {
	DocumentID   string            `json:"documentId"`
	DocumentDate time.Time         `json:"documentDate"`
	PolicyNumber string            `json:"policyNumber"`
	CustomerID   string            `json:"customerId"`
	Properties   map[string]string `json:"properties"`
}

// BiPROError represents BiPRO error information
type BiPROError struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
	Field    string `json:"field,omitempty"`
}

// ============================================================
// HELPER FUNCTIONS
// ============================================================

// generateMessageID creates a unique BiPRO message identifier
func generateMessageID() string {
	timestamp := time.Now().Format("20060102150405")
	hash := sha256.Sum256([]byte(timestamp + "CLIENT-UX"))
	return fmt.Sprintf("CUX-%s-%s", timestamp, hex.EncodeToString(hash[:4]))
}

// generateResponseID creates a response ID based on request ID
func generateResponseID(requestID string) string {
	return fmt.Sprintf("RESP-%s-%d", requestID, time.Now().Unix())
}

// generateTransferID creates a unique transfer identifier
func generateTransferID() string {
	return fmt.Sprintf("TXF-%d", time.Now().UnixNano())
}

// generateSessionID creates a unique session identifier
func generateSessionID() string {
	return fmt.Sprintf("SES-%d", time.Now().UnixNano())
}

// calculateNCD calculates No Claims Discount based on claims history
func calculateNCD(claims []BiPROClaim) float64 {
	if len(claims) == 0 {
		return 0.65 // 65% NCD for no claims
	}

	// Reduce NCD based on number of claims
	ncd := 0.65 - (float64(len(claims)) * 0.15)
	if ncd < 0 {
		ncd = 0
	}

	return ncd
}

// generateSecureAccessURL creates a secure deep link URL
func (bs *BiPROService) generateSecureAccessURL(system, function string, params map[string]string) string {
	baseURL := fmt.Sprintf("https://client-ux.example.com/bipro/access/%s/%s", system, function)

	// Add parameters as query string
	if len(params) > 0 {
		var paramPairs []string
		for key, value := range params {
			paramPairs = append(paramPairs, fmt.Sprintf("%s=%s", key, value))
		}
		baseURL += "?" + strings.Join(paramPairs, "&")
	}

	return baseURL
}

// ============================================================
// PROCESSING FUNCTIONS
// ============================================================

// processGDVTransfer processes GDV format data transfers
func (bs *BiPROService) processGDVTransfer(req Norm430TransferRequest) error {
	return bs.gdvProcessor.ProcessGDVData(req.Data, req.Metadata)
}

// processPaymentIrregularities processes payment irregularity notifications
func (bs *BiPROService) processPaymentIrregularities(req Norm430TransferRequest) error {
	// Process payment reminders, dunning notices, etc.
	// Implementation would integrate with payment processing system
	return nil
}

// processContractTransactions processes contract-related business transactions
func (bs *BiPROService) processContractTransactions(req Norm430TransferRequest) error {
	// Process contract changes, renewals, cancellations, etc.
	// Implementation would integrate with policy management system
	return nil
}

// processClaimsData processes claims and benefit-related data
func (bs *BiPROService) processClaimsData(req Norm430TransferRequest) error {
	// Process claims notifications, status updates, settlements, etc.
	// Implementation would integrate with claims processing system
	return nil
}
