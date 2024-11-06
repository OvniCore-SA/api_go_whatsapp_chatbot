package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/config"
)

type OllamaService struct{}

func NewOllamaService() *OllamaService {
	return &OllamaService{}
}

func (service *OllamaService) SendMessageToChat(userMessage string, promptForOllama string) (string, error) {
	url := "https://oi.telco.com.ar/ollama/api/chat"
	token := config.OLLAMA_TOKEN

	// Construcción del cuerpo de la solicitud
	body := map[string]interface{}{
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": promptForOllama,
			},
			{
				"role":    "user",
				"content": userMessage,
			},
		},
		"model": "codestral:22b",
		"options": map[string]interface{}{
			"temperature": 0.8,
			"num_predict": 100,
		},
	}

	// Serialización del cuerpo a JSON
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("error serializing body: %v", err)
	}

	// Creación de la solicitud HTTP
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Cliente HTTP
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	// Ejecución de la solicitud
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Verificación del estado HTTP
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Lectura del stream de respuesta
	fullResponse := ""
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		// Procesamiento de cada fragmento de la respuesta
		line := scanner.Text()
		var jsonResponse map[string]interface{}
		if err := json.Unmarshal([]byte(line), &jsonResponse); err != nil {
			return "", fmt.Errorf("error parsing response: %v", err)
		}

		// Verificación si el mensaje contiene contenido
		if message, exists := jsonResponse["message"]; exists {
			if msgMap, ok := message.(map[string]interface{}); ok {
				if content, ok := msgMap["content"].(string); ok && content != "" {
					fullResponse += content
				}
			}
		}

		// Si el campo "done" es true, se termina la lectura
		if done, exists := jsonResponse["done"].(bool); exists && done {
			break
		}
	}

	// Verificación de errores de escaneo
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	// Verificación si la respuesta está vacía
	if fullResponse == "" {
		return "", errors.New("empty response from assistant")
	}

	// Retornar la respuesta completa
	return fullResponse, nil
}
