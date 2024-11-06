package mysql_client

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"gorm.io/gorm"
)

// BussinessRepository is the repository for Bussiness entities
type BussinessRepository struct {
	db *gorm.DB
}

// NewBussinessRepository creates a new instance of BussinessRepository
func NewBussinessRepository(db *gorm.DB) *BussinessRepository {
	return &BussinessRepository{db: db}
}

// Create inserts a new bussiness record into the database
func (r *BussinessRepository) Create(record entities.Bussines) error {
	return r.db.Create(&record).Error
}

// FindByID retrieves a bussiness record by its ID
func (r *BussinessRepository) FindByID(id int64) (entities.Bussines, error) {
	var record entities.Bussines
	err := r.db.First(&record, id).Error
	return record, err
}

// Update modifies an existing bussiness record
func (r *BussinessRepository) Update(id int64, record entities.Bussines) error {
	return r.db.Model(&record).Where("id = ?", id).Updates(record).Error
}

// Delete removes a bussiness record from the database
func (r *BussinessRepository) Delete(id int64) error {
	return r.db.Delete(&entities.Bussines{}, id).Error
}

// List retrieves all bussiness records
func (r *BussinessRepository) List() ([]entities.Bussines, error) {
	var records []entities.Bussines
	err := r.db.Find(&records).Error
	return records, err
}

// FindByUserId retrieves all bussiness records associated with a specific user ID
func (r *BussinessRepository) FindByUserId(userId int64) ([]entities.Bussines, error) {
	var records []entities.Bussines
	err := r.db.Where("users_id = ?", userId).
		Preload("Assistants"). // Carga anticipada de asistentes, si es necesario
		Find(&records).Error

	return records, err
}
