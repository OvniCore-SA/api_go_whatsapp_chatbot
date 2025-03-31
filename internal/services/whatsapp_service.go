package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/config"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	googlecalendar "github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/googleCalendar"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/openaiassistantdtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/openaiassistantdtos/openairuns"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/whatsapp"
	metaapi "github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/whatsapp/metaApi"
	whatsappservicedto "github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/whatsapp_service_DTO"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities/filters"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
	"golang.org/x/exp/rand"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

type WhatsappService struct {
	usersService           *UsersService
	logsService            *LogsService
	openAIAssistantService *OpenAIAssistantService
	utilService            *UtilService
	userSessions           map[string]*whatsapp.UserSession // Mapa para almacenar sesiones por PhoneNumberID
	numberPhone            *NumberPhonesService
	messagesRepository     *mysql_client.MessagesRepository
	assistantService       *AssistantService
	configurationService   *ConfigurationsService
	googleCalendarService  *GoogleCalendarService
	oauthConfig            *oauth2.Config
	eventsService          EventsService
	threadService          *ThreadService
}

func NewWhatsappService(usersService *UsersService, logsService *LogsService, openAIAssistantService *OpenAIAssistantService, utilService *UtilService, numberPhone *NumberPhonesService, messagesRepository *mysql_client.MessagesRepository, assistantService *AssistantService, configurationService *ConfigurationsService, googleCalendarService *GoogleCalendarService, oauthConfig *oauth2.Config, eventsService EventsService, threadService *ThreadService) *WhatsappService {
	return &WhatsappService{
		usersService:           usersService,
		logsService:            logsService,
		openAIAssistantService: openAIAssistantService,
		utilService:            utilService,
		userSessions:           make(map[string]*whatsapp.UserSession),
		numberPhone:            numberPhone,
		messagesRepository:     messagesRepository,
		assistantService:       assistantService,
		configurationService:   configurationService,
		googleCalendarService:  googleCalendarService,
		oauthConfig:            oauthConfig,
		eventsService:          eventsService,
		threadService:          threadService,
	}
}

// Funcion para pruebas rapidas de funciones

func (service *WhatsappService) DemoFunctionWhatsappService(templateName, numberPhoneIDWhatsapp, tokenWhatsapp string) error {
	// Crear el mensaje template
	messageTemplate := metaapi.NewBodyWhatsappTemplateCRUD(
		"summary",
		"startDate",
		"formattedEnd",
		"contactToString",
		"CodeEvent",
		"5493794869394",
		templateName,
	)

	// Enviar mensaje template
	err := service.SendMessageTemplate(
		messageTemplate,
		numberPhoneIDWhatsapp,
		tokenWhatsapp,
	)

	if err != nil {
		fmt.Printf("ERROR AL NOTIFICAR EVENTO AL CLIENTE,\nERROR: %s \nCódigo de evento: %s\n", err, "adsfasd")
	}
	return nil
}

func (service *WhatsappService) HandleIncomingMessageWithAssistant(response whatsapp.ResponseComplet) error {
	for _, entry := range response.Entry {
		fmt.Println("Len Entry: ", len(response.Entry))
		fmt.Println("ID Entry: ", entry.ID)
		for _, change := range entry.Changes {
			for _, message := range change.Value.Messages {

				fmt.Println("ID Message: ", message.ID)
				fmt.Print(message.Text.Body)

				// Verificar si el mensaje ya existe por su message_id_whatsapp
				exists, err := service.messagesRepository.ExistsByMessageID(message.ID)
				if err != nil {
					fmt.Println("MENSAJE REPETIDO")
					return fmt.Errorf("failed to check message existence: %w", err)
				}

				// Si ya se proceso el mensaje se retorna nil
				if exists {
					fmt.Printf("Message with ID %s already processed\n", message.ID)
					return nil
				}

				// Extraer información básica
				sender, text, _, _, phoneNumberID, err := extractMessageInfo(change.Value)
				if err != nil {
					log.Printf("Error extracting message info: %v", err)
					return err
				}

				// Buscar el número de teléfono asociado
				numberPhone, err := service.findNumberPhoneByPhoneID(phoneNumberID)
				if err != nil {
					log.Printf("Error finding NumberPhone: %v", err)
					return err
				}

				// Buscar el contacto asociado
				contact, err := service.findOrCreateContact(numberPhone, sender)
				if err != nil {
					log.Printf("Error finding or creating contact: %v", err)
					return err
				}

				if contact.IsBlocked {
					log.Printf("contacto bloqueado")
					return nil
				}

				// Manejar el mensaje con OpenAI
				err = service.handleMessageWithOpenAI(contact, text, message.ID, numberPhone)
				if err != nil {
					log.Printf("Error handling message with OpenAI: %v", err)
					return err
				}
			}
		}
	}
	return nil
}

func (service *WhatsappService) findNumberPhoneByPhoneID(WhatsappNumberPhoneID string) (*entities.NumberPhone, error) {
	// Busca en la base de datos el número de teléfono asociado al ID recibido
	var numberPhone entities.NumberPhone
	err := config.DB.Where("whatsapp_number_phone_id = ?", WhatsappNumberPhoneID).First(&numberPhone).Error
	if err != nil {
		return nil, fmt.Errorf("number phone not found: %v", err)
	}
	return &numberPhone, nil
}

func (service *WhatsappService) findOrCreateContact(numberPhone *entities.NumberPhone, sender string) (*entities.Contact, error) {
	// Busca el contacto en la base de datos o crea uno nuevo
	var contact entities.Contact
	err := config.DB.Where("number_phones_id = ? AND number_phone = ?", numberPhone.ID, sender).First(&contact).Error
	if err == nil {
		return &contact, nil
	}

	senderInt64, err := strconv.Atoi(sender)
	if err == nil && contact.ID > 0 {
		return &contact, nil
	}

	// Crear un nuevo contacto si no existe
	contact = entities.Contact{
		NumberPhonesID: numberPhone.ID,
		NumberPhone:    int64(senderInt64),
	}
	err = config.DB.Create(&contact).Error
	if err != nil {
		return nil, fmt.Errorf("error creating contact: %v", err)
	}
	return &contact, nil
}

func (service *WhatsappService) handleMessageWithOpenAI(contact *entities.Contact, text, messageID string, numberPhone *entities.NumberPhone) error {
	// Configurar el asistente
	assistant, err := service.assistantService.FindAssistantById(numberPhone.AssistantsID)
	if err != nil {
		return fmt.Errorf("assistant not found: %v", err)
	}

	// Crear o usar el Thread existente
	thread, err := service.threadService.GetOrCreateThread(entities.MapEntityToContactDto(*contact), assistant)
	if err != nil {
		return fmt.Errorf("error obteniendo o creando thread: %v", err)
	}

	// Obtengo el vector_store que usa el assistant
	files, err := service.assistantService.serviceFile.GetFileByAssistantID(assistant.ID)
	if err != nil {
		return fmt.Errorf("error GetFileByAssistantID: %v", err)
	}

	var fileAssistant dtos.FileDto
	if len(files) > 0 {
		fileAssistant = files[len(files)-1]
		// Asigno el archivo al hilo.
		err = service.openAIAssistantService.EjecutarThread(thread.OpenaiThreadsId, []string{fileAssistant.OpenaiVectorStoreIDs})
		if err != nil {
			return fmt.Errorf("error EjecutarThread: %v", err)
		}
	} else {
		fmt.Println("Assistant sin file")
	}

	// Guardar el mensaje del contacto en la base de datos
	err = service.messagesRepository.Create(entities.Message{
		NumberPhonesID:    numberPhone.ID,
		ContactsID:        contact.ID,
		MessageText:       text,
		MessageIdWhatsapp: messageID,
		IsFromBot:         false,
	})
	if err != nil {
		return fmt.Errorf("error saving contact message: %v", err)
	}

	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		return fmt.Errorf("error cargando la zona horaria: %v", err)
	}

	currentTime := time.Now().In(loc)

	// Obtener el día de la semana en inglés (e.g. "Monday")
	weekdayEnglish := currentTime.Weekday().String()

	// Mapa para traducir el nombre del día al español
	daysInSpanish := map[string]string{
		"Sunday":    "domingo",
		"Monday":    "lunes",
		"Tuesday":   "martes",
		"Wednesday": "miércoles",
		"Thursday":  "jueves",
		"Friday":    "viernes",
		"Saturday":  "sábado",
	}

	// Obtener el día en español; si no existe en el mapa (caso muy raro), se usa el valor en inglés
	weekdaySpanish, ok := daysInSpanish[weekdayEnglish]
	if !ok {
		weekdaySpanish = weekdayEnglish
	}

	// Formatear fecha y hora en el formato habitual, agregando el día de la semana al inicio
	// Ejemplo: "lunes, 31/01/2025 15:04"
	formattedTime := fmt.Sprintf("%s, %s", strings.Title(weekdaySpanish), currentTime.Format("02/01/2006 15:04"))

	// Formatear los días de apertura y los horarios de trabajo
	availableDaysText := service.utilService.FormatOpeningDays(assistant.OpeningDays)
	workingHoursText := service.utilService.FormatWorkingHours(assistant.WorkingHours)

	// Crear el texto final
	text += fmt.Sprintf("\n\nFecha y hora actual en Argentina: %s\n%s\n%s", formattedTime, availableDaysText, workingHoursText)

	// Enviar el mensaje a OpenAI
	response, err := service.InteractWithAssistant(thread.OpenaiThreadsId, assistant.OpenaiAssistantsID, text)
	if err != nil {
		return fmt.Errorf("error sending message to OpenAI: %v", err)
	}

	// Este valor lo usaremos como respuesta "por defecto" en caso de error o fallback
	responseUser := "Podrías ser más específico, por favor?"

	// 2. Parsear la respuesta en la estructura AssistantResponse
	assistantResp, err := parseAssistantResponse(response)
	if err != nil {
		// Si hay un error en el parseo, lo mostramos y devolvemos nil (o podrías manejarlo distinto)
		fmt.Println("Error parseando la respuesta del Assistant:", err.Error())
		// Guardamos en la base de datos la respuestaUser (fallback) y retornamos
		if saveErr := saveMessageWithUniqueID(service, int(numberPhone.NumberPhone), int(contact.ID), response); saveErr != nil {
			return fmt.Errorf("error saving fallback contact message: %v", saveErr)
		}

		// Enviar el fallback al usuario
		contactToString := strconv.Itoa(int(contact.NumberPhone))
		message := metaapi.NewSendMessageWhatsappBasic(response, contactToString)
		sendErr := service.sendMessageBasic(message, strconv.FormatInt(numberPhone.WhatsappNumberPhoneId, 10), numberPhone.TokenPermanent)
		if sendErr != nil {
			return fmt.Errorf("error sending fallback response to user: %v", sendErr)
		}

		return nil
	}

	// 3. Lógica según assistantResp.Function
	switch assistantResp.Function {

	case "getMeetingDetails":
		fmt.Print("getMeetingDetails")
	case "getMeetings":

		var eventsDB []entities.Events
		startDateStr := assistantResp.UserData.DateToSearch
		// Si el contacto es el dueño configurado para recibir notificaciones, se buscan todos los eventos que estan relacionados a ese NumberPhone. Si no, se obtienen los eventos de un contactos.
		if numberPhone.NumberPhoneToNotify == contact.NumberPhone {
			eventsDB, err = service.eventsService.GetEventsByContactDateAndNumberPhone(contact.ID, startDateStr, assistant.ID)
			if err != nil {
				return fmt.Errorf("error retrieving events by contact, date, and numberPhoneID: %v", err)
			}
		} else {

			startDateStr := assistantResp.UserData.DateToSearch
			formattedTime := currentTime.Format("2006-01-02 15:04:05")
			eventsDB, err = service.eventsService.GetEventByContactAndDate(contact.ID, startDateStr, formattedTime)
			if err != nil {
				return fmt.Errorf("error retrieving events by contact, date, and time: %v", err)
			}
		}

		parsedStart, err := service.utilService.ParseDateString(startDateStr, FormatTypeDateForPersona)
		if err != nil {
			parsedStart = time.Now() // fallback en caso de error
		}

		formattedStartStr := parsedStart.Format("02-01-2006")
		if len(eventsDB) == 0 {
			responseUser = fmt.Sprintf("⚠️ *No hay eventos programados para el %s.*\n\nParece que no tienes eventos agendados para esta fecha. Si necesitas crear alguno o tienes alguna consulta, no dudes en contactarnos. 😊", formattedStartStr)
			break
		}

		// Formatear la cabecera con un mensaje claro y visualmente atractivo
		responseUser = fmt.Sprintf("🌟 *Eventos programados para el %s*:\n\n", formattedStartStr)

		for i, event := range eventsDB {
			// Intentar parsear las fechas (se asume formato RFC3339)
			startTime, err := time.Parse(time.RFC3339, event.StartDate)
			if err != nil {
				startTime = time.Now() // fallback en caso de error
			}
			endTime, err := time.Parse(time.RFC3339, event.EndDate)
			if err != nil {
				endTime = time.Now() // fallback en caso de error
			}
			// Formatear las fechas en un formato más amigable (solo hora)
			formattedStart := startTime.Format("15:04")
			formattedEnd := endTime.Format("15:04")

			// Agregar los detalles del evento con un formato más visual y emojis para claridad
			responseUser += fmt.Sprintf(
				"🔶 *Evento #%d:* %s\n⏰ *Hora de Inicio:* %s\n⏳ *Hora de Fin:* %s\n\n🔏 *Código:* %s\n────────────────\n",
				i+1,
				event.Summary,
				formattedStart,
				formattedEnd,
				event.CodeEvent,
			)
		}

	case "createMeeting":
		fmt.Println("Creando evento...")
		fmt.Print("========= RESPONSE ASSISTANT ==========")
		fmt.Print(response)
		fmt.Print("========== RESPONSE PARSED ============")
		fmt.Print(assistantResp)
		fmt.Print("===================")

		formattedTime := currentTime.Format("2006-01-02 15:04:05")

		endDateStrToDate, err := time.Parse("2006-01-02T15:04:05", assistantResp.UserData.MeetingDate)
		if err != nil {
			fmt.Println("Error al parsear la fecha:", err)
			return err
		}

		// Verificar si la fecha y hora están dentro del rango de trabajo del asistente
		isAvailable, err := service.assistantService.IsWithinWorkingHours(assistant.ID, endDateStrToDate)
		if err != nil {
			fmt.Println("Error al verificar las horas de trabajo:", err)
			return err
		}
		fmt.Print("\nIsAvailable?: ", isAvailable)
		if !isAvailable {
			responseUser = fmt.Sprintf(" Lo siento 🙁, no estamos disponible en el horario seleccionado. Recuerda que el horario de atencion es de: %s Por favor, elige otro horario 💪", assistant.WorkingHours)
			break
		}

		dateToSearch := endDateStrToDate.Format("2006-01-02")

		// Verificar si el contacto ya tiene la cantidad maxima de eventos posibles en un dia.
		eventsInDate, err := service.eventsService.GetEventsByContactDateAndNumberPhone(contact.ID, dateToSearch, assistant.ID)
		if err != nil {
			return fmt.Errorf("error retrieving events by contact, date, and numberPhoneID: %v", err)
		}
		if len(eventsInDate) >= int(assistant.EventCountPerDay) {
			responseUser = fmt.Sprintf("Lo siento 🙁. La cantidad máxima de %s por día es de %d. Si deseas cancelar o modificar alguno, solo hazmelo saber.💪", assistant.EventType, assistant.EventCountPerDay)
			break
		}

		eventsExists, err := service.eventsService.GetEventByContactAndDate(contact.ID, dateToSearch, formattedTime)
		if err != nil {
			log.Printf("Error retrieving event by contact and date: %v", err)
			return fmt.Errorf("failed to create event: %v", err)
		}
		if len(eventsExists) > 0 && assistant.EventCountPerDay == 1 {
			event := eventsExists[0] // Tomamos el primer evento porque solo debe haber uno activo

			// Intentamos parsear las fechas del evento
			startTime, err := time.Parse(time.RFC3339, event.StartDate)
			if err != nil {
				startTime = time.Now() // fallback en caso de error
			}
			endTime, err := time.Parse(time.RFC3339, event.EndDate)
			if err != nil {
				endTime = time.Now() // fallback en caso de error
			}

			// Formateamos las fechas en un formato amigable
			formattedStart := startTime.Format("15:04")
			formattedEnd := endTime.Format("15:04")

			// Mensaje para informar al usuario que ya tiene un turno en esa fecha
			responseUser = fmt.Sprintf(
				"⚠️ *Ya tienes un turno programado para esta fecha.*\n\n📅 *Fecha:* %s\n⏰ *Hora de Inicio:* %s\n⏳ *Hora de Fin:* %s\n🔏 *Código:* %s\n\nSi deseas cambiar la fecha y hora de tu turno o cancelarlo, por favor háznoslo saber.",
				endDateStrToDate.Format("02-01-2006"), formattedStart, formattedEnd, event.CodeEvent,
			)

			break
		}

		// Generamos un código único para el evento
		code, err := service.eventsService.GenerateUniqueCode()
		if err != nil {
			log.Printf("Error generating unique event code: %v", err)
			return fmt.Errorf("error creating event: %v", err)
		}

		startDateStrToDate, err := time.Parse("2006-01-02T15:04:05", assistantResp.UserData.MeetingDate)
		if err != nil {
			fmt.Println("Error al parsear la fecha:", err)
			return err
		}

		startDateToStr := startDateStrToDate.Format("2006-01-02 15:04:05")

		// Sumar 30 minutos
		endDate := endDateStrToDate.Add(time.Duration(assistant.EventDuration) * time.Minute)
		endDateStr := endDate.Format("2006-01-02T15:04:05")

		eventDTO := dtos.EventsDto{
			Summary:      assistantResp.UserData.UserName,
			Description:  "Contacto: " + assistantResp.UserData.UserEmail + "\n Tel: " + strconv.Itoa(int(contact.NumberPhone)),
			StartDate:    startDateToStr,
			EndDate:      endDateStr,
			AssistantsID: assistant.ID,
			ContactsID:   contact.ID,
			CodeEvent:    code, // Genero un codigo único para el evento
		}

		// Creo el evento en la base de datos

		if assistant.AccountGoogle {
			context := context.Background()
			token, err := service.googleCalendarService.GetOrRefreshToken(int(assistant.ID), service.oauthConfig, context)
			if err != nil {
				return err
			}

			// Crear el evento
			event := &calendar.Event{
				Summary:     assistant.EventType + " - " + eventDTO.Summary,
				Description: assistantResp.UserData.UserName + ", " + assistantResp.UserData.UserEmail,
				Start: &calendar.EventDateTime{
					DateTime: assistantResp.UserData.MeetingDate,
					TimeZone: "UTC-3",
				},
				End: &calendar.EventDateTime{
					DateTime: endDateStr,
					TimeZone: "UTC-3",
				},
				Attendees: []*calendar.EventAttendee{
					{Email: assistantResp.UserData.UserEmail},
				},
				ConferenceData: &calendar.ConferenceData{
					CreateRequest: &calendar.CreateConferenceRequest{
						RequestId: fmt.Sprintf("meet-%d", time.Now().Unix()), // ID único
						ConferenceSolutionKey: &calendar.ConferenceSolutionKey{
							Type: "hangoutsMeet",
						},
					},
				},
			}
			eventGoogleCalendar, err := service.googleCalendarService.CreateGoogleCalendarEvent(token, context, event)
			if err != nil {
				log.Println("Error al crear evento en google calendar. " + err.Error())
			}

			if eventGoogleCalendar != nil {
				eventDTO.EventGoogleCalendarID = eventGoogleCalendar.Id
			} else {
				eventDTO.EventGoogleCalendarID = "eventGoogleCalendar_default"
			}

		}

		err = service.eventsService.Create(eventDTO)
		if err != nil {
			log.Printf("\nCreate(eventDTO): %s", err.Error())
			return fmt.Errorf("error creating event: %v", err)
		}

		// Parsear las fechas en el formato esperado
		startDateTime, err := time.Parse("2006-01-02T15:04:05", assistantResp.UserData.MeetingDate)
		if err != nil {
			fmt.Printf("Error al procesar la fecha de inicio: %v\n", err)
			return err
		}

		endDateTime, err := time.Parse("2006-01-02T15:04:05", endDateStr)
		if err != nil {
			fmt.Printf("Error al procesar la fecha de finalización: %v\n", err)
			return err
		}

		// Extraer componentes de la fecha
		formattedStart := startDateTime.Format("02/01/2006 15:04")
		formattedEnd := endDateTime.Format("02/01/2006 15:04")

		// Mensaje de respuesta
		responseUser = fmt.Sprintf(
			"✅ ¡Tu nuevo evento se agendó con éxito! 📅\n\n🕒 Inicio: %s \n🕒 Fin: %s.\n🔏 Código: %s\n\nTe esperamos... ¡Que tengas un excelente día! 😊",
			formattedStart, formattedEnd, eventDTO.CodeEvent)

		// Notificar al cliente
		//  Enviar la notificacion al cliente de que un usuario registró un turno o reunion
		contactToString := strconv.Itoa(int(contact.NumberPhone))
		contactForNotifyToString := strconv.Itoa(int(numberPhone.NumberPhoneToNotify))

		// Crear el mensaje template
		messageTemplate := metaapi.NewBodyWhatsappTemplateCRUD(
			eventDTO.Summary,
			startDateToStr,
			formattedEnd,
			contactToString,
			eventDTO.CodeEvent,
			contactForNotifyToString,
			metaapi.TemplateEventoCreado,
		)

		// Enviar mensaje template
		err = service.SendMessageTemplate(
			messageTemplate,
			strconv.FormatInt(numberPhone.WhatsappNumberPhoneId, 10),
			numberPhone.TokenPermanent,
		)

		if err != nil {
			fmt.Printf("ERROR AL NOTIFICAR EVENTO AL CLIENTE,\nERROR: %s \nCódigo de evento: %s\n", err, eventDTO.CodeEvent)
		}

	case "updateEvents":
		// Convertir la fecha actual al formato adecuado (se asume RFC3339)
		currentTimeStr := time.Now().Format(time.RFC3339)

		// Se obtiene el evento del contacto para la fecha indicada y con hora >= a la actual
		eventFound, err := service.eventsService.GetEventByCodeEvent(contact.ID, assistantResp.UserData.EventCode)
		if err != nil {
			responseUser = "Lo siento, pero no pudimos encontrar el turno que mencionas. Te puedes ayudar viendo los turnos que tenes en la fecha que quieres consultar. 😊"
			break
		}
		if eventFound.ID <= 0 {
			return fmt.Errorf("no se encontró un evento para el contacto %d en la fecha %s con hora mayor o igual a %s", contact.ID, assistantResp.UserData.MeetingDate, currentTimeStr)
		}

		newDateStrToDate, err := time.Parse("2006-01-02T15:04:05", assistantResp.UserData.NewDate)
		if err != nil {
			fmt.Println("Error al parsear la fecha:", err)
			return err
		}

		// Verificar si el asistente tiene disponibilidad en la nueva fecha y hora
		isAvailable, err := service.assistantService.IsWithinWorkingHours(assistant.ID, newDateStrToDate)
		if err != nil {
			fmt.Println("Error al verificar las horas de trabajo:", err)
			return err
		}

		if !isAvailable {
			responseUser = "⚠️ Lo siento, el asistente no está disponible en el horario seleccionado para actualizar el turno. Por favor, elige otro horario."
			break
		}

		// Sumar 30 minutos
		endDate := newDateStrToDate.Add(time.Duration(assistant.EventDuration) * time.Minute)
		endDateStr := endDate.Format("2006-01-02T15:04:05")

		// Actualizar el evento en la base de datos con la información nueva
		eventDTO := dtos.EventsDto{
			ID:           eventFound.ID,
			Summary:      eventFound.Summary,
			Description:  eventFound.Description,
			StartDate:    assistantResp.UserData.NewDate,
			EndDate:      endDateStr,
			AssistantsID: assistant.ID,
			ContactsID:   contact.ID,
			CodeEvent:    eventFound.CodeEvent,
			CreatedAt:    eventFound.CreatedAt,
		}

		err = service.eventsService.Update(eventDTO)
		if err != nil {
			return err
		}

		if assistant.AccountGoogle {
			context := context.Background()
			token, err := service.googleCalendarService.GetOrRefreshToken(int(assistant.ID), service.oauthConfig, context)
			if err != nil {
				return err
			}

			// Crear el evento
			event := &googlecalendar.EventRequest{
				Summary:     eventDTO.Summary,
				Description: assistantResp.UserData.UserName + ", " + assistantResp.UserData.UserEmail,
				Start:       eventDTO.StartDate,
				End:         eventDTO.EndDate,
				ContactsID:  uint(contact.ID),
			}

			if _, err := service.googleCalendarService.UpdateGoogleCalendarEvent(token, context, eventDTO.EventGoogleCalendarID, event); err != nil {
				log.Println("Error al crear evento en google calendar. " + err.Error())
			}
		}

		responseUser = "✅ Tu evento ha sido modificada con éxito. Si nesesitas cualquier otra cosa, estoy acá para ayudarte 😊"

		// Notificar al cliente
		//  Enviar la notificacion al cliente de que un usuario registró un turno o reunion
		contactToString := strconv.Itoa(int(numberPhone.NumberPhoneToNotify))
		endDateStr = endDate.Format("2006-01-02 15:04:05")

		// Crear el mensaje template
		messageTemplate := metaapi.NewBodyWhatsappTemplateCRUD(
			eventDTO.Summary,
			assistantResp.UserData.NewDate,
			endDateStr,
			contactToString,
			eventDTO.CodeEvent,
			contactToString,
			metaapi.TemplateEventoModificado,
		)

		// Enviar mensaje template
		err = service.SendMessageTemplate(
			messageTemplate,
			strconv.FormatInt(numberPhone.WhatsappNumberPhoneId, 10),
			numberPhone.TokenPermanent,
		)

		if err != nil {
			fmt.Printf("ERROR AL NOTIFICAR EVENTO AL CLIENTE,\nERROR: %s \nCódigo de evento: %s\n", err, eventDTO.CodeEvent)
		}

	case "deleteEvent":
		event, err := service.eventsService.GetEventByCodeEvent(contact.ID, assistantResp.UserData.EventCode)
		if err != nil {
			fmt.Println(err.Error())
			responseUser = fmt.Sprint("Lo siento, ocurrió un error al eliminar reunion.\n Podrías intentar mas tarde o simplemente modificar tu reunión. \nSi nesesitas cualquier otra cosa, estoy acá para ayudarte 😊")
		}

		err = service.eventsService.Cancel(assistantResp.UserData.EventCode)
		if err != nil {
			return err
		}
		responseUser = fmt.Sprintf("✅ Su turno con el código '%s' ha sido cancelado con éxito. Si nesesitas cualquier otra cosa, estoy acá para ayudarte 😊", assistantResp.UserData.EventCode)

		if assistant.AccountGoogle {
			context := context.Background()
			token, err := service.googleCalendarService.GetOrRefreshToken(int(assistant.ID), service.oauthConfig, context)
			if err != nil {
				return err
			}

			if err := service.googleCalendarService.DeleteGoogleCalendarEvent(token, context, event.EventGoogleCalendarID); err != nil {
				log.Println("no se pudo eliminar el evento de google: " + err.Error())
			}
		}

		// // Notificar al cliente
		// //  Enviar la notificacion al cliente de que un usuario registró un turno o reunion
		contactToString := strconv.Itoa(int(numberPhone.NumberPhoneToNotify))
		textNotifyClient := fmt.Sprintf("  \n\n 🔴 *Cancelación de Turno* \n\n🔶 %s : *%s*\n⏰ *Hora de Inicio:* %sHs.\n⏳ *Hora de Fin:* %sHs.\n🔏 *Código:* %s.\n\n⚠️ Su turno ha sido cancelado. Para más información, por favor, contáctenos.",
			service.utilService.CapitalizeFirstLetter(assistant.EventType),
			event.Summary,
			event.StartDate,
			event.EndDate,
			event.CodeEvent,
		)

		message := metaapi.NewSendMessageWhatsappBasic(textNotifyClient, contactToString)
		err = service.sendMessageBasic(message, strconv.FormatInt(numberPhone.WhatsappNumberPhoneId, 10), numberPhone.TokenPermanent)
		if err != nil {
			fmt.Printf("ERROR AL NOTIFICAR CANCELACIÓN DE EVENTO AL CLIENTE,\nERROR: %s \nCódigo de evento: %s", err, event.CodeEvent)
		}

	default:
		// El texto a mostrar al usuario viene en assistantResp.Message
		responseUser = assistantResp.Message
	}

	// 4. Guardar la respuesta que le daremos al usuario
	err = saveMessageWithUniqueID(service, int(numberPhone.ID), int(contact.ID), responseUser)
	if err != nil {
		return fmt.Errorf("error saving contact message: %v", err)
	}

	// 5. Enviar la respuesta al usuario
	contactToString := strconv.Itoa(int(contact.NumberPhone))
	message := metaapi.NewSendMessageWhatsappBasic(responseUser, contactToString)
	err = service.sendMessageBasic(message, strconv.FormatInt(numberPhone.WhatsappNumberPhoneId, 10), numberPhone.TokenPermanent)
	if err != nil {
		return fmt.Errorf("error sending response to user: %v", err)
	}

	return nil
}

func parseAssistantResponse(response string) (assistantResp *openaiassistantdtos.AssistantJSONResponse, err error) {
	// Elimina espacios en blanco extra
	clean := strings.TrimSpace(response)

	// Si empieza con ```json, quítalo
	if strings.HasPrefix(clean, "```json") {
		clean = strings.TrimPrefix(clean, "```json")
	}
	// Si acaba con ```, quítalo
	if strings.HasSuffix(strings.TrimSpace(clean), "```") {
		clean = strings.TrimSuffix(strings.TrimSpace(clean), "```")
	}

	// Vuelve a recortar espacios
	clean = strings.TrimSpace(clean)

	// Parseamos el JSON

	err = json.Unmarshal([]byte(clean), &assistantResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal: %v", err)
	}

	return assistantResp, nil
}

// Función para generar un string único
func generateUniqueID() string {
	rand.Seed(uint64(time.Now().UnixNano())) // Conversión explícita
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := make([]byte, 12) // Generará un ID de longitud 12
	for i := range id {
		id[i] = charset[rand.Intn(len(charset))]
	}
	return string(id)
}

// Lógica para guardar un mensaje con un ID único
func saveMessageWithUniqueID(service *WhatsappService, numberPhoneID, contactID int, response string) error {
	var messageID string
	var isUnique bool

	// Intentar generar un ID único
	for !isUnique {
		messageID = generateUniqueID()

		// Comprobar si el ID ya existe en la base de datos
		exists, err := service.messagesRepository.ExistsByMessageID(messageID)
		if err != nil {
			return fmt.Errorf("error checking message ID uniqueness: %v", err)
		}

		if !exists {
			isUnique = true
		}
	}

	// Guardar el mensaje en la base de datos
	err := service.messagesRepository.Create(entities.Message{
		NumberPhonesID:    int64(numberPhoneID),
		ContactsID:        int64(contactID),
		MessageIdWhatsapp: messageID,
		MessageText:       response,
		IsFromBot:         true,
	})
	if err != nil {
		return fmt.Errorf("error saving assistant response: %v", err)
	}

	return nil
}

func (s *WhatsappService) InteractWithAssistant(threadID, assistantID, message string) (response string, err error) {

	// Verificar si es seguro proceder (sin runs activos)
	safeToProceed, err := s.CheckForActiveRuns(threadID)
	if err != nil {
		return "", err // Devuelve el error si no es seguro proceder
	}
	if !safeToProceed {
		return "", fmt.Errorf("an active run is still processing")
	}

	// Crear un run para el thread con la conversación completa
	err = s.openAIAssistantService.SendMessageToThread(threadID, message, true)
	if err != nil {
		return "", fmt.Errorf("error creating message with conversation: %v", err)
	}

	runID, err := s.openAIAssistantService.CreateRunForThreadWithConversation(threadID, assistantID, nil)
	if err != nil {
		return "", fmt.Errorf("error creating run with conversation: %v", err)
	}
	fmt.Printf("Run created: %s\n", runID)

	// Esperar a que el run esté completado
	requiredAction, err := s.openAIAssistantService.WaitForRunCompletion(threadID, runID, 30, 2*time.Second)
	if err != nil {
		return "", fmt.Errorf("error waiting for run completion: %v", err)
	}

	// el assistant manda la funcion pero en un estado de required action, se pasa a JSON el objeto y se devuelve como string para ser decodificado y usado luego
	if requiredAction.RequiredAction.Type == "submit_tool_outputs" && len(requiredAction.RequiredAction.SubmitToolOutputs.ToolCalls) > 0 {
		response, err = convertToAssistantJSON(requiredAction)
		if err != nil {
			return "", fmt.Errorf("error getting messages from thread: %v", err)
		}
	} else {
		// Obtener los mensajes del thread y encontrar la respuesta del asistente
		response, err = s.openAIAssistantService.GetMessagesFromThread(threadID)
		if err != nil {
			return "", fmt.Errorf("error getting messages from thread: %v", err)
		}
	}

	return response, nil
}

// ConvertToAssistantJSON convierte OpenAIRunResponse en AssistantJSONResponse
func convertToAssistantJSON(response openairuns.OpenAIRunResponse) (string, error) {
	// Verifica que haya una acción requerida

	// Obtiene la primera toolCall (ajusta si pueden haber varias)
	toolCall := response.RequiredAction.SubmitToolOutputs.ToolCalls[0]

	// Decodifica los argumentos JSON
	var args map[string]string
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
		fmt.Printf("error al decodificar los argumentos (convertToAssistantJSON): %v", err)
		return "", fmt.Errorf("error al decodificar los argumentos: %v", err)
	}

	// Crea la respuesta
	assistantResponse := openaiassistantdtos.AssistantJSONResponse{
		Function: toolCall.Function.Name,
	}

	// Mapea los argumentos decodificados a la estructura AssistantJSONResponse
	if val, exists := args["user_name"]; exists {
		assistantResponse.UserData.UserName = val
	}
	if val, exists := args["user_email"]; exists {
		assistantResponse.UserData.UserEmail = val
	}
	if val, exists := args["meeting_date"]; exists {
		assistantResponse.UserData.MeetingDate = val
	}
	if val, exists := args["event_code"]; exists {
		assistantResponse.UserData.EventCode = val
	}
	if val, exists := args["new_date"]; exists {
		assistantResponse.UserData.NewDate = val
	}
	if val, exists := args["date_to_search"]; exists {
		assistantResponse.UserData.DateToSearch = val
	}

	// Convierte la estructura a JSON
	jsonResponse, err := json.Marshal(assistantResponse)
	if err != nil {
		return "", fmt.Errorf("error al serializar AssistantJSONResponse: %v", err)
	}

	return string(jsonResponse), nil
}

// CheckForActiveRuns checks if there are any active runs for a given thread and retries a few times if there are.
// Returns true if it is safe to proceed (no active runs), false and an error otherwise.
func (s *WhatsappService) CheckForActiveRuns(threadID string) (bool, error) {
	maxRetries := 5
	retryInterval := time.Second * 10 // 10 segundos de intervalo entre reintento

	for i := 0; i < maxRetries; i++ {
		runs, err := s.openAIAssistantService.ListRunsForThread(threadID, 20, "desc", "", "")
		if err != nil {
			return false, fmt.Errorf("error listing runs for thread: %v", err)
		}

		activeRunFound := false
		for _, run := range runs {
			if run.Status == "active" {
				activeRunFound = true
				fmt.Printf("Waiting, found active run: %s\n", run.ID)
				break
			}
			if run.Status == "requires_action" {
				activeRunFound = true
				fmt.Printf("Waiting, found requires_action run: %s\n", run.ID)
				break
			}
		}

		if !activeRunFound {
			return true, nil // No active runs, safe to proceed
		}

		// If active runs are found and it's not the last retry attempt, wait before retrying
		if i < maxRetries-1 {
			fmt.Printf("Retrying in %v seconds...\n", retryInterval.Seconds())
			time.Sleep(retryInterval)
		}
	}

	return false, fmt.Errorf("cannot create a new run, active run did not complete after %d retries", maxRetries)
}

func (s *WhatsappService) NotifyInteractions(horasAtras uint) error {
	// Obtener los números de teléfono asociados al asistente
	filter := filters.AssistantsFiltro{
		NumberPhoneToNotifyNotEmpty: true,
		UpladContacts:               true,
	}
	numbers, err := s.numberPhone.repository.ListByFilter(filter)
	if err != nil {
		return fmt.Errorf("error retrieving numbers: %v", err)
	}

	var interactionSummaries []whatsappservicedto.InteractionSummary
	var sixHoursAgo time.Time

	for _, number := range numbers {
		var usersContactos []whatsappservicedto.UserContactInfo
		// Filtrar mensajes de las últimas 6 horas
		for _, contact := range number.Contacts {
			if horasAtras == 18 {
				sixHoursAgo = time.Now().Add(-18 * time.Hour)
			} else {
				sixHoursAgo = time.Now().Add(-6 * time.Hour)
			}

			messages, err := s.messagesRepository.GetMessagesByNumber(number.ID, contact.ID, sixHoursAgo)
			if err != nil {
				return fmt.Errorf("error retrieving messages: %v", err)
			}

			if len(messages) <= 0 {
				continue
			}

			// messagePromt := "Quiero que revises este hilo de conversación y, dentro de un lapso de 6 horas atrás, identifiques si algún usuario acordó una reunión. Para que sea válida, debe haber proporcionado explícitamente los siguientes datos: su nombre y correo electrónico, siendo opcional el número de celular.\nSi encuentras estos datos, respóndeme únicamente en el siguiente formato y nada más:\n\n[Nombre del usuario]\n[Correo del usuario]\n[Número de celular] (indica 'No proporcionado' si no lo dio). SOLO LOS DATOS, NO EL NOMBRE DE CADA DATO.\n\nSi no encuentras ninguna reunión dentro del lapso indicado o si los datos no están completos, responde únicamente con 'No se encontró información de reuniones.' No proporciones nada adicional ni interpretes respuestas parciales."
			// response, err := s.InteractWithAssistant(contact.OpenaiThreadsID, number.Assistant.OpenaiAssistantsID, messagePromt)
			// if err != nil {
			// 	return fmt.Errorf("error interacting with assistant: %v", err)
			// }
			response := ""

			// Construir el mensaje para el asistente

			//response := "3872937497\nemanuel.garcia@ovnix.com"

			// Procesar la respuesta
			lines := strings.Split(response, "\n")

			// Verificar que haya al menos 4 líneas (nombre, teléfono, email, fecha/hora)
			if len(lines) < 3 {
				fmt.Println("Error: Respuesta incompleta o mal formateada")
				fmt.Print(lines)
				continue
			}

			// Verificar que cada campo no esté vacío
			for i, line := range lines[:3] {
				if strings.TrimSpace(line) == "" {
					fmt.Printf("Error: Campo vacío en la línea %d\n", i+1)
					continue
				}
			}

			// Asignar los datos solo si pasan las verificaciones
			contactInfo := whatsappservicedto.UserContactInfo{
				Nombre:   strings.TrimSpace(lines[0]),
				Telefono: strings.TrimSpace(lines[1]),
				Email:    strings.TrimSpace(lines[2]),
			}

			usersContactos = append(usersContactos, contactInfo)

		}
		// Agregar la interacción al resumen
		interactionSummaries = append(interactionSummaries, whatsappservicedto.InteractionSummary{
			NumberPhoneEntity: number,
			NumberPhoneID:     number.ID,
			NumberPhone:       number.NumberPhone,
			Contacts:          usersContactos,
		})
	}

	// Enviar notificaciones a los números designados
	for _, summary := range interactionSummaries {
		var message strings.Builder
		message.WriteString("👋 *Hola,*\n\n")

		message.WriteString("Espero que estés muy bien. 😊 Desde el equipo de *OvniCore*, queremos informarte que los siguientes usuarios han solicitado coordinar una reunión:\n\n")

		if len(summary.Contacts) <= 0 {
			continue
		}

		for _, contact := range summary.Contacts {
			message.WriteString(fmt.Sprintf("👤 *Nombre:* %s \n", contact.Nombre))
			message.WriteString(fmt.Sprintf("📞 *Teléfono:* %s\n", contact.Telefono))
			message.WriteString(fmt.Sprintf("✉️ *Correo:* %s\n", contact.Email))
			message.WriteString(fmt.Sprintf("📅 *Fecha y Hora de la Reunión:* %s\n\n", contact.FechaHora))

		}

		message.WriteString("📅 *Por favor, te pedimos que contactes a estos usuarios para coordinar una reunión en el horario que sea más conveniente para ambas partes.*\n\n")
		message.WriteString("🙏 ¡Gracias por tu atención y apoyo! Si necesitas alguna información adicional, no dudes en contactarnos. 🌟\n\n")
		message.WriteString("🤝 *Saludos cordiales,*\n")
		message.WriteString("El equipo de *OvniCore* 🚀")

		err := s.SendWhatsappNotification(summary.NumberPhoneEntity, message.String())
		if err != nil {
			return fmt.Errorf("error sending WhatsApp notification: %v", err)
		}
	}

	return nil
}

func (s *WhatsappService) SendWhatsappNotification(numberPhone entities.NumberPhone, message string) error {
	// Lógica para enviar notificaciones por WhatsApp
	// Puedes usar Baileys o la API que tengas configurada para este propósito
	fmt.Printf("Enviando mensaje a %d: %s\n", numberPhone.NumberPhoneToNotify, message)

	// Enviar la respuesta al usuario
	contactToString := strconv.Itoa(int(numberPhone.NumberPhoneToNotify))
	messageForBody := metaapi.NewSendMessageWhatsappBasic(message, contactToString)
	err := s.sendMessageBasic(messageForBody, strconv.FormatInt(numberPhone.WhatsappNumberPhoneId, 10), numberPhone.TokenPermanent)
	if err != nil {
		return fmt.Errorf("error sending response to user: %v", err)
	}
	return nil
}

func (service *WhatsappService) sendMessageBasic(message metaapi.SendMessageBasic, phoneNumberId string, tokenApiWhatsapp string) error {

	reqJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error al convertir el cuerpo de la solicitud a JSON:", err)
		return err
	}

	client := &http.Client{}

	base, err := url.Parse(os.Getenv("WHATSAPP_URL") + "/" + os.Getenv("WHATSAPP_VERSION") + "/" + phoneNumberId + "/messages")
	fmt.Println("Url: ", base.String())
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	req, err := http.NewRequest("POST", base.String(), bytes.NewBuffer(reqJSON))
	if err != nil {
		fmt.Println("Error al construir solicitud POST:", err)
		return err
	}

	req.Header.Set("Authorization", "Bearer "+tokenApiWhatsapp)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if resp.StatusCode != 200 {
		respuesta, erro := VerCuerpoRespuesta(resp)
		if erro != nil {
			fmt.Println("Error al mostrar respuesta erronea")
		}
		fmt.Println(respuesta)
		return fmt.Errorf("Codigo error: " + resp.Status)
	}

	defer resp.Body.Close()

	return err
}

// Función que muestra la respuesta que devuelve una API.
func VerCuerpoRespuesta(resp *http.Response) (body []byte, err error) {
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//log.Println(string(body))
	//fmt.Print()
	return
}

// Extrae la información del mensaje
func extractMessageInfo(value whatsapp.Value) (sender, text, fechaMessage, messageType, phoneNumberID string, err error) {

	fmt.Println("Fecha: " + value.Messages[0].Timestamp)
	sender = value.Messages[0].From
	text = value.Messages[0].Text.Body
	messageType = value.Messages[0].Type
	phoneNumberID = value.Metadata.PhoneNumberID

	return sender, text, fechaMessage, messageType, phoneNumberID, nil
}

func (service *WhatsappService) SendMessageTemplate(message metaapi.SendMessageTemplate, phoneNumberId, tokenApiWhatsapp string) error {
	reqJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error al convertir la plantilla a JSON:", err)
		return err
	}

	client := &http.Client{}

	base, err := url.Parse(os.Getenv("WHATSAPP_URL") + "/" + os.Getenv("WHATSAPP_VERSION") + "/" + phoneNumberId + "/messages")
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	req, err := http.NewRequest("POST", base.String(), bytes.NewBuffer(reqJSON))
	if err != nil {
		fmt.Println("Error al construir solicitud POST:", err)
		return err
	}

	req.Header.Set("Authorization", "Bearer "+tokenApiWhatsapp)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error al realizar la petición:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respuesta, erro := VerCuerpoRespuesta(resp)
		if erro != nil {
			fmt.Println("Error al mostrar respuesta erronea")
		}
		fmt.Println(respuesta)
		return fmt.Errorf("Codigo error: %s", resp.Status)
	}

	return nil
}
