package mysql_client

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"gorm.io/gorm"
)

// MessagesRepository handles operations related to messages
type MessagesRepository struct {
	db *gorm.DB
}

// NewMessagesRepository creates a new instance of MessagesRepository
func NewMessagesRepository(db *gorm.DB) *MessagesRepository {
	return &MessagesRepository{db: db}
}

// Create inserts a new message record into the database
func (r *MessagesRepository) Create(record entities.Message) error {
	return r.db.Create(&record).Error
}

// Obtener todos los mensajes entre un assistant y un contact
func (r *MessagesRepository) GetMessagesByAssistantAndContact(assistantID, contactID int64) ([]entities.Message, error) {
	var messages []entities.Message
	err := r.db.Where("assistant_id = ? AND contact_id = ?", assistantID, contactID).Order("created_at ASC").Find(&messages).Error
	return messages, err
}

func (r *MessagesRepository) GetConversation(assistantID, contactID int64, sinceMinutes int) ([]entities.Message, error) {
	var messages []entities.Message
	query := r.db.Where("assistants_id = ? AND contacts_id = ?", assistantID, contactID).Order("created_at ASC")

	if sinceMinutes > 0 {
		threshold := time.Now().Add(-time.Duration(sinceMinutes) * time.Minute)
		query = query.Where("created_at >= ?", threshold)
	}

	err := query.Find(&messages).Error
	return messages, err
}
