package postgres_client

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"gorm.io/gorm"
)

// LogsRepository is the repository for Logs entities
type LogsRepository struct {
	db *gorm.DB
}

// NewLogsRepository creates a new instance of LogsRepository
func NewLogsRepository(db *gorm.DB) *LogsRepository {
	return &LogsRepository{db: db}
}

// Create inserts a new record into the database
func (r *LogsRepository) Create(record entities.Logs) error {
	return r.db.Create(&record).Error
}

// FindByID retrieves a record by its ID
func (r *LogsRepository) FindByID(id string) (entities.Logs, error) {
	var record entities.Logs
	err := r.db.First(&record, id).Error
	return record, err
}

// Update modifies an existing record
func (r *LogsRepository) Update(id string, record entities.Logs) error {
	return r.db.Model(&record).Where("id = ?", id).Updates(record).Error
}

// Delete removes a record from the database
func (r *LogsRepository) Delete(id string) error {
	return r.db.Delete(&entities.Logs{}, id).Error
}

// List retrieves all records
func (r *LogsRepository) List() ([]entities.Logs, error) {
	var records []entities.Logs
	err := r.db.Find(&records).Error
	return records, err
}
