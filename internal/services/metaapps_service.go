package services

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
)

type MetaAppsService struct {
	repository *mysql_client.MetaAppsRepository
}

func NewMetaAppsService(repository *mysql_client.MetaAppsRepository) *MetaAppsService {
	return &MetaAppsService{repository: repository}
}

func (s *MetaAppsService) GetAll() ([]dtos.MetaAppsDto, error) {
	records, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	dtos := make([]dtos.MetaAppsDto, len(records))
	for i, record := range records {
		dtos[i] = entities.MapEntitiesToMetaAppsDto(record)
	}

	return dtos, nil
}

func (s *MetaAppsService) GetById(id int64) (dtos.MetaAppsDto, error) {
	record, err := s.repository.FindByID(id)
	if err != nil {
		return dtos.MetaAppsDto{}, err
	}

	return entities.MapEntitiesToMetaAppsDto(record), nil
}

func (s *MetaAppsService) Create(dto dtos.MetaAppsDto) error {
	record := entities.MapDtoToMetaApps(dto)
	return s.repository.Create(record)
}

func (s *MetaAppsService) Update(id string, dto dtos.MetaAppsDto) error {
	record := entities.MapDtoToMetaApps(dto)
	return s.repository.Update(id, record)
}

func (s *MetaAppsService) Delete(id string) error {
	return s.repository.Delete(id)
}
