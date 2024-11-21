package controllers

import (
	"strconv"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2"
)

type FileController struct {
	service *services.FileService
}

func NewFileController(service *services.FileService) *FileController {
	return &FileController{service: service}
}

func (controller *FileController) CreateFile(c *fiber.Ctx) error {
	// Obtener assistantsID y purpose desde el formulario
	assistantsID, err := strconv.ParseInt(c.FormValue("assistants_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid assistants_id"})
	}

	purpose := c.FormValue("purpose")

	// Obtener el archivo de la solicitud
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "No file uploaded"})
	}

	// Llamar a CreateFile en el servicio, pasando fileHeader, assistantsID y purpose
	newFile, err := controller.service.CreateFile(fileHeader, assistantsID, purpose, "", "")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error creating file"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": true, "message": "File created successfully", "data": newFile})
}

func (controller *FileController) GetAllFiles(c *fiber.Ctx) error {
	files, err := controller.service.GetAllFiles()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error fetching files"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "data": files})
}

func (controller *FileController) GetFileById(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	file, err := controller.service.GetFileById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "File not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "data": file})
}

func (controller *FileController) UpdateFile(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	assistantsID, _ := strconv.ParseInt(c.FormValue("assistants_id"), 10, 64)
	purpose := c.FormValue("purpose")

	updatedFile, err := controller.service.UpdateFile(id, assistantsID, purpose)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "File not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "message": "File updated successfully", "data": updatedFile})
}

func (controller *FileController) DeleteFile(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	err := controller.service.DeleteFile(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "File not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "message": "File deleted successfully"})
}
