package mysql_client

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"gorm.io/gorm"
)

type AssistantRepository struct {
	db *gorm.DB
}

func NewAssistantRepository(db *gorm.DB) *AssistantRepository {
	return &AssistantRepository{db: db}
}

func (r *AssistantRepository) Create(data *entities.Assistant) error {
	// GORM automáticamente asigna el ID a data.ID después de la creación
	if err := r.db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func (r *AssistantRepository) FindAll() ([]entities.Assistant, error) {
	var assistants []entities.Assistant
	err := r.db.Preload("Bussiness").Find(&assistants).Error
	return assistants, err
}

func (r *AssistantRepository) FindById(id int64) (entities.Assistant, error) {
	var assistant entities.Assistant
	err := r.db.Preload("Bussiness").First(&assistant, id).Error
	return assistant, err
}

func (r *AssistantRepository) Update(id int64, data entities.Assistant) error {
	return r.db.Model(&data).Where("id = ?", id).Updates(data).Error
}

func (r *AssistantRepository) Delete(id int64) error {
	return r.db.Delete(&entities.Assistant{}, id).Error
}

func (r *AssistantRepository) GetAllAssistantsByBussinessId(businessId int64) ([]entities.Assistant, error) {
	var assistants []entities.Assistant
	err := r.db.Where("bussiness_id = ?", businessId).
		Preload("Bussiness").
		Find(&assistants).Error
	return assistants, err
}
