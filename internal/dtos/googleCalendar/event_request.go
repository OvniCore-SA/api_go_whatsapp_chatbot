package googlecalendar

import (
	"errors"
	"time"
)

type EventRequest struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Start       string `json:"start"` // Fecha y hora en formato RFC3339
	End         string `json:"end"`   // Fecha y hora en formato RFC3339
}

// Validate verifica que el evento tenga valores válidos
func (e *EventRequest) Validate() error {
	// Validar que los campos Summary y Description no estén vacíos
	if e.Summary == "" {
		return errors.New("el campo 'summary' no puede estar vacío")
	}
	if e.Description == "" {
		return errors.New("el campo 'description' no puede estar vacío")
	}

	// Validar el formato de Start y End
	startTime, err := time.Parse(time.RFC3339, e.Start)
	if err != nil {
		return errors.New("el campo 'start' debe estar en formato RFC3339.(2006-01-02T15:04:05Z07:00)")
	}

	endTime, err := time.Parse(time.RFC3339, e.End)
	if err != nil {
		return errors.New("el campo 'end' debe estar en formato RFC3339.(2006-01-02T15:04:05Z07:00)")
	}

	// Validar que Start sea anterior a End
	if !startTime.Before(endTime) {
		return errors.New("el campo 'start' debe ser anterior a 'end'")
	}

	return nil
}
