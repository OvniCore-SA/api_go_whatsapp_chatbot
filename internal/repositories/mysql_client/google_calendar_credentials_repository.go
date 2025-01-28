package mysql_client

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"gorm.io/gorm"
)

type GoogleCalendarConfigsRepository struct {
	db *gorm.DB
}

func NewGoogleCalendarConfigsRepository(db *gorm.DB) *GoogleCalendarConfigsRepository {
	return &GoogleCalendarConfigsRepository{db: db}
}

// Create saves Google Calendar credentials into the database
func (r *GoogleCalendarConfigsRepository) Create(data *entities.GoogleCalendarConfig) error {
	if err := r.db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

// FindByAssistantID retrieves credentials for a specific assistant
func (r *GoogleCalendarConfigsRepository) FindByAssistantID(assistantID int) (*entities.GoogleCalendarConfig, error) {
	var credential entities.GoogleCalendarConfig
	err := r.db.Where("assistants_id = ?", assistantID).First(&credential).Error
	if err != nil {
		return nil, err
	}
	return &credential, nil
}

// Update updates the credentials for a specific assistant
func (r *GoogleCalendarConfigsRepository) Update(data *entities.GoogleCalendarConfig) error {
	if err := r.db.Save(data).Error; err != nil {
		return err
	}
	return nil
}

// Delete deletes the credentials for a specific assistant
func (r *GoogleCalendarConfigsRepository) Delete(assistantID int) error {
	err := r.db.Where("assistants_id = ?", assistantID).Delete(&entities.GoogleCalendarConfig{}).Error
	if err != nil {
		return err
	}
	return nil
}
