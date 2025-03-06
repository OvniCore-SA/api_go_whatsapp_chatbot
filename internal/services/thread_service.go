package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities/mappers"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
	"gorm.io/gorm"
)

// ThreadService define la lógica de negocio para hilos
type ThreadService struct {
	threadRepo             *mysql_client.ThreadRepository
	openAIAssistantService *OpenAIAssistantService
}

// NewThreadService crea una nueva instancia del servicio
func NewThreadService(threadRepo *mysql_client.ThreadRepository, openAIAssistantService *OpenAIAssistantService) *ThreadService {
	return &ThreadService{threadRepo: threadRepo, openAIAssistantService: openAIAssistantService}
}

// Obtener el último Thread de un contacto o crear uno nuevo si tiene más de 12 horas
func (s *ThreadService) GetOrCreateThread(contact dtos.ContactDto, assistant dtos.AssistantDto) (*dtos.ThreadResponse, error) {
	// Buscar el último Thread del contacto que no esté eliminado
	lastThread, err := s.threadRepo.FindLastActiveByContactID(contact.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error buscando thread del contacto: %v", err)
	}

	// Si hay un thread y tiene menos de 12 horas, devolverlo
	if lastThread != nil && time.Since(lastThread.CreatedAt) < 12*time.Hour {
		thread := mappers.ToThreadResponse(*lastThread)
		return &thread, nil
	}

	// Si no hay un thread reciente, crear uno nuevo en OpenAI
	newThreadID, err := s.openAIAssistantService.CreateThread(assistant.Model, assistant.Instructions)
	if err != nil {
		return nil, fmt.Errorf("error creando thread en OpenAI: %v", err)
	}

	// Crear el nuevo Thread en la base de datos
	newThread := entities.Thread{
		ThreadsId:  newThreadID,
		ContactsID: contact.ID,
		Active:     true,
		CreatedAt:  time.Now(),
	}

	createdThread, err := s.threadRepo.Create(newThread)
	if err != nil {
		return nil, fmt.Errorf("error guardando thread en base de datos: %v", err)
	}

	if lastThread != nil && lastThread.ID > 0 {
		if err := s.threadRepo.Delete(lastThread.ID); err != nil {
			return nil, fmt.Errorf("error guardando thread en base de datos: %v", err)
		}
	}

	thread := mappers.ToThreadResponse(*createdThread)
	return &thread, nil
}

// Crear un nuevo hilo
func (s *ThreadService) CreateThread(dto dtos.ThreadCreateRequest) (*dtos.ThreadResponse, error) {
	// Convertir DTO a entidad
	threadEntity := mappers.ToThreadEntity(dto)

	// Guardar en la base de datos
	savedThread, err := s.threadRepo.Create(threadEntity)
	if err != nil {
		return nil, err
	}

	// Convertir a DTO de respuesta
	response := mappers.ToThreadResponse(*savedThread)
	return &response, nil
}

// Obtener un hilo por ID
func (s *ThreadService) GetThreadByID(id int64) (*dtos.ThreadResponse, error) {
	thread, err := s.threadRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("hilo no encontrado")
	}
	response := mappers.ToThreadResponse(*thread)
	return &response, nil
}

// Obtener un hilo por ThreadsId
func (s *ThreadService) GetThreadByThreadsId(threadsId string) (*dtos.ThreadResponse, error) {
	thread, err := s.threadRepo.FindByThreadsId(threadsId)
	if err != nil {
		return nil, errors.New("hilo no encontrado")
	}
	response := mappers.ToThreadResponse(*thread)
	return &response, nil
}

// Obtener todos los hilos
func (s *ThreadService) GetAllThreads() ([]dtos.ThreadResponse, error) {
	threads, err := s.threadRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return mappers.ToThreadResponseList(threads), nil
}

// Actualizar un hilo
func (s *ThreadService) UpdateThread(id int64, dto dtos.ThreadCreateRequest) error {
	threadEntity := mappers.ToThreadEntity(dto)
	return s.threadRepo.Update(id, threadEntity)
}

// Eliminar un hilo
func (s *ThreadService) DeleteThread(id int64) error {
	return s.threadRepo.Delete(id)
}
