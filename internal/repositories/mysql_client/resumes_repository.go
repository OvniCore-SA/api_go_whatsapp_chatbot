package mysql_client

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"gorm.io/gorm"
)

type ResumesRepository struct {
	db *gorm.DB
}

func NewResumesRepository(db *gorm.DB) *ResumesRepository {
	return &ResumesRepository{db: db}
}

// GetResumeByChatbotID busca un Resume basado en el ChatbotID
func (r *ResumesRepository) GetResumeByChatbotID(chatbotID int64) (*entities.Resumes, error) {
	var resume entities.Resumes
	err := r.db.Where("chatbots_id = ?", chatbotID).Preload("Chatbot").First(&resume).Error
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

func (r *ResumesRepository) List() ([]entities.Resumes, error) {
	var resumes []entities.Resumes
	err := r.db.Preload("Chatbot").Find(&resumes).Error
	return resumes, err
}

func (r *ResumesRepository) GetByID(id int64) (entities.Resumes, error) {
	var resume entities.Resumes
	err := r.db.Preload("Chatbot").First(&resume, id).Error
	return resume, err
}

func (r *ResumesRepository) Create(record entities.Resumes) error {
	return r.db.Create(&record).Error
}

func (r *ResumesRepository) Update(record entities.Resumes) error {
	return r.db.Save(&record).Error
}

func (r *ResumesRepository) Delete(id int64) error {
	return r.db.Delete(&entities.Resumes{}, id).Error
}
