package filters

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type EventsFilter struct {
	ID                    int    `json:"id"`
	StartDate             string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate               string `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
	EventGoogleCalendarID string `json:"event_google_calendar_id" validate:"omitempty"`
	AssistantsID          *int64 `json:"assistants_id" validate:"omitempty,gt=0"`
	ContactsID            *int64 `json:"contacts_id" validate:"omitempty,gt=0"`
	CodeEvent             string `json:"code_event" validate:"omitempty"`
	CreatedAt             string `json:"created_at" validate:"omitempty,datetime=2006-01-02"`
	MonthYear             string `json:"month_year" validate:"omitempty,len=7,datetime=2006-01"`
}

// Validador de datos del DTO
var validate = validator.New()

// Validate valida la estructura de EventsDto
func (e *EventsFilter) Validate() error {
	err := validate.Struct(e)
	if err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	if e.StartDate != "" && e.EndDate != "" {
		start, err := time.Parse("2006-01-02", e.StartDate)
		if err != nil {
			return fmt.Errorf("invalid start date format, use YYYY-MM-DD: %v", err)
		}

		end, err := time.Parse("2006-01-02", e.EndDate)
		if err != nil {
			return fmt.Errorf("invalid end date format, use YYYY-MM-DD: %v", err)
		}

		if start.After(end) {
			return fmt.Errorf("start date must be before end date")
		}
	}

	return nil
}
