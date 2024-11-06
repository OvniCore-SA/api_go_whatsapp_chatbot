package mysql_client

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"gorm.io/gorm"
)

// ChatbotsRepository is the repository for Chatbots entities
type ChatbotsRepository struct {
	db *gorm.DB
}

// NewChatbotsRepository creates a new instance of ChatbotsRepository
func NewChatbotsRepository(db *gorm.DB) *ChatbotsRepository {
	return &ChatbotsRepository{db: db}
}

// Create inserts a new record into the database
func (r *ChatbotsRepository) Create(record entities.Chatbots) error {
	return r.db.Create(&record).Error
}

// Buscar Chatbot por PhoneNumberID
func (r *ChatbotsRepository) FindByPhoneNumberID(phoneNumberID string) (*entities.Chatbots, error) {
	var chatbot entities.Chatbots
	if err := r.db.Where("phone_number_id = ?", phoneNumberID).Preload("MetaApps").First(&chatbot).Error; err != nil {
		return nil, err
	}
	return &chatbot, nil
}

// FindByID retrieves a record by its ID
func (r *ChatbotsRepository) FindByID(id int64) (entities.Chatbots, error) {
	var record entities.Chatbots
	err := r.db.First(&record, id).Error
	return record, err
}

// Update modifies an existing record
func (r *ChatbotsRepository) Update(id string, record entities.Chatbots) error {
	return r.db.Model(&record).Where("id = ?", id).Updates(record).Error
}

// Delete removes a record from the database
func (r *ChatbotsRepository) Delete(id string) error {
	return r.db.Delete(&entities.Chatbots{}, id).Error
}

// List retrieves all records
func (r *ChatbotsRepository) List() ([]entities.Chatbots, error) {
	var records []entities.Chatbots
	err := r.db.Preload("MetaApps").Find(&records).Error
	return records, err
}
