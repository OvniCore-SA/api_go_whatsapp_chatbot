package controllers

import (
	"strconv"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2"
)

type BussinessController struct {
	service *services.BussinessService
}

func NewBussinessController(service *services.BussinessService) *BussinessController {
	return &BussinessController{service: service}
}

// Crear un nuevo negocio
func (controller *BussinessController) CreateBussiness(c *fiber.Ctx) error {
	var bussinessDto dtos.BussinessDto
	if err := c.BodyParser(&bussinessDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Error al analizar la solicitud",
			"error":   err.Error(),
		})
	}

	bussiness, err := controller.service.CreateBussiness(bussinessDto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Error al crear el negocio",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  true,
		"message": "Negocio creado con éxito",
		"data":    bussiness,
	})
}

// Obtener todos los negocios
func (controller *BussinessController) GetAllBussinesses(c *fiber.Ctx) error {
	bussinesses, err := controller.service.GetAllBussinesses()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Error al obtener los negocios",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Negocios obtenidos con éxito",
		"data":    bussinesses,
	})
}

// Obtener un negocio por ID
func (controller *BussinessController) GetBussinessById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "ID inválido",
		})
	}

	bussiness, err := controller.service.GetBussinessById(int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Negocio no encontrado",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Negocio obtenido con éxito",
		"data":    bussiness,
	})
}

// Actualizar un negocio
func (controller *BussinessController) UpdateBussiness(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "ID inválido",
		})
	}

	var bussinessDto dtos.BussinessDto
	if err := c.BodyParser(&bussinessDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Error al analizar la solicitud",
			"error":   err.Error(),
		})
	}

	updatedBussiness, err := controller.service.UpdateBussiness(int64(id), bussinessDto)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Negocio no encontrado",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Negocio actualizado con éxito",
		"data":    updatedBussiness,
	})
}

// Eliminar un negocio
func (controller *BussinessController) DeleteBussiness(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "ID inválido",
		})
	}

	if err := controller.service.DeleteBussiness(int64(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Negocio no encontrado",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Negocio eliminado con éxito",
	})
}

// Obtener negocios por usuario
func (controller *BussinessController) GetBussinessUser(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("userId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "ID de usuario inválido",
		})
	}

	bussinesses, err := controller.service.FindByUserId(int64(userId))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Negocios no encontrados",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Negocios obtenidos con éxito",
		"data":    bussinesses,
	})
}
