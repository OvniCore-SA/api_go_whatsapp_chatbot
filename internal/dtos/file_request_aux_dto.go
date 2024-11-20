package dtos

import "mime/multipart"

type FileRequestAuxDTO struct {
	File        multipart.File // Contenido del archivo
	FileName    string         // Nombre del archivo
	FileSize    int64          // Tamaño del archivo
	ContentType string         // Tipo de contenido del archivo (MIME)
}
