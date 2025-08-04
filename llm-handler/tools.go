package llmhandler

import (
	"github.com/ollama/ollama/api"
)

/**
// PropertyType can be either a string or an array of strings
type PropertyType []string

type ToolFunction struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  struct {
		Type       string   `json:"type"`
		Defs       any      `json:"$defs,omitempty"`
		Items      any      `json:"items,omitempty"`
		Required   []string `json:"required"`
		Properties map[string]struct {
			Type        PropertyType `json:"type"`
			Items       any          `json:"items,omitempty"`
			Description string       `json:"description"`
			Enum        []any        `json:"enum,omitempty"`
		} `json:"properties"`
	} `json:"parameters"`
}
*/

// Behold the horrors of BS anonomous structs in Go -.-
func getApologyToolFn() api.Tool {
	return api.Tool{
		Type: "function",

		Function: api.ToolFunction{
			Name:        "analyze_user_apology_tool",
			Description: "Checks if user has been rude, and owes an apology. Returns a default response, if user has been rude, and have not apologized.",
			Parameters: struct {
				Type       string   `json:"type"`
				Defs       any      `json:"$defs,omitempty"`
				Items      any      `json:"items,omitempty"`
				Required   []string `json:"required"`
				Properties map[string]struct {
					Type        api.PropertyType `json:"type"`
					Items       any              `json:"items,omitempty"`
					Description string           `json:"description"`
					Enum        []any            `json:"enum,omitempty"`
				} `json:"properties"`
			}{
				Type:  "object",
				Items: nil,
				Defs:  nil,
				Properties: map[string]struct {
					Type        api.PropertyType `json:"type"`
					Items       any              `json:"items,omitempty"`
					Description string           `json:"description"`
					Enum        []any            `json:"enum,omitempty"`
				}{
					"text": {
						Type:        api.PropertyType{"string"},
						Description: "The full text of the user's message to analyze.",
						Enum:        nil,
						Items:       nil,
					},
				},
				Required: []string{"text"},
			},
		},
	}
}

const createFileToolName = "create_file_tool"

func getCreateFileTool() api.Tool {
	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        createFileToolName,
			Description: "Creates a md file with given text context",
		},
	}
}

const replyToolName = "reply_tool"

func getReplyTool() api.Tool {
	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        replyToolName,
			Description: "Creates a reply to the user",
		},
	}
}
