package services

import (
	"errors"
	"time"

	"golang.org/x/exp/rand"
)

const (
	FormatTypeDateForBaseDeDatos          DateFormat = "base_de_datos"
	FormatTypeDateForBaseDeDatosSoloFecha DateFormat = "base_de_datos_solo_fecha"
	FormatTypeDateForPersona              DateFormat = "persona"
)

type UtilService struct {
}

func NewUtilService() *UtilService {
	return &UtilService{}
}

type DateFormat string

func (utilService *UtilService) ParseDateString(dateStr string, formatType DateFormat) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	var parsedTime time.Time
	var err error

	// Intentamos parsear con cada formato
	for _, format := range formats {
		parsedTime, err = time.Parse(format, dateStr)
		if err == nil {
			break
		}
	}

	// Si no se pudo parsear, devolvemos un error
	if err != nil {
		return time.Time{}, errors.New("formato de fecha no válido")
	}

	// Definimos los formatos de salida
	switch formatType {
	case FormatTypeDateForBaseDeDatos:
		return parsedTime, nil
	case FormatTypeDateForBaseDeDatosSoloFecha:
		return time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 0, 0, 0, 0, parsedTime.Location()), nil
	case FormatTypeDateForPersona:
		// Si la fecha original no tenía hora, agregamos "00:00:00"
		if len(dateStr) == 10 { // Formato YYYY-MM-DD
			dateStr += " 00:00:00"
		}
		formattedStr := parsedTime.Format("02-01-2006 15:04:05")
		return time.Parse("02-01-2006 15:04:05", formattedStr)
	default:
		return time.Time{}, errors.New("tipo de formato no reconocido")
	}
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
