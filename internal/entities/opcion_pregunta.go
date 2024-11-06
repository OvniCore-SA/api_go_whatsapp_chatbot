package entities

import (
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"gorm.io/gorm"
)

type OpcionPreguntas struct {
	ID                   int64 `gorm:"primaryKey"` // ID como clave primaria
	Nombre               string
	OpcionPregunta       string
	Activo               bool
	UltimaOpcion         bool
	ChatbotsID           int64             // Clave foránea que referencia a Chatbots
	Chatbot              Chatbots          `gorm:"foreignKey:ChatbotsID"` // Relación con Chatbots
	OpcionPreguntasId    *uint             // Relación recursiva con la misma tabla, representa el ID del padre
	ChildOpcionPreguntas []OpcionPreguntas `gorm:"foreignKey:OpcionPreguntasId"` // Relación uno a muchos con las subpreguntas
	CreatedAt            time.Time         // Fecha de creación
	UpdatedAt            time.Time         // Fecha de actualización
	DeletedAt            gorm.DeletedAt    `gorm:"index"` // Fecha de eliminación suave (soft delete)
}

func MapEntitiesToOpcionPreguntaDto(opcionPregunta OpcionPreguntas) dtos.OpcionPreguntasDto {
	// Mapeo de los campos básicos
	dto := dtos.OpcionPreguntasDto{
		ID:             opcionPregunta.ID,
		Nombre:         opcionPregunta.Nombre,
		OpcionPregunta: opcionPregunta.OpcionPregunta,
		Activo:         opcionPregunta.Activo,
		UltimaOpcion:   opcionPregunta.UltimaOpcion,
		ChatbotsID:     opcionPregunta.ChatbotsID,
		CreatedAt:      opcionPregunta.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      opcionPregunta.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Mapeo de las preguntas hijas
	for _, child := range opcionPregunta.ChildOpcionPreguntas {
		dto.ChildOpcionPreguntas = append(dto.ChildOpcionPreguntas, MapEntitiesToOpcionPreguntaDto(child))
	}

	return dto
}

func MapDtoToOpcionPreguntas(dto dtos.OpcionPreguntasDto) OpcionPreguntas {
	// Parseo de las fechas desde el string
	createdAt, _ := time.Parse("2006-01-02 15:04:05", dto.CreatedAt)
	updatedAt, _ := time.Parse("2006-01-02 15:04:05", dto.UpdatedAt)

	// Mapeo de los campos básicos
	opcionPregunta := OpcionPreguntas{
		ID:                dto.ID,
		Nombre:            dto.Nombre,
		OpcionPregunta:    dto.OpcionPregunta,
		OpcionPreguntasId: dto.OpcionPreguntaID,
		Activo:            dto.Activo,
		UltimaOpcion:      dto.UltimaOpcion,
		ChatbotsID:        dto.ChatbotsID,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}

	// Mapeo de las preguntas hijas
	for _, childDto := range dto.ChildOpcionPreguntas {
		opcionPregunta.ChildOpcionPreguntas = append(opcionPregunta.ChildOpcionPreguntas, MapDtoToOpcionPreguntas(childDto))
	}

	return opcionPregunta
}
