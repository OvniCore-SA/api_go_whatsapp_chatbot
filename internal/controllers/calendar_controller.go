package controllers

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	googlecalendar "github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/googleCalendar"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

// GetRequestDetails obtiene todos los parámetros enviados en la solicitud.
func GetRequestDetails() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Obtener los headers
		headers := make(map[string]string)
		c.Request().Header.VisitAll(func(key, value []byte) {
			headers[string(key)] = string(value)
		})

		// Obtener los parámetros de la URL
		params := make(map[string]string)
		for _, param := range c.Route().Params {
			params[param] = c.Params(param)
		}

		// Obtener los parámetros de consulta (query params)
		queryParams := make(map[string]string)
		for key, values := range c.Queries() {
			queryParams[key] = values
		}

		// Obtener el cuerpo de la solicitud (body)
		var body map[string]interface{}
		if err := c.BodyParser(&body); err != nil {
			body = map[string]interface{}{
				"error": "No se pudo parsear el body",
			}
		}

		// Retornar todos los datos en JSON
		return c.JSON(fiber.Map{
			"headers":     headers,
			"params":      params,
			"queryParams": queryParams,
			"body":        body,
		})
	}
}

// GetCalendarEvents obtiene los eventos del calendario.
func GetCalendarEventsByDate(service *services.GoogleCalendarService, config *oauth2.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		assistantID, err := parseAssistantID(c.Query("assistant_id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		startDateStr := c.Query("start_date") // Fecha inicial en formato "dd-mm-aaaa"
		endDateStr := c.Query("end_date")     // Fecha final en formato "dd-mm-aaaa"

		startDate, err := time.Parse("02-01-2006", startDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid start_date format.(sended dd-mm-aaaa)"})
		}

		endDate, err := time.Parse("02-01-2006", endDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid end_date format.(sended dd-mm-aaaa)"})
		}

		token, err := service.GetOrRefreshToken(assistantID, config, c.Context())
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		events, err := service.FetchGoogleCalendarEventsByDate(token, c.Context(), startDate, endDate)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if len(events.Items) == 0 {
			return c.JSON(fiber.Map{
				"status":  true,
				"data":    nil,
				"message": "No se encontraron eventos en el rango especificado.",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Eventos obtenidos con éxito.",
			"data":    formatEvents(events.Items),
			"status":  true,
		})
	}
}

// AddCalendarEvent crea un nuevo evento en Google Calendar.
func AddCalendarEvent(service *services.GoogleCalendarService, config *oauth2.Config, contactService *services.ContactsService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		assistantID, err := parseAssistantID(c.Query("assistant_id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		var eventRequest googlecalendar.EventRequest
		if err := c.BodyParser(&eventRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		err = eventRequest.Validate()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// Configurar el asistente
		assistantWihtGoogleConfig, err := service.AssistantService.FindAssistantById(int64(assistantID))
		if err != nil {
			return fmt.Errorf("assistant not found: %v", err)
		}

		var createdEvent *calendar.Event

		if assistantWihtGoogleConfig.AccountGoogle {
			token, err := service.GetOrRefreshToken(assistantID, config, c.Context())
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
			}

			// Crear el evento
			event := &calendar.Event{
				Summary:     eventRequest.Summary,
				Description: eventRequest.Description,
				Start: &calendar.EventDateTime{
					DateTime: eventRequest.Start,
					TimeZone: "UTC-3",
				},
				End: &calendar.EventDateTime{
					DateTime: eventRequest.End,
					TimeZone: "UTC-3",
				},
			}

			createdEvent, err = service.CreateGoogleCalendarEvent(token, c.Context(), event)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		} else {
			// cargo el dto para la creacion solo en nuestra db
			createdEvent = &calendar.Event{
				Summary:     eventRequest.Summary,
				Description: eventRequest.Description,
				Start: &calendar.EventDateTime{
					DateTime: eventRequest.Start,
					TimeZone: "UTC-3",
				},
				End: &calendar.EventDateTime{
					DateTime: eventRequest.End,
					TimeZone: "UTC-3",
				},
			}
		}

		eventDto, err := dtos.MapCalendarEventToEventsDto(createdEvent)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		codeUnique, err := service.EventsService.GenerateUniqueCode()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		eventDto.ContactsID = int64(eventRequest.ContactsID)
		eventDto.AssistantsID = assistantWihtGoogleConfig.ID
		eventDto.CodeEvent = codeUnique

		// Guardar en la base de datos usando EventsService
		err = service.EventsService.Create(eventDto)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"message": "Evento creado con éxito.",
			"data":    createdEvent,
			"status":  true,
		})
	}
}

func InsertCalendarEvent(service *services.GoogleCalendarService, config *oauth2.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		assistantID, err := parseAssistantID(c.Query("assistant_id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		eventDateStr := c.Query("event_date")      // Fecha del evento en formato "dd-mm-aaaa"
		eventTitle := c.Query("title")             // Título del evento
		eventDescription := c.Query("description") // Descripción del evento

		eventDate, err := time.Parse("02-01-2006", eventDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event_date format"})
		}

		token, err := service.GetOrRefreshToken(assistantID, config, c.Context())
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		event := &calendar.Event{
			Summary:     eventTitle,
			Description: eventDescription,
			Start: &calendar.EventDateTime{
				DateTime: eventDate.Format(time.RFC3339),
			},
			End: &calendar.EventDateTime{
				DateTime: eventDate.Add(1 * time.Hour).Format(time.RFC3339),
			},
		}

		err = service.InsertGoogleCalendarEvent(token, c.Context(), event)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"message": "Evento registrado con éxito.",
			"status":  true,
		})
	}
}

// DeleteCalendarEvent elimina un evento del Google Calendar
func DeleteCalendarEvent(service *services.GoogleCalendarService, config *oauth2.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		assistantID, err := parseAssistantID(c.Query("assistant_id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		eventID := c.Params("event_id")
		if eventID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "El event_id es requerido"})
		}

		token, err := service.GetOrRefreshToken(assistantID, config, c.Context())
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		// Llamar al servicio para eliminar el evento
		err = service.DeleteGoogleCalendarEvent(token, c.Context(), eventID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"message": "Evento eliminado con éxito.",
			"status":  true,
		})
	}
}

// UpdateCalendarEvent actualiza un evento en Google Calendar
func UpdateCalendarEvent(service *services.GoogleCalendarService, config *oauth2.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		assistantID, err := parseAssistantID(c.Query("assistant_id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		eventID := c.Params("event_id")
		if eventID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "El event_id es requerido"})
		}

		var eventRequest googlecalendar.EventRequest
		if err := c.BodyParser(&eventRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
		}

		err = eventRequest.Validate()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		token, err := service.GetOrRefreshToken(assistantID, config, c.Context())
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		// Llamar al servicio para actualizar el evento
		updatedEvent, err := service.UpdateGoogleCalendarEvent(token, c.Context(), eventID, &eventRequest)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"message": "Evento actualizado con éxito.",
			"data":    updatedEvent,
			"status":  true,
		})
	}
}

func parseAssistantID(assistantIDStr string) (int, error) {
	if assistantIDStr == "" {
		return 0, errors.New("el parámetro assistant_id es obligatorio")
	}
	return strconv.Atoi(assistantIDStr)
}

// func formatEvents(events []*calendar.Event) []fiber.Map {
// 	var result []fiber.Map
// 	for _, item := range events {
// 		start := item.Start.DateTime
// 		if start == "" {
// 			start = item.Start.Date
// 		}
// 		end := item.End.DateTime
// 		if end == "" {
// 			end = item.End.Date
// 		}

// 		result = append(result, fiber.Map{
// 			"id":            item.Id,
// 			"title":         item.Summary,
// 			"description":   item.Description,
// 			"location":      item.Location,
// 			"start":         start,
// 			"end":           end,
// 			"attendees":     extractAttendees(item.Attendees),
// 			"organizer":     item.Organizer.DisplayName,
// 			"html_link":     item.HtmlLink,
// 			"conference":    extractConferenceData(item.ConferenceData),
// 			"event_type":    item.EventType,
// 			"visibility":    item.Visibility,
// 			"creation_time": item.Created,
// 		})
// 	}
// 	return result
// }

func extractAttendees(attendees []*calendar.EventAttendee) []string {
	var attendeeNames []string
	for _, attendee := range attendees {
		if attendee.DisplayName != "" {
			attendeeNames = append(attendeeNames, attendee.DisplayName)
		} else {
			attendeeNames = append(attendeeNames, attendee.Email)
		}
	}
	return attendeeNames
}

func extractConferenceData(data *calendar.ConferenceData) string {
	if data == nil || data.EntryPoints == nil {
		return ""
	}
	for _, entry := range data.EntryPoints {
		if entry.Uri != "" {
			return entry.Uri
		}
	}
	return ""
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

func formatEvents(events []*calendar.Event) []fiber.Map {
	var result []fiber.Map
	for _, event := range events {
		start := parseEventDate(event.Start)
		end := parseEventDate(event.End)

		result = append(result, fiber.Map{
			"id":          event.Id,
			"title":       event.Summary,
			"description": event.Description,
			"start":       start,
			"end":         end,
			"location":    event.Location,
			"link":        event.HtmlLink,
		})
	}
	return result
}

func parseEventDate(eventDateTime *calendar.EventDateTime) string {
	if eventDateTime.DateTime != "" {
		t, _ := time.Parse(time.RFC3339, eventDateTime.DateTime)
		return t.Format("02-01-2006 15:04")
	}
	if eventDateTime.Date != "" {
		t, _ := time.Parse("2006-01-02", eventDateTime.Date)
		return t.Format("02-01-2006")
	}
	return "Fecha desconocida"
}
