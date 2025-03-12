package entities

import "time"

type GoogleCalendarCredential struct {
	ID           int       `gorm:"primaryKey"`
	AssistantsID int       `gorm:"not null"`
	GoogleUserID string    `gorm:"not null"`
	AccessToken  string    `gorm:"type:text;not null"`
	RefreshToken string    `gorm:"type:text;not null"`
	Email        string    `gorm:"type:text;null"`
	TokenExpiry  time.Time `gorm:"not null"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt *time.Time
}
