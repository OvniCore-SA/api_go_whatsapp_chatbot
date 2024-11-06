package ollamaDtos

type RequestPayload struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
	Options  Options   `json:"options"`
	Format   string    `json:"format"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Options struct {
	Temperature float64 `json:"temperature"`
}

type Response struct {
	Model      string         `json:"model"`
	CreatedAt  string         `json:"created_at"`
	Message    MessageContent `json:"message"`
	Done       bool           `json:"done"`
	DoneReason string         `json:"done_reason,omitempty"`
}

type MessageContent struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
