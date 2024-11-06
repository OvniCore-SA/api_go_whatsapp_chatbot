package dtos

type OpcionPreguntasDto struct {
	ID                   int64                `json:"id"`
	Nombre               string               `json:"nombre"`
	OpcionPregunta       string               `json:"opcionPregunta"`
	OpcionPreguntaID     *uint                `json:"opcionPreguntaID"`
	Activo               bool                 `json:"activo"`
	UltimaOpcion         bool                 `json:"ultimaOpcion"`
	ChatbotsID           int64                `json:"chatbotsID"`
	ParentID             *int64               `json:"parentID,omitempty"`
	CreatedAt            string               `json:"createdAt"`
	UpdatedAt            string               `json:"updatedAt"`
	ChildOpcionPreguntas []OpcionPreguntasDto `json:"childOpcionPreguntas,omitempty"` // Relación uno a muchos con las subpreguntas
}
