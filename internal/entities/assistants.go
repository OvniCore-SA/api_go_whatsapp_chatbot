package entities

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"gorm.io/gorm"
)

type Assistant struct {
	ID                       int64    `gorm:"primaryKey;autoIncrement"`
	BussinessID              int64    `gorm:"not null"`
	Bussiness                Bussines `gorm:"foreignKey:BussinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name                     string   `gorm:"not null"`
	OpenaiAssistantsID       string
	Description              string
	Model                    string
	Instructions             string
	Active                   bool                     `gorm:"not null;default:true"`
	NumberPhones             []NumberPhone            `gorm:"foreignKey:AssistantsID"`                                                // Relación de uno a muchos con NumberPhone
	GoogleCalendarCredential GoogleCalendarCredential `gorm:"foreignKey:AssistantsID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Relación uno a uno con GoogleCalendarCredential
	CreatedAt                time.Time
	UpdatedAt                time.Time
	DeletedAt                gorm.DeletedAt `gorm:"index"` // Soft delete
}

func MapAssistantToDto(a Assistant) dtos.AssistantDto {
	var bussiness dtos.BussinessDto
	if a.Bussiness.ID != 0 {
		bussiness = MapEntitiesToBussinessDto(a.Bussiness)
	}

	// Mapeo de GoogleCalendarCredential a DTO
	var googleCalendarCredential *dtos.GoogleCalendarCredentialDto
	if a.GoogleCalendarCredential.ID != 0 {
		googleCalendarCredential = &dtos.GoogleCalendarCredentialDto{
			ID:           a.GoogleCalendarCredential.ID,
			AssistantsID: a.GoogleCalendarCredential.AssistantsID,
			GoogleUserID: a.GoogleCalendarCredential.GoogleUserID,
			AccessToken:  a.GoogleCalendarCredential.AccessToken,
			RefreshToken: a.GoogleCalendarCredential.RefreshToken,
			TokenExpiry:  a.GoogleCalendarCredential.TokenExpiry,
		}
	}

	return dtos.AssistantDto{
		ID:                       a.ID,
		BussinessID:              a.BussinessID,
		Name:                     a.Name,
		OpenaiAssistantsID:       a.OpenaiAssistantsID,
		Description:              a.Description,
		Model:                    a.Model,
		Instructions:             a.Instructions,
		Active:                   a.Active,
		Bussiness:                bussiness, // cargo el bussiness
		GoogleCalendarCredential: googleCalendarCredential,
	}
}

func MapDtoToAssistant(dto dtos.AssistantDto) Assistant {
	var googleCalendarCredential GoogleCalendarCredential
	if dto.GoogleCalendarCredential != nil {
		googleCalendarCredential = GoogleCalendarCredential{
			ID:           dto.GoogleCalendarCredential.ID,
			AssistantsID: dto.GoogleCalendarCredential.AssistantsID,
			GoogleUserID: dto.GoogleCalendarCredential.GoogleUserID,
			AccessToken:  dto.GoogleCalendarCredential.AccessToken,
			RefreshToken: dto.GoogleCalendarCredential.RefreshToken,
			TokenExpiry:  dto.GoogleCalendarCredential.TokenExpiry,
		}
	}

	return Assistant{
		ID:                       dto.ID,
		BussinessID:              dto.BussinessID,
		Name:                     dto.Name,
		OpenaiAssistantsID:       dto.OpenaiAssistantsID,
		Description:              dto.Description,
		Model:                    dto.Model,
		Instructions:             dto.Instructions,
		Active:                   dto.Active,
		GoogleCalendarCredential: googleCalendarCredential,
	}
}
