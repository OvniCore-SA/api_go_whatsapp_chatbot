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

// AddFileToAssistant asocia un archivo existente a un asistente
func (s *OpenAIAssistantService) AddFileToAssistant(assistantID, fileID string) error {
	// Prepara el cuerpo de la solicitud con el file_id
	data := map[string]string{
		"file_id": fileID,
	}

	body, _ := json.Marshal(data)
	url := fmt.Sprintf("%s/assistants/%s/files", baseURL, assistantID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	// Añadir encabezados requeridos
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v1")
	req.Header.Set("Content-Type", "application/json")

	// Realiza la solicitud
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Verifica el estado de la respuesta
	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
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

// UploadFileToVectorStore sube archivos a un vector store existente
func (s *OpenAIAssistantService) UploadFileToVectorStore(vectorStoreID string, fileContent io.Reader) error {

	// Primero, subimos el archivo a OpenAI
	fileID, err := s.UploadFileToGPT(fileContent, "")
	if err != nil {
		return err
	}

	// Luego, asociamos el archivo subido con el vector store
	err = s.addFileToVectorStore(vectorStoreID, fileID)
	if err != nil {
		return err
	}

	return nil
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

// AssociateAssistantWithVectorStore asocia un vector store con un asistente
func (s *OpenAIAssistantService) AssociateAssistantWithVectorStore(assistantID, vectorStoreID string) error {
	data := map[string]interface{}{
		"tool_resources": map[string]interface{}{
			"file_search": map[string][]string{"vector_store_ids": {vectorStoreID}},
		},
	}

	body, _ := json.Marshal(data)
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/assistants/%s", baseURL, assistantID), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	_, err = s.doRequest(req)
	return err
}
