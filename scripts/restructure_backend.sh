#!/bin/bash

# Backend Restructuring Script for CLIENT-UX
# Creates clean Go architecture with proper separation of concerns

echo "🏗️ Restructuring Go backend for better AI manageability..."

# Create the new directory structure
echo "📁 Creating directory structure..."
mkdir -p {cmd/server,internal/{api/{handlers,middleware,routes},services/{ocr,validation,gdpr,ontology},models,repository,config}}

# Create main.go entry point (minimal)
echo "📝 Creating minimal main.go..."
cat > cmd/server/main.go << 'EOF'
package main

import (
	"log"
	"net/http"

	"github.com/client-ux/internal/api/routes"
	"github.com/client-ux/internal/config"
)

func main() {
	cfg := config.Load()
	
	router := routes.SetupRoutes()
	
	log.Printf("🚀 CLIENT-UX server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
EOF

# Create config package
echo "📝 Creating config package..."
cat > internal/config/config.go << 'EOF'
package config

import (
	"os"
)

type Config struct {
	Port           string
	TessDataPrefix string
	StaticDir      string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "3000"),
		TessDataPrefix: getEnv("TESSDATA_PREFIX", "/opt/homebrew/share/tessdata"),
		StaticDir:      getEnv("STATIC_DIR", "static"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
EOF

# Create router
echo "📝 Creating router..."
cat > internal/api/routes/router.go << 'EOF'
package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/client-ux/internal/api/handlers"
	"github.com/client-ux/internal/api/middleware"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	
	// Middleware
	r.Use(middleware.CORS)
	r.Use(middleware.Logging)
	r.Use(middleware.GDPR)
	
	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	api.HandleFunc("/ontology", handlers.GetOntology).Methods("GET")
	api.HandleFunc("/process-document", handlers.ProcessDocument).Methods("POST")
	api.HandleFunc("/save-session", handlers.SaveSession).Methods("POST")
	api.HandleFunc("/sessions/{id}", handlers.GetSession).Methods("GET")
	
	// Static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	
	return r
}
EOF

# Create handlers
echo "📝 Creating handlers..."
cat > internal/api/handlers/health.go << 'EOF'
package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "2.0.0",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
EOF

cat > internal/api/handlers/ontology.go << 'EOF'
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/client-ux/internal/services/ontology"
)

func GetOntology(w http.ResponseWriter, r *http.Request) {
	ontologyService := ontology.NewService()
	
	data, err := ontologyService.GetFormDefinitions()
	if err != nil {
		http.Error(w, "Failed to load ontology", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
EOF

cat > internal/api/handlers/documents.go << 'EOF'
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/client-ux/internal/services/ocr"
)

func ProcessDocument(w http.ResponseWriter, r *http.Request) {
	ocrService := ocr.NewService()
	
	result, err := ocrService.ProcessUpload(r)
	if err != nil {
		http.Error(w, "OCR processing failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
EOF

cat > internal/api/handlers/sessions.go << 'EOF'
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/client-ux/internal/services/session"
)

func SaveSession(w http.ResponseWriter, r *http.Request) {
	sessionService := session.NewService()
	
	result, err := sessionService.Save(r)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func GetSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]
	
	sessionService := session.NewService()
	
	session, err := sessionService.Get(sessionID)
	if err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}
EOF

# Create middleware
echo "📝 Creating middleware..."
cat > internal/api/middleware/cors.go << 'EOF'
package middleware

import (
	"net/http"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}
EOF

cat > internal/api/middleware/logging.go << 'EOF'
package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		next.ServeHTTP(w, r)
		
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}
EOF

cat > internal/api/middleware/gdpr.go << 'EOF'
package middleware

import (
	"net/http"
)

func GDPR(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add GDPR compliance headers
		w.Header().Set("X-GDPR-Compliant", "true")
		w.Header().Set("X-Data-Protection", "field-level")
		
		// Log data access for audit trail
		if r.Method == "POST" || r.Method == "PUT" {
			// TODO: Implement audit logging
		}
		
		next.ServeHTTP(w, r)
	})
}
EOF

# Create service stubs
echo "📝 Creating service stubs..."
cat > internal/services/ocr/service.go << 'EOF'
package ocr

import (
	"net/http"
)

type Service struct {
	// OCR service dependencies
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ProcessUpload(r *http.Request) (interface{}, error) {
	// TODO: Move OCR logic from document_processor.go here
	return map[string]interface{}{
		"status": "processing",
	}, nil
}
EOF

cat > internal/services/ontology/service.go << 'EOF'
package ontology

type Service struct {
	// Ontology service dependencies
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetFormDefinitions() (interface{}, error) {
	// TODO: Move TTL parsing logic from ttl_parser.go here
	return map[string]interface{}{
		"drivers":  map[string]interface{}{"fields": []interface{}{}},
		"vehicles": map[string]interface{}{"fields": []interface{}{}},
		"claims":   map[string]interface{}{"fields": []interface{}{}},
	}, nil
}
EOF

cat > internal/services/session/service.go << 'EOF'
package session

import (
	"net/http"
)

type Service struct {
	// Session service dependencies
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Save(r *http.Request) (interface{}, error) {
	// TODO: Move session logic here
	return map[string]interface{}{
		"status": "saved",
	}, nil
}

func (s *Service) Get(sessionID string) (interface{}, error) {
	// TODO: Move session retrieval logic here
	return map[string]interface{}{
		"id": sessionID,
	}, nil
}
EOF

# Create models
echo "📝 Creating models..."
cat > internal/models/driver.go << 'EOF'
package models

import "time"

type Driver struct {
	ID                 string       `json:"id"`
	Classification     string       `json:"classification"`
	FirstName          string       `json:"firstName"`
	LastName           string       `json:"lastName"`
	DateOfBirth        time.Time    `json:"dateOfBirth"`
	Email              string       `json:"email"`
	Phone              string       `json:"phone"`
	LicenceNumber      string       `json:"licenceNumber"`
	LicenceIssueDate   time.Time    `json:"licenceIssueDate"`
	LicenceExpiryDate  time.Time    `json:"licenceExpiryDate"`
	LicenceValidUntil  time.Time    `json:"licenceValidUntil"`
	Convictions        []Conviction `json:"convictions"`
}

type Conviction struct {
	ID           string    `json:"id"`
	Date         time.Time `json:"date"`
	OffenceCode  string    `json:"offenceCode"`
	Description  string    `json:"description"`
	PenaltyPoints int      `json:"penaltyPoints"`
	FineAmount   float64   `json:"fineAmount"`
}
EOF

cat > internal/models/vehicle.go << 'EOF'
package models

type Vehicle struct {
	ID              string   `json:"id"`
	Registration    string   `json:"registration"`
	Make            string   `json:"make"`
	Model           string   `json:"model"`
	Year            int      `json:"year"`
	EngineSize      string   `json:"engineSize"`
	FuelType        string   `json:"fuelType"`
	Transmission    string   `json:"transmission"`
	EstimatedValue  float64  `json:"estimatedValue"`
	Modifications   []string `json:"modifications"`
}
EOF

cat > internal/models/session.go << 'EOF'
package models

import "time"

type Session struct {
	ID        string    `json:"id"`
	Language  string    `json:"language"`
	Drivers   []Driver  `json:"drivers"`
	Vehicles  []Vehicle `json:"vehicles"`
	Claims    Claims    `json:"claims"`
	Policy    Policy    `json:"policy"`
	Documents []Document `json:"documents"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Claims struct {
	Claims    []Claim    `json:"claims"`
	Accidents []Accident `json:"accidents"`
}

type Claim struct {
	ID          string  `json:"id"`
	Date        string  `json:"date"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Settled     bool    `json:"settled"`
}

type Accident struct {
	ID            string  `json:"id"`
	Date          string  `json:"date"`
	Type          string  `json:"type"`
	Description   string  `json:"description"`
	EstimatedCost float64 `json:"estimatedCost"`
	FaultClaim    bool    `json:"faultClaim"`
}

type Policy struct {
	StartDate    string  `json:"startDate"`
	CoverType    string  `json:"coverType"`
	Excess       float64 `json:"excess"`
	NCDYears     int     `json:"ncdYears"`
	NCDProtected bool    `json:"ncdProtected"`
}

type Document struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Path string `json:"path"`
}
EOF

echo "✅ Backend restructuring complete!"
echo ""
echo "📋 New Go Architecture:"
echo "  cmd/server/main.go              # Entry point (50 lines)"
echo "  internal/api/handlers/          # HTTP handlers (100-150 lines each)"
echo "  internal/api/middleware/        # Middleware (50-100 lines each)"
echo "  internal/api/routes/            # Route configuration"
echo "  internal/services/              # Business logic (200-300 lines each)"
echo "  internal/models/                # Data structures"
echo "  internal/config/                # Configuration"
echo ""
echo "🎯 Next Steps:"
echo "1. Move logic from main.go to appropriate handlers"
echo "2. Move OCR logic from document_processor.go to services/ocr/"
echo "3. Move TTL parsing from ttl_parser.go to services/ontology/"
echo "4. Update imports and test compilation"
echo "5. Add proper error handling and logging"
echo ""
echo "⚠️  Manual Migration Required:"
echo "- Extract actual business logic from existing files"
echo "- Update import paths"
echo "- Add dependency injection"
echo "- Implement proper error handling"
