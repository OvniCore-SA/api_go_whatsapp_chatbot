package controllers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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
		// Obtén parámetros adicionales desde la URL de la solicitud
		assistantID := c.Query("assistant_id")
		redirectURL := c.Query("redirect_url")

		// Verifica que los parámetros sean válidos
		if assistantID == "" || redirectURL == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "assistant_id y redirect_url son obligatorios",
			})
		}

		// Construir el estado personalizado con assistant_id y redirect_url
		state := fmt.Sprintf("assistant_id=%s&redirect_url=%s", assistantID, redirectURL)

		// Generar la URL de autenticación con el estado personalizado
		authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)

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

func GetAuthToken(config *oauth2.Config, googleCalendarService *services.GoogleCalendarService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Capturar el código y el estado desde los parámetros de consulta
		code := c.Query("code")
		state := c.Query("state")

		if code == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "El parámetro 'code' es obligatorio",
			})
		}

		if state == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "El parámetro 'state' es obligatorio",
			})
		}

		// Procesar el estado para extraer los parámetros adicionales
		params := make(map[string]string)
		for _, param := range strings.Split(state, "&") {
			kv := strings.SplitN(param, "=", 2)
			if len(kv) == 2 {
				params[kv[0]] = kv[1]
			}
		}

		assistantID := params["assistant_id"]
		redirectURL := params["redirect_url"]

		if assistantID == "" || redirectURL == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "No se encontraron los parámetros 'assistant_id' y 'redirect_url' en el estado",
			})
		}

		// Intercambiar el código por un token
		token, err := config.Exchange(c.Context(), code)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "No se pudo generar el token",
			})
		}
		fmt.Println("AccessToken GOOGLE: " + token.AccessToken)
		// Obtener información del usuario con el token
		client := config.Client(c.Context(), token)
		googleUserID, err := services.GetGoogleUserID(client, token)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "No se pudo obtener el ID del usuario de Google" + err.Error(),
			})
		}

		// Guardar las credenciales
		assistantIDInt, _ := strconv.Atoi(assistantID)
		err = googleCalendarService.SaveCredentials(assistantIDInt, token, googleUserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "No se pudieron guardar las credenciales",
			})
		}

		// Redirigir al usuario a la URL proporcionada con los parámetros adicionales
		return c.Redirect(fmt.Sprintf("%s?status=success", redirectURL))
	}
}
