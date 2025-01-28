package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

type GoogleCalendarService struct {
	repository *mysql_client.GoogleCalendarConfigsRepository
}

func NewGoogleCalendarService(repository *mysql_client.GoogleCalendarConfigsRepository) *GoogleCalendarService {
	return &GoogleCalendarService{
		repository: repository,
	}
}

// GetCredentials obtiene las credenciales de un asistente
func (s *GoogleCalendarService) GetCredentials(assistantID int) (*entities.GoogleCalendarConfig, error) {
	return s.repository.FindByAssistantID(assistantID)
}

// SaveCredentials guarda las credenciales en la base de datos
func (s *GoogleCalendarService) SaveCredentials(assistantID int, token *oauth2.Token, googleUserID string) error {
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
		return s.repository.Update(existingCredential)
	}

	// Crear nuevas credenciales
	newCredential := &entities.GoogleCalendarConfig{
		AssistantsID: assistantID,
		GoogleUserID: googleUserID,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenExpiry:  token.Expiry,
	}
	return s.repository.Create(newCredential)
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
	createdEvent, err := srv.Events.Insert("primary", event).Do()
	if err != nil {
		return nil, err
	}

	return createdEvent, nil
}

// GetGoogleUserID obtiene el ID único del usuario de Google
func GetGoogleUserID(client *http.Client, token *oauth2.Token) (string, error) {
	// Crear una solicitud HTTP con el token de autenticación
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo?alt=json", nil)
	if err != nil {
		return "", err
	}

	// Agregar el token de autenticación en el encabezado Authorization
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	// Ejecutar la solicitud
	resp, err := client.Do(req)
	if err != nil {
		return "", err
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
			return "", err
		}
		log.Println("Segunda consulta: %d", resp.StatusCode)

		return "", errors.New("error al obtener la información del usuario de Google")
	}

	// Decodificar la respuesta JSON
	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Retornar el ID del usuario
	return result.ID, nil
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

		err = s.SaveCredentials(assistantID, newToken, credentials.GoogleUserID)
		if err != nil {
			return nil, err
		}

		return newToken, nil
	}

	return token, nil
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
