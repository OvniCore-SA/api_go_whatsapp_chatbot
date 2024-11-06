package controllers

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2"
)

type ResumesController struct {
	service *services.ResumesService
}

func NewResumesController(service *services.ResumesService) *ResumesController {
	return &ResumesController{service: service}
}

func (controller *ResumesController) GetAll(c *fiber.Ctx) error {
	items, err := controller.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.JSON(items)
}

func (controller *ResumesController) Create(c *fiber.Ctx) error {
	var dto dtos.ResumesDto
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
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success"})
}

func (controller *ResumesController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	item, err := controller.service.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.JSON(item)
}

func (controller *ResumesController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var dto dtos.ResumesDto
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
	return c.JSON(fiber.Map{"status": "success"})
}

func (controller *ResumesController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	err := controller.service.Delete(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"status": "success"})
}
