package controllers

import (
	"fmt"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/openaiassistantdtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/openaiassistantdtos/openaivectorfiles"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos/openaiassistantdtos/openaivectorsdtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2"
)

// OpenIAAssistantController struct para manejar los endpoints del asistente de OpenAI.
type OpenIAAssistantController struct {
	service *services.OpenIAAssistantService
}

// NewOpenIAAssistantController crea una nueva instancia de OpenIAAssistantController.
func NewOpenIAAssistantController(service *services.OpenIAAssistantService) *OpenIAAssistantController {
	return &OpenIAAssistantController{service: service}
}

// PostCreateAssistant maneja la creación de un asistente.
func (controller *OpenIAAssistantController) PostCreateAssistant(c *fiber.Ctx) error {
	var request openaiassistantdtos.AssistantWithFileSearchRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  false,
			"message": "Invalid request format",
			"data":    nil,
		})
	}

	assistant, err := controller.service.CreateAssistant(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  false,
			"message": "Error creating assistant",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  true,
		"message": "Assistant created successfully",
		"data":    assistant,
	})
}

// PostSaveFileToOpenAI guarda un archivo en openAI. Generalmente luego se utiliza para asociar a un vector, el asistente nesesita sacar informacion de algun lado si le queremos dar de donde sacar.
func (controller *OpenIAAssistantController) PostSaveFileToOpenAI(c *fiber.Ctx) error {

	// Recibir el archivo .txt
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  false,
			"message": "No file was uploaded",
			"data":    nil,
		})
	}

	// Guardar el archivo en la ubicación `pkg/prompts`
	filePath := fmt.Sprintf("../pkg/files/%s", file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  false,
			"message": fmt.Sprintf("Error saving file: %s", err.Error()),
			"data":    nil,
		})
	}

	query := openaiassistantdtos.AssistantUploadFileQuery{
		Purpose:  "fine-tune",
		PathFile: filePath,
	}

	// Crear un cliente para interactuar con la API de OpenAI
	assistantResponse, err := controller.service.PostUploadFileForAssistant(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  false,
			"message": fmt.Sprintf("Error generating response from assistant: %s", err.Error()),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  true,
		"message": "Response generated successfully",
		"data":    assistantResponse,
	})
}

// PostCreateVectorStorage crea un vector storage para la cuenta de open ai principal (usa el token de variable de entorno).
func (controller *OpenIAAssistantController) PostCreateFiles(c *fiber.Ctx) error {
	var request openaiassistantdtos.AssistantUploadFileQuery
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  false,
			"message": "Invalid request format",
			"data":    nil,
		})
	}

	response, err := controller.service.PostUploadFileForAssistant(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  false,
			"message": "Error upload file",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  true,
		"message": "Upload successfully",
		"data":    response,
	})
}

// PostCreateVectorStorage crea un vector storage para la cuenta de open ai principal (usa el token de variable de entorno).
func (controller *OpenIAAssistantController) PostCreateVectorStorageFiles(c *fiber.Ctx) error {
	var request openaivectorfiles.RequestVectorStorage
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  false,
			"message": "Invalid request format",
			"data":    nil,
		})
	}

	response, err := controller.service.AddFileToVectorStoreService(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  false,
			"message": "Error create vector storage",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  true,
		"message": "Create vector storage successfully",
		"data":    response,
	})
}

// PostCreateVectorStorage crea un vector storage para la cuenta de open ai principal (usa el token de variable de entorno).
func (controller *OpenIAAssistantController) GetCreateVectorStorageFiles(c *fiber.Ctx) error {
	var request openaivectorfiles.RequestVectorStoreFiles
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  false,
			"message": "Invalid request format",
			"data":    nil,
		})
	}

	response, err := controller.service.GetVectorStoreFileService(request.VectorStoreID, request.FileID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  false,
			"message": "Error retreive vector file storage",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  true,
		"message": "Retreive vector storage file successfully",
		"data":    response,
	})
}

// PostCreateVectorStorage crea un vector storage para la cuenta de open ai principal (usa el token de variable de entorno).
func (controller *OpenIAAssistantController) PostCreateVectorStorage(c *fiber.Ctx) error {
	var request openaivectorsdtos.CreateVectorStoreRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  false,
			"message": "Invalid request format",
			"data":    nil,
		})
	}

	vectorStorage, err := controller.service.CreateVectorSorage(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  false,
			"message": "Error creating assistant",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  true,
		"message": "Vector storage created successfully",
		"data":    vectorStorage,
	})
}

// PostCreateVectorStorage crea un vector storage para la cuenta de open ai principal (usa el token de variable de entorno).
func (controller *OpenIAAssistantController) GetFile(c *fiber.Ctx) error {
	fileID := c.Query("file_id")
	if fileID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  false,
			"message": "Missing 'assistant_id' parameter",
			"data":    nil,
		})
	}

	fileResponse, err := controller.service.GetFile(fileID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  false,
			"message": "Error creating assistant",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  true,
		"message": "File retrieve successfully",
		"data":    fileResponse,
	})
}
