package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"client-ux/internal/services/llm"
)

type AssistantHandler struct{ model llm.BaseLLM }

func NewAssistantHandler() *AssistantHandler { return &AssistantHandler{model: llm.NewOpenAI()} }

type AssistantRequest struct {
	Field     string `json:"field"`
	Prompt    string `json:"prompt"`
	UserInput string `json:"userInput"`
}

type AssistantResponse struct {
	Reply string `json:"reply"`
}

// simple redactor for emails/postcodes/dates
func redact(s string) string {
	s = strings.ReplaceAll(s, "@", "[at]")
	s = strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return '*'
		}
		return r
	}, s)
	return s
}

func ensureLogDir() string {
	dir := filepath.Join("logs")
	_ = os.MkdirAll(dir, 0o755)
	return filepath.Join(dir, "assistant.log")
}

func appendLog(line string) { _ = os.WriteFile(ensureLogDir(), []byte(line+"\n"), 0o644) }

func (h *AssistantHandler) Generate(w http.ResponseWriter, r *http.Request) {
	var req AssistantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	system := "You are an insurance form assistant. Be concise, structured, and request only information relevant to underwriting."
	user := req.Prompt + "\nUser: " + req.UserInput
	reply, err := h.model.Generate(ctx, system, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// redacted file log
	ts := time.Now().Format(time.RFC3339)
	logLine := ts + "\tfield=" + req.Field + "\tuser=\"" + redact(req.UserInput) + "\"\treply=\"" + redact(reply) + "\""
	appendLog(logLine)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(AssistantResponse{Reply: reply})
}
