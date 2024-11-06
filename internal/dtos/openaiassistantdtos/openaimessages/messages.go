package openaimessages

type SendMessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type SendMessageResponse struct {
	ID          string  `json:"id"`
	Object      string  `json:"object"`
	CreatedAt   int64   `json:"created_at"`
	AssistantID *string `json:"assistant_id"`
	ThreadID    string  `json:"thread_id"`
	RunID       *string `json:"run_id"`
	Role        string  `json:"role"`
	Content     []struct {
		Type string `json:"type"`
		Text struct {
			Value       string        `json:"value"`
			Annotations []interface{} `json:"annotations"`
		} `json:"text"`
	} `json:"content"`
	Attachments []interface{}          `json:"attachments"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// ==============================================

type GetMessageResponse struct {
	ID          string  `json:"id"`
	Object      string  `json:"object"`
	CreatedAt   int64   `json:"created_at"`
	AssistantID *string `json:"assistant_id"`
	ThreadID    string  `json:"thread_id"`
	RunID       *string `json:"run_id"`
	Role        string  `json:"role"`
	Content     []struct {
		Type string `json:"type"`
		Text struct {
			Value       string        `json:"value"`
			Annotations []interface{} `json:"annotations"`
		} `json:"text"`
	} `json:"content"`
	Attachments []interface{}          `json:"attachments"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// ==============================================

type UpdateMessageRequest struct {
	Metadata map[string]string `json:"metadata"`
}

type UpdateMessageResponse struct {
	ID          string  `json:"id"`
	Object      string  `json:"object"`
	CreatedAt   int64   `json:"created_at"`
	AssistantID *string `json:"assistant_id"`
	ThreadID    string  `json:"thread_id"`
	RunID       *string `json:"run_id"`
	Role        string  `json:"role"`
	Content     []struct {
		Type string `json:"type"`
		Text struct {
			Value       string        `json:"value"`
			Annotations []interface{} `json:"annotations"`
		} `json:"text"`
	} `json:"content"`
	FileIDs  []string          `json:"file_ids"`
	Metadata map[string]string `json:"metadata"`
}

// ==============================================

type DeleteMessageResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// ==============================================
