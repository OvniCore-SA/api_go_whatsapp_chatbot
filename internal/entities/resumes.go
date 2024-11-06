package entities

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"gorm.io/gorm"
)

type Resumes struct {
	ID               int64  `gorm:"primaryKey"`
	RequestToResolve string `gorm:"type:text"`
	ChatbotsID       int64
	Chatbot          Chatbots `gorm:"foreignKey:ChatbotsID"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

// MapEntitiesToResumesDto convierte un entities.Resumes a un dtos.ResumesDto
func MapEntitiesToResumesDto(record Resumes) dtos.ResumesDto {
	return dtos.ResumesDto{
		ID:               record.ID,
		RequestToResolve: record.RequestToResolve,
		ChatbotID:        record.ChatbotsID,
	}
}

// MapDtoToResumes convierte un dtos.ResumesDto a un entities.Resumes
func MapDtoToResumes(dto dtos.ResumesDto) Resumes {
	return Resumes{
		ID:               dto.ID,
		RequestToResolve: dto.RequestToResolve,
		ChatbotsID:       dto.ChatbotID,
	}
}
