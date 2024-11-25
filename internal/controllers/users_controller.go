package controllers

import (
	"fmt"
	"strconv"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2"
)

type UsersController struct {
	service *services.UsersService
}

func NewUsersController(service *services.UsersService) *UsersController {
	return &UsersController{service: service}
}

func (controller *UsersController) GetAll(c *fiber.Ctx) error {
	items, err := controller.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Usuarios obtenidos exitosamente.",
		"data":    items,
	})
}

func (controller *UsersController) GetById(c *fiber.Ctx) error {
	idString := c.Params("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		fmt.Println("Error al parcear el id del usuario: " + err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	item, err := controller.service.GetById(int64(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.JSON(item)
}

func (controller *UsersController) Create(c *fiber.Ctx) error {
	var dto dtos.UsersDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	// Validar los datos para creación
	if err := dto.ValidateUsersDto(true); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	err := controller.service.Create(dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Item created successfully",
	})
}

func (controller *UsersController) Update(c *fiber.Ctx) error {

	var dto dtos.UsersDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	idString := c.Params("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		fmt.Println("Error al parcear el id del usuario: " + err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	dto.ID = int64(id)
	// Validar los datos para creación
	if err := dto.ValidateUsersDto(false); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	err = controller.service.Update(int64(id), dto)
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

func (controller *UsersController) Delete(c *fiber.Ctx) error {
	idString := c.Params("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		fmt.Println("Error al parcear el id del usuario: " + err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	err = controller.service.Delete(int64(id))
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
