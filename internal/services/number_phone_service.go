package services

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
)

type NumberPhonesService struct {
	repository *mysql_client.NumberPhonesRepository
}

func NewNumberPhonesService(repository *mysql_client.NumberPhonesRepository) *NumberPhonesService {
	return &NumberPhonesService{repository: repository}
}

func (s *NumberPhonesService) GetAll() ([]dtos.NumberPhoneDto, error) {
	records, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	dtos := make([]dtos.NumberPhoneDto, len(records))
	for i, record := range records {
		dtos[i] = entities.MapEntityToNumberPhoneDto(record)
	}

	return dtos, nil
}

func (s *NumberPhonesService) GetById(id string) (dtos.NumberPhoneDto, error) {
	record, err := s.repository.FindByID(id)
	if err != nil {
		return dtos.NumberPhoneDto{}, err
	}
	return entities.MapEntityToNumberPhoneDto(record), nil
}

func (s *NumberPhonesService) Create(dto dtos.NumberPhoneDto) error {
	record := entities.MapDtoToNumberPhone(dto)
	return s.repository.Create(record)
}

func (s *NumberPhonesService) Update(id string, dto dtos.NumberPhoneDto) error {
	record := entities.MapDtoToNumberPhone(dto)
	return s.repository.Update(id, record)
}

func (s *NumberPhonesService) Delete(id string) error {
	return s.repository.Delete(id)
}