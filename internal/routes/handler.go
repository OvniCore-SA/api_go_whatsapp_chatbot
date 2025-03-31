package routes

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/api/middlewares"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/controllers"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

func Setup(app *fiber.App,
	middleware *middlewares.MiddlewareManager,
	AuthController *controllers.AuthController,
	FileController *controllers.FileController,
	AssistantController *controllers.AssistantController,
	BussinessController *controllers.BussinessController,
	UsersController *controllers.UsersController,
	LogsController *controllers.LogsController,
	Password_resetsController *controllers.Password_resetsController,
	RolesController *controllers.RolesController,
	PermissionsController *controllers.PermissionsController,
	WhatsappController *controllers.WhatsappController,
	NumberPhonesController *controllers.NumberPhonesController,
	TelegramService *controllers.TelegramController,
	OauthConfig *oauth2.Config,
	GoogleCalendarService *services.GoogleCalendarService,
	MessageController *controllers.MessagesController,
	ContactController *controllers.ContactsController,
	ContactService *services.ContactsService,
	EventController *controllers.EventsController) {

	app.Get("/", middleware.ValidarPermiso("assistants.create"), func(c *fiber.Ctx) error {
		return c.Send([]byte("Api chatbot whatsapp by OVNICORE  ®️ "))
	})
	// ROUTE GENERAL ("url_base/api")
	api := app.Group("/api")
	//auth := api.Group("/auth")

	api.Post("/assistants/add", middleware.ValidarPermiso("assistants.create"), AssistantController.AddAssistant)
	api.Get("/assistants/", middleware.ValidarPermiso("assistants.index"), AssistantController.GetAllAssistants)
	api.Get("/assistants/:id", middleware.ValidarPermiso("assistants.show"), AssistantController.GetAssistant)
	api.Put("/assistants/:id", middleware.ValidarPermiso("assistants.edit"), AssistantController.EditAssistant)
	api.Patch("/assistants/:id", middleware.ValidarPermiso("assistants.edit"), AssistantController.UpdateAssistant)
	api.Get("/assistants/getAssistantsByBussiness/:id", middleware.ValidarPermiso("assistants.show"), AssistantController.GetAllAssistantsByBussinessId)
	api.Delete("/assistants/:id", middleware.ValidarPermiso("assistants.delete"), AssistantController.DeleteAssistant)

	api.Post("/files/create", middleware.ValidarPermiso("assistants.create"), FileController.CreateFile)
	api.Get("/files/", middleware.ValidarPermiso("assistants.index"), FileController.GetAllFiles)
	api.Get("/files/:id", middleware.ValidarPermiso("assistants.show"), FileController.GetFileById)
	api.Put("/files/:id", middleware.ValidarPermiso("assistants.edit"), FileController.UpdateFile)
	api.Delete("/files/:id", middleware.ValidarPermiso("assistants.delete"), FileController.DeleteFile)

	// MESSAGE
	api.Get("/messages/:number_phone_id", middleware.ValidarPermiso("messages.index"), MessageController.GetMessagesByNumberPhone)

	api.Post("/login", AuthController.Login)
	api.Post("/restore-password", AuthController.RestorePassword)
	api.Post("/reset-password", AuthController.ResetPassword)

	// Rutas de autenticación
	app.Get("/auth/url", middleware.ValidarPermiso("assistants.google_account"), controllers.GetAuthURL(OauthConfig))
	app.Get("/auth/callback-auth", controllers.SaveOrUpdateAuthToken(OauthConfig, GoogleCalendarService))
	app.Get("/demo-redirect-url-post-auth", middleware.ValidarPermiso("assistants.create"), controllers.GetRequestDetails())

	// GOOGLE CALENDAR
	api.Get("/calendar/events", middleware.ValidarPermiso("events.index"), controllers.GetCalendarEventsByDate(GoogleCalendarService, OauthConfig))
	api.Post("/calendar/events", middleware.ValidarPermiso("events.create"), controllers.AddCalendarEvent(GoogleCalendarService, OauthConfig, ContactService))
	api.Put("/calendar/events/:event_id", middleware.ValidarPermiso("events.edit"), controllers.UpdateCalendarEvent(GoogleCalendarService, OauthConfig))
	api.Delete("/calendar/events/:event_id", middleware.ValidarPermiso("events.delete"), controllers.DeleteCalendarEvent(GoogleCalendarService, OauthConfig))

	// Events (BOT-CORE)
	api.Post("/events", middleware.ValidarPermiso("events.create"), EventController.CreateEvent)                                                             // Crear un evento
	api.Get("/events/:id", middleware.ValidarPermiso("events.index"), EventController.GetEventByID)                                                          // Obtener un evento por ID
	api.Get("/events/", middleware.ValidarPermiso("events.index"), EventController.GetAllEvents)                                                             // Obtener todos los eventos
	api.Put("/events/", middleware.ValidarPermiso("events.edit"), EventController.UpdateEvent)                                                               // Actualizar un evento
	api.Delete("/events/:id", middleware.ValidarPermiso("events.delete"), EventController.DeleteEvent)                                                       // Eliminar un evento por ID
	api.Delete("/events/cancel/:codeEvent", middleware.ValidarPermiso("events.delete"), EventController.CancelEvent)                                         // Cancelar un evento por código
	api.Get("/events/contact/:contactID/date/:date/time/:currentTime", middleware.ValidarPermiso("events.index"), EventController.GetEventsByContactAndDate) // Obtener eventos por contacto y fecha

	// CONTACTS
	api.Get("/contacts/number_phone/:number_phone_id", middleware.ValidarPermiso("contacts.index"), ContactController.GetMessagesByNumberPhone)
	api.Patch("/contacts/:id/number_phone/:number_phone_id", middleware.ValidarPermiso("contacts.block"), ContactController.UpdateIsBlocked)

	api.Get("/users", middleware.ValidarPermiso("users.index"), UsersController.GetAll)
	api.Get("/users/:id", middleware.ValidarPermiso("users.show"), UsersController.GetById)
	api.Post("/users", middleware.ValidarPermiso("users.create"), UsersController.Create)
	api.Put("/users/:id", middleware.ValidarPermiso("users.edit"), UsersController.Update)
	api.Delete("/users/:id", middleware.ValidarPermiso("users.delete"), UsersController.Delete)

	// Rutas de negocios (bussiness)
	api.Post("/bussiness", middleware.ValidarPermiso("bussiness.create"), BussinessController.CreateBussiness)
	api.Get("/bussiness", middleware.ValidarPermiso("bussiness.index"), BussinessController.GetAllBussinesses)
	api.Get("/bussinessUser/:userId", middleware.ValidarPermiso("bussiness.show"), BussinessController.GetBussinessUser)
	api.Get("/bussiness/:id", middleware.ValidarPermiso("bussiness.show"), BussinessController.GetBussinessById)
	api.Put("/bussiness/:id", middleware.ValidarPermiso("bussiness.edit"), BussinessController.UpdateBussiness)
	api.Delete("/bussiness/:id", middleware.ValidarPermiso("bussiness.delete"), BussinessController.DeleteBussiness)

	api.Get("/webhook", WhatsappController.GetWhatsapp)
	api.Post("/webhook", WhatsappController.PostWhatsapp)
	api.Post("/notificar-datos-clientes", middleware.ValidarPermiso("events.index"), WhatsappController.DemoNotifyInteractions)
	api.Post("/send-message-basic", middleware.ValidarPermiso("whatsapp.send_message"), WhatsappController.PostSendMessageWhatsapp)
	api.Post("/send-message-template", WhatsappController.DemoFunctionWhatsappController)

	api.Post("/telegram/send-message", middleware.ValidarPermiso("telegram.send_message"), TelegramService.SendMessageBasic)

	api.Get("/logs", middleware.ValidarPermiso("logs.index"), LogsController.GetAll)
	api.Get("/logs/:id", middleware.ValidarPermiso("logs.index"), LogsController.GetById)
	api.Post("/logs", middleware.ValidarPermiso("logs.index"), LogsController.Create)
	api.Put("/logs/:id", middleware.ValidarPermiso("logs.index"), LogsController.Update)
	api.Delete("/logs/:id", middleware.ValidarPermiso("logs.index"), LogsController.Delete)

	api.Get("/password_resets", middleware.ValidarPermiso("events.index"), Password_resetsController.GetAll)
	api.Get("/password_resets/:id", middleware.ValidarPermiso("events.index"), Password_resetsController.GetById)
	api.Post("/password_resets", middleware.ValidarPermiso("events.index"), Password_resetsController.Create)
	api.Put("/password_resets/:id", middleware.ValidarPermiso("events.index"), Password_resetsController.Update)
	api.Delete("/password_resets/:id", middleware.ValidarPermiso("events.index"), Password_resetsController.Delete)

	api.Get("/roles", middleware.ValidarPermiso("roles.index"), RolesController.GetAll)
	api.Get("/roles/:id", middleware.ValidarPermiso("roles.index"), RolesController.GetById)
	api.Post("/roles", middleware.ValidarPermiso("roles.index"), RolesController.Create)
	api.Put("/roles/:id", middleware.ValidarPermiso("roles.index"), RolesController.Update)
	api.Delete("/roles/:id", middleware.ValidarPermiso("roles.index"), RolesController.Delete)

	api.Get("/permissions", middleware.ValidarPermiso("permissions.index"), PermissionsController.GetAll)
	api.Get("/permissions/:id", middleware.ValidarPermiso("permissions.index"), PermissionsController.GetById)
	api.Post("/permissions", middleware.ValidarPermiso("permissions.index"), PermissionsController.Create)
	api.Put("/permissions/:id", middleware.ValidarPermiso("permissions.index"), PermissionsController.Update)
	api.Delete("/permissions/:id", middleware.ValidarPermiso("permissions.index"), PermissionsController.Delete)

	// Añadir estas rutas en tu archivo principal de rutas:
	api.Get("/number-phones", middleware.ValidarPermiso("events.index"), NumberPhonesController.GetAll)
	api.Get("/number-phones/get-by-assistantID/:id", middleware.ValidarPermiso("events.index"), NumberPhonesController.GetAllByAssistantID)

	api.Get("/number-phones/:id", middleware.ValidarPermiso("events.index"), NumberPhonesController.GetById)
	api.Post("/number-phones", middleware.ValidarPermiso("events.index"), NumberPhonesController.Create)
	api.Put("/number-phones/:id", middleware.ValidarPermiso("events.index"), NumberPhonesController.Update)
	api.Delete("/number-phones/:id", middleware.ValidarPermiso("events.index"), NumberPhonesController.Delete)

}
