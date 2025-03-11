package entities

import (
	"time"

	"gorm.io/gorm"
)

type BussinessHasUsers struct {
	BussinessID int64 `gorm:"primaryKey"` // Clave primaria compuesta
	UsersID     int64 `gorm:"primaryKey"` // Clave primaria compuesta
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"` // Soft delete
}
