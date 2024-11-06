package mysql_client

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities/filters"
	"gorm.io/gorm"
)

// PrompsRepository is the repository for Promps entities
type PrompsRepository struct {
	db *gorm.DB
}

// NewPrompsRepository creates a new instance of PrompsRepository
func NewPrompsRepository(db *gorm.DB) *PrompsRepository {
	return &PrompsRepository{db: db}
}

// GetByFilter busca un registro de Promps en base a los filtros proporcionados
func (r *PrompsRepository) GetByFilter(filter filters.PrompsFiltro) (entities.Promps, error) {
	var record entities.Promps
	query := r.db.Model(&entities.Promps{})

	if filter.MetaAppsId != 0 {
		query = query.Where("meta_apps_id = ?", filter.MetaAppsId)
	}

	if filter.Activo != nil {
		query = query.Where("activo = ?", *filter.Activo)
	}

	err := query.First(&record).Error
	if err != nil {
		return entities.Promps{}, err
	}

	return record, nil
}

// Create inserts a new record into the database
func (r *PrompsRepository) Create(record entities.Promps) error {
	return r.db.Create(&record).Error
}

// FindByID retrieves a record by its ID
func (r *PrompsRepository) FindByID(id int64) (entities.Promps, error) {
	var record entities.Promps
	err := r.db.First(&record, id).Error
	return record, err
}

// Update modifies an existing record
func (r *PrompsRepository) Update(id string, record entities.Promps) error {
	return r.db.Model(&record).Where("id = ?", id).Updates(record).Error
}

// Delete removes a record from the database
func (r *PrompsRepository) Delete(id string) error {
	return r.db.Delete(&entities.Promps{}, id).Error
}

// List retrieves all records
func (r *PrompsRepository) List() ([]entities.Promps, error) {
	var records []entities.Promps
	err := r.db.Find(&records).Error
	return records, err
}
