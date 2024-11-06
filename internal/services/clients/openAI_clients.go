package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/openaiassistantdtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/openaiassistantdtos/openaimessages"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/openaiassistantdtos/openairuns"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/openaiassistantdtos/openaithreads"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/openaiassistantdtos/openaivectorfiles"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/openaiassistantdtos/openaivectorsdtos"
)

type OpenAIClient struct {
	apiBaseURL string
	apiKey     string
}

func NewOpenAIClient(apiBaseURL, apiKey string) *OpenAIClient {
	return &OpenAIClient{
		apiBaseURL: apiBaseURL,
		apiKey:     apiKey,
	}
}

func (client *OpenAIClient) GetResponseFromAssistant(ctx context.Context, query openaiassistantdtos.AssistantQuery) (*openaiassistantdtos.AssistantResponse, error) {
	url := client.apiBaseURL + "/assistants/" + query.AssistantID + "/responses"

	reqBody, _ := json.Marshal(query)
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error generating assistant response")
	}

	var response openaiassistantdtos.AssistantResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (client *OpenAIClient) CreateAssistantWithFileSearch(ctx context.Context, request openaiassistantdtos.AssistantWithFileSearchRequest) (ResponseCreateAssistant *openaiassistantdtos.ResponseCreateAssistant, err error) {
	url := client.apiBaseURL + "/assistants"

	// Convertir la solicitud a JSON
	reqBody, _ := json.Marshal(request)
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, err
	}

	// Establecer cabeceras
	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	// Realizar la solicitud HTTP
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Manejo de errores HTTP
	if resp.StatusCode != http.StatusOK {
		var apiErr struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil {
			return nil, errors.New("error creating assistant: " + apiErr.Error)
		}
		return nil, errors.New("error creating assistant: status code " + resp.Status)
	}

	// Decodificar la respuesta
	if err := json.NewDecoder(resp.Body).Decode(&ResponseCreateAssistant); err != nil {
		return nil, err
	}

	return ResponseCreateAssistant, nil
}

// UploadFile guarda un archivo en OpenAI
func (client *OpenAIClient) UploadFile(ctx context.Context, request openaiassistantdtos.AssistantUploadFileQuery) (*openaiassistantdtos.FileUploadResponse, error) {
	url := client.apiBaseURL + "/files"

	// Abrir el archivo
	file, err := os.Open(request.PathFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Crear una tubería para conectar la escritura de multipart con la lectura de la solicitud HTTP
	pr, pw := io.Pipe()
	body := multipart.NewWriter(pw)

	// Crear la solicitud HTTP
	req, err := http.NewRequestWithContext(ctx, "POST", url, pr)
	if err != nil {
		return nil, err
	}

	// Establecer las cabeceras
	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("Content-Type", body.FormDataContentType())

	// Escribir en el cuerpo de la solicitud en una goroutine
	go func() {
		defer pw.Close()
		defer body.Close()

		// Adjuntar el archivo a la solicitud
		part, err := body.CreateFormFile("file", filepath.Base(file.Name()))
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		_, err = io.Copy(part, file)
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		// Adjuntar el propósito a la solicitud
		if err := body.WriteField("purpose", request.Purpose); err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	// Hacer la petición
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Manejar errores HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error uploading file")
	}

	// Decodificar la respuesta
	var fileResp openaiassistantdtos.FileUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&fileResp); err != nil {
		return nil, err
	}

	return &fileResp, nil
}

// DeleteFile elimina un archivo en OpenAI.
func (client *OpenAIClient) DeleteFile(ctx context.Context, fileID string) (*openaiassistantdtos.FileDeleteResponse, error) {
	url := client.apiBaseURL + "/files/" + fileID

	// Crear la solicitud HTTP
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	// Establecer las cabeceras
	req.Header.Set("Authorization", "Bearer "+client.apiKey)

	// Hacer la petición
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Manejar errores HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error deleting file")
	}

	// Decodificar la respuesta
	var fileResp openaiassistantdtos.FileDeleteResponse
	if err := json.NewDecoder(resp.Body).Decode(&fileResp); err != nil {
		return nil, err
	}

	return &fileResp, nil
}

// Recupera un archivo de OpenAI en base a su ID
func (client *OpenAIClient) RetrieveFile(ctx context.Context, fileID string) (*openaiassistantdtos.FileRetrieveResponse, error) {
	url := client.apiBaseURL + "/files/" + fileID

	// Crear la solicitud HTTP
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Establecer las cabeceras
	req.Header.Set("Authorization", "Bearer "+client.apiKey)

	// Hacer la petición
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Manejar errores HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error retrieving file")
	}

	// Decodificar la respuesta
	var fileResp openaiassistantdtos.FileRetrieveResponse
	if err := json.NewDecoder(resp.Body).Decode(&fileResp); err != nil {
		return nil, err
	}

	return &fileResp, nil
}

// Crea un nuevo vector.
// Este se utiliza luego para cargar un archivo a este.
func (client *OpenAIClient) CreateVectorStore(ctx context.Context, request openaivectorsdtos.CreateVectorStoreRequest) (*openaivectorsdtos.VectorStore, error) {
	url := client.apiBaseURL + "/vector_stores"

	// Convertir la solicitud a JSON
	reqBody, _ := json.Marshal(request)
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, err
	}

	// Establecer cabeceras
	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	// Realizar la solicitud HTTP
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Manejo de errores HTTP
	if resp.StatusCode != http.StatusOK {
		var apiErr struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil {
			return nil, errors.New("error creating vector store: " + apiErr.Error)
		}
		return nil, errors.New("error creating vector store: status code " + resp.Status)
	}

	// Decodificar la respuesta
	var vectorStore openaivectorsdtos.VectorStore
	if err := json.NewDecoder(resp.Body).Decode(&vectorStore); err != nil {
		return nil, err
	}

	return &vectorStore, nil
}

// VECTOR STORAGE FILES
//
// Crea un vector storage en base a un archivo ya subido. Esto sirve para que el asistente trabaje con este vector storage
func (client *OpenAIClient) AddFileToVectorStore(ctx context.Context, vectorStoreID string, request openaivectorfiles.AddFileRequest) (*openaivectorfiles.VectorStoreFile, error) {
	url := client.apiBaseURL + "/vector_stores/" + vectorStoreID + "/files"

	// Convertir la solicitud a JSON
	reqBody, _ := json.Marshal(request)
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, err
	}

	// Establecer cabeceras
	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	// Realizar la solicitud HTTP
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Manejo de errores HTTP
	if resp.StatusCode != http.StatusOK {
		var apiErr struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil {
			return nil, errors.New("error adding file to vector store: " + apiErr.Error)
		}
		return nil, errors.New("error adding file to vector store: status code " + resp.Status)
	}

	// Decodificar la respuesta
	var vectorStoreFile openaivectorfiles.VectorStoreFile
	if err := json.NewDecoder(resp.Body).Decode(&vectorStoreFile); err != nil {
		return nil, err
	}

	return &vectorStoreFile, nil
}

// Obtiene un vector storage file
func (service *OpenAIClient) GetVectorStoreFile(ctx context.Context, vectorStoreID, fileID string) (*openaivectorfiles.RequestVectorStoreFile, error) {
	// Construir la URL del endpoint
	url := fmt.Sprintf("%s/vector_stores/%s/files/%s", service.apiBaseURL, vectorStoreID, fileID)

	// Crear la solicitud HTTP GET
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Agregar los headers necesarios
	req.Header.Set("Authorization", "Bearer "+service.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	// Ejecutar la solicitud
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado de la respuesta
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error retrieving file, status code: %d", resp.StatusCode)
	}

	// Decodificar la respuesta JSON en un objeto
	var file openaivectorfiles.RequestVectorStoreFile
	if err := json.NewDecoder(resp.Body).Decode(&file); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &file, nil
}

// HILOS
//
// CreateThread crea un nuevo hilo (thread) con o sin mensajes iniciales.
func (client *OpenAIClient) CreateThread(ctx context.Context, request openaithreads.ThreadRequest) (*openaithreads.ThreadResponse, error) {
	// Construir la URL del endpoint
	url := fmt.Sprintf("%s/threads", client.apiBaseURL)

	// Marshal de la solicitud a JSON
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	// Crear la solicitud HTTP POST
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Agregar los headers necesarios
	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	// Ejecutar la solicitud
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making POST request: %w", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado de la respuesta
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error creating thread, status code: " + resp.Status)
	}

	// Decodificar la respuesta JSON en un objeto
	var threadResponse openaithreads.ThreadResponse
	if err := json.NewDecoder(resp.Body).Decode(&threadResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &threadResponse, nil
}

// GetThreadDetails obtiene los detalles de un hilo específico.
func (client *OpenAIClient) GetThreadDetails(ctx context.Context, threadID string) (*openaithreads.ThreadDetailsResponse, error) {
	// Construir la URL del endpoint
	url := fmt.Sprintf("%s/threads/%s", client.apiBaseURL, threadID)

	// Crear la solicitud HTTP GET
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Agregar los headers necesarios
	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	// Ejecutar la solicitud
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado de la respuesta
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error getting thread details, status code: " + resp.Status)
	}

	// Decodificar la respuesta JSON en un objeto
	var threadDetails openaithreads.ThreadDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&threadDetails); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &threadDetails, nil
}

// UpdateThread actualiza un hilo específico con los datos proporcionados.
func (client *OpenAIClient) UpdateThread(ctx context.Context, threadID string, updateRequest openaithreads.ThreadUpdateRequest) (*openaithreads.ThreadUpdateResponse, error) {
	// Construir la URL del endpoint
	url := fmt.Sprintf("%s/threads/%s", client.apiBaseURL, threadID)

	// Marshal de la solicitud a JSON
	reqBody, err := json.Marshal(updateRequest)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	// Crear la solicitud HTTP PATCH
	req, err := http.NewRequestWithContext(ctx, "PATCH", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Agregar los headers necesarios
	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	// Ejecutar la solicitud
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making PATCH request: %w", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado de la respuesta
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error updating thread, status code: " + resp.Status)
	}

	// Decodificar la respuesta JSON en un objeto
	var threadUpdateResponse openaithreads.ThreadUpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&threadUpdateResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &threadUpdateResponse, nil
}

// DeleteThread elimina un hilo específico.
func (client *OpenAIClient) DeleteThread(ctx context.Context, threadID string) (*openaithreads.ThreadDeleteResponse, error) {
	// Construir la URL del endpoint
	url := fmt.Sprintf("%s/threads/%s", client.apiBaseURL, threadID)

	// Crear la solicitud HTTP DELETE
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Agregar los headers necesarios
	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	// Ejecutar la solicitud
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making DELETE request: %w", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado de la respuesta
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error deleting thread, status code: " + resp.Status)
	}

	// Decodificar la respuesta JSON en un objeto
	var threadDeleteResponse openaithreads.ThreadDeleteResponse
	if err := json.NewDecoder(resp.Body).Decode(&threadDeleteResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &threadDeleteResponse, nil
}

// MENSAJES
//
// SendMessageToThread crea un mensaje cuando el usuario manda
func (client *OpenAIClient) SendMessageToThread(threadID string, request openaimessages.SendMessageRequest) (*openaimessages.SendMessageResponse, error) {
	url := fmt.Sprintf("%s/threads/%s/messages", client.apiBaseURL, threadID)

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error encoding request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response status: %s", resp.Status)
	}

	var response openaimessages.SendMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response body: %w", err)
	}

	return &response, nil
}

func (client *OpenAIClient) GetMessageFromThread(threadID, messageID string) (*openaimessages.GetMessageResponse, error) {
	url := fmt.Sprintf("%s/threads/%s/messages/%s", client.apiBaseURL, threadID, messageID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response status: %s", resp.Status)
	}

	var response openaimessages.GetMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response body: %w", err)
	}

	return &response, nil
}

func (client *OpenAIClient) DeleteMessageFromThread(threadID, messageID string) (*openaimessages.DeleteMessageResponse, error) {
	url := fmt.Sprintf("%s/threads/%s/messages/%s", client.apiBaseURL, threadID, messageID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response status: %s", resp.Status)
	}

	var response openaimessages.DeleteMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response body: %w", err)
	}

	return &response, nil
}

func (client *OpenAIClient) UpdateMessageMetadata(threadID, messageID string, request openaimessages.UpdateMessageRequest) (*openaimessages.UpdateMessageResponse, error) {
	url := fmt.Sprintf("%s/threads/%s/messages/%s", client.apiBaseURL, threadID, messageID)

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error encoding request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response status: %s", resp.Status)
	}

	var response openaimessages.UpdateMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response body: %w", err)
	}

	return &response, nil
}

// CORRIDAS
//
// StartRun corre un hilo
func (client *OpenAIClient) StartRun(threadID string, request openairuns.StartRunRequest) (*openairuns.StartRunResponse, error) {
	url := fmt.Sprintf("%s/threads/%s/runs", client.apiBaseURL, threadID)

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error encoding request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response status: %s", resp.Status)
	}

	var response openairuns.StartRunResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response body: %w", err)
	}

	return &response, nil
}

func (client *OpenAIClient) GetRun(threadID, runID string) (*openairuns.GetRunResponse, error) {
	url := fmt.Sprintf("%s/threads/%s/runs/%s", client.apiBaseURL, threadID, runID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response status: %s", resp.Status)
	}

	var response openairuns.GetRunResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response body: %w", err)
	}

	return &response, nil
}

func (client *OpenAIClient) CancelRun(threadID, runID string) (*openairuns.CancelRunResponse, error) {
	url := fmt.Sprintf("%s/threads/%s/runs/%s/cancel", client.apiBaseURL, threadID, runID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response status: %s", resp.Status)
	}

	var response openairuns.CancelRunResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response body: %w", err)
	}

	return &response, nil
}

func (client *OpenAIClient) UpdateRunMetadata(threadID, runID string, request openairuns.UpdateRunMetadataRequest) (*openairuns.UpdateRunMetadataResponse, error) {
	url := fmt.Sprintf("%s/threads/%s/runs/%s", client.apiBaseURL, threadID, runID)

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error encoding request body: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response status: %s", resp.Status)
	}

	var response openairuns.UpdateRunMetadataResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response body: %w", err)
	}

	return &response, nil
}
