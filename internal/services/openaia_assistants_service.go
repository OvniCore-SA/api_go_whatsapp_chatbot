package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
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

// Helper para realizar peticiones
func (s *OpenAIAssistantService) doRequest(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+s.apiKey)
	req.Header.Add("OpenAI-Beta", "assistants=v2")
	req.Header.Add("Content-Type", "application/json")

	fmt.Println("Ejecutando endpoint: " + req.URL.Path)

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
		"name":         name,
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
	req, err := http.NewRequest("POST", os.Getenv("OPENAI_API_URL")+"/assistants", bytes.NewBuffer(body))
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

// EditAssistant edita un asistente existente
func (s *OpenAIAssistantService) EditAssistant(assistantID, name, instructions, model, vectorStoreID string) (string, error) {
	data := map[string]interface{}{
		"instructions": instructions,
		"name":         name,
		"tools": []map[string]string{
			{"type": "file_search"},
		},
		"model": model,
	}

	body, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", os.Getenv("OPENAI_API_URL")+"/assistants/"+assistantID, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	// Enviar la solicitud
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
	url := fmt.Sprintf("%s/assistants/%s", os.Getenv("OPENAI_API_URL"), assistantID)
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
	req, err := http.NewRequest("POST", os.Getenv("OPENAI_API_URL")+"/vector_stores", bytes.NewBuffer(body))
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
	req, err := http.NewRequest("POST", os.Getenv("OPENAI_API_URL")+"/files", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
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
	url := fmt.Sprintf("%s/files/%s", os.Getenv("OPENAI_API_URL"), fileID)
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
	url := fmt.Sprintf("%s/vector_stores/%s/files/%s", os.Getenv("OPENAI_API_URL"), vectorStoreID, fileID)
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/vector_stores/%s/files", os.Getenv("OPENAI_API_URL"), vectorStoreID), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	_, err = s.doRequest(req)
	return err
}

// SendMessageToThread envía un mensaje a un Thread existente en OpenAI
func (s *OpenAIAssistantService) SendMessageToThread(threadID, message string, user bool) error {
	// Preparar el cuerpo de la solicitud con el mensaje especificado
	rol := "user"
	if !user {
		rol = "assistant"
	}
	data := map[string]interface{}{
		"role": rol,
		"content": []map[string]interface{}{
			{
				"text": message,
				"type": "text",
			},
		},
	}

	// Serializar el cuerpo de la solicitud a JSON
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %v", err)
	}

	// Crear la solicitud HTTP
	url := fmt.Sprintf("%s/threads/%s/messages", os.Getenv("OPENAI_API_URL"), threadID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Enviar la solicitud y procesar la respuesta

	resp, err := s.doRequest(req)
	if err != nil {
		return fmt.Errorf("error sending request to OpenAI: %v", err)
	}
	defer resp.Body.Close()

	// Verificar el código de respuesta
	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send messages to thread, status code: %d, response: %s", resp.StatusCode, string(responseBody))
	}

	fmt.Println("Mensaje enviado exitosamente")
	return nil
}

func (s *OpenAIAssistantService) CreateRunForThreadWithConversation(threadID, assistantID string, conversation []map[string]interface{}) (string, error) {
	// Estructura del cuerpo de la solicitud
	data := map[string]interface{}{
		"assistant_id": assistantID,
		// "additional_messages": conversation,
	}

	// Serializar el cuerpo de la solicitud
	body, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error marshalling request body: %v", err)
	}

	// Crear la solicitud HTTP
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/threads/%s/runs", os.Getenv("OPENAI_API_URL"), threadID), bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Enviar la solicitud y procesar la respuesta
	resp, err := s.doRequest(req)
	if err != nil {
		return "", fmt.Errorf("error sending request to OpenAI: %v", err)
	}
	defer resp.Body.Close()

	// Decodificar la respuesta
	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	// Retornar el ID del run creado
	return result.ID, nil
}

func (s *OpenAIAssistantService) WaitForRunCompletion(threadID, runID string, maxRetries int, retryInterval time.Duration) error {
	for i := 0; i < maxRetries; i++ {
		// Crear la solicitud HTTP para consultar el estado del run
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/threads/%s/runs/%s", os.Getenv("OPENAI_API_URL"), threadID, runID), nil)
		if err != nil {
			return fmt.Errorf("error creating request: %v", err)
		}

		// Enviar la solicitud y procesar la respuesta
		resp, err := s.doRequest(req)
		if err != nil {
			return fmt.Errorf("error sending request to OpenAI: %v", err)
		}
		defer resp.Body.Close()

		// Decodificar la respuesta
		var result struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return fmt.Errorf("error decoding response: %v", err)
		}

		// Verificar el estado del run
		switch result.Status {
		case "completed":
			// Si está completado, retornar sin error
			return nil
		case "incomplete":
			// Si está incomplete, retornar sin error
			return nil
		case "failed", "cancelled", "expired":
			// Si el run falla o es cancelado, retornar error
			return fmt.Errorf("run ended with status: %s", result.Status)
		default:
			// Si está en progreso o en cola, esperar y reintentar
			fmt.Printf("Run status: %s, retrying...\n", result.Status)
			time.Sleep(retryInterval)
		}
	}

	// Si se agotan los reintentos, retornar error
	return fmt.Errorf("run did not complete after %d retries", maxRetries)
}

func (s *OpenAIAssistantService) GetMessagesFromThread(threadID string) (string, error) {
	// Crear la solicitud HTTP
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/threads/%s/messages", os.Getenv("OPENAI_API_URL"), threadID), nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Enviar la solicitud y procesar la respuesta
	resp, err := s.doRequest(req)
	if err != nil {
		return "", fmt.Errorf("error sending request to OpenAI: %v", err)
	}
	defer resp.Body.Close()

	// Decodificar la respuesta
	var result struct {
		Data []struct {
			Role    string `json:"role"`
			Id      string `json:"id"`
			Content []struct {
				Type string `json:"type"`
				Text struct {
					Value string `json:"value"`
				} `json:"text"`
			} `json:"content"`
		} `json:"data"`
		FirstId string `json:"first_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	// Buscar el último mensaje del asistente
	for i := len(result.Data) - 1; i >= 0; i-- {
		if result.Data[i].Role == "assistant" && result.Data[i].Id == result.FirstId {
			for _, content := range result.Data[i].Content {
				if content.Type == "text" {
					return content.Text.Value, nil
				}
			}
		}
	}

	return "", fmt.Errorf("no assistant response found in thread messages")
}

// CreateThread crea un nuevo Thread en OpenAI con un asistente específico
func (s *OpenAIAssistantService) CreateThread(model, instructions string) (string, error) {
	// Estructura del cuerpo de la solicitud

	// Crear la solicitud HTTP
	req, err := http.NewRequest("POST", os.Getenv("OPENAI_API_URL")+"/threads", nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Enviar la solicitud y procesar la respuesta
	resp, err := s.doRequest(req)
	if err != nil {
		return "", fmt.Errorf("error sending request to OpenAI: %v", err)
	}
	defer resp.Body.Close()

	// Decodificar la respuesta
	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	return result.ID, nil
}

// EjecutarThread edita un thread colocandole un vectorstore.
func (s *OpenAIAssistantService) EjecutarThread(threadID string, vectorStoreIDs []string) error {
	// Construir el cuerpo de la solicitud
	payload := map[string]interface{}{
		"tool_resources": map[string]interface{}{
			"file_search": map[string]interface{}{
				"vector_store_ids": vectorStoreIDs,
			},
		},
	}

	// Codificar el cuerpo en JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error encoding payload: %v", err)
	}

	// Crear la solicitud HTTP
	url := fmt.Sprintf("%s/threads/%s", os.Getenv("OPENAI_API_URL"), threadID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Enviar la solicitud y procesar la respuesta
	resp, err := s.doRequest(req)
	if err != nil {
		return fmt.Errorf("error sending request to OpenAI: %v", err)
	}
	defer resp.Body.Close()

	// Verificar el código de respuesta
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Leer y procesar la respuesta (opcional)
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	fmt.Printf("Respuesta de OpenAI: %v\n", response)
	return nil
}
