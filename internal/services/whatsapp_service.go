package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/config"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/whatsapp"
	metaapi "github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/whatsapp/metaApi"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities/filters"
)

type WhatsappService struct {
	usersService             *UsersService
	prompsService            *PrompsService
	logsService              *LogsService
	opcionesPreguntasService *OpcionesPreguntasService
	metaApps                 *MetaAppsService
	chatbotService           *ChatbotsService
	openAIService            *OpenAIService
	resumesService           *ResumesService
	utilService              *UtilService
	ollamaService            *OllamaService
	userSessions             map[string]*whatsapp.UserSession // Mapa para almacenar sesiones por PhoneNumberID
}

func NewWhatsappService(usersService *UsersService, prompsService *PrompsService, logsService *LogsService, opcionesPreguntasService *OpcionesPreguntasService, metaApps *MetaAppsService, chatbotService *ChatbotsService, openAIService *OpenAIService, utilService *UtilService, resumeService *ResumesService, ollamaService *OllamaService) *WhatsappService {
	return &WhatsappService{
		usersService:             usersService,
		prompsService:            prompsService,
		logsService:              logsService,
		opcionesPreguntasService: opcionesPreguntasService,
		metaApps:                 metaApps,
		chatbotService:           chatbotService,
		openAIService:            openAIService, // Inicialización del servicio de OpenAI
		resumesService:           resumeService,
		utilService:              utilService,
		ollamaService:            ollamaService,
		userSessions:             make(map[string]*whatsapp.UserSession),
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

func (service *WhatsappService) HandleMessageOllama(prompt string) string {
	// Enviar el mensaje recibido a Ollama
	ollamaResponse, err := service.ollamaService.SendMessageToChat(prompt, "Responde brevemente")
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return err.Error()
	}
	return ollamaResponse
}

func (service *WhatsappService) sendMessageBasic(message metaapi.SendMessageBasic, phoneNumberId string, tokenApiWhatsapp string) error {

	reqJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error al convertir el cuerpo de la solicitud a JSON:", err)
		return err
	}

	client := &http.Client{}

	base, err := url.Parse(config.WHATSAPP_URL + "/" + config.WHATSAPP_VERSION + "/" + phoneNumberId + "/messages")
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

	// Obtener el chatbot relacionado
	chatbot, err := service.chatbotService.GetChatbotByPhoneNumberID(phoneNumberID)
	if err != nil {
		return fmt.Errorf("error retrieving chatbot for phoneNumberID %s: %v", phoneNumberID, err)
	}

	if !chatbot.OptionsMenu {
		filterForMetaApp := filters.PrompsFiltro{
			MetaAppsId: chatbot.MetaApps.ID,
		}
		promptForMetaApp, err := service.prompsService.GetByFilter(filterForMetaApp)
		if err != nil {
			return fmt.Errorf("error retrieving metaApps: " + err.Error())
		}
		if !exists || isSessionExpired(session, currentTime) {
			return service.handleNewSessionForChatbotOptionsMenuFalse(phoneNumberID, sender, text, currentTime, &chatbot, &promptForMetaApp)
		}
		return service.handleExistingSessionForChatbotOptionsMenuFalse(sender, text, &chatbot, &promptForMetaApp)
	}

	if !exists || isSessionExpired(session, currentTime) {
		return service.handleNewSession(phoneNumberID, sender, currentTime)
	}

	// Procesar mensaje en sesión activa
	return service.processSessionMessage(session, sender, text, phoneNumberID, currentTime)
}

func (service *WhatsappService) handleExistingSessionForChatbotOptionsMenuFalse(sender string, text string, chatbot *dtos.ChatbotsDto, promps *dtos.PrompsDto) error {
	// Buscar el Resume relacionado al chatbot
	resume, err := service.resumesService.GetResumeByChatbotID(chatbot.ID)
	if err != nil {
		return fmt.Errorf("error retrieving resume for chatbot ID %d: %v", chatbot.ID, err)
	}

	// Crear el mensaje para OpenAI sin saludo
	requestUserAndPrompt := "NO SALUDES AL USUARIO PORQUE NO ES SU PRIMER MENSAJE.\n\nConsulta del usuario: " + text + ".\n" + resume.RequestToResolve

	var openAiResponse string
	if chatbot.GptUse == "CHATGPT" {
		// Enviar el prompt a OpenAI junto con el Resume y la consulta del usuario
		openAiResponse, err = service.openAIService.SendMessageBasicToOpenAI(requestUserAndPrompt, promps.Descripcion)
		if err != nil {
			return fmt.Errorf("error retrieving response from OpenAI: %v", err)
		}
	} else {
		// Enviar el prompt a OpenAI junto con el Resume y la consulta del usuario
		openAiResponse, err = service.ollamaService.SendMessageToChat(requestUserAndPrompt, promps.Descripcion)
		if err != nil {
			return fmt.Errorf("error retrieving response from OpenAI: %v", err)
		}
	}

	response := metaapi.NewSendMessageWhatsappBasic(openAiResponse, sender)

	// Responder al usuario con la respuesta de OpenAI
	err = service.sendMessageBasic(response, chatbot.PhoneNumberId, chatbot.TokenApiWhatsapp)
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
	resume, err := service.resumesService.GetResumeByChatbotID(chatbot.ID)
	if err != nil {
		return fmt.Errorf("error retrieving resume for chatbot ID %d: %v", chatbot.ID, err)
	}

	// Crear el mensaje para OpenAI
	requestUserAndPrompt := "Consulta del usuario: " + text + ".\n" + resume.RequestToResolve

	var openAiResponse string
	if chatbot.GptUse == "CHATGPT" {
		// Enviar el prompt a OpenAI junto con el Resume y la consulta del usuario
		openAiResponse, err = service.openAIService.SendMessageBasicToOpenAI(requestUserAndPrompt, promps.Descripcion)
		if err != nil {
			return fmt.Errorf("error retrieving response from OpenAI: %v", err)
		}
	} else {
		// Enviar el prompt a OpenAI junto con el Resume y la consulta del usuario
		openAiResponse, err = service.ollamaService.SendMessageToChat(requestUserAndPrompt, promps.Descripcion)
		if err != nil {
			return fmt.Errorf("error retrieving response from OpenAI: %v", err)
		}
	}

	response := metaapi.NewSendMessageWhatsappBasic(openAiResponse, sender)

	// Responder al usuario con la respuesta de OpenAI
	err = service.sendMessageBasic(response, chatbot.PhoneNumberId, chatbot.TokenApiWhatsapp)
	if err != nil {
		return fmt.Errorf("error sending OpenAI response to user: %v", err)
	}

	return nil
}

// Extrae la información del mensaje
func extractMessageInfo(value whatsapp.Value) (sender, text, fechaMessage, messageType, phoneNumberID string, err error) {

	fmt.Println("Fecha: " + value.Messages[0].Timestamp)
	sender = value.Messages[0].From
	text = value.Messages[0].Text.Body
	messageType = value.Messages[0].Type
	phoneNumberID = value.Metadata.PhoneNumberID

	// El timestamp en formato Unix
	timestamp, _ := strconv.Atoi(value.Messages[0].Timestamp)

	// Convertir el timestamp a una fecha y hora legible en UTC
	t := time.Unix(int64(timestamp), 0)

	fechaMessage = t.Format("2006-01-02 15:04:05")

	// Verificar si el tiempo del mensaje es posterior al tiempo actual + 5 minutos
	currentTime := time.Now()
	if t.Before(currentTime.Add(-5 * time.Minute)) {
		err = fmt.Errorf("el timestamp del mensaje (%d) excede la fecha y hora actual menos 5 minutos", timestamp)
		return sender, text, "", "", "", err
	}

	return sender, text, fechaMessage, messageType, phoneNumberID, nil
}

// Verifica si la sesión ha expirado
func isSessionExpired(session *whatsapp.UserSession, currentTime time.Time) bool {
	return currentTime.Sub(session.HoraConsulta) > time.Minute*5
}

// Manejar la creación de una nueva sesión y el envío del primer menú
func (service *WhatsappService) handleNewSession(phoneNumberID, sender string, currentTime time.Time) error {
	chatbot, err := service.chatbotService.GetChatbotByPhoneNumberID(phoneNumberID)
	if err != nil {
		return fmt.Errorf("error retrieving chatbot for phoneNumberID %s: %v", phoneNumberID, err)
	}

	// Obtener y enviar el menú inicial
	firstMenu, menuMap := service.getAndSendInitialMenu(&chatbot, phoneNumberID, sender)
	if firstMenu == "" {
		return fmt.Errorf("error retrieving initial menu for phoneNumberID %s", phoneNumberID)
	}

	// Crear o actualizar la sesión
	service.createOrUpdateSession(phoneNumberID, currentTime, menuMap)
	return nil
}

// Obtener y enviar el menú inicial
func (service *WhatsappService) getAndSendInitialMenu(chatbot *dtos.ChatbotsDto, phoneNumberID, sender string) (string, map[int]int) {
	filterOpcionesPreguntas := filters.OpcionPreguntasFiltro{
		OpcionPreguntaID: 0,
		ChatbotsID:       chatbot.ID,
		PrimerMenu:       true,
	}
	firstMenu, menuMap := service.opcionesPreguntasService.GenerarMenuInicial(filterOpcionesPreguntas)

	// Enviar el menú inicial
	message := metaapi.NewSendMessageWhatsappBasic(firstMenu, sender)
	fmt.Println(firstMenu)
	err := service.sendMessageBasic(message, phoneNumberID, chatbot.TokenApiWhatsapp)
	if err != nil {
		fmt.Printf("\nERROR SENDING MESSAGE: %v \n", err)
	}

	return firstMenu, menuMap
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
	chatbot, err := service.chatbotService.GetChatbotByPhoneNumberID(phoneNumberID)
	if err != nil {
		return fmt.Errorf("error retrieving chatbot for phoneNumberID %s: %v", phoneNumberID, err)
	}

	userSelection, err := parseUserSelection(text)
	if err != nil {
		// Si la selección del usuario es inválida, dejamos que OpenAI determine la mejor respuesta
		return service.handleWithOpenAI(session, text, &chatbot, sender, phoneNumberID, currentTime)
	}

	// Si el usuario selecciona una opción válida, manejamos esa opción
	if session.MenuOpciones != nil {
		return service.handleUserSelection(session, userSelection, &chatbot, sender, phoneNumberID, currentTime)
	}

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

	opciones, err := service.getOptionsForSelectedOption(selectedOptionID, chatbot)
	if err != nil {
		return fmt.Errorf("error retrieving options for selected option %d: %v", selectedOptionID, err)
	}

	// Enviar las opciones disponibles al usuario o manejar el final de la conversación
	if len(opciones) > 0 {
		return service.sendOptionsToUser(session, opciones, chatbot, sender, phoneNumberID)
	}

	// Si no hay más opciones, obtener y enviar el contenido real de la opción seleccionada
	_, err = service.opcionesPreguntasService.GetRealContentForOption(uint64(selectedOptionID))
	if err != nil {
		return fmt.Errorf("error retrieving real content for option %d: %v", selectedOptionID, err)
	}

	// Enviar el contenido real al usuario
	err = service.sendRealContentToUser(session, uint64(selectedOptionID), chatbot, sender, phoneNumberID)
	if err != nil {
		return fmt.Errorf("error sending real content to user: %v", err)
	}

	// Marcar como la última opción y finalizar la sesión
	session.EsUltimaOpcion = true
	return nil
}

func (service *WhatsappService) sendRealContentToUser(session *whatsapp.UserSession, optionID uint64, chatbot *dtos.ChatbotsDto, sender, phoneNumberID string) error {
	// Obtener el contenido real para la opción seleccionada
	realContent, err := service.opcionesPreguntasService.GetRealContentForOption(optionID)
	if err != nil {
		return fmt.Errorf("error retrieving real content for option ID %d: %v", optionID, err)
	}

	// Construir y enviar el mensaje con el contenido real al usuario
	message := metaapi.NewSendMessageWhatsappBasic(realContent, sender)
	fmt.Println(realContent)
	err = service.sendMessageBasic(message, phoneNumberID, chatbot.TokenApiWhatsapp)
	if err != nil {
		return fmt.Errorf("error sending real content to user: %v", err)
	}

	// Marcar la sesión como completada y limpiar
	session.EsUltimaOpcion = true
	delete(service.userSessions, phoneNumberID) // Limpiar la sesión después de enviar el contenido final

	return nil
}

// Obtener las opciones relacionadas a la opción seleccionada
func (service *WhatsappService) getOptionsForSelectedOption(optionID int, chatbot *dtos.ChatbotsDto) ([]dtos.OpcionPreguntasDto, error) {
	filterOpcionesPreguntas := filters.OpcionPreguntasFiltro{
		OpcionPreguntaID: int64(optionID),
		ChatbotsID:       chatbot.ID,
	}
	return service.opcionesPreguntasService.ListByIDOpcionPreguntaService(filterOpcionesPreguntas)
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
	filterOpcionesPreguntas := filters.OpcionPreguntasFiltro{
		OpcionPreguntaID: optionID,
		ChatbotsID:       chatbotID,
	}
	opcionPuedeSerUltima, err := service.opcionesPreguntasService.ListByIDOpcionPreguntaService(filterOpcionesPreguntas)
	if err != nil || len(opcionPuedeSerUltima) == 0 {
		return false
	}
	return len(opcionPuedeSerUltima[0].ChildOpcionPreguntas) == 0
}

// Nueva función que maneja la lógica de respuesta utilizando OpenAI
func (service *WhatsappService) handleWithOpenAI(session *whatsapp.UserSession, userMessage string, chatbot *dtos.ChatbotsDto, sender, phoneNumberID string, currentTime time.Time) error {
	// Construir el prompt para OpenAI basado en el mensaje del usuario y las opciones del menú actual desde la base de datos
	optionsPrompt, err := service.buildOptionsPromptFromDB(session)
	if err != nil {
		return fmt.Errorf("error building options prompt from DB: %v", err)
	}

	userMessage += "\n\nNo saludes. No es un primer mensaje."
	// Enviar el prompt a OpenAI y recibir la respuesta
	openAIResponse, err := service.openAIService.SendMessageToOpenAI(userMessage, optionsPrompt)
	if err != nil {
		return fmt.Errorf("error generating response from OpenAI: %v", err)
	}

	// Enviar la respuesta de OpenAI al usuario
	message := metaapi.NewSendMessageWhatsappBasic(openAIResponse, sender)
	err = service.sendMessageBasic(message, phoneNumberID, chatbot.TokenApiWhatsapp)
	if err != nil {
		return fmt.Errorf("error sending OpenAI response to user: %v", err)
	}

	// Actualizar la sesión y el menú
	service.updateSessionWithOpenAIResponse(session, currentTime)

	return nil
}

// Construir el prompt para OpenAI basado en las opciones del menú obtenidas desde la base de datos
func (service *WhatsappService) buildOptionsPromptFromDB(session *whatsapp.UserSession) (string, error) {
	var optionsText string

	for optionNumber, optionID := range session.MenuOpciones {
		// Consultar el nombre de la opción en la base de datos
		opcion, err := service.opcionesPreguntasService.GetById(uint64(optionID))
		if err != nil {
			return "", fmt.Errorf("error retrieving option %d from DB: %v", optionID, err)
		}

		// Agregar la opción al texto del prompt
		optionsText += fmt.Sprintf("%d. %s\n", optionNumber, opcion.Nombre)
	}

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

	// Obtener el chatbot relacionado
	chatbot, err := service.chatbotService.GetChatbotByPhoneNumberID(phoneNumberID)
	if err != nil {
		return fmt.Errorf("error retrieving chatbot for phoneNumberID %s: %v", phoneNumberID, err)
	}

	if len(sender) >= 3 {
		// Remover el tercer dígito del número de teléfono
		sender = sender[:2] + sender[3:]
	}

	base, err := url.Parse(config.WHATSAPP_URL + "/" + config.WHATSAPP_VERSION + "/" + phoneNumberID + "/messages")
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
	req.Header.Set("Authorization", "Bearer "+chatbot.TokenApiWhatsapp)
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
