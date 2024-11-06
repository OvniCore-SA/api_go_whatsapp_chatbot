package mysql_client

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"gorm.io/gorm"
)

// PermissionsRepository is the repository for Permissions entities
type PermissionsRepository struct {
	db *gorm.DB
}

// NewPermissionsRepository creates a new instance of PermissionsRepository
func NewPermissionsRepository(db *gorm.DB) *PermissionsRepository {
	return &PermissionsRepository{db: db}
}

// Create inserts a new record into the database
func (r *PermissionsRepository) Create(record entities.Permissions) error {
	return r.db.Create(&record).Error
}

// FindByID retrieves a record by its ID
func (r *PermissionsRepository) FindByID(id string) (entities.Permissions, error) {
	var record entities.Permissions
	err := r.db.First(&record, id).Error
	return record, err
}

// Update modifies an existing record
func (r *PermissionsRepository) Update(id string, record entities.Permissions) error {
	return r.db.Model(&record).Where("id = ?", id).Updates(record).Error
}

// Delete removes a record from the database
func (r *PermissionsRepository) Delete(id string) error {
	return r.db.Delete(&entities.Permissions{}, id).Error
}

// List retrieves all records
func (r *PermissionsRepository) List() ([]entities.Permissions, error) {
	var records []entities.Permissions
	err := r.db.Find(&records).Error
	return records, err
}
