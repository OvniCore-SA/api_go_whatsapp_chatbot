package controllers

import (
	"path/filepath"
	"strconv"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2"
)

type AssistantController struct {
	service *services.AssistantService
}

func NewAssistantController(service *services.AssistantService) *AssistantController {
	return &AssistantController{service: service}
}

// UploadFileToGPT recibe un archivo .jsonl y lo sube a la API de GPT
func (controller *AssistantController) UploadFileToGPT(c *fiber.Ctx) error {
	// Obtener el archivo de la solicitud
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File is required"})
	}

	// Validar que el archivo tenga extensión .jsonl
	if filepath.Ext(file.Filename) != ".jsonl" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File must be in .jsonl format"})
	}

	// Abrir el archivo
	fileContent, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to open file"})
	}
	defer fileContent.Close()

	// Subir el archivo a la API de GPT
	fileID, err := controller.service.UploadFileToGPT(fileContent, file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error uploading file to GPT", "details": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "File uploaded successfully", "file_id": fileID})
}

// Agregar un nuevo asistente
func (controller *AssistantController) AddAssistant(c *fiber.Ctx) error {
	var assistantDto dtos.AssistantDto
	if err := c.BodyParser(&assistantDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Obtener el archivo de la solicitud
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File is required"})
	}

	// Validar que el archivo tenga extensión .jsonl o .txt
	if filepath.Ext(fileHeader.Filename) != ".txt" && filepath.Ext(fileHeader.Filename) != ".pdf" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File must be in .jsonl or .txt format"})
	}

	// Llamar al servicio AssistantService para crear el asistente, pasando el fileHeader
	newAssistant, err := controller.service.CreateAssistantWithFile(assistantDto, fileHeader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating assistant. " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Assistant created successfully", "data": newAssistant})
}

// Obtener todos los asistentes
func (controller *AssistantController) GetAllAssistants(c *fiber.Ctx) error {
	assistants, err := controller.service.FindAllAssistants()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving assistants"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "message": "Assistentes obtenidos con éxito.", "data": assistants})
}

// Obtener un asistente por ID
func (controller *AssistantController) GetAssistant(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	assistant, err := controller.service.FindAssistantById(int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Assistant not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "message": "Assistente obtenido con éxito.", "data": assistant})
}

// Actualizar un asistente
func (controller *AssistantController) UpdateAssistant(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var assistantDto dtos.AssistantDto
	if err := c.BodyParser(&assistantDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	updatedAssistant, err := controller.service.UpdateAssistant(int64(id), assistantDto)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Assistant not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Assistant updated successfully", "data": updatedAssistant, "status": true})
}

// Eliminar un asistente de open ai y de la base de datos(soft delete)
func (controller *AssistantController) DeleteAssistant(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = controller.service.DeleteAssistant(int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Assistant not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "message": "Assistant deleted successfully", "data": id})
}

func (controller *AssistantController) GetAllAssistantsByBussinessId(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Bussiness ID"})
	}

	assistants, err := controller.service.GetAllAssistantsByBussinessId(int64(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving assistants"})
	}
	if len(assistants) == 0 {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"status":  true,
			"message": "Assistants no encontrados.",
			"data":    fiber.Map{"assistants": assistants},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Assistants obtenidos exitosamente.",
		"data":    fiber.Map{"assistants": assistants},
	})
}