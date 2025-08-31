package handlers

import (
	"net/http"

	"client-ux/internal/services/ai_validation"
)

func ValidateAIInput(w http.ResponseWriter, r *http.Request) {
	service := ai_validation.NewService()
	service.HandleValidation(w, r)
}
