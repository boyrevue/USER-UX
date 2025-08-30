package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/client-ux/internal/services/ocr"
)

func ProcessDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ocrService := ocr.NewService()

	result, err := ocrService.ProcessUpload(r)
	if err != nil {
		http.Error(w, "OCR processing failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(result)
}
