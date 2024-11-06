package entities

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"gorm.io/gorm"
)

type Chatbots struct {
	ID                int64 `gorm:"primaryKey"` // ID como clave primaria
	Name              string
	NumberPhone       string
	Apikey            string
	PhoneNumberId     string
	TokenApiWhatsapp  string
	OptionsMenu       bool
	MetaAppsID        int64    // Clave foránea que apunta al Chatbot
	MetaApps          MetaApps `gorm:"foreignKey:MetaAppsID"`
	Resume            *Resumes `gorm:"foreignKey:ChatbotsID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Relación uno a uno con Resumes
	WhatsappCompanyId string
	GptUse            string         // Define con que GPT trabaja el chatbot (Disponible: 'CHATGPT', 'OLLAMA')
	CreatedAt         time.Time      // Fecha de creación
	UpdatedAt         time.Time      // Fecha de actualización
	DeletedAt         gorm.DeletedAt `gorm:"index"` // Fecha de eliminación suave (soft delete)
}

// MapEntitiesToChatbotsDto convierte un entities.Chatbots a un dtos.ChatbotsDto
func MapEntitiesToChatbotsDto(record Chatbots) dtos.ChatbotsDto {
	// Mapeo de campos desde la entidad Chatbots al DTO

	return dtos.ChatbotsDto{
		ID:                record.ID,
		Name:              record.Name,
		NumberPhone:       record.NumberPhone,
		Apikey:            record.Apikey,
		PhoneNumberId:     record.PhoneNumberId,
		TokenApiWhatsapp:  record.TokenApiWhatsapp,
		OptionsMenu:       record.OptionsMenu,
		WhatsappCompanyId: record.WhatsappCompanyId,
		GptUse:            record.GptUse,
		CreatedAt:         record.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         record.UpdatedAt.Format(time.RFC3339),
		DeletedAt:         record.DeletedAt.Time.Format(time.RFC3339), // Si es no nulo
		MetaApps:          MapEntitiesToMetaAppsDto(record.MetaApps),  // Mapeo de MetaApps a DTOs
	}
}

// MapDtoToChatbots convierte un dtos.ChatbotsDto a un entities.Chatbots
func MapDtoToChatbots(dto dtos.ChatbotsDto) Chatbots {
	createdAt, _ := time.Parse(time.RFC3339, dto.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, dto.UpdatedAt)

	return Chatbots{
		ID:                dto.ID,
		Name:              dto.Name,
		MetaAppsID:        dto.MetaApps.ID,
		NumberPhone:       dto.NumberPhone,
		Apikey:            dto.Apikey,
		PhoneNumberId:     dto.PhoneNumberId,
		OptionsMenu:       dto.OptionsMenu,
		TokenApiWhatsapp:  dto.TokenApiWhatsapp,
		WhatsappCompanyId: dto.WhatsappCompanyId,
		GptUse:            dto.GptUse,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}
