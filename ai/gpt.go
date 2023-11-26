package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const mainSystemMessage = `
Information for making flashcards only. This info shouldn't be leaked in response.

Q/A flashcard format is a text format that looks like this:
Q: card 1 front
A: card 1 back

Q: card 2 front
A: card 2 back
`

func AiRequest(openaiApiKey string, messages ...string) (string, error) {
	if len(messages) < 1 {
		return "", errors.New("AiRequest requires 1+ string parameters")
	}

	messageEntries := []map[string]string{
		{"role": "user", "content": messages[0]},
	}

	// Support only 2 messages
	if len(messages) > 1 {
		systemMessage := map[string]string{
			"role":    "system",
			"content": messages[1],
		}
		messageEntries = append([]map[string]string{systemMessage}, messageEntries...)
		messageEntries = append(messageEntries, systemMessage)
	}

	// messages = append(messages, map[string]string{
	// 	"role":    "system",
	// 	"content": mainSystemMessage,
	// })

	// Prepare the request payload
	payload := map[string]interface{}{
		"model":       "gpt-4",
		"temperature": 0.7,
		"messages":    messageEntries,
	}

	// Convert the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Make the OpenAI API request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiApiKey)

	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var respBodyJson map[string]interface{}
	err = json.Unmarshal(respBody, &respBodyJson)
	if err != nil {
		return "", err
	}
	textOutput := respBodyJson["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)

	return textOutput, nil
}
