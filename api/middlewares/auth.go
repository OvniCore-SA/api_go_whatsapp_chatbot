package middlewares

import (
	"fmt"
	"log"
	"net/http"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/config"
	"github.com/gofiber/fiber/v2"
)

type MiddlewareManager struct {
	HTTPClient *http.Client
}

func (m *MiddlewareManager) ValidarApikey() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("ApiKey")

		if apiKey != config.API_SYSTEM_KEY {
			err := fmt.Errorf("acceso denegado o permisos insuficientes")
			return fiber.NewError(403, err.Error())
		}

		return c.Next()
	}
}

func RequestLogger(c *fiber.Ctx) error {
	// Registrar la ruta, método y tiempo de ejecución
	log.Printf("%s: %s ", c.Method(), c.Path())

	// Continuar con la siguiente capa del middleware
	return c.Next()
}
