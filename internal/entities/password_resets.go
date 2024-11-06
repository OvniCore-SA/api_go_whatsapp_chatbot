package entities

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"gorm.io/gorm"
)

type PasswordResets struct {
	ID        int64          `gorm:"primaryKey"`         // ID como clave primaria
	UsersID   int64          `gorm:"not null"`           // ID del usuario que solicita el restablecimiento de contraseña (clave foránea)
	User      Users          `gorm:"foreignKey:UsersID"` // Relación con el modelo Users
	Token     string         `gorm:"unique"`             // Token único para el restablecimiento de contraseña
	CreatedAt time.Time      // Fecha de creación
	UpdatedAt time.Time      // Fecha de actualización
	DeletedAt gorm.DeletedAt `gorm:"index"` // Fecha de eliminación suave (soft delete)
}

// MapEntitiesToPassword_resetsDto convierte un entities.Password_resets a un dtos.Password_resetsDto
func MapEntitiesToPassword_resetsDto(record PasswordResets) dtos.Password_resetsDto {
	return dtos.Password_resetsDto{
		// Mapear campos desde el dominio al DTO aquí
	}
}

// MapDtoToPassword_resets convierte un dtos.Password_resetsDto a un entities.Password_resets
func MapDtoToPassword_resets(dto dtos.Password_resetsDto) PasswordResets {
	return PasswordResets{
		// Mapear campos desde el DTO al dominio aquí
	}
}
