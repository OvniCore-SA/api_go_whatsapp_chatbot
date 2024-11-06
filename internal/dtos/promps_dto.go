package dtos

type PrompsDto struct {
	ID          int64  `gorm:"primaryKey"` // ID como clave primaria
	Name        string // Nombre del prompt
	Descripcion string // Descripción del prompt
	Activo      bool   // Estado del prompt
	MetaAppsId  int64  // Clave foránea que apunta a MetaApps
	CreatedAt   string // Fecha de creación
	UpdatedAt   string // Fecha de actualización
}
