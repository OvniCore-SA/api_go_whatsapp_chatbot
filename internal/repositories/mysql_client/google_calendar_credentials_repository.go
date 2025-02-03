package mysql_client

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"gorm.io/gorm"
)

type GoogleCalendarCredentialsRepository struct {
	db *gorm.DB
}

func NewGoogleCalendarConfigsRepository(db *gorm.DB) *GoogleCalendarCredentialsRepository {
	return &GoogleCalendarCredentialsRepository{db: db}
}

// Create saves Google Calendar credentials into the database
func (r *GoogleCalendarCredentialsRepository) Create(data *entities.GoogleCalendarCredential) error {
	if err := r.db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

// FindByAssistantID retrieves credentials for a specific assistant
func (r *GoogleCalendarCredentialsRepository) FindByAssistantID(assistantID int) (*entities.GoogleCalendarCredential, error) {
	var credential entities.GoogleCalendarCredential
	err := r.db.Where("assistants_id = ?", assistantID).First(&credential).Error
	if err != nil {
		return nil, err
	}
	return &credential, nil
}

// Update updates the credentials for a specific assistant
func (r *GoogleCalendarCredentialsRepository) Update(data *entities.GoogleCalendarCredential) error {
	if err := r.db.Save(data).Error; err != nil {
		return err
	}
	return nil
}

// Delete deletes the credentials for a specific assistant
func (r *GoogleCalendarCredentialsRepository) Delete(assistantID int) error {
	err := r.db.Where("assistants_id = ?", assistantID).Delete(&entities.GoogleCalendarCredential{}).Error
	if err != nil {
		return err
	}
	return nil
}
