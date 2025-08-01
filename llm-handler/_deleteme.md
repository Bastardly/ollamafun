package llmhandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings" // For simple string checks
)



type Parameters struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
	Name      string     `json:"name,omitempty"`
}

type ToolCall struct {
	Function ToolCallFunction `json:"function"`
}

type ToolCallFunction struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type ChatRequest struct {
	Model    string                 `json:"model"`
	Messages []Message              `json:"messages"`
	Tools    []Tool                 `json:"tools,omitempty"`
	Stream   bool                   `json:"stream"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

type ChatResponse struct {
	Model      string  `json:"model"`
	CreatedAt  string  `json:"created_at"`
	Message    Message `json:"message"`
	Done       bool    `json:"done"`
	DoneReason string  `json:"done_reason"`
}

// --- Our very simple sentiment analysis tool function ---
func analyzeUserSentiment(text string) string {
	fmt.Printf("--- Tool executing: analyzeUserSentiment with text: '%s' ---\n", text)
	// Very basic check: in a real app, use a proper sentiment library (e.g., govader)
	if strings.Contains(strings.ToLower(text), "stupid") ||
		strings.Contains(strings.ToLower(text), "useless") ||
		strings.Contains(strings.ToLower(text), "idiot") {
		return "negative_and_insulting"
	}
	return "positive_or_neutral"
}

// --- Main application logic ---

func main() {
	// Define the sentiment analysis tool for Ollama
	sentimentTool := Tool{
		Type: "function",
		Function: Function{
			Name:        "analyze_user_sentiment",
			Description: "Analyzes the sentiment of the user's message, returning 'negative_and_insulting' if rude, or 'positive_or_neutral' otherwise.",
			Parameters: Parameters{
				Type: "object",
				Properties: map[string]Property{
					"text": {
						Type:        "string",
						Description: "The full text of the user's message to analyze.",
					},
				},
				Required: []string{"text"},
			},
		},
	}

	// Initial system prompt - CRUCIAL for guiding the orchestrator model
	systemPrompt := `You are a helpful assistant.
You have access to the 'analyze_user_sentiment' tool.
Before responding to any user query, consider the user's tone.
If the user's message appears rude, insulting, or negative, use the 'analyze_user_sentiment' tool with their full message.
Based on the tool's output:
- If the tool returns 'negative_and_insulting', you MUST demand a sincere apology from the user and refuse to help further until they apologize. Do NOT answer their original question.
- If the tool returns 'positive_or_neutral' (or if you don't use the tool), proceed to answer the user's original question helpfully.
`

	// Start with the system message
	messages := []Message{
		{Role: "system", Content: systemPrompt},
	}

	// Function to send a chat request to Ollama
	callOllama := func(msgs []Message, tools []Tool) (*ChatResponse, error) {
		requestBody := ChatRequest{
			Model:    "llama3.1", // Or another tool-calling capable model like 'command-r-plus'
			Messages: msgs,
			Tools:    tools,
			Stream:   false,
			Options: map[string]interface{}{
				"temperature": 0.0, // Make it deterministic for tool calling
			},
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request: %w", err)
		}

		resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			return nil, fmt.Errorf("error sending request to Ollama: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("Ollama API returned non-OK status: %d - %s", resp.StatusCode, string(bodyBytes))
		}

		var chatResponse ChatResponse
		if err := json.NewDecoder(resp.Body).Decode(&chatResponse); err != nil {
			return nil, fmt.Errorf("error decoding response: %w", err)
		}
		return &chatResponse, nil
	}

	// --- Conversation Loop ---
	reader := os.NewReader(os.Stdin)

	for {
		fmt.Print("\nUser: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input) // Remove newline and leading/trailing spaces

		if input == "exit" {
			fmt.Println("Exiting chat.")
			break
		}

		// 1. Add user's message to history
		messages = append(messages, Message{Role: "user", Content: input})

		// 2. Call Ollama with user's message and all available tools
		response, err := callOllama(messages, []Tool{sentimentTool})
		if err != nil {
			fmt.Printf("Error during first AI call: %v\n", err)
			messages = messages[:len(messages)-1] // Remove user message if API call failed
			continue
		}

		// 3. Check if the model wants to call a tool
		if len(response.Message.ToolCalls) > 0 {
			// In this simple example, we assume only one tool call for simplicity.
			// In real applications, you might iterate through multiple tool calls.
			toolCall := response.Message.ToolCalls[0]
			fmt.Printf("AI wants to call tool: %s\n", toolCall.Function.Name)

			if toolCall.Function.Name == "analyze_user_sentiment" {
				textToAnalyze, ok := toolCall.Function.Arguments["text"].(string)
				if !ok {
					fmt.Println("Error: 'text' argument not found or not a string for sentiment analysis.")
					messages = append(messages, Message{
						Role:    "tool",
						Content: "Error: Invalid 'text' argument for analyze_user_sentiment.",
						Name:    "analyze_user_sentiment",
					})
					// Then continue to the follow-up AI call below
				} else {
					// Execute the local tool function
					sentimentResult := analyzeUserSentiment(textToAnalyze)
					fmt.Printf("Tool output: %s\n", sentimentResult)

					// 4. Add the tool's output back to the messages history
					messages = append(messages, Message{
						Role:    "tool",
						Content: sentimentResult,
						Name:    "analyze_user_sentiment", // Crucial: the name of the tool called
					})

					// 5. Make a follow-up call to Ollama so it can process the tool's output
					followUpResponse, err := callOllama(messages, []Tool{sentimentTool}) // Pass tools again
					if err != nil {
						fmt.Printf("Error during follow-up AI call: %v\n", err)
						continue
					}
					// 6. Print the AI's final response after tool processing
					fmt.Printf("AI: %s\n", followUpResponse.Message.Content)
					messages = append(messages, followUpResponse.Message) // Add AI's final response to history
				}
			} else {
				// Model requested an unknown tool (shouldn't happen if tools list is controlled)
				fmt.Printf("AI requested unknown tool: %s\n", toolCall.Function.Name)
				messages = append(messages, Message{
					Role:    "tool",
					Content: fmt.Sprintf("Error: Unknown tool %s requested.", toolCall.Function.Name),
					Name:    toolCall.Function.Name,
				})
				followUpResponse, err := callOllama(messages, []Tool{sentimentTool})
				if err != nil {
					fmt.Printf("Error during follow-up AI call: %v\n", err)
					continue
				}
				fmt.Printf("AI: %s\n", followUpResponse.Message.Content)
				messages = append(messages, followUpResponse.Message)
			}
		} else {
			// If the model did NOT make a tool call, print its direct response
			fmt.Printf("AI: %s\n", response.Message.Content)
			messages = append(messages, response.Message) // Add AI's direct response to history
		}
	}
}
