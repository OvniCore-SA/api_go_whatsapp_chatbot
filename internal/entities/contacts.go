package entities

import (
	"time"

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
