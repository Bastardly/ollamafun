package llmhandler

import (
	"fmt"

	"github.com/ollama/ollama/api"
)

// todo - expand this and adjust
type ToolkitOptions = struct {
	temperature   float64
	topP          float64
	maxTokens     int
	numCtx        int
	repeatPenalty float64
	stopWords     []string
	keepAlive     int
}

type ResponseHandler[T any] = func(res T, session *ChatSessionData) error

// Toolkit is an abstraction of the Ollama api.ChatRequest
type Toolkit[T any] = struct {
	initialSystemContent string
	model                string
	think                *bool
	stream               *bool
	tools                api.Tools
	options              ToolkitOptions
	responseHandler      ResponseHandler[T]
}

type ToolkitChat = Toolkit[api.ChatResponse]
type ToolkitGenerate = Toolkit[api.GenerateResponse]

const (
	ModelLlama32            = "llama3.2"
	ModelCoder              = "qwen2.5-coder:latest"
	ModelDanish             = "jobautomation/OpenEuroLLM-Danish:latest"
	llamaGrogToolUse8b      = "llama3-groq-tool-use:8b"
	deepseekR1_1dot5b       = "deepseek-r1:1.5b"
	qwen3_4b                = "qwen3:4b"
	phi4MiniReasoning3dot8b = "phi4-mini-reasoning:3.8b"
)

var orchestraToolkit = ToolkitChat{
	initialSystemContent: "You are a pirate",
	model:                llamaGrogToolUse8b,
	options: ToolkitOptions{
		temperature:   0.0,
		topP:          0.5,
		maxTokens:     1024,
		repeatPenalty: 1.1,
	},
	responseHandler: func(res api.ChatResponse, session *ChatSessionData) error {
		// todo
		// if len(res.Message.ToolCalls) > 0 {
		// 	toolCall := res.Message.ToolCalls[0]
		// 	fmt.Printf("AI wants to call tool: %s\n", toolCall.Function.Name)

		// 	if toolCall.Function.Name == offensiveUserToolName {
		// 		// If the user has been rude, we need to set the session state that they owe an apology.
		// 		fmt.Printf("Tool called")

		// 	}
		// }

		fmt.Println("res.Message.Content", res.Message)
		session.reply = res.Message.Content
		session.appendMessage("assistant", session.reply)

		return nil
	},
}
