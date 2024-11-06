package dtos

type ChatbotsDto struct {
	ID                int64 // ID como clave primaria
	Name              string
	NumberPhone       string // Número de teléfono
	Apikey            string // API key
	PhoneNumberId     string // ID del número de teléfono
	TokenApiWhatsapp  string // Token de API de Whatsapp
	OptionsMenu       bool
	WhatsappCompanyId string      // ID de la compañía de Whatsapp
	GptUse            string      // Define con que GPT trabaja el chatbot (Disponible: 'CHATGPT', 'OLLAMA')
	CreatedAt         string      // Fecha de creación como string
	UpdatedAt         string      // Fecha de actualización como string
	DeletedAt         string      // Fecha de eliminación suave (soft delete) como string
	MetaApps          MetaAppsDto // Relación de uno a muchos con MetaAppsDto
}
