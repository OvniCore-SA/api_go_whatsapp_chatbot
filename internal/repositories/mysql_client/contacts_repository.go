package mysql_client

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"gorm.io/gorm"
)

// ContactsRepository is the repository for Contacts entities
type ContactsRepository struct {
	db *gorm.DB
}

// NewContactsRepository creates a new instance of ContactsRepository
func NewContactsRepository(db *gorm.DB) *ContactsRepository {
	return &ContactsRepository{db: db}
}

// Create inserts a new contact record into the database
func (r *ContactsRepository) Create(record entities.Contact) error {
	return r.db.Create(&record).Error
}

// FindByID retrieves a contact record by its ID
func (r *ContactsRepository) FindByID(id string) (entities.Contact, error) {
	var record entities.Contact
	err := r.db.First(&record, id).Error
	return record, err
}

// Update modifies an existing contact record
func (r *ContactsRepository) Update(id string, record entities.Contact) error {
	return r.db.Model(&record).Where("id = ?", id).Updates(record).Error
}

// Delete removes a contact record from the database
func (r *ContactsRepository) Delete(id string) error {
	return r.db.Delete(&entities.Contact{}, id).Error
}

// List retrieves all contact records
func (r *ContactsRepository) List() ([]entities.Contact, error) {
	var records []entities.Contact
	err := r.db.Find(&records).Error
	return records, err
}
