package dtos

type Pagination struct {
	Number   uint32 `json:"number"`    // Número de la página actual
	Size     uint32 `json:"size"`      // Cantidad de registros por página
	Total    int64  `json:"total"`     // Total de registros en la base de datos
	LastPage uint32 `json:"last_page"` // Número total de páginas
	From     uint32 `json:"from"`      // Primer registro en la página actual
	To       uint32 `json:"to"`        // Último registro en la página actual
}

// Setea por defecto number en 1 y size en 10.
func (p *Pagination) SetDefaults() {
	if p.Number == 0 {
		p.Number = 1
	}
	if p.Size == 0 {
		p.Size = 10
	}
}
