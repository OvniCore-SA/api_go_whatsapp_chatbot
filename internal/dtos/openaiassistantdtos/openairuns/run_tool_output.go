package openairuns

// OpenAIRunResponse representa la respuesta completa de un run en OpenAI
type OpenAIRunResponse struct {
	ID             string         `json:"id"`
	ThreadID       string         `json:"thread_id"`
	Status         string         `json:"status"`
	RequiredAction RequiredAction `json:"required_action"`
}

// RequiredAction representa la acción requerida cuando el estado es "requires_action"
type RequiredAction struct {
	Type              string            `json:"type"`
	SubmitToolOutputs SubmitToolOutputs `json:"submit_tool_outputs"`
}

// SubmitToolOutputs contiene los tool_calls requeridos
type SubmitToolOutputs struct {
	ToolCalls []ToolCall `json:"tool_calls"`
}

// ToolCall representa una llamada a una función específica dentro de OpenAI
type ToolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function representa los detalles de la función invocada por OpenAI
type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ToolCallOutput representa la estructura esperada para la respuesta de `submit_tool_outputs`
type ToolCallOutput struct {
	ToolCallID string `json:"tool_call_id"`
	Output     string `json:"output"`
}

// OpenAIToolOutput representa la estructura de salida para submit_tool_outputs
type OpenAIToolOutput struct {
	ToolCallID string `json:"tool_call_id"`
	Output     string `json:"output"`
}
