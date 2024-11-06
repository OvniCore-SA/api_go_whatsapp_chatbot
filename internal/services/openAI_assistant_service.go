package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/openaiassistantdtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/openaiassistantdtos/openaivectorfiles"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/openaiassistantdtos/openaivectorsdtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services/clients"
)

// OpenIAAssistantService struct que encapsula la lógica de negocio para el asistente de OpenAI.
type OpenIAAssistantService struct {
	openAIClient *clients.OpenAIClient
}

// NewOpenIAAssistantService crea una nueva instancia de OpenIAAssistantService.
func NewOpenIAAssistantService(openAIClient *clients.OpenAIClient) *OpenIAAssistantService {
	return &OpenIAAssistantService{
		openAIClient: openAIClient,
	}
}

// CreateAssistant crea un nuevo asistente.
func (service *OpenIAAssistantService) CreateAssistant(request openaiassistantdtos.AssistantWithFileSearchRequest) (*openaiassistantdtos.ResponseCreateAssistant, error) {
	assistant, err := service.openAIClient.CreateAssistantWithFileSearch(context.Background(), request)
	if err != nil {
		log.Printf("Error creating assistant: %v", err)
		return nil, err
	}
	return assistant, nil
}

// CreateVectorSorage crea un nuevo vector storage, se nesesita para asociar con un archivo.
func (service *OpenIAAssistantService) CreateVectorSorage(request openaivectorsdtos.CreateVectorStoreRequest) (*openaivectorsdtos.VectorStore, error) {
	vectorStorage, err := service.openAIClient.CreateVectorStore(context.Background(), request)
	if err != nil {
		log.Printf("Error creating assistant: %v", err)
		return nil, err
	}
	return vectorStorage, nil
}

// GetResponseFromAssistant obtiene una respuesta del asistente.
func (service *OpenIAAssistantService) GetResponseFromAssistant(query openaiassistantdtos.AssistantQuery) (*openaiassistantdtos.AssistantResponse, error) {
	response, err := service.openAIClient.GetResponseFromAssistant(context.Background(), query)
	if err != nil {
		log.Printf("Error generating assistant response: %v", err)
		return nil, err
	}
	return response, nil
}

// GetResponseFromAssistant obtiene una respuesta del asistente.
func (service *OpenIAAssistantService) PostUploadFileForAssistant(request openaiassistantdtos.AssistantUploadFileQuery) (*openaiassistantdtos.FileUploadResponse, error) {
	response, err := service.openAIClient.UploadFile(context.Background(), request)
	if err != nil {
		log.Printf("Error generating assistant response: %v", err)
		return nil, err
	}

	return response, nil
}

// GetUploadFile recupera un file
func (service *OpenIAAssistantService) GetFile(fileID string) (*openaiassistantdtos.FileRetrieveResponse, error) {
	response, err := service.openAIClient.RetrieveFile(context.Background(), fileID)
	if err != nil {
		log.Printf("Error retrieve file response: %v", err)
		return nil, err
	}

	return response, nil
}

// SaveTextToJsonFile toma un texto, lo convierte a formato JSON Lines y lo guarda en un archivo .jsonl
func (service *OpenIAAssistantService) SaveTextToJsonFile(text string) (string, error) {
	// Definir la ruta del directorio y el nombre del archivo
	dir := "../pkg/files"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("error creating directory: %w", err)
	}

	// Obtener el timestamp actual
	timestamp := time.Now().Format("20060102150405") // Formato: YYYYMMDDHHMMSS

	// Generar el nombre de archivo usando el timestamp
	fileName := fmt.Sprintf("data_%s.json", timestamp)
	filePath := filepath.Join(dir, fileName)

	// Crear o abrir el archivo para escritura
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// Aquí simplemente creamos un objeto JSON con el texto
	jsonData := map[string]string{"text": text}

	// Convertir el objeto a JSON
	jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %w", err)
	}

	// Escribir el JSON en el archivo
	_, err = file.Write(jsonBytes)
	if err != nil {
		return "", fmt.Errorf("error writing to file: %w", err)
	}

	return filePath, nil
}

// PostFileToVectorStorage agrega un un archivo a un vector storage.
func (service *OpenIAAssistantService) AddFileToVectorStoreService(request openaivectorfiles.RequestVectorStorage) (*openaivectorfiles.VectorStoreFile, error) {

	requestFileName := openaivectorfiles.AddFileRequest{
		FileID: request.FileID,
	}
	response, err := service.openAIClient.AddFileToVectorStore(context.Background(), request.VectorStoreID, requestFileName)
	if err != nil {
		log.Printf("Error retrieve file response: %v", err)
		return nil, err
	}

	return response, nil
}

// GetVectorStoreFileService recupera vector storage file.
func (service *OpenIAAssistantService) GetVectorStoreFileService(vectorStoreID string, fileID string) (*openaivectorfiles.RequestVectorStoreFile, error) {

	response, err := service.openAIClient.GetVectorStoreFile(context.Background(), vectorStoreID, fileID)
	if err != nil {
		log.Printf("Error retrieve file response: %v", err)
		return nil, err
	}

	return response, nil
}
