package entities

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"gorm.io/gorm"
)

type Users struct {
	ID             int64 `gorm:"primaryKey"` // ID como clave primaria
	Name           string
	Email          string
	Password       string
	RememberToken  string
	Activo         bool
	Telefono       string
	CuilCuit       string
	Bussiness      []Bussines       `gorm:"many2many:bussiness_has_users;"` // Relación muchos a muchos
	RolesID        int64            `gorm:"foreignKey:RolesID"`
	Rol            Roles            `gorm:"foreignKey:RolesID"`
	PasswordResets []PasswordResets `gorm:"foreignKey:UsersID"` // Relación de uno a muchos con PasswordResets
	CreatedAt      time.Time        // Fecha de creación
	UpdatedAt      time.Time        // Fecha de actualización
	DeletedAt      gorm.DeletedAt   `gorm:"index"` // Fecha de eliminación suave (soft delete)
}

// MapEntitiesToUsersDto convierte un entities.Users a un dtos.UsersDto
func MapEntitiesToUsersDto(record Users) dtos.UsersDto {
	// Mapeo de permisos en RolesDto
	permissions := []dtos.PermissionsDto{}
	for _, perm := range record.Rol.Permissions {
		permissions = append(permissions, dtos.PermissionsDto{
			ID:          perm.ID,
			Permission:  perm.Permission,
			Description: perm.Description,
		})
	}

	// Mapeo del rol en UsersDto
	rolDto := dtos.RolesDto{
		ID:          record.Rol.ID,
		Rol:         record.Rol.Rol,
		Description: record.Rol.Description,
		Permissions: permissions,
	}

	return dtos.UsersDto{
		ID:            record.ID,
		Name:          record.Name,
		Email:         record.Email,
		Password:      record.Password,
		RememberToken: record.RememberToken,
		Activo:        record.Activo,
		Telefono:      record.Telefono,
		CuilCuit:      record.CuilCuit,
		RolesID:       record.RolesID,
		Rol:           rolDto,
		CreatedAt:     record.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     record.UpdatedAt.Format(time.RFC3339),
		DeletedAt:     record.DeletedAt.Time.Format(time.RFC3339), // Manejo de soft delete
	}
}

// MapDtoToUsers convierte un dtos.UsersDto a un entities.Users
func MapDtoToUsers(dto dtos.UsersDto) Users {
	return Users{
		// Mapear campos desde el DTO al dominio aquí
		ID:            dto.ID,
		Name:          dto.Name,
		Email:         dto.Email,
		Password:      dto.Password,
		RememberToken: dto.RememberToken,
		Activo:        dto.Activo,
		Telefono:      dto.Telefono,
		CuilCuit:      dto.CuilCuit,
		RolesID:       dto.RolesID,
	}
}
