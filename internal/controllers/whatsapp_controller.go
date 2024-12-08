package controllers

import (
	"fmt"
	"os"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/whatsapp"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2"
)

type WhatsappController struct {
	service *services.WhatsappService
}

func NewWhatsappController(service *services.WhatsappService) *WhatsappController {
	return &WhatsappController{service: service}
}

// Se encarga de toda la logica de recepcion y envio de mensajes
func (controller *WhatsappController) PostWhatsapp(c *fiber.Ctx) error {

	// Recibir la solicitud y parsear el cuerpo
	var responseComplet whatsapp.ResponseComplet
	if err := c.BodyParser(&responseComplet); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	go func() {

		// Lógica principal delegada al servicio de WhatsApp
		err := controller.service.HandleIncomingMessageWithAssistant(responseComplet)
		if err != nil {
			fmt.Println("ERROR PROCESANDING THE MESSAGE")
			fmt.Println(err.Error())
		}
	}()

	// Respuesta de éxito
	return c.Status(fiber.StatusOK).SendString("Message processed successfully")

}

// Envia un mensaje de whatsapp a un numero específico.
func (controller *WhatsappController) PostSendMessageWhatsapp(c *fiber.Ctx) error {

	go func() {
		// Lógica principal delegada al servicio de WhatsApp

	}()

	// Respuesta de éxito
	return c.Status(fiber.StatusOK).SendString("response")
}

// Se usa para vincular el webhook de la API Whatsapp META
func (controller *WhatsappController) GetWhatsapp(c *fiber.Ctx) error {

	tokenApiWhatsapp := c.Query("hub.verify_token")
	requestHubChallenge := c.Query("hub.challenge")

	var ReqWhats whatsapp.RequestWhatsapp
	if len(tokenApiWhatsapp) > 0 && len(requestHubChallenge) > 0 && tokenApiWhatsapp == os.Getenv("WEBHOOK_TOKEN_WHATSAPP") {
		ReqWhats.Hub_challenge = requestHubChallenge
		ct := c
		ct.Set("Content-Type", "text/plain")
		ct.Status(fiber.StatusOK)

		return c.Send([]byte(ReqWhats.Hub_challenge))
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  false,
			"message": "Error en obtención de parametros.",
			"data":    nil,
		})
	}
}

func (controller *WhatsappController) DemoNotifyInteractions(c *fiber.Ctx) error {

	err := controller.service.NotifyInteractions(6)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  false,
			"message": "Error: " + err.Error(),
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  true,
		"message": "Notificaciones enviadas con éxito.",
		"data":    nil,
	})

}
