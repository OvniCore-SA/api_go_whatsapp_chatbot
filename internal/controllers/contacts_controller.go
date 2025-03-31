package controllers

import (
	"strconv"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2"
)

type ContactsController struct {
	service *services.ContactsService
}

func NewContactsController(service *services.ContactsService) *ContactsController {
	return &ContactsController{service: service}
}

func (controller *ContactsController) GetMessagesByNumberPhone(c *fiber.Ctx) error {
	numberPhoneID, err := strconv.ParseInt(c.Params("number_phone_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Número de teléfono inválido",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	contacts, total, err := controller.service.GetContactsByNumberPhone(numberPhoneID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Calcular total de páginas
	totalPages := (total + limit - 1) / limit // Redondeo hacia arriba

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"data":    contacts,
		"message": "Contactos obtenidos exitosamente",
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

func (controller *ContactsController) UpdateIsBlocked(c *fiber.Ctx) error {
	// Obtener el ID del número de teléfono desde los parámetros
	contactID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Número de teléfono inválido",
		})
	}

	// Obtener el ID del número de teléfono desde los parámetros
	numberPhoneID, err := strconv.ParseInt(c.Params("number_phone_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Número de teléfono inválido",
		})
	}

	block := c.QueryBool("block", false)

	// Llamar al servicio para actualizar el campo IsBlocked
	err = controller.service.UpdateIsBlocked(contactID, numberPhoneID, block)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "El contacto ha sido bloqueado exitosamente",
		"data":    nil,
	})
}
