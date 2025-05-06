package services

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/postgres_client"
)

type RolesService struct {
	repository *postgres_client.RolesRepository
}

func NewRolesService(repository *postgres_client.RolesRepository) *RolesService {
	return &RolesService{repository: repository}
}

func (s *RolesService) GetAll() ([]dtos.RolesDto, error) {
	records, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	dtos := make([]dtos.RolesDto, len(records))
	for i, record := range records {
		dtos[i] = entities.MapEntitiesToRolesDto(record)
	}

	return dtos, nil
}

func (s *RolesService) GetById(id string) (dtos.RolesDto, error) {
	record, err := s.repository.FindByID(id)
	if err != nil {
		return dtos.RolesDto{}, err
	}

	return entities.MapEntitiesToRolesDto(record), nil
}

func (s *RolesService) GetByRol(rol string) (dtos.RolesDto, error) {
	rolDto, err := s.repository.GetByRol(rol)
	if err != nil {
		return dtos.RolesDto{}, err
	}

	return entities.MapEntitiesToRolesDto(rolDto), nil
}

func (s *RolesService) Create(dto dtos.RolesDto) error {
	record := entities.MapDtoToRoles(dto)
	return s.repository.Create(record)
}

func (s *RolesService) Update(id string, dto dtos.RolesDto) error {
	record := entities.MapDtoToRoles(dto)
	return s.repository.Update(id, record)
}

func (s *RolesService) Delete(id string) error {
	return s.repository.Delete(id)
}
