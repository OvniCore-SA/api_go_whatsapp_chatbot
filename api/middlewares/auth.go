package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

type MiddlewareManager struct {
	HTTPClient *http.Client
}

// SecureHeadersMiddleware agrega cabeceras de seguridad a todas las respuestas
func (m *MiddlewareManager) SecureHeadersMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains") // HSTS
		c.Set("Content-Security-Policy", "default-src 'self'")                    // Política de contenido
		c.Set("X-Content-Type-Options", "nosniff")                                // Prevenir el sniffing de contenido
		c.Set("X-Frame-Options", "DENY")                                          // Prevenir ataques de clickjacking
		c.Set("X-XSS-Protection", "1; mode=block")                                // Protección contra XSS
		return c.Next()
	}
}

func (m *MiddlewareManager) ValidarApikey() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("ApiKey")

		if apiKey != os.Getenv("API_SYSTEM_KEY") {
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
