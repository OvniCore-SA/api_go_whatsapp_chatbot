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
	"golang.org/x/exp/rand"
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

func (service *WhatsappService) HandleIncomingMessageWithAssistant(response whatsapp.ResponseComplet) error {
	for _, entry := range response.Entry {
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

	// Enviar el mensaje a OpenAI
	response, err := service.InteractWithAssistant(threadID, assistant.OpenaiAssistantsID, text)
	if err != nil {
		return fmt.Errorf("error sending message to OpenAI: %v", err)
	}

	// Guardar el mensaje del contacto en la base de datos
	err = saveMessageWithUniqueID(service, int(numberPhone.NumberPhone), int(contact.ID), response)
	if err != nil {
		return fmt.Errorf("error saving contact message: %v", err)
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
			response, err := s.InteractWithAssistant(threadIDGpt.Value, assistantIDGpt.Value, conversation.String())
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

		for _, contact := range summary.Contacts {
			message.WriteString(fmt.Sprintf("📞 *Teléfono:* %s\n✉️ *Correo:* %s\n\n", contact.Telefono, contact.Email))
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
