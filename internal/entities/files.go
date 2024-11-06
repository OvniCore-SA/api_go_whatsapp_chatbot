package entities

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID                       int64     `gorm:"primaryKey;autoIncrement"`
	AssistantsID             int64     `gorm:"not null"`                                                              // Clave foránea hacia Assistant
	Assistant                Assistant `gorm:"foreignKey:AssistantsID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Relación con Assistant
	OpenaiFilesID            string
	Filename                 string
	Purpose                  string
	OpenaiVectorStoreIDs     string
	OpenaiVectorStoreFileIDs string
	CreatedAt                time.Time
	UpdatedAt                time.Time
	DeletedAt                gorm.DeletedAt `gorm:"index"` // Soft delete
}
