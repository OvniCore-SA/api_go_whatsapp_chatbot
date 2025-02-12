package openaiassistantdtos

type AssistantJSONResponse struct {
	Function string `json:"function"` // Nombre de la función a ejecutar
	UserData struct {
		// Para "createMeeting"
		UserName    string `json:"user_name,omitempty"`    // Nombre del usuario
		UserEmail   string `json:"user_email,omitempty"`   // Correo del usuario
		MeetingDate string `json:"meeting_date,omitempty"` // Fecha y hora de la reunión en formato ISO 8601
		UserPhone   string `json:"user_phone,omitempty"`   // Telefono del usuario

		// Para "updateEvents"
		// Para "deleteEvent" y "getMeetingDetails"
		EventCode string `json:"event_code,omitempty"` // Código del evento a actualizar
		NewDate   string `json:"new_date,omitempty"`   // Nueva fecha y hora del evento

		// Para "getMeetings"
		DateToSearch string `json:"date_to_search,omitempty"` // Fecha interpretada automáticamente por el asistente
	} `json:"user_data"`

	Message string `json:"message"` // Mensaje de error o solicitud de más información
}
