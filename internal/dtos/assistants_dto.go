package dtos

import "time"

type AssistantDto struct {
	ID                 int64            `json:"id"`
	BussinessID        int64            `json:"bussiness_id"`
	Name               string           `json:"name"`
	OpenaiAssistantsID string           `json:"openai_assistants_id,omitempty"`
	Description        string           `json:"description,omitempty"`
	Model              string           `json:"model,omitempty"`
	Instructions       string           `json:"instructions,omitempty"`
	Active             bool             `json:"active"`
	Bussiness          BussinessDto     `json:"bussiness,omitempty"`
	NumberPhones       []NumberPhoneDto `json:"number_phones,omitempty"` // DTO de NumberPhone para la relación de uno a muchos
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
}
