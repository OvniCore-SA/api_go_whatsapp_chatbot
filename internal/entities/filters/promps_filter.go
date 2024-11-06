package filters

type PrompsFiltro struct {
	MetaAppsId int64
	Activo     *bool // Utilizamos un puntero para que sea opcional
}
