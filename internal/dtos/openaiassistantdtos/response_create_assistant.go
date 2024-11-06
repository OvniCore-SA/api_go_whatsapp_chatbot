package openaiassistantdtos

// Estructura de respuesta, misma que en la función anterior
type ResponseCreateAssistant struct {
	ID             string                 `json:"id"`
	Object         string                 `json:"object"`
	CreatedAt      int64                  `json:"created_at"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Model          string                 `json:"model"`
	Instructions   string                 `json:"instructions"`
	Tools          []AssistantTool        `json:"tools"`
	ToolResources  ToolResources          `json:"tool_resources"`
	Metadata       map[string]interface{} `json:"metadata"`
	TopP           float64                `json:"top_p"`
	Temperature    float64                `json:"temperature"`
	ResponseFormat string                 `json:"response_format"`
}

type ToolResources struct {
	FileSearch FileSearchResource `json:"file_search"`
}

type FileSearchResource struct {
	VectorStoreIDs []string `json:"vector_store_ids"`
}
