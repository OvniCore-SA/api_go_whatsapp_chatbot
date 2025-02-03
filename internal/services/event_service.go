package services

import (
	"errors"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
)

// Interfaz de servicio con DTOs
type EventsService interface {
	Create(eventDTO dtos.EventsDto) error
	GetByID(id int) (*dtos.EventsDto, error)
	GetAll() ([]dtos.EventsDto, error)
	Update(eventDTO dtos.EventsDto) error
	Delete(id int) error
}

// Implementación del servicio
type eventsServiceImpl struct {
	repo mysql_client.EventsRepository
}

func NewEventsService(repo mysql_client.EventsRepository) EventsService {
	return &eventsServiceImpl{repo: repo}
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
func (s *eventsServiceImpl) Delete(id int) error {
	return s.repo.Delete(id)
}
