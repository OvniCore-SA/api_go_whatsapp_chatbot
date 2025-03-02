package mysql_client

import (
	"fmt"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities/filters"
	"gorm.io/gorm"
)

// NumberPhonesRepository is the repository for NumberPhone entities
type NumberPhonesRepository struct {
	db *gorm.DB
}

// NewNumberPhonesRepository creates a new instance of NumberPhonesRepository
func NewNumberPhonesRepository(db *gorm.DB) *NumberPhonesRepository {
	return &NumberPhonesRepository{db: db}
}

// Create inserts a new number phone record into the database
func (r *NumberPhonesRepository) Create(record entities.NumberPhone) error {
	return r.db.Create(&record).Error
}

// FindByID retrieves a number phone record by its ID
func (r *NumberPhonesRepository) FindByID(id string) (entities.NumberPhone, error) {
	var record entities.NumberPhone
	err := r.db.First(&record, id).Error
	return record, err
}

// Update modifies an existing number phone record
func (r *NumberPhonesRepository) Update(id string, record entities.NumberPhone) error {
	return r.db.Model(&record).Where("id = ?", id).Updates(record).Error
}

// Delete removes a number phone record from the database
func (r *NumberPhonesRepository) Delete(id string) error {
	return r.db.Delete(&entities.NumberPhone{}, id).Error
}

// List retrieves all number phone records
func (r *NumberPhonesRepository) List() ([]entities.NumberPhone, error) {
	var records []entities.NumberPhone
	err := r.db.Find(&records).Error
	return records, err
}

// GetNumberPhonesByAssistantID retrieves all number phones associated with a specific assistant
func (r *NumberPhonesRepository) GetNumberPhonesByAssistantID(assistantID int64) ([]entities.NumberPhone, error) {
	var numberPhones []entities.NumberPhone
	err := r.db.Where("assistants_id = ?", assistantID).Find(&numberPhones).Error
	if err != nil {
		return nil, fmt.Errorf("error retrieving number phones for assistant ID %d: %w", assistantID, err)
	}
	return numberPhones, nil
}

// FindByAssistantID retrieves all number phones associated with a specific assistant
func (r *NumberPhonesRepository) FindByAssistantID(assistantID int64) ([]entities.NumberPhone, error) {
	var records []entities.NumberPhone
	err := r.db.Where("assistants_id = ?", assistantID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r *NumberPhonesRepository) ListByFilter(filter filters.AssistantsFiltro) ([]entities.NumberPhone, error) {
	var records []entities.NumberPhone

	// Base de la consulta
	query := r.db.Where("number_phone_to_notify > ?", 0)

	// Aplicar Preload para relaciones si es necesario
	if filter.UpladContacts {
		query = query.Preload("Contacts") // Precarga la relación "Contacts"
	}

	query = query.Preload("Assistant")

	// Ejecutar la consulta
	err := query.Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("error retrieving number phones with contacts: %w", err)
	}
	return records, err
}

func (r *NumberPhonesRepository) UUIDExists(uuid string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.NumberPhone{}).Where("uuid = ?", uuid).Count(&count).Error
	if err != nil {
		return false, err
	}

	// Si count > 0, significa que ya existe un número de teléfono con ese UUID
	return count > 0, nil
}
