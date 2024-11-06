package entities

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"gorm.io/gorm"
)

type Roles struct {
	ID          int64 `gorm:"primaryKey"` // ID como clave primaria
	Rol         string
	Description string
	Users       []Users        `gorm:"foreignKey:RolesID"`
	Permissions []Permissions  `gorm:"many2many:roles_has_permissions;"` // Relación de muchos a muchos
	CreatedAt   time.Time      // Fecha de creación
	UpdatedAt   time.Time      // Fecha de actualización
	DeletedAt   gorm.DeletedAt `gorm:"index"` // Fecha de eliminación suave (soft delete)
}

// MapEntitiesToRolesDto convierte un entities.Roles a un dtos.RolesDto
func MapEntitiesToRolesDto(record Roles) dtos.RolesDto {
	permissionsDto := []dtos.PermissionsDto{}
	for _, permission := range record.Permissions {
		permissionsDto = append(permissionsDto, MapEntitiesToPermissionsDto(permission))
	}

	return dtos.RolesDto{
		ID:          record.ID,
		Rol:         record.Rol,
		Description: record.Description,
		Permissions: permissionsDto,
	}
}

// MapDtoToRoles convierte un dtos.RolesDto a un entities.Roles
func MapDtoToRoles(dto dtos.RolesDto) Roles {
	permissions := []Permissions{}
	for _, permissionDto := range dto.Permissions {
		permissions = append(permissions, MapDtoToPermissions(permissionDto))
	}

	return Roles{
		ID:          dto.ID,
		Rol:         dto.Rol,
		Description: dto.Description,
		Permissions: permissions,
	}
}
