package routes

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/api/middlewares"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App,
	middleware *middlewares.MiddlewareManager,
	AuthController *controllers.AuthController,
	FileController *controllers.FileController,
	AssistantController *controllers.AssistantController,
	BussinessController *controllers.BussinessController,
	UsersController *controllers.UsersController,
	PrompsController *controllers.PrompsController,
	LogsController *controllers.LogsController,
	Password_resetsController *controllers.Password_resetsController,
	RolesController *controllers.RolesController,
	PermissionsController *controllers.PermissionsController,
	WhatsappController *controllers.WhatsappController,
	NumberPhonesController *controllers.NumberPhonesController) {

	app.Get("/", middleware.ValidarApikey(), func(c *fiber.Ctx) error {
		return c.Send([]byte("Api chatbot whatsapp by OVNICORE  ®️ "))
	})
	// ROUTE GENERAL ("url_base/api")
	api := app.Group("/api")
	auth := api.Group("/auth")

	api.Post("/assistants/add", AssistantController.AddAssistant)
	api.Get("/assistants/", AssistantController.GetAllAssistants)
	api.Get("/assistants/:id", AssistantController.GetAssistant)
	api.Put("/assistants/:id", AssistantController.UpdateAssistant)
	api.Get("/assistants/getAssistantsByBussiness/:id", AssistantController.GetAllAssistantsByBussinessId)
	api.Delete("/assistants/:id", AssistantController.DeleteAssistant)

	app.Post("/files/create", FileController.CreateFile)
	app.Get("/files/", FileController.GetAllFiles)
	app.Get("/files/:id", FileController.GetFileById)
	app.Put("/files/:id", FileController.UpdateFile)
	app.Delete("/files/:id", FileController.DeleteFile)

	auth.Post("/login", AuthController.Login)
	auth.Post("/restore-password", AuthController.RestorePassword)
	auth.Post("/reset-password", AuthController.ResetPassword)

	api.Get("/users", middleware.ValidarApikey(), UsersController.GetAll)
	api.Get("/users/:id", middleware.ValidarApikey(), UsersController.GetById)
	api.Post("/users", middleware.ValidarApikey(), UsersController.Create)
	api.Put("/users/:id", middleware.ValidarApikey(), UsersController.Update)
	api.Delete("/users/:id", middleware.ValidarApikey(), UsersController.Delete)

	// Rutas de negocios (bussiness)
	api.Post("/bussiness", middleware.ValidarApikey(), BussinessController.CreateBussiness)
	api.Get("/bussiness", middleware.ValidarApikey(), BussinessController.GetAllBussinesses)
	api.Get("/bussinessUser/:userId", middleware.ValidarApikey(), BussinessController.GetBussinessUser)
	api.Get("/bussiness/:id", middleware.ValidarApikey(), BussinessController.GetBussinessById)
	api.Put("/bussiness/:id", middleware.ValidarApikey(), BussinessController.UpdateBussiness)
	api.Delete("/bussiness/:id", middleware.ValidarApikey(), BussinessController.DeleteBussiness)

	api.Get("/webhook", WhatsappController.GetWhatsapp)
	api.Post("/webhook", WhatsappController.PostWhatsapp)
	api.Post("/send-message-basic", middleware.ValidarApikey(), WhatsappController.PostSendMessageWhatsapp)

	api.Get("/logs", middleware.ValidarApikey(), LogsController.GetAll)
	api.Get("/logs/:id", middleware.ValidarApikey(), LogsController.GetById)
	api.Post("/logs", middleware.ValidarApikey(), LogsController.Create)
	api.Put("/logs/:id", middleware.ValidarApikey(), LogsController.Update)
	api.Delete("/logs/:id", middleware.ValidarApikey(), LogsController.Delete)

	api.Get("/password_resets", middleware.ValidarApikey(), Password_resetsController.GetAll)
	api.Get("/password_resets/:id", middleware.ValidarApikey(), Password_resetsController.GetById)
	api.Post("/password_resets", middleware.ValidarApikey(), Password_resetsController.Create)
	api.Put("/password_resets/:id", middleware.ValidarApikey(), Password_resetsController.Update)
	api.Delete("/password_resets/:id", middleware.ValidarApikey(), Password_resetsController.Delete)

	api.Get("/roles", middleware.ValidarApikey(), RolesController.GetAll)
	api.Get("/roles/:id", middleware.ValidarApikey(), RolesController.GetById)
	api.Post("/roles", middleware.ValidarApikey(), RolesController.Create)
	api.Put("/roles/:id", middleware.ValidarApikey(), RolesController.Update)
	api.Delete("/roles/:id", middleware.ValidarApikey(), RolesController.Delete)

	api.Get("/Permissions", middleware.ValidarApikey(), PermissionsController.GetAll)
	api.Get("/Permissions/:id", middleware.ValidarApikey(), PermissionsController.GetById)
	api.Post("/Permissions", middleware.ValidarApikey(), PermissionsController.Create)
	api.Put("/Permissions/:id", middleware.ValidarApikey(), PermissionsController.Update)
	api.Delete("/Permissions/:id", middleware.ValidarApikey(), PermissionsController.Delete)

	api.Get("/Permissions", middleware.ValidarApikey(), PermissionsController.GetAll)
	api.Get("/Permissions/:id", middleware.ValidarApikey(), PermissionsController.GetById)
	api.Post("/Permissions", middleware.ValidarApikey(), PermissionsController.Create)
	api.Put("/Permissions/:id", middleware.ValidarApikey(), PermissionsController.Update)
	api.Delete("/Permissions/:id", middleware.ValidarApikey(), PermissionsController.Delete)

	// Añadir estas rutas en tu archivo principal de rutas:
	api.Get("/number-phones", middleware.ValidarApikey(), NumberPhonesController.GetAll)
	api.Get("/number-phones/:id", middleware.ValidarApikey(), NumberPhonesController.GetById)
	api.Post("/number-phones", middleware.ValidarApikey(), NumberPhonesController.Create)
	api.Put("/number-phones/:id", middleware.ValidarApikey(), NumberPhonesController.Update)
	api.Delete("/number-phones/:id", middleware.ValidarApikey(), NumberPhonesController.Delete)

}
