package controllers

import (
	"log"
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// GetCalendarEvents obtiene los eventos del calendario.
func GetCalendarEvents(config *oauth2.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Cargar el token desde el archivo
		token, err := services.LoadToken("token.json")
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No estás autenticado",
			})
		}

		// Crear cliente HTTP autenticado
		client := config.Client(c.Context(), token)

		// Crear servicio de Google Calendar
		srv, err := calendar.NewService(c.Context(), option.WithHTTPClient(client))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "No se pudo crear el servicio de Google Calendar",
			})
		}

		// Obtener eventos
		now := time.Now().Format(time.RFC3339)
		events, err := srv.Events.List("primary").
			ShowDeleted(false).
			SingleEvents(true).
			TimeMin(now).
			MaxResults(10).
			OrderBy("startTime").
			Do()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "No se pudieron obtener los eventos",
			})
		}

		if len(events.Items) == 0 {
			return c.JSON(fiber.Map{
				"events": "No se encontraron próximos eventos.",
			})
		}

		// Formatear la respuesta
		var result []fiber.Map
		for _, item := range events.Items {
			dateTime := item.Start.DateTime
			if dateTime == "" {
				dateTime = item.Start.Date
			}
			result = append(result, fiber.Map{
				"summary": item.Summary,
				"start":   dateTime,
			})
		}

		return c.JSON(fiber.Map{
			"events": result,
		})
	}
}

// GetAuthURL genera la URL de autenticación de Google.
func GetAuthURL(config *oauth2.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		return c.JSON(fiber.Map{
			"auth_url": authURL,
		})
	}
}

// HandleAuthCallback maneja el callback de Google después de la autorización.
func HandleAuthCallback(config *oauth2.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Obtener el código de la URL
		code := c.Query("code")
		if code == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "El parámetro 'code' es obligatorio",
			})
		}

		// Intercambiar el código por un token
		token, err := services.ExchangeCodeForToken(config, code)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "No se pudo obtener el token",
			})
		}

		// Guardar el token en un archivo (o en memoria/cache si prefieres)
		if err := services.SaveToken("token.json", token); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "No se pudo guardar el token",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Autenticación exitosa",
			"token":   token.AccessToken,
			"expiry":  token.Expiry,
		})
	}
}

// GetAuthToken procesa el código de autorización y genera el token de acceso.
func GetAuthToken(config *oauth2.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Imprime la URL completa de la solicitud
		log.Println("URL completa de la solicitud:")
		log.Println(c.OriginalURL())

		// Imprime el cuerpo de la solicitud
		log.Println("Body recibido:")
		log.Println(string(c.Body()))

		// Captura el código de autorización desde los parámetros de consulta
		code := c.Query("code")
		if code == "" {
			log.Println("code: ", code)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "El parámetro 'code' es obligatorio",
			})
		}

		// Imprime las cabeceras de la solicitud
		log.Println("Headers:")
		c.Request().Header.VisitAll(func(key, value []byte) {
			log.Printf("%s: %s\n", string(key), string(value))
		})

		var request struct {
			Code string `json:"code"`
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Código de autorización inválido",
			})
		}

		token, err := services.ExchangeCodeForToken(config, request.Code)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "No se pudo generar el token",
			})
		}

		return c.JSON(fiber.Map{
			"access_token": token.AccessToken,
			"expiry":       token.Expiry,
		})
	}
}
