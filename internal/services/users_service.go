package services

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
)

type UsersService struct {
	repository *mysql_client.UsersRepository
}

func NewUsersService(repository *mysql_client.UsersRepository) *UsersService {
	return &UsersService{repository: repository}
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
	record, err := s.repository.FindByID(id)
	if err != nil {
		return dtos.UsersDto{}, err
	}

	return entities.MapEntitiesToUsersDto(record), nil
}

func (s *UsersService) Create(dto dtos.UsersDto) error {
	record := entities.MapDtoToUsers(dto)
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
