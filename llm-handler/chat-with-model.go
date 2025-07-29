package llmhandler

import (
	"context"
	"net/http"
	"time"

	"github.com/ollama/ollama/api"
)

func (g generateInput) chatWithModel(w http.ResponseWriter, r *http.Request) (string, error) {
	session.mu.Lock()
	defer session.mu.Unlock()

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	// Add the current mood to the session
	session.Messages = append(session.Messages, api.Message{
		Role:    "system",
		Content: "You know that the user's current intent is: " + session.CurrentUserIntent,
	})

	// Add the new user message
	session.Messages = append(session.Messages, api.Message{
		Role:    "user",
		Content: g.Prompt,
	})

	req := &api.ChatRequest{
		Model:    ModelLlama32,
		Messages: session.Messages,
		Options: map[string]interface{}{
			"temperature": 0.7,
			"max_tokens":  1024,
		},
	}

	var reply string
	err := session.Client.Chat(ctx, req, func(res api.ChatResponse) error {
		reply += res.Message.Content
		return nil
	})

	// Save the assistant reply to the session
	session.Messages = append(session.Messages, api.Message{
		Role:    "assistant",
		Content: reply,
	})

	return reply, err
}
