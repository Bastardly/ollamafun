package llmhandler

import (
	"sync"

	"github.com/ollama/ollama/api"
)

type Session struct {
	CurrentUserIntent string
	Messages          []api.Message
	mu                sync.Mutex
	Client            *api.Client
}

func NewSession() *Session {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		panic("Failed to initialize Ollama client: " + err.Error())
	}
	return &Session{
		Client:            client,
		CurrentUserIntent: "unknown",
		Messages: []api.Message{
			{
				Role:    "system",
				Content: initialSystemContent,
			},
		},
	}
}

var session = NewSession()
