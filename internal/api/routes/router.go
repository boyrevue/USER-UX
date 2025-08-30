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
