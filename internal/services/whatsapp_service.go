package services

import (
	"bytes"
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
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/whatsapp"
	metaapi "github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/whatsapp/metaApi"
	whatsappservicedto "github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/whatsapp_service_DTO"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities/filters"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
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
}

func NewWhatsappService(usersService *UsersService, logsService *LogsService, openAIAssistantService *OpenAIAssistantService, utilService *UtilService, numberPhone *NumberPhonesService, messagesRepository *mysql_client.MessagesRepository, assistantService *AssistantService, configurationService *ConfigurationsService) *WhatsappService {
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
	}
}

func (service *WhatsappService) HandleIncomingMessage(response whatsapp.ResponseComplet) error {
	// Procesar cada entrada en el mensaje recibido
	for _, entry := range response.Entry {
		for _, change := range entry.Changes {
			for _, message := range change.Value.Messages {
				// Procesar el mensaje dependiendo de su tipo
				fmt.Printf("Mensaje de tipo: %s", message.Type)
				err := service.processMessage(change.Value)
				if err != nil {
					fmt.Println("ERROR: " + err.Error())
					return err
				}
			}
		}
	}
	return nil
}

func (service *WhatsappService) HandleIncomingMessageWithAssistant(response whatsapp.ResponseComplet) error {
	for _, entry := range response.Entry {
		for _, change := range entry.Changes {
			for _, message := range change.Value.Messages {

				fmt.Print(message.Text.Body)
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
				err = service.handleMessageWithOpenAI(contact, text, numberPhone)
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

func (service *WhatsappService) handleMessageWithOpenAI(contact *entities.Contact, text string, numberPhone *entities.NumberPhone) error {
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
		NumberPhonesID: numberPhone.ID,
		ContactsID:     contact.ID,
		MessageText:    text,
		IsFromBot:      false,
	})
	if err != nil {
		return fmt.Errorf("error saving contact message: %v", err)
	}

	// Enviar el mensaje a OpenAI
	response, err := service.InteractWithAssistant(threadID, assistant.OpenaiAssistantsID, text)
	if err != nil {
		return fmt.Errorf("error sending message to OpenAI: %v", err)
	}

	// Guardar la respuesta del asistente en la base de datos
	err = service.messagesRepository.Create(entities.Message{
		NumberPhonesID: numberPhone.ID,
		ContactsID:     contact.ID,
		MessageText:    response,
		IsFromBot:      true,
	})
	if err != nil {
		return fmt.Errorf("error saving assistant response: %v", err)
	}

	// Enviar la respuesta al usuario
	contactToString := strconv.Itoa(int(contact.NumberPhone))
	message := metaapi.NewSendMessageWhatsappBasic(response, contactToString)
	err = service.sendMessageBasic(message, strconv.FormatInt(numberPhone.WhatsappNumberPhoneId, 10), numberPhone.TokenPermanent)
	if err != nil {
		return fmt.Errorf("error sending response to user: %v", err)
	}

	return nil
}

func (s *WhatsappService) InteractWithAssistant(threadID, assistantID, message string) (string, error) {
	// Crear un run para el thread con la conversación completa
	err := s.openAIAssistantService.SendMessageToThread(threadID, message, true)
	if err != nil {
		return "", fmt.Errorf("error creating run with conversation: %v", err)
	}

	runID, err := s.openAIAssistantService.CreateRunForThreadWithConversation(threadID, assistantID, nil)
	if err != nil {
		return "", fmt.Errorf("error creating run with conversation: %v", err)
	}
	fmt.Printf("Run created: %s\n", runID)

	// Esperar a que el run esté completado
	err = s.openAIAssistantService.WaitForRunCompletion(threadID, runID, 10, 2*time.Second)
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

func (s *WhatsappService) NotifyInteractions(assistantIDDB int64, threadID string) error {
	// Obtener los números de teléfono asociados al asistente
	filter := filters.AssistantsFiltro{
		NumberPhoneToNotifyNotEmpty: true,
	}
	numbers, err := s.numberPhone.GetByFilter(filter)
	if err != nil {
		return fmt.Errorf("error retrieving numbers: %v", err)
	}

	var interactionSummaries []whatsappservicedto.InteractionSummary

	for _, number := range numbers {
		// Filtrar mensajes de las últimas 6 horas
		sixHoursAgo := time.Now().Add(-6 * time.Hour)
		messages, err := s.messagesRepository.GetMessagesByNumber(number.ID, sixHoursAgo)
		if err != nil {
			return fmt.Errorf("error retrieving messages: %v", err)
		}

		// Construir el mensaje para el asistente
		var conversation strings.Builder
		for _, message := range messages {
			quienEnvioMensaje := "usuario"
			if message.IsFromBot {
				quienEnvioMensaje = "assistente"
			}
			conversation.WriteString(fmt.Sprintf("Emisor: %s, Mensaje: %s\n", quienEnvioMensaje, message.MessageText))
		}
		conversation.WriteString("\n¿El usuario envió su teléfono y correo electrónico? Responde solo en el siguiente formato:\n\ntelefono\nemail")

		assistantIDGpt, err := s.configurationService.FindByKey("ASSISTANT_NOTIFY_CLIENTS")
		if err != nil {
			return fmt.Errorf("error retrieving assistant ID: %v", err)
		}
		threadIDGpt, err := s.configurationService.FindByKey("THREAD_NOTIFY_CLIENTS")
		if err != nil {
			return fmt.Errorf("error retrieving thread ID: %v", err)
		}
		fmt.Print(conversation.String())
		// Interactuar con el asistente
		response, err := s.InteractWithAssistant(assistantIDGpt.Value, threadIDGpt.Value, conversation.String())
		if err != nil {
			return fmt.Errorf("error interacting with assistant: %v", err)
		}

		//response := "3872937497\nemanuel.garcia@ovnix.com"

		// Procesar la respuesta
		lines := strings.Split(response, "\n")
		if len(lines) < 2 {
			continue
		}
		contactInfo := whatsappservicedto.UserContactInfo{
			Telefono: lines[0],
			Email:    lines[1],
		}

		// Agregar la interacción al resumen
		interactionSummaries = append(interactionSummaries, whatsappservicedto.InteractionSummary{
			NumberPhoneID: number.ID,
			NumberPhone:   number.NumberPhone,
			Contacts:      []whatsappservicedto.UserContactInfo{contactInfo},
		})
	}

	// Enviar notificaciones a los números designados
	/* 	for _, summary := range interactionSummaries {
		var message strings.Builder
		message.WriteString("Estos usuarios pidieron una reunión:\n")
		for _, contact := range summary.Contacts {
			message.WriteString(fmt.Sprintf("%s, %s\n", contact.Telefono, contact.Email))
		}
		message.WriteString("\nPor favor coordina una reunión a convenir con el usuario.")

		err := s.SendWhatsappNotification(summary.NumberPhone, message.String())
		if err != nil {
			return fmt.Errorf("error sending WhatsApp notification: %v", err)
		}
	} */

	return nil
}

func (s *AssistantService) SendWhatsappNotification(number int64, message string) error {
	// Lógica para enviar notificaciones por WhatsApp
	// Puedes usar Baileys o la API que tengas configurada para este propósito
	fmt.Printf("Enviando mensaje a %d: %s\n", number, message)
	return nil
}

func (s *WhatsappService) getConversationHistory(assistantID, contactID int64) ([]map[string]interface{}, error) {

	messages, err := s.messagesRepository.GetConversation(assistantID, contactID, 10) // Historial de los últimos 5 minutos
	if err != nil {
		return nil, err
	}

	var history []map[string]interface{}
	for _, msg := range messages {
		role := "user"
		if msg.IsFromBot {
			role = "assistant"
		}
		history = append(history, map[string]interface{}{
			"role":    role,
			"content": msg.MessageText,
		})
	}

	return history, nil
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

// parseUserSelection extrae la opción seleccionada del mensaje recibido
func parseUserSelection(messageReceived string) (int, error) {
	// Eliminar espacios en blanco al inicio y final del mensaje
	trimmedMessage := strings.TrimSpace(messageReceived)

	// Intentar convertir el mensaje en un número entero (la opción seleccionada)
	userSelection, err := strconv.Atoi(trimmedMessage)
	if err != nil {
		return 0, fmt.Errorf("el mensaje recibido '%s' no es una opción válida", messageReceived)
	}

	// Retornar el número de la opción seleccionada
	return userSelection, nil
}

func (service *WhatsappService) processMessage(value whatsapp.Value) error {
	// Extraer la información principal del mensaje
	sender, text, fecha, _, phoneNumberID, err := extractMessageInfo(value)
	// Imprimo numero del emisor
	fmt.Println("\nFecha message: " + fecha)
	fmt.Println("\nEmisor: " + sender)
	fmt.Println("\nMessage: " + text)
	if err != nil {
		return err
	}

	currentTime := time.Now()
	// Verificar si existe una sesión para este phoneNumberID
	session, exists := service.userSessions[phoneNumberID]

	if !exists || isSessionExpired(session, currentTime) {
		return service.handleNewSession(phoneNumberID, sender, currentTime)
	}

	// Procesar mensaje en sesión activa
	return service.processSessionMessage(session, sender, text, phoneNumberID, currentTime)
}

func (service *WhatsappService) handleExistingSessionForChatbotOptionsMenuFalse(sender string, text string, chatbot *dtos.ChatbotsDto, promps *dtos.PrompsDto) error {
	// Buscar el Resume relacionado al chatbot

	// Crear el mensaje para OpenAI sin saludo

	var openAiResponse string

	response := metaapi.NewSendMessageWhatsappBasic(openAiResponse, sender)

	// Responder al usuario con la respuesta de OpenAI
	err := service.sendMessageBasic(response, chatbot.PhoneNumberId, chatbot.TokenApiWhatsapp)
	if err != nil {
		return fmt.Errorf("error sending OpenAI response to user: %v", err)
	}

	return nil
}

func (service *WhatsappService) handleNewSessionForChatbotOptionsMenuFalse(phoneNumberID string, sender string, text string, currentTime time.Time, chatbot *dtos.ChatbotsDto, promps *dtos.PrompsDto) error {
	// Crear nueva sesión
	menuMap := make(map[int]int)
	// Crear o actualizar la sesión
	service.createOrUpdateSession(phoneNumberID, currentTime, menuMap)

	// Buscar el Resume relacionado al chatbot

	// Crear el mensaje para OpenAI

	return nil
}

// Extrae la información del mensaje
func extractMessageInfo(value whatsapp.Value) (sender, text, fechaMessage, messageType, phoneNumberID string, err error) {

	fmt.Println("Fecha: " + value.Messages[0].Timestamp)
	sender = value.Messages[0].From
	text = value.Messages[0].Text.Body
	messageType = value.Messages[0].Type
	phoneNumberID = value.Metadata.PhoneNumberID

	// // El timestamp en formato Unix
	// timestamp, _ := strconv.Atoi(value.Messages[0].Timestamp)

	// // Convertir el timestamp a una fecha y hora legible en UTC
	// t := time.Unix(int64(timestamp), 0)

	// fechaMessage = t.Format("2006-01-02 15:04:05")

	// // Verificar si el tiempo del mensaje es posterior al tiempo actual + 5 minutos
	// currentTime := time.Now()
	// if t.Before(currentTime.Add(-5 * time.Minute)) {
	// 	err = fmt.Errorf("el timestamp del mensaje (%d) excede la fecha y hora actual menos 5 minutos", timestamp)
	// 	return sender, text, "", "", phoneNumberID, err
	// }

	return sender, text, fechaMessage, messageType, phoneNumberID, nil
}

// Verifica si la sesión ha expirado
func isSessionExpired(session *whatsapp.UserSession, currentTime time.Time) bool {
	return currentTime.Sub(session.HoraConsulta) > time.Minute*5
}

// Manejar la creación de una nueva sesión y el envío del primer menú
func (service *WhatsappService) handleNewSession(phoneNumberID, sender string, currentTime time.Time) error {

	// Obtener y enviar el menú inicial

	// Crear o actualizar la sesión

	return nil
}

// Obtener y enviar el menú inicial
func (service *WhatsappService) getAndSendInitialMenu(chatbot *dtos.ChatbotsDto, phoneNumberID, sender string) (string, map[int]int) {

	// Enviar el menú inicial
	return "", nil
}

// Crear o actualizar la sesión
func (service *WhatsappService) createOrUpdateSession(phoneNumberID string, currentTime time.Time, menuMap map[int]int) {
	service.userSessions[phoneNumberID] = &whatsapp.UserSession{
		Opcion:         0,
		HoraConsulta:   currentTime,
		MenuEnviado:    1,
		EsUltimaOpcion: false,
		MenuOpciones:   menuMap, // Guardamos el mapa de menú en la sesión
	}
}

// Procesar el mensaje cuando ya existe una sesión
func (service *WhatsappService) processSessionMessage(session *whatsapp.UserSession, sender, text, phoneNumberID string, currentTime time.Time) error {

	return nil
}

// Maneja una selección de usuario no válida y envía el menú inicial de nuevo
func (service *WhatsappService) handleInvalidUserSelection(chatbot *dtos.ChatbotsDto, phoneNumberID, sender string, currentTime time.Time) error {
	firstMenu, menuMap := service.getAndSendInitialMenu(chatbot, phoneNumberID, sender)
	if firstMenu == "" {
		return fmt.Errorf("error retrieving initial menu for phoneNumberID %s", phoneNumberID)
	}

	service.createOrUpdateSession(phoneNumberID, currentTime, menuMap)
	return nil
}

// Manejar la selección del usuario y enviar el nuevo menú o finalizar la conversación
func (service *WhatsappService) handleUserSelection(session *whatsapp.UserSession, userSelection int, chatbot *dtos.ChatbotsDto, sender, phoneNumberID string, currentTime time.Time) error {
	selectedOptionID, exists := session.MenuOpciones[userSelection]
	if !exists {
		return service.handleInvalidUserSelection(chatbot, phoneNumberID, sender, currentTime)
	}

	// Si no hay más opciones, obtener y enviar el contenido real de la opción seleccionada

	// Enviar el contenido real al usuario
	err := service.sendRealContentToUser(session, uint64(selectedOptionID), chatbot, sender, phoneNumberID)
	if err != nil {
		return fmt.Errorf("error sending real content to user: %v", err)
	}

	// Marcar como la última opción y finalizar la sesión
	session.EsUltimaOpcion = true
	return nil
}

func (service *WhatsappService) sendRealContentToUser(session *whatsapp.UserSession, optionID uint64, chatbot *dtos.ChatbotsDto, sender, phoneNumberID string) error {
	// Obtener el contenido real para la opción seleccionada

	// Construir y enviar el mensaje con el contenido real al usuario

	// Marcar la sesión como completada y limpiar
	session.EsUltimaOpcion = true
	delete(service.userSessions, phoneNumberID) // Limpiar la sesión después de enviar el contenido final

	return nil
}

// Enviar las opciones disponibles al usuario
func (service *WhatsappService) sendOptionsToUser(session *whatsapp.UserSession, opciones []dtos.OpcionPreguntasDto, chatbot *dtos.ChatbotsDto, sender, phoneNumberID string) error {
	var messageText string
	session.MenuOpciones = make(map[int]int) // Reiniciar el mapa para el nuevo conjunto de opciones

	for i, opcion := range opciones {
		optionNumber := i + 1
		emojiOption := service.utilService.GetNumberEmoji(i + 1)
		messageText += fmt.Sprintf("%s %s \n\n", emojiOption, opcion.Nombre)
		session.MenuOpciones[optionNumber] = int(opcion.ID)

		if service.isLastOption(opcion.ID, chatbot.ID) {
			delete(service.userSessions, phoneNumberID)
		}
	}

	message := metaapi.NewSendMessageWhatsappBasic(messageText, sender)
	fmt.Println(messageText)
	err := service.sendMessageBasic(message, phoneNumberID, chatbot.TokenApiWhatsapp)
	if err != nil {
		fmt.Printf("\nERROR SENDING MESSAGE: %v \n", err)
	}

	return nil
}

// Verifica si una opción es la última opción
func (service *WhatsappService) isLastOption(optionID int64, chatbotID int64) bool {

	return false
}

// Construir el prompt para OpenAI basado en las opciones del menú obtenidas desde la base de datos
func (service *WhatsappService) buildOptionsPromptFromDB(session *whatsapp.UserSession) (string, error) {
	var optionsText string

	return optionsText, nil
}

// Actualizar la sesión basada en la respuesta de OpenAI
func (service *WhatsappService) updateSessionWithOpenAIResponse(session *whatsapp.UserSession, currentTime time.Time) {
	// Aquí puedes actualizar el estado de la sesión o el menú en función de la respuesta de OpenAI
	session.HoraConsulta = currentTime
}

func (service *WhatsappService) SendTypingIndicator(responseComplet whatsapp.ResponseComplet) error {

	sender := ""
	phoneNumberID := ""
	if len(responseComplet.Entry) > 0 && len(responseComplet.Entry[0].Changes) > 0 && len(responseComplet.Entry[0].Changes[0].Value.Messages) > 0 {
		sender = responseComplet.Entry[0].Changes[0].Value.Messages[0].From
		phoneNumberID = responseComplet.Entry[0].Changes[0].Value.Metadata.PhoneNumberID
	} else {
		return nil
	}

	if len(sender) >= 3 {
		// Remover el tercer dígito del número de teléfono
		sender = sender[:2] + sender[3:]
	}

	base, err := url.Parse(os.Getenv("WHATSAPP_URL") + "/" + os.Getenv("WHATSAPP_VERSION") + "/" + phoneNumberID + "/messages")
	if err != nil {
		return fmt.Errorf("error parcer url: %v", err)
	}

	// Construcción del cuerpo de la solicitud para enviar el evento "escribiendo"
	body := map[string]interface{}{
		"recipient_type": "individual",
		"to":             sender, // Número del destinatario
		"type":           "typing",
		"typing": map[string]interface{}{
			"status": "active",
		},
	}

	// Serialización del cuerpo a JSON
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error serializing body: %v", err)
	}

	// Creación de la solicitud HTTP
	req, err := http.NewRequest("POST", base.String(), bytes.NewBuffer(bodyJSON))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	// Cliente HTTP
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	// Ejecución de la solicitud
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Verificación del estado HTTP
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
