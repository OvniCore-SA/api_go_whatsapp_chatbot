package services

import (
	"time"

	"golang.org/x/exp/rand"
)

type UtilService struct {
}

func NewUtilService() *UtilService {
	return &UtilService{}
}

func (utilService *UtilService) GetNumberEmoji(number int) string {
	emojis := map[int]string{
		1:  "1️⃣",
		2:  "2️⃣",
		3:  "3️⃣",
		4:  "4️⃣",
		5:  "5️⃣",
		6:  "6️⃣",
		7:  "7️⃣",
		8:  "8️⃣",
		9:  "9️⃣",
		10: "🔟",
		0:  "0️⃣",
	}
	return emojis[number]
}

func (utilService *UtilService) GenerateUniqueCode() string {
	// Caracteres permitidos en el código.
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var codeLength = 6
	rand.Seed(uint64(time.Now().UnixNano())) // Establece la semilla del generador aleatorio.

	// Crear un slice de bytes para almacenar los caracteres del código.
	code := make([]byte, codeLength)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))] // Selecciona un carácter aleatorio de la lista.
	}

	return string(code) // Convierte el slice de bytes a string y lo retorna.
}
