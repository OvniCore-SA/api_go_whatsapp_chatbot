package postgres_client

import (
	"fmt"
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities/filters"
	"gorm.io/gorm"
)

// Interfaz para definir métodos del repositorio
type EventsRepository interface {
	Create(event *entities.Events) error
	FindByID(id int) (*entities.Events, error)
	FindAll(request *filters.EventsFilter, pagination *dtos.Pagination) (events []entities.Events, total int64, err error)
	Update(event *entities.Events) error
	Delete(id int) error
	Cancel(codeEvent string) error

	FindByContactAndDateAndTime(contactID int64, date string, currentTime string) ([]entities.Events, error)
	ExistsByCode(code string) (bool, error)
	FindByContactAndCodeEvent(contactID int64, codeEvent string) (entities.Events, error)
	FindByContactDateAndNumberPhone(contactID int64, date string, assistantID int64) ([]entities.Events, error)
}

// Implementación del repositorio
type eventsRepositoryImpl struct {
	db *gorm.DB
}

func NewEventsRepository(db *gorm.DB) EventsRepository {
	return &eventsRepositoryImpl{db: db}
}

func (r *eventsRepositoryImpl) FindByContactDateAndNumberPhone(contactID int64, date string, assistantID int64) ([]entities.Events, error) {
	var events []entities.Events

	// Realizamos la consulta filtrando por contacts_id, fecha y number_phones_id
	err := r.db.
		Where("contacts_id = ? AND DATE(start_date) = ? AND assistants_id = ?", contactID, date, assistantID).
		Find(&events).Error

	if err != nil {
		return nil, fmt.Errorf("error finding events by contact, date, and number_phone_id: %v", err)
	}
	return events, nil
}

func (r *eventsRepositoryImpl) FindByContactAndCodeEvent(contactID int64, codeEvent string) (entities.Events, error) {
	var event entities.Events

	// Realizamos la consulta para obtener un evento por contactID y code_event
	err := r.db.
		Where("contacts_id = ? AND code_event = ?", contactID, codeEvent).
		First(&event).Error // Usamos First() porque esperamos solo un evento
	if err != nil {
		return entities.Events{}, fmt.Errorf("error finding event by code_event: %v", err)
	}
	return event, nil
}

func (r *eventsRepositoryImpl) FindByContactAndDateAndTime(contactID int64, date string, currentTime string) ([]entities.Events, error) {
	var events []entities.Events

	// La función espera:
	// - 'date' en formato "YYYY-MM-DD" (solo la parte de la fecha) para comparar con DATE(start_date).
	// - 'currentTime' en formato "YYYY-MM-DD HH:MM:SS", que se usará para filtrar los eventos a partir de esa hora.
	// Sin embargo, currentTime puede venir en distintos formatos. Por ello, se intenta parsearlo usando varios formatos:
	formats := []string{
		time.RFC3339,          // Ejemplo: "2006-01-02T15:04:05Z07:00"
		"2006-01-02 15:04:05", // Ejemplo: "2025-02-03 22:43:00"
		"2006-01-02",          // Ejemplo: "2025-02-03"
	}

	var ct time.Time
	var err error
	parsed := false
	for _, format := range formats {
		ct, err = time.Parse(format, currentTime)
		if err == nil {
			parsed = true
			break
		}
	}
	if !parsed {
		return nil, fmt.Errorf("error parsing currentTime: %v", err)
	}

	// Convertir la hora parseada al formato DATETIME deseado: "YYYY-MM-DD HH:MM:SS"
	formattedCurrentTime := ct.Format("2006-01-02 15:04:05")

	// Realizar la consulta: se buscan los eventos donde:
	// - contacts_id coincide con el parámetro contactID.
	// - La fecha de start_date (sin la hora) coincide con 'date' (formato "YYYY-MM-DD").
	// - start_date es mayor o igual que formattedCurrentTime.
	err = r.db.
		Where("contacts_id = ? AND DATE(start_date) = ? AND start_date >= ?", contactID, date, formattedCurrentTime).
		Order("start_date ASC").
		Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (r *eventsRepositoryImpl) ExistsByCode(code string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.Events{}).Where("code_event = ?", code).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *eventsRepositoryImpl) Create(event *entities.Events) error {
	return r.db.Create(event).Error
}

func (r *eventsRepositoryImpl) FindByID(id int) (*entities.Events, error) {
	var event entities.Events
	err := r.db.Preload("Contact").First(&event, id).Error
	return &event, err
}

func (r *eventsRepositoryImpl) FindAll(request *filters.EventsFilter, pagination *dtos.Pagination) (events []entities.Events, total int64, err error) {
	query := r.db.Model(&entities.Events{}).Preload("Contact")

	if request.AssistantsID != nil {
		if *request.AssistantsID != 0 {
			query = query.Where("assistants_id = ?", request.AssistantsID)
		}

	}
	if request.MonthYear != "" {
		query = query.Where("DATE_FORMAT(start_date, '%Y-%m') = ?", request.MonthYear)
	}
	if request.StartDate != "" {
		query = query.Where("DATE(start_date) = ?", request.StartDate)
	}
	// Obtener total sin paginación
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación si corresponde
	if pagination.Number > 0 && pagination.Size > 0 {
		offset := int((pagination.Number - 1) * pagination.Size)
		query = query.Offset(offset).Limit(int(pagination.Size))
	}

	if err := query.Order("start_date DESC").Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (r *eventsRepositoryImpl) Update(event *entities.Events) error {
	return r.db.Save(event).Error
}

func (r *eventsRepositoryImpl) Delete(id int) error {
	return r.db.Delete(&entities.Events{}, id).Error
}

func (r *eventsRepositoryImpl) Cancel(codeEvent string) error {
	var event entities.Events
	// Primero obtenemos el primer registro que coincida con el código
	err := r.db.Where("code_event = ?", codeEvent).First(&event).Error
	if err != nil {
		return fmt.Errorf("no se pudo eliminar el evento con el código '%s': %v", codeEvent, err)
	}
	// Luego eliminamos el registro encontrado
	err = r.db.Delete(&event).Error
	if err != nil {
		return fmt.Errorf("error cancelando evento: %v", err)
	}
	return nil
}
