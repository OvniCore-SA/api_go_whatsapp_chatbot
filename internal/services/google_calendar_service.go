package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type GoogleCalendarService struct {
	repository *mysql_client.GoogleCalendarCredentialsRepository
}

func NewGoogleCalendarService(repository *mysql_client.GoogleCalendarCredentialsRepository) *GoogleCalendarService {
	return &GoogleCalendarService{
		repository: repository,
	}
}

// GetCredentials obtiene las credenciales de un asistente
func (s *GoogleCalendarService) GetCredentials(assistantID int) (*entities.GoogleCalendarCredential, error) {
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
	newCredential := &entities.GoogleCalendarCredential{
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

// GetGoogleUserID obtiene el ID único del usuario de Google
func GetGoogleUserID(client *http.Client) (string, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("error al obtener la información del usuario de Google")
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.ID, nil
}
