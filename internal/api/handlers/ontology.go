package handlers

import (
	"encoding/json"
	"net/http"

	"client-ux/internal/services/ontology"
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
