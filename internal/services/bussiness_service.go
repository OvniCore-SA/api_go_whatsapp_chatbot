package services

import (
	"errors"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
)

type BussinessService struct {
	repository *mysql_client.BussinessRepository
}

func NewBussinessService(repository *mysql_client.BussinessRepository) *BussinessService {
	return &BussinessService{repository: repository}
}

// Crear un nuevo negocio
func (s *BussinessService) CreateBussiness(data dtos.BussinessDto) (dtos.BussinessDto, error) {
	// Convertir el DTO a entidad
	bussiness := entities.MapDtoToBussiness(data)

	// Llamar al repositorio para crear el negocio
	idBussiness, err := s.repository.Create(bussiness)
	if err != nil {
		return dtos.BussinessDto{}, err
	}

	bussinessDTO := entities.MapEntitiesToBussinessDto(bussiness)
	bussinessDTO.ID = int64(idBussiness)
	// Devolver el DTO del negocio creado
	return bussinessDTO, nil
}

// Obtener un negocio por ID
func (s *BussinessService) GetBussinessById(id int64) (dtos.BussinessDto, error) {
	bussiness, err := s.repository.FindByID(id)
	if err != nil {
		return dtos.BussinessDto{}, errors.New("bussiness no encontrado")
	}
	return entities.MapEntitiesToBussinessDto(bussiness), nil
}

// Obtener todos los negocios
func (s *BussinessService) GetAllBussinesses() ([]dtos.BussinessDto, error) {
	records, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	// Convertir las entidades a DTOs
	dtos := make([]dtos.BussinessDto, len(records))
	for i, record := range records {
		dtos[i] = entities.MapEntitiesToBussinessDto(record)
	}

	return dtos, nil
}

// Actualizar un negocio
func (s *BussinessService) UpdateBussiness(id int64, data dtos.BussinessDto) (dtos.BussinessDto, error) {
	// Convertir el DTO a entidad
	bussiness := entities.MapDtoToBussiness(data)

	// Llamar al repositorio para actualizar el negocio
	if err := s.repository.Update(id, bussiness); err != nil {
		return dtos.BussinessDto{}, errors.New("bussiness no encontrado")
	}

	// Retornar el DTO del negocio actualizado
	return entities.MapEntitiesToBussinessDto(bussiness), nil
}

// Eliminar un negocio
func (s *BussinessService) DeleteBussiness(id int64) error {
	if err := s.repository.Delete(id); err != nil {
		return errors.New("bussiness no encontrado")
	}
	return nil
}

// Obtener negocios por ID de usuario
func (s *BussinessService) FindByUserId(userId int64) ([]dtos.BussinessDto, error) {
	records, err := s.repository.FindByUserId(userId)
	if err != nil {
		return nil, errors.New("el usuario no posee bussiness registrados")
	}

	// Convertir las entidades a DTOs
	dtos := make([]dtos.BussinessDto, len(records))
	for i, record := range records {
		dtos[i] = entities.MapEntitiesToBussinessDto(record)
	}

	return dtos, nil
}
