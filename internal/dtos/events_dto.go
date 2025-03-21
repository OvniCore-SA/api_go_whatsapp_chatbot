package dtos

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"google.golang.org/api/calendar/v3"
)

// EventsDto representa la estructura de un evento con validaciones
type EventsDto struct {
	ID                    int    `json:"id"`
	Summary               string `json:"summary" validate:"required,min=3,max=255"`
	Description           string `json:"description" validate:"required,min=5,max=500"`
	StartDate             string `json:"start_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndDate               string `json:"end_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00,gtfield=StartDate"`
	EventGoogleCalendarID string `json:"event_google_calendar_id" validate:"omitempty"`
	AssistantsID          int64  `json:"assistants_id" validate:"required,gt=0"`
	ContactsID            int64  `json:"contacts_id" validate:"required,gt=0"`
	CodeEvent             string `json:"code_event" validate:"omitempty"`
	CreatedAt             string `json:"created_at"`
	MonthYear             string `json:"month_year" validate:"required,len=7,datetime=2006-01"`
}

// Validador de datos del DTO
var validate = validator.New()

// Validate valida la estructura de EventsDto
func (e *EventsDto) Validate() error {
	err := validate.Struct(e)
	if err != nil {
		return fmt.Errorf("validation error: %v", err)
	}
	return nil
}

// MapCalendarEventToEventsDto convierte un evento de Google Calendar a EventsDto
func MapCalendarEventToEventsDto(event *calendar.Event) (EventsDto, error) {
	if event == nil {
		return EventsDto{}, fmt.Errorf("event don't field in MapCalendarEventToEventsDto")
	}

	return EventsDto{
		Summary:               event.Summary,
		Description:           event.Description,
		StartDate:             event.Start.DateTime,
		EndDate:               event.End.DateTime,
		EventGoogleCalendarID: event.Id,
	}, nil
}
