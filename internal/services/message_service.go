package services

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/postgres_client"
)

type MessagesService struct {
	repository *postgres_client.MessagesRepository
}

func NewMessagesService(repository *postgres_client.MessagesRepository) *MessagesService {
	return &MessagesService{repository: repository}
}

// GetMessagesByNumberPhone - Obtiene los mensajes asociados a un número de teléfono específico con paginación
func (s *MessagesService) GetMessagesByNumberPhone(numberPhoneID int64, page int, limit int) ([]dtos.MessageDto, int, error) {
	messages, total, err := s.repository.GetMessagesByNumberPhone(numberPhoneID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	dtos := make([]dtos.MessageDto, len(messages))
	for i, message := range messages {
		dtos[i] = entities.MapEntityToMessageDto(message)
	}

	return dtos, total, nil
}

// GetMessagesByNumberPhoneAndContact - Obtiene los mensajes asociados a un número de teléfono específico y un contacto con paginación
func (s *MessagesService) GetMessagesByNumberPhoneAndContact(numberPhoneID int64, contactID int64, page int, limit int) ([]dtos.MessageDto, int, error) {
	messages, total, err := s.repository.GetMessagesByNumberPhoneAndContact(numberPhoneID, contactID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	dtos := make([]dtos.MessageDto, len(messages))
	for i, message := range messages {
		dtos[i] = entities.MapEntityToMessageDto(message)
	}

	return dtos, total, nil
}

// DoesNumberPhoneExist - Verifica si un number_phones_id existe en la base de datos
func (s *MessagesService) DoesNumberPhoneExist(numberPhoneID int64) (bool, error) {
	return s.repository.DoesNumberPhoneExist(numberPhoneID)
}
