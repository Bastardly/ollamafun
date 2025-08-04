package llmhandler

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/ollama/ollama/api"
)

type ChatSessionData struct {
	reply    string
	messages []api.Message
	mu       sync.Mutex
}

// todo - This won't scale at all. At least it needs some cleanup if used public
var sessions = map[string]*ChatSessionData{}

var client = getClient()

func (s *ChatSessionData) updateReply(reply string) {
	s.reply = reply
}

func (s *ChatSessionData) appendMessage(role, content string) {
	s.messages = append(s.messages, api.Message{
		Role:    role,
		Content: content,
	})
}

// generateReply generates a reply based on chat history (messages)
func (s *ChatSessionData) getChatReply(r *http.Request, toolkit ToolkitChat, sessionID string) error {
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	chatReq := &api.ChatRequest{
		Messages: sessions[sessionID].messages,
		Tools:    toolkit.tools,
		Model:    toolkit.model,
		Think:    toolkit.think,
		Stream:   toolkit.stream,
		Options: map[string]interface{}{
			"temperature":    toolkit.options.temperature,
			"top_p":          toolkit.options.topP,
			"max_tokens":     toolkit.options.maxTokens,
			"repeat_penalty": toolkit.options.repeatPenalty,
			// "num_ctx":     toolkit.options.numCtx,
			// "stop":        toolkit.options.stopWords,
			// "keep_alive":  toolkit.options.keepAlive,
		}}

	return client.Chat(ctx, chatReq, func(res api.ChatResponse) error {
		return toolkit.responseHandler(res, sessionID)
	})
}

// generateReply generates a single reply
func (s *ChatSessionData) generateReply(w http.ResponseWriter, r *http.Request, toolkit ToolkitGenerate, prompt string, sessionID string) (string, error) {
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	req := &api.GenerateRequest{
		Model:  toolkit.model,
		Prompt: prompt,
		Stream: func(b bool) *bool { return &b }(false),
		Options: map[string]interface{}{
			"temperature": toolkit.options.temperature,
			"max_tokens":  toolkit.options.maxTokens,
		},
	}

	var reply string
	err := client.Generate(ctx, req, func(res api.GenerateResponse) error {
		return toolkit.responseHandler(res, sessionID)
	})

	return reply, err
}
