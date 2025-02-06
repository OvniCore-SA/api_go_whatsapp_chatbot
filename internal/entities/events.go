package entities

import (
	"fmt"
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"gorm.io/gorm"
)

type Events struct {
	ID                    int    `gorm:"primaryKey"`
	Summary               string `gorm:"not null"`
	Description           string `gorm:"not null"`
	StartDate             string `gorm:"type:text;not null"`
	EndDate               string `gorm:"type:text;not null"`
	EventGoogleCalendarID string
	CodeEvent             string `gorm:"type:text;not null"`

	AssistantsID int64     `gorm:"not null"` // Relación con Assistant (un asistente tiene muchos eventos)
	Assistant    Assistant `gorm:"foreignKey:AssistantsID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ContactsID int64   `gorm:"not null"` // Relación con Contact (un contacto tiene muchos eventos)
	Contact    Contact `gorm:"foreignKey:ContactsID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func MapEntityToEventsDto(entity Events) dtos.EventsDto {
	createdAtToString := entity.CreatedAt.Format(time.RFC3339)
	return dtos.EventsDto{
		ID:                    entity.ID,
		Summary:               entity.Summary,
		Description:           entity.Description,
		StartDate:             entity.StartDate,
		EndDate:               entity.EndDate,
		EventGoogleCalendarID: entity.EventGoogleCalendarID,
		AssistantsID:          entity.AssistantsID,
		ContactsID:            entity.ContactsID,
		CodeEvent:             entity.CodeEvent,
		CreatedAt:             createdAtToString,
	}
}

func MapDtoToEvents(dto dtos.EventsDto) Events {
	var createdAtToTime time.Time
	if len(dto.CreatedAt) > 0 {
		var err error
		createdAtToTime, err = time.Parse("2006-01-02T15:04:05-03:00", dto.CreatedAt)
		if err != nil {
			fmt.Println("Error al parsear la fecha:", err)
			return Events{}
		}
	}
	return Events{
		ID:                    dto.ID,
		Summary:               dto.Summary,
		Description:           dto.Description,
		StartDate:             dto.StartDate,
		EndDate:               dto.EndDate,
		EventGoogleCalendarID: dto.EventGoogleCalendarID,
		AssistantsID:          dto.AssistantsID,
		ContactsID:            dto.ContactsID,
		CodeEvent:             dto.CodeEvent,
		CreatedAt:             createdAtToTime,
	}
}
