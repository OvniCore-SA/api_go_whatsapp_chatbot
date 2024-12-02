package controllers

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type NumberPhonesController struct {
	service *services.NumberPhonesService
}

func NewNumberPhonesController(service *services.NumberPhonesService) *NumberPhonesController {
	return &NumberPhonesController{service: service}
}

func (controller *NumberPhonesController) GetAll(c *fiber.Ctx) error {
	items, err := controller.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.JSON(items)
}

func (controller *NumberPhonesController) GetById(c *fiber.Ctx) error {
	id := c.Params("id")
	item, err := controller.service.GetById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Item not found",
		})
	}
	return c.JSON(item)
}

func (controller *NumberPhonesController) Create(c *fiber.Ctx) error {
	var dto dtos.NumberPhoneDto
	// Parseo del cuerpo de la solicitud
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Intentar generar un UUID único
	var uniqueUUID string
	for {
		// Generar un nuevo UUID
		uniqueUUID = uuid.New().String()

		// Verificar si el UUID ya existe en la base de datos
		exists, err := controller.service.UUIDExists(uniqueUUID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		// Si el UUID no existe, salir del bucle
		if !exists {
			break
		}
	}

	// Asignar el UUID único al DTO
	dto.UUID = uniqueUUID

	// Llamada al servicio para crear el número de teléfono
	err := controller.service.Create(dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Number phone created successfully",
	})
}

func (controller *NumberPhonesController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var dto dtos.NumberPhoneDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	err := controller.service.Update(id, dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Item updated successfully",
	})
}

func (controller *NumberPhonesController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	err := controller.service.Delete(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Item deleted successfully",
	})
}
