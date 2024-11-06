package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/config"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/entities"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/repositories/mysql_client"
	"github.com/minio/minio-go/v7"
)

type FileService struct {
	repository  *mysql_client.FileRepository
	minioClient *minio.Client
}

func NewFileService(repository *mysql_client.FileRepository, minioClient *minio.Client) *FileService {
	return &FileService{repository: repository, minioClient: minioClient}
}

func (s *FileService) CreateFile(file *multipart.FileHeader, assistantsID int64, purpose string, fileIDOpenAI string) (entities.File, error) {
	// Subir archivo a MinIO y crear el registro en la base de datos
	// Suponiendo una función `uploadToMinIO` que maneja la subida
	filePath, err := uploadToMinIO(s.minioClient, file)
	if err != nil {
		return entities.File{}, err
	}

	fileRecord := entities.File{
		AssistantsID:  assistantsID,
		Filename:      filePath,
		OpenaiFilesID: fileIDOpenAI,
		Purpose:       purpose,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repository.Create(fileRecord); err != nil {
		return entities.File{}, err
	}

	return fileRecord, nil
}

func (s *FileService) GetAllFiles() ([]entities.File, error) {
	return s.repository.FindAll()
}

func (s *FileService) GetFileById(id int64) (entities.File, error) {
	return s.repository.FindById(id)
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
func uploadToMinIO(client *minio.Client, file *multipart.FileHeader) (string, error) {
	fileContent, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileContent.Close()

	filename := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
	_, err = client.PutObject(
		context.Background(),
		config.MINIO_BUCKET_NAME, // Cambia esto por el nombre de tu bucket
		filename,
		fileContent,
		file.Size,
		minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")},
	)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func deleteFromMinIO(client *minio.Client, filename string) error {
	return client.RemoveObject(context.Background(), "bucket-name", filename, minio.RemoveObjectOptions{})
}
