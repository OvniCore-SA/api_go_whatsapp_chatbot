package services

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/postgres_client"
)

type Password_resetsService struct {
	repository *postgres_client.PasswordResetsRepository
}

func NewPassword_resetsService(repository *postgres_client.PasswordResetsRepository) *Password_resetsService {
	return &Password_resetsService{repository: repository}
}

func (s *Password_resetsService) GetAll() ([]dtos.Password_resetsDto, error) {
	records, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	dtos := make([]dtos.Password_resetsDto, len(records))
	for i, record := range records {
		dtos[i] = entities.MapEntitiesToPassword_resetsDto(record)
	}

	return dtos, nil
}

func (s *Password_resetsService) GetById(id string) (dtos.Password_resetsDto, error) {
	record, err := s.repository.FindByID(id)
	if err != nil {
		return dtos.Password_resetsDto{}, err
	}

	return entities.MapEntitiesToPassword_resetsDto(record), nil
}

func (s *Password_resetsService) Create(dto dtos.Password_resetsDto) error {
	record := entities.MapDtoToPassword_resets(dto)
	return s.repository.Create(record)
}

func (s *Password_resetsService) Update(id string, dto dtos.Password_resetsDto) error {
	record := entities.MapDtoToPassword_resets(dto)
	return s.repository.Update(id, record)
}

func (s *Password_resetsService) Delete(id string) error {
	return s.repository.Delete(id)
}
