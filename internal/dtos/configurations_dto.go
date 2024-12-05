package dtos

import "time"

type ConfigurationDto struct {
	ID          int64     `json:"id"`
	KeyName     string    `json:"key_name"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at,omitempty"`
}
