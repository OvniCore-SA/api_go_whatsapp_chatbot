package filters

type OpcionPreguntasFiltro struct {
	OpcionPreguntaID int64
	ChatbotsID       int64
	PrimerMenu       bool // Se utiliza para que en el repository se tenga en cuenta si se busca el primer menu o no de este chatbot
}
