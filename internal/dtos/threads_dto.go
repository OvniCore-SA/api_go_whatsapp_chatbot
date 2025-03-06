package dtos

// ThreadCreateRequest representa la estructura de una solicitud para crear un hilo
type ThreadCreateRequest struct {
	ThreadsId string `json:"threads_id" validate:"required"`
	Active    bool   `json:"active"`
}

// ThreadResponse representa la estructura de la respuesta de un hilo
type ThreadResponse struct {
	ID        int64  `json:"id"`
	ThreadsId string `json:"threads_id"`
	Active    bool   `json:"active"`
}
