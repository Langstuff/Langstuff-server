package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

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


func ConvertTextToSpeech(openaiApiKey string, text string, voice string) ([]byte, error) {
	fmt.Println("key: ", openaiApiKey)
	url := "https://api.openai.com/v1/audio/speech"

	payload := struct {
		Model string `json:"model"`
		Input string `json:"input"`
		Voice string `json:"voice"`
	}{
		Model: "tts-1",
		Input: text,
		Voice: voice,
	}

	payloadBytes, err := json.Marshal(payload)

	fmt.Println("payloadBytes: ", string(payloadBytes))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+openaiApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("request failed with status: " + resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
