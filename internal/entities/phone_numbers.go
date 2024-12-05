package entities

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"gorm.io/gorm"
)

type NumberPhone struct {
	ID                    int64     `gorm:"primaryKey;autoIncrement"`
	AssistantsID          int64     `gorm:"not null"`                                                              // Clave foránea hacia Assistant
	Assistant             Assistant `gorm:"foreignKey:AssistantsID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Relación con Assistant
	NumberPhone           int64     `gorm:"not null"`
	UUID                  string    `gorm:"not null;unique"`
	NumberPhoneToNotify   int64     `gorm:"default:0"`
	TokenPermanent        string    `gorm:"not null;unique"`
	WhatsappNumberPhoneId int64     `gorm:"not null;unique"`
	Active                bool      `gorm:"default:false"`             // Activo cuando el usuario escanea con éxito el QR
	Contacts              []Contact `gorm:"foreignKey:NumberPhonesID"` // Relación de uno a muchos con Contact
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             gorm.DeletedAt `gorm:"index"` // Soft delete
}

func MapEntityToNumberPhoneDto(entity NumberPhone) dtos.NumberPhoneDto {

	return dtos.NumberPhoneDto{
		ID:                    entity.ID,
		AssistantsID:          entity.AssistantsID,
		NumberPhone:           entity.NumberPhone,
		UUID:                  entity.UUID,
		TokenPermanent:        entity.TokenPermanent,
		NumberPhoneToNotify:   entity.NumberPhoneToNotify,
		WhatsappNumberPhoneId: entity.AssistantsID,
		Active:                entity.Active,
	}
}

func MapDtoToNumberPhone(dto dtos.NumberPhoneDto) NumberPhone {
	return NumberPhone{
		ID:                    dto.ID,
		AssistantsID:          dto.AssistantsID,
		NumberPhone:           dto.NumberPhone,
		UUID:                  dto.UUID,
		NumberPhoneToNotify:   dto.NumberPhoneToNotify,
		TokenPermanent:        dto.TokenPermanent,
		WhatsappNumberPhoneId: dto.WhatsappNumberPhoneId,
		Active:                dto.Active,
	}
}
