package mysql_client

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities/filters"
	"gorm.io/gorm"
)

// OpcionesPreguntasRepository is the repository for OpcionPreguntas entities
type OpcionesPreguntasRepository struct {
	db *gorm.DB
}

// NewOpcionesPreguntasRepository creates a new instance of OpcionesPreguntasRepository
func NewOpcionesPreguntasRepository(db *gorm.DB) *OpcionesPreguntasRepository {
	return &OpcionesPreguntasRepository{db: db}
}

// Create inserts a new record into the database
func (r *OpcionesPreguntasRepository) Create(record entities.OpcionPreguntas) error {
	return r.db.Create(&record).Error
}

// FindByID retrieves a record by its ID
func (r *OpcionesPreguntasRepository) FindByID(id uint64) (entities.OpcionPreguntas, error) {
	var record entities.OpcionPreguntas
	err := r.db.First(&record, id).Error
	return record, err
}

// FindByID retrieves a record by its ID
func (r *OpcionesPreguntasRepository) ListByIDOpcionPreguntaRepository(filter filters.OpcionPreguntasFiltro) ([]entities.OpcionPreguntas, error) {
	var record []entities.OpcionPreguntas
	var err error

	if !filter.PrimerMenu {

		err = r.db.Model(&entities.OpcionPreguntas{}).
			Where("opcion_preguntas_id = ?", filter.OpcionPreguntaID).
			Where("chatbots_id = ?", filter.ChatbotsID).
			Preload("Chatbot").
			Find(&record).Error

	} else {
		// Encontramos el padre cuya opcion_preguntas_id es NULL
		var parent entities.OpcionPreguntas
		err = r.db.Model(&entities.OpcionPreguntas{}).
			Select("id").
			Where("opcion_preguntas_id IS NULL").
			Where("chatbots_id = ?", filter.ChatbotsID).
			First(&parent).Error

		if err != nil {
			return record, err
		}

		// Ahora encontramos los hijos del padre
		err = r.db.Model(&entities.OpcionPreguntas{}).
			Where("opcion_preguntas_id = ?", parent.ID).
			Preload("ChildOpcionPreguntas"). // Preload de los hijos si es necesario
			Find(&record).Error
	}

	return record, err
}

// Update modifies an existing record
func (r *OpcionesPreguntasRepository) Update(id string, record entities.OpcionPreguntas) error {
	return r.db.Model(&record).Where("id = ?", id).Updates(record).Error
}

// Return OpcionPregunta ultima_opcion = true
func (r *OpcionesPreguntasRepository) GetUltimaOpcionByID(id string, record entities.OpcionPreguntas) error {
	return r.db.Model(&entities.OpcionPreguntas{}).Where("id = ?", id).First(record).Error
}

// Delete removes a record from the database
func (r *OpcionesPreguntasRepository) Delete(id string) error {
	return r.db.Delete(&entities.OpcionPreguntas{}, id).Error
}

// List retrieves all records
func (r *OpcionesPreguntasRepository) List() ([]entities.OpcionPreguntas, error) {
	var records []entities.OpcionPreguntas
	err := r.db.Find(&records).Error
	return records, err
}
