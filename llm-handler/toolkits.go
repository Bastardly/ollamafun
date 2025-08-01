package llmhandler

import (
	"fmt"

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
//
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

type ResponseHandler[T any] = func(res T, sessionID string) error

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
	stream:               getPointBool(false),
	options: ToolkitOptions{
		temperature:   0.0,
		topP:          0.5,
		maxTokens:     1024,
		repeatPenalty: 1.1,
	},
	responseHandler: func(res api.ChatResponse, sessionID string) error {
		// todo
		// if len(res.Message.ToolCalls) > 0 {
		// 	toolCall := res.Message.ToolCalls[0]

		// 	if toolCall.Function.Name == offensiveUserToolName {
		// 		// If the user has been rude, we need to set the session state that they owe an apology.
		// 		fmt.Printf("Tool called")

		// 	}
		// }

		sessions[sessionID].updateReply(res.Message.Content)
		sessions[sessionID].appendMessage("assistant", sessions[sessionID].reply)

		return nil
	},
}
