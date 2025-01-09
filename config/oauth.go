package config

import (
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// LoadOAuthConfig carga la configuración de OAuth desde el archivo JSON.
func LoadOAuthConfig() *oauth2.Config {
	credentials, err := os.ReadFile(os.Getenv("PATH_CREDENTIAL_GOOGLE"))
	if err != nil {
		log.Fatalf("No se pudo leer el archivo de credenciales: %v", err)
	}

	config, err := google.ConfigFromJSON(credentials, "https://www.googleapis.com/auth/calendar.readonly")
	if err != nil {
		log.Fatalf("No se pudo parsear el archivo de credenciales: %v", err)
	}
	return config
}
