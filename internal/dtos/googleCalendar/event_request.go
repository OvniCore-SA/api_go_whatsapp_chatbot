package googlecalendar

import (
	"errors"
	"time"

	"google.golang.org/api/calendar/v3"
)

type EventRequest struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Start       string `json:"start"` // Fecha y hora en formato RFC3339
	End         string `json:"end"`   // Fecha y hora en formato RFC3339
	ContactsID  uint   `json:"contacts_id"`
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

	startTime, endTime, err := ValidateDateEventGoogleCalendar(e.Start, e.End)
	if err != nil {
		return err
	}

	e.Start = startTime.Format("2006-01-02T15:04:05")
	e.End = endTime.Format("2006-01-02T15:04:05")

	if e.ContactsID <= 0 {
		return errors.New("el campo 'contacts_id' no puede estar vacío")
	}

	return nil
}

func ValidateDateEventGoogleCalendar(startTimeString, endTimeString string) (startTime, endTime time.Time, err error) {
	// Definir posibles formatos para aceptar RFC3339 con y sin Z
	formats := []string{ // "2006-01-02T15:04:05Z07:00" (con Z o offset)
		"2006-01-02T15:04:05",       // Sin Z ni offset
		"2006-01-02T15:04:05-07:00", // Con offset explícito
	}

	// Función para intentar parsear en distintos formatos
	parseTime := func(value string) (time.Time, error) {
		var err error
		var t time.Time
		for _, format := range formats {
			t, err = time.Parse(format, value)
			if err == nil {
				return t, nil
			}
		}
		return t, err
	}

	// Validar Start
	startTime, err = parseTime(startTimeString)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("el campo 'start' debe estar en formato RFC3339 sin 'Z' (2006-01-02T15:04:05) o con offset (2006-01-02T15:04:05-07:00)")
	}

	// Validar End
	endTime, err = parseTime(endTimeString)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("el campo 'end' debe estar en formato RFC3339 sin 'Z' (2006-01-02T15:04:05) o con offset (2006-01-02T15:04:05-07:00)")
	}

	// Validar que End sea después de Start
	if !endTime.After(startTime) {
		return time.Time{}, time.Time{}, errors.New("el campo 'end' debe ser posterior al campo 'start'")
	}

	return
}

// Carga un request para enviar a la creacion de un evento en calendar.
func UploadEventRequestCalendar(summary, description, startDate, endDate string) (calendarRequest *calendar.Event, err error) {

	startDateTime, endDateTime, err := ValidateDateEventGoogleCalendar(startDate, endDate)
	if err != nil {
		return
	}

	startDate = startDateTime.Format("2006-01-02T15:04:05")
	endDate = endDateTime.Format("2006-01-02T15:04:05")

	return &calendar.Event{
		Summary:     summary,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: startDate,
			TimeZone: "UTC-3",
		},
		End: &calendar.EventDateTime{
			DateTime: endDate,
			TimeZone: "UTC-3",
		},
	}, nil
}

// Carga un request para enviar a la creacion de un evento en calendar.
func UploadEventRequestEditCalendar(summary, description, startDate, endDate string) (eventChatbotRequest *EventRequest, err error) {

	startDateTime, endDateTime, err := ValidateDateEventGoogleCalendar(startDate, endDate)
	if err != nil {
		return
	}

	startDate = startDateTime.Format("2006-01-02T15:04:05")
	endDate = endDateTime.Format("2006-01-02T15:04:05")

	return &EventRequest{
		Summary:     summary,
		Description: description,
		Start:       startDate,
		End:         endDate,
	}, nil
}
