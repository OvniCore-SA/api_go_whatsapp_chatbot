package dtos

type PermissionsDto struct {
	ID          int64 // ID como clave primaria
	Permission  string
	Description string
	Roles       []RolesDto
	CreatedAt   string // Fecha de creación
	UpdatedAt   string // Fecha de actualización
	DeletedAt   string
}
