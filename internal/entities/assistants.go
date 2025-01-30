package entities

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"gorm.io/gorm"
)

type Assistant struct {
	ID                   int64    `gorm:"primaryKey;autoIncrement"`
	BussinessID          int64    `gorm:"not null"`
	Bussiness            Bussines `gorm:"foreignKey:BussinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name                 string   `gorm:"not null"`
	OpenaiAssistantsID   string
	Description          string
	Model                string
	Instructions         string
	Active               bool                 `gorm:"not null;default:true"`
	AccountGoogle        bool                 `gorm:"default:false"`
	NumberPhones         []NumberPhone        `gorm:"foreignKey:AssistantsID"`                                                // Relación de uno a muchos con NumberPhone
	GoogleCalendarConfig GoogleCalendarConfig `gorm:"foreignKey:AssistantsID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Relación uno a uno con GoogleCalendarConfig
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt `gorm:"index"` // Soft delete
}

func MapAssistantToDto(a Assistant) dtos.AssistantDto {
	var bussiness dtos.BussinessDto
	if a.Bussiness.ID != 0 {
		bussiness = MapEntitiesToBussinessDto(a.Bussiness)
	}

	// Mapeo de GoogleCalendarConfig a DTO
	var googleCalendarCredential *dtos.GoogleCalendarConfigDto
	if a.GoogleCalendarConfig.ID != 0 {
		googleCalendarCredential = &dtos.GoogleCalendarConfigDto{
			ID:           a.GoogleCalendarConfig.ID,
			AssistantsID: a.GoogleCalendarConfig.AssistantsID,
			GoogleUserID: a.GoogleCalendarConfig.GoogleUserID,
			AccessToken:  a.GoogleCalendarConfig.AccessToken,
			RefreshToken: a.GoogleCalendarConfig.RefreshToken,
			TokenExpiry:  a.GoogleCalendarConfig.TokenExpiry,
		}
	}

	return dtos.AssistantDto{
		ID:                   a.ID,
		BussinessID:          a.BussinessID,
		Name:                 a.Name,
		OpenaiAssistantsID:   a.OpenaiAssistantsID,
		Description:          a.Description,
		Model:                a.Model,
		Instructions:         a.Instructions,
		Active:               a.Active,
		Bussiness:            bussiness, // cargo el bussiness
		GoogleCalendarConfig: googleCalendarCredential,
		AccountGoogle:        a.AccountGoogle,
	}
}

func MapDtoToAssistant(dto dtos.AssistantDto) Assistant {
	var googleCalendarCredential GoogleCalendarConfig
	if dto.GoogleCalendarConfig != nil {
		googleCalendarCredential = GoogleCalendarConfig{
			ID:           dto.GoogleCalendarConfig.ID,
			AssistantsID: dto.GoogleCalendarConfig.AssistantsID,
			GoogleUserID: dto.GoogleCalendarConfig.GoogleUserID,
			AccessToken:  dto.GoogleCalendarConfig.AccessToken,
			RefreshToken: dto.GoogleCalendarConfig.RefreshToken,
			TokenExpiry:  dto.GoogleCalendarConfig.TokenExpiry,
		}
	}

	return Assistant{
		ID:                   dto.ID,
		BussinessID:          dto.BussinessID,
		Name:                 dto.Name,
		OpenaiAssistantsID:   dto.OpenaiAssistantsID,
		Description:          dto.Description,
		Model:                dto.Model,
		Instructions:         dto.Instructions,
		Active:               dto.Active,
		GoogleCalendarConfig: googleCalendarCredential,
	}
}
