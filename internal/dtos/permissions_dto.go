package dtos

type PermissionsDto struct {
	ID          int64 // ID como clave primaria
	Permission  string
	Description string
	Roles       []RolesDto `json:"Roles,omitempty"`
	CreatedAt   string     `json:"CreatedAt,omitempty"` // Fecha de creación
	UpdatedAt   string     `json:"UpdatedAt,omitempty"` // Fecha de actualización
	DeletedAt   string     `json:"DeletedAt,omitempty"`
}
