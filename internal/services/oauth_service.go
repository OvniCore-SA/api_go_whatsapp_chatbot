package services

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"golang.org/x/oauth2"
)

// ExchangeCodeForToken intercambia el código de autorización por un token.
func ExchangeCodeForToken(config *oauth2.Config, code string) (*oauth2.Token, error) {
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Error al intercambiar el código por un token: %v", err)
		return nil, err
	}
	// Log para depurar el contenido del token
	log.Printf("Token recibido: %+v\n", token)
	return token, nil
}

// SaveToken guarda un token en un archivo local.
func SaveToken(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

// LoadToken carga un token desde un archivo local.
func LoadToken(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}
