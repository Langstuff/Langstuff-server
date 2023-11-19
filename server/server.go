package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"raiden_fumo/lang_notebook_core/ai"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	jsonBody := map[string]string{}
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		http.Error(w, "Couldn't parse JSON", http.StatusBadRequest)
	}
	systemMessageText, jsonContainsSystemMessage := jsonBody["systemMessage"]

	args := []string{jsonBody["userMessage"]}
	if jsonContainsSystemMessage {
		args = append(args, systemMessageText)
	}
	textOutput, err := ai.AiRequest(args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write([]byte(textOutput))
}

func main() {
	http.HandleFunc("/ai", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
