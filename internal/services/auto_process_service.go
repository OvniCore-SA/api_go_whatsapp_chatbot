package services

import (
	"fmt"
	"log"
	"os"

	"github.com/robfig/cron/v3"
)

// AutoProcessService estructura para manejar procesos automáticos
type AutoProcessService struct {
	whatsappService *WhatsappService
}

// NewAutoProcessService inicializa un nuevo AutoProcessService
func NewAutoProcessService(whatsappService *WhatsappService) *AutoProcessService {
	return &AutoProcessService{
		whatsappService: whatsappService,
	}
}

// Start inicia los procesos automáticos
func (s *AutoProcessService) Start() error {
	if os.Getenv("HOST_API") == "http://127.0.0.1:5001" {
		return nil
	}
	c := cron.New()

	//
	// Programar tarea para las 11:00
	_, err := c.AddFunc("0 11 * * *", func() {
		log.Println("Ejecutando NotifyInteractions a las 11:00")
		if err := s.whatsappService.NotifyInteractions(11); err != nil {
			log.Printf("Error en NotifyInteractions: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("error scheduling 11:00 process: %v", err)
	}

	_, err = c.AddFunc("13 14 * * *", func() {
		log.Println("Ejecutando NotifyInteractions a las 14:13")

	})
	if err != nil {
		return fmt.Errorf("error scheduling 11:00 process: %v", err)
	}

	// Programar tarea para las 17:00
	_, err = c.AddFunc("0 17 * * *", func() {
		log.Println("Ejecutando NotifyInteractions a las 17:00")
		if err := s.whatsappService.NotifyInteractions(6); err != nil {
			log.Printf("Error en NotifyInteractions: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("error scheduling 17:00 process: %v", err)
	}

	// Iniciar el cron
	c.Start()

	log.Println("AutoProcessService iniciado con éxito.")
	return nil
}
