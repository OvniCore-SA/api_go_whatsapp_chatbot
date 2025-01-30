package openaiassistantdtos

type AssistantJSONResponse struct {
	Function string `json:"function"`
	UserData struct {
		Nombre    string `json:"nombre"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	} `json:"user_data"`
	Message string `json:"message"`
}
