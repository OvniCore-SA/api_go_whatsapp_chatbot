package mysql_client

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"gorm.io/gorm"
)

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(file entities.File) error {
	return r.db.Create(&file).Error
}

func (r *FileRepository) FindAll() ([]entities.File, error) {
	var files []entities.File
	err := r.db.Find(&files).Error
	return files, err
}

func (r *FileRepository) FindById(id int64) (entities.File, error) {
	var file entities.File
	err := r.db.First(&file, id).Error
	return file, err
}

func (r *FileRepository) FindByAssistantID(id int64) ([]entities.File, error) {
	var file []entities.File
	err := r.db.Where("assistants_id = ?", id).Find(&file).Error
	return file, err
}

func (r *FileRepository) Update(file entities.File) error {
	return r.db.Save(&file).Error
}

func (r *FileRepository) Delete(id int64) error {
	return r.db.Delete(&entities.File{}, id).Error
}
