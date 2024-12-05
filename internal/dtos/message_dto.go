package dtos

import "time"

type MessageDto struct {
	ID             int64     `json:"id"`
	NumberPhonesID int64     `json:"number_phones_id"`
	ContactsID     int64     `json:"contacts_id"`
	MessageText    string    `json:"message_text"`
	IsFromBot      bool      `json:"is_from_bot"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
