package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

// BaseLLM defines a minimal provider-agnostic interface
// for generating responses.
type BaseLLM interface {
	Generate(ctx context.Context, system string, user string) (string, error)
}

// OpenAI implements BaseLLM using the Chat Completions API.
type OpenAI struct {
	APIKey string
	Model  string
}

type oaMessage struct{ Role, Content string }

type oaBody struct {
	Model       string      `json:"model"`
	Messages    []oaMessage `json:"messages"`
	Temperature float64     `json:"temperature"`
}

type oaChoice struct{ Message oaMessage }

type oaResp struct{ Choices []oaChoice }

func NewOpenAI() *OpenAI {
	key := os.Getenv("OPENAI_API_KEY")
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &OpenAI{APIKey: key, Model: model}
}

func (o *OpenAI) Generate(ctx context.Context, system string, user string) (string, error) {
	body := oaBody{
		Model:       o.Model,
		Messages:    []oaMessage{{Role: "system", Content: system}, {Role: "user", Content: user}},
		Temperature: 0.2,
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+o.APIKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var out oaResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", nil
	}
	return out.Choices[0].Message.Content, nil
}
