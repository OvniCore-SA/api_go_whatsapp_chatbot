package entities

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"gorm.io/gorm"
)

type Promps struct {
	ID          int64          `gorm:"primaryKey"` // ID como clave primaria
	MetaAppsID  int64          // Clave foránea que apunta a MetaApps
	Name        string         // Nombre del prompt
	Descripcion string         // Descripción del prompt
	Activo      bool           // Estado del prompt
	CreatedAt   time.Time      // Fecha de creación
	UpdatedAt   time.Time      // Fecha de actualización
	DeletedAt   gorm.DeletedAt `gorm:"index"` // Fecha de eliminación suave (soft delete)
}

// MapEntitiesToPrompsDto convierte un entities.Promps a un dtos.PrompsDto
func MapEntitiesToPrompsDto(record Promps) dtos.PrompsDto {
	return dtos.PrompsDto{
		ID:          record.ID,
		Name:        record.Name,
		Descripcion: record.Descripcion,
		Activo:      record.Activo,
		MetaAppsId:  record.MetaAppsID,
		CreatedAt:   record.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   record.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// MapDtoToPromps convierte un dtos.PrompsDto a un entities.Promps
func MapDtoToPromps(dto dtos.PrompsDto) Promps {
	return Promps{
		// Mapear campos desde el DTO al dominio aquí
		Name:        dto.Name,
		MetaAppsID:  dto.MetaAppsId,
		Activo:      dto.Activo,
		Descripcion: dto.Descripcion,
	}
}
