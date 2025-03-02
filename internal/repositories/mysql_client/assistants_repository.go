package mysql_client

import (
	"errors"
	"fmt"
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"gorm.io/gorm"
)

type AssistantRepository struct {
	db *gorm.DB
}

func NewAssistantRepository(db *gorm.DB) *AssistantRepository {
	return &AssistantRepository{db: db}
}

func (r *AssistantRepository) Create(data *entities.Assistant) error {
	// GORM automáticamente asigna el ID a data.ID después de la creación
	if err := r.db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

// IsWithinWorkingHours verifica si una fecha y hora están dentro del horario de atención de un asistente
func (r *AssistantRepository) IsWithinWorkingHours(assistantID int64, dateTime time.Time) (bool, error) {
	var assistant entities.Assistant
	err := r.db.First(&assistant, assistantID).Error
	if err != nil {
		return false, fmt.Errorf("assistant not found: %v", err)
	}

	// Obtener el día de la semana (0 = Domingo, 1 = Lunes, ..., 6 = Sábado)
	weekday := int(dateTime.Weekday())

	// Verificar si el bit del día está activado en OpeningDays
	if (assistant.OpeningDays & (1 << weekday)) == 0 {
		return false, nil // El asistente no trabaja en este día
	}

	// Parsear el horario de trabajo
	var openHour, openMin, closeHour, closeMin int
	_, err = fmt.Sscanf(assistant.WorkingHours, "%d:%d-%d:%d", &openHour, &openMin, &closeHour, &closeMin)
	if err != nil {
		return false, errors.New("invalid WorkingHours format")
	}

	// Crear objetos de tiempo para los límites de horario
	openTime := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), openHour, openMin, 0, 0, dateTime.Location())
	closeTime := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), closeHour, closeMin, 0, 0, dateTime.Location())

	// Verificar si la hora actual está dentro del rango
	if dateTime.Before(openTime) || dateTime.After(closeTime) {
		return false, nil
	}

	return true, nil
}

// FindByAssistantID retrieves all number phones associated with a specific assistant
func (r *AssistantRepository) FindByAssistantID(assistantID int64) ([]entities.NumberPhone, error) {
	var records []entities.NumberPhone
	err := r.db.Where("assistants_id = ?", assistantID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r *AssistantRepository) FindAll() ([]entities.Assistant, error) {
	var assistants []entities.Assistant
	err := r.db.Preload("Bussiness").Find(&assistants).Error
	return assistants, err
}

func (r *AssistantRepository) FindById(id int64) (entities.Assistant, error) {
	var assistant entities.Assistant
	err := r.db.Preload("Bussiness").First(&assistant, id).Error
	return assistant, err
}

func (r *AssistantRepository) Update(id int64, data entities.Assistant) error {
	return r.db.Model(&data).Where("id = ?", id).Updates(data).Error
}

func (r *AssistantRepository) Delete(id int64) error {
	return r.db.Delete(&entities.Assistant{}, id).Error
}

func (r *AssistantRepository) GetAllAssistantsByBussinessId(businessId int64) ([]entities.Assistant, error) {
	var assistants []entities.Assistant
	err := r.db.Where("bussiness_id = ?", businessId).
		Preload("Bussiness").
		Find(&assistants).Error
	return assistants, err
}
