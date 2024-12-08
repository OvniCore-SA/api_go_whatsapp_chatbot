package services

import (
	"fmt"
	"log"

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
	c := cron.New()

	// Programar tarea para las 11:00
	_, err := c.AddFunc("0 11 * * *", func() {
		log.Println("Ejecutando NotifyInteractions a las 11:00")
		if err := s.whatsappService.NotifyInteractions(); err != nil {
			log.Printf("Error en NotifyInteractions: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("error scheduling 11:00 process: %v", err)
	}

	// Programar tarea para las 17:00
	_, err = c.AddFunc("0 17 * * *", func() {
		log.Println("Ejecutando NotifyInteractions a las 17:00")
		if err := s.whatsappService.NotifyInteractions(); err != nil {
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
