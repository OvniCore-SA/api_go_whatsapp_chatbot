package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/config"
	"github.com/sashabaranov/go-openai"
)

type OpenAIService struct {
	openAIClient *openai.Client
}

func NewOpenAIService(clientOpenAI *openai.Client) *OpenAIService {
	return &OpenAIService{
		openAIClient: clientOpenAI,
	}
}

type OpenAIRequest struct {
	Model    string      `json:"model"`
	Messages []ChatEntry `json:"messages"`
}

type ChatEntry struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (service *OpenAIService) SendMessageToOpenAI(userMessage string, menuOptions string) (string, error) {

	maxTokensGeneralMessages, err := strconv.Atoi(config.MAX_TOKENS_GENERAL_MESSAGE)
	if err != nil {
		fmt.Println("ERROR maxTokensGeneralMessages: " + err.Error())
		return "", err
	}

	// Crear mensaje de entrada para OpenAI
	prompt := fmt.Sprintf("Usuario dice: %s. Opciones de menú disponibles: %v. Responde cuál opción del menú se debe seleccionar o da una respuesta genérica.", userMessage, menuOptions)

	resp, err := service.openAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: "Actúa como un asistente en un chatbot de WhatsApp."},
				{Role: openai.ChatMessageRoleUser, Content: prompt},
			},
			MaxTokens: maxTokensGeneralMessages,
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, err

}

func (service *OpenAIService) SendMessageBasicToOpenAI(prompt string, promptForOpenAI string) (string, error) {

	_, err := strconv.Atoi(config.MAX_TOKENS_INITIAL_MESSAGE)
	if err != nil {
		fmt.Println("ERROR maxTokensInitialMessages: " + err.Error())
		return "", err
	}

	resp, err := service.openAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: promptForOpenAI},
				{Role: openai.ChatMessageRoleUser, Content: prompt},
			},
			//MaxTokens: maxTokensInitialMessages,
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}
	return resp.Choices[0].Message.Content, err
}

func (service *OpenAIService) SendMessageForHttpToOpenAI(requestBody *OpenAIRequest) (string, error) {
	// Convertir a JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error encoding OpenAI request: %v", err)
	}

	// Hacer la petición a OpenAI
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating OpenAI request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.OPENAI_API_KEY)
	req.Header.Set("Content-Type", "application/json")

	// Hacer la solicitud HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request to OpenAI: %v", err)
	}
	defer resp.Body.Close()

	// Parsear la respuesta
	var openAIResp OpenAIResponse
	err = json.NewDecoder(resp.Body).Decode(&openAIResp)
	if err != nil {
		return "", fmt.Errorf("error decoding OpenAI response: %v", err)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no valid response from OpenAI")
	}

	return openAIResp.Choices[0].Message.Content, nil
}
