package dtos

type RolesDto struct {
	ID          int64
	Rol         string
	Description string
	Permissions []PermissionsDto `json:"Permissions,omitempty"`
	CreatedAt   string           `json:"CreatedAt,omitempty"`
	UpdatedAt   string           `json:"UpdatedAt,omitempty"`
	DeletedAt   string           `json:"DeletedAt,omitempty"`
}
