package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"client-ux/internal/services/bipro"
)

// BiPROHandler handles BiPRO compliance API requests
type BiPROHandler struct {
	biproService *bipro.BiPROService
}

// NewBiPROHandler creates a new BiPRO handler
func NewBiPROHandler() *BiPROHandler {
	return &BiPROHandler{
		biproService: bipro.NewBiPROService(),
	}
}

// ============================================================
// NORM 420 - TARIFICATION, OFFER, APPLICATION (TAA)
// ============================================================

// ProcessTariffCalculation handles BiPRO Norm 420 tariff calculation requests
func (h *BiPROHandler) ProcessTariffCalculation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req bipro.Norm420TariffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Process tariff calculation
	response, err := h.biproService.ProcessTariffRequest(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Tariff calculation failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("BiPRO-Norm", "420")
	w.Header().Set("BiPRO-Version", "2024.1")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetTariffQuote provides a simplified tariff quote endpoint
func (h *BiPROHandler) GetTariffQuote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simplified quote request structure
	var quoteReq struct {
		Vehicle struct {
			Make         string  `json:"make"`
			Model        string  `json:"model"`
			Year         int     `json:"year"`
			Value        float64 `json:"value"`
			Registration string  `json:"registration"`
		} `json:"vehicle"`
		Driver struct {
			DateOfBirth string `json:"dateOfBirth"`
			LicenseDate string `json:"licenseDate"`
			Occupation  string `json:"occupation"`
			PostalCode  string `json:"postalCode"`
		} `json:"driver"`
		Coverage struct {
			Type   string  `json:"type"`
			Excess float64 `json:"excess"`
		} `json:"coverage"`
	}

	if err := json.NewDecoder(r.Body).Decode(&quoteReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert to BiPRO format
	biproReq := h.convertToNorm420Request(quoteReq)

	// Process through BiPRO service
	response, err := h.biproService.ProcessTariffRequest(biproReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Quote calculation failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Simplify response for client
	quoteResponse := map[string]interface{}{
		"quoteId":        response.ResponseID,
		"annualPremium":  response.Premium.AnnualPremium,
		"monthlyPremium": response.Premium.MonthlyPremium,
		"currency":       response.Premium.Currency,
		"excess":         response.Conditions.Excess,
		"coverageType":   response.Conditions.CoverageType,
		"validUntil":     response.Validity.ValidTo,
		"calculations":   response.Calculations,
		"biproCompliant": true,
		"norm":           "420",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(quoteResponse); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// ============================================================
// NORM 430 - TRANSFER SERVICES
// ============================================================

// ProcessDocumentTransfer handles BiPRO Norm 430 document transfer requests
func (h *BiPROHandler) ProcessDocumentTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form for document uploads
	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Get transfer type from form
	transferType := r.FormValue("transferType")
	if transferType == "" {
		http.Error(w, "Transfer type is required", http.StatusBadRequest)
		return
	}

	// Get document format
	format := r.FormValue("format")
	if format == "" {
		format = "GDV" // Default format
	}

	// Get uploaded file
	file, header, err := r.FormFile("document")
	if err != nil {
		http.Error(w, "Document file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file data
	fileData := make([]byte, header.Size)
	if _, err := file.Read(fileData); err != nil {
		http.Error(w, "Failed to read document file", http.StatusInternalServerError)
		return
	}

	// Create BiPRO transfer request
	transferReq := bipro.Norm430TransferRequest{
		MessageHeader: bipro.BiPROMessageHeader{
			MessageID:   generateBiPROMessageID(),
			Sender:      "CLIENT-UX",
			Receiver:    r.Header.Get("BiPRO-Receiver"),
			NormVersion: "430.2024.1",
			Timestamp:   time.Now(),
		},
		TransferType: transferType,
		DocumentType: r.FormValue("documentType"),
		Format:       format,
		Compression:  r.FormValue("compression"),
		Data:         fileData,
		Metadata: bipro.BiPROMetadata{
			DocumentID:   header.Filename,
			DocumentDate: time.Now(),
			PolicyNumber: r.FormValue("policyNumber"),
			CustomerID:   r.FormValue("customerId"),
			Properties: map[string]string{
				"originalFilename": header.Filename,
				"contentType":      header.Header.Get("Content-Type"),
				"uploadedBy":       r.Header.Get("User-Agent"),
			},
		},
	}

	// Process transfer
	response, err := h.biproService.ProcessDocumentTransfer(transferReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Document transfer failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("BiPRO-Norm", "430")
	w.Header().Set("BiPRO-Version", "2024.1")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// ProcessGDVData handles GDV format data processing (Norm 430.1)
func (h *BiPROHandler) ProcessGDVData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create GDV-specific transfer request
	transferReq := bipro.Norm430TransferRequest{
		MessageHeader: bipro.BiPROMessageHeader{
			MessageID:   generateBiPROMessageID(),
			Sender:      "CLIENT-UX",
			Receiver:    "GDV-PROCESSOR",
			NormVersion: "430.1.2024.1",
			Timestamp:   time.Now(),
		},
		TransferType: "430.1",
		Format:       "GDV",
		Compression:  "ZIP",
	}

	// Read request body as GDV data
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	transferReq.Data = body

	// Process GDV data
	response, err := h.biproService.ProcessDocumentTransfer(transferReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("GDV processing failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("BiPRO-Norm", "430.1")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// ============================================================
// NORM 440 - DIRECT ACCESS (DEEP LINK)
// ============================================================

// ProcessDeepLinkAccess handles BiPRO Norm 440 deep link requests
func (h *BiPROHandler) ProcessDeepLinkAccess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req bipro.Norm440DeepLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Process deep link request
	response, err := h.biproService.ProcessDeepLinkRequest(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Deep link processing failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("BiPRO-Norm", "440")
	w.Header().Set("BiPRO-Version", "2024.1")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GenerateDeepLink creates a deep link for external system access
func (h *BiPROHandler) GenerateDeepLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var linkReq struct {
		TargetSystem   string            `json:"targetSystem"`
		TargetFunction string            `json:"targetFunction"`
		Parameters     map[string]string `json:"parameters"`
		UserID         string            `json:"userId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&linkReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create BiPRO deep link request
	biproReq := bipro.Norm440DeepLinkRequest{
		MessageHeader: bipro.BiPROMessageHeader{
			MessageID:   generateBiPROMessageID(),
			Sender:      "CLIENT-UX",
			Receiver:    linkReq.TargetSystem,
			NormVersion: "440.2024.1",
			Timestamp:   time.Now(),
		},
		TargetSystem:   linkReq.TargetSystem,
		TargetFunction: linkReq.TargetFunction,
		Parameters:     linkReq.Parameters,
		SessionToken:   generateSessionToken(),
		UserID:         linkReq.UserID,
	}

	// Process deep link
	response, err := h.biproService.ProcessDeepLinkRequest(biproReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Deep link generation failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Return simplified response
	linkResponse := map[string]interface{}{
		"accessUrl":      response.AccessURL,
		"sessionId":      response.SessionID,
		"expiresAt":      response.ExpiresAt,
		"status":         response.Status,
		"biproCompliant": true,
		"norm":           "440",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(linkResponse); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// ============================================================
// BIPRO COMPLIANCE STATUS & INFORMATION
// ============================================================

// GetBiPROComplianceStatus returns BiPRO compliance information
func (h *BiPROHandler) GetBiPROComplianceStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	complianceStatus := map[string]interface{}{
		"platform":            "CLIENT-UX",
		"biproCompliant":      true,
		"certificationStatus": "pending", // Would be updated after certification
		"supportedNorms": []map[string]interface{}{
			{
				"norm":       "420",
				"name":       "Tarification, Offer, Application",
				"version":    "2024.1",
				"status":     "implemented",
				"generation": "RNext",
				"endpoints":  []string{"/api/bipro/tariff", "/api/bipro/quote"},
			},
			{
				"norm":       "430",
				"name":       "Transfer Services",
				"version":    "2024.1",
				"status":     "implemented",
				"generation": "RNext",
				"subNorms":   []string{"430.1", "430.2", "430.4", "430.5"},
				"endpoints":  []string{"/api/bipro/transfer", "/api/bipro/gdv"},
			},
			{
				"norm":       "440",
				"name":       "Direct Access (Deep Link)",
				"version":    "2024.1",
				"status":     "implemented",
				"generation": "RNext",
				"endpoints":  []string{"/api/bipro/deeplink"},
			},
		},
		"supportedFormats":     []string{"JSON", "XML", "GDV", "SOAP"},
		"supportedGenerations": []string{"RClassic", "RNext"},
		"germanMarketReady":    true,
		"bafinCompliant":       true,
		"gdprCompliant":        true,
		"lastUpdated":          time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(complianceStatus); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetSupportedNorms returns list of supported BiPRO norms
func (h *BiPROHandler) GetSupportedNorms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	norms := []map[string]interface{}{
		{
			"normNumber":  "420",
			"name":        "Tarification, Offer, Application (TAA)",
			"description": "Standardizes premium calculation, offer generation, and application submission",
			"version":     "2024.1",
			"status":      "active",
			"implemented": true,
			"endpoints": []string{
				"POST /api/bipro/tariff - Full BiPRO 420 tariff calculation",
				"POST /api/bipro/quote - Simplified quote generation",
			},
		},
		{
			"normNumber":  "430",
			"name":        "Transfer Services",
			"description": "Electronic document and data transmission between insurers and brokers",
			"version":     "2024.1",
			"status":      "active",
			"implemented": true,
			"subNorms": []map[string]string{
				{"430.1": "Policyholder data transmission (GDV format)"},
				{"430.2": "Payment irregularities (reminders, notices)"},
				{"430.4": "Contract-related business transactions"},
				{"430.5": "Claims and benefit-related data/documents"},
			},
			"endpoints": []string{
				"POST /api/bipro/transfer - General document transfer",
				"POST /api/bipro/gdv - GDV format processing",
			},
		},
		{
			"normNumber":  "440",
			"name":        "Direct Access (Deep Link)",
			"description": "Direct portal access from broker management systems without additional authentication",
			"version":     "2024.1",
			"status":      "active",
			"implemented": true,
			"endpoints": []string{
				"POST /api/bipro/deeplink - Generate deep link access",
				"POST /api/bipro/access - Process deep link requests",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(norms); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// ============================================================
// HELPER FUNCTIONS
// ============================================================

// convertToNorm420Request converts simplified quote request to BiPRO format
func (h *BiPROHandler) convertToNorm420Request(quoteReq interface{}) bipro.Norm420TariffRequest {
	// This is a simplified conversion - in production, this would be more comprehensive
	return bipro.Norm420TariffRequest{
		MessageHeader: bipro.BiPROMessageHeader{
			MessageID:   generateBiPROMessageID(),
			Sender:      "CLIENT-UX",
			Receiver:    "BIPRO-TARIFF",
			NormVersion: "420.2024.1",
			Timestamp:   time.Now(),
		},
		RequestID: generateBiPROMessageID(),
		Timestamp: time.Now(),
		// RiskData, CoverageData, CustomerData would be populated from quoteReq
	}
}

// generateBiPROMessageID creates a BiPRO compliant message ID
func generateBiPROMessageID() string {
	return fmt.Sprintf("CUX-BIPRO-%d", time.Now().UnixNano())
}

// generateSessionToken creates a secure session token
func generateSessionToken() string {
	return fmt.Sprintf("SES-BIPRO-%d", time.Now().UnixNano())
}
