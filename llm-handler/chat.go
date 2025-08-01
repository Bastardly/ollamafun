package llmhandler

import (
	"encoding/json"
	"net/http"
)

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	http.Error(w, message, statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"response": message,
	})
}

const ChatSessionName = "chat_session"

func chatWithModel(w http.ResponseWriter, r *http.Request, prompt, sessionID string) error {
	session := sessions[sessionID]
	session.mu.Lock()
	defer session.mu.Unlock()

	// Add the new user message
	session.appendMessage("user", prompt)
	return session.getChatReply(r, orchestraToolkit, session)
}

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
	if sessions[sessionID] == nil {
		sessions[sessionID] = createChatSessionData(orchestraToolkit.initialSystemContent)
	}

	if err := chatWithModel(w, r, input.Prompt, sessionID); err != nil {
		sendErrorResponse(w, "failed to generate chat response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"response": sessions[sessionID].reply,
	})
}
