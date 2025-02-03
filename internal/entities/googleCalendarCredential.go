package entities

import "time"

type GoogleCalendarCredential struct {
	ID           int       `gorm:"primaryKey"`
	AssistantsID int       `gorm:"not null"`
	GoogleUserID string    `gorm:"not null"`
	AccessToken  string    `gorm:"type:text;not null"`
	RefreshToken string    `gorm:"type:text;not null"`
	TokenExpiry  time.Time `gorm:"not null"`
	// Horarios de trabajo
	WorkStartTime string    `gorm:"type:varchar(5);not null"`  // Hora de inicio (HH:mm)
	WorkEndTime   string    `gorm:"type:varchar(5);not null"`  // Hora de fin (HH:mm)
	DaysOpen      string    `gorm:"type:varchar(60);not null"` // Días abiertos como una lista (ejemplo: "lunes,martes,miercoles,jueves,viernes")
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
	DeletedAt     *time.Time
}
