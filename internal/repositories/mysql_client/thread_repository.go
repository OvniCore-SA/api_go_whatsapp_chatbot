package mysql_client

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"gorm.io/gorm"
)

// ThreadRepository define las operaciones con la base de datos
type ThreadRepository struct {
	db *gorm.DB
}

// NewThreadRepository crea una nueva instancia del repositorio
func NewThreadRepository(db *gorm.DB) *ThreadRepository {
	return &ThreadRepository{db: db}
}

// Crear un nuevo hilo en la base de datos
func (r *ThreadRepository) Create(thread entities.Thread) (*entities.Thread, error) {
	if err := r.db.Create(&thread).Error; err != nil {
		return nil, err
	}
	return &thread, nil
}

// Obtener el último Thread activo asociado a un contacto
func (r *ThreadRepository) FindLastActiveByContactID(contactID int64) (*entities.Thread, error) {
	var thread entities.Thread
	err := r.db.Where("contacts_id = ? AND deleted_at IS NULL", contactID).
		Order("created_at DESC").
		First(&thread).Error

	if err != nil {
		return nil, err
	}
	return &thread, nil
}

// Buscar hilo por ID
func (r *ThreadRepository) FindByID(id int64) (*entities.Thread, error) {
	var thread entities.Thread
	if err := r.db.First(&thread, id).Error; err != nil {
		return nil, err
	}
	return &thread, nil
}

// Buscar hilo por ThreadsId
func (r *ThreadRepository) FindByThreadsId(threadsId string) (*entities.Thread, error) {
	var thread entities.Thread
	if err := r.db.Where("threads_id = ?", threadsId).First(&thread).Error; err != nil {
		return nil, err
	}
	return &thread, nil
}

// Obtener todos los hilos
func (r *ThreadRepository) GetAll() ([]entities.Thread, error) {
	var threads []entities.Thread
	if err := r.db.Find(&threads).Error; err != nil {
		return nil, err
	}
	return threads, nil
}

// Actualizar un hilo
func (r *ThreadRepository) Update(id int64, updatedThread entities.Thread) error {
	return r.db.Model(&entities.Thread{}).Where("id = ?", id).Updates(updatedThread).Error
}

// Eliminar un hilo (Soft Delete)
func (r *ThreadRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&entities.Thread{}).Error
}
