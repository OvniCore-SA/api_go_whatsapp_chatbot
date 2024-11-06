package dtos

type NumberPhoneDto struct {
	ID           int64  `json:"id"`
	AssistantsID int64  `json:"assistants_id"`
	NumberPhone  string `json:"number_phone"`
	UUID         string `json:"uuid"`
	Active       bool   `json:"active"`
}
