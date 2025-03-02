package services

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
)

type ContactsService struct {
	repository *mysql_client.ContactsRepository
}

func NewContactsService(repository *mysql_client.ContactsRepository) *ContactsService {
	return &ContactsService{repository: repository}
}

func (s *ContactsService) GetAll() ([]dtos.ContactDto, error) {
	records, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	dtos := make([]dtos.ContactDto, len(records))
	for i, record := range records {
		dtos[i] = entities.MapEntityToContactDto(record)
	}

	return dtos, nil
}

func (s *ContactsService) GetById(id string) (dtos.ContactDto, error) {
	record, err := s.repository.FindByID(id)
	if err != nil {
		return dtos.ContactDto{}, err
	}
	return entities.MapEntityToContactDto(record), nil
}

func (s *ContactsService) GetByPhoneNumber(phoneNumber int64) (dtos.ContactDto, error) {
	record, err := s.repository.FindByPhoneNumber(phoneNumber)
	if err != nil {
		return dtos.ContactDto{}, err
	}
	return entities.MapEntityToContactDto(record), nil
}

func (s *ContactsService) Create(dto dtos.ContactDto) error {
	record := entities.MapDtoToContact(dto)
	return s.repository.Create(record)
}

func (s *ContactsService) Update(id string, dto dtos.ContactDto) error {
	record := entities.MapDtoToContact(dto)
	return s.repository.Update(id, record)
}

func (s *ContactsService) Delete(id string) error {
	return s.repository.Delete(id)
}

func (s *ContactsService) GetContactsByNumberPhone(numberPhoneID int64, page int, limit int) ([]dtos.ContactDto, int, error) {
	contacts, total, err := s.repository.GetContactsByNumberPhone(numberPhoneID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	dtos := make([]dtos.ContactDto, len(contacts))
	for i, contact := range contacts {
		dtos[i] = entities.MapEntityToContactDto(contact)
	}

	return dtos, total, nil
}
