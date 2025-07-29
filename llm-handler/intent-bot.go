package llmhandler

import "net/http"

func (g generateInput) getUserIntent(w http.ResponseWriter, r *http.Request) (string, error) {
	options := llmOptions{
		model:       ModelLlama32,
		temperature: 0.4,
		prompt: `Analyze the following user input and return a short description of the user's intent, mood, or tone in one sentence.

	Input: """` + g.Prompt + `"""
	`}

	return g.GetLLMReply(w, r, options)
}
