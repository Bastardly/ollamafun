package llmhandler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ollama/ollama/api"
)

const (
	ModelLlama32 = "llama3.2"
	ModelPhi3    = "phi3:latest"
	ModelMulti   = "llava"
)

func HandleGenerate(w http.ResponseWriter, r *http.Request) {
	var input generateInput
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
		Model:  ModelLlama32,
		Prompt: input.getInsultJSON(),
		Stream: func(b bool) *bool { return &b }(false),
		Options: map[string]interface{}{
			"temperature": 0.8,
			"max_tokens":  2000,
		},
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
