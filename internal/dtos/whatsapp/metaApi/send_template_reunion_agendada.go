package metaapi

type SendMessageTemplate struct {
	MessagingProduct string   `json:"messaging_product"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Template         Template `json:"template"`
}

type Template struct {
	Name       string      `json:"name"`
	Language   Language    `json:"language"`
	Components []Component `json:"components"`
}

type Language struct {
	Code string `json:"code"`
}

type Component struct {
	Type       string      `json:"type"`
	Parameters []Parameter `json:"parameters"`
}

type Parameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func NewBodyWhatsappTemplateCRUD(summary, startTime, endTime, contact, eventCode, numberPhone, templateName string) SendMessageTemplate {
	// Remover tercer dígito del número de teléfono
	if len(numberPhone) >= 3 {
		numberPhone = numberPhone[:2] + numberPhone[3:]
	}

	return SendMessageTemplate{
		MessagingProduct: "whatsapp",
		To:               numberPhone,
		Type:             "template",
		Template: Template{
			Name: templateName,
			Language: Language{
				Code: "es",
			},
			Components: []Component{
				{
					Type: "body",
					Parameters: []Parameter{
						{Type: "text", Text: summary},
						{Type: "text", Text: startTime},
						{Type: "text", Text: endTime},
						{Type: "text", Text: contact},
						{Type: "text", Text: eventCode},
					},
				},
			},
		},
	}
}
