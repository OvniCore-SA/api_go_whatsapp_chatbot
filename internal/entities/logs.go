package entities

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"gorm.io/gorm"
)

type Logs struct {
	ID        int64          `gorm:"primaryKey"` // ID como clave primaria
	CreatedAt time.Time      // Fecha de creación
	UpdatedAt time.Time      // Fecha de actualización
	DeletedAt gorm.DeletedAt `gorm:"index"` // Fecha de eliminación suave (soft delete)
	// Añadir más campos de la entidad aquí
}

// MapEntitiesToLogsDto convierte un entities.Logs a un dtos.LogsDto
func MapEntitiesToLogsDto(record Logs) dtos.LogsDto {
	return dtos.LogsDto{
		// Mapear campos desde el dominio al DTO aquí
	}
}

// MapDtoToLogs convierte un dtos.LogsDto a un entities.Logs
func MapDtoToLogs(dto dtos.LogsDto) Logs {
	return Logs{
		// Mapear campos desde el DTO al dominio aquí
	}
}
