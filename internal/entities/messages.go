package entities

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"gorm.io/gorm"
)

type Message struct {
	ID             int64          `gorm:"primaryKey;autoIncrement"`
	NumberPhonesID int64          `gorm:"not null"`                                                                // Clave foránea hacia Assistant
	NumberPhone    NumberPhone    `gorm:"foreignKey:NumberPhonesID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Relación con Assistant
	ContactsID     int64          `gorm:"not null"`                                                                // Clave foránea hacia Contact
	Contact        Contact        `gorm:"foreignKey:ContactsID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`     // Relación con Contact
	MessageText    string         `gorm:"type:text;not null"`                                                      // Texto del mensaje
	IsFromBot      bool           `gorm:"not null"`                                                                // Indica si el mensaje fue enviado por el bot (Assistant)
	CreatedAt      time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`                                      // Fecha de creación
	UpdatedAt      time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`          // Fecha de última actualización
	DeletedAt      gorm.DeletedAt `gorm:"index"`                                                                   // Soft delete
}

func MapEntityToMessageDto(entity Message) dtos.MessageDto {
	return dtos.MessageDto{
		ID:             entity.ID,
		NumberPhonesID: entity.NumberPhonesID,
		ContactsID:     entity.ContactsID,
		MessageText:    entity.MessageText,
		IsFromBot:      entity.IsFromBot,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
	}
}

func MapDtoToMessage(dto dtos.MessageDto) Message {
	return Message{
		ID:             dto.ID,
		NumberPhonesID: dto.NumberPhonesID,
		ContactsID:     dto.ContactsID,
		MessageText:    dto.MessageText,
		IsFromBot:      dto.IsFromBot,
		CreatedAt:      dto.CreatedAt,
		UpdatedAt:      dto.UpdatedAt,
	}
}
