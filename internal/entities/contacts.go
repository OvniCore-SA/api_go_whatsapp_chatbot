package entities

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"gorm.io/gorm"
)

type Contact struct {
	ID              int64       `gorm:"primaryKey;autoIncrement"`
	NumberPhonesID  int64       `gorm:"not null"`                                                                // Clave foránea hacia NumberPhone
	NumberPhone     NumberPhone `gorm:"foreignKey:NumberPhonesID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Relación con NumberPhone
	ContactNumber   int64       `gorm:"not null"`                                                                // Número de teléfono del contacto
	OpenaiThreadsID string
	CountTokens     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"` // Soft delete
}

func MapEntityToContactDto(entity Contact) dtos.ContactDto {
	return dtos.ContactDto{
		ID:              entity.ID,
		NumberPhonesID:  entity.NumberPhonesID,
		ContactNumber:   entity.ContactNumber,
		OpenaiThreadsID: entity.OpenaiThreadsID,
		CountTokens:     entity.CountTokens,
	}
}

func MapDtoToContact(dto dtos.ContactDto) Contact {
	return Contact{
		ID:              dto.ID,
		NumberPhonesID:  dto.NumberPhonesID,
		ContactNumber:   dto.ContactNumber,
		OpenaiThreadsID: dto.OpenaiThreadsID,
		CountTokens:     dto.CountTokens,
	}
}
