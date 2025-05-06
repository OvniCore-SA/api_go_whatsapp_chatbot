package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/postgres_client"
	"github.com/minio/minio-go/v7"
)

type FileService struct {
	repository  *postgres_client.FileRepository
	minioClient *minio.Client
}

func NewFileService(repository *postgres_client.FileRepository, minioClient *minio.Client) *FileService {
	return &FileService{repository: repository, minioClient: minioClient}
}

func (s *FileService) CreateFile(fileHeader *multipart.FileHeader, assistantsID int64, purpose, fileIDOpenAI, vectorStoreID string) (entities.File, error) {
	// Abrir el archivo para obtener su contenido
	file, err := fileHeader.Open()
	if err != nil {
		return entities.File{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Obtener el nombre del archivo, tamaño y tipo de contenido
	fileName := fileHeader.Filename
	fileSize := fileHeader.Size
	contentType := fileHeader.Header.Get("Content-Type")

	// Subir archivo a MinIO
	filePath, err := uploadToMinIO(s.minioClient, file, fileName, fileSize, contentType)
	if err != nil {
		return entities.File{}, fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	// Crear el registro en la base de datos
	fileRecord := entities.File{
		AssistantsID:         assistantsID,
		Filename:             filePath,
		OpenaiFilesID:        fileIDOpenAI,
		OpenaiVectorStoreIDs: vectorStoreID,
		Purpose:              purpose,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := s.repository.Create(fileRecord); err != nil {
		return entities.File{}, fmt.Errorf("failed to save file record to database: %w", err)
	}

	return fileRecord, nil
}

func (s *FileService) GetAllFiles() ([]entities.File, error) {
	return s.repository.FindAll()
}

func (s *FileService) GetFileById(id int64) (entities.File, error) {
	return s.repository.FindById(id)
}

func (s *FileService) GetFileByAssistantID(assistantID int64) ([]dtos.FileDto, error) {

	files, err := s.repository.FindByAssistantID(assistantID)
	if err != nil {
		return nil, errors.New("files not found")
	}

	var fileDTOs []dtos.FileDto
	for _, fileEntitie := range files {
		fileDTO := entities.MapEntityToFileDto(fileEntitie)
		fileDTOs = append(fileDTOs, fileDTO)
	}

	return fileDTOs, nil
}

func (s *FileService) UpdateFile(id, assistantsID int64, purpose string) (entities.File, error) {
	file, err := s.repository.FindById(id)
	if err != nil {
		return entities.File{}, errors.New("file not found")
	}

	file.AssistantsID = assistantsID
	file.Purpose = purpose
	file.UpdatedAt = time.Now()

	if err := s.repository.Update(file); err != nil {
		return entities.File{}, err
	}

	return file, nil
}

func (s *FileService) DeleteFile(id int64) error {
	file, err := s.repository.FindById(id)
	if err != nil {
		return errors.New("file not found")
	}

	// Eliminar archivo en MinIO y en base de datos
	if err := deleteFromMinIO(s.minioClient, file.Filename); err != nil {
		return err
	}

	return s.repository.Delete(id)
}
func uploadToMinIO(client *minio.Client, file multipart.File, fileName string, fileSize int64, contentType string) (string, error) {
	// Generar un nombre único para el archivo
	uniqueFileName := fmt.Sprintf("%d-%s", time.Now().Unix(), fileName)

	// Subir el archivo al bucket de MinIO
	_, err := client.PutObject(
		context.Background(),
		os.Getenv("MINIO_BUCKET_NAME"), // Cambia esto por el nombre de tu bucket
		uniqueFileName,
		file,
		fileSize,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	return uniqueFileName, nil
}

func deleteFromMinIO(client *minio.Client, filename string) error {
	return client.RemoveObject(context.Background(), os.Getenv("MINIO_BUCKET_NAME"), filename, minio.RemoveObjectOptions{})
}
