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

// Create inserts a new bussiness record into the database and returns the ID
func (r *BussinessRepository) Create(record entities.Bussines) (uint, error) {
	if err := r.db.Create(&record).Error; err != nil {
		return 0, err // Si ocurre un error, devolver el valor cero y el error
	}

	return uint(record.ID), nil // Devolver el ID del registro creado
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
	err := r.db.Preload("Users").Find(&records).Error
	return records, err
}

// FindByUserId retrieves all businesses associated with a specific user ID
func (r *BussinessRepository) FindByUserId(userId int64) ([]entities.Bussines, error) {
	var businesses []entities.Bussines

	err := r.db.Joins("JOIN bussiness_has_users ON bussiness_has_users.bussiness_id = bussiness.id").
		Where("bussiness_has_users.users_id = ?", userId).
		Find(&businesses).Error

	if err != nil {
		return nil, err
	}

	return businesses, nil
}

// AddUserToBussiness associates a user with a business in the many-to-many table
func (r *BussinessRepository) AddUserToBussiness(businessID, userID int64) error {
	association := entities.BussinessHasUsers{
		BussinessID: businessID,
		UsersID:     userID,
	}
	return r.db.Create(&association).Error
}

// RemoveUserFromBussiness removes the association between a user and a business
func (r *BussinessRepository) RemoveUserFromBussiness(businessID, userID int64) error {
	return r.db.Where("bussiness_id = ? AND users_id = ?", businessID, userID).
		Delete(&entities.BussinessHasUsers{}).Error
}
