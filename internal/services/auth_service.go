package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/smtp"
	"os"
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
	// Plantilla de correo
	htmlTemplate := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="es">
		<head>
			<meta charset="UTF-8">
			<title>Restaurar Contraseña</title>
		</head>
		<body>
			<p>Hola <b>%s</b>,</p>
			<p>Utiliza el siguiente enlace para restaurar tu contraseña:</p>
			<p><a href="%s/api/reset-password?token=%s">Restaurar Contraseña</a></p>
			<p>Si no solicitaste este cambio, puedes ignorar este correo.</p>
		</body>
		</html>
	`, userName, os.Getenv("HOST_VIEW"), token)

	// Configuración de email
	from := os.Getenv("USER_EMAIL")
	pass := os.Getenv("PASSWORD_EMAIL")
	to := email

	// Configuración SMTP (Gmail)
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", from, pass, smtpHost)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: Restaurar contraseña\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
		htmlTemplate)

	// Enviar el email
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
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
