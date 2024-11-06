package services

import (
	"strconv"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
)

type ResumesService struct {
	repository *mysql_client.ResumesRepository
}

func NewResumesService(repository *mysql_client.ResumesRepository) *ResumesService {
	return &ResumesService{repository: repository}
}

func (s *ResumesService) GetAll() ([]dtos.ResumesDto, error) {
	records, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	dtos := make([]dtos.ResumesDto, len(records))
	for i, record := range records {
		dtos[i] = entities.MapEntitiesToResumesDto(record)
	}

	return dtos, nil
}

// GetResumeByChatbotID retorna un resume basado en el ChatbotID
func (s *ResumesService) GetResumeByChatbotID(chatbotID int64) (*dtos.ResumesDto, error) {
	record, err := s.repository.GetResumeByChatbotID(chatbotID)
	if err != nil {
		return nil, err
	}

	// Convertir la entidad a DTO
	resumeDto := entities.MapEntitiesToResumesDto(*record)
	return &resumeDto, nil
}

func (s *ResumesService) GetByID(id string) (dtos.ResumesDto, error) {
	recordID, _ := strconv.ParseInt(id, 10, 64)
	record, err := s.repository.GetByID(recordID)
	if err != nil {
		return dtos.ResumesDto{}, err
	}

	return entities.MapEntitiesToResumesDto(record), nil
}

func (s *ResumesService) Create(dto dtos.ResumesDto) error {
	record := entities.MapDtoToResumes(dto)
	return s.repository.Create(record)
}

func (s *ResumesService) Update(id string, dto dtos.ResumesDto) error {
	recordID, _ := strconv.ParseInt(id, 10, 64)
	record := entities.MapDtoToResumes(dto)
	record.ID = recordID
	return s.repository.Update(record)
}

func (s *ResumesService) Delete(id string) error {
	recordID, _ := strconv.ParseInt(id, 10, 64)
	return s.repository.Delete(recordID)
}
