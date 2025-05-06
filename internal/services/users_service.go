package services

import (
	"errors"
	"os"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/postgres_client"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	repository   *postgres_client.UsersRepository
	rolesService *RolesService
}

func NewUsersService(repository *postgres_client.UsersRepository, rolesService *RolesService) *UsersService {
	return &UsersService{repository: repository, rolesService: rolesService}
}

func (s *UsersService) GetAll() ([]dtos.UsersDto, error) {
	records, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	dtos := make([]dtos.UsersDto, len(records))
	for i, record := range records {
		dtos[i] = entities.MapEntitiesToUsersDto(record)
	}

	return dtos, nil
}

func (s *UsersService) GetById(id int64) (dtos.UsersDto, error) {
	record, err := s.repository.FindByID(uint(id))
	if err != nil {
		return dtos.UsersDto{}, err
	}

	return entities.MapEntitiesToUsersDto(*record), nil
}

func (s *UsersService) Create(dto dtos.UsersDto) error {
	// Validar que la contraseña no esté vacía
	if dto.Password == "" {
		return errors.New("password is required")
	}

	// Asigno el rol normal (USER).
	rolUser, _ := s.rolesService.GetByRol(os.Getenv("ROL_USER"))
	dto.RolesID = rolUser.ID

	// Encriptar la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Reemplazar la contraseña en el DTO con la versión encriptada
	dto.Password = string(hashedPassword)

	// Mapear el DTO a la entidad
	record := entities.MapDtoToUsers(dto)

	// Guardar el usuario en el repositorio
	return s.repository.Create(record)
}

func (s *UsersService) Update(id int64, dto dtos.UsersDto) error {
	record := entities.MapDtoToUsers(dto)
	return s.repository.Update(id, record)
}

func (s *UsersService) Delete(id int64) error {
	return s.repository.Delete(id)
}

// FindByEmail busca un usuario por email (Carga el rol y los permisos)
func (s *UsersService) FindByEmail(email string) (dtos.UsersDto, error) {
	record, err := s.repository.FindByEmail(email)
	if err != nil {
		return dtos.UsersDto{}, err
	}

	return entities.MapEntitiesToUsersDto(record), nil
}
