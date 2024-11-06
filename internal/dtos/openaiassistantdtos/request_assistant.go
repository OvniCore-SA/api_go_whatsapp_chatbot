package openaiassistantdtos

// Request para crear un asistente
type AssistantRequest struct {
	Instructions string          `json:"instructions"`
	Name         string          `json:"name"`
	Tools        []AssistantTool `json:"tools"`
	Model        string          `json:"model"` // Modelo de AI a utilizar:  "gpt-3.5-turbo" ,"gpt-4o", ...
}

type AssistantWithFileSearchRequest struct {
	Name          string          `json:"name"`
	Instructions  string          `json:"instructions"`
	Tools         []AssistantTool `json:"tools"`
	ToolResources ToolResources   `json:"tool_resources"`
	Model         string          `json:"model"`
}

// Forma parte del AssistantRequest, indica el tipo de dato con el que va trabajar el asistente: "file_search", "code_interpreter", "function"
type AssistantTool struct {
	Type string `json:"type"`
}

type AssistantQuery struct {
	AssistantID string `json:"assistant_id"`
	Message     string `json:"message"`
}

type Assistant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AssistantResponse struct {
	Response string `json:"response"`
}

type AssistantUploadFileQuery struct {
	PathFile string `json:"path_file"`
	Purpose  string `json:"purpose"`
}

// Struct en respuesta a la eliminación de un archivo en OpenAI
type FileDeleteResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

type FileRetrieveResponse struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Bytes     int    `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	FileName  string `json:"filename"`
	Purpose   string `json:"purpose"`
}
