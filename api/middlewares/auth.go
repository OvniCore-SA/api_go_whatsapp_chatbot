package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
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
			fmt.Println(os.Getenv("API_SYSTEM_KEY"))
			fmt.Println(apiKey)
			err := fmt.Errorf("acceso denegado o permisos insuficientes")
			return fiber.NewError(403, err.Error())
		}

		return c.Next()
	}
}

// ValidarPermiso verifica que el token incluya el rol y permisos adecuados
func (m *MiddlewareManager) ValidarPermiso(scope string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Obtener el token desde el header Authorization
		bearer := c.Get("Authorization")
		if len(bearer) <= 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "Acceso no autorizado, debe enviar un token de autenticación",
			})
		}

		// Eliminar el prefijo "Bearer "
		tokenString := strings.TrimPrefix(bearer, "Bearer ")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "Token inválido o no proporcionado",
			})
		}

		// Verificar y decodificar el token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil // Usa tu clave secreta aquí
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "Token inválido o expirado",
			})
		}

		// Obtener los claims del token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "No se pudieron obtener los claims del token",
			})
		}

		// Verificar expiración del token
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"status":  false,
					"message": "Token expirado",
				})
			}
		}

		// Obtener roles y permisos del token
		var permissions []string
		rol, ok := claims["role"].(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  false,
				"message": "No tienes permiso para acceder a esta ruta",
			})
		}

		if p, ok := claims["permissions"].([]interface{}); ok {
			for _, permiso := range p {
				permissions = append(permissions, fmt.Sprintf("%v", permiso))
			}
		}

		// Verificar si el usuario tiene el permiso requerido para esta ruta
		if !m.tienePermiso(permissions, scope) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  false,
				"message": "No tienes permiso para acceder a esta ruta",
			})
		}

		// Pasar los claims al contexto para usarlos en el handler si es necesario
		c.Locals("user_id", claims["id"])
		c.Locals("role", rol)
		c.Locals("permissions", permissions)

		return c.Next()
	}
}

// Función auxiliar para verificar permisos
func (m *MiddlewareManager) tienePermiso(permisos []string, permisoRequerido string) bool {
	for _, p := range permisos {
		if p == permisoRequerido {
			return true
		}
	}
	return false
}

func RequestLogger(c *fiber.Ctx) error {
	// Registrar la ruta, método y tiempo de ejecución
	log.Printf("%s: %s ", c.Method(), c.Path())

	// Continuar con la siguiente capa del middleware
	return c.Next()
}
