package services

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
)

type ConfigurationsService struct {
	repository *mysql_client.ConfigurationsRepository
}

func NewConfigurationsService(repository *mysql_client.ConfigurationsRepository) *ConfigurationsService {
	return &ConfigurationsService{repository: repository}
}

// Create creates a new configuration record
func (s *ConfigurationsService) Create(dto dtos.ConfigurationDto) error {
	record := entities.MapDtoToConfiguration(dto)
	return s.repository.Create(record)
}

// FindByID retrieves a configuration record by its ID
func (s *ConfigurationsService) FindByID(id int64) (dtos.ConfigurationDto, error) {
	record, err := s.repository.FindByID(id)
	if err != nil {
		return dtos.ConfigurationDto{}, err
	}
	return entities.MapConfigurationToDto(record), nil
}

// FindByKey retrieves a configuration record by its key name
func (s *ConfigurationsService) FindByKey(keyName string) (dtos.ConfigurationDto, error) {
	record, err := s.repository.FindByKey(keyName)
	if err != nil {
		return dtos.ConfigurationDto{}, err
	}
	return entities.MapConfigurationToDto(record), nil
}

// GetAll retrieves all configuration records
func (s *ConfigurationsService) GetAll() ([]dtos.ConfigurationDto, error) {
	records, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	dtos := make([]dtos.ConfigurationDto, len(records))
	for i, record := range records {
		dtos[i] = entities.MapConfigurationToDto(record)
	}
	return dtos, nil
}

// Update updates an existing configuration record
func (s *ConfigurationsService) Update(dto dtos.ConfigurationDto) error {
	record := entities.MapDtoToConfiguration(dto)
	return s.repository.Update(record)
}

// Delete removes a configuration record by ID
func (s *ConfigurationsService) Delete(id int64) error {
	return s.repository.Delete(id)
}
