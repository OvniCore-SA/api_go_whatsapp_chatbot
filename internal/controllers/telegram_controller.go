package controllers

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2"
)

type TelegramController struct {
	service          *services.TelegramService
	instanceTelegram *services.InstanceTelegram
}

func NewTelegramController(service *services.TelegramService, instanceTelegram *services.InstanceTelegram) *TelegramController {
	return &TelegramController{service: service, instanceTelegram: instanceTelegram}
}

// Se encarga de toda la logica de recepcion y envio de mensajes
func (controller *TelegramController) SendMessageBasic(c *fiber.Ctx) error {

	// Recibir la solicitud y parsear el cuerpo
	var sendMessageTelegramRequest services.SendMessageTelegramRequest
	if err := c.BodyParser(&sendMessageTelegramRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Enviar mensaje por telegram en segundo plano
	go controller.service.SendMessageTelegram(sendMessageTelegramRequest, controller.instanceTelegram)

	// Respuesta de éxito
	return c.Status(fiber.StatusOK).SendString("Message processed successfully")

}
