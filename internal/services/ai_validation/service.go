package ai_validation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type ValidationRequest struct {
	FieldName        string `json:"fieldName"`
	UserInput        string `json:"userInput"`
	ValidationPrompt string `json:"validationPrompt"`
}

type ValidationResponse struct {
	IsValid      bool   `json:"isValid"`
	Message      string `json:"message"`
	Suggestions  string `json:"suggestions,omitempty"`
	RequiredInfo string `json:"requiredInfo,omitempty"`
}

type Service struct {
	// Could integrate with OpenAI, Claude, or other AI services
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ValidateUserInput(req ValidationRequest) (*ValidationResponse, error) {
	// For now, implement basic validation rules for insurance-specific fields
	// This can be enhanced with actual AI integration later

	userInput := strings.TrimSpace(strings.ToLower(req.UserInput))

	// Check for obviously invalid responses
	invalidResponses := []string{
		"fish and chips", "n/a", "not applicable", "none", "nothing",
		"personal reasons", "private", "don't want to say", "refuse to answer",
		"test", "testing", "asdf", "qwerty", "123", "abc",
	}

	for _, invalid := range invalidResponses {
		if strings.Contains(userInput, invalid) {
			return &ValidationResponse{
				IsValid:      false,
				Message:      "Please provide specific details relevant to your insurance application.",
				RequiredInfo: s.getRequiredInfoForField(req.FieldName),
			}, nil
		}
	}

	// Check minimum word count for meaningful responses (temporarily allow 1 word)
	wordCount := len(strings.Fields(strings.TrimSpace(req.UserInput)))
	if wordCount < 1 {
		return &ValidationResponse{
			IsValid:      false,
			Message:      "Please provide at least 1 meaningful word describing your situation.",
			RequiredInfo: s.getRequiredInfoForField(req.FieldName),
		}, nil
	}

	// Field-specific validation
	switch req.FieldName {
	case "revocationOtherReason":
		return s.validateRevocationReason(req.UserInput)
	case "licenceRefusalOtherReason":
		return s.validateRefusalReason(req.UserInput)
	case "licenceRestrictionOtherDetails":
		return s.validateRestrictionDetails(req.UserInput)
	case "medicalConditionOtherDetails":
		return s.validateMedicalConditionDetails(req.UserInput)
	case "claimsAccidentsOtherDetails":
		return s.validateClaimsAccidentsDetails(req.UserInput)
	case "disqualificationOtherReason":
		return s.validateDisqualificationOther(req.UserInput)
	}

	// Default validation passed
	return &ValidationResponse{
		IsValid: true,
		Message: "Thank you for providing this information.",
	}, nil
}

func (s *Service) validateRevocationReason(input string) (*ValidationResponse, error) {
	inputLower := strings.ToLower(input)
	inputOriginal := strings.TrimSpace(input)

	// Check for specific offence codes (high value indicators)
	offenceCodes := []string{
		"dr10", "dr20", "dr30", "dr40", "dr50", "dr60", "dr70", "dr80", // Drink/drug driving
		"in10",         // No insurance
		"sp30", "sp50", // Speeding (for new driver revocations)
	}

	hasOffenceCode := false
	for _, code := range offenceCodes {
		if strings.Contains(inputLower, code) {
			hasOffenceCode = true
			break
		}
	}

	// Look for specific revocation reasons and circumstances
	specificReasons := []string{
		"drink driving", "drug driving", "drink-driving", "drug-driving",
		"no insurance", "without insurance", "uninsured",
		"medical condition", "epilepsy", "diabetes", "vision", "eyesight",
		"new driver", "6 points", "12 points", "totting up",
		"fraud", "false information", "false documents",
		"failed medical", "failed eyesight", "court order",
		"dvla investigation", "compliance", "surrender licence",
		"foreign licence", "age related", "dementia",
	}

	hasSpecificReason := false
	for _, reason := range specificReasons {
		if strings.Contains(inputLower, reason) {
			hasSpecificReason = true
			break
		}
	}

	// Look for date indicators (year or specific date)
	datePattern := regexp.MustCompile(`(19|20)\d{2}|january|february|march|april|may|june|july|august|september|october|november|december`)
	hasDate := datePattern.MatchString(inputLower)

	// Look for current status indicators
	statusTerms := []string{
		"reinstated", "regained", "got back", "licence back", "retook", "re-passed",
		"still revoked", "currently revoked", "medical clearance", "extended test",
	}

	hasStatus := false
	for _, status := range statusTerms {
		if strings.Contains(inputLower, status) {
			hasStatus = true
			break
		}
	}

	// Check for vague or nonsensical responses (immediate rejection)
	vaguePhrases := []string{
		"fish and chips", "personal reasons", "don't want to say", "private matter",
		"no reason", "just because", "dvla decision", "administrative reasons",
		"complications", "licence issues", "problems with dvla", "legal issues",
		"prefer not to discuss", "random decision", "no reason given",
	}

	for _, phrase := range vaguePhrases {
		if strings.Contains(inputLower, phrase) {
			return &ValidationResponse{
				IsValid:      false,
				Message:      "❌ INVALID: Vague responses are not acceptable for licence revocation. Insurance underwriters require specific details about this serious matter.",
				RequiredInfo: s.getRequiredInfoForField("revocationOtherReason"),
			}, nil
		}
	}

	// Detect high-risk scenarios that need extra validation
	highRiskTerms := []string{
		"drink", "drug", "alcohol", "dr10", "dr20", "dr30",
		"no insurance", "in10", "fraud", "false",
	}

	isHighRisk := false
	for _, term := range highRiskTerms {
		if strings.Contains(inputLower, term) {
			isHighRisk = true
			break
		}
	}

	// Scoring system for response quality
	score := 0
	if hasOffenceCode {
		score += 3 // Official codes are very valuable
	}
	if hasSpecificReason {
		score += 2
	}
	if hasDate {
		score += 2
	}
	if hasStatus {
		score += 1
	}
	if len(inputOriginal) > 50 {
		score += 1 // Detailed responses
	}

	// Validation decision based on score and content
	if score >= 4 || (hasOffenceCode && hasSpecificReason) {
		message := "✅ ACCEPTED: Thank you for providing specific details about your licence revocation."
		if isHighRisk {
			message += " This is a serious matter that will significantly impact your insurance premium."
		}
		return &ValidationResponse{
			IsValid: true,
			Message: message,
		}, nil
	}

	if score >= 2 && hasSpecificReason {
		return &ValidationResponse{
			IsValid:      false,
			Message:      "⚠️ NEEDS MORE DETAIL: Please provide the approximate date of revocation and your current licence status (reinstated/still revoked).",
			RequiredInfo: s.getRequiredInfoForField("revocationOtherReason"),
		}, nil
	}

	// Insufficient information
	return &ValidationResponse{
		IsValid:      false,
		Message:      "❌ INSUFFICIENT: Please provide specific details including the DVLA reason for revocation, approximate date, and current licence status.",
		RequiredInfo: s.getRequiredInfoForField("revocationOtherReason"),
	}, nil
}

func (s *Service) validateRefusalReason(input string) (*ValidationResponse, error) {
	input = strings.ToLower(input)

	// Look for specific refusal-related terms
	relevantTerms := []string{
		"medical", "eyesight", "test", "theory", "practical", "age",
		"residency", "conviction", "disqualification", "information",
		"dvla", "application", "refused", "failed", "disability",
	}

	hasRelevantTerm := false
	for _, term := range relevantTerms {
		if strings.Contains(input, term) {
			hasRelevantTerm = true
			break
		}
	}

	if !hasRelevantTerm {
		return &ValidationResponse{
			IsValid:      false,
			Message:      "Please provide specific details about why your licence application was refused.",
			RequiredInfo: "We need to know: 1) The specific reason DVLA gave for refusal, 2) What stage of the application process, 3) Current status of your licence.",
		}, nil
	}

	return &ValidationResponse{
		IsValid: true,
		Message: "Thank you for providing details about your licence refusal.",
	}, nil
}

// validateDisqualificationOther validates free-text "Other" disqualification reasons
func (s *Service) validateDisqualificationOther(input string) (*ValidationResponse, error) {
	inputLower := strings.ToLower(strings.TrimSpace(input))

	// Immediate rejection for obviously vague/joke content
	vague := []string{"ice cream", "test", "asdf", "n/a", "none", "no idea", "prefer not"}
	for _, v := range vague {
		if strings.Contains(inputLower, v) {
			return &ValidationResponse{
				IsValid:      false,
				Message:      "❌ INVALID: Please provide the exact court reason and, if known, the offence code (e.g., BA10, IN10, CD10, AC10, DR10, DR80, DD40, DD80, MS90, TT99).",
				RequiredInfo: "Required: specific offence/context, official code if known, ban length/dates.",
			}, nil
		}
	}

	// Look for recognised offence codes and keywords
	codes := []string{"ba10", "in10", "cd10", "cd20", "cd30", "ac10", "dr10", "dr30", "dr80", "dd40", "dd80", "ms90", "tt99"}
	keywords := []string{"driving while disqualified", "without insurance", "careless driving", "failing to stop", "fail to stop", "dangerous driving", "drink driving", "drug driving", "motor racing"}

	hasCode := false
	for _, c := range codes {
		if strings.Contains(inputLower, c) {
			hasCode = true
			break
		}
	}

	hasKeyword := false
	for _, k := range keywords {
		if strings.Contains(inputLower, k) {
			hasKeyword = true
			break
		}
	}

	words := len(strings.Fields(inputLower))
	if hasCode || (hasKeyword && words >= 6) {
		return &ValidationResponse{IsValid: true, Message: "✅ ACCEPTED: Thank you for providing a specific disqualification reason."}, nil
	}

	return &ValidationResponse{
		IsValid:      false,
		Message:      "⚠️ NEEDS MORE DETAIL: Specify the offence and include the official code if known (BA10, IN10, CD10/CD20/CD30, AC10, DR10/DR80, DD40/DD80, MS90, TT99).",
		RequiredInfo: "Provide offence, context, and ensure dates/duration are completed in the date fields.",
	}, nil
}

func (s *Service) validateRestrictionDetails(input string) (*ValidationResponse, error) {
	input = strings.ToLower(input)

	// Look for restriction code patterns or relevant terms
	relevantTerms := []string{
		"code", "restriction", "01", "02", "03", "78", "79", "40",
		"glasses", "contacts", "hearing", "automatic", "medical",
		"modified", "steering", "brakes", "transmission", "dvla",
	}

	hasRelevantTerm := false
	for _, term := range relevantTerms {
		if strings.Contains(input, term) {
			hasRelevantTerm = true
			break
		}
	}

	if !hasRelevantTerm {
		return &ValidationResponse{
			IsValid:      false,
			Message:      "Please provide specific DVLA restriction codes and what they mean.",
			RequiredInfo: "We need to know: 1) The exact restriction code numbers (e.g., 01, 78), 2) What the restrictions mean for your driving, 3) Any equipment or modifications required.",
		}, nil
	}

	return &ValidationResponse{
		IsValid: true,
		Message: "Thank you for providing details about your licence restrictions.",
	}, nil
}

func (s *Service) validateMedicalConditionDetails(input string) (*ValidationResponse, error) {
	input = strings.ToLower(input)

	// Look for medical condition-related terms (based on AI_Medical_Conditions_Guide.ttl)
	relevantTerms := []string{
		"condition", "diagnosis", "disease", "disorder", "syndrome", "medical",
		"doctor", "specialist", "consultant", "hospital", "treatment", "medication",
		"symptoms", "affects", "driving", "dvla", "declared", "notified",
		// Specific conditions from ontology
		"diabetes", "epilepsy", "heart", "vision", "hearing", "mobility", "mental",
		"narcolepsy", "parkinson", "multiple sclerosis", "stroke", "tia",
		"chronic", "disability", "impairment", "restriction", "adaptation",
		// Medication effects
		"sedating", "drowsiness", "alertness", "concentration", "cognition",
	}

	hasRelevantTerm := false
	for _, term := range relevantTerms {
		if strings.Contains(input, term) {
			hasRelevantTerm = true
			break
		}
	}

	// Check for medical condition indicators
	medicalIndicators := []string{
		"diagnosed", "suffer", "have", "experience", "condition", "problem",
		"issue", "difficulty", "unable", "limited", "restricted", "affects",
	}

	hasMedicalIndicator := false
	for _, indicator := range medicalIndicators {
		if strings.Contains(input, indicator) {
			hasMedicalIndicator = true
			break
		}
	}

	if !hasRelevantTerm || !hasMedicalIndicator {
		return &ValidationResponse{
			IsValid:      false,
			Message:      "Please provide specific medical information relevant to driving and insurance.",
			RequiredInfo: "We need to know: 1) The specific medical condition or diagnosis, 2) Whether declared to DVLA (legally required if affects driving), 3) How condition affects driving ability, 4) Current treatment/medication status, 5) Any licence restrictions or medical reviews required. Key question: Could this condition cause sudden loss of control while driving?",
		}, nil
	}

	// Check for mandatory DVLA declaration conditions (based on AI_Medical_Conditions_Guide.ttl)
	mandatoryDeclarationConditions := []string{
		"epilepsy", "seizure", "blackout", "fit", "convulsion",
		"insulin", "diabetes", "diabetic", "hypoglycaemic",
		"heart attack", "cardiac", "pacemaker", "defibrillator", "arrhythmia",
		"stroke", "tia", "transient ischaemic",
		"sleep apnea", "narcolepsy", "excessive sleepiness",
		"substance misuse", "alcohol dependency", "drug dependency",
	}

	hasMandatoryCondition := false
	for _, condition := range mandatoryDeclarationConditions {
		if strings.Contains(input, condition) {
			hasMandatoryCondition = true
			break
		}
	}

	// Additional validation for mandatory declaration conditions
	if hasMandatoryCondition {
		dvlaTerms := []string{"dvla", "declared", "notified", "medical review", "clearance", "licence valid"}
		hasDVLAMention := false
		for _, term := range dvlaTerms {
			if strings.Contains(input, term) {
				hasDVLAMention = true
				break
			}
		}

		if !hasDVLAMention {
			return &ValidationResponse{
				IsValid:      false,
				Message:      "This condition requires mandatory DVLA declaration. Please confirm your DVLA status.",
				RequiredInfo: "CONDITIONS REQUIRING MANDATORY DVLA DECLARATION: Epilepsy/seizures, insulin-treated diabetes, heart conditions causing sudden disablement, stroke/TIA, severe mental health episodes, sleep disorders causing excessive sleepiness, substance misuse dependencies. Please provide: 1) DVLA declaration status, 2) Current licence validity, 3) Any restrictions imposed.",
			}, nil
		}
	}

	return &ValidationResponse{
		IsValid: true,
		Message: "Thank you for providing comprehensive medical condition information.",
	}, nil
}

func (s *Service) validateClaimsAccidentsDetails(input string) (*ValidationResponse, error) {
	input = strings.ToLower(input)

	// Look for claims/accidents-related terms (based on AI_Claims_Accidents_Guide.ttl)
	relevantTerms := []string{
		"damage", "accident", "incident", "collision", "hit", "crash",
		"claim", "repair", "dent", "scratch", "broken", "cracked",
		// Specific scenario terms
		"falling", "tree", "branch", "debris", "rock", "landslide",
		"keyed", "scratched", "trolley", "garage", "lift", "fire",
		"theft", "stolen", "vandalism", "malicious", "brake", "engine",
		"electrical", "mechanical", "failure", "spillage", "chemical",
		// Location and time indicators
		"parked", "driving", "location", "address", "time", "date",
		"street", "car park", "road", "motorway",
	}

	hasRelevantTerm := false
	for _, term := range relevantTerms {
		if strings.Contains(input, term) {
			hasRelevantTerm = true
			break
		}
	}

	// Check for incident description indicators
	incidentIndicators := []string{
		"happened", "occurred", "caused", "resulted", "discovered",
		"found", "noticed", "damaged", "broken", "hit", "struck",
		"fell", "dropped", "spilled", "leaked", "failed",
	}

	hasIncidentIndicator := false
	for _, indicator := range incidentIndicators {
		if strings.Contains(input, indicator) {
			hasIncidentIndicator = true
			break
		}
	}

	if !hasRelevantTerm || !hasIncidentIndicator {
		return &ValidationResponse{
			IsValid:      false,
			Message:      "Please provide specific details about the accident or claim incident.",
			RequiredInfo: "We need to know: 1) What exactly happened (specific incident description), 2) Where and when it occurred (exact location and time), 3) What damage was caused (precise description), 4) Any third-party involvement, 5) Whether reported to police if criminal/malicious damage.",
		}, nil
	}

	// Check for vague or speculative language that should be avoided
	vagueTerms := []string{
		"something happened", "accident occurred", "damage found", "hit something",
		"car got damaged", "incident happened", "mystery damage", "someone hit",
		"i don't know", "not sure", "might have", "could be", "probably",
		"i think", "maybe", "possibly",
	}

	hasVagueTerm := false
	for _, term := range vagueTerms {
		if strings.Contains(input, term) {
			hasVagueTerm = true
			break
		}
	}

	if hasVagueTerm {
		return &ValidationResponse{
			IsValid:      false,
			Message:      "Please avoid vague descriptions. Provide specific, factual details about what happened.",
			RequiredInfo: "AVOID: 'something happened', 'mystery damage', 'I think...', 'probably...'. PROVIDE: Specific facts - what, where, when, how much damage, any third parties involved. Example: 'Falling tree branch during storm damaged roof while parked at [address]'",
		}, nil
	}

	// Check for fault-admitting language that should be avoided
	faultAdmittingTerms := []string{
		"i crashed into", "i hit", "my fault", "i caused", "i was speeding",
		"i wasn't paying attention", "i made a mistake",
	}

	hasFaultAdmission := false
	for _, term := range faultAdmittingTerms {
		if strings.Contains(input, term) {
			hasFaultAdmission = true
			break
		}
	}

	if hasFaultAdmission {
		return &ValidationResponse{
			IsValid:      false,
			Message:      "Avoid admitting fault unnecessarily. Stick to factual descriptions and let the insurance company determine fault.",
			RequiredInfo: "Instead of 'I crashed into a wall', say 'Vehicle collided with wall on [Street Name]'. Focus on facts: what happened, where, when, what damage occurred.",
		}, nil
	}

	// Check for criminal/malicious damage that should mention police
	criminalTerms := []string{
		"keyed", "vandalism", "malicious", "deliberate", "theft", "stolen",
		"break in", "criminal", "sabotage", "tampering",
	}

	hasCriminalTerm := false
	for _, term := range criminalTerms {
		if strings.Contains(input, term) {
			hasCriminalTerm = true
			break
		}
	}

	if hasCriminalTerm {
		policeTerms := []string{"police", "reported", "crime reference", "officer", "statement"}
		hasPoliceReference := false
		for _, term := range policeTerms {
			if strings.Contains(input, term) {
				hasPoliceReference = true
				break
			}
		}

		if !hasPoliceReference {
			return &ValidationResponse{
				IsValid:      false,
				Message:      "Criminal or malicious damage should be reported to police. Please confirm police involvement.",
				RequiredInfo: "For theft, vandalism, or malicious damage, please provide: 1) Police crime reference number, 2) When reported to police, 3) Any evidence collected, 4) Extent of damage/theft.",
			}, nil
		}
	}

	return &ValidationResponse{
		IsValid: true,
		Message: "Thank you for providing detailed accident/claim information.",
	}, nil
}

func (s *Service) getRequiredInfoForField(fieldName string) string {
	switch fieldName {
	case "revocationOtherReason":
		return `⚠️ LICENCE REVOCATION - CRITICAL INSURANCE INFORMATION REQUIRED:

MANDATORY DETAILS NEEDED:
1. ✅ EXACT DVLA reason for revocation (not vague terms like 'medical reasons')
2. ✅ Approximate date/year of revocation
3. ✅ Current licence status (revoked/reinstated/restricted)
4. ✅ Official offence codes if known (DR10, IN10, SP30, etc.)
5. ✅ Reinstatement process completed (if applicable)

EXAMPLES OF ACCEPTABLE RESPONSES:
• "Licence revoked 2019 for drink-driving conviction (DR10). Disqualified 24 months. Completed medical assessment and regained licence 2021."
• "Medical revocation 2020 due to epilepsy. Reinstated 2022 after 2 years seizure-free with DVLA clearance."
• "Revoked under New Driver Act 2018 for 6 points (two SP30 speeding). Re-passed theory and practical tests 2019."

⚠️ WARNING: Licence revocation is EXTREMELY serious for insurance. Vague responses like 'DVLA decision', 'personal reasons', or 'medical issues' are not acceptable. Insurers WILL check DVLA records and non-disclosure is fraud.`
	case "licenceRefusalOtherReason":
		return "Please provide: 1) Specific DVLA reason for refusal, 2) Stage of application when refused, 3) Current licence status, 4) Steps taken since refusal."
	case "licenceRestrictionOtherDetails":
		return "Please provide: 1) Exact DVLA restriction codes, 2) What each code means practically, 3) Any equipment/modifications required, 4) How it affects your driving."
	case "medicalConditionOtherDetails":
		return "Please provide: 1) Specific medical condition name/diagnosis, 2) Whether declared to DVLA (legally required if affects driving), 3) How condition affects driving ability, 4) Current treatment/medication status, 5) Any licence restrictions or medical reviews required. Key question: Could this condition cause sudden loss of control while driving?"
	case "claimsAccidentsOtherDetails":
		return "Please provide: 1) What exactly happened (specific incident description), 2) Where and when it occurred (exact location and time), 3) What damage was caused (precise description), 4) Any third-party involvement or damage, 5) Whether reported to police if criminal/malicious damage. Be factual, specific, and avoid admitting fault."
	default:
		return "Please provide specific, detailed information relevant to your insurance application."
	}
}

// HTTP Handler for AI validation endpoint
func (s *Service) HandleValidation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := s.ValidateUserInput(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
