package entities

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"gorm.io/gorm"
)

type Bussines struct {
	ID         int64  `gorm:"primaryKey;autoIncrement"` // ID como clave primaria
	Name       string `gorm:"not null"`
	Address    string `gorm:"not null"`
	CuilCuit   string
	WebSite    string
	Users      []Users     `gorm:"many2many:bussiness_has_users;"` // Relación muchos a muchos
	Assistants []Assistant `gorm:"foreignKey:BussinessID"`         // Relación de uno a muchos con Assistant
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"` // Soft delete
}

// TableName establece el nombre de la tabla en la base de datos
func (Bussines) TableName() string {
	return "bussiness"
}

// MapEntitiesToBussinessDto convierte un entities.Bussiness a un dtos.BussinessDto
func MapEntitiesToBussinessDto(record Bussines) dtos.BussinessDto {
	// Mapeo de asistentes en AssistantDto
	assistants := []dtos.AssistantDto{}
	for _, assistant := range record.Assistants {
		assistants = append(assistants, dtos.AssistantDto{
			ID:                 assistant.ID,
			Name:               assistant.Name,
			OpenaiAssistantsID: assistant.OpenaiAssistantsID,
			Description:        assistant.Description,
			Model:              assistant.Model,
			Active:             assistant.Active,
		})
	}

	usersDto := []dtos.UsersDto{}
	for _, u := range record.Users {
		usersDto = append(usersDto, dtos.UsersDto{
			ID:            u.ID,
			Name:          u.Name,
			Email:         u.Email,
			Password:      u.Password,
			CuilCuit:      u.CuilCuit,
			Activo:        u.Activo,
			RememberToken: u.RememberToken,
			RolesID:       u.RolesID,
			Telefono:      u.Telefono,
			CreatedAt:     u.CreatedAt.String(),
			UpdatedAt:     u.UpdatedAt.String(),
		})
	}

	return dtos.BussinessDto{
		ID:         record.ID,
		Name:       record.Name,
		Address:    record.Address,
		CuilCuit:   record.CuilCuit,
		Users:      usersDto,
		WebSite:    record.WebSite,
		Assistants: assistants,
		CreatedAt:  record.CreatedAt,
		UpdatedAt:  record.UpdatedAt,
	}
}

// MapDtoToBussiness convierte un dtos.BussinessDto a un entities.Bussiness
func MapDtoToBussiness(dto dtos.BussinessDto) Bussines {
	// Mapeo de asistentes desde AssistantDto
	assistants := []Assistant{}
	for _, assistantDto := range dto.Assistants {
		assistants = append(assistants, Assistant{
			ID:                 assistantDto.ID,
			Name:               assistantDto.Name,
			OpenaiAssistantsID: assistantDto.OpenaiAssistantsID,
			Description:        assistantDto.Description,
			Model:              assistantDto.Model,
			Active:             assistantDto.Active,
		})
	}

	return Bussines{
		ID:         dto.ID,
		Name:       dto.Name,
		Address:    dto.Address,
		CuilCuit:   dto.CuilCuit,
		WebSite:    dto.WebSite,
		Assistants: assistants,
		CreatedAt:  dto.CreatedAt,
		UpdatedAt:  dto.UpdatedAt,
	}
}
