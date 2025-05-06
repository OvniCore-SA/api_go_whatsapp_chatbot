package config

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Argentina/Buenos_Aires  search_path=chatbot_whatsapp",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Obtener la conexión SQL para configurar el pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Configuración del pool de conexiones
	sqlDB.SetMaxOpenConns(25)                  // Máximo de conexiones abiertas al mismo tiempo
	sqlDB.SetMaxIdleConns(10)                  // Máximo de conexiones en espera
	sqlDB.SetConnMaxLifetime(10 * time.Minute) // Reiniciar conexiones cada 10 minutos
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)  // Tiempo máximo de inactividad antes de cerrar una conexión

	DB = db
	return db, nil
}
