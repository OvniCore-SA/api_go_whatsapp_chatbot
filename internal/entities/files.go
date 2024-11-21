package entities

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
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

func MapEntityToFileDto(entity File) dtos.FileDto {
	return dtos.FileDto{
		ID:                       entity.ID,
		AssistantsID:             entity.AssistantsID,
		OpenaiFilesID:            entity.OpenaiFilesID,
		Filename:                 entity.Filename,
		Purpose:                  entity.Purpose,
		OpenaiVectorStoreIDs:     entity.OpenaiVectorStoreIDs,
		OpenaiVectorStoreFileIDs: entity.OpenaiVectorStoreFileIDs,
	}
}

func MapDtoToFile(dto dtos.FileDto) File {
	return File{
		ID:                       dto.ID,
		AssistantsID:             dto.AssistantsID,
		OpenaiFilesID:            dto.OpenaiFilesID,
		Filename:                 dto.Filename,
		Purpose:                  dto.Purpose,
		OpenaiVectorStoreIDs:     dto.OpenaiVectorStoreIDs,
		OpenaiVectorStoreFileIDs: dto.OpenaiVectorStoreFileIDs,
	}
}
