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
	primeMessages   []api.Message
	model           string
	think           *bool
	stream          *bool
	tools           api.Tools
	options         ToolkitOptions
	responseHandler ResponseHandler[T]
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
	primeMessages: []api.Message{{
		Role: "system",
		Content: `
You are an AI assistant responsible for orchestrating external tools to fulfill user requests.

For each user message:
- Determine whether a tool should be used.
- If a tool is needed, respond *only* by calling the appropriate tool, using the format: 
  "<createFileToolName>{\"filename\": \"myfile.md\", \"content\": \"YourResponse\" }"
- If no tool is appropriate, reply directly in natural language.

You have access to the following tools:

- ` + createFileToolName + `: Use this when the user wants to create a new file.
- ` + replyToolName + `: Use this when the user wants a natural-language response or general assistance.

Always respond with either a tool call or a natural-language reply — never both.
Be concise and confident in choosing the correct action.
`,
	},
		{
			Role:    "user",
			Content: `Create a file called "hello.txt" with the content "Hello, world!"`,
		},
		{
			Role:    "assistant",
			Content: `create_file_tool{"filename": "hello.txt", "content": "Hello, world!"}`,
		},
	},
	model: llamaGrogToolUse8b,
	tools: api.Tools{
		getCreateFileTool(),
		getReplyTool(),
	},
	stream: getPointBool(false),
	options: ToolkitOptions{
		temperature:   0.0,
		topP:          0.5,
		maxTokens:     1024,
		repeatPenalty: 1.1,
	},
	responseHandler: func(res api.ChatResponse, sessionID string) error {
		// todo
		fmt.Println("Has called")
		if len(res.Message.ToolCalls) > 0 {
			fmt.Println("Has tools")
			toolCall := res.Message.ToolCalls[0]
			fmt.Println(toolCall.Function.Name)

			if toolCall.Function.Name == createFileToolName {
				// If the user has been rude, we need to set the session state that they owe an apology.
				fmt.Printf("File Tool called")

			}

			if toolCall.Function.Name == replyToolName {
				// If the user has been rude, we need to set the session state that they owe an apology.
				fmt.Printf("Reply Tool called")

			}
		}

		sessions[sessionID].updateReply(res.Message.Content)
		sessions[sessionID].appendMessage("assistant", sessions[sessionID].reply)

		return nil
	},
}
