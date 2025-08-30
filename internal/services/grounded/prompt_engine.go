package grounded

import (
	"fmt"
	"math"
	"strings"
	"time"

	"client-ux/internal/services/reserve"
)

// GroundedPromptEngine implements the grounded AI pattern from the ontology
type GroundedPromptEngine struct {
	reserveCalc   *reserve.ReserveCalculator
	fraudScorer   *reserve.FraudScorer
	ontologyGraph map[string]interface{} // Simplified graph representation
}

// NewGroundedPromptEngine creates a new grounded prompt processing engine
func NewGroundedPromptEngine() *GroundedPromptEngine {
	return &GroundedPromptEngine{
		reserveCalc:   reserve.NewReserveCalculator(),
		fraudScorer:   reserve.NewFraudScorer(),
		ontologyGraph: initializeOntologyGraph(),
	}
}

// GroundedRequest represents an incoming request for grounded AI processing
type GroundedRequest struct {
	UserQuery     string                 `json:"userQuery"`
	GraphContext  map[string]interface{} `json:"graphContext"`
	RequiredTools []string               `json:"requiredTools"`
	Language      string                 `json:"language,omitempty"`
}

// GroundedResponse represents the structured response with fact grounding
type GroundedResponse struct {
	Answer                string            `json:"answer"`
	FactsUsed             []FactReference   `json:"factsUsed"`
	CalculationsPerformed []CalculationRef  `json:"calculationsPerformed"`
	FollowUpNeeded        bool              `json:"followUpNeeded"`
	FollowUpQuestion      string            `json:"followUpQuestion,omitempty"`
	ConfidenceLevel       float64           `json:"confidenceLevel"`
	ProcessingTime        time.Duration     `json:"processingTime"`
	ValidationErrors      []ValidationError `json:"validationErrors,omitempty"`
}

// FactReference represents a specific graph fact used in reasoning
type FactReference struct {
	IRI        string  `json:"iri"`
	Value      string  `json:"value"`
	Property   string  `json:"property"`
	Source     string  `json:"source"`
	Confidence float64 `json:"confidence"`
}

// CalculationRef represents a calculation performed during processing
type CalculationRef struct {
	Type       string                 `json:"type"`
	Input      map[string]interface{} `json:"input"`
	Result     interface{}            `json:"result"`
	Confidence float64                `json:"confidence"`
	Timestamp  time.Time              `json:"timestamp"`
}

// ValidationError represents SHACL validation failures
type ValidationError struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
	Field    string `json:"field,omitempty"`
}

// ProcessGroundedQuery processes a user query using the grounded AI pattern
func (gpe *GroundedPromptEngine) ProcessGroundedQuery(req GroundedRequest) GroundedResponse {
	startTime := time.Now()

	response := GroundedResponse{
		FactsUsed:             []FactReference{},
		CalculationsPerformed: []CalculationRef{},
		FollowUpNeeded:        false,
		ValidationErrors:      []ValidationError{},
	}

	// Step 1: Validate request against SHACL rules
	validationErrors := gpe.validateRequest(req)
	if len(validationErrors) > 0 {
		response.ValidationErrors = validationErrors
		response.Answer = "Request validation failed. Please correct the errors and try again."
		response.ProcessingTime = time.Since(startTime)
		return response
	}

	// Step 2: Execute SPARQL queries to gather facts
	facts, err := gpe.executeSPARQLQueries(req)
	if err != nil {
		response.Answer = fmt.Sprintf("Failed to retrieve facts from knowledge graph: %v", err)
		response.ProcessingTime = time.Since(startTime)
		return response
	}
	response.FactsUsed = facts

	// Step 3: Perform required calculations
	calculations := gpe.performCalculations(req, facts)
	response.CalculationsPerformed = calculations

	// Step 4: Check for missing critical information
	missingInfo := gpe.checkMissingInformation(req, facts)
	if len(missingInfo) > 0 {
		response.FollowUpNeeded = true
		response.FollowUpQuestion = gpe.generateFollowUpQuestion(missingInfo)
	}

	// Step 5: Generate grounded response
	response.Answer = gpe.generateGroundedAnswer(req, facts, calculations)
	response.ConfidenceLevel = gpe.calculateConfidence(facts, calculations)
	response.ProcessingTime = time.Since(startTime)

	return response
}

// validateRequest validates the request against SHACL shapes
func (gpe *GroundedPromptEngine) validateRequest(req GroundedRequest) []ValidationError {
	var errors []ValidationError

	// Check if user query is provided
	if strings.TrimSpace(req.UserQuery) == "" {
		errors = append(errors, ValidationError{
			Code:     "MISSING_USER_QUERY",
			Message:  "User query is required for grounded processing",
			Severity: "ERROR",
			Field:    "userQuery",
		})
	}

	// Check if required tools are specified
	if len(req.RequiredTools) == 0 {
		errors = append(errors, ValidationError{
			Code:     "MISSING_REQUIRED_TOOLS",
			Message:  "At least one tool must be specified for grounded processing",
			Severity: "ERROR",
			Field:    "requiredTools",
		})
	}

	// Validate tool availability
	availableTools := []string{"SPARQL_SELECT", "COVERAGE_CALC", "RESERVE_CALC", "FRAUD_SCORER"}
	for _, tool := range req.RequiredTools {
		if !gpe.isToolAvailable(tool, availableTools) {
			errors = append(errors, ValidationError{
				Code:     "INVALID_TOOL",
				Message:  fmt.Sprintf("Tool '%s' is not available. Available tools: %v", tool, availableTools),
				Severity: "ERROR",
				Field:    "requiredTools",
			})
		}
	}

	return errors
}

// executeSPARQLQueries retrieves facts from the knowledge graph
func (gpe *GroundedPromptEngine) executeSPARQLQueries(req GroundedRequest) ([]FactReference, error) {
	var facts []FactReference

	// Simulate SPARQL queries based on context
	// In a real implementation, this would use a proper SPARQL engine

	if claimData, exists := req.GraphContext["claim"]; exists {
		if claimMap, ok := claimData.(map[string]interface{}); ok {
			// Extract claim facts
			if claimType, exists := claimMap["claimType"]; exists {
				facts = append(facts, FactReference{
					IRI:        "autoins:claim_001#claimType",
					Value:      fmt.Sprintf("%v", claimType),
					Property:   "autoins:claimType",
					Source:     "SPARQL_SELECT",
					Confidence: 1.0,
				})
			}

			if claimAmount, exists := claimMap["claimAmount"]; exists {
				facts = append(facts, FactReference{
					IRI:        "autoins:claim_001#claimAmount",
					Value:      fmt.Sprintf("%v", claimAmount),
					Property:   "autoins:claimAmount",
					Source:     "SPARQL_SELECT",
					Confidence: 1.0,
				})
			}
		}
	}

	if policyData, exists := req.GraphContext["policy"]; exists {
		if policyMap, ok := policyData.(map[string]interface{}); ok {
			// Extract policy facts
			if coverageLimit, exists := policyMap["coverageLimit"]; exists {
				facts = append(facts, FactReference{
					IRI:        "autoins:policy_001#coverageLimit",
					Value:      fmt.Sprintf("%v", coverageLimit),
					Property:   "autoins:coverageLimit",
					Source:     "SPARQL_SELECT",
					Confidence: 1.0,
				})
			}
		}
	}

	return facts, nil
}

// performCalculations executes the required calculation tools
func (gpe *GroundedPromptEngine) performCalculations(req GroundedRequest, facts []FactReference) []CalculationRef {
	var calculations []CalculationRef

	for _, tool := range req.RequiredTools {
		switch tool {
		case "RESERVE_CALC":
			calc := gpe.performReserveCalculation(req, facts)
			if calc != nil {
				calculations = append(calculations, *calc)
			}

		case "FRAUD_SCORER":
			calc := gpe.performFraudScoring(req, facts)
			if calc != nil {
				calculations = append(calculations, *calc)
			}

		case "COVERAGE_CALC":
			calc := gpe.performCoverageCalculation(req, facts)
			if calc != nil {
				calculations = append(calculations, *calc)
			}
		}
	}

	return calculations
}

// performReserveCalculation executes reserve calculation using the reserve service
func (gpe *GroundedPromptEngine) performReserveCalculation(req GroundedRequest, facts []FactReference) *CalculationRef {
	// Extract claim data for reserve calculation
	claimData := reserve.ClaimData{}

	// Get loss type from facts
	for _, fact := range facts {
		if fact.Property == "autoins:claimType" {
			claimData.LossType = fact.Value
		}
	}

	// Get vehicle ACV from context
	if vehicleData, exists := req.GraphContext["vehicle"]; exists {
		if vehicleMap, ok := vehicleData.(map[string]interface{}); ok {
			if acv, exists := vehicleMap["actualCashValue"]; exists {
				if acvFloat, ok := acv.(float64); ok {
					claimData.VehicleACV = acvFloat
				}
			}
		}
	}

	// Get fraud indicators from context
	if fraudData, exists := req.GraphContext["fraudIndicators"]; exists {
		if fraudMap, ok := fraudData.(map[string]interface{}); ok {
			claimData.HasFraudSignals = gpe.getBoolValue(fraudMap, "hasFraudSignals")
			claimData.LiabilityUncertain = gpe.getBoolValue(fraudMap, "liabilityUncertain")
			claimData.PartsBackorder = gpe.getBoolValue(fraudMap, "partsBackorder")
		}
	}

	// Perform calculation
	result := gpe.reserveCalc.CalculateReserve(claimData)

	return &CalculationRef{
		Type: "RESERVE_CALC",
		Input: map[string]interface{}{
			"lossType":           claimData.LossType,
			"vehicleACV":         claimData.VehicleACV,
			"hasFraudSignals":    claimData.HasFraudSignals,
			"liabilityUncertain": claimData.LiabilityUncertain,
			"partsBackorder":     claimData.PartsBackorder,
		},
		Result:     result,
		Confidence: 0.95,
		Timestamp:  time.Now(),
	}
}

// performFraudScoring executes fraud risk assessment
func (gpe *GroundedPromptEngine) performFraudScoring(req GroundedRequest, facts []FactReference) *CalculationRef {
	// Extract fraud indicators from context
	fraudData := make(map[string]interface{})

	if indicators, exists := req.GraphContext["fraudIndicators"]; exists {
		if indicatorMap, ok := indicators.(map[string]interface{}); ok {
			fraudData = indicatorMap
		}
	}

	// Perform fraud assessment
	assessment := gpe.fraudScorer.AssessFraudRisk(fraudData)

	return &CalculationRef{
		Type:       "FRAUD_SCORER",
		Input:      fraudData,
		Result:     assessment,
		Confidence: 0.85,
		Timestamp:  time.Now(),
	}
}

// performCoverageCalculation calculates policy coverage amounts
func (gpe *GroundedPromptEngine) performCoverageCalculation(req GroundedRequest, facts []FactReference) *CalculationRef {
	// Extract policy and claim data
	var claimAmount, coverageLimit float64

	for _, fact := range facts {
		if fact.Property == "autoins:claimAmount" {
			fmt.Sscanf(fact.Value, "%f", &claimAmount)
		} else if fact.Property == "autoins:coverageLimit" {
			fmt.Sscanf(fact.Value, "%f", &coverageLimit)
		}
	}

	// Calculate coverage
	coveredAmount := claimAmount
	if claimAmount > coverageLimit {
		coveredAmount = coverageLimit
	}

	result := map[string]interface{}{
		"claimAmount":    claimAmount,
		"coverageLimit":  coverageLimit,
		"coveredAmount":  coveredAmount,
		"isFullyCovered": claimAmount <= coverageLimit,
		"shortfall":      math.Max(0, claimAmount-coverageLimit),
	}

	return &CalculationRef{
		Type: "COVERAGE_CALC",
		Input: map[string]interface{}{
			"claimAmount":   claimAmount,
			"coverageLimit": coverageLimit,
		},
		Result:     result,
		Confidence: 1.0,
		Timestamp:  time.Now(),
	}
}

// checkMissingInformation identifies critical missing facts
func (gpe *GroundedPromptEngine) checkMissingInformation(req GroundedRequest, facts []FactReference) []string {
	var missing []string

	// Check for required claim information
	hasClaimType := false
	hasClaimAmount := false

	for _, fact := range facts {
		if fact.Property == "autoins:claimType" {
			hasClaimType = true
		}
		if fact.Property == "autoins:claimAmount" {
			hasClaimAmount = true
		}
	}

	if !hasClaimType {
		missing = append(missing, "claim type")
	}
	if !hasClaimAmount {
		missing = append(missing, "claim amount")
	}

	// Check for policy information if coverage calculation is required
	for _, tool := range req.RequiredTools {
		if tool == "COVERAGE_CALC" {
			hasCoverageLimit := false
			for _, fact := range facts {
				if fact.Property == "autoins:coverageLimit" {
					hasCoverageLimit = true
					break
				}
			}
			if !hasCoverageLimit {
				missing = append(missing, "policy coverage limit")
			}
		}
	}

	return missing
}

// generateFollowUpQuestion creates a targeted follow-up question
func (gpe *GroundedPromptEngine) generateFollowUpQuestion(missingInfo []string) string {
	if len(missingInfo) == 1 {
		return fmt.Sprintf("I need additional information to process your request. Could you please provide the %s?", missingInfo[0])
	}

	return fmt.Sprintf("I need additional information to process your request. Could you please provide the following: %s?", strings.Join(missingInfo, ", "))
}

// generateGroundedAnswer creates the final grounded response
func (gpe *GroundedPromptEngine) generateGroundedAnswer(req GroundedRequest, facts []FactReference, calculations []CalculationRef) string {
	var answer strings.Builder

	// Start with the main response
	answer.WriteString("Based on the facts retrieved from the knowledge graph and calculations performed:\n\n")

	// Summarize key facts
	if len(facts) > 0 {
		answer.WriteString("**Facts Used:**\n")
		for _, fact := range facts {
			answer.WriteString(fmt.Sprintf("- %s: %s (IRI: %s)\n", fact.Property, fact.Value, fact.IRI))
		}
		answer.WriteString("\n")
	}

	// Summarize calculations
	if len(calculations) > 0 {
		answer.WriteString("**Calculations Performed:**\n")
		for _, calc := range calculations {
			switch calc.Type {
			case "RESERVE_CALC":
				if result, ok := calc.Result.(reserve.ReserveResult); ok {
					answer.WriteString(fmt.Sprintf("- Reserve Calculation: %s (Band: %s)\n", result.Breakdown, result.ReserveBand))
				}
			case "FRAUD_SCORER":
				if result, ok := calc.Result.(reserve.FraudAssessment); ok {
					answer.WriteString(fmt.Sprintf("- Fraud Assessment: %s risk (Score: %.2f)\n", result.RiskLevel, result.RiskScore))
				}
			case "COVERAGE_CALC":
				if result, ok := calc.Result.(map[string]interface{}); ok {
					answer.WriteString(fmt.Sprintf("- Coverage Analysis: £%.2f covered of £%.2f claimed\n",
						result["coveredAmount"], result["claimAmount"]))
				}
			}
		}
		answer.WriteString("\n")
	}

	// Add specific analysis based on the query
	answer.WriteString("**Analysis:**\n")
	answer.WriteString("This response is grounded in the insurance ontology graph and verified calculations. ")
	answer.WriteString("All facts and calculations are traceable to their sources as listed above.")

	return answer.String()
}

// calculateConfidence determines the overall confidence level
func (gpe *GroundedPromptEngine) calculateConfidence(facts []FactReference, calculations []CalculationRef) float64 {
	if len(facts) == 0 {
		return 0.0
	}

	// Calculate average confidence from facts
	var factConfidence float64
	for _, fact := range facts {
		factConfidence += fact.Confidence
	}
	factConfidence /= float64(len(facts))

	// Calculate average confidence from calculations
	var calcConfidence float64 = 1.0
	if len(calculations) > 0 {
		calcConfidence = 0.0
		for _, calc := range calculations {
			calcConfidence += calc.Confidence
		}
		calcConfidence /= float64(len(calculations))
	}

	// Return weighted average
	return (factConfidence*0.6 + calcConfidence*0.4)
}

// Helper functions
func (gpe *GroundedPromptEngine) isToolAvailable(tool string, available []string) bool {
	for _, t := range available {
		if t == tool {
			return true
		}
	}
	return false
}

func (gpe *GroundedPromptEngine) getBoolValue(data map[string]interface{}, key string) bool {
	if val, exists := data[key]; exists {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return false
}

// initializeOntologyGraph creates a simplified graph representation
func initializeOntologyGraph() map[string]interface{} {
	return map[string]interface{}{
		"classes": map[string]interface{}{
			"autoins:Claim":    map[string]string{"label": "Insurance Claim"},
			"autoins:Policy":   map[string]string{"label": "Insurance Policy"},
			"autoins:Vehicle":  map[string]string{"label": "Vehicle"},
			"autoins:Accident": map[string]string{"label": "Accident"},
		},
		"properties": map[string]interface{}{
			"autoins:claimType":     map[string]string{"domain": "autoins:Claim", "range": "xsd:string"},
			"autoins:claimAmount":   map[string]string{"domain": "autoins:Claim", "range": "xsd:decimal"},
			"autoins:coverageLimit": map[string]string{"domain": "autoins:Policy", "range": "xsd:decimal"},
		},
	}
}
