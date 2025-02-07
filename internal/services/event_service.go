package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
	"golang.org/x/exp/rand"
)

// Interfaz de servicio con DTOs
type EventsService interface {
	Create(eventDTO dtos.EventsDto) error
	GetByID(id int) (*dtos.EventsDto, error)
	GetAll() ([]dtos.EventsDto, error)
	Update(eventDTO dtos.EventsDto) error
	Delete(id int) error
	Cancel(codeEvent string) error
	// método para buscar un evento por contacto, fecha y hora
	GetEventByContactAndDate(contactID int64, date, currentTime string) ([]entities.Events, error)
	// Verifica si un codigo existe en un evento.
	IsCodeUnique(code string) (bool, error)
	// Genera un codigo unico para un nuevo evento.
	GenerateUniqueCode() (string, error)
	GetEventByCodeEvent(contactID int64, codeEvent string) (dtos.EventsDto, error)
	GetEventsByContactDateAndNumberPhone(contactID int64, date string, assistantID int64) ([]entities.Events, error)
}

// Implementación del servicio
type eventsServiceImpl struct {
	repo        mysql_client.EventsRepository
	utilService UtilService
}

func NewEventsService(repo mysql_client.EventsRepository, utilService UtilService) EventsService {
	return &eventsServiceImpl{
		repo:        repo,
		utilService: utilService,
	}
}

// GenerateUniqueCode generates a unique code for event, ensuring it does not exist in the database
func (s *eventsServiceImpl) GenerateUniqueCode() (string, error) {
	rand.Seed(uint64(time.Now().UnixNano())) // Seed the random number generator
	attempts := 0
	maxAttempts := 10 // Limita el número de intentos para evitar un bucle infinito

	for {
		if attempts >= maxAttempts {
			return "", fmt.Errorf("failed to generate a unique code after %d attempts", maxAttempts)
		}
		code := s.utilService.GenerateUniqueCode() // Asume que UtilService tiene un método GenerateUniqueCode
		unique, err := s.repo.ExistsByCode(code)
		if err != nil {
			return "", err
		}
		if !unique {
			return code, nil
		}
		attempts++
	}
}
func (s *eventsServiceImpl) GetEventByCodeEvent(contactID int64, codeEvent string) (dtos.EventsDto, error) {
	// Llamamos al repositorio pasando el contactID y el codeEvent
	event, err := s.repo.FindByContactAndCodeEvent(contactID, codeEvent)
	if err != nil {
		return dtos.EventsDto{}, fmt.Errorf("error fetching event by code_event: %v", err)
	}
	return entities.MapEntityToEventsDto(event), nil
}

func (s *eventsServiceImpl) GetEventByContactAndDate(contactID int64, date, currentTime string) ([]entities.Events, error) {
	return s.repo.FindByContactAndDateAndTime(contactID, date, currentTime)
}

func (s *eventsServiceImpl) GetEventsByContactDateAndNumberPhone(contactID int64, date string, assistantID int64) ([]entities.Events, error) {
	// Validamos que la fecha tenga el formato adecuado
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, use YYYY-MM-DD: %v", err)
	}

	// Llamamos al repositorio para obtener los eventos
	events, err := s.repo.FindByContactDateAndNumberPhone(contactID, parsedDate.Format("2006-01-02"), assistantID)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *eventsServiceImpl) IsCodeUnique(code string) (bool, error) {
	unique, err := s.repo.ExistsByCode(code)
	if err != nil {
		return false, err
	}
	return unique, nil
}

// Crear un evento desde DTO
func (s *eventsServiceImpl) Create(eventDTO dtos.EventsDto) error {
	if eventDTO.Summary == "" || eventDTO.Description == "" {
		return errors.New("el resumen y la descripción son obligatorios")
	}

	// Convertir DTO a entidad
	event := entities.MapDtoToEvents(eventDTO)
	return s.repo.Create(&event)
}

// Obtener un evento por ID y devolver DTO
func (s *eventsServiceImpl) GetByID(id int) (*dtos.EventsDto, error) {
	event, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Convertir entidad a DTO
	eventDTO := entities.MapEntityToEventsDto(*event)
	return &eventDTO, nil
}

// Obtener todos los eventos y devolver DTOs
func (s *eventsServiceImpl) GetAll() ([]dtos.EventsDto, error) {
	events, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// Convertir lista de entidades a lista de DTOs
	var eventsDTOs []dtos.EventsDto
	for _, event := range events {
		eventsDTOs = append(eventsDTOs, entities.MapEntityToEventsDto(event))
	}

	return eventsDTOs, nil
}

// Actualizar un evento desde un DTO
func (s *eventsServiceImpl) Update(eventDTO dtos.EventsDto) error {
	if eventDTO.ID == 0 {
		return errors.New("el ID del evento es obligatorio")
	}

	// Convertir DTO a entidad
	event := entities.MapDtoToEvents(eventDTO)
	return s.repo.Update(&event)
}

// Eliminar un evento por ID
func (s *eventsServiceImpl) Cancel(codeEvent string) error {
	return s.repo.Cancel(codeEvent)
}

// Eliminar un evento por ID
func (s *eventsServiceImpl) Delete(id int) error {
	return s.repo.Delete(id)
}
