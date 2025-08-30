package sivi

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ============================================================
// SIVI AFS (NETHERLANDS) ADAPTER SERVICE
// ============================================================
// Implements SIVI AFS JSON Schemas + OpenAPI for Netherlands
// Supports: Policy and Schade (Claims) processing
// Format: JSON with OpenAPI specification
// ============================================================

// SIVIService handles Netherlands SIVI AFS compliance
type SIVIService struct {
	endpoint    string
	apiKey      string
	version     string
	openAPISpec string
}

// NewSIVIService creates a new SIVI service instance
func NewSIVIService() *SIVIService {
	return &SIVIService{
		endpoint:    "https://api.sivi.nl/afs/v1",
		version:     "2024.1",
		openAPISpec: "https://api.sivi.nl/afs/openapi.json",
	}
}

// ============================================================
// SIVI POLICY STRUCTURES (AFS JSON Schema)
// ============================================================

// SIVIPolicy represents a policy in SIVI AFS format
type SIVIPolicy struct {
	PolicyHeader SIVIPolicyHeader  `json:"policyHeader"`
	Policyholder SIVIPolicyholder  `json:"policyholder"`
	Vehicle      SIVIVehicle       `json:"vehicle"`
	Drivers      []SIVIDriver      `json:"drivers"`
	Coverage     SIVICoverage      `json:"coverage"`
	Premium      SIVIPremium       `json:"premium"`
	Terms        SIVITerms         `json:"terms"`
	Endorsements []SIVIEndorsement `json:"endorsements,omitempty"`
	Metadata     SIVIMetadata      `json:"metadata"`
}

// SIVIPolicyHeader represents policy header information
type SIVIPolicyHeader struct {
	PolicyNumber   string    `json:"policyNumber"`
	ProductCode    string    `json:"productCode"`
	ProductName    string    `json:"productName"`
	EffectiveDate  time.Time `json:"effectiveDate"`
	ExpirationDate time.Time `json:"expirationDate"`
	IssueDate      time.Time `json:"issueDate"`
	Status         string    `json:"status"`   // Active, Cancelled, Expired
	Channel        string    `json:"channel"`  // Direct, Broker, Online
	Source         string    `json:"source"`   // Application source
	Currency       string    `json:"currency"` // EUR
	Language       string    `json:"language"` // nl, en
	TestIndicator  bool      `json:"testIndicator"`
}

// SIVIPolicyholder represents policyholder information
type SIVIPolicyholder struct {
	PersonalDetails SIVIPersonalDetails `json:"personalDetails"`
	Address         SIVIAddress         `json:"address"`
	ContactInfo     SIVIContactInfo     `json:"contactInfo"`
	Identification  SIVIIdentification  `json:"identification"`
	Occupation      string              `json:"occupation"`
	MaritalStatus   string              `json:"maritalStatus"`
	Nationality     string              `json:"nationality"`
}

// SIVIPersonalDetails represents personal details
type SIVIPersonalDetails struct {
	Title       string    `json:"title,omitempty"`
	FirstName   string    `json:"firstName"`
	MiddleName  string    `json:"middleName,omitempty"`
	LastName    string    `json:"lastName"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Gender      string    `json:"gender"` // M, F, X
	Initials    string    `json:"initials,omitempty"`
}

// SIVIAddress represents Dutch address format
type SIVIAddress struct {
	Street         string `json:"street"`
	HouseNumber    string `json:"houseNumber"`
	HouseNumberExt string `json:"houseNumberExt,omitempty"`
	PostalCode     string `json:"postalCode"` // Dutch format: 1234AB
	City           string `json:"city"`
	Province       string `json:"province"`
	Country        string `json:"country"`     // NL
	AddressType    string `json:"addressType"` // Residential, Business, Correspondence
}

// SIVIContactInfo represents contact information
type SIVIContactInfo struct {
	PhoneNumber     string `json:"phoneNumber,omitempty"`
	MobileNumber    string `json:"mobileNumber,omitempty"`
	EmailAddress    string `json:"emailAddress"`
	PreferredMethod string `json:"preferredMethod"` // Email, Phone, Post
}

// SIVIIdentification represents identification documents
type SIVIIdentification struct {
	BSN            string `json:"bsn"` // Burgerservicenummer (Dutch SSN)
	PassportNumber string `json:"passportNumber,omitempty"`
	IDCardNumber   string `json:"idCardNumber,omitempty"`
	DrivingLicence string `json:"drivingLicence,omitempty"`
}

// SIVIVehicle represents vehicle information in SIVI format
type SIVIVehicle struct {
	VehicleDetails SIVIVehicleDetails `json:"vehicleDetails"`
	Registration   SIVIRegistration   `json:"registration"`
	Technical      SIVITechnical      `json:"technical"`
	Usage          SIVIUsage          `json:"usage"`
	Security       SIVISecurity       `json:"security"`
	Modifications  []SIVIModification `json:"modifications,omitempty"`
	Valuation      SIVIValuation      `json:"valuation"`
}

// SIVIVehicleDetails represents basic vehicle details
type SIVIVehicleDetails struct {
	VIN               string `json:"vin,omitempty"`
	LicensePlate      string `json:"licensePlate"` // Dutch format
	Make              string `json:"make"`
	Model             string `json:"model"`
	ModelVariant      string `json:"modelVariant,omitempty"`
	YearOfManufacture int    `json:"yearOfManufacture"`
	FirstRegistration string `json:"firstRegistration"`
	VehicleType       string `json:"vehicleType"` // Personenauto, Bedrijfsauto, Motor
	BodyType          string `json:"bodyType"`
	Color             string `json:"color"`
	NumberOfDoors     int    `json:"numberOfDoors"`
	NumberOfSeats     int    `json:"numberOfSeats"`
}

// SIVIRegistration represents registration information
type SIVIRegistration struct {
	RegistrationDate    time.Time `json:"registrationDate"`
	RegistrationCountry string    `json:"registrationCountry"`
	PreviousOwners      int       `json:"previousOwners"`
	ImportedVehicle     bool      `json:"importedVehicle"`
	APKValid            bool      `json:"apkValid"` // Dutch MOT equivalent
	APKExpiryDate       string    `json:"apkExpiryDate,omitempty"`
}

// SIVITechnical represents technical specifications
type SIVITechnical struct {
	EngineCapacity int     `json:"engineCapacity"` // In CC
	EnginePower    int     `json:"enginePower"`    // In kW
	FuelType       string  `json:"fuelType"`       // Benzine, Diesel, Elektro, Hybride
	Transmission   string  `json:"transmission"`   // Handgeschakeld, Automaat
	EmissionClass  string  `json:"emissionClass"`  // Euro 6, etc.
	CO2Emission    int     `json:"co2Emission"`    // g/km
	Weight         int     `json:"weight"`         // In kg
	MaximumSpeed   int     `json:"maximumSpeed,omitempty"`
	Acceleration   float64 `json:"acceleration,omitempty"` // 0-100 km/h
}

// SIVIUsage represents vehicle usage
type SIVIUsage struct {
	MainUse            string `json:"mainUse"`       // Privé, Zakelijk, Woon-werk
	AnnualMileage      int    `json:"annualMileage"` // Kilometers per year
	CommutingDistance  int    `json:"commutingDistance,omitempty"`
	BusinessUse        bool   `json:"businessUse"`
	BusinessPercentage int    `json:"businessPercentage,omitempty"`
	ParkingLocation    string `json:"parkingLocation"` // Garage, Oprit, Straat
	OvernightLocation  string `json:"overnightLocation"`
	PostalCodeUsage    string `json:"postalCodeUsage"` // Where vehicle is primarily used
}

// SIVISecurity represents security features
type SIVISecurity struct {
	AlarmSystem     bool   `json:"alarmSystem"`
	AlarmType       string `json:"alarmType,omitempty"`
	Immobilizer     bool   `json:"immobilizer"`
	ImmobilizerType string `json:"immobilizerType,omitempty"`
	TrackingSystem  bool   `json:"trackingSystem"`
	TrackingType    string `json:"trackingType,omitempty"`
	SecurityMarking bool   `json:"securityMarking"`
	SecurityLevel   string `json:"securityLevel"` // VbV classification
}

// SIVIModification represents vehicle modifications
type SIVIModification struct {
	ModificationType string  `json:"modificationType"`
	Description      string  `json:"description"`
	Value            float64 `json:"value"`
	InstallationDate string  `json:"installationDate,omitempty"`
	Approved         bool    `json:"approved"`
	ApprovalNumber   string  `json:"approvalNumber,omitempty"`
}

// SIVIValuation represents vehicle valuation
type SIVIValuation struct {
	MarketValue      float64   `json:"marketValue"`
	PurchasePrice    float64   `json:"purchasePrice,omitempty"`
	PurchaseDate     string    `json:"purchaseDate,omitempty"`
	ValuationDate    time.Time `json:"valuationDate"`
	ValuationMethod  string    `json:"valuationMethod"` // Taxatie, Catalogus, Eigen opgave
	FinancingType    string    `json:"financingType"`   // Eigen, Lease, Financiering
	FinancingCompany string    `json:"financingCompany,omitempty"`
}

// SIVIDriver represents driver information
type SIVIDriver struct {
	DriverType      string              `json:"driverType"` // Hoofdbestuurder, Medebest
	PersonalDetails SIVIPersonalDetails `json:"personalDetails"`
	Address         SIVIAddress         `json:"address,omitempty"`
	ContactInfo     SIVIContactInfo     `json:"contactInfo,omitempty"`
	Identification  SIVIIdentification  `json:"identification"`
	LicenceDetails  SIVILicenceDetails  `json:"licenceDetails"`
	DrivingHistory  SIVIDrivingHistory  `json:"drivingHistory"`
	Occupation      string              `json:"occupation"`
	MaritalStatus   string              `json:"maritalStatus,omitempty"`
}

// SIVILicenceDetails represents driving licence information
type SIVILicenceDetails struct {
	LicenceNumber    string    `json:"licenceNumber"`
	LicenceType      string    `json:"licenceType"` // B, BE, etc.
	IssueDate        time.Time `json:"issueDate"`
	ExpiryDate       time.Time `json:"expiryDate"`
	IssuingCountry   string    `json:"issuingCountry"`
	YearsHeld        float64   `json:"yearsHeld"`
	PenaltyPoints    int       `json:"penaltyPoints"`
	LicenceSuspended bool      `json:"licenceSuspended"`
	SuspensionPeriod string    `json:"suspensionPeriod,omitempty"`
}

// SIVIDrivingHistory represents driving history
type SIVIDrivingHistory struct {
	NoClaimsYears     float64               `json:"noClaimsYears"`
	NoClaimsProof     bool                  `json:"noClaimsProof"`
	PreviousClaims    []SIVIHistoricalClaim `json:"previousClaims,omitempty"`
	Convictions       []SIVIConviction      `json:"convictions,omitempty"`
	AccidentFreeYears float64               `json:"accidentFreeYears"`
	PreviousInsurer   string                `json:"previousInsurer,omitempty"`
	PreviousPolicyEnd string                `json:"previousPolicyEnd,omitempty"`
}

// SIVICoverage represents insurance coverage
type SIVICoverage struct {
	CoverageType       string           `json:"coverageType"` // WA, Beperkt Casco, Volledig Casco
	Deductible         SIVIDeductible   `json:"deductible"`
	Limits             SIVILimits       `json:"limits"`
	AdditionalCoverage []SIVIAdditional `json:"additionalCoverage,omitempty"`
	Exclusions         []string         `json:"exclusions,omitempty"`
	GeographicalScope  string           `json:"geographicalScope"` // Nederland, Europa, Wereldwijd
}

// SIVIDeductible represents deductible amounts
type SIVIDeductible struct {
	Comprehensive   float64 `json:"comprehensive"`   // Volledig casco
	PartialCoverage float64 `json:"partialCoverage"` // Beperkt casco
	Theft           float64 `json:"theft"`
	Fire            float64 `json:"fire"`
	Glass           float64 `json:"glass"`
	Vandalism       float64 `json:"vandalism"`
}

// SIVILimits represents coverage limits
type SIVILimits struct {
	ThirdPartyLiability float64 `json:"thirdPartyLiability"` // WA dekking
	PersonalAccident    float64 `json:"personalAccident,omitempty"`
	MedicalExpenses     float64 `json:"medicalExpenses,omitempty"`
	LegalExpenses       float64 `json:"legalExpenses,omitempty"`
	PassengerAccident   float64 `json:"passengerAccident,omitempty"`
}

// SIVIAdditional represents additional coverage options
type SIVIAdditional struct {
	CoverageCode        string  `json:"coverageCode"`
	CoverageName        string  `json:"coverageName"`
	CoverageDescription string  `json:"coverageDescription"`
	Limit               float64 `json:"limit,omitempty"`
	Deductible          float64 `json:"deductible,omitempty"`
	Premium             float64 `json:"premium"`
}

// SIVIPremium represents premium information
type SIVIPremium struct {
	BasePremium        float64         `json:"basePremium"`
	Discounts          []SIVIDiscount  `json:"discounts,omitempty"`
	Surcharges         []SIVISurcharge `json:"surcharges,omitempty"`
	NetPremium         float64         `json:"netPremium"`
	Tax                SIVITax         `json:"tax"`
	TotalPremium       float64         `json:"totalPremium"`
	PaymentFrequency   string          `json:"paymentFrequency"` // Jaarlijks, Maandelijks
	PaymentMethod      string          `json:"paymentMethod"`    // Incasso, Factuur
	FirstPayment       float64         `json:"firstPayment"`
	SubsequentPayments float64         `json:"subsequentPayments"`
}

// SIVIDiscount represents premium discounts
type SIVIDiscount struct {
	DiscountCode   string  `json:"discountCode"`
	DiscountName   string  `json:"discountName"`
	DiscountType   string  `json:"discountType"` // Percentage, Amount
	DiscountValue  float64 `json:"discountValue"`
	DiscountAmount float64 `json:"discountAmount"`
}

// SIVISurcharge represents premium surcharges
type SIVISurcharge struct {
	SurchargeCode   string  `json:"surchargeCode"`
	SurchargeName   string  `json:"surchargeName"`
	SurchargeType   string  `json:"surchargeType"` // Percentage, Amount
	SurchargeValue  float64 `json:"surchargeValue"`
	SurchargeAmount float64 `json:"surchargeAmount"`
}

// SIVITax represents tax information
type SIVITax struct {
	VATRate      float64 `json:"vatRate"` // BTW percentage
	VATAmount    float64 `json:"vatAmount"`
	InsuranceTax float64 `json:"insuranceTax"` // Assurantiebelasting
	TotalTax     float64 `json:"totalTax"`
}

// SIVITerms represents policy terms and conditions
type SIVITerms struct {
	PolicyConditions  string   `json:"policyConditions"`
	SpecialConditions []string `json:"specialConditions,omitempty"`
	CancellationTerms string   `json:"cancellationTerms"`
	RenewalTerms      string   `json:"renewalTerms"`
	ClaimsHandling    string   `json:"claimsHandling"`
	DisputeResolution string   `json:"disputeResolution"`
}

// SIVIEndorsement represents policy endorsements
type SIVIEndorsement struct {
	EndorsementNumber string    `json:"endorsementNumber"`
	EndorsementType   string    `json:"endorsementType"`
	EffectiveDate     time.Time `json:"effectiveDate"`
	Description       string    `json:"description"`
	PremiumAdjustment float64   `json:"premiumAdjustment"`
}

// SIVIMetadata represents metadata
type SIVIMetadata struct {
	CreatedDate      time.Time `json:"createdDate"`
	LastModified     time.Time `json:"lastModified"`
	Version          string    `json:"version"`
	Source           string    `json:"source"`
	ProcessingID     string    `json:"processingId"`
	ValidationStatus string    `json:"validationStatus"`
	Errors           []string  `json:"errors,omitempty"`
	Warnings         []string  `json:"warnings,omitempty"`
}

// ============================================================
// SIVI SCHADE (CLAIMS) STRUCTURES
// ============================================================

// SIVISchade represents a claim in SIVI format
type SIVISchade struct {
	ClaimHeader SIVIClaimHeader    `json:"claimHeader"`
	Incident    SIVIIncident       `json:"incident"`
	Parties     []SIVIParty        `json:"parties"`
	Vehicles    []SIVIClaimVehicle `json:"vehicles"`
	Damages     []SIVIDamage       `json:"damages"`
	Injuries    []SIVIInjury       `json:"injuries,omitempty"`
	Liability   SIVILiability      `json:"liability"`
	Settlement  SIVISettlement     `json:"settlement,omitempty"`
	Reserves    *SIVIReserves      `json:"reserves,omitempty"`
	Payments    []SIVIPayment      `json:"payments,omitempty"`
	Documents   []SIVIDocument     `json:"documents,omitempty"`
	Metadata    SIVIMetadata       `json:"metadata"`
}

// SIVIClaimHeader represents claim header
type SIVIClaimHeader struct {
	ClaimNumber       string          `json:"claimNumber"`
	PolicyNumber      string          `json:"policyNumber"`
	ClaimType         string          `json:"claimType"` // Schade, Diefstal, Brand, etc.
	ReportDate        time.Time       `json:"reportDate"`
	Status            string          `json:"status"`   // Gemeld, Onderzoek, Afgehandeld
	Priority          string          `json:"priority"` // Laag, Normaal, Hoog, Urgent
	Handler           string          `json:"handler"`
	HandlerContact    SIVIContactInfo `json:"handlerContact"`
	ExternalReference string          `json:"externalReference,omitempty"`
}

// SIVIIncident represents incident details
type SIVIIncident struct {
	IncidentDate      time.Time     `json:"incidentDate"`
	IncidentTime      string        `json:"incidentTime"`
	Location          SIVILocation  `json:"location"`
	Description       string        `json:"description"`
	Circumstances     string        `json:"circumstances"`
	WeatherConditions string        `json:"weatherConditions"`
	RoadConditions    string        `json:"roadConditions"`
	LightConditions   string        `json:"lightConditions"`
	TrafficSituation  string        `json:"trafficSituation"`
	PoliceInvolved    bool          `json:"policeInvolved"`
	PoliceReport      string        `json:"policeReport,omitempty"`
	Witnesses         []SIVIWitness `json:"witnesses,omitempty"`
}

// SIVILocation represents incident location
type SIVILocation struct {
	Street              string          `json:"street"`
	HouseNumber         string          `json:"houseNumber,omitempty"`
	PostalCode          string          `json:"postalCode"`
	City                string          `json:"city"`
	Province            string          `json:"province"`
	Country             string          `json:"country"`
	Coordinates         SIVICoordinates `json:"coordinates,omitempty"`
	LocationDescription string          `json:"locationDescription,omitempty"`
}

// SIVICoordinates represents GPS coordinates
type SIVICoordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy,omitempty"`
}

// SIVIWitness represents witness information
type SIVIWitness struct {
	PersonalDetails SIVIPersonalDetails `json:"personalDetails"`
	ContactInfo     SIVIContactInfo     `json:"contactInfo"`
	WitnessType     string              `json:"witnessType"` // Ooggetuige, Deskundige
	Statement       string              `json:"statement,omitempty"`
}

// SIVIParty represents involved party
type SIVIParty struct {
	PartyType        string               `json:"partyType"` // Verzekerde, Tegenpartij, Derde
	PersonalDetails  SIVIPersonalDetails  `json:"personalDetails"`
	Address          SIVIAddress          `json:"address"`
	ContactInfo      SIVIContactInfo      `json:"contactInfo"`
	Identification   SIVIIdentification   `json:"identification"`
	InsuranceDetails SIVIInsuranceDetails `json:"insuranceDetails,omitempty"`
	Role             string               `json:"role"` // Bestuurder, Passagier, Voetganger
}

// SIVIInsuranceDetails represents insurance information
type SIVIInsuranceDetails struct {
	InsuranceCompany string    `json:"insuranceCompany"`
	PolicyNumber     string    `json:"policyNumber"`
	CoverageType     string    `json:"coverageType"`
	ValidFrom        time.Time `json:"validFrom"`
	ValidUntil       time.Time `json:"validUntil"`
	Deductible       float64   `json:"deductible"`
}

// SIVIClaimVehicle represents vehicle in claim
type SIVIClaimVehicle struct {
	VehicleDetails      SIVIVehicleDetails `json:"vehicleDetails"`
	DamageDescription   string             `json:"damageDescription"`
	DamageLocation      []string           `json:"damageLocation"` // Voor, Achter, Links, Rechts
	Repairable          bool               `json:"repairable"`
	TotalLoss           bool               `json:"totalLoss"`
	EstimatedRepairCost float64            `json:"estimatedRepairCost"`
	SalvageValue        float64            `json:"salvageValue,omitempty"`
	RecoveryLocation    string             `json:"recoveryLocation,omitempty"`
}

// SIVIDamage represents damage information
type SIVIDamage struct {
	DamageType     string  `json:"damageType"` // Materieel, Letsel, Gevolgschade
	DamageCategory string  `json:"damageCategory"`
	Description    string  `json:"description"`
	EstimatedCost  float64 `json:"estimatedCost"`
	ActualCost     float64 `json:"actualCost,omitempty"`
	Repairer       string  `json:"repairer,omitempty"`
	RepairDate     string  `json:"repairDate,omitempty"`
	RepairStatus   string  `json:"repairStatus,omitempty"`
}

// SIVIInjury represents injury information
type SIVIInjury struct {
	InjuredParty        string  `json:"injuredParty"`
	InjuryType          string  `json:"injuryType"`
	InjuryDescription   string  `json:"injuryDescription"`
	BodyPart            string  `json:"bodyPart"`
	Severity            string  `json:"severity"` // Licht, Matig, Ernstig
	MedicalTreatment    bool    `json:"medicalTreatment"`
	HospitalAdmission   bool    `json:"hospitalAdmission"`
	RecoveryPeriod      string  `json:"recoveryPeriod,omitempty"`
	PermanentDisability bool    `json:"permanentDisability"`
	MedicalCosts        float64 `json:"medicalCosts,omitempty"`
}

// SIVILiability represents liability assessment
type SIVILiability struct {
	LiabilityAssessment string `json:"liabilityAssessment"` // Volledig, Gedeeld, Geen
	LiabilityPercentage int    `json:"liabilityPercentage"`
	LiabilityReason     string `json:"liabilityReason"`
	DisputedLiability   bool   `json:"disputedLiability"`
	LiabilityDate       string `json:"liabilityDate,omitempty"`
}

// SIVISettlement represents settlement information
type SIVISettlement struct {
	SettlementType    string    `json:"settlementType"` // Minnelijk, Gerechtelijk
	SettlementAmount  float64   `json:"settlementAmount"`
	SettlementDate    time.Time `json:"settlementDate"`
	PaymentTerms      string    `json:"paymentTerms"`
	SettlementDetails string    `json:"settlementDetails"`
}

// SIVIReserves represents claim reserves
type SIVIReserves struct {
	InitialReserve float64   `json:"initialReserve"`
	CurrentReserve float64   `json:"currentReserve"`
	ReserveType    string    `json:"reserveType"` // Schade, Kosten, Totaal
	LastUpdate     time.Time `json:"lastUpdate"`
	ReserveReason  string    `json:"reserveReason"`
}

// SIVIPayment represents claim payments
type SIVIPayment struct {
	PaymentType        string    `json:"paymentType"` // Voorschot, Eindafrekening
	PaymentAmount      float64   `json:"paymentAmount"`
	PaymentDate        time.Time `json:"paymentDate"`
	PaymentMethod      string    `json:"paymentMethod"`
	Recipient          string    `json:"recipient"`
	PaymentReference   string    `json:"paymentReference"`
	PaymentDescription string    `json:"paymentDescription"`
}

// SIVIDocument represents claim documents
type SIVIDocument struct {
	DocumentType    string    `json:"documentType"`
	DocumentName    string    `json:"documentName"`
	DocumentDate    time.Time `json:"documentDate"`
	DocumentSize    int       `json:"documentSize"`
	DocumentFormat  string    `json:"documentFormat"`
	DocumentURL     string    `json:"documentUrl,omitempty"`
	DocumentContent string    `json:"documentContent,omitempty"`
}

// SIVIHistoricalClaim represents previous claims in driving history
type SIVIHistoricalClaim struct {
	ClaimDate   time.Time `json:"claimDate"`
	ClaimType   string    `json:"claimType"`
	ClaimAmount float64   `json:"claimAmount"`
	FaultStatus string    `json:"faultStatus"` // Fault/Non-fault/Split
	Settled     bool      `json:"settled"`
	ClaimNumber string    `json:"claimNumber,omitempty"`
}

// SIVIConviction represents traffic convictions
type SIVIConviction struct {
	ConvictionDate    time.Time `json:"convictionDate"`
	ConvictionType    string    `json:"convictionType"`
	ConvictionCode    string    `json:"convictionCode"`
	Description       string    `json:"description"`
	FineAmount        float64   `json:"fineAmount,omitempty"`
	Points            int       `json:"points,omitempty"`
	LicenceSuspension bool      `json:"licenceSuspension"`
	SuspensionPeriod  string    `json:"suspensionPeriod,omitempty"`
}

// ============================================================
// SERVICE METHODS
// ============================================================

// ProcessPolicy processes a policy using SIVI AFS standards
func (ss *SIVIService) ProcessPolicy(policy SIVIPolicy) (*SIVIPolicy, error) {
	// Validate policy against SIVI AFS JSON schema
	if err := ss.validatePolicy(policy); err != nil {
		return nil, fmt.Errorf("SIVI policy validation failed: %w", err)
	}

	// Apply Dutch insurance regulations
	if err := ss.applyDutchRegulations(&policy); err != nil {
		return nil, fmt.Errorf("Dutch regulations compliance failed: %w", err)
	}

	// Calculate premium using Dutch rating factors
	premium, err := ss.calculateDutchPremium(policy)
	if err != nil {
		return nil, fmt.Errorf("premium calculation failed: %w", err)
	}
	policy.Premium = premium

	// Update metadata
	policy.Metadata.LastModified = time.Now()
	policy.Metadata.ValidationStatus = "Valid"
	policy.Metadata.ProcessingID = generateSIVIProcessingID()

	return &policy, nil
}

// ProcessSchade processes a claim using SIVI standards
func (ss *SIVIService) ProcessSchade(schade SIVISchade) (*SIVISchade, error) {
	// Validate claim against SIVI schema
	if err := ss.validateSchade(schade); err != nil {
		return nil, fmt.Errorf("SIVI schade validation failed: %w", err)
	}

	// Apply Dutch claims handling procedures
	if err := ss.applyDutchClaimsHandling(&schade); err != nil {
		return nil, fmt.Errorf("Dutch claims handling failed: %w", err)
	}

	// Calculate reserves using Dutch methods
	reserves, err := ss.calculateDutchReserves(schade)
	if err != nil {
		return nil, fmt.Errorf("reserve calculation failed: %w", err)
	}
	schade.Reserves = &reserves

	// Update metadata
	schade.Metadata.LastModified = time.Now()
	schade.Metadata.ValidationStatus = "Valid"
	schade.Metadata.ProcessingID = generateSIVIProcessingID()

	return &schade, nil
}

// ============================================================
// VALIDATION METHODS
// ============================================================

// validatePolicy validates policy against SIVI AFS JSON schema
func (ss *SIVIService) validatePolicy(policy SIVIPolicy) error {
	// Policy number validation
	if policy.PolicyHeader.PolicyNumber == "" {
		return fmt.Errorf("policy number is required")
	}

	// Dutch postal code validation (1234AB format)
	if !ss.isValidDutchPostalCode(policy.Policyholder.Address.PostalCode) {
		return fmt.Errorf("invalid Dutch postal code: %s", policy.Policyholder.Address.PostalCode)
	}

	// BSN validation (Dutch social security number)
	if !ss.isValidBSN(policy.Policyholder.Identification.BSN) {
		return fmt.Errorf("invalid BSN: %s", policy.Policyholder.Identification.BSN)
	}

	// Dutch license plate validation
	if !ss.isValidDutchLicensePlate(policy.Vehicle.VehicleDetails.LicensePlate) {
		return fmt.Errorf("invalid Dutch license plate: %s", policy.Vehicle.VehicleDetails.LicensePlate)
	}

	// Currency must be EUR for Netherlands
	if policy.PolicyHeader.Currency != "EUR" {
		return fmt.Errorf("currency must be EUR for Dutch policies")
	}

	return nil
}

// validateSchade validates claim against SIVI schema
func (ss *SIVIService) validateSchade(schade SIVISchade) error {
	// Claim number validation
	if schade.ClaimHeader.ClaimNumber == "" {
		return fmt.Errorf("claim number is required")
	}

	// Policy number validation
	if schade.ClaimHeader.PolicyNumber == "" {
		return fmt.Errorf("policy number is required")
	}

	// Incident date validation
	if schade.Incident.IncidentDate.IsZero() {
		return fmt.Errorf("incident date is required")
	}

	// Location validation (must be in Netherlands or covered territory)
	if schade.Incident.Location.Country == "" {
		return fmt.Errorf("incident country is required")
	}

	return nil
}

// ============================================================
// DUTCH REGULATIONS COMPLIANCE
// ============================================================

// applyDutchRegulations applies Dutch insurance regulations
func (ss *SIVIService) applyDutchRegulations(policy *SIVIPolicy) error {
	// Mandatory WA (third-party liability) coverage
	if policy.Coverage.CoverageType == "" {
		policy.Coverage.CoverageType = "WA" // Minimum required
	}

	// Minimum third-party liability limit (€6.07M for personal injury, €1.22M for property)
	if policy.Coverage.Limits.ThirdPartyLiability < 6070000 {
		policy.Coverage.Limits.ThirdPartyLiability = 6070000
	}

	// Dutch tax rates
	policy.Premium.Tax.VATRate = 21.0      // 21% BTW
	policy.Premium.Tax.InsuranceTax = 21.0 // 21% Assurantiebelasting

	// APK (Dutch MOT) requirement for vehicles > 4 years
	vehicleAge := time.Now().Year() - policy.Vehicle.VehicleDetails.YearOfManufacture
	if vehicleAge > 4 && !policy.Vehicle.Registration.APKValid {
		return fmt.Errorf("APK required for vehicles older than 4 years")
	}

	return nil
}

// applyDutchClaimsHandling applies Dutch claims handling procedures
func (ss *SIVIService) applyDutchClaimsHandling(schade *SIVISchade) error {
	// Dutch liability assessment rules
	if schade.Liability.LiabilityAssessment == "" {
		schade.Liability.LiabilityAssessment = "Onderzoek" // Under investigation
	}

	// Mandatory reporting timeframes
	reportDelay := time.Since(schade.Incident.IncidentDate)
	if reportDelay > 48*time.Hour {
		schade.Metadata.Warnings = append(schade.Metadata.Warnings,
			"Claim reported more than 48 hours after incident")
	}

	return nil
}

// ============================================================
// PREMIUM CALCULATION
// ============================================================

// calculateDutchPremium calculates premium using Dutch rating factors
func (ss *SIVIService) calculateDutchPremium(policy SIVIPolicy) (SIVIPremium, error) {
	basePremium := 400.0 // Base premium in EUR

	// Vehicle factors
	vehicleMultiplier := 1.0
	if policy.Vehicle.VehicleDetails.YearOfManufacture < 2015 {
		vehicleMultiplier += 0.15
	}

	// Engine capacity factor (Dutch system)
	engineCC := policy.Vehicle.Technical.EngineCapacity
	if engineCC > 2000 {
		vehicleMultiplier += 0.25
	} else if engineCC > 1600 {
		vehicleMultiplier += 0.10
	}

	// Fuel type factor
	switch policy.Vehicle.Technical.FuelType {
	case "Elektro":
		vehicleMultiplier -= 0.10 // Electric discount
	case "Hybride":
		vehicleMultiplier -= 0.05 // Hybrid discount
	case "Diesel":
		vehicleMultiplier += 0.05 // Diesel surcharge
	}

	// Driver factors
	driverMultiplier := 1.0
	if len(policy.Drivers) > 0 {
		mainDriver := policy.Drivers[0]
		age := time.Now().Year() - mainDriver.PersonalDetails.DateOfBirth.Year()

		if age < 25 {
			driverMultiplier += 0.50 // Young driver surcharge
		} else if age > 65 {
			driverMultiplier += 0.15 // Senior driver surcharge
		}

		// No claims discount
		ncdDiscount := mainDriver.DrivingHistory.NoClaimsYears * 0.05
		if ncdDiscount > 0.70 {
			ncdDiscount = 0.70 // Maximum 70% NCD
		}
		driverMultiplier = driverMultiplier * (1.0 - ncdDiscount)

		// Penalty points surcharge
		if mainDriver.LicenceDetails.PenaltyPoints > 0 {
			driverMultiplier += float64(mainDriver.LicenceDetails.PenaltyPoints) * 0.05
		}
	}

	// Usage factors
	usageMultiplier := 1.0
	if policy.Vehicle.Usage.AnnualMileage > 20000 {
		usageMultiplier += 0.20
	} else if policy.Vehicle.Usage.AnnualMileage > 15000 {
		usageMultiplier += 0.10
	}

	if policy.Vehicle.Usage.BusinessUse {
		usageMultiplier += 0.15
	}

	// Postal code risk factor (simplified)
	postalCode := policy.Policyholder.Address.PostalCode
	if ss.isHighRiskPostalCode(postalCode) {
		usageMultiplier += 0.25
	}

	// Security discount
	securityDiscount := 0.0
	if policy.Vehicle.Security.AlarmSystem {
		securityDiscount += 0.05
	}
	if policy.Vehicle.Security.Immobilizer {
		securityDiscount += 0.03
	}
	if policy.Vehicle.Security.TrackingSystem {
		securityDiscount += 0.07
	}

	// Calculate net premium
	netPremium := basePremium * vehicleMultiplier * driverMultiplier * usageMultiplier
	netPremium = netPremium * (1.0 - securityDiscount)

	// Calculate taxes
	vatAmount := netPremium * 0.21    // 21% BTW
	insuranceTax := netPremium * 0.21 // 21% Assurantiebelasting
	totalTax := vatAmount + insuranceTax
	totalPremium := netPremium + totalTax

	return SIVIPremium{
		BasePremium: basePremium,
		NetPremium:  netPremium,
		Tax: SIVITax{
			VATRate:      21.0,
			VATAmount:    vatAmount,
			InsuranceTax: insuranceTax,
			TotalTax:     totalTax,
		},
		TotalPremium:       totalPremium,
		PaymentFrequency:   "Maandelijks",
		PaymentMethod:      "Incasso",
		FirstPayment:       totalPremium / 12,
		SubsequentPayments: totalPremium / 12,
	}, nil
}

// calculateDutchReserves calculates claim reserves using Dutch methods
func (ss *SIVIService) calculateDutchReserves(schade SIVISchade) (SIVIReserves, error) {
	totalReserve := 0.0

	// Calculate damage reserves
	for _, damage := range schade.Damages {
		if damage.EstimatedCost > 0 {
			totalReserve += damage.EstimatedCost
		}
	}

	// Calculate injury reserves
	for _, injury := range schade.Injuries {
		if injury.MedicalCosts > 0 {
			totalReserve += injury.MedicalCosts
		}

		// Add reserve for potential compensation based on severity
		switch injury.Severity {
		case "Licht":
			totalReserve += 2500 // Light injury reserve
		case "Matig":
			totalReserve += 10000 // Moderate injury reserve
		case "Ernstig":
			totalReserve += 50000 // Severe injury reserve
		}
	}

	// Add handling costs reserve (10% of damage reserve)
	handlingCosts := totalReserve * 0.10
	totalReserve += handlingCosts

	return SIVIReserves{
		InitialReserve: totalReserve,
		CurrentReserve: totalReserve,
		ReserveType:    "Totaal",
		LastUpdate:     time.Now(),
		ReserveReason:  "Initial reserve calculation based on estimated costs",
	}, nil
}

// ============================================================
// VALIDATION UTILITIES
// ============================================================

// isValidDutchPostalCode validates Dutch postal code format (1234AB)
func (ss *SIVIService) isValidDutchPostalCode(postalCode string) bool {
	if len(postalCode) != 6 {
		return false
	}

	// First 4 characters should be digits
	for i := 0; i < 4; i++ {
		if postalCode[i] < '0' || postalCode[i] > '9' {
			return false
		}
	}

	// Last 2 characters should be letters
	for i := 4; i < 6; i++ {
		if postalCode[i] < 'A' || postalCode[i] > 'Z' {
			return false
		}
	}

	return true
}

// isValidBSN validates Dutch BSN (Burgerservicenummer)
func (ss *SIVIService) isValidBSN(bsn string) bool {
	if len(bsn) != 9 {
		return false
	}

	// All characters should be digits
	for _, char := range bsn {
		if char < '0' || char > '9' {
			return false
		}
	}

	// BSN checksum validation (11-test)
	sum := 0
	for i := 0; i < 8; i++ {
		digit := int(bsn[i] - '0')
		sum += digit * (9 - i)
	}

	lastDigit := int(bsn[8] - '0')
	checksum := sum % 11

	return checksum == lastDigit
}

// isValidDutchLicensePlate validates Dutch license plate format
func (ss *SIVIService) isValidDutchLicensePlate(plate string) bool {
	// Remove spaces and hyphens
	plate = strings.ReplaceAll(plate, " ", "")
	plate = strings.ReplaceAll(plate, "-", "")

	if len(plate) != 6 {
		return false
	}

	// Various Dutch formats: 12-AB-34, AB-12-CD, 12-ABC-4, etc.
	// Simplified validation - in production would use comprehensive regex
	return true
}

// isHighRiskPostalCode determines if postal code is in high-risk area
func (ss *SIVIService) isHighRiskPostalCode(postalCode string) bool {
	// Simplified risk assessment based on first 2 digits
	// In production would use comprehensive postal code risk database
	highRiskAreas := []string{"10", "20", "30", "40"} // Amsterdam, Rotterdam, etc.

	if len(postalCode) >= 2 {
		prefix := postalCode[:2]
		for _, riskArea := range highRiskAreas {
			if prefix == riskArea {
				return true
			}
		}
	}

	return false
}

// ============================================================
// UTILITY FUNCTIONS
// ============================================================

// generateSIVIProcessingID generates unique processing ID
func generateSIVIProcessingID() string {
	return fmt.Sprintf("SIVI-%d", time.Now().UnixNano())
}

// ToJSON converts SIVI data to JSON format
func (ss *SIVIService) ToJSON(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

// ValidateJSONSchema validates data against SIVI AFS JSON schema
func (ss *SIVIService) ValidateJSONSchema(data interface{}, schemaType string) error {
	// Mock implementation - in production would validate against actual SIVI AFS JSON schemas
	fmt.Printf("Validating %s against SIVI AFS JSON schema\n", schemaType)
	return nil
}
