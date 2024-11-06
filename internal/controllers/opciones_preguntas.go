package controllers

import (
	"strconv"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2"
)

type OpcionesPreguntasController struct {
	service *services.OpcionesPreguntasService
}

func NewOpcionesPreguntasController(service *services.OpcionesPreguntasService) *OpcionesPreguntasController {
	return &OpcionesPreguntasController{service: service}
}

func (controller *OpcionesPreguntasController) GetAll(c *fiber.Ctx) error {
	items, err := controller.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.JSON(items)
}

func (controller *OpcionesPreguntasController) GetById(c *fiber.Ctx) error {
	id := c.Params("id")
	idToInt64, err := strconv.ParseInt(id, 10, 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	item, err := controller.service.GetById(uint64(idToInt64))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.JSON(item)
}

func (controller *OpcionesPreguntasController) Create(c *fiber.Ctx) error {
	var dto dtos.OpcionPreguntasDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	err := controller.service.Create(dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Item created successfully",
	})
}

func (controller *OpcionesPreguntasController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var dto dtos.OpcionPreguntasDto
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

func (controller *OpcionesPreguntasController) Delete(c *fiber.Ctx) error {
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
