package dtos

type UsersDto struct {
	ID            int64
	Name          string
	Email         string
	Password      string
	RememberToken string
	Activo        bool
	Telefono      string
	CuilCuit      string
	ChatbotsID    int64
	RolesID       int64
	Rol           RolesDto `json:"rol"`
	CreatedAt     string
	UpdatedAt     string
	DeletedAt     string
}
