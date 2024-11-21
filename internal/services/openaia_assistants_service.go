package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/config"
)

type OpenAIAssistantService struct {
	apiKey string
	client *http.Client
}

func NewOpenAIAssistantService(apiKey string) *OpenAIAssistantService {
	return &OpenAIAssistantService{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

const baseURL = "https://api.openai.com/v1"

// Helper para realizar peticiones
func (s *OpenAIAssistantService) doRequest(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+s.apiKey)
	req.Header.Add("OpenAI-Beta", "assistants=v2")
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return resp, nil
}

// CreateAssistant crea un nuevo asistente con búsqueda de archivos activada
func (s *OpenAIAssistantService) CreateAssistant(name, instructions, model, vectorStoreID string) (string, error) {
	data := map[string]interface{}{
		"instructions": instructions,
		"tools": []map[string]string{
			{"type": "file_search"},
		},
		"tool_resources": map[string]interface{}{
			"file_search": map[string]interface{}{
				"vector_store_ids": []string{vectorStoreID},
			},
		},
		"model": model,
	}

	body, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", baseURL+"/assistants", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	resp, err := s.doRequest(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.ID, nil
}

// DeleteAssistant elimina un asistente específico de OpenAI por su ID
func (s *OpenAIAssistantService) DeleteAssistant(assistantID string) error {
	// Crear la solicitud DELETE con la URL del asistente
	url := fmt.Sprintf("%s/assistants/%s", baseURL, assistantID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	// Realizar la solicitud utilizando el helper doRequest
	resp, err := s.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Leer la respuesta y verificar errores en el cuerpo si es necesario
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete assistant, status: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}

// CreateVectorStore crea un nuevo vector store
func (s *OpenAIAssistantService) CreateVectorStore(name string) (string, error) {
	data := map[string]string{"name": name}

	body, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", baseURL+"/vector_stores", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	resp, err := s.doRequest(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.ID, nil
}

func (s *OpenAIAssistantService) UploadFileToGPT(fileContent io.Reader, filename string) (string, error) {
	// Preparar el cuerpo de la solicitud como multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Crear la parte de archivo en la solicitud multipart
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}

	// Copiar el contenido desde io.Reader a la parte del archivo en el formulario
	if _, err = io.Copy(part, fileContent); err != nil {
		return "", err
	}

	// Añadir el propósito requerido por OpenAI
	writer.WriteField("purpose", "assistants")
	writer.Close()

	// Crear la solicitud HTTP
	req, err := http.NewRequest("POST", baseURL+"/files", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+config.OPENAI_API_KEY)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	//req.Header.Set("OpenAI-Beta", "assistants=v1")

	// Realizar la solicitud
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Verificar si la solicitud fue exitosa
	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed: %s", string(bodyBytes))
	}

	// Decodificar la respuesta para obtener el file_id
	var result struct {
		ID       string `json:"id"`
		Filename string `json:"filename"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.ID, nil
}

// DeleteFile elimina un archivo específico de OpenAI por su ID
func (s *OpenAIAssistantService) DeleteFile(fileID string) error {
	// Crear la URL para la solicitud DELETE
	url := fmt.Sprintf("%s/files/%s", baseURL, fileID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	// Realizar la solicitud utilizando el helper doRequest
	resp, err := s.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Leer la respuesta y verificar errores en el cuerpo si es necesario
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete file, status: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}

// DeleteFileFromVectorStore elimina un archivo específico de un vector_store por su ID y el ID del archivo
func (s *OpenAIAssistantService) DeleteFileFromVectorStore(vectorStoreID, fileID string) error {
	// Crear la URL para la solicitud DELETE
	url := fmt.Sprintf("%s/vector_stores/%s/files/%s", baseURL, vectorStoreID, fileID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	// Realizar la solicitud utilizando el helper doRequest
	resp, err := s.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Leer la respuesta y verificar errores en el cuerpo si es necesario
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete file from vector_store, status: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}

// addFileToVectorStore asocia un archivo con un vector store
func (s *OpenAIAssistantService) addFileToVectorStore(vectorStoreID, fileID string) error {
	data := map[string]string{
		"file_id": fileID,
	}

	body, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/vector_stores/%s/files", baseURL, vectorStoreID), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	_, err = s.doRequest(req)
	return err
}
