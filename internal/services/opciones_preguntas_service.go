package services

import (
	"fmt"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities/filters"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
)

type OpcionesPreguntasService struct {
	repository      *mysql_client.OpcionesPreguntasRepository
	chatbotService  ChatbotsService
	prompsService   PrompsService
	metaAppsService MetaAppsService
	openAIService   OpenAIService
	utilService     UtilService
}

func NewOpcionesPreguntasService(repository *mysql_client.OpcionesPreguntasRepository, chatbotService *ChatbotsService, prompsService *PrompsService, metaAppsService *MetaAppsService, openAIService *OpenAIService, utilservice *UtilService) *OpcionesPreguntasService {
	return &OpcionesPreguntasService{
		repository:      repository,
		chatbotService:  *chatbotService,
		prompsService:   *prompsService,
		metaAppsService: *metaAppsService,
		openAIService:   *openAIService,
		utilService:     *utilservice,
	}
}

func (s *OpcionesPreguntasService) GetAll() ([]dtos.OpcionPreguntasDto, error) {
	records, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	dtos := make([]dtos.OpcionPreguntasDto, len(records))
	for i, record := range records {
		dtos[i] = entities.MapEntitiesToOpcionPreguntaDto(record)
	}

	return dtos, nil
}

func (s *OpcionesPreguntasService) GenerarMenuInicial(filter filters.OpcionPreguntasFiltro) (firstMenu string, menuMap map[int]int) {

	chatbot, err := s.chatbotService.repository.FindByID(filter.ChatbotsID)
	if err != nil {
		return "", nil // Si hay un error, retornamos cadenas vacías
	}

	// Traigo el metaApp asociado al chatbot. El metaApp ya tiene cargado a que promt está relacionado.
	metaApp, err := s.metaAppsService.GetById(chatbot.MetaAppsID)
	if err != nil {
		return "", nil // Si hay un error, retornamos cadenas vacías
	}

	// Prompt de bienvenida
	promptWelcome := "Genera un mensaje de bienvenida simple, sencillo, amigable, corto, que no salga de lo formal. No menciones nada robótico. No menciones 'bienvenido', 'te doy la bienvenida' o frases similares; usa frases naturales como: 'Buen día...', 'Hola, ¿cómo estás?', 'Buenas tardes', etc.Puedes usar iconos o emojis para responder. La idea es que el usuario tenga la sensación de estar interactuando con un humano más. IMPORTANTE: Ten en cuenta que este mensaje estará encima de un listado con opciones de temas o menús disponibles para que el usuario elija, así que te sugiero que no hagas preguntas como '¿En qué te puedo ayudar?', en cambio, puedes decir (TE RECOMIENDO): 'selecciona una opción de las siguientes', 'abajo te dejo algunas opciones para que elijas la que te interesa', etc. Nosotros vamos a estar brindando los temas (en ítems) para que el usuario elija uno. Este mensaje se va a usar como mensaje inicial cuando un usuario mande un mensaje al asistente de WhatsApp, por ende, no debe ser largo. Responde solo el texto del mensaje y nada más."

	responseOpenAI, err := s.openAIService.SendMessageBasicToOpenAI(promptWelcome, metaApp.Promps.Descripcion)
	if err != nil {
		fmt.Println("ERROR al intentar mensaje de bienvenida en OpenAI" + err.Error())
		return "", nil // Si hay un error, retornamos cadenas vacías
	}

	firstMenu += responseOpenAI
	firstMenu += "\n\n"

	// Inicializar el mapa de menú
	menuMap = make(map[int]int)

	// Obtener los registros desde el repositorio
	records, err := s.repository.ListByIDOpcionPreguntaRepository(filter)
	if err != nil {
		return "", nil // Si hay un error, retornamos cadenas vacías
	}

	// Concatenar las opciones y mapearlas a sus IDs
	for i, record := range records {
		emoji := s.utilService.GetNumberEmoji(i + 1)
		opcionText := fmt.Sprintf("%s  %s\n", emoji, record.Nombre) // Formato: "1. Texto de la opción"
		firstMenu += opcionText + "\n"                              // Concatenamos las opciones
		menuMap[i+1] = int(record.ID)                               // Mapeamos el índice con la opción
	}

	return firstMenu, menuMap
}

func (s *OpcionesPreguntasService) ListByIDOpcionPreguntaService(filter filters.OpcionPreguntasFiltro) ([]dtos.OpcionPreguntasDto, error) {
	records, err := s.repository.ListByIDOpcionPreguntaRepository(filter)
	if err != nil {
		return nil, err
	}

	dtos := make([]dtos.OpcionPreguntasDto, len(records))
	for i, record := range records {
		dtos[i] = entities.MapEntitiesToOpcionPreguntaDto(record)
	}

	return dtos, nil
}

func (s *OpcionesPreguntasService) GetById(id uint64) (dtos.OpcionPreguntasDto, error) {
	record, err := s.repository.FindByID(id)
	if err != nil {
		return dtos.OpcionPreguntasDto{}, err
	}

	return entities.MapEntitiesToOpcionPreguntaDto(record), nil
}

func (s *OpcionesPreguntasService) GetRealContentForOption(id uint64) (string, error) {
	// Obtener la opción por su ID
	opcion, err := s.GetById(id)
	if err != nil {
		return "", fmt.Errorf("error retrieving option by ID %d: %v", id, err)
	}

	// Aquí se asume que el contenido real que quieres enviar está en el campo "Contenido" de la opción
	if opcion.OpcionPregunta == "" {
		return "", fmt.Errorf("no real content found for option ID %d", id)
	}

	// Retornar el contenido real de la opción
	return opcion.OpcionPregunta, nil
}

func (s *OpcionesPreguntasService) Create(dto dtos.OpcionPreguntasDto) error {
	record := entities.MapDtoToOpcionPreguntas(dto)
	return s.repository.Create(record)
}

func (s *OpcionesPreguntasService) Update(id string, dto dtos.OpcionPreguntasDto) error {
	record := entities.MapDtoToOpcionPreguntas(dto)
	return s.repository.Update(id, record)
}

func (s *OpcionesPreguntasService) Delete(id string) error {
	return s.repository.Delete(id)
}
