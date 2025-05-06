package postgres_client

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"gorm.io/gorm"
)

// UsersRepository is the repository for Users entities
type UsersRepository struct {
	db *gorm.DB
}

// NewUsersRepository creates a new instance of UsersRepository
func NewUsersRepository(db *gorm.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

// Create inserts a new record into the database
func (r *UsersRepository) Create(record entities.Users) error {
	return r.db.Create(&record).Error
}

// FindByID retrieves a record by its ID
func (r *UsersRepository) FindByID(id uint) (*entities.Users, error) {
	var user entities.Users
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update modifies an existing record
func (r *UsersRepository) Update(id int64, record entities.Users) error {
	return r.db.Model(&record).Where("id = ?", id).Updates(record).Error
}

// Delete removes a record from the database
func (r *UsersRepository) Delete(id int64) error {
	return r.db.Delete(&entities.Users{}, id).Error
}

// List retrieves all records
func (r *UsersRepository) List() ([]entities.Users, error) {
	var records []entities.Users
	res := r.db.Model(entities.Users{})
	res = res.Preload("Rol")
	res = res.Preload("PasswordResets")
	err := res.Find(&records).Error
	return records, err
}

// FindByEmail retrieves a user record by email
func (r *UsersRepository) FindByEmail(email string) (entities.Users, error) {
	var record entities.Users
	err := r.db.Where("email = ?", email).Preload("Rol.Permissions").First(&record).Error
	return record, err
}
