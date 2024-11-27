package controllers

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2"
)

type RolesController struct {
	service *services.RolesService
}

func NewRolesController(service *services.RolesService) *RolesController {
	return &RolesController{service: service}
}

func (controller *RolesController) GetAll(c *fiber.Ctx) error {
	items, err := controller.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.JSON(items)
}

func (controller *RolesController) GetById(c *fiber.Ctx) error {
	id := c.Params("id")
	item, err := controller.service.GetById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"data":    item,
		"message": "Item created successfully",
	})
}

func (controller *RolesController) Create(c *fiber.Ctx) error {
	var dto dtos.RolesDto
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

func (controller *RolesController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var dto dtos.RolesDto
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
		"data":    dto,
		"message": "Item updated successfully",
	})
}

func (controller *RolesController) Delete(c *fiber.Ctx) error {
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
