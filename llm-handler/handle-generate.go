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

func HandleGenerate(w http.ResponseWriter, r *http.Request) {
	var input generateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	userIntent, err := input.getUserIntent(w, r)

	if err != nil {
		sendErrorResponse(w, "failed to determine user's intent: "+err.Error(), http.StatusInternalServerError)
	}

	println("Current mood:", userIntent)
	result, _ := input.chatWithModel(w, r)

	if err != nil {
		sendErrorResponse(w, "failed to generate chat response: "+err.Error(), http.StatusInternalServerError)
	}

	println("Current result:", result)
	json.NewEncoder(w).Encode(map[string]string{
		"response": result,
	})
}
