package openaithreads

// Message representa un mensaje enviado o recibido por el asistente.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ThreadRequest representa la solicitud para crear un nuevo hilo (thread).
type ThreadRequest struct {
	Messages []Message `json:"messages,omitempty"` // Campo opcional
}

// ThreadResponse representa la respuesta recibida después de crear un hilo.
type ThreadResponse struct {
	ID            string                 `json:"id"`
	Object        string                 `json:"object"`
	CreatedAt     int64                  `json:"created_at"`
	Metadata      map[string]interface{} `json:"metadata"`
	ToolResources map[string]interface{} `json:"tool_resources"`
}

// ===========================================
// ThreadDetailsResponse representa la respuesta al obtener los detalles de un hilo.
type ThreadDetailsResponse struct {
	ID            string                 `json:"id"`
	Object        string                 `json:"object"`
	CreatedAt     int64                  `json:"created_at"`
	Metadata      map[string]interface{} `json:"metadata"`
	ToolResources map[string]interface{} `json:"tool_resources"`
}

// ===========================================

// ThreadUpdateRequest representa la solicitud para actualizar un hilo.
type ThreadUpdateRequest struct {
	Metadata map[string]interface{} `json:"metadata"`
}

// ThreadUpdateResponse representa la respuesta al actualizar un hilo.
type ThreadUpdateResponse struct {
	ID            string                 `json:"id"`
	Object        string                 `json:"object"`
	CreatedAt     int64                  `json:"created_at"`
	Metadata      map[string]interface{} `json:"metadata"`
	ToolResources map[string]interface{} `json:"tool_resources"`
}

// ===========================================
// ThreadDeleteResponse representa la respuesta al eliminar un hilo.
type ThreadDeleteResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}
