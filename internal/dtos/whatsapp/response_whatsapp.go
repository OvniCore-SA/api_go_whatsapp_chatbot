package whatsapp

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

type SendWhatsappToUserRequest struct {
	PhoneNumberID string `json:"phone_number_id"`
	To            string `json:"to"`
	Message       string `json:"message"`
}

type ResponseBasicWhatsapp struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberID      string `json:"phone_number_id"`
				} `json:"metadata"`
				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					} `json:"profile"`
					WaID string `json:"wa_id"`
				} `json:"contacts"`
				Messages []struct {
					From      string `json:"from"`
					ID        string `json:"id"`
					Timestamp string `json:"timestamp"`
					Text      struct {
						Body string `json:"body"`
					} `json:"text"`
					Type string `json:"type"`
				} `json:"messages"`
			} `json:"value"`
			Field string `json:"field"`
		} `json:"changes"`
	} `json:"entry"`
}

type ResponseSent struct {
	Statuses []struct {
		ID           string `json:"id"`
		Status       string `json:"status"`
		Timestamp    string `json:"timestamp"`
		RecipientID  string `json:"recipient_id"`
		Conversation struct {
			ID                  string `json:"id"`
			ExpirationTimestamp string `json:"expiration_timestamp"`
			Origin              struct {
				Type string `json:"type"`
			} `json:"origin"`
		} `json:"conversation"`
		Pricing struct {
			Billable     bool   `json:"billable"`
			PricingModel string `json:"pricing_model"`
			Category     string `json:"category"`
		} `json:"pricing"`
	} `json:"statuses"`
}

/* Objeto general que captura las respuestas de la API WPP */
type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type ContactProfile struct {
	Name string `json:"name"`
}

type Contact struct {
	Profile ContactProfile `json:"profile"`
	WAID    string         `json:"wa_id"`
}

type InteractiveListReply struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type InteractiveMessage struct {
	Type      string               `json:"type"`
	ListReply InteractiveListReply `json:"list_reply"`
}

type Message struct {
	Context     map[string]string  `json:"context"`
	From        string             `json:"from"`
	ID          string             `json:"id"`
	Timestamp   string             `json:"timestamp"`
	Type        string             `json:"type"`
	Interactive InteractiveMessage `json:"interactive"`
	Text        Tex                `json:"text"`
}
type Tex struct {
	Body string `json:"body"`
}

type Value struct {
	MessagingProduct string    `json:"messaging_product"`
	Metadata         Metadata  `json:"metadata"`
	Contacts         []Contact `json:"contacts"`
	Messages         []Message `json:"messages"`
}

type Change struct {
	Value Value  `json:"value"`
	Field string `json:"field"`
}

type Entry struct {
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
}

type ResponseComplet struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}
