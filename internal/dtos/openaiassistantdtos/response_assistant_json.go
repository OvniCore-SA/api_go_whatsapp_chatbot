package openaiassistantdtos

type AssistantJSONResponse struct {
	Function string `json:"function"` // Nombre de la función a ejecutar, por ejemplo: "getEvents", "insertEvents", "updateEvents", "deleteEvent"
	UserData struct {
		ID               string `json:"id,omitempty"`                 // (Opcional) ID del evento, si es necesario identificarlo
		Nombre           string `json:"nombre,omitempty"`             // Título o nombre del evento
		Email            string `json:"email,omitempty"`              // Email del usuario (para agendar reuniones)
		Phone            string `json:"phone,omitempty"`              // Teléfono del usuario (para agendar reuniones)
		CurrentStartDate string `json:"current_start_date,omitempty"` // Para edición de eventos: la fecha y hora actuales del evento, que se usa para buscarlo en la base de datos
		StartDate        string `json:"start_date,omitempty"`         // Fecha y hora de inicio del evento (nuevo o a registrar)
		EndDate          string `json:"end_date,omitempty"`           // Fecha y hora de fin del evento (se calcula automáticamente sumando 30 minutos a start_date)
		EventCode        string `json:"event_code,omitempty"`         // Codigo de evento, para cuando se nesesite
	} `json:"user_data"`
	Message string `json:"message"` // Mensaje a enviar al usuario, o indicación de error/falta de datos
}
