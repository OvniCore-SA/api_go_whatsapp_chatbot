package postgres_client

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"gorm.io/gorm"
)

// RolesRepository is the repository for Roles entities
type RolesRepository struct {
	db *gorm.DB
}

// NewRolesRepository creates a new instance of RolesRepository
func NewRolesRepository(db *gorm.DB) *RolesRepository {
	return &RolesRepository{db: db}
}

// Create inserts a new record into the database
func (r *RolesRepository) Create(record entities.Roles) error {
	return r.db.Create(&record).Error
}

// FindByID retrieves a record by its ID
func (r *RolesRepository) FindByID(id string) (entities.Roles, error) {
	var record entities.Roles
	err := r.db.First(&record, id).Error
	return record, err
}

// Update modifies an existing record
func (r *RolesRepository) Update(id string, record entities.Roles) error {
	return r.db.Model(&record).Where("id = ?", id).Updates(record).Error
}

func (r *RolesRepository) GetByRol(rol string) (entities.Roles, error) {
	var record entities.Roles
	err := r.db.Model(&record).Where("rol = ?", rol).Find(&record).Error
	return record, err
}

// Delete removes a record from the database
func (r *RolesRepository) Delete(id string) error {
	return r.db.Delete(&entities.Roles{}, id).Error
}

// List retrieves all records
func (r *RolesRepository) List() ([]entities.Roles, error) {
	var records []entities.Roles
	err := r.db.Find(&records).Error
	return records, err
}
