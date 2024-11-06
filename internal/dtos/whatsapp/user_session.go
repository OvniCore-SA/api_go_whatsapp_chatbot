package whatsapp

import "time"

type UserSession struct {
	Opcion         int
	HoraConsulta   time.Time
	MenuEnviado    int
	EsUltimaOpcion bool
	MenuOpciones   map[int]int // Mapeo de IDs a opciones enviadas
}
