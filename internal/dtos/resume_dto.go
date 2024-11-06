package dtos

type ResumesDto struct {
	ID               int64  `json:"id"`
	RequestToResolve string `json:"request_to_resolve"`
	ChatbotID        int64  `json:"chatbot_id"`
}
