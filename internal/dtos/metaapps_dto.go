package dtos

type MetaAppsDto struct {
	ID           int64 // ID como clave primaria
	Promps       PrompsDto
	Chatbots     []ChatbotsDto // Clave foránea que apunta al Chatbot
	Name         string        // Nombre de la app en meta
	AplicationId string        // Identificador de la app en meta
	Company      string        // Nombre de empresa en meta
	CreatedAt    string        // Fecha de creación
	UpdatedAt    string        // Fecha de actualización
	DeletedAt    string
}
