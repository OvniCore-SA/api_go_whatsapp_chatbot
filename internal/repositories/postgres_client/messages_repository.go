package postgres_client

import (
	"fmt"
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

// Verifyca si existe un registro con el messageID.
func (r *MessagesRepository) ExistsByMessageID(messageID string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.Message{}).Where("message_id_whatsapp = ?", messageID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Obtener todos los mensajes entre un assistant y un contact
func (r *MessagesRepository) GetMessagesByAssistantAndContact(assistantID, contactID int64) ([]entities.Message, error) {
	var messages []entities.Message
	err := r.db.Where("assistant_id = ? AND contact_id = ?", assistantID, contactID).Order("created_at ASC").Find(&messages).Error
	return messages, err
}

// GetMessagesByNumber retrieves all messages associated with a specific number within a given time range
func (r *MessagesRepository) GetMessagesByNumber(numberID, contacID int64, since time.Time) ([]entities.Message, error) {
	var messages []entities.Message
	err := r.db.Where("number_phones_id = ? AND contacts_id = ? AND  created_at >= ?", numberID, contacID, since).Order("created_at ASC").Preload("Contact").Find(&messages).Error

	if err != nil {
		return nil, err
	}
	return messages, nil
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

func (r *MessagesRepository) GetMessagesWithContacts(numberIDs []int64, since time.Time) ([]entities.Message, error) {
	var messages []entities.Message

	err := r.db.
		Where("number_phones_id IN ? AND created_at >= ?", numberIDs, since).
		Preload("NumberPhone").
		Preload("Contact").
		Order("created_at DESC").
		Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching messages with contacts: %w", err)
	}

	return messages, nil
}

// DoesNumberPhoneExist - Verifica si un número de teléfono existe en la base de datos
func (r *MessagesRepository) DoesNumberPhoneExist(numberPhoneID int64) (bool, error) {
	var count int64

	err := r.db.Model(&entities.NumberPhone{}).
		Where("id = ?", numberPhoneID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetMessagesByNumberPhone - Obtiene los mensajes asociados a un número de teléfono específico con paginación
func (r *MessagesRepository) GetMessagesByNumberPhone(numberPhoneID int64, page int, limit int) ([]entities.Message, int, error) {
	var messages []entities.Message
	var total int64

	// Contar el total de registros antes de aplicar paginación
	err := r.db.Model(&entities.Message{}).
		Joins("JOIN contacts ON contacts.id = messages.contacts_id").
		Where("messages.number_phones_id = ? AND contacts.deleted_at IS NULL", numberPhoneID).
		Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	// Evitar valores inválidos en paginación
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Obtener los registros paginados
	err = r.db.
		Joins("JOIN contacts ON contacts.id = messages.contacts_id").
		Where("messages.number_phones_id = ? AND contacts.deleted_at IS NULL", numberPhoneID).
		Order("messages.created_at DESC").
		Limit(limit).
		Offset(offset).
		Preload("Contact").
		Find(&messages).Error

	if err != nil {
		return nil, 0, err
	}

	return messages, int(total), nil
}

// GetMessagesByNumberPhoneAndContact - Obtiene los mensajes asociados a un número de teléfono y contacto específico con paginación
func (r *MessagesRepository) GetMessagesByNumberPhoneAndContact(numberPhoneID int64, contactID int64, page int, limit int) ([]entities.Message, int, error) {
	var messages []entities.Message
	var total int64

	// Contar el total de registros antes de aplicar paginación
	err := r.db.Model(&entities.Message{}).
		Joins("JOIN contacts ON contacts.id = messages.contacts_id").
		Where("messages.number_phones_id = ? AND messages.contacts_id = ? AND contacts.deleted_at IS NULL", numberPhoneID, contactID).
		Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	// Evitar valores inválidos en paginación
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Obtener los registros paginados
	err = r.db.
		Joins("JOIN contacts ON contacts.id = messages.contacts_id").
		Where("messages.number_phones_id = ? AND messages.contacts_id = ? AND contacts.deleted_at IS NULL", numberPhoneID, contactID).
		Order("messages.created_at DESC").
		Limit(limit).
		Offset(offset).
		Preload("Contact").
		Find(&messages).Error

	if err != nil {
		return nil, 0, err
	}

	return messages, int(total), nil
}
