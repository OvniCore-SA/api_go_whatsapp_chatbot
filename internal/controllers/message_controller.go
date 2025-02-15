package controllers

import (
	"strconv"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2"
)

type MessagesController struct {
	service *services.MessagesService
}

func NewMessagesController(service *services.MessagesService) *MessagesController {
	return &MessagesController{service: service}
}

// GetMessagesByNumberPhone - Devuelve los mensajes asociados a un número de teléfono específico con paginación
func (controller *MessagesController) GetMessagesByNumberPhone(c *fiber.Ctx) error {
	numberPhoneID, err := strconv.ParseInt(c.Params("number_phone_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Número de teléfono inválido",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	// Verificar si el numberPhoneID existe en la base de datos
	exists, err := controller.service.DoesNumberPhoneExist(numberPhoneID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al verificar el número de teléfono",
		})
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "El ID de número de teléfono no es válido",
		})
	}

	messages, total, err := controller.service.GetMessagesByNumberPhone(numberPhoneID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Calcular total de páginas
	totalPages := (total + limit - 1) / limit

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"data":    messages,
		"message": "Mensajes obtenidos exitosamente",
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
			"has_next":    page < totalPages,
			"has_prev":    page > 1,
		},
	})
}
