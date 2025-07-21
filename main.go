package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ollama/ollama/api"
)

// REST input structure
type GenerateInput struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"` // e.g., temperature, max_tokens
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	var input GenerateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	client, err := api.ClientFromEnvironment()
	if err != nil {
		http.Error(w, "failed to init Ollama client", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	req := &api.GenerateRequest{
		Model:   input.Model,
		Prompt:  input.Prompt,
		Stream:  &input.Stream,
		Options: input.Options,
	}

	// Set headers for streaming
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if input.Stream {
		// Stream results chunk-by-chunk
		enc := json.NewEncoder(w)
		err := client.Generate(ctx, req, api.GenerateResponseFunc(func(res api.GenerateResponse) error {
			if err := enc.Encode(res); err != nil {
				return fmt.Errorf("streaming encode failed: %w", err)
			}
			// Make sure it's flushed
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			return nil
		}))
		if err != nil {
			log.Printf("error during generation: %v", err)
		}
	} else {
		// Non-stream: collect and return once
		var full string
		err := client.Generate(ctx, req, api.GenerateResponseFunc(func(res api.GenerateResponse) error {
			full += res.Response
			return nil
		}))
		if err != nil {
			http.Error(w, fmt.Sprintf("generation failed: %v", err), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"response": full})
	}
}

func startUI() {
	var err error
	tpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}

	http.HandleFunc("/", handleUI)
}

func main() {
	startUI()
	http.HandleFunc("/generate", handleGenerate)
	port := os.Getenv("PORT")
	if port == "" {
		port = "6789"
	}
	log.Printf("✅ Server running at http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
