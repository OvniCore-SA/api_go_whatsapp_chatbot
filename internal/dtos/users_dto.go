package dtos

import (
	"errors"
	"regexp"
	"strings"
)

type UsersDto struct {
	ID            int64
	Name          string
	Email         string
	Password      string
	RememberToken string
	Activo        bool
	Telefono      string
	CuilCuit      string
	Dni           string
	RolesID       int64
	Rol           RolesDto `json:"rol"`
	CreatedAt     string
	UpdatedAt     string
	DeletedAt     string
}

func (dto *UsersDto) ValidateUsersDto(isCreate bool) error {
	// Validaciones comunes para creación y edición
	if strings.TrimSpace(dto.Name) == "" {
		return errors.New("el nombre es obligatorio")
	}
	if len(dto.Name) > 100 {
		return errors.New("el nombre no debe exceder los 100 caracteres")
	}

	if strings.TrimSpace(dto.Email) == "" {
		return errors.New("el email es obligatorio")
	}
	if !isValidEmail(dto.Email) {
		return errors.New("el email no tiene un formato válido")
	}

	if len(dto.Telefono) > 20 {
		return errors.New("el teléfono no debe exceder los 20 caracteres")
	}

	if len(dto.CuilCuit) > 20 {
		return errors.New("el CUIL/CUIT no debe exceder los 20 caracteres")
	}

	if len(dto.Dni) > 20 {
		return errors.New("el DNI no debe exceder los 20 caracteres")
	}

	// if dto.RolesID <= 0 {
	// 	return errors.New("el ID del rol es obligatorio y debe ser mayor que 0")
	// }

	// Validaciones específicas para creación
	if isCreate {
		if strings.TrimSpace(dto.Password) == "" {
			return errors.New("la contraseña es obligatoria")
		}
		if len(dto.Password) < 8 {
			return errors.New("la contraseña debe tener al menos 8 caracteres")
		}
	}

	// Validaciones específicas para edición
	if !isCreate {
		if dto.ID <= 0 {
			return errors.New("el ID es obligatorio al editar un usuario")
		}
		if strings.TrimSpace(dto.Password) != "" {
			return errors.New("no puedes actualizar la contraseña, solo puedes restablecerla.")
		}
	}

	return nil
}

// Helper para validar el formato de email
func isValidEmail(email string) bool {
	// Expresión regular para validar email
	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
