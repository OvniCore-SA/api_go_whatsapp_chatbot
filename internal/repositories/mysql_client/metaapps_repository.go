package mysql_client

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"gorm.io/gorm"
)

// MetaAppsRepository is the repository for MetaApps entities
type MetaAppsRepository struct {
	db *gorm.DB
}

// NewMetaAppsRepository creates a new instance of MetaAppsRepository
func NewMetaAppsRepository(db *gorm.DB) *MetaAppsRepository {
	return &MetaAppsRepository{db: db}
}

// Create inserts a new record into the database
func (r *MetaAppsRepository) Create(record entities.MetaApps) error {
	return r.db.Create(&record).Error
}

// FindByID retrieves a record by its ID
func (r *MetaAppsRepository) FindByID(id int64) (entities.MetaApps, error) {
	var record entities.MetaApps
	err := r.db.Preload("Promps").First(&record, id).Error
	return record, err
}

// Update modifies an existing record
func (r *MetaAppsRepository) Update(id string, record entities.MetaApps) error {
	return r.db.Model(&record).Where("id = ?", id).Updates(record).Error
}

// Delete removes a record from the database
func (r *MetaAppsRepository) Delete(id string) error {
	return r.db.Delete(&entities.MetaApps{}, id).Error
}

// List retrieves all records
func (r *MetaAppsRepository) List() ([]entities.MetaApps, error) {
	var records []entities.MetaApps
	err := r.db.Find(&records).Error
	return records, err
}
