package dtos

import (
	"time"
)

type BussinessDto struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	CuilCuit   string `json:"cuil_cuit,omitempty"`
	WebSite    string `json:"web_site,omitempty"`
	Users      []UsersDto
	Assistants []AssistantDto `json:"assistants,omitempty"` // DTO de asistente para la relación de uno a muchos
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}
