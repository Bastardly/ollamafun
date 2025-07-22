package llmhandler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ollama/ollama/api"
)

type GenerateInput struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

func HandleGenerate(w http.ResponseWriter, r *http.Request) {
	var input GenerateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	client, err := api.ClientFromEnvironment()
	if err != nil {
		http.Error(w, "failed to create Ollama client: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	req := &api.GenerateRequest{
		Model:   input.Model,
		Prompt:  input.Prompt,
		Stream:  &input.Stream,
		Options: input.Options,
	}

	var result string
	err = client.Generate(ctx, req, api.GenerateResponseFunc(func(res api.GenerateResponse) error {
		result += res.Response
		return nil
	}))
	if err != nil {
		http.Error(w, "generation error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"response": result,
	})
}
