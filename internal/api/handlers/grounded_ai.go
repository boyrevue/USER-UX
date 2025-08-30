package handlers

import (
	"encoding/json"
	"net/http"

	"client-ux/internal/services/grounded"
	"client-ux/internal/services/reserve"
)

// GroundedAIHandler handles grounded AI processing requests
type GroundedAIHandler struct {
	promptEngine *grounded.GroundedPromptEngine
	reserveCalc  *reserve.ReserveCalculator
	fraudScorer  *reserve.FraudScorer
}

// NewGroundedAIHandler creates a new grounded AI handler
func NewGroundedAIHandler() *GroundedAIHandler {
	return &GroundedAIHandler{
		promptEngine: grounded.NewGroundedPromptEngine(),
		reserveCalc:  reserve.NewReserveCalculator(),
		fraudScorer:  reserve.NewFraudScorer(),
	}
}

// ProcessGroundedQuery handles grounded AI query processing
func (h *GroundedAIHandler) ProcessGroundedQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req grounded.GroundedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Process the grounded query
	response := h.promptEngine.ProcessGroundedQuery(req)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// CalculateReserve handles reserve calculation requests
func (h *GroundedAIHandler) CalculateReserve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var claimData reserve.ClaimData
	if err := json.NewDecoder(r.Body).Decode(&claimData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Calculate reserve
	result := h.reserveCalc.CalculateReserve(claimData)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// AssessFraud handles fraud risk assessment requests
func (h *GroundedAIHandler) AssessFraud(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var fraudData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&fraudData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Assess fraud risk
	assessment := h.fraudScorer.AssessFraudRisk(fraudData)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(assessment); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// ValidateFNOL handles FNOL completeness validation
func (h *GroundedAIHandler) ValidateFNOL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var claimData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&claimData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate FNOL completeness
	validation := h.validateFNOLCompleteness(claimData)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(validation); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// FNOLValidationResult represents FNOL validation results
type FNOLValidationResult struct {
	IsComplete       bool     `json:"isComplete"`
	MissingFields    []string `json:"missingFields"`
	ValidationErrors []string `json:"validationErrors"`
	Completeness     float64  `json:"completeness"`
	Recommendations  []string `json:"recommendations"`
}

// validateFNOLCompleteness implements FNOL completeness validation from the ontology
func (h *GroundedAIHandler) validateFNOLCompleteness(claimData map[string]interface{}) FNOLValidationResult {
	var missingFields []string
	var validationErrors []string
	var recommendations []string

	// Check required FNOL fields based on SHACL shapes
	requiredFields := map[string]string{
		"hasIncident":      "Incident/accident record",
		"relatesToPolicy":  "Policy reference",
		"claimDate":        "Claim date",
		"incidentDate":     "Incident date",
		"claimType":        "Claim type",
		"claimDescription": "Claim description",
	}

	for field, description := range requiredFields {
		if _, exists := claimData[field]; !exists {
			missingFields = append(missingFields, field)
			validationErrors = append(validationErrors,
				"🚨 FNOL ERROR: "+description+" is required for First Notice of Loss processing")
		}
	}

	// Validate claim description length (minimum 10 characters)
	if desc, exists := claimData["claimDescription"]; exists {
		if descStr, ok := desc.(string); ok && len(descStr) < 10 {
			validationErrors = append(validationErrors,
				"📝 FNOL ERROR: Detailed claim description (minimum 10 characters) is required for processing")
			recommendations = append(recommendations,
				"Please provide a more detailed description of the incident")
		}
	}

	// Check for incident linkage
	if _, hasIncident := claimData["hasIncident"]; !hasIncident {
		recommendations = append(recommendations,
			"Link this claim to the underlying incident/accident record")
	}

	// Check for policy linkage
	if _, hasPolicy := claimData["relatesToPolicy"]; !hasPolicy {
		recommendations = append(recommendations,
			"Verify and link the claim to the appropriate insurance policy")
	}

	// Calculate completeness percentage
	totalRequired := len(requiredFields)
	completed := totalRequired - len(missingFields)
	completeness := float64(completed) / float64(totalRequired) * 100

	isComplete := len(missingFields) == 0 && len(validationErrors) == 0

	return FNOLValidationResult{
		IsComplete:       isComplete,
		MissingFields:    missingFields,
		ValidationErrors: validationErrors,
		Completeness:     completeness,
		Recommendations:  recommendations,
	}
}

// GetSystemPrompt returns the grounded AI system prompt from the ontology
func (h *GroundedAIHandler) GetSystemPrompt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get language preference from query params
	language := r.URL.Query().Get("lang")
	if language == "" {
		language = "en"
	}

	// Return the grounded claims processing prompt from the ontology
	prompt := map[string]interface{}{
		"id":       "autoins:GroundedClaimsPrompt",
		"label":    "Grounded Claims Processing Prompt",
		"language": language,
		"promptText": `You are a grounded AI assistant for insurance claims processing. You must answer only from the knowledge graph and attached calculations.

MANDATORY TOOLS AVAILABLE:
- SPARQL_SELECT: Query the insurance ontology graph
- COVERAGE_CALC: Calculate policy coverage amounts
- RESERVE_CALC: Calculate claim reserves using severity tables
- FRAUD_SCORER: Assess fraud risk indicators

GROUNDING REQUIREMENTS:
1. You MUST use SPARQL_SELECT to retrieve facts from the graph before answering
2. You MUST use appropriate calculation tools for any numerical assessments
3. If a fact is missing from the graph, ask ONE targeted follow-up question, then re-check
4. Your final answer MUST list the IRIs of all facts used in your reasoning
5. You MUST cite specific calculation results with their input parameters

RESPONSE FORMAT:
- Answer: [Your grounded response]
- Facts Used: [List of IRIs from graph queries]
- Calculations: [List of calculations performed with results]
- Follow-up Needed: [Yes/No - if missing critical information]

FORBIDDEN ACTIONS:
- Never make assumptions not supported by graph data
- Never provide estimates without using RESERVE_CALC tool
- Never assess coverage without using COVERAGE_CALC tool
- Never ignore fraud indicators without using FRAUD_SCORER tool`,
		"complianceLevel":       "GDPR_Strict",
		"allowedDataCategories": []string{"ClaimsData", "PolicyData"},
		"requiredTools":         []string{"SPARQL_SELECT", "COVERAGE_CALC", "RESERVE_CALC", "FRAUD_SCORER"},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(prompt); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
