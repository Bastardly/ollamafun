package llmhandler

import "net/http"

// not used, but kept for reference
func (g generateInput) getMood(w http.ResponseWriter, r *http.Request) (string, error) {
	options := llmOptions{
		model:       ModelLlama32,
		temperature: 0.4,
		prompt: `
	You are a mood bot, you will determine your own mood based in the user's input.
	You will remember the mood you are in, and weight that towards your response.
	Your mood with not shift quickly. If the user is rude, you will remember this until the user apologises or somehow makes up for it.

	Analyze the following input and describe your mood as consice as possible.
	The reply should only be a single sentence of your current mood, and why you are in that mood: e.g. 'I am angry because the user called me an idiot should apologize".

	Input: """` + g.Prompt + `"""
	`}

	return g.GetLLMReply(w, r, options)
}
