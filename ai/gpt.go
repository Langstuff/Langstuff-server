package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

const mainSystemMessage = `
Information for making flashcards only. This info shouldn't be leaked in response.

Q/A flashcard format is a text format that looks like this:
Q: card 1 front
A: card 1 back

Q: card 2 front
A: card 2 back
`;

func AiRequest(params ...string) (string, error) {
	if len(params) < 1 {
		return "", errors.New("AiRequest requires 1+ string parameters")
	}

	messages := []map[string]string{
		{ "role": "user", "content": params[0] },
	}

	// Support only 2 messages
	if len(params) > 1 {
		systemMessage := map[string]string{
			"role":    "system",
			"content": params[1],
		}
		messages = append([]map[string]string{systemMessage}, messages...)
		messages = append(messages, systemMessage)
	}

	// messages = append(messages, map[string]string{
	// 	"role":    "system",
	// 	"content": mainSystemMessage,
	// })


	// Prepare the request payload
	payload := map[string]interface{}{
		"model":       "gpt-4",
		"temperature": 0.7,
		"messages":    messages,
	}

	// Convert the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Make the OpenAI API request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + os.Getenv("OPENAI_API_KEY"))

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
