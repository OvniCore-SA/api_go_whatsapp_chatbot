package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	repositories "github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// AuthService estructura para los métodos de autenticación
type AuthService struct {
	userService        *UsersService
	passwordResetsRepo *repositories.PasswordResetsRepository // Si decides crear uno
}

// NewAuthService inicializa un nuevo AuthService
func NewAuthService(userService *UsersService, passwordResetsRepo *repositories.PasswordResetsRepository) *AuthService {
	return &AuthService{
		userService:        userService,
		passwordResetsRepo: passwordResetsRepo,
	}
}

// // Login autentica al usuario y genera un token JWT
func (s *AuthService) Login(email, password string) (string, error) {
	// Recupera el usuario por email
	user, err := s.userService.FindByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Verificar la contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	// Crear el token JWT
	token, err := s.generateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// generateJWT genera el token JWT con los datos del usuario
func (s *AuthService) generateJWT(user dtos.UsersDto) (string, error) {
	permissions := []string{}
	for _, perm := range user.Rol.Permissions {
		permissions = append(permissions, perm.Permission)
	}

	// Configuración del token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":      user.ID,
		"email":       user.Email,
		"role":        user.Rol.Rol,
		"permissions": permissions,
		"exp":         time.Now().Add(time.Hour * 1).Unix(), // Expiración de 1 hora
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// RestorePassword genera un token de reseteo de contraseña y envía el correo electrónico
func (s *AuthService) RestorePassword(email string) error {
	// Buscar el usuario por correo electrónico
	user, err := s.userService.FindByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	// Generar el token de reseteo
	token, err := generateResetToken()
	if err != nil {
		return err
	}

	// Guardar el token de reseteo en la base de datos
	passwordReset := entities.PasswordResets{
		UsersID:   user.ID,
		Token:     token,
		CreatedAt: time.Now(),
	}
	if err := s.passwordResetsRepo.Create(passwordReset); err != nil {
		return err
	}

	// Enviar el email de reseteo de contraseña
	return s.sendResetPasswordEmail(user.Email, token, user.Name)
}

// sendResetPasswordEmail envía un email de reseteo de contraseña
func (s *AuthService) sendResetPasswordEmail(email, token, userName string) error {
	// Verificar que las variables de entorno requeridas estén configuradas
	hostView := os.Getenv("HOST_VIEW")
	from := os.Getenv("USER_EMAIL")
	pass := os.Getenv("PASSWORD_EMAIL")

	if hostView == "" || from == "" || pass == "" {
		return fmt.Errorf("variables de entorno HOST_VIEW, USER_EMAIL o PASSWORD_EMAIL no configuradas")
	}

	// Nueva plantilla de correo
	htmlTemplate := getTemplate(userName, token, hostView)
	// Configuración del email
	to := email
	subject := "Restaurar contraseña"
	headers := map[string]string{
		"From":         from,
		"To":           to,
		"Subject":      subject,
		"Content-Type": "text/html; charset=UTF-8",
	}

	// Construir el mensaje con cabeceras
	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(htmlTemplate)

	// Configuración SMTP (Gmail)
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", from, pass, smtpHost)

	// Enviar el email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg.String()))
	if err != nil {
		return fmt.Errorf("error enviando el correo: %w", err)
	}

	return nil
}

// Devuelve el template para enviar por correo electronico.
func getTemplate(userName, token, hostView string) string {
	htmlTemplate := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="es">
		<head>
			<meta charset="UTF-8" />
			<meta http-equiv="X-UA-Compatible" content="IE=edge" />
			<meta name="viewport" content="width=device-width, initial-scale=1.0" />
			<title>Restaurar Contraseña</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
					color: #333;
					margin: 0;
					padding: 0;
				}
				.container {
					max-width: 600px;
					margin: 0 auto;
					background-color: #fff;
					padding: 20px;
					border-radius: 8px;
					box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
				}
				.header {
					background-color: #8a2be2;
					color: #fff;
					padding: 10px 0;
					text-align: center;
					border-radius: 8px 8px 0 0;
				}
				.content {
					margin: 20px 0;
					text-align: center;
				}
				.button {
					display: inline-block;
					background-color: #6a5acd;
					padding: 15px 25px;
					text-decoration: none;
					border-radius: 5px;
					margin-top: 20px;
					font-size: 16px;
					font-weight: bold;
					color: white !important;
				}
				.button:hover {
					background-color: #6495ed;
				}
				.footer {
					font-size: 12px;
					color: #777;
					text-align: center;
					margin-top: 20px;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>Restaurar tu contraseña</h1>
				</div>
				<div class="content">
					<p>Hola <b>%s</b>,</p>
					<p>Has solicitado restaurar tu contraseña. Utiliza el siguiente enlace para hacerlo:</p>
					<a href="%s/api/reset-password?token=%s" class="button">Restaurar contraseña</a>
					<p>Si no solicitaste este cambio, puedes ignorar este correo.</p>
				</div>
				<div class="footer">
					<p>© 2024 Ovnicore. Todos los derechos reservados.</p>
				</div>
			</div>
		</body>
		</html>
	`, userName, hostView, token)
	return htmlTemplate
}

// ResetPassword restablece la contraseña del usuario
func (s *AuthService) ResetPassword(token, newPassword string) (bool, error) {
	// Busca el registro de reseteo de contraseña por token
	passwordReset, err := s.passwordResetsRepo.FindByToken(token)
	if err != nil {
		return false, errors.New("invalid or expired token")
	}

	// Buscar el usuario por ID y actualizar la contraseña
	user, err := s.userService.GetById(passwordReset.UsersID)
	if err != nil {
		return false, errors.New("user not found")
	}

	// Hashear la nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}

	// Actualizar la contraseña del usuario
	user.Password = string(hashedPassword)
	if err := s.userService.Update(user.ID, user); err != nil {
		return false, err
	}

	// Eliminar el token de reseteo
	if err := s.passwordResetsRepo.DeleteByToken(token); err != nil {
		return false, err
	}

	return true, nil
}

// generateResetToken genera un token seguro aleatorio
func generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
