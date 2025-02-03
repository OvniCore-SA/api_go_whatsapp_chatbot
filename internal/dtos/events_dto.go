package dtos

import (
	"fmt"

	"google.golang.org/api/calendar/v3"
)

type EventsDto struct {
	ID                    int    `json:"id"`
	Summary               string `json:"summary"`
	Description           string `json:"description"`
	StartDate             string `json:"start_date"`
	EndDate               string `json:"end_date"`
	EventGoogleCalendarID string `json:"event_google_calendar_id"`
	AssistantsID          int64  `json:"assistants_id"`
	ContactsID            int64  `json:"contacts_id"`
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
