package llmhandler

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/ollama/ollama/api"
)

// settings defines configuration options for controlling model behavior.
//
// Fields:
//   - think: Enables or disables thinking mode.
//   - stream: Enables or disables a streaming response.
//   - temperature: Controls response creativity; lower values (e.g., 0.2-0.5) make responses more deterministic and focused.
//   - topP: Limits deviation from core instructions; lower values (e.g., 0.5-0.7) reduce off-topic responses.
//   - maxTokens: Sets the maximum number of output tokens to generate in the response.
//   - numCtx: Sets the context limit (e.g., 100,000 tokens) for considering previous conversation history.
//   - repeatPenalty: Discourages repetition; a value of 1.1 promotes varied responses.
//   - stopWords: Specifies stop words to prevent the model from generating certain phrases or continuing past a point.
//   - keepAlive: Keeps the connection alive indefinitely.

type ChatSessionData struct {
	reply    string
	messages []api.Message
	mu       sync.Mutex
}

type generateInput struct {
	Prompt string `json:"prompt"`
	Method string `json:"method"` // currently not used.
}

// todo - This won't scale. Need some cleanup if used public
var sessions = map[string]*ChatSessionData{}

func getClient() *api.Client {
	baseURL, err := url.Parse("http://127.0.0.1:8003")
	if err != nil {
		panic("Invalid Ollama base URL: " + err.Error())
	}
	httpClient := &http.Client{}
	return api.NewClient(baseURL, httpClient)
}

var client = getClient()

// Just a mock session handler for local single user
func createChatSessionData(initialSystemContent string) *ChatSessionData {
	session := &ChatSessionData{}

	session.appendMessage("system", initialSystemContent)

	return session
}

func (s ChatSessionData) appendMessage(role, content string) {
	s.messages = append(s.messages, api.Message{
		Role:    role,
		Content: content,
	})
}

func (s ChatSessionData) getChatReply(r *http.Request, toolkit ToolkitChat, session *ChatSessionData) error {
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	session.reply = ""

	chatReq := &api.ChatRequest{
		Messages: session.messages,
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
		return toolkit.responseHandler(res, session)
	})
}

// generateReply generates a reply using the specified LLM options. That is fresh on every new input since model holds no memory.
func (s ChatSessionData) generateReply(w http.ResponseWriter, r *http.Request, toolkit ToolkitGenerate, prompt string, session *ChatSessionData) (string, error) {
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
		return toolkit.responseHandler(res, session)
	})

	return reply, err
}
