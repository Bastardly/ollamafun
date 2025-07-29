package llmhandler

import (
	"context"
	"net/http"
	"time"

	"github.com/ollama/ollama/api"
)

type generateInput struct {
	Prompt string `json:"prompt"`
	Method string `json:"method"`
}

const (
	ModelLlama32 = "llama3.2"
	ModelCoder   = "qwen2.5-coder:latest"
	ModelDanish  = "jobautomation/OpenEuroLLM-Danish:latest"
)

type llmOptions struct {
	model       string
	prompt      string
	temperature float32
}

// GetLLMReply generates a reply using the specified LLM options. That is fresh on every new input since model holds no memory.
func (g generateInput) GetLLMReply(w http.ResponseWriter, r *http.Request, options llmOptions) (string, error) {
	session.mu.Lock()
	defer session.mu.Unlock()
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	req := &api.GenerateRequest{
		Model:  options.model,
		Prompt: options.prompt,
		Stream: func(b bool) *bool { return &b }(false),
		Options: map[string]interface{}{
			"temperature": options.temperature,
			"max_tokens":  2000,
		},
	}

	var reply string
	err := session.Client.Generate(ctx, req, api.GenerateResponseFunc(func(res api.GenerateResponse) error {
		reply += res.Response
		return nil
	}))

	return reply, err
}
