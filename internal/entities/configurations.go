package entities

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"gorm.io/gorm"
)

type Configuration struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	KeyName     string         `gorm:"not null;unique"` // Clave única para identificar la configuración
	Value       string         `gorm:"type:text"`       // Valor de la configuración
	Description string         `gorm:"type:text"`       // Descripción opcional para explicar el propósito
	CreatedAt   time.Time      `gorm:"autoCreateTime"`  // Fecha de creación
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`  // Fecha de última actualización
	DeletedAt   gorm.DeletedAt `gorm:"index"`           // Soft delete
}

// MapConfigurationToDto convierte una entidad Configuration a un DTO
func MapConfigurationToDto(c Configuration) dtos.ConfigurationDto {
	return dtos.ConfigurationDto{
		ID:          c.ID,
		KeyName:     c.KeyName,
		Value:       c.Value,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		DeletedAt:   c.DeletedAt.Time,
	}
}

// MapDtoToConfiguration convierte un DTO a la entidad Configuration
func MapDtoToConfiguration(dto dtos.ConfigurationDto) Configuration {
	return Configuration{
		ID:          dto.ID,
		KeyName:     dto.KeyName,
		Value:       dto.Value,
		Description: dto.Description,
	}
}
