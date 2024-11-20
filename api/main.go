package main

import (
	"log"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/api/middlewares"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/config"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/controllers"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/routes"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/sashabaranov/go-openai"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Configuración de la base de datos
	db, err := config.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Instanciar el cliente de MinIO
	minioClient, err := minio.New(config.MINIO_ENDPOINT, &minio.Options{
		Creds: credentials.NewStaticV4(config.MINIO_ACCESS_KEY, config.MINIO_SECRET_KEY, ""),
		//Secure: config.MINIO_USE_SSL,
	})
	if err != nil {
		log.Fatal("Error initializing MinIO client:", err)
	}

	// Instancio api de OPEN AI
	OpenAIClient := openai.NewClient(config.OPENAI_API_KEY)
	OpenAIAssistantClient := services.NewOpenAIAssistantService(config.OPENAI_API_KEY)

	OpenAIService := services.NewOpenAIService(OpenAIClient)
	OllamaService := services.NewOllamaService()

	// Inicialización de repositorios y servicios

	UtilService := services.NewUtilService()
	UsersRepository := mysql_client.NewUsersRepository(db)
	UsersService := services.NewUsersService(UsersRepository)
	UsersController := controllers.NewUsersController(UsersService)
	PrompsRepository := mysql_client.NewPrompsRepository(db)
	PrompsService := services.NewPrompsService(PrompsRepository)
	PrompsController := controllers.NewPrompsController(PrompsService)
	LogsRepository := mysql_client.NewLogsRepository(db)
	LogsService := services.NewLogsService(LogsRepository)
	LogsController := controllers.NewLogsController(LogsService)
	Password_resetsRepository := mysql_client.NewPasswordResetsRepository(db)
	Password_resetsService := services.NewPassword_resetsService(Password_resetsRepository)
	Password_resetsController := controllers.NewPassword_resetsController(Password_resetsService)
	RolesRepository := mysql_client.NewRolesRepository(db)
	RolesService := services.NewRolesService(RolesRepository)
	RolesController := controllers.NewRolesController(RolesService)
	PermissionsRepository := mysql_client.NewPermissionsRepository(db)
	PermissionsService := services.NewPermissionsService(PermissionsRepository)
	PermissionsController := controllers.NewPermissionsController(PermissionsService)

	// AUTH
	AuthService := services.NewAuthService(UsersService, Password_resetsRepository)
	AuthController := controllers.NewAuthController(AuthService)

	WhatsappService := services.NewWhatsappService(UsersService, PrompsService, LogsService, OpenAIService, UtilService, OllamaService)
	WhatsappController := controllers.NewWhatsappController(WhatsappService)

	BussinessRepository := mysql_client.NewBussinessRepository(db)
	BussinessService := services.NewBussinessService(BussinessRepository)
	BussinessController := controllers.NewBussinessController(BussinessService)

	FileRepository := mysql_client.NewFileRepository(db)
	FileService := services.NewFileService(FileRepository, minioClient)
	FileController := controllers.NewFileController(FileService)

	AssistantRepository := mysql_client.NewAssistantRepository(db)
	AssistantService := services.NewAssistantService(AssistantRepository, FileService, OpenAIAssistantClient)
	AssistantController := controllers.NewAssistantController(AssistantService)

	NumberPhonesRepository := mysql_client.NewNumberPhonesRepository(db)
	NumberPhonesService := services.NewNumberPhonesService(NumberPhonesRepository)
	NumberPhonesController := controllers.NewNumberPhonesController(NumberPhonesService)

	meddlewares := middlewares.MiddlewareManager{}

	// Registrar el middleware para registrar las rutas consultadas
	// Configuración del middleware logger con formato personalizado
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path}\n",
		TimeFormat: "15:04:05", // Formato de hora
		TimeZone:   "Local",
	}))

	// Configuración de TODAS las rutas
	routes.Setup(app, &meddlewares, AuthController, FileController, AssistantController, BussinessController, UsersController, PrompsController, LogsController, Password_resetsController, RolesController, PermissionsController, WhatsappController, NumberPhonesController)

	log.Fatal(app.Listen(":" + config.APP_PORT))
}
