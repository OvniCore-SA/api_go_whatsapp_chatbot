package services

import (
	"errors"
	"fmt"
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

func (utilService *UtilService) FormatOpeningDays(openingDays uint8) string {
	// Los días de la semana (0 = Domingo, 1 = Lunes, ..., 6 = Sábado)
	daysOfWeek := []string{"Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado"}

	var availableDays []string
	for i := 0; i < 7; i++ {
		// Si el bit correspondiente está en 1, agregar el día a la lista
		if (openingDays & (1 << i)) != 0 {
			availableDays = append(availableDays, daysOfWeek[i])
		}
	}
	return fmt.Sprintf("Abierto los días: %s", fmt.Sprint(availableDays))
}

func (utilService *UtilService) FormatWorkingHours(workingHours string) string {
	// Se asume que el formato de WorkingHours es "HH:MM-HH:MM"
	return fmt.Sprintf("Horario de trabajo: %s", workingHours)
}

// ConvertDateFormat convierte una fecha en string de cualquier formato válido a otro formato especificado.
func (utilService *UtilService) ConvertDateFormat(dateStr, outputFormat string) (string, error) {
	// Posibles formatos de entrada que puede recibir
	inputFormats := []string{
		"2006-01-02T15:04:05",
		"2006-01-02",
		"2006-01-02 15:04:05",
		"01-02-2006",
		"02-01-2006",
		"02/01/2006",
		"2006/01/02",
		"02-01-2006 15:04:05",
		"01/02/2006 15:04:05",
		"Mon, 02 Jan 2006 15:04:05 MST",
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
		time.RFC3339,
		time.RFC3339Nano,
	}

	var parsedTime time.Time
	var err error

	// Intentar parsear con cada uno de los formatos hasta encontrar uno que funcione
	for _, format := range inputFormats {
		parsedTime, err = time.Parse(format, dateStr)
		if err == nil {
			break
		}
	}

	// Si no se pudo parsear, devolver error
	if err != nil {
		return "", errors.New("formato de fecha no reconocido")
	}

	// Retornar la fecha en el formato deseado
	return parsedTime.Format(outputFormat), nil
}

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
