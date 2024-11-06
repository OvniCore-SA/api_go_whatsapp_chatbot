package openairuns

type StartRunRequest struct {
	AssistantID string `json:"assistant_id"`
}

type StartRunResponse struct {
	ID                string  `json:"id"`
	Object            string  `json:"object"`
	CreatedAt         int64   `json:"created_at"`
	AssistantID       string  `json:"assistant_id"`
	ThreadID          string  `json:"thread_id"`
	Status            string  `json:"status"`
	StartedAt         int64   `json:"started_at"`
	ExpiresAt         *int64  `json:"expires_at"`
	CancelledAt       *int64  `json:"cancelled_at"`
	FailedAt          *int64  `json:"failed_at"`
	CompletedAt       *int64  `json:"completed_at"`
	LastError         *string `json:"last_error"`
	Model             string  `json:"model"`
	Instructions      *string `json:"instructions"`
	IncompleteDetails *string `json:"incomplete_details"`
	Tools             []struct {
		Type string `json:"type"`
	} `json:"tools"`
	Metadata map[string]interface{} `json:"metadata"`
	Usage    *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Temperature         float64 `json:"temperature"`
	TopP                float64 `json:"top_p"`
	MaxPromptTokens     int     `json:"max_prompt_tokens"`
	MaxCompletionTokens int     `json:"max_completion_tokens"`
	TruncationStrategy  struct {
		Type         string `json:"type"`
		LastMessages *int   `json:"last_messages"`
	} `json:"truncation_strategy"`
	ResponseFormat    string `json:"response_format"`
	ToolChoice        string `json:"tool_choice"`
	ParallelToolCalls bool   `json:"parallel_tool_calls"`
}

// ===================================================================

type GetRunResponse struct {
	ID                string  `json:"id"`
	Object            string  `json:"object"`
	CreatedAt         int64   `json:"created_at"`
	AssistantID       string  `json:"assistant_id"`
	ThreadID          string  `json:"thread_id"`
	Status            string  `json:"status"`
	StartedAt         int64   `json:"started_at"`
	ExpiresAt         *int64  `json:"expires_at"`
	CancelledAt       *int64  `json:"cancelled_at"`
	FailedAt          *int64  `json:"failed_at"`
	CompletedAt       *int64  `json:"completed_at"`
	LastError         *string `json:"last_error"`
	Model             string  `json:"model"`
	Instructions      *string `json:"instructions"`
	IncompleteDetails *string `json:"incomplete_details"`
	Tools             []struct {
		Type string `json:"type"`
	} `json:"tools"`
	ToolResources map[string]interface{} `json:"tool_resources"`
	Metadata      map[string]interface{} `json:"metadata"`
	Usage         *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Temperature         float64 `json:"temperature"`
	TopP                float64 `json:"top_p"`
	MaxPromptTokens     int     `json:"max_prompt_tokens"`
	MaxCompletionTokens int     `json:"max_completion_tokens"`
	TruncationStrategy  struct {
		Type         string `json:"type"`
		LastMessages *int   `json:"last_messages"`
	} `json:"truncation_strategy"`
	ResponseFormat    string `json:"response_format"`
	ToolChoice        string `json:"tool_choice"`
	ParallelToolCalls bool   `json:"parallel_tool_calls"`
}

// ===================================================================
// ===================================================================

type CancelRunResponse struct {
	ID            string                 `json:"id"`
	Object        string                 `json:"object"`
	CreatedAt     int64                  `json:"created_at"`
	AssistantID   string                 `json:"assistant_id"`
	ThreadID      string                 `json:"thread_id"`
	Status        string                 `json:"status"`
	StartedAt     int64                  `json:"started_at"`
	ExpiresAt     *int64                 `json:"expires_at"`
	CancelledAt   *int64                 `json:"cancelled_at"`
	FailedAt      *int64                 `json:"failed_at"`
	CompletedAt   *int64                 `json:"completed_at"`
	LastError     *string                `json:"last_error"`
	Model         string                 `json:"model"`
	Instructions  *string                `json:"instructions"`
	ToolResources map[string]interface{} `json:"tool_resources"`
	Metadata      map[string]interface{} `json:"metadata"`
	Usage         *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Temperature       float64 `json:"temperature"`
	TopP              float64 `json:"top_p"`
	ResponseFormat    string  `json:"response_format"`
	ToolChoice        string  `json:"tool_choice"`
	ParallelToolCalls bool    `json:"parallel_tool_calls"`
}

// ===================================================================
// ===================================================================

type UpdateRunMetadataRequest struct {
	Metadata map[string]string `json:"metadata"`
}

type UpdateRunMetadataResponse struct {
	ID                string  `json:"id"`
	Object            string  `json:"object"`
	CreatedAt         int64   `json:"created_at"`
	AssistantID       string  `json:"assistant_id"`
	ThreadID          string  `json:"thread_id"`
	Status            string  `json:"status"`
	StartedAt         int64   `json:"started_at"`
	ExpiresAt         *int64  `json:"expires_at"`
	CancelledAt       *int64  `json:"cancelled_at"`
	FailedAt          *int64  `json:"failed_at"`
	CompletedAt       *int64  `json:"completed_at"`
	LastError         *string `json:"last_error"`
	Model             string  `json:"model"`
	Instructions      *string `json:"instructions"`
	IncompleteDetails *string `json:"incomplete_details"`
	Tools             []struct {
		Type string `json:"type"`
	} `json:"tools"`
	ToolResources map[string]interface{} `json:"tool_resources"`
	Metadata      map[string]string      `json:"metadata"`
	Usage         *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Temperature         float64 `json:"temperature"`
	TopP                float64 `json:"top_p"`
	MaxPromptTokens     int     `json:"max_prompt_tokens"`
	MaxCompletionTokens int     `json:"max_completion_tokens"`
	TruncationStrategy  struct {
		Type         string `json:"type"`
		LastMessages *int   `json:"last_messages"`
	} `json:"truncation_strategy"`
	ResponseFormat    string `json:"response_format"`
	ToolChoice        string `json:"tool_choice"`
	ParallelToolCalls bool   `json:"parallel_tool_calls"`
}
