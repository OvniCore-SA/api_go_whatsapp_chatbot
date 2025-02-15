package config

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: true, // Optimiza la reutilización de sentencias preparadas
	})
	if err != nil {
		fmt.Println("====")
		fmt.Println(os.Getenv("DB_HOST"))
		fmt.Println("====")
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
