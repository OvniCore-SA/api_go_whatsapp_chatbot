package dtos

type RolesDto struct {
	ID          int64
	Rol         string
	Description string
	Permissions []PermissionsDto
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
}
