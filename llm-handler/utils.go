package llmhandler

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
)

func getPointBool(val bool) *bool {
	return &val
}

func getClient() *api.Client {
	baseURL, err := url.Parse("http://127.0.0.1:8003")
	if err != nil {
		panic("Invalid Ollama base URL: " + err.Error())
	}
	httpClient := &http.Client{}
	return api.NewClient(baseURL, httpClient)
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	http.Error(w, message, statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"response": message,
	})
}

func sendReplyResponse(w http.ResponseWriter, reply string) {
	json.NewEncoder(w).Encode(map[string]string{
		"response": reply,
	})
}

func createChatSessionData(primeMessages []api.Message) *ChatSessionData {
	session := &ChatSessionData{}

	for _, msg := range primeMessages {
		session.appendMessage(msg.Role, msg.Content)
	}

	return session
}
