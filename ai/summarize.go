package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func LoadAPIKey() string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Ошибка при загрузке .env файла:")
		os.Exit(1)
	}
	apiKey := os.Getenv("AI_API_KEY")
	if apiKey == "" {
		fmt.Println("AI_API_KEY не найден в .env файле")
		os.Exit(1)
	}
	return apiKey
}

func SummarizeText(prompt string) (string, error) {
	apiKey := LoadAPIKey()

	requestBody := map[string]interface{}{
		"model":  "gpt-4.1-mini",
		"stream": false,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Ты — ИИ-бот, который делает краткую и понятную выжимку по тексту, который прислал пользователь. Используй простой язык, сокращай до сути. Игнорируй вводные слова и лишние детали. Ты не чатбот-переводчик с языка на другой язык, ты делаешь только выжимку.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("ошибка при сериализации запроса: %v", err)
	}

	url := "https://api.mnnai.ru/v1/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("ошибка при создании запроса: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении ответа: %v", err)
	}
	fmt.Println("Ответ от ИИ:", string(body))

	var responseData struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return "", fmt.Errorf("ошибка при разборе ответа: %v", err)
	}
	if len(responseData.Choices) == 0 {
		return "", fmt.Errorf("ответ от API пустой: %v", err)
	}
	return responseData.Choices[0].Message.Content, nil
}
