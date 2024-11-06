package services

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
)

type ChatbotsService struct {
	repository *mysql_client.ChatbotsRepository
}

func NewChatbotsService(repository *mysql_client.ChatbotsRepository) *ChatbotsService {
	return &ChatbotsService{repository: repository}
}

func (s *ChatbotsService) GetAll() ([]dtos.ChatbotsDto, error) {
	records, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	dtos := make([]dtos.ChatbotsDto, len(records))
	for i, record := range records {
		dtos[i] = entities.MapEntitiesToChatbotsDto(record)
	}

	return dtos, nil
}

// Buscar Chatbot por PhoneNumberID
func (s *ChatbotsService) GetChatbotByPhoneNumberID(phoneNumberID string) (chatbotDTO dtos.ChatbotsDto, err error) {
	chatbot, err := s.repository.FindByPhoneNumberID(phoneNumberID)
	if err != nil {
		return chatbotDTO, err
	}
	chatbotDTO = entities.MapEntitiesToChatbotsDto(*chatbot)
	return chatbotDTO, nil
}

func (s *ChatbotsService) GetById(id int64) (dtos.ChatbotsDto, error) {
	record, err := s.repository.FindByID(id)
	if err != nil {
		return dtos.ChatbotsDto{}, err
	}

	return entities.MapEntitiesToChatbotsDto(record), nil
}

func (s *ChatbotsService) Create(dto dtos.ChatbotsDto) error {
	record := entities.MapDtoToChatbots(dto)
	return s.repository.Create(record)
}

func (s *ChatbotsService) Update(id string, dto dtos.ChatbotsDto) error {
	record := entities.MapDtoToChatbots(dto)
	return s.repository.Update(id, record)
}

func (s *ChatbotsService) Delete(id string) error {
	return s.repository.Delete(id)
}
