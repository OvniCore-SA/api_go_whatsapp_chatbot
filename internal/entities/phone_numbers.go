package entities

import (
	"time"

	"gorm.io/gorm"
)

type NumberPhone struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	AssistantsID int64     `gorm:"not null"`                                                              // Clave foránea hacia Assistant
	Assistant    Assistant `gorm:"foreignKey:AssistantsID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Relación con Assistant
	PhoneNumber  int64     `gorm:"not null"`
	UUID         string    `gorm:"not null;unique"`
	Active       bool      `gorm:"default:false"`             // Activo cuando el usuario escanea con éxito el QR
	Contacts     []Contact `gorm:"foreignKey:NumberPhonesID"` // Relación de uno a muchos con Contact
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"` // Soft delete
}
