package dtos

// ThreadCreateRequest representa la estructura de una solicitud para crear un hilo
type ThreadCreateRequest struct {
	OpenaiThreadsId string `json:"openai_threads_id" validate:"required"`
	Active          bool   `json:"active"`
	ContactsID      int64  `json:"contacts_id" validate:"required"`
}

// ThreadResponse representa la estructura de la respuesta de un hilo
type ThreadResponse struct {
	ID              int64  `json:"id"`
	OpenaiThreadsId string `json:"openai_threads_id"`
	ContactsID      int64  `json:"contacts_id"`
	Active          bool   `json:"active"`
}
