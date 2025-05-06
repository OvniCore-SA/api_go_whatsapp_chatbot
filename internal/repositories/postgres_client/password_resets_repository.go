package postgres_client

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"gorm.io/gorm"
)

// PasswordResetsRepository is the repository for PasswordResets entities
type PasswordResetsRepository struct {
	db *gorm.DB
}

// NewPasswordResetsRepository creates a new instance of PasswordResetsRepository
func NewPasswordResetsRepository(db *gorm.DB) *PasswordResetsRepository {
	return &PasswordResetsRepository{db: db}
}

// Create inserts a new record into the database
func (r *PasswordResetsRepository) Create(record entities.PasswordResets) error {
	return r.db.Create(&record).Error
}

// FindByID retrieves a record by its ID
func (r *PasswordResetsRepository) FindByID(id string) (entities.PasswordResets, error) {
	var record entities.PasswordResets
	err := r.db.First(&record, id).Error
	return record, err
}

// Update modifies an existing record
func (r *PasswordResetsRepository) Update(id string, record entities.PasswordResets) error {
	return r.db.Model(&record).Where("id = ?", id).Updates(record).Error
}

// Delete removes a record from the database
func (r *PasswordResetsRepository) Delete(id string) error {
	return r.db.Delete(&entities.PasswordResets{}, id).Error
}

// List retrieves all records
func (r *PasswordResetsRepository) List() ([]entities.PasswordResets, error) {
	var records []entities.PasswordResets
	err := r.db.Find(&records).Error
	return records, err
}

// FindByToken retrieves a password reset record by its token
func (r *PasswordResetsRepository) FindByToken(token string) (entities.PasswordResets, error) {
	var record entities.PasswordResets
	err := r.db.Where("token = ?", token).First(&record).Error
	return record, err
}

// DeleteByToken removes a password reset record from the database by its token
func (r *PasswordResetsRepository) DeleteByToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&entities.PasswordResets{}).Error
}
