package entities

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"gorm.io/gorm"
)

type Contact struct {
	ID                int64       `gorm:"primaryKey;autoIncrement"`
	NumberPhonesID    int64       `gorm:"not null"`
	NumberPhoneEntity NumberPhone `gorm:"foreignKey:NumberPhonesID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	NumberPhone       int64       `gorm:"not null"`
	IsBlocked         bool

	CountTokens string
	Events      []Events `gorm:"foreignKey:ContactsID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Relación con Events
	Threads     []Thread `gorm:"foreignKey:ContactsID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Relación de uno a muchos con Thread
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func MapEntityToContactDto(entity Contact) dtos.ContactDto {
	return dtos.ContactDto{
		ID:             entity.ID,
		NumberPhonesID: entity.NumberPhonesID,
		ContactNumber:  entity.NumberPhone,
		CountTokens:    entity.CountTokens,
		IsBlocked:      entity.IsBlocked,
	}
}

func MapDtoToContact(dto dtos.ContactDto) Contact {
	return Contact{
		ID:             dto.ID,
		NumberPhonesID: dto.NumberPhonesID,
		NumberPhone:    dto.ContactNumber,
		CountTokens:    dto.CountTokens,
		IsBlocked:      dto.IsBlocked,
	}
}
