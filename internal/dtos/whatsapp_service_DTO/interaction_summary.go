package whatsappservicedto

import "github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"

type UserContactInfo struct {
	Nombre    string
	Telefono  string
	Email     string
	FechaHora string
}

type InteractionSummary struct {
	NumberPhoneEntity entities.NumberPhone
	NumberPhoneID     int64
	NumberPhone       int64
	Contacts          []UserContactInfo
}
