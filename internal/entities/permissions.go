package entities

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"gorm.io/gorm"
)

type Permissions struct {
	ID          int64 `gorm:"primaryKey"` // ID como clave primaria
	Permission  string
	Description string
	Roles       []Roles        `gorm:"many2many:roles_has_permissions;"` // Relación de muchos a muchos
	CreatedAt   time.Time      // Fecha de creación
	UpdatedAt   time.Time      // Fecha de actualización
	DeletedAt   gorm.DeletedAt `gorm:"index"` // Fecha de eliminación suave (soft delete)
}

// MapEntitiesToPermissionsDto convierte un entities.Permissions a un dtos.PermissionsDto
func MapEntitiesToPermissionsDto(record Permissions) dtos.PermissionsDto {
	rolesDto := []dtos.RolesDto{}
	for _, role := range record.Roles {
		rolesDto = append(rolesDto, MapEntitiesToRolesDto(role))
	}

	return dtos.PermissionsDto{
		ID:          record.ID,
		Permission:  record.Permission,
		Description: record.Description,
		Roles:       rolesDto,
	}
}

// MapDtoToPermissions convierte un dtos.PermissionsDto a un entities.Permissions
func MapDtoToPermissions(dto dtos.PermissionsDto) Permissions {
	roles := []Roles{}
	for _, roleDto := range dto.Roles {
		roles = append(roles, MapDtoToRoles(roleDto))
	}

	return Permissions{
		ID:          dto.ID,
		Permission:  dto.Permission,
		Description: dto.Description,
		Roles:       roles,
	}
}
