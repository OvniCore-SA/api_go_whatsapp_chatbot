package dtos

type NumberPhoneDto struct {
	ID                    int64  `json:"id"`
	AssistantsID          int64  `json:"assistants_id"`
	NumberPhone           int64  `json:"number_phone"`
	UUID                  string `json:"uuid"`
	NumberPhoneToNotify   int64  `json:"number_phone_to_notify"`
	TokenPermanent        string `json:"token_permanent"`
	WhatsappNumberPhoneId int64  `json:"whatsapp_number_phone_id"`
	Active                bool   `json:"active"`
}
