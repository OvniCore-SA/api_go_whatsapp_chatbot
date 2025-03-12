package dtos

import (
	"errors"
	"strings"
)

type AssistantDto struct {
	ID                 int64            `json:"id"`
	BussinessID        int64            `json:"bussiness_id"`                   // Id de la empresa a la que pertenece
	Name               string           `json:"name"`                           // Openai
	OpenaiAssistantsID string           `json:"openai_assistants_id,omitempty"` // Openai
	Description        string           `json:"description,omitempty"`          // Openai
	Model              string           `json:"model,omitempty"`                // Openai
	Instructions       string           `json:"instructions,omitempty"`         // Openai
	Active             bool             `json:"active"`                         // Indica si el assistant está activo
	EventDuration      int64            `json:"event_duration"`                 // Duracion de cada evento que se programa ej: 15, 30, 60. (esto se cuenta como minutos)
	AccountGoogle      bool             `json:"account_google"`                 // Indica si el assistant tiene asociada credenciales para registrar eventos en google calendar
	EventType          string           `json:"event_type"`                     // Typo de evento que se guarda en google calendar ej: turno, reunion
	EventCountPerDay   int16            `json:"event_count_per_day"`            // Cantidad de eventos por dia que puede tener una persona
	Bussiness          BussinessDto     `json:"bussiness,omitempty"`            //
	NumberPhones       []NumberPhoneDto `json:"number_phones,omitempty"`        //
	Events             []EventsDto      `json:"events,omitempty"`               //
	OpeningDays        uint8            `json:"opening_days"`                   // Días de apertura representados en un entero de 7 bits
	WorkingHours       string           `json:"working_hours"`                  // Horarios de trabajo en formato "HH:MM-HH:MM,HH:MM-HH:MM"
}

func (dto *AssistantDto) ValidateAssistantDto(isCreate bool) error {
	// Validaciones comunes para creación y edición
	if dto.BussinessID <= 0 {
		return errors.New("bussiness_id es obligatorio y debe ser mayor que 0")
	}

	if strings.TrimSpace(dto.Name) == "" {
		return errors.New("el nombre es obligatorio")
	}

	if len(dto.Name) > 100 {
		return errors.New("el nombre no debe exceder los 100 caracteres")
	}

	if dto.Model != "" && len(dto.Model) > 50 {
		return errors.New("el modelo no debe exceder los 50 caracteres")
	}

	if dto.Instructions != "" && len(dto.Instructions) > 10000 {
		return errors.New("las instrucciones no deben exceder los 10000 caracteres")
	}

	if len(dto.Instructions) < 10 {
		return errors.New("la instrucción debe tener por lo menos 10 caracteres")
	}

	if len(dto.Description) < 10 {
		return errors.New("la descripción debe tener por lo menos 10 caracteres")
	}

	if strings.TrimSpace(dto.OpenaiAssistantsID) == "" {
		return errors.New("openai_assistants_id es obligatorio")
	}

	if len(dto.Model) < 2 {
		return errors.New("el modelo debe tener por lo menos 2 caracteres")
	}

	// Validaciones específicas para edición
	if !isCreate {
		if dto.ID <= 0 {
			return errors.New("el ID del asistente es obligatorio al editar")
		}
	}

	// Validar GoogleCalendarConfig si está presente
	// if dto.GoogleCalendarConfig != nil {
	// 	if dto.GoogleCalendarConfig.GoogleUserID == "" {
	// 		return errors.New("google_user_id es obligatorio en google_calendar_credential")
	// 	}

	// 	if dto.GoogleCalendarConfig.AccessToken == "" {
	// 		return errors.New("access_token es obligatorio en google_calendar_credential")
	// 	}

	// 	if dto.GoogleCalendarConfig.RefreshToken == "" {
	// 		return errors.New("refresh_token es obligatorio en google_calendar_credential")
	// 	}

	// 	if dto.GoogleCalendarConfig.TokenExpiry.IsZero() {
	// 		return errors.New("token_expiry es obligatorio en google_calendar_credential")
	// 	}
	// }

	return nil
}
