package dtos

type ContactDto struct {
	ID              int64       `json:"id"`
	NumberPhonesID  int64       `json:"number_phones_id"`
	ContactNumber   int64       `json:"contact_number"`
	OpenaiThreadsID string      `json:"openai_threads_id"`
	CountTokens     string      `json:"count_tokens"`
	Events          []EventsDto `json:"events,omitempty"`
}
