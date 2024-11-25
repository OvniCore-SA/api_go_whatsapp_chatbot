package dtos

import (
	"errors"
	"strings"
	"time"
)

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

func (dto *AssistantDto) ValidateAssistantDto(isCreate bool) error {
	// Validaciones comunes para creación y edición
	if dto.BussinessID <= 0 {
		return errors.New("bussiness_id es obligatorio y debe ser mayor que 0")
	}

	if strings.TrimSpace(dto.Name) == "" {
		return errors.New("el nombre es obligatorio")
	}

	if len(dto.Name) > 100 {
		return errors.New("el nombre no debe exceder los 100 caracteres")
	}

	if dto.Model != "" && len(dto.Model) > 50 {
		return errors.New("el modelo no debe exceder los 50 caracteres")
	}

	if dto.Instructions != "" && len(dto.Instructions) > 1000 {
		return errors.New("las instrucciones no deben exceder los 1000 caracteres")
	}

	if len(dto.Instructions) < 10 {
		return errors.New("la instruccion debe tener por lo menos 10 caracteres")
	}

	if len(dto.Description) < 10 {
		return errors.New("la descripción debe tener por lo menos 10 caracteres")
	}

	// Validaciones específicas para edición
	if !isCreate {
		if strings.TrimSpace(dto.OpenaiAssistantsID) == "" {
			return errors.New("openai_assistants_id es obligatorio al editar un asistente")
		}

		if len(dto.Model) < 2 {
			return errors.New("el model debe tener por lo menos 2 caracteres")
		}

		if dto.ID <= 0 {
			return errors.New("el ID del asistente es obligatorio al editar")
		}
	}

	return nil
}
