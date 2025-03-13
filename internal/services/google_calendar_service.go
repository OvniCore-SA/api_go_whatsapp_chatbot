package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	googlecalendar "github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/googleCalendar"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

type GoogleCalendarService struct {
	repository       *mysql_client.GoogleCalendarCredentialsRepository
	AssistantService AssistantService
	EventsService    EventsService
}

func NewGoogleCalendarService(repository *mysql_client.GoogleCalendarCredentialsRepository, assistantService AssistantService, eventsService EventsService) *GoogleCalendarService {
	return &GoogleCalendarService{
		repository:       repository,
		AssistantService: assistantService,
		EventsService:    eventsService,
	}
}

// GetCredentials obtiene las credenciales de un asistente
func (s *GoogleCalendarService) GetCredentials(assistantID int) (*entities.GoogleCalendarCredential, error) {
	return s.repository.FindByAssistantID(assistantID)
}

// SaveCredentials guarda las credenciales en la base de datos
func (s *GoogleCalendarService) SaveCredentials(assistantID int, token *oauth2.Token, googleUserID, googleUserEmail string) error {
	// Validar si ya existen credenciales para este asistente
	existingCredential, err := s.repository.FindByAssistantID(assistantID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if existingCredential != nil {
		// Actualizar las credenciales existentes
		existingCredential.AccessToken = token.AccessToken
		existingCredential.RefreshToken = token.RefreshToken
		existingCredential.TokenExpiry = token.Expiry
		existingCredential.GoogleUserID = googleUserID
		if err := s.repository.Update(existingCredential); err != nil {
			return err
		}

		// Actualizo el campo google calendar al assistant
		if err := s.AssistantService.repository.Update(int64(assistantID), entities.Assistant{AccountGoogle: true}); err != nil {
			return err
		}

		return nil
	}

	// Crear nuevas credenciales
	newCredential := &entities.GoogleCalendarCredential{
		AssistantsID: assistantID,
		GoogleUserID: googleUserID,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenExpiry:  token.Expiry,
	}

	if err := s.repository.Create(newCredential); err != nil {
		return err
	}

	// Actualizo el campo google calendar al assistant
	if err := s.AssistantService.repository.Update(int64(assistantID), entities.Assistant{AccountGoogle: true}); err != nil {
		return err
	}

	return nil
}

// DeleteCredentials removes Google Calendar credentials for an assistant
func (s *GoogleCalendarService) DeleteCredentials(assistantID int) error {
	return s.repository.Delete(assistantID)
}

func (s *GoogleCalendarService) CreateGoogleCalendarEvent(token *oauth2.Token, ctx context.Context, event *calendar.Event) (*calendar.Event, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	// Registrar el evento en el calendario principal
	createdEvent, err := srv.Events.
		Insert("primary", event).
		SendNotifications(true).
		SendUpdates("all").
		ConferenceDataVersion(1).
		Do()

	if err != nil {
		return nil, err
	}

	return createdEvent, nil
}

// DeleteGoogleCalendarEvent elimina un evento del Google Calendar
func (s *GoogleCalendarService) DeleteGoogleCalendarEvent(token *oauth2.Token, ctx context.Context, eventID string) error {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	// Eliminar el evento del calendario principal
	err = srv.Events.Delete("primary", eventID).Do()
	if err != nil {
		return err
	}

	return nil
}

// UpdateGoogleCalendarEvent actualiza un evento en Google Calendar
func (s *GoogleCalendarService) UpdateGoogleCalendarEvent(token *oauth2.Token, ctx context.Context, eventID string, eventRequest *googlecalendar.EventRequest) (*calendar.Event, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	// Crear un objeto solo con los campos que deseas actualizar
	event := &calendar.Event{
		Summary:     eventRequest.Summary,
		Description: eventRequest.Description,
		Start: &calendar.EventDateTime{
			DateTime: eventRequest.Start,
			TimeZone: "UTC",
		},
		End: &calendar.EventDateTime{
			DateTime: eventRequest.End,
			TimeZone: "UTC",
		},
	}

	// Si quieres asegurar que se mantenga el evento de Meet existente (sin eliminarlo), no hace falta pasar ConferenceData.
	// Pero si deseas asegurarte de que el evento siempre tenga conferencia de Meet, inclúyelo así:
	event.ConferenceData = &calendar.ConferenceData{
		CreateRequest: &calendar.CreateConferenceRequest{
			RequestId: fmt.Sprintf("meet-%s", eventID), // siempre único
			ConferenceSolutionKey: &calendar.ConferenceSolutionKey{
				Type: "hangoutsMeet",
			},
		},
	}

	// Aplicar cambios usando Patch en vez de Update
	updatedEvent, err := srv.Events.Patch("primary", eventID, event).
		ConferenceDataVersion(1).
		SendUpdates("all").
		Do()

	if err != nil {
		return nil, err
	}

	return updatedEvent, nil
}

// GetGoogleUserID obtiene el ID único del usuario de Google
func GetGoogleUserID(client *http.Client, token *oauth2.Token) (string, string, error) {
	// Crear una solicitud HTTP con el token de autenticación
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo?alt=json", nil)
	if err != nil {
		return "", "", err
	}

	// Agregar el token de autenticación en el encabezado Authorization
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	// Ejecutar la solicitud
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// Verificar el código de estado de la respuesta
	if resp.StatusCode != http.StatusOK {
		log.Printf("\n\n respStatusCode: %d", resp.StatusCode)

		req.Header.Set("Authorization", "Bearer "+token.RefreshToken)
		log.Printf("\nRefreshToken: %s", token.RefreshToken)
		// Ejecutar la solicitud
		resp, err = client.Do(req)
		if err != nil {
			log.Println(resp.StatusCode)
			return "", "", err
		}
		log.Printf("Segunda consulta: %d", resp.StatusCode)

		return "", "", errors.New("error al obtener la información del usuario de Google")
	}

	// Decodificar la respuesta JSON
	var result struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	// Retornar el ID del usuario
	return result.ID, result.Email, nil
}

func (s *GoogleCalendarService) FetchGoogleCalendarEventsByDate(token *oauth2.Token, ctx context.Context, startDate, endDate time.Time) (*calendar.Events, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return srv.Events.List("primary").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(startDate.Format(time.RFC3339)).
		TimeMax(endDate.Format(time.RFC3339)).
		OrderBy("startTime").
		Do()
}

func (s *GoogleCalendarService) InsertGoogleCalendarEvent(token *oauth2.Token, ctx context.Context, event *calendar.Event) error {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	_, err = srv.Events.Insert("primary", event).Do()
	return err
}

func (s *GoogleCalendarService) FetchGoogleCalendarEvents(token *oauth2.Token, ctx context.Context) (*calendar.Events, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	now := time.Now().Format(time.RFC3339)
	return srv.Events.List("primary").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(now).
		MaxResults(10).
		OrderBy("startTime").
		Do()
}

func (s *GoogleCalendarService) GetOrRefreshToken(assistantID int, config *oauth2.Config, ctx context.Context) (*oauth2.Token, error) {
	credentials, err := s.GetCredentials(assistantID)
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{
		AccessToken:  credentials.AccessToken,
		RefreshToken: credentials.RefreshToken,
		Expiry:       credentials.TokenExpiry,
	}

	if token.Expiry.Before(time.Now()) {
		tokenSource := config.TokenSource(ctx, token)
		newToken, err := tokenSource.Token()
		if err != nil {
			return nil, err
		}

		err = s.SaveCredentials(assistantID, newToken, credentials.GoogleUserID, credentials.Email)
		if err != nil {
			return nil, err
		}

		return newToken, nil
	}

	return token, nil
}
