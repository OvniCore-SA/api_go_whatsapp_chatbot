package mysql_client

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"gorm.io/gorm"
)

// ConfigurationsRepository handles operations related to configurations
type ConfigurationsRepository struct {
	db *gorm.DB
}

// NewConfigurationsRepository creates a new instance of ConfigurationsRepository
func NewConfigurationsRepository(db *gorm.DB) *ConfigurationsRepository {
	return &ConfigurationsRepository{db: db}
}

// Create inserts a new configuration record into the database
func (r *ConfigurationsRepository) Create(record entities.Configuration) error {
	return r.db.Create(&record).Error
}

// FindByID retrieves a configuration record by its ID
func (r *ConfigurationsRepository) FindByID(id int64) (entities.Configuration, error) {
	var record entities.Configuration
	err := r.db.First(&record, id).Error
	return record, err
}

// FindByKey retrieves a configuration record by its key name
func (r *ConfigurationsRepository) FindByKey(keyName string) (entities.Configuration, error) {
	var record entities.Configuration
	err := r.db.Where("key_name = ?", keyName).First(&record).Error
	return record, err
}

// List retrieves all configuration records
func (r *ConfigurationsRepository) List() ([]entities.Configuration, error) {
	var records []entities.Configuration
	err := r.db.Find(&records).Error
	return records, err
}

// Update updates an existing configuration record in the database
func (r *ConfigurationsRepository) Update(record entities.Configuration) error {
	return r.db.Save(&record).Error
}

// Delete removes a configuration record from the database by ID
func (r *ConfigurationsRepository) Delete(id int64) error {
	return r.db.Delete(&entities.Configuration{}, id).Error
}
