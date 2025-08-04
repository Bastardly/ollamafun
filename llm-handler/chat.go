package llmhandler

import (
	"encoding/json"
	"net/http"
)

type generateInput struct {
	Prompt string `json:"prompt"`
	Method string `json:"method"` // currently not used.
}

func chatWithModel(w http.ResponseWriter, r *http.Request, prompt, sessionID string) error {
	sessions[sessionID].mu.Lock()
	defer sessions[sessionID].mu.Unlock()

	// Add the new user message
	sessions[sessionID].appendMessage("user", prompt)
	return sessions[sessionID].getChatReply(r, orchestraToolkit, sessionID)
}

const ChatSessionName = "chat_session"

func Chat(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(ChatSessionName)

	if err != nil {
		sendErrorResponse(w, "failed to get chat session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if cookie.Value == "" {
		sendErrorResponse(w, "chat sessionID is empty: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var input generateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendErrorResponse(w, "invalid JSON body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sessionID := cookie.Value
	if _, ok := sessions[sessionID]; !ok {
		sessions[sessionID] = createChatSessionData(orchestraToolkit.primeMessages)
	}

	if err := chatWithModel(w, r, input.Prompt, sessionID); err != nil {
		sendErrorResponse(w, "failed to generate chat response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sendReplyResponse(w, sessions[sessionID].reply)
}
