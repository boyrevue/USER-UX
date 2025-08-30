package reserve

import (
	"fmt"
	"math"
)

// ReserveCalculator implements the reserve calculation logic from the ontology
type ReserveCalculator struct {
	severityTables map[string]map[string]float64
}

// NewReserveCalculator creates a new reserve calculator with predefined severity tables
func NewReserveCalculator() *ReserveCalculator {
	return &ReserveCalculator{
		severityTables: initializeSeverityTables(),
	}
}

// ClaimData represents the input data for reserve calculation
type ClaimData struct {
	LossType           string  // "Collision", "Theft", "Vandalism", etc.
	VehicleACV         float64 // Actual Cash Value
	HasFraudSignals    bool
	LiabilityUncertain bool
	PartsBackorder     bool
}

// ReserveResult represents the calculated reserve with breakdown
type ReserveResult struct {
	BaseReserve       float64 `json:"baseReserve"`
	FraudModifier     float64 `json:"fraudModifier"`
	LiabilityModifier float64 `json:"liabilityModifier"`
	PartsModifier     float64 `json:"partsModifier"`
	FinalReserve      float64 `json:"finalReserve"`
	ReserveBand       string  `json:"reserveBand"`
	Breakdown         string  `json:"breakdown"`
}

// CalculateReserve implements the reserve calculation algorithm from the ontology
func (rc *ReserveCalculator) CalculateReserve(claim ClaimData) ReserveResult {
	// Step 1: Get base reserve from severity table
	base := rc.getSeverityTableAmount(claim.LossType, claim.VehicleACV)

	// Step 2: Apply modifiers
	var mods float64 = 0
	var fraudMod, liabilityMod, partsMod float64

	// Fraud signals: +10% of base
	if claim.HasFraudSignals {
		fraudMod = 0.1 * base
		mods += fraudMod
	}

	// Liability uncertain: +15% of base
	if claim.LiabilityUncertain {
		liabilityMod = 0.15 * base
		mods += liabilityMod
	}

	// Parts backorder: +5% of base
	if claim.PartsBackorder {
		partsMod = 0.05 * base
		mods += partsMod
	}

	// Step 3: Calculate final reserve
	finalReserve := base + mods

	// Step 4: Quantize to band
	reserveBand := rc.quantizeToBand(finalReserve)

	// Step 5: Generate breakdown explanation
	breakdown := rc.generateBreakdown(base, fraudMod, liabilityMod, partsMod, finalReserve)

	return ReserveResult{
		BaseReserve:       base,
		FraudModifier:     fraudMod,
		LiabilityModifier: liabilityMod,
		PartsModifier:     partsMod,
		FinalReserve:      finalReserve,
		ReserveBand:       reserveBand,
		Breakdown:         breakdown,
	}
}

// getSeverityTableAmount looks up base reserve from severity tables
func (rc *ReserveCalculator) getSeverityTableAmount(lossType string, vehicleACV float64) float64 {
	// Get the severity table for this loss type
	table, exists := rc.severityTables[lossType]
	if !exists {
		// Default to "Other" if loss type not found
		table = rc.severityTables["Other"]
	}

	// Determine ACV band and get base amount
	acvBand := rc.getACVBand(vehicleACV)
	baseAmount, exists := table[acvBand]
	if !exists {
		// Default to medium value if band not found
		baseAmount = table["£10k-25k"]
	}

	return baseAmount
}

// getACVBand categorizes vehicle ACV into bands
func (rc *ReserveCalculator) getACVBand(acv float64) string {
	switch {
	case acv < 5000:
		return "£0-5k"
	case acv < 10000:
		return "£5k-10k"
	case acv < 25000:
		return "£10k-25k"
	case acv < 50000:
		return "£25k-50k"
	case acv < 100000:
		return "£50k-100k"
	default:
		return "£100k+"
	}
}

// quantizeToBand converts final reserve to standardized bands
func (rc *ReserveCalculator) quantizeToBand(amount float64) string {
	switch {
	case amount < 2000:
		return "£0-2k"
	case amount < 5000:
		return "£2k-5k"
	case amount < 10000:
		return "£5k-10k"
	case amount < 25000:
		return "£10k-25k"
	case amount < 50000:
		return "£25k-50k"
	default:
		return "£50k+"
	}
}

// generateBreakdown creates a human-readable explanation of the calculation
func (rc *ReserveCalculator) generateBreakdown(base, fraud, liability, parts, final float64) string {
	breakdown := fmt.Sprintf("Base Reserve: £%.2f", base)

	if fraud > 0 {
		breakdown += fmt.Sprintf(" + Fraud Risk: £%.2f", fraud)
	}
	if liability > 0 {
		breakdown += fmt.Sprintf(" + Liability Uncertainty: £%.2f", liability)
	}
	if parts > 0 {
		breakdown += fmt.Sprintf(" + Parts Backorder: £%.2f", parts)
	}

	breakdown += fmt.Sprintf(" = Final Reserve: £%.2f", final)
	return breakdown
}

// initializeSeverityTables creates the lookup tables for base reserves
func initializeSeverityTables() map[string]map[string]float64 {
	return map[string]map[string]float64{
		"Collision": {
			"£0-5k":     1500,
			"£5k-10k":   2500,
			"£10k-25k":  4000,
			"£25k-50k":  6000,
			"£50k-100k": 8000,
			"£100k+":    12000,
		},
		"Theft": {
			"£0-5k":     3000,
			"£5k-10k":   5000,
			"£10k-25k":  8000,
			"£25k-50k":  12000,
			"£50k-100k": 18000,
			"£100k+":    25000,
		},
		"Vandalism": {
			"£0-5k":     800,
			"£5k-10k":   1200,
			"£10k-25k":  1800,
			"£25k-50k":  2500,
			"£50k-100k": 3500,
			"£100k+":    5000,
		},
		"Fire": {
			"£0-5k":     2500,
			"£5k-10k":   4000,
			"£10k-25k":  6500,
			"£25k-50k":  10000,
			"£50k-100k": 15000,
			"£100k+":    22000,
		},
		"Flood": {
			"£0-5k":     2000,
			"£5k-10k":   3500,
			"£10k-25k":  5500,
			"£25k-50k":  8500,
			"£50k-100k": 12000,
			"£100k+":    18000,
		},
		"Glass Damage": {
			"£0-5k":     200,
			"£5k-10k":   300,
			"£10k-25k":  400,
			"£25k-50k":  600,
			"£50k-100k": 800,
			"£100k+":    1200,
		},
		"Third Party": {
			"£0-5k":     5000,
			"£5k-10k":   7500,
			"£10k-25k":  12000,
			"£25k-50k":  18000,
			"£50k-100k": 25000,
			"£100k+":    35000,
		},
		"Comprehensive": {
			"£0-5k":     1800,
			"£5k-10k":   2800,
			"£10k-25k":  4500,
			"£25k-50k":  7000,
			"£50k-100k": 10000,
			"£100k+":    15000,
		},
		"Other": {
			"£0-5k":     1000,
			"£5k-10k":   1500,
			"£10k-25k":  2500,
			"£25k-50k":  4000,
			"£50k-100k": 6000,
			"£100k+":    8000,
		},
	}
}

// FraudScorer assesses fraud risk indicators
type FraudScorer struct {
	riskFactors map[string]float64
}

// NewFraudScorer creates a new fraud scoring service
func NewFraudScorer() *FraudScorer {
	return &FraudScorer{
		riskFactors: initializeFraudRiskFactors(),
	}
}

// FraudAssessment represents fraud risk analysis
type FraudAssessment struct {
	RiskScore      float64            `json:"riskScore"`
	RiskLevel      string             `json:"riskLevel"`
	Indicators     []string           `json:"indicators"`
	Factors        map[string]float64 `json:"factors"`
	Recommendation string             `json:"recommendation"`
}

// AssessFraudRisk evaluates fraud indicators for a claim
func (fs *FraudScorer) AssessFraudRisk(claimData map[string]interface{}) FraudAssessment {
	var totalScore float64
	var indicators []string
	factors := make(map[string]float64)

	// Check various fraud indicators
	if fs.checkIndicator(claimData, "late_reporting") {
		score := fs.riskFactors["late_reporting"]
		totalScore += score
		factors["late_reporting"] = score
		indicators = append(indicators, "Late claim reporting (>7 days)")
	}

	if fs.checkIndicator(claimData, "multiple_claims") {
		score := fs.riskFactors["multiple_claims"]
		totalScore += score
		factors["multiple_claims"] = score
		indicators = append(indicators, "Multiple recent claims")
	}

	if fs.checkIndicator(claimData, "high_value_claim") {
		score := fs.riskFactors["high_value_claim"]
		totalScore += score
		factors["high_value_claim"] = score
		indicators = append(indicators, "Unusually high claim value")
	}

	if fs.checkIndicator(claimData, "inconsistent_story") {
		score := fs.riskFactors["inconsistent_story"]
		totalScore += score
		factors["inconsistent_story"] = score
		indicators = append(indicators, "Inconsistent incident description")
	}

	if fs.checkIndicator(claimData, "no_police_report") {
		score := fs.riskFactors["no_police_report"]
		totalScore += score
		factors["no_police_report"] = score
		indicators = append(indicators, "No police report for significant incident")
	}

	// Determine risk level and recommendation
	riskLevel := fs.getRiskLevel(totalScore)
	recommendation := fs.getRecommendation(riskLevel, totalScore)

	return FraudAssessment{
		RiskScore:      math.Round(totalScore*100) / 100,
		RiskLevel:      riskLevel,
		Indicators:     indicators,
		Factors:        factors,
		Recommendation: recommendation,
	}
}

// checkIndicator evaluates if a specific fraud indicator is present
func (fs *FraudScorer) checkIndicator(data map[string]interface{}, indicator string) bool {
	switch indicator {
	case "late_reporting":
		// Check if claim was reported more than 7 days after incident
		return fs.getBoolValue(data, "late_reporting")
	case "multiple_claims":
		// Check if claimant has multiple recent claims
		return fs.getBoolValue(data, "multiple_recent_claims")
	case "high_value_claim":
		// Check if claim value is unusually high for incident type
		return fs.getBoolValue(data, "high_value_for_type")
	case "inconsistent_story":
		// Check for inconsistencies in incident description
		return fs.getBoolValue(data, "story_inconsistencies")
	case "no_police_report":
		// Check if police report missing for significant incident
		return fs.getBoolValue(data, "missing_police_report")
	default:
		return false
	}
}

// getBoolValue safely extracts boolean values from claim data
func (fs *FraudScorer) getBoolValue(data map[string]interface{}, key string) bool {
	if val, exists := data[key]; exists {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return false
}

// getRiskLevel categorizes the total fraud risk score
func (fs *FraudScorer) getRiskLevel(score float64) string {
	switch {
	case score >= 0.8:
		return "HIGH"
	case score >= 0.5:
		return "MEDIUM"
	case score >= 0.2:
		return "LOW"
	default:
		return "MINIMAL"
	}
}

// getRecommendation provides action recommendations based on risk assessment
func (fs *FraudScorer) getRecommendation(riskLevel string, score float64) string {
	switch riskLevel {
	case "HIGH":
		return "IMMEDIATE INVESTIGATION REQUIRED - Refer to Special Investigation Unit (SIU)"
	case "MEDIUM":
		return "ENHANCED REVIEW - Additional documentation and verification required"
	case "LOW":
		return "STANDARD PROCESSING - Monitor for additional indicators"
	default:
		return "NORMAL PROCESSING - No additional fraud concerns identified"
	}
}

// initializeFraudRiskFactors defines the scoring weights for different fraud indicators
func initializeFraudRiskFactors() map[string]float64 {
	return map[string]float64{
		"late_reporting":      0.3,
		"multiple_claims":     0.4,
		"high_value_claim":    0.25,
		"inconsistent_story":  0.5,
		"no_police_report":    0.2,
		"suspicious_timing":   0.35,
		"unusual_location":    0.15,
		"prior_fraud_history": 0.8,
	}
}
