package mysql_client

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"gorm.io/gorm"
)

// Interfaz para definir métodos del repositorio
type EventsRepository interface {
	Create(event *entities.Events) error
	FindByID(id int) (*entities.Events, error)
	FindAll() ([]entities.Events, error)
	Update(event *entities.Events) error
	Delete(id int) error
}

// Implementación del repositorio
type eventsRepositoryImpl struct {
	db *gorm.DB
}

func NewEventsRepository(db *gorm.DB) EventsRepository {
	return &eventsRepositoryImpl{db: db}
}

func (r *eventsRepositoryImpl) Create(event *entities.Events) error {
	return r.db.Create(event).Error
}

func (r *eventsRepositoryImpl) FindByID(id int) (*entities.Events, error) {
	var event entities.Events
	err := r.db.Preload("GoogleCalendarConfig").Preload("Contact").First(&event, id).Error
	return &event, err
}

func (r *eventsRepositoryImpl) FindAll() ([]entities.Events, error) {
	var events []entities.Events
	err := r.db.Preload("GoogleCalendarConfig").Preload("Contact").Find(&events).Error
	return events, err
}

func (r *eventsRepositoryImpl) Update(event *entities.Events) error {
	return r.db.Save(event).Error
}

func (r *eventsRepositoryImpl) Delete(id int) error {
	return r.db.Delete(&entities.Events{}, id).Error
}
