package mappers

import (
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
)

// Convertir de DTO a entidad
func ToThreadEntity(dto dtos.ThreadCreateRequest) entities.Thread {
	return entities.Thread{
		Active:          dto.Active,
		OpenaiThreadsId: dto.OpenaiThreadsId,
		ContactsID:      dto.ContactsID,
	}
}

// Convertir de entidad a DTO de respuesta
func ToThreadResponse(entity entities.Thread) dtos.ThreadResponse {
	return dtos.ThreadResponse{
		ID:              entity.ID,
		OpenaiThreadsId: entity.OpenaiThreadsId,
		ContactsID:      entity.ContactsID,
		Active:          entity.Active,
	}
}

// Convertir lista de entidades a lista de DTOs de respuesta
func ToThreadResponseList(entities []entities.Thread) []dtos.ThreadResponse {
	var responseList []dtos.ThreadResponse
	for _, entity := range entities {
		responseList = append(responseList, ToThreadResponse(entity))
	}
	return responseList
}
