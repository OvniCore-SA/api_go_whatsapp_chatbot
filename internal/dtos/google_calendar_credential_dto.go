package dtos

import "time"

type GoogleCalendarConfigDto struct {
	ID           int       `json:"id"`
	AssistantsID int       `json:"assistants_id"`
	GoogleUserID string    `json:"google_user_id"` // ID único del usuario en Google
	AccessToken  string    `json:"access_token"`   // Token de acceso para la API de Google Calendar
	RefreshToken string    `json:"refresh_token"`  // Token de refresco para regenerar el access_token
	TokenExpiry  time.Time `json:"token_expiry"`   // Expiración del token de acceso
}
