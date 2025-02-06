package config

import (
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// LoadOAuthConfig carga la configuración de OAuth desde el archivo JSON.
func LoadOAuthConfig() *oauth2.Config {
	credentials, err := os.ReadFile(os.Getenv("PATH_CREDENTIAL_GOOGLE"))
	if err != nil {
		fmt.Printf("No se pudo leer el archivo de credenciales: %v", err)
		//log.Fatalf("No se pudo leer el archivo de credenciales: %v", err)
	}

	fmt.Println("Archivo de autenticación GOOGLE cargado exitosamente.")

	config, err := google.ConfigFromJSON(credentials, calendar.CalendarEventsScope, calendar.CalendarReadonlyScope, calendar.CalendarScope, "https://www.googleapis.com/auth/calendar.events.owned", "https://www.googleapis.com/auth/calendar.app.created", "https://www.googleapis.com/auth/userinfo.profile")
	if err != nil {
		fmt.Printf("No se pudo parsear el archivo de credenciales: %v", err)
		//log.Fatalf("No se pudo parsear el archivo de credenciales: %v", err)
	}
	return config
}
