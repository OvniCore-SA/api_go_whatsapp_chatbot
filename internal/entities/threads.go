package entities

import (
	"time"

	"gorm.io/gorm"
)

type Thread struct {
	ID              int64  `gorm:"primaryKey;autoIncrement"` // ID como clave primaria
	OpenaiThreadsId string `gorm:"not null;unique"`
	Active          bool   `gorm:"not null"`
	ContactsID      int64  `gorm:"not null"` // FK a Contact

	Contact   Contact `gorm:"foreignKey:ContactsID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Relación con Contact
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete
}

// TableName establece el nombre de la tabla en la base de datos
func (Thread) TableName() string {
	return "threads"
}
