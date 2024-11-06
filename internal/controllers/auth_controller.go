package controllers

import (
	"log"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}

// Login maneja la lógica de autenticación
func (a *AuthController) Login(c *fiber.Ctx) error {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parsear el cuerpo de la solicitud
	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	// Llama al servicio de autenticación
	token, err := a.service.Login(loginData.Email, loginData.Password)
	if err != nil {
		log.Println("Login error:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"data":    token,
		"message": "Inicio de sesión exitoso.",
	})
}

// RestorePassword maneja el envío del token de reseteo de contraseña
func (a *AuthController) RestorePassword(c *fiber.Ctx) error {
	var requestData struct {
		Email string `json:"email"`
	}

	// Parsear el cuerpo de la solicitud
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	// Llama al servicio para enviar el correo de reseteo
	err := a.service.RestorePassword(requestData.Email)
	if err != nil {
		log.Println("Restore password error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"data":    nil,
		"message": "Password reset email sent.",
	})
}

// ResetPassword maneja el reseteo de contraseña usando un token
func (a *AuthController) ResetPassword(c *fiber.Ctx) error {
	var resetData struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	// Parsear el cuerpo de la solicitud
	if err := c.BodyParser(&resetData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	// Llama al servicio de autenticación para resetear la contraseña
	success, err := a.service.ResetPassword(resetData.Token, resetData.NewPassword)
	if err != nil {
		log.Println("Reset password error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
	}

	if !success {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid or expired token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"data":    nil,
		"message": "Password reset successful.",
	})
}
