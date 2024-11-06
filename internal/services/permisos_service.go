package services

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
)

type PermissionsService struct {
	repository *mysql_client.PermissionsRepository
}

func NewPermissionsService(repository *mysql_client.PermissionsRepository) *PermissionsService {
	return &PermissionsService{repository: repository}
}

func (s *PermissionsService) GetAll() ([]dtos.PermissionsDto, error) {
	records, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	dtos := make([]dtos.PermissionsDto, len(records))
	for i, record := range records {
		dtos[i] = entities.MapEntitiesToPermissionsDto(record)
	}

	return dtos, nil
}

func (s *PermissionsService) GetById(id string) (dtos.PermissionsDto, error) {
	record, err := s.repository.FindByID(id)
	if err != nil {
		return dtos.PermissionsDto{}, err
	}

	return entities.MapEntitiesToPermissionsDto(record), nil
}

func (s *PermissionsService) Create(dto dtos.PermissionsDto) error {
	record := entities.MapDtoToPermissions(dto)
	return s.repository.Create(record)
}

func (s *PermissionsService) Update(id string, dto dtos.PermissionsDto) error {
	record := entities.MapDtoToPermissions(dto)
	return s.repository.Update(id, record)
}

func (s *PermissionsService) Delete(id string) error {
	return s.repository.Delete(id)
}
