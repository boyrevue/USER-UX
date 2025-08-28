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
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type App struct {
	Ontology *OntologyData
	Sessions map[string]*QuoteSession
	Store    *sessions.CookieStore
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

	app := &App{
		Ontology: LoadOntology(),
		Sessions: make(map[string]*QuoteSession),
		Store:    store,
	}
	_ = app.loadSessionsFromDisk()

	r := mux.NewRouter()

	// Static files with proper MIME types
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set proper MIME types for different file extensions
		switch {
		case strings.HasSuffix(r.URL.Path, ".css"):
			w.Header().Set("Content-Type", "text/css")
		case strings.HasSuffix(r.URL.Path, ".js"):
			w.Header().Set("Content-Type", "application/javascript")
		case strings.HasSuffix(r.URL.Path, ".png"):
			w.Header().Set("Content-Type", "image/png")
		case strings.HasSuffix(r.URL.Path, ".jpg"), strings.HasSuffix(r.URL.Path, ".jpeg"):
			w.Header().Set("Content-Type", "image/jpeg")
		case strings.HasSuffix(r.URL.Path, ".svg"):
			w.Header().Set("Content-Type", "image/svg+xml")
		case strings.HasSuffix(r.URL.Path, ".ico"):
			w.Header().Set("Content-Type", "image/x-icon")
		case strings.HasSuffix(r.URL.Path, ".json"):
			w.Header().Set("Content-Type", "application/json")
		}
		http.FileServer(http.Dir("./static")).ServeHTTP(w, r)
	})))

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
	api.HandleFunc("/process-document", ProcessDocumentHandler).Methods("POST")
	api.HandleFunc("/validate-document", ValidateDocumentHandler).Methods("POST")

	// AI Form Analysis routes
	api.HandleFunc("/analyze-form", AnalyzeFormHandler).Methods("POST")
	api.HandleFunc("/map-fields", MapFieldsHandler).Methods("POST")
	api.HandleFunc("/shacl-transform", SHACLTransformHandler).Methods("POST")

	// Web Spider routes
	api.HandleFunc("/web-spider", WebSpiderHandler).Methods("POST")
	api.HandleFunc("/extract-data", ExtractDataHandler).Methods("POST")
	api.HandleFunc("/fill-form", FillFormHandler).Methods("POST")

	// Stealth Browser routes
	api.HandleFunc("/stealth-browser", StealthBrowserHandler).Methods("POST")

	// Bitwarden Integration routes
	api.HandleFunc("/bitwarden/status", BitwardenStatusHandler).Methods("GET")
	api.HandleFunc("/bitwarden/login", BitwardenLoginHandler).Methods("POST")
	api.HandleFunc("/bitwarden/login-apikey", BitwardenAPIKeyLoginHandler).Methods("POST")
	api.HandleFunc("/bitwarden/unlock", BitwardenUnlockHandler).Methods("POST")
	api.HandleFunc("/bitwarden/store-site", BitwardenStoreSiteHandler).Methods("POST")
	api.HandleFunc("/bitwarden/store-banking", BitwardenStoreBankingHandler).Methods("POST")
	api.HandleFunc("/bitwarden/get-site", BitwardenGetSiteHandler).Methods("GET")
	api.HandleFunc("/bitwarden/get-banking", BitwardenGetBankingHandler).Methods("GET")
	api.HandleFunc("/bitwarden/list-credentials", BitwardenListCredentialsHandler).Methods("GET")
	api.HandleFunc("/bitwarden/sync", BitwardenSyncHandler).Methods("POST")
	api.HandleFunc("/bitwarden/setup-templates", BitwardenSetupTemplatesHandler).Methods("POST")

	// Money Supermarket Infiltration
	api.HandleFunc("/infiltrate/moneysupermarket", MoneySupermarketInfiltrationHandler).Methods("POST")

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

	fmt.Println("====================================")
	fmt.Println("✅ Enhanced Insurance Quote App v2.2")
	fmt.Printf("🌐 Open: http://localhost:%s\n", port)
	fmt.Println("🌍 Languages: EN/DE, DVLA lookup, masks")
	fmt.Println("====================================")
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
		if _, ok := app.Sessions[sessionID]; !ok {
			app.Sessions[sessionID] = &QuoteSession{
				ID:           sessionID,
				Language:     "en",
				Drivers:      []Driver{},
				Progress:     map[string]bool{},
				FormData:     map[string]map[string]interface{}{},
				CreatedAt:    time.Now(),
				LastAccessed: time.Now(),
			}
			_ = app.saveSessionToDisk(sessionID)
		} else {
			app.Sessions[sessionID].LastAccessed = time.Now()
		}
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
	mods, ok := app.Ontology.Subforms["modifications"]
	if !ok {
		http.Error(w, "Not configured", http.StatusNotFound)
		return
	}
	sub := mods.Subforms[typ]
	if sub.ID == "" {
		http.Error(w, "Unknown modification type", http.StatusNotFound)
		return
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
		s.FormData = map[string]map[string]interface{}{}
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
	fields, ok := app.Ontology.Fields[category]
	if !ok {
		return result
	}
	for _, f := range fields {
		if f.Required {
			v, exists := data[f.Property]
			if !exists || v == "" {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:   f.Property,
					Message: "This field is required",
				})
			}
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
		Progress:     map[string]bool{},
		FormData:     map[string]map[string]interface{}{},
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

// Utility
func base64Key(key []byte) string { return base64.StdEncoding.EncodeToString(key) }
