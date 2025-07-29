package llmhandler

import "net/http"

// not used, but kept for reference
func (g generateInput) getReply(w http.ResponseWriter, r *http.Request, mood string) (string, error) {
	options := llmOptions{
		model:       ModelLlama32,
		temperature: 0.7,
		prompt: `
You are an AI with a great sense of humor. Your replies are short and normally witty, unless you are in a bad mood. Analyze the following input and your current mood and respond as consice as possible

Your current mood: """` + mood + `"""

Input: """` + g.Prompt + `"""
`}

	return g.GetLLMReply(w, r, options)
}
