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
