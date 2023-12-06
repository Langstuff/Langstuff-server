package main

import (
	"encoding/json"
	"io"
	"net/http"
	"raiden_fumo/lang_notebook_core/ai"
)

type AiRequestHandler struct {
	openaiApiKey string
}

func (handler AiRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	jsonBody := map[string]string{}
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		http.Error(w, "Couldn't parse request data", http.StatusBadRequest)
		return
	}
	systemMessageText, jsonContainsSystemMessage := jsonBody["systemMessage"]

	args := []string{jsonBody["userMessage"]}
	if jsonContainsSystemMessage {
		args = append(args, systemMessageText)
	}
	textOutput, err := ai.AiRequest(handler.openaiApiKey, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write([]byte(textOutput))
}

func makeAiRequestServer(openaiApiKey string) *AiRequestHandler {
	return &AiRequestHandler{openaiApiKey: openaiApiKey}
}
