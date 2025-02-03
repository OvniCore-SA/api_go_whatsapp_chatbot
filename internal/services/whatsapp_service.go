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
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/whatsapp"
	metaapi "github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/whatsapp/metaApi"
	whatsappservicedto "github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/whatsapp_service_DTO"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities/filters"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
	"golang.org/x/exp/rand"
	"golang.org/x/oauth2"
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
}

func NewWhatsappService(usersService *UsersService, logsService *LogsService, openAIAssistantService *OpenAIAssistantService, utilService *UtilService, numberPhone *NumberPhonesService, messagesRepository *mysql_client.MessagesRepository, assistantService *AssistantService, configurationService *ConfigurationsService, googleCalendarService *GoogleCalendarService, oauthConfig *oauth2.Config, eventsService EventsService) *WhatsappService {
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
	}
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
	threadID := contact.OpenaiThreadsID
	if threadID == "" {
		threadID, err = service.createNewThread(assistant)
		if err != nil {
			return fmt.Errorf("error creating new thread: %v", err)
		}
		contact.NumberPhonesID = numberPhone.ID
		contact.OpenaiThreadsID = threadID
		config.DB.Save(contact)

		// Obtengo el vector_store que usa el assistant
		files, err := service.assistantService.serviceFile.GetFileByAssistantID(assistant.ID)
		if err != nil {
			return fmt.Errorf("error GetFileByAssistantID: %v", err)
		}

		var vectorStoreID dtos.FileDto
		if len(files) > 0 {
			vectorStoreID = files[len(files)-1]
		}

		// Asigno el archivo al hilo.
		err = service.openAIAssistantService.EjecutarThread(threadID, []string{vectorStoreID.OpenaiVectorStoreIDs})
		if err != nil {
			return fmt.Errorf("error EjecutarThread: %v", err)
		}
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

	// Formatear la fecha y hora en un formato legible: "31 de enero de 2025 a las 15:04"
	formattedTime := currentTime.Format("02/01/2006 15:04")
	text += fmt.Sprintf("%s\n\n\nFecha y hora actual en Argentina: %s", text, formattedTime)

	// Enviar el mensaje a OpenAI
	response, err := service.InteractWithAssistant(threadID, assistant.OpenaiAssistantsID, text)
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

	// start_date y end_date de Google Calendar
	startDateStr := assistantResp.UserData.StartDate
	endDateStr := assistantResp.UserData.EndDate

	// 3. Lógica según assistantResp.Function
	switch assistantResp.Function {

	case "getEvents":
		if assistant.AccountGoogle {

			// Convertir strings a time.Time
			startDate, err := time.Parse(time.RFC3339, startDateStr)
			if err != nil {
				return fmt.Errorf("error parsing start_date: %v", err)
			}
			endDate, err := time.Parse(time.RFC3339, endDateStr)
			if err != nil {
				return fmt.Errorf("error parsing end_date: %v", err)
			}

			eventsDB, err := service.eventsService.GetAll()
			if err != nil {
				return fmt.Errorf("error parsing end_date: %v", err)
			}

			fmt.Println(eventsDB)

			// Obtener token Google
			token, err := service.googleCalendarService.GetOrRefreshToken(
				int(assistant.ID),
				service.oauthConfig,
				context.Background(),
			)
			if err != nil {
				return err
			}

			// Buscar eventos en el rango
			events, err := service.googleCalendarService.
				FetchGoogleCalendarEventsByDate(token, context.Background(), startDate, endDate)
			if err != nil {
				return err
			}

			// Aquí puedes formatear la respuesta como quieras (por ejemplo, un string con los eventos)
			// Ejemplo: responseUser = "Los eventos encontrados son: ..."
			// En este ejemplo, simplemente los imprimimos.
			fmt.Println(events)
			responseUser = fmt.Sprintf("Eventos entre %s y %s: %v", startDateStr, endDateStr, events)
		} else {
			// Si no hay cuenta Google, devolvemos un mensaje por defecto
			responseUser = "Lo siento, no tengo acceso a Google Calendar."
		}

	case "insertEvents":
		// Solo ejecutamos si assistant.AccountGoogle == true
		if assistant.AccountGoogle {

			token, err := service.googleCalendarService.GetOrRefreshToken(
				int(assistant.ID),
				service.oauthConfig,
				context.Background(),
			)
			if err != nil {
				return err
			}

			// Podrías también usar assistantResp.UserData.Nombre, Email y Phone
			// si los necesitas para el "summary" o "description".
			event, err := googlecalendar.UploadEventRequestCalendar(
				assistantResp.UserData.Nombre,
				"Contacto: "+assistantResp.UserData.Email+"\n Tel: "+assistantResp.UserData.Phone,
				startDateStr,
				endDateStr,
			)
			if err != nil {
				return err
			}

			event, err = service.googleCalendarService.CreateGoogleCalendarEvent(token, context.Background(), event)
			if err != nil {
				return err
			}

			eventDTO := dtos.EventsDto{
				Summary:               assistantResp.UserData.Nombre,
				Description:           "Contacto: " + assistantResp.UserData.Email + "\n Tel: " + assistantResp.UserData.Phone,
				StartDate:             startDateStr,
				EndDate:               endDateStr,
				AssistantsID:          assistant.ID,
				ContactsID:            contact.ID,
				EventGoogleCalendarID: event.Id,
			}

			// Creo el evento en la base de datos
			err = service.eventsService.Create(eventDTO)
			if err != nil {
				log.Printf("\nCreate(eventDTO): %s", err.Error())
				return fmt.Errorf("error creating event: %v", err)
			}

			// Parsear las fechas en el formato esperado
			startDateTime, err := time.Parse("2006-01-02T15:04:05", startDateStr)
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
				"✅ ¡Tu turno ha sido agendado con éxito! 📅\n\n🕒 Inicio: %s \n🕒 Fin: %s.\n\nTe esperamos... ¡Que tengas un excelente día! 😊",
				formattedStart, formattedEnd,
			)

		} else {
			// Si no hay cuenta Google, devolvemos un mensaje por defecto
			responseUser = "Lo siento, no tengo acceso a Google Calendar."
		}

	case "updateEvent":
		if assistant.AccountGoogle {
			// start_date y end_date de Google Calendar
			startDateStr := assistantResp.UserData.StartDate
			endDateStr := assistantResp.UserData.EndDate

			token, err := service.googleCalendarService.GetOrRefreshToken(int(assistant.ID), service.oauthConfig, context.Background())
			if err != nil {
				return err
			}

			event, err := googlecalendar.UploadEventRequestEditCalendar(
				assistantResp.UserData.Nombre,
				"Contacto: "+assistantResp.UserData.Email+"\n Tel: "+assistantResp.UserData.Phone,
				startDateStr,
				endDateStr,
			)
			if err != nil {
				return err
			}

			updatedEvent, err := service.googleCalendarService.UpdateGoogleCalendarEvent(token, context.Background(), "asdfasdf", event)
			if err != nil {
				return err
			}
			fmt.Print(updatedEvent)

		}
		eventDTO := dtos.EventsDto{
			ID:           0,
			Summary:      assistantResp.UserData.Nombre,
			Description:  "Contacto: " + assistantResp.UserData.Email + "\n Tel: " + assistantResp.UserData.Phone,
			StartDate:    startDateStr,
			EndDate:      endDateStr,
			AssistantsID: assistant.ID,
			ContactsID:   contact.ID,
		}
		err = service.eventsService.Update(eventDTO)
		if err != nil {
			return err
		}

	case "deleteEvent":
		if assistant.AccountGoogle {
			token, err := service.googleCalendarService.GetOrRefreshToken(int(assistant.ID), service.oauthConfig, context.Background())
			if err != nil {
				return err
			}
			err = service.googleCalendarService.DeleteGoogleCalendarEvent(token, context.Background(), "")
			if err != nil {
				return err
			}

			responseUser = "✅ Su turno ha sido eliminado con éxito. Si nesesitas cualquier otra cosa, estoy acá para ayudarte 😊"
		}

	case "responseText":
		// Respuesta normal (sin Calendar)
		// El texto a mostrar al usuario viene en assistantResp.Message
		responseUser = assistantResp.Message

	default:
		// Acción desconocida
		return fmt.Errorf("unknown action in response: %s", assistantResp.Function)
	}

	// 4. Guardar la respuesta que le daremos al usuario
	err = saveMessageWithUniqueID(service, int(numberPhone.NumberPhone), int(contact.ID), responseUser)
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

func (s *WhatsappService) InteractWithAssistant(threadID, assistantID, message string) (string, error) {
	// Crear un run para el thread con la conversación completa
	err := s.openAIAssistantService.SendMessageToThread(threadID, message, true)
	if err != nil {
		return "", fmt.Errorf("error creating message with conversation: %v", err)
	}

	runID, err := s.openAIAssistantService.CreateRunForThreadWithConversation(threadID, assistantID, nil)
	if err != nil {
		return "", fmt.Errorf("error creating run with conversation: %v", err)
	}
	fmt.Printf("Run created: %s\n", runID)

	// Esperar a que el run esté completado
	err = s.openAIAssistantService.WaitForRunCompletion(threadID, runID, 30, 2*time.Second)
	if err != nil {
		return "", fmt.Errorf("error waiting for run completion: %v", err)
	}

	// Obtener los mensajes del thread y encontrar la respuesta del asistente
	response, err := s.openAIAssistantService.GetMessagesFromThread(threadID)
	if err != nil {
		return "", fmt.Errorf("error getting messages from thread: %v", err)
	}

	return response, nil
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

			messagePromt := "Quiero que revises este hilo de conversación y, dentro de un lapso de 6 horas atrás, identifiques si algún usuario acordó una reunión. Para que sea válida, debe haber proporcionado explícitamente los siguientes datos: su nombre y correo electrónico, siendo opcional el número de celular.\nSi encuentras estos datos, respóndeme únicamente en el siguiente formato y nada más:\n\n[Nombre del usuario]\n[Correo del usuario]\n[Número de celular] (indica 'No proporcionado' si no lo dio). SOLO LOS DATOS, NO EL NOMBRE DE CADA DATO.\n\nSi no encuentras ninguna reunión dentro del lapso indicado o si los datos no están completos, responde únicamente con 'No se encontró información de reuniones.' No proporciones nada adicional ni interpretes respuestas parciales."
			response, err := s.InteractWithAssistant(contact.OpenaiThreadsID, number.Assistant.OpenaiAssistantsID, messagePromt)
			if err != nil {
				return fmt.Errorf("error interacting with assistant: %v", err)
			}

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

func (service *WhatsappService) createNewThread(assistant dtos.AssistantDto) (string, error) {
	// Crear un nuevo Thread en OpenAI
	threadID, err := service.openAIAssistantService.CreateThread(assistant.Model, assistant.Instructions)
	if err != nil {
		return "", fmt.Errorf("error creating thread: %v", err)
	}
	return threadID, nil
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
