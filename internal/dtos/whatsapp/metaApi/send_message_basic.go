package metaapi

type SendMessageBasic struct {
	MessagingProduct string `json:"messaging_product"`
	RecipientType    string `json:"recipient_type"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Text             Texto  `json:"text"`
}
type Texto struct {
	PreviewURL bool   `json:"preview_url"`
	Body       string `json:"body"`
}

// Se encarga de devolver un body listo para ser enviado por api de whatsapp (Se quita el tercer dígito del numero.)
func NewSendMessageWhatsappBasic(messsage, numberPhone string) (messageBasic SendMessageBasic) {
	// Verificar que el número de teléfono tenga al menos 3 dígitos
	if len(numberPhone) >= 3 {
		// Remover el tercer dígito del número de teléfono
		numberPhone = numberPhone[:2] + numberPhone[3:]
	}

	messageBasic.MessagingProduct = "whatsapp"
	messageBasic.To = numberPhone
	messageBasic.Type = "text"
	messageBasic.RecipientType = "individual"
	messageBasic.Text.PreviewURL = false
	messageBasic.Text.Body = messsage

	return
}
