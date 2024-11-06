package services

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities/filters"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
)

type PrompsService struct {
	repository *mysql_client.PrompsRepository
}

func NewPrompsService(repository *mysql_client.PrompsRepository) *PrompsService {
	return &PrompsService{repository: repository}
}

// GetByFilter busca un prompt en base a los filtros proporcionados
func (s *PrompsService) GetByFilter(filter filters.PrompsFiltro) (dtos.PrompsDto, error) {
	record, err := s.repository.GetByFilter(filter)
	if err != nil {
		return dtos.PrompsDto{}, err
	}

	return entities.MapEntitiesToPrompsDto(record), nil
}

func (s *PrompsService) GetAll() ([]dtos.PrompsDto, error) {
	records, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	dtos := make([]dtos.PrompsDto, len(records))
	for i, record := range records {
		dtos[i] = entities.MapEntitiesToPrompsDto(record)
	}

	return dtos, nil
}

func (s *PrompsService) GetById(id int64) (dtos.PrompsDto, error) {
	record, err := s.repository.FindByID(id)
	if err != nil {
		return dtos.PrompsDto{}, err
	}

	return entities.MapEntitiesToPrompsDto(record), nil
}

func (s *PrompsService) Create(dto dtos.PrompsDto) error {
	record := entities.MapDtoToPromps(dto)
	return s.repository.Create(record)
}

func (s *PrompsService) Update(id string, dto dtos.PrompsDto) error {
	record := entities.MapDtoToPromps(dto)
	return s.repository.Update(id, record)
}

func (s *PrompsService) Delete(id string) error {
	return s.repository.Delete(id)
}
