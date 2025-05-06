package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/postgres_client"
)

type AssistantService struct {
	repository             *postgres_client.AssistantRepository
	serviceFile            *FileService
	openAIAssistantService *OpenAIAssistantService
	client                 *http.Client
}

func NewAssistantService(repository *postgres_client.AssistantRepository, serviceFile *FileService, openAIAssistantService *OpenAIAssistantService) *AssistantService {
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
	req, err := http.NewRequest("POST", os.Getenv("OPENAI_API_URL")+"/files", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
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
	assistantID, err := m.openAIAssistantService.CreateAssistant(data.Name, data.Instructions, data.Model, vectorStoreID)
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
	_, err = m.serviceFile.CreateFile(fileHeader, assistantDB.ID, "assistants", fileIDOpenAI, vectorStoreID)
	if err != nil {
		return dtos.AssistantDto{}, err
	}

	data.ID = assistantDB.ID

	return data, nil
}

func (s *AssistantService) CreateAssistant(data dtos.AssistantDto) (dtos.AssistantDto, error) {
	// Crear el asistente en OpenAI
	assistantID, err := s.openAIAssistantService.CreateAssistant(data.Name, data.Instructions, data.Model, "")
	if err != nil {
		return dtos.AssistantDto{}, err
	}

	data.OpenaiAssistantsID = assistantID

	assistant := entities.MapDtoToAssistant(data)
	if err := s.repository.Create(&assistant); err != nil {
		return dtos.AssistantDto{}, err
	}
	return entities.MapAssistantToDto(assistant), nil
}

func (s *AssistantService) FindNumberPhonesByAssistantID(assistantID uint64) (dtos.NumberPhoneDto, error) {
	numberPhone, err := s.repository.FindByAssistantID(int64(assistantID))
	if err != nil {
		return dtos.NumberPhoneDto{}, err
	}
	if len(numberPhone) <= 0 {
		return dtos.NumberPhoneDto{}, fmt.Errorf("No hay numeros de telefonos para este assistant.")
	}

	return entities.MapEntityToNumberPhoneDto(numberPhone[0]), nil
}

func (s *AssistantService) DeleteOpenAIAssistant(assistantID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/assistants/%s", os.Getenv("OPENAI_API_URL"), assistantID), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
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

// IsWithinWorkingHours verifica si la fecha y hora están dentro del horario de trabajo del asistente
func (s *AssistantService) IsWithinWorkingHours(assistantID int64, dateTime time.Time) (bool, error) {
	// Llamar al repositorio para obtener los datos del asistente
	isAvailable, err := s.repository.IsWithinWorkingHours(assistantID, dateTime)
	if err != nil {
		return false, fmt.Errorf("error al verificar las horas de trabajo: %v", err)
	}

	return isAvailable, nil
}

// PartialUpdateAssistant actualiza solo los campos proporcionados en la solicitud
func (s *AssistantService) PartialUpdateAssistant(id int64, data dtos.AssistantDto) (dtos.AssistantDto, error) {
	// Obtener el registro actual
	existingAssistant, err := s.repository.FindById(id)
	if err != nil {
		return dtos.AssistantDto{}, errors.New("assistant not found")
	}

	EditNameOpenAI := false
	EditInstructionsOpenAI := false
	EditModelOpenAI := false

	// Solo actualizamos los valores enviados en el body
	if data.Name != "" {
		existingAssistant.Name = data.Name
		EditNameOpenAI = true
	}
	if data.Description != "" {
		existingAssistant.Description = data.Description
	}
	if data.Model != "" {
		existingAssistant.Model = data.Model
		EditModelOpenAI = true
	}
	if data.Instructions != "" {
		existingAssistant.Instructions = data.Instructions
		EditInstructionsOpenAI = true
	}
	if data.OpenaiAssistantsID != "" {
		existingAssistant.OpenaiAssistantsID = data.OpenaiAssistantsID
	}
	existingAssistant.Active = data.Active
	if data.EventDuration > 0 {
		existingAssistant.EventDuration = data.EventDuration
	}
	existingAssistant.AccountGoogle = data.AccountGoogle
	if data.EventType != "" {
		existingAssistant.EventType = data.EventType
	}
	if data.EventCountPerDay > 0 {
		existingAssistant.EventCountPerDay = data.EventCountPerDay
	}
	if data.OpeningDays > 0 {
		existingAssistant.OpeningDays = data.OpeningDays
	}
	if data.WorkingHours != "" {
		existingAssistant.WorkingHours = data.WorkingHours
	}

	if EditInstructionsOpenAI || EditModelOpenAI || EditNameOpenAI {
		// Se actualiza solo el campo que vino y si no se le coloca el que ya tenía porque se envia a actualizar a openAI
		if !EditInstructionsOpenAI {
			data.Instructions = existingAssistant.Instructions
		}
		if !EditModelOpenAI {
			data.Model = existingAssistant.Model
		}
		if !EditNameOpenAI {
			data.Name = existingAssistant.Name
		}
		data.OpenaiAssistantsID = existingAssistant.OpenaiAssistantsID
		// Actualizo los datos del assistant en OPEN AI
		if _, err := s.openAIAssistantService.EditAssistant(data.OpenaiAssistantsID, data.Name, data.Instructions, data.Model); err != nil {
			return dtos.AssistantDto{}, err
		}
	}

	// Guardamos los cambios en la base de datos
	if err := s.repository.Update(id, existingAssistant); err != nil {
		return dtos.AssistantDto{}, errors.New("failed to update assistant")
	}

	return entities.MapAssistantToDto(existingAssistant), nil
}

func (s *AssistantService) UpdateAssistant(id int64, data dtos.AssistantDto) (dtos.AssistantDto, error) {

	// Actualizo los datos del assistant en OPEN AI
	if _, err := s.openAIAssistantService.EditAssistant(data.OpenaiAssistantsID, data.Name, data.Instructions, data.Model); err != nil {
		return dtos.AssistantDto{}, err
	}

	assistant := entities.MapDtoToAssistant(data)
	if err := s.repository.Update(id, assistant); err != nil {
		return dtos.AssistantDto{}, errors.New("assistant not found")
	}
	return entities.MapAssistantToDto(assistant), nil
}

func (s *AssistantService) UpdateAssistantWithFile(id int64, data dtos.AssistantDto, fileHeader *multipart.FileHeader) (dtos.AssistantDto, error) {
	// Busco los files asociados a este assistente, por ahora solo debe tener uno. La relacion es para tener un respaldo de los otros solamente y porque un assistente puede tener muchos archivos en GPT pero no lo usamos así.

	files, err := s.serviceFile.GetFileByAssistantID(id)
	if err != nil {
		return dtos.AssistantDto{}, errors.New("failed to find files with assistantID")
	}

	if len(files) > 1 || len(files) == 0 {
		return dtos.AssistantDto{}, errors.New("failed to find files with assistantID. >1 <0")
	}

	// Desvinculo el archivo con el vector store OpenAI. Con esto se logra que el archivo quede vivo por si se quiere usar en otra ocación
	err = s.openAIAssistantService.DeleteFileFromVectorStore(files[0].OpenaiVectorStoreIDs, files[0].OpenaiFilesID)
	if err != nil {
		return dtos.AssistantDto{}, err
	}

	// Abrir el archivo
	fileContent, err := fileHeader.Open()
	if err != nil {
		return dtos.AssistantDto{}, fmt.Errorf("unable to open file: %w", err)
	}
	defer fileContent.Close()

	// Subir archivo a OpenAI
	fileIDOpenAI, err := s.openAIAssistantService.UploadFileToGPT(fileContent, fileHeader.Filename)
	if err != nil {
		return dtos.AssistantDto{}, err
	}

	// Asignar archivo al vector store
	err = s.openAIAssistantService.addFileToVectorStore(files[0].OpenaiVectorStoreIDs, fileIDOpenAI)
	if err != nil {
		// Si falla la desvinculación, elimino el archivo de OpenAI directamente y con estó logramos que se desvincule también el archvo.
		err = s.openAIAssistantService.DeleteFile(files[0].OpenaiFilesID)
		if err != nil {
			return dtos.AssistantDto{}, err
		}
		fmt.Printf("FILE '%s' eliminado.", fileIDOpenAI)
		return dtos.AssistantDto{}, errors.New("failed to assign file to vector store")
	}

	// Subir a MinIO y registrar en DB
	_, err = s.serviceFile.CreateFile(fileHeader, id, "assistants", fileIDOpenAI, files[0].OpenaiVectorStoreIDs)
	if err != nil {
		return dtos.AssistantDto{}, err
	}

	// Luego de hacer todas las operaciones anteriores con respecto al file y de haber creado y asociado el nuevo, procedo a eliminar de la DB
	err = s.serviceFile.DeleteFile(files[0].ID)
	if err != nil {
		return dtos.AssistantDto{}, err
	}

	// Actualizo los datos del assistant en OPEN AI
	_, err = s.openAIAssistantService.EditAssistant(data.OpenaiAssistantsID, data.Name, data.Instructions, data.Model)
	if err != nil {
		return dtos.AssistantDto{}, err
	}

	// Actualizo los otros campos del assistente
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
		fmt.Println(err.Error())
		if !strings.Contains(err.Error(), "No assistant found with id") {
			return fmt.Errorf("failed to delete assistant from OpenAI: %w", err)
		}
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
