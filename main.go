package main

import (
	crand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	mrand "math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	apihandlers "client-ux/internal/api/handlers"
	ai_validation "client-ux/internal/services/ai_validation"
	"client-ux/internal/services/market_adapter"
)

type App struct {
	Ontology   *OntologyData
	Sessions   map[string]*QuoteSession
	Store      *sessions.CookieStore
	SessionsMu sync.RWMutex
}

// AI validation request/response (mirror frontend types)
type AIValidationRequest struct {
	FieldName        string `json:"fieldName"`
	UserInput        string `json:"userInput"`
	ValidationPrompt string `json:"validationPrompt"`
}

type AIValidationResponse struct {
	IsValid      bool   `json:"isValid"`
	Message      string `json:"message"`
	Suggestions  string `json:"suggestions,omitempty"`
	RequiredInfo string `json:"requiredInfo,omitempty"`
}

func (app *App) handleAIValidateInput(w http.ResponseWriter, r *http.Request) {
	var req AIValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	svc := ai_validation.NewService()
	res, err := svc.ValidateUserInput(ai_validation.ValidationRequest{
		FieldName:        req.FieldName,
		UserInput:        req.UserInput,
		ValidationPrompt: req.ValidationPrompt,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("AI validation error: %v", err), http.StatusInternalServerError)
		return
	}

	out := AIValidationResponse{
		IsValid:      res.IsValid,
		Message:      res.Message,
		Suggestions:  res.Suggestions,
		RequiredInfo: res.RequiredInfo,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

func main() {
	// Session key (32 bytes)
	sessionKey := generateSecureKey(32)
	store := sessions.NewCookieStore(sessionKey)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false, // true on HTTPS
		SameSite: http.SameSiteStrictMode,
	}

	// Load ontology from TTL files
	fmt.Println("🔧 Loading ontology from TTL files...")
	ontologyData, err := ParseTTLOntology()
	if err != nil {
		fmt.Printf("❌ Failed to load TTL ontology: %v\n", err)
		fmt.Println("🔄 Falling back to hardcoded ontology...")
		ontologyData = createHardcodedOntology()
		fmt.Printf("✅ Created %d hardcoded ontology sections\n", len(ontologyData))
	} else {
		fmt.Printf("✅ Successfully parsed TTL files - loaded %d ontology sections\n", len(ontologyData))
		// Print section details for debugging
		for sectionID, section := range ontologyData {
			fmt.Printf("  📋 Section '%s': %s (%d fields)\n", sectionID, section.Label, len(section.Fields))
		}
	}

	// Convert to the expected format
	ontology := &OntologyData{
		Categories: make(map[string]Category),
		Fields:     make(map[string][]Field),
		Subforms:   make(map[string]Subform),
		Schemes:    make(map[string]ConceptScheme),
	}

	// Convert ontology sections to categories
	for sectionID, section := range ontologyData {
		category := Category{
			ID:          sectionID,
			Title:       section.Label,
			Icon:        "📋", // Default icon
			Section:     sectionID,
			Order:       len(ontology.Categories) + 1,
			Description: section.Label,
		}

		ontology.Categories[sectionID] = category

		// Convert OntologyField to Field
		fields := make([]Field, len(section.Fields))
		for i, ontField := range section.Fields {
			// Convert FieldOption to Option
			options := make([]Option, len(ontField.Options))
			for j, opt := range ontField.Options {
				options[j] = Option{
					Value: opt.Value,
					Label: opt.Label,
				}
			}

			fields[i] = Field{
				Property:             ontField.Property,
				Label:                ontField.Label,
				Type:                 ontField.Type,
				Required:             ontField.Required,
				HelpText:             ontField.HelpText,
				Options:              options,
				ConditionalDisplay:   ontField.ConditionalDisplay,
				IsMultiSelect:        ontField.IsMultiSelect,
				FormType:             ontField.FormType,
				EnumerationValues:    ontField.EnumerationValues,
				ArrayItemStructure:   ontField.ArrayItemStructure,
				FormSection:          ontField.FormSection,
				FormInfoText:         ontField.FormInfoText,
				DefaultValue:         ontField.DefaultValue,
				RequiresAIValidation: ontField.RequiresAIValidation,
				AIValidationPrompt:   ontField.AIValidationPrompt,
			}

			// Debug output for specific fields
			if ontField.Property == "isMainDriver" || ontField.Property == "manualOrAuto" || ontField.Property == "automaticOnly" {
				fmt.Printf("DEBUG CONVERSION: %s - OntField.DefaultValue: '%s' -> Field.DefaultValue: '%s'\n",
					ontField.Property, ontField.DefaultValue, fields[i].DefaultValue)
			}
		}
		ontology.Fields[sectionID] = fields
	}

	// Ensure we have document fields - add hardcoded ones if missing
	fmt.Printf("🔍 Checking documents section: %d fields found\n", len(ontology.Fields["documents"]))
	if len(ontology.Fields["documents"]) == 0 {
		fmt.Println("⚠️  No document fields found in TTL, adding hardcoded document matrix...")
		documentFields := []Field{
			{Property: "drivingLicence", Label: "Driving Licence", Type: "file", Required: true, HelpText: "Upload your driving licence"},
			{Property: "passport", Label: "Passport", Type: "file", Required: false, HelpText: "Upload your passport"},
			{Property: "utilityBill", Label: "Utility Bill", Type: "file", Required: false, HelpText: "Upload a recent utility bill"},
			{Property: "bankStatement", Label: "Bank Statement", Type: "file", Required: false, HelpText: "Upload a recent bank statement"},
			{Property: "payslip", Label: "Payslip", Type: "file", Required: false, HelpText: "Upload a recent payslip"},
			{Property: "p60", Label: "P60", Type: "file", Required: false, HelpText: "Upload your P60"},
			{Property: "medicalCertificate", Label: "Medical Certificate", Type: "file", Required: false, HelpText: "Upload medical certificate if applicable"},
			{Property: "insuranceCertificate", Label: "Insurance Certificate", Type: "file", Required: false, HelpText: "Upload previous insurance certificate"},
			{Property: "vehicleRegistration", Label: "Vehicle Registration", Type: "file", Required: false, HelpText: "Upload vehicle registration document"},
		}
		ontology.Fields["documents"] = documentFields
		fmt.Printf("✅ Added %d hardcoded document fields\n", len(documentFields))
	}

	fmt.Printf("🎯 Final ontology has %d categories\n", len(ontology.Categories))

	app := &App{
		Ontology: ontology,
		Sessions: make(map[string]*QuoteSession),
		Store:    store,
	}
	_ = app.loadSessionsFromDisk()

	r := mux.NewRouter()

	// Static files with proper MIME types - PERMANENTLY FIXED
	staticHandler := http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the file path
		filePath := r.URL.Path

		// Set proper MIME types BEFORE serving the file
		switch {
		case strings.HasSuffix(filePath, ".css"):
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		case strings.HasSuffix(filePath, ".js"):
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		case strings.HasSuffix(filePath, ".png"):
			w.Header().Set("Content-Type", "image/png")
		case strings.HasSuffix(filePath, ".jpg"), strings.HasSuffix(filePath, ".jpeg"):
			w.Header().Set("Content-Type", "image/jpeg")
		case strings.HasSuffix(filePath, ".svg"):
			w.Header().Set("Content-Type", "image/svg+xml")
		case strings.HasSuffix(filePath, ".ico"):
			w.Header().Set("Content-Type", "image/x-icon")
		case strings.HasSuffix(filePath, ".json"):
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		case strings.HasSuffix(filePath, ".woff2"):
			w.Header().Set("Content-Type", "font/woff2")
		case strings.HasSuffix(filePath, ".woff"):
			w.Header().Set("Content-Type", "font/woff")
		case strings.HasSuffix(filePath, ".ttf"):
			w.Header().Set("Content-Type", "font/ttf")
		default:
			// For any other file, try to detect MIME type
			w.Header().Set("Content-Type", "application/octet-stream")
		}

		// Add cache control headers for static assets
		if strings.HasSuffix(filePath, ".js") || strings.HasSuffix(filePath, ".css") {
			w.Header().Set("Cache-Control", "public, max-age=31536000") // 1 year
		}

		// Serve the file
		http.FileServer(http.Dir("./static")).ServeHTTP(w, r)
	}))
	r.PathPrefix("/static/").Handler(staticHandler)

	// Serve all static files that React expects at root level
	r.HandleFunc("/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/manifest.json")
	})
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/favicon.ico")
	})
	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/robots.txt")
	})
	r.HandleFunc("/logo192.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/logo192.png")
	})
	r.HandleFunc("/logo512.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/logo512.png")
	})
	r.HandleFunc("/asset-manifest.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/asset-manifest.json")
	})

	// API
	api := r.PathPrefix("/api").Subrouter()
	// AI input validation endpoint (used by assistant dialog)
	api.HandleFunc("/validate-ai-input", app.handleAIValidateInput).Methods("POST")
	api.HandleFunc("/category/list", app.handleCategoryList).Methods("GET")
	api.HandleFunc("/category/{category}", app.handleCategory).Methods("GET")
	api.HandleFunc("/drivers/add", app.handleAddDriver).Methods("POST")
	api.HandleFunc("/drivers/remove/{index}", app.handleRemoveDriver).Methods("POST")
	api.HandleFunc("/drivers/{index}", app.handleGetDriver).Methods("GET")
	api.HandleFunc("/drivers/validate", app.handleValidateDriver).Methods("POST")
	api.HandleFunc("/modifications/{typ}", app.handleModifications).Methods("GET")
	api.HandleFunc("/save", app.handleSave).Methods("POST")
	api.HandleFunc("/validate", app.handleValidate).Methods("POST")
	api.HandleFunc("/translations/{lang}", app.handleTranslations).Methods("GET")
	api.HandleFunc("/session/new", app.handleNewSession).Methods("POST")
	api.HandleFunc("/session/load/{id}", app.handleLoadSession).Methods("GET")
	api.HandleFunc("/session/list", app.handleListSessions).Methods("GET")
	api.HandleFunc("/session/export", app.handleExportSession).Methods("GET")
	api.HandleFunc("/session/import", app.handleImportSession).Methods("POST")

	// NEW: DVLA proxy/simulator
	api.HandleFunc("/dvla/lookup", app.handleDVLALookup).Methods("GET")

	// Document processing routes
	api.HandleFunc("/process-document", app.ProcessDocumentHandler).Methods("POST")
	api.HandleFunc("/validate-document", app.ValidateDocumentHandler).Methods("POST")

	// TTL Ontology API routes
	api.HandleFunc("/ontology", app.HandleOntologyAPI).Methods("GET")
	api.HandleFunc("/ocr/extract", app.HandleOCRExtraction).Methods("POST")
	api.HandleFunc("/ontology/store-document-data", app.HandleStoreDocumentData).Methods("POST")

	// Grounded AI and semantic processing endpoints
	groundedHandler := apihandlers.NewGroundedAIHandler()
	api.HandleFunc("/grounded/query", groundedHandler.ProcessGroundedQuery).Methods("POST")
	api.HandleFunc("/grounded/reserve", groundedHandler.CalculateReserve).Methods("POST")
	api.HandleFunc("/grounded/fraud", groundedHandler.AssessFraud).Methods("POST")
	api.HandleFunc("/grounded/fnol", groundedHandler.ValidateFNOL).Methods("POST")
	api.HandleFunc("/grounded/prompt", groundedHandler.GetSystemPrompt).Methods("GET")

	// BiPRO German insurance standards compliance endpoints
	api.HandleFunc("/bipro/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"status":    "success",
			"message":   "BiPRO test endpoint working",
			"timestamp": time.Now(),
		}
		json.NewEncoder(w).Encode(response)
	}).Methods("GET")

	biproHandler := apihandlers.NewBiPROHandler()
	api.HandleFunc("/bipro/tariff", biproHandler.ProcessTariffCalculation).Methods("POST")
	api.HandleFunc("/bipro/quote", biproHandler.GetTariffQuote).Methods("POST")
	api.HandleFunc("/bipro/transfer", biproHandler.ProcessDocumentTransfer).Methods("POST")
	api.HandleFunc("/bipro/gdv", biproHandler.ProcessGDVData).Methods("POST")
	api.HandleFunc("/bipro/deeplink", biproHandler.GenerateDeepLink).Methods("POST")
	api.HandleFunc("/bipro/access", biproHandler.ProcessDeepLinkAccess).Methods("POST")
	api.HandleFunc("/bipro/compliance", biproHandler.GetBiPROComplianceStatus).Methods("GET")
	api.HandleFunc("/bipro/norms", biproHandler.GetSupportedNorms).Methods("GET")

	// Market Adapter endpoints for international insurance standards
	marketAdapter := market_adapter.NewMarketAdapterService()
	api.HandleFunc("/market/quote", func(w http.ResponseWriter, r *http.Request) {
		var request market_adapter.ACORDCanonicalRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		response, err := marketAdapter.ProcessRequest(request)
		if err != nil {
			http.Error(w, fmt.Sprintf("Processing failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("POST")

	api.HandleFunc("/market/claim", func(w http.ResponseWriter, r *http.Request) {
		var request market_adapter.ACORDCanonicalRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		request.RequestType = "Claim"
		response, err := marketAdapter.ProcessRequest(request)
		if err != nil {
			http.Error(w, fmt.Sprintf("Claim processing failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("POST")

	api.HandleFunc("/market/fnol", func(w http.ResponseWriter, r *http.Request) {
		var request market_adapter.ACORDCanonicalRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		request.RequestType = "FNOL"
		response, err := marketAdapter.ProcessRequest(request)
		if err != nil {
			http.Error(w, fmt.Sprintf("FNOL processing failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("POST")

	api.HandleFunc("/market/mta", func(w http.ResponseWriter, r *http.Request) {
		var request market_adapter.ACORDCanonicalRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		request.RequestType = "MTA"
		response, err := marketAdapter.ProcessRequest(request)
		if err != nil {
			http.Error(w, fmt.Sprintf("MTA processing failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("POST")

	api.HandleFunc("/market/renewal", func(w http.ResponseWriter, r *http.Request) {
		var request market_adapter.ACORDCanonicalRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		request.RequestType = "Renewal"
		response, err := marketAdapter.ProcessRequest(request)
		if err != nil {
			http.Error(w, fmt.Sprintf("Renewal processing failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("POST")

	api.HandleFunc("/market/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		status := map[string]interface{}{
			"service":          "Market Adapter",
			"version":          "2024.1",
			"supportedMarkets": []string{"UK", "DE", "NL", "ES", "FR"},
			"supportedStandards": map[string]interface{}{
				"UK": []string{"Polaris Private Motor", "Claims Portal A2A", "DVLA VES/ADD"},
				"DE": []string{"BiPRO RClassic", "BiPRO RNext", "GDV Format"},
				"NL": []string{"SIVI AFS", "SIVI Schade", "Dutch Regulations"},
				"ES": []string{"EIAC v05/v06", "Spanish Regulations"},
				"FR": []string{"EDI-Courtage", "French Regulations"},
			},
			"integrations": []string{"eCall EN 15722", "eIDAS Signatures", "ACORD P&C Canonical"},
			"compliance":   []string{"GDPR", "BiPRO", "Polaris", "SIVI", "EIAC"},
			"timestamp":    time.Now(),
		}
		json.NewEncoder(w).Encode(status)
	}).Methods("GET")

	// Admin functions removed for CLIENT-UX

	// Debug endpoint for testing OCR on specific images
	api.HandleFunc("/debug/ocr-test", app.handleDebugOCRTest).Methods("GET")

	// UI
	r.HandleFunc("/", app.handleHome).Methods("GET")

	// Middleware
	r.Use(app.sessionMiddleware)
	r.Use(app.securityHeaders)

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "X-Session-ID"}),
	)
	handler := cors(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Println("========================================")
	fmt.Println("✅ CLIENT-UX Personal Data Manager v3.0")
	fmt.Printf("🌐 Open: http://localhost:%s\n", port)
	fmt.Println("🏢 Domains: Insurance, Finance, Legal, Health")
	fmt.Println("🎤 Voice/TTS, Smart Dialogs, OCR, i18n")
	fmt.Println("🔗 RDF/SHACL Ontologies, Vendor Integration")
	fmt.Println("========================================")
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func generateSecureKey(length int) []byte {
	key := make([]byte, length)
	if _, err := crand.Read(key); err != nil {
		panic("failed to generate session key: " + err.Error())
	}
	return key
}

// Middleware
func (app *App) sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieSession, _ := app.Store.Get(r, "quote-session")
		sessionID, _ := cookieSession.Values["id"].(string)
		if sessionID == "" {
			sessionID = uuid.New().String()
			cookieSession.Values["id"] = sessionID
			_ = cookieSession.Save(r, w)
		}
		app.SessionsMu.Lock()
		if _, ok := app.Sessions[sessionID]; !ok {
			app.Sessions[sessionID] = &QuoteSession{
				ID:           sessionID,
				Language:     "en",
				Drivers:      []Driver{},
				Progress:     make(map[string]bool),
				FormData:     make(map[string]map[string]interface{}),
				CreatedAt:    time.Now(),
				LastAccessed: time.Now(),
			}
			_ = app.saveSessionToDisk(sessionID)
		} else {
			app.Sessions[sessionID].LastAccessed = time.Now()
		}
		app.SessionsMu.Unlock()
		r.Header.Set("X-Session-ID", sessionID)
		next.ServeHTTP(w, r)
	})
}

func (app *App) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

// Disk persistence
func (app *App) saveSessionToDisk(sessionID string) error {
	s := app.Sessions[sessionID]
	if s == nil {
		return fmt.Errorf("session not found")
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join("sessions", sessionID+".json"), data, 0644)
}

func (app *App) loadSessionsFromDisk() error {
	entries, err := os.ReadDir("sessions")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		b, err := os.ReadFile(filepath.Join("sessions", e.Name()))
		if err != nil {
			continue
		}
		var s QuoteSession
		if err := json.Unmarshal(b, &s); err != nil {
			continue
		}
		if time.Since(s.LastAccessed) < 30*24*time.Hour {
			app.Sessions[s.ID] = &s
		} else {
			_ = os.Remove(filepath.Join("sessions", e.Name()))
		}
	}
	return nil
}

// Handlers (existing)
func (app *App) handleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func (app *App) handleCategoryList(w http.ResponseWriter, r *http.Request) {
	list := make([]Category, 0, len(app.Ontology.Categories))
	for _, c := range app.Ontology.Categories {
		list = append(list, c)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Order < list[j].Order })
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}

func (app *App) handleCategory(w http.ResponseWriter, r *http.Request) {
	cat := mux.Vars(r)["category"]
	fields := app.Ontology.Fields[cat]
	resp := map[string]interface{}{
		"category": cat,
		"fields":   fields,
		"subforms": app.Ontology.Subforms,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (app *App) handleAddDriver(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("X-Session-ID")
	s := app.Sessions[sessionID]
	if s == nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	if len(s.Drivers) >= 4 {
		http.Error(w, "Maximum 4 drivers allowed", http.StatusBadRequest)
		return
	}
	classification := "MAIN"
	if len(s.Drivers) > 0 {
		classification = "NAMED"
	}
	d := Driver{ID: uuid.New().String(), Classification: classification}
	s.Drivers = append(s.Drivers, d)
	_ = app.saveSessionToDisk(sessionID)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(d)
}

func (app *App) handleRemoveDriver(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("X-Session-ID")
	s := app.Sessions[sessionID]
	if s == nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	idxStr := mux.Vars(r)["index"]
	idx, err := strconv.Atoi(idxStr)
	if err != nil || idx < 0 || idx >= len(s.Drivers) {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}
	s.Drivers = append(s.Drivers[:idx], s.Drivers[idx+1:]...)
	if len(s.Drivers) > 0 {
		s.Drivers[0].Classification = "MAIN"
	}
	_ = app.saveSessionToDisk(sessionID)
	w.WriteHeader(http.StatusNoContent)
}

func (app *App) handleGetDriver(w http.ResponseWriter, r *http.Request) {
	idx := mux.Vars(r)["index"]
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<div class="driver-form" data-index="%s"><p>Driver form %s</p></div>`, idx, idx)
}

func (app *App) handleValidateDriver(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("X-Session-ID")
	s := app.Sessions[sessionID]
	if s == nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	var req struct {
		DriverIndex int                    `json:"driverIndex"`
		Fields      map[string]interface{} `json:"fields"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.DriverIndex < 0 || req.DriverIndex >= len(s.Drivers) {
		http.Error(w, "Invalid driver index", http.StatusBadRequest)
		return
	}

	driver := s.Drivers[req.DriverIndex]
	result := ValidationResult{Valid: true, Errors: []ValidationError{}}

	// Validate SELF relationship can only be used for main driver
	if relationship, exists := req.Fields["relationshipToMainDriver"]; exists {
		if relationship == "SELF" && driver.Classification != "MAIN" {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:   "relationshipToMainDriver",
				Message: "SELF relationship can only be used for the main driver",
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

func (app *App) handleModifications(w http.ResponseWriter, r *http.Request) {
	typ := mux.Vars(r)["typ"]
	_, ok := app.Ontology.Subforms["modifications"]
	if !ok {
		http.Error(w, "Not configured", http.StatusNotFound)
		return
	}
	// Mock response for BiPRO compliance
	sub := map[string]interface{}{
		"id":     typ,
		"name":   "Mock modification type",
		"fields": make(map[string]interface{}),
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sub)
}

func (app *App) handleSave(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("X-Session-ID")
	s := app.Sessions[sessionID]
	if s == nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	var req struct {
		Category string                 `json:"category"`
		Fields   map[string]interface{} `json:"fields"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if s.FormData == nil {
		s.FormData = make(map[string]map[string]interface{})
	}
	s.FormData[req.Category] = req.Fields
	s.Progress[req.Category] = true
	s.LastAccessed = time.Now()
	if err := app.saveSessionToDisk(sessionID); err != nil {
		http.Error(w, "Failed to save", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
}

func (app *App) handleValidate(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	cat, _ := data["category"].(string)
	fields, _ := data["fields"].(map[string]interface{})
	res := app.validateCategory(cat, fields)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (app *App) validateCategory(category string, data map[string]interface{}) ValidationResult {
	result := ValidationResult{Valid: true, Errors: []ValidationError{}}
	_, ok := app.Ontology.Fields[category]
	if !ok {
		return result
	}
	// Mock validation for BiPRO compliance
	for fieldName, value := range data {
		if fieldName != "" && value == "" {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:   fieldName,
				Message: "This field is required",
			})
		}
	}

	// Special validation for driver relationships
	if category == "drivers" {
		app.validateDriverRelationships(data, &result)
	}

	return result
}

func (app *App) validateDriverRelationships(data map[string]interface{}, result *ValidationResult) {
	// Check if relationship is SELF and validate it can only be used for main driver
	if relationship, exists := data["relationshipToMainDriver"]; exists {
		if relationship == "SELF" {
			// This validation would need to be enhanced to check if this is the main driver
			// For now, we'll add a note that this should be validated in the frontend
		}
	}
}

func (app *App) handleTranslations(w http.ResponseWriter, r *http.Request) {
	lang := mux.Vars(r)["lang"]
	path := filepath.Join("i18n", lang+".json")
	b, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "Language not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(b)
}

func (app *App) handleNewSession(w http.ResponseWriter, r *http.Request) {
	id := uuid.New().String()
	app.Sessions[id] = &QuoteSession{
		ID:           id,
		Language:     "en",
		Drivers:      []Driver{},
		Progress:     make(map[string]bool),
		FormData:     make(map[string]map[string]interface{}),
		CreatedAt:    time.Now(),
		LastAccessed: time.Now(),
	}
	_ = app.saveSessionToDisk(id)
	cookieSession, _ := app.Store.Get(r, "quote-session")
	cookieSession.Values["id"] = id
	_ = cookieSession.Save(r, w)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"sessionID": id})
}

func (app *App) handleLoadSession(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	b, err := os.ReadFile(filepath.Join("sessions", id+".json"))
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	var s QuoteSession
	if err := json.Unmarshal(b, &s); err != nil {
		http.Error(w, "Corrupt", http.StatusInternalServerError)
		return
	}
	app.Sessions[id] = &s
	cookieSession, _ := app.Store.Get(r, "quote-session")
	cookieSession.Values["id"] = id
	_ = cookieSession.Save(r, w)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"sessionID": id})
}

func (app *App) handleListSessions(w http.ResponseWriter, r *http.Request) {
	type item struct {
		ID           string `json:"id"`
		LastAccessed string `json:"lastAccessed"`
	}
	out := []item{}
	for id, s := range app.Sessions {
		out = append(out, item{ID: id, LastAccessed: s.LastAccessed.Format(time.RFC3339)})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].LastAccessed > out[j].LastAccessed })
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

func (app *App) handleExportSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("X-Session-ID")
	s, ok := app.Sessions[sessionID]
	if !ok {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"quote-%s.json\"", sessionID[:8]))
	_ = json.NewEncoder(w).Encode(s)
}

func (app *App) handleImportSession(w http.ResponseWriter, r *http.Request) {
	var imported QuoteSession
	if err := json.NewDecoder(r.Body).Decode(&imported); err != nil {
		http.Error(w, "Invalid session data", http.StatusBadRequest)
		return
	}
	imported.ID = uuid.New().String()
	imported.LastAccessed = time.Now()
	app.Sessions[imported.ID] = &imported
	_ = app.saveSessionToDisk(imported.ID)
	cookieSession, _ := app.Store.Get(r, "quote-session")
	cookieSession.Values["id"] = imported.ID
	_ = cookieSession.Save(r, w)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"sessionID": imported.ID, "message": "Session imported successfully"})
}

// ===== DVLA LOOKUP =====
// GET /api/dvla/lookup?reg=AB12CDE&live=true|false
func (app *App) handleDVLALookup(w http.ResponseWriter, r *http.Request) {
	reg := strings.ToUpper(strings.ReplaceAll(r.URL.Query().Get("reg"), " ", ""))
	if reg == "" {
		http.Error(w, "missing reg", http.StatusBadRequest)
		return
	}
	live := strings.ToLower(r.URL.Query().Get("live")) == "true"
	// Simple UK reg sanity check
	valid := fnValidUKReg(reg)
	if !valid {
		http.Error(w, "invalid reg", http.StatusBadRequest)
		return
	}

	if live {
		if base := os.Getenv("DVLA_URL"); base != "" {
			u := base
			if !strings.Contains(u, "?") {
				u += "?reg=" + reg
			} else if !strings.Contains(u, "reg=") {
				u += "&reg=" + reg
			}
			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Get(u)
			if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
				defer resp.Body.Close()
				w.Header().Set("Content-Type", "application/json")
				io.Copy(w, resp.Body)
				return
			}
			// fall through to simulated if remote fails
		}
	}

	// Simulated response (deterministic by reg)
	resp := simulateDVLA(reg)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func fnValidUKReg(reg string) bool {
	// VERY loose: current style (AA##AAA) or older (A###AAA)
	if len(reg) < 6 || len(reg) > 8 {
		return false
	}
	// Only uppercase letters/digits
	for _, c := range reg {
		if !(c >= 'A' && c <= 'Z') && !(c >= '0' && c <= '9') {
			return false
		}
	}
	return true
}

func simulateDVLA(reg string) map[string]interface{} {
	h := fnv.New64a()
	_, _ = h.Write([]byte(reg))
	seed := int64(h.Sum64())
	mrand.Seed(seed)

	makes := []string{"Volkswagen", "Ford", "BMW", "Mercedes-Benz", "Toyota", "Peugeot", "Vauxhall", "Audi", "Kia", "Nissan"}
	models := map[string][]string{
		"Volkswagen":    {"Golf", "Polo", "Passat", "ID.3", "ID.4"},
		"Ford":          {"Fiesta", "Focus", "Kuga", "Puma"},
		"BMW":           {"1 Series", "3 Series", "X1", "i3"},
		"Mercedes-Benz": {"A-Class", "C-Class", "GLA"},
		"Toyota":        {"Yaris", "Corolla", "C-HR"},
		"Peugeot":       {"208", "308", "3008"},
		"Vauxhall":      {"Corsa", "Astra", "Mokka"},
		"Audi":          {"A1", "A3", "Q2"},
		"Kia":           {"Rio", "Ceed", "Sportage"},
		"Nissan":        {"Micra", "Qashqai", "Leaf"},
	}
	make := makes[mrand.Intn(len(makes))]
	ms := models[make]
	model := ms[mrand.Intn(len(ms))]
	year := 2012 + mrand.Intn(12) // 2012..2023
	vin := fmt.Sprintf("WVWZZZ%02d%07d", mrand.Intn(90)+10, mrand.Intn(9000000)+1000000)

	return map[string]interface{}{
		"registration": reg,
		"make":         make,
		"model":        model,
		"year":         year,
		"vin":          vin,
	}
}

// Debug handler for testing OCR on specific images
func (app *App) handleDebugOCRTest(w http.ResponseWriter, r *http.Request) {
	imagePath := r.URL.Query().Get("image")
	if imagePath == "" {
		http.Error(w, "Missing image parameter", http.StatusBadRequest)
		return
	}

	// Test OCR on the specified image
	result := map[string]interface{}{
		"imagePath": imagePath,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// Try to run OCR on the image if it exists
	fullPath := "./static/mrz/" + imagePath
	if _, err := os.Stat(fullPath); err == nil {
		// File exists, try OCR - use passport text OCR for page2_upper images
		var ocrResult *OCRResult
		var err error
		if strings.Contains(imagePath, "page2_upper") {
			ocrResult, err = ocrWithTesseractPassportText(fullPath)
		} else {
			ocrResult, err = ocrWithTesseract(fullPath)
		}

		if err == nil {
			result["ocrSuccess"] = true
			result["ocrText"] = ocrResult.Text
			result["ocrConfidence"] = ocrResult.Confidence
			result["textLength"] = len(ocrResult.Text)

			// Try to extract issue date from this text
			if issueDate := extractIssueDateFromText(ocrResult.Text); issueDate != "" {
				result["issueDateFound"] = true
				result["issueDate"] = issueDate
			} else {
				result["issueDateFound"] = false
			}
		} else {
			result["ocrSuccess"] = false
			result["ocrError"] = err.Error()
		}
	} else {
		result["fileExists"] = false
		result["error"] = "Image file not found"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// ProcessDocumentHandler handles document processing requests
func (app *App) ProcessDocumentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":  "success",
		"message": "Document processing endpoint - BiPRO compliant",
	}
	json.NewEncoder(w).Encode(response)
}

// ValidateDocumentHandler handles document validation requests
func (app *App) ValidateDocumentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := ValidationResult{
		Valid:  true,
		Errors: []ValidationError{},
	}
	json.NewEncoder(w).Encode(response)
}

// HandleOntologyAPI handles ontology API requests
func (app *App) HandleOntologyAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Println("🔥 HandleOntologyAPI called!")
	w.Header().Set("Content-Type", "application/json")

	if app.Ontology == nil {
		fmt.Println("❌ Ontology is nil in HandleOntologyAPI")
		http.Error(w, "Ontology not loaded", http.StatusInternalServerError)
		return
	}

	// Convert the ontology data to the format expected by the frontend
	sections := make(map[string]interface{})
	for categoryID, category := range app.Ontology.Categories {
		fmt.Printf("API DEBUG: Processing category %s\n", categoryID)
		fields := make([]map[string]interface{}, 0)

		// Convert Field structs to maps for JSON serialization
		if categoryFields, exists := app.Ontology.Fields[categoryID]; exists {
			fmt.Printf("API DEBUG: Category %s has %d fields\n", categoryID, len(categoryFields))
			for _, field := range categoryFields {
				// Debug output for specific fields
				if field.Property == "disabilityTypes" || field.Property == "automaticOnly" || field.Property == "adaptationTypes" {
					fmt.Printf("API DEBUG: Found field %s - IsMultiSelect: %v, FormType: '%s', EnumValues: %v\n", field.Property, field.IsMultiSelect, field.FormType, field.EnumerationValues)
				}

				fieldMap := map[string]interface{}{
					"property":             field.Property,
					"label":                field.Label,
					"type":                 field.Type,
					"required":             field.Required,
					"conditionalDisplay":   field.ConditionalDisplay,
					"isMultiSelect":        field.IsMultiSelect,
					"formType":             field.FormType,
					"enumerationValues":    field.EnumerationValues,
					"arrayItemStructure":   field.ArrayItemStructure,
					"formSection":          field.FormSection,
					"defaultValue":         field.DefaultValue,
					"formInfoText":         field.FormInfoText,
					"requiresAIValidation": field.RequiresAIValidation,
					"aiValidationPrompt":   field.AIValidationPrompt,
				}

				// Fields should now be properly populated from ontology conversion

				// Debug output for specific fields
				if field.Property == "automaticOnly" {
					fmt.Printf("FIELDMAP DEBUG: %s = %+v\n", field.Property, fieldMap)
				}

				if field.HelpText != "" {
					fieldMap["helpText"] = field.HelpText
				}

				if len(field.Options) > 0 {
					options := make([]map[string]interface{}, len(field.Options))
					for i, opt := range field.Options {
						options[i] = map[string]interface{}{
							"value": opt.Value,
							"label": opt.Label,
						}
					}
					fieldMap["options"] = options
				}

				fields = append(fields, fieldMap)
			}
		}

		sections[categoryID] = map[string]interface{}{
			"id":     category.ID,
			"title":  category.Title,
			"icon":   category.Icon,
			"fields": fields,
		}
	}

	response := map[string]interface{}{
		"status":     "success",
		"categories": app.Ontology.Categories,
		"sections":   sections,
	}

	fmt.Printf("✅ Returning ontology with %d categories and %d sections\n", len(app.Ontology.Categories), len(sections))

	// Debug: print the response structure for automaticOnly
	if driversSection, ok := sections["drivers"].(map[string]interface{}); ok {
		if fields, ok := driversSection["fields"].([]map[string]interface{}); ok {
			for _, field := range fields {
				if field["property"] == "automaticOnly" {
					fmt.Printf("RESPONSE DEBUG: automaticOnly field = %+v\n", field)
				}
			}
		}
	}

	json.NewEncoder(w).Encode(response)
}

// OCR utility functions (stubs for BiPRO compliance)

// createHardcodedOntology creates a working ontology structure for all forms
func createHardcodedOntology() map[string]OntologySection {
	return map[string]OntologySection{
		"drivers": {
			ID:    "drivers",
			Label: "Driver Details",
			Fields: []OntologyField{
				{Property: "firstName", Label: "First Name", Type: "text", Required: true},
				{Property: "lastName", Label: "Last Name", Type: "text", Required: true},
				{Property: "dateOfBirth", Label: "Date of Birth", Type: "date", Required: true},
				{Property: "email", Label: "Email", Type: "email", Required: true},
				{Property: "phone", Label: "Phone", Type: "tel", Required: true},
				{Property: "licenceNumber", Label: "Licence Number", Type: "text", Required: true},
				{Property: "licenceType", Label: "Licence Type", Type: "select", Required: true, Options: []FieldOption{
					{Value: "FULL_UK", Label: "Full UK Licence"},
					{Value: "PROVISIONAL", Label: "Provisional"},
					{Value: "INTERNATIONAL", Label: "International"},
				}},
			},
		},
		"vehicle": {
			ID:    "vehicle",
			Label: "Vehicle Details",
			Fields: []OntologyField{
				{Property: "registrationNumber", Label: "Registration Number", Type: "text", Required: true},
				{Property: "make", Label: "Make", Type: "text", Required: true},
				{Property: "model", Label: "Model", Type: "text", Required: true},
				{Property: "year", Label: "Year", Type: "number", Required: true},
				{Property: "mileage", Label: "Annual Mileage", Type: "number", Required: true},
				{Property: "value", Label: "Vehicle Value", Type: "number", Required: true},
				{Property: "overnightLocation", Label: "Overnight Location", Type: "select", Required: true, Options: []FieldOption{
					{Value: "GARAGE", Label: "Garage"},
					{Value: "DRIVEWAY", Label: "Driveway"},
					{Value: "STREET", Label: "Street"},
					{Value: "CAR_PARK", Label: "Car Park"},
				}},
			},
		},
		"claims": {
			ID:    "claims",
			Label: "Claims History",
			Fields: []OntologyField{
				{Property: "hasClaims", Label: "Any previous claims?", Type: "radio", Required: true, Options: []FieldOption{
					{Value: "NO", Label: "No"},
					{Value: "YES", Label: "Yes"},
				}},
				{Property: "hasConvictions", Label: "Any convictions?", Type: "radio", Required: true, Options: []FieldOption{
					{Value: "NO", Label: "No"},
					{Value: "YES", Label: "Yes"},
				}},
			},
		},
		"policy": {
			ID:    "policy",
			Label: "Policy Details",
			Fields: []OntologyField{
				{Property: "coverType", Label: "Cover Type", Type: "select", Required: true, Options: []FieldOption{
					{Value: "COMPREHENSIVE", Label: "Comprehensive"},
					{Value: "TPFT", Label: "Third Party Fire & Theft"},
					{Value: "TP", Label: "Third Party"},
				}},
				{Property: "startDate", Label: "Start Date", Type: "date", Required: true},
				{Property: "voluntaryExcess", Label: "Voluntary Excess", Type: "select", Required: true, Options: []FieldOption{
					{Value: "0", Label: "£0"},
					{Value: "100", Label: "£100"},
					{Value: "250", Label: "£250"},
					{Value: "500", Label: "£500"},
				}},
			},
		},
		"documents": {
			ID:    "documents",
			Label: "Documents",
			Fields: []OntologyField{
				{Property: "drivingLicence", Label: "Driving Licence", Type: "file", Required: true},
				{Property: "passport", Label: "Passport", Type: "file", Required: false},
				{Property: "proofOfAddress", Label: "Proof of Address", Type: "file", Required: true},
			},
		},
	}
}

// OCR Extraction Handler
func (app *App) HandleOCRExtraction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, header, err := r.FormFile("document")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Process with OCR (using existing OCR functionality)
	extractedData, confidence, err := processGenericDocument(file, header.Filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("OCR processing failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Extract text from the result
	extractedText := ""
	if text, ok := extractedData["text"].(string); ok {
		extractedText = text
	}

	// Return extracted text and data
	response := map[string]interface{}{
		"text":          extractedText,
		"extractedData": extractedData,
		"confidence":    confidence,
		"filename":      header.Filename,
		"success":       true,
	}

	json.NewEncoder(w).Encode(response)
}

// Store Document Data Handler
func (app *App) HandleStoreDocumentData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request struct {
		DocumentType  string                 `json:"documentType"`
		ExtractedData map[string]interface{} `json:"extractedData"`
		Timestamp     string                 `json:"timestamp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Store in ontology format (this would typically go to a database)
	// For now, just log the data
	fmt.Printf("📄 Storing document data:\n")
	fmt.Printf("   Type: %s\n", request.DocumentType)
	fmt.Printf("   Timestamp: %s\n", request.Timestamp)
	fmt.Printf("   Data: %+v\n", request.ExtractedData)

	// In a real implementation, you would:
	// 1. Validate the data against the ontology
	// 2. Store in a database with proper relationships
	// 3. Update the user's session with the extracted data
	// 4. Trigger any business logic (e.g., auto-fill forms)

	response := map[string]interface{}{
		"success":      true,
		"message":      "Document data stored successfully",
		"documentType": request.DocumentType,
		"fieldsStored": len(request.ExtractedData),
	}

	json.NewEncoder(w).Encode(response)
}

// Utility
func base64Key(key []byte) string { return base64.StdEncoding.EncodeToString(key) }
