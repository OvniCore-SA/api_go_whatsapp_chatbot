package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/api/middlewares"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/config"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/controllers"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/routes"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	// Configurar zona horaria global
	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		log.Fatalf("Error cargando la zona horaria: %v", err)
	}
	time.Local = loc

	// Cargar las variables del archivo .env
	err = godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error cargando el archivo .env")
		return
	}

	// Cargar configuración de OAuth
	OauthConfig := config.LoadOAuthConfig()

	app := fiber.New()

	// Configuración de la base de datos
	db, err := config.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Instanciar el cliente de MinIO
	minioClient, err := minio.New(os.Getenv("MINIO_ENDPOINT"), &minio.Options{
		Creds: credentials.NewStaticV4(os.Getenv("MINIO_ACCESS_KEY"), os.Getenv("MINIO_SECRET_KEY"), ""),
		//Secure: config.MINIO_USE_SSL,
	})
	if err != nil {
		log.Fatal("Error initializing MinIO client:", err)
	}

	// Telegram
	InstanceTelegram := services.InstanceTelegram{
		TelegramBotRunning: false,
		InstanceBot:        nil,
		Bot:                nil,
	}
	// Funcion para que se escuchen los eventos de Telegram
	TelegramService := services.NewTelegramService()
	TelegramController := controllers.NewTelegramController(TelegramService, &InstanceTelegram)

	if os.Getenv("HOST_API") == "https://api-botcore.ovnicore.com" {
		go TelegramService.RunTelegramGoRoutine(&InstanceTelegram)
	}

	// Instancio api de OPEN AI
	OpenAIAssistantClient := services.NewOpenAIAssistantService(os.Getenv("OPENAI_API_KEY"))

	// Inicialización de repositorios y servicios
	RolesRepository := mysql_client.NewRolesRepository(db)
	RolesService := services.NewRolesService(RolesRepository)
	UtilService := services.NewUtilService()
	UsersRepository := mysql_client.NewUsersRepository(db)
	UsersService := services.NewUsersService(UsersRepository, RolesService)
	UsersController := controllers.NewUsersController(UsersService)
	LogsRepository := mysql_client.NewLogsRepository(db)
	LogsService := services.NewLogsService(LogsRepository)
	LogsController := controllers.NewLogsController(LogsService)
	Password_resetsRepository := mysql_client.NewPasswordResetsRepository(db)
	Password_resetsService := services.NewPassword_resetsService(Password_resetsRepository)
	Password_resetsController := controllers.NewPassword_resetsController(Password_resetsService)
	MessageRepository := mysql_client.NewMessagesRepository(db)
	MessageService := services.NewMessagesService(MessageRepository)
	MessageController := controllers.NewMessagesController(MessageService)

	ContactRepository := mysql_client.NewContactsRepository(db)
	ContactService := services.NewContactsService(ContactRepository)
	ContactController := controllers.NewContactsController(ContactService)

	RolesController := controllers.NewRolesController(RolesService)
	PermissionsRepository := mysql_client.NewPermissionsRepository(db)
	PermissionsService := services.NewPermissionsService(PermissionsRepository)
	PermissionsController := controllers.NewPermissionsController(PermissionsService)
	NumberPhonesRepository := mysql_client.NewNumberPhonesRepository(db)
	NumberPhonesService := services.NewNumberPhonesService(NumberPhonesRepository)
	NumberPhonesController := controllers.NewNumberPhonesController(NumberPhonesService)
	FileRepository := mysql_client.NewFileRepository(db)
	FileService := services.NewFileService(FileRepository, minioClient)
	FileController := controllers.NewFileController(FileService)
	ConfigurationRepository := mysql_client.NewConfigurationsRepository(db)
	ConfigurationService := services.NewConfigurationsService(ConfigurationRepository)
	AssistantRepository := mysql_client.NewAssistantRepository(db)
	AssistantService := services.NewAssistantService(AssistantRepository, FileService, OpenAIAssistantClient)
	AssistantController := controllers.NewAssistantController(AssistantService)
	EventsRepository := mysql_client.NewEventsRepository(db)
	EventsService := services.NewEventsService(EventsRepository, *UtilService)
	GoogleCalendarRepository := mysql_client.NewGoogleCalendarConfigsRepository(db)
	GoogleCalendarService := services.NewGoogleCalendarService(GoogleCalendarRepository, *AssistantService, EventsService)
	WhatsappService := services.NewWhatsappService(UsersService, LogsService, OpenAIAssistantClient, UtilService, NumberPhonesService, MessageRepository, AssistantService, ConfigurationService, GoogleCalendarService, OauthConfig, EventsService)
	WhatsappController := controllers.NewWhatsappController(WhatsappService)
	BussinessRepository := mysql_client.NewBussinessRepository(db)
	BussinessService := services.NewBussinessService(BussinessRepository)
	BussinessController := controllers.NewBussinessController(BussinessService)

	// Start procesos automaticos
	autoProcess := services.NewAutoProcessService(WhatsappService)
	err = autoProcess.Start()
	if err != nil {
		log.Fatal(err)
	}

	// AUTH
	AuthService := services.NewAuthService(UsersService, Password_resetsRepository)
	AuthController := controllers.NewAuthController(AuthService)

	meddlewares := middlewares.MiddlewareManager{}

	// Registrar el middleware para registrar las rutas consultadas
	// Configuración del middleware logger con formato personalizado
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path}\n",
		TimeFormat: "15:04:05", // Formato de hora
		TimeZone:   "Local",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("HOST_VIEW"),
		AllowHeaders: "",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	// Registrar el middleware de cabeceras seguras
	app.Use(meddlewares.SecureHeadersMiddleware())

	// Configuración de TODAS las rutas
	routes.Setup(app, &meddlewares, AuthController, FileController, AssistantController, BussinessController, UsersController, LogsController, Password_resetsController, RolesController, PermissionsController, WhatsappController, NumberPhonesController, TelegramController, OauthConfig, GoogleCalendarService, MessageController, ContactController, ContactService)

	log.Fatal(app.Listen(":" + os.Getenv("APP_PORT")))
}
