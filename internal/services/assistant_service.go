package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/config"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
)

type AssistantService struct {
	repository             *mysql_client.AssistantRepository
	serviceFile            *FileService
	openAIAssistantService *OpenAIAssistantService
	client                 *http.Client
}

func NewAssistantService(repository *mysql_client.AssistantRepository, serviceFile *FileService, openAIAssistantService *OpenAIAssistantService) *AssistantService {
	return &AssistantService{
		repository:             repository,
		serviceFile:            serviceFile,
		openAIAssistantService: openAIAssistantService,
		client:                 &http.Client{},
	}
}

func (s *AssistantService) UploadFileToGPT(fileContent io.Reader, filename string) (string, error) {
	// Preparar el cuerpo de la solicitud como multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Crear la parte de archivo en la solicitud multipart
	part, err := writer.CreateFormFile("file", ".txt")
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(part, fileContent); err != nil {
		return "", err
	}

	// Añadir el propósito requerido por OpenAI
	writer.WriteField("purpose", "fine-tune")
	writer.Close()

	// Crear la solicitud HTTP
	req, err := http.NewRequest("POST", baseURL+"/files", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+config.OPENAI_API_KEY)
	req.Header.Set("Content-Type", writer.FormDataContentType())

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
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.ID, nil
}

func (m *AssistantService) CreateAssistantWithFile(data dtos.AssistantDto, fileHeader *multipart.FileHeader) (dtos.AssistantDto, error) {
	// Abrir el archivo
	fileContent, err := fileHeader.Open()
	if err != nil {
		return dtos.AssistantDto{}, fmt.Errorf("unable to open file: %w", err)
	}
	defer fileContent.Close()

	// Subir archivo a OpenAI
	fileIDOpenAI, err := m.openAIAssistantService.UploadFileToGPT(fileContent, fileHeader.Filename)
	if err != nil {
		return dtos.AssistantDto{}, err
	}

	// Restablecer el cursor después de subir a OpenAI
	if seeker, ok := fileContent.(io.Seeker); ok {
		_, err := seeker.Seek(0, io.SeekStart)
		if err != nil {
			return dtos.AssistantDto{}, fmt.Errorf("failed to reset file cursor after OpenAI upload: %w", err)
		}
	}

	// Crear vector store en OpenAI
	vectorStoreID, err := m.openAIAssistantService.CreateVectorStore(fileHeader.Filename)
	if err != nil {
		return dtos.AssistantDto{}, err
	}

	// Asignar archivo al vector store
	err = m.openAIAssistantService.addFileToVectorStore(vectorStoreID, fileIDOpenAI)
	if err != nil {
		return dtos.AssistantDto{}, err
	}

	// Crear el asistente en OpenAI
	assistantID, err := m.openAIAssistantService.CreateAssistant(data.Name, data.Instructions, config.OPENAI_MODEL_USE, vectorStoreID)
	if err != nil {
		return dtos.AssistantDto{}, err
	}

	data.OpenaiAssistantsID = assistantID
	assistantDB := entities.MapDtoToAssistant(data)

	// Guardar el asistente en la base de datos
	err = m.repository.Create(&assistantDB)
	if err != nil {
		return dtos.AssistantDto{}, err
	}

	// Subir a MinIO y registrar en DB
	_, err = m.serviceFile.CreateFile(fileHeader, assistantDB.ID, "assistants", fileIDOpenAI)
	if err != nil {
		return dtos.AssistantDto{}, err
	}

	return data, nil
}

func (s *AssistantService) CreateAssistant(data dtos.AssistantDto) (dtos.AssistantDto, error) {
	assistant := entities.MapDtoToAssistant(data)
	if err := s.repository.Create(&assistant); err != nil {
		return dtos.AssistantDto{}, err
	}
	return entities.MapAssistantToDto(assistant), nil
}

func (s *AssistantService) DeleteOpenAIAssistant(assistantID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/assistants/%s", baseURL, assistantID), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+config.OPENAI_API_KEY)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Verifica si la solicitud fue exitosa
	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request to delete assistant failed: %s", string(bodyBytes))
	}

	return nil
}

func (s *AssistantService) FindAllAssistants() ([]dtos.AssistantDto, error) {
	assistants, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}
	var assistantDtos []dtos.AssistantDto
	for _, assistant := range assistants {
		assistantDtos = append(assistantDtos, entities.MapAssistantToDto(assistant))
	}
	return assistantDtos, nil
}

func (s *AssistantService) FindAssistantById(id int64) (dtos.AssistantDto, error) {
	assistant, err := s.repository.FindById(id)
	if err != nil {
		return dtos.AssistantDto{}, err
	}
	return entities.MapAssistantToDto(assistant), nil
}

func (s *AssistantService) UpdateAssistant(id int64, data dtos.AssistantDto) (dtos.AssistantDto, error) {
	assistant := entities.MapDtoToAssistant(data)
	if err := s.repository.Update(id, assistant); err != nil {
		return dtos.AssistantDto{}, errors.New("assistant not found")
	}
	return entities.MapAssistantToDto(assistant), nil
}

func (s *AssistantService) DeleteAssistant(id int64) error {
	assistant, err := s.FindAssistantById(int64(id))
	if err != nil {
		return errors.New("assistant not found")
	}
	// Eliminar el asistente en OpenAI
	err = s.DeleteOpenAIAssistant(assistant.OpenaiAssistantsID)
	if err != nil {
		return fmt.Errorf("failed to delete assistant from OpenAI: %w", err)
	}
	return s.repository.Delete(id)
}

func (s *AssistantService) GetAllAssistantsByBussinessId(businessId int64) ([]dtos.AssistantDto, error) {
	assistants, err := s.repository.GetAllAssistantsByBussinessId(businessId)
	if err != nil {
		return nil, err
	}

	assistantDtos := make([]dtos.AssistantDto, len(assistants))
	for i, assistant := range assistants {
		assistantDtos[i] = entities.MapAssistantToDto(assistant)
	}
	return assistantDtos, nil
}
