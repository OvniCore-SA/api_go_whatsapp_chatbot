package openaiassistantdtos

// FileUploadResponse define la estructura para la respuesta después de subir un archivo
type FileUploadResponse struct {
	ID       string `json:"id"`
	Object   string `json:"object"`
	Bytes    int    `json:"bytes"`
	Created  int    `json:"created_at"`
	Filename string `json:"filename"`
	Purpose  string `json:"purpose"`
}
