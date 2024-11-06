package services

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
)

type LogsService struct {
	repository *mysql_client.LogsRepository
}

func NewLogsService(repository *mysql_client.LogsRepository) *LogsService {
	return &LogsService{repository: repository}
}

func (s *LogsService) GetAll() ([]dtos.LogsDto, error) {
	records, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	dtos := make([]dtos.LogsDto, len(records))
	for i, record := range records {
		dtos[i] = entities.MapEntitiesToLogsDto(record)
	}

	return dtos, nil
}

func (s *LogsService) GetById(id string) (dtos.LogsDto, error) {
	record, err := s.repository.FindByID(id)
	if err != nil {
		return dtos.LogsDto{}, err
	}

	return entities.MapEntitiesToLogsDto(record), nil
}

func (s *LogsService) Create(dto dtos.LogsDto) error {
	record := entities.MapDtoToLogs(dto)
	return s.repository.Create(record)
}

func (s *LogsService) Update(id string, dto dtos.LogsDto) error {
	record := entities.MapDtoToLogs(dto)
	return s.repository.Update(id, record)
}

func (s *LogsService) Delete(id string) error {
	return s.repository.Delete(id)
}
