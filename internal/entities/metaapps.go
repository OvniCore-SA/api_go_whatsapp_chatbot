package entities

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"gorm.io/gorm"
)

// Entidad MetaApps
type MetaApps struct {
	ID           int64          `gorm:"primaryKey"` // ID como clave primaria
	Chatbots     []Chatbots     `gorm:"foreignKey:MetaAppsID"`
	Promps       Promps         `gorm:"foreignKey:MetaAppsID"` // Relación uno a uno con Promps
	Name         string         // Nombre de la app en meta
	AplicationId string         // Identificador de la app en meta
	Company      string         // Nombre de empresa en meta
	CreatedAt    time.Time      // Fecha de creación
	UpdatedAt    time.Time      // Fecha de actualización
	DeletedAt    gorm.DeletedAt `gorm:"index"` // Fecha de eliminación suave (soft delete)
}

// MapEntitiesToMetaAppsDto convierte un entities.MetaApps a un dtos.MetaAppsDto
func MapEntitiesToMetaAppsDto(record MetaApps) dtos.MetaAppsDto {
	chatbotsDto := []dtos.ChatbotsDto{}
	for _, chatbot := range record.Chatbots {
		chatbotsDto = append(chatbotsDto, MapEntitiesToChatbotsDto(chatbot))
	}

	prompDto := MapEntitiesToPrompsDto(record.Promps)

	return dtos.MetaAppsDto{
		ID:           record.ID,
		AplicationId: record.AplicationId,
		Name:         record.Name,
		Chatbots:     chatbotsDto,
		Company:      record.Company,
		Promps:       prompDto,
		CreatedAt:    record.CreatedAt.Format(time.RFC3339),      // Formatear como string
		UpdatedAt:    record.UpdatedAt.Format(time.RFC3339),      // Formatear como string
		DeletedAt:    record.DeletedAt.Time.Format(time.RFC3339), // Formatear DeletedAt si no es nulo
	}
}

// MapDtoToMetaApps convierte un dtos.MetaAppsDto a un entities.MetaApps
func MapDtoToMetaApps(dto dtos.MetaAppsDto) MetaApps {

	createdAt, _ := time.Parse(time.RFC3339, dto.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, dto.UpdatedAt)

	return MetaApps{
		ID:           dto.ID,
		AplicationId: dto.AplicationId,
		Name:         dto.Name,
		Company:      dto.Company,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

}
