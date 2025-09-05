package service

import (
	"FitByte/internal/models"
	"FitByte/internal/repositories"
	"FitByte/pkg/log"
	"context"
)

type FileService interface {
	SaveFileUpload(ctx context.Context, userID int64, file models.UploadFile) (string, error)
}

type fileService struct {
	fileRepo  repositories.FileRepository
	minioRepo repositories.MinioRepository
}

func NewFileService(fileRepo repositories.FileRepository, storageRepo repositories.MinioRepository) FileService {
	return &fileService{
		fileRepo:  fileRepo,
		minioRepo: storageRepo,
	}
}

func (s *fileService) SaveFileUpload(ctx context.Context, userID int64, file models.UploadFile) (string, error) {
	key, err := s.minioRepo.UploadFile(ctx, file)
	if err != nil {
		log.Logger.Error().Err(err).Msg("minioRepo.UploadFile")
		return "", err
	}

	fileTableSchema := models.File{
		UserID:   userID,
		FileName: file.FileName,
		FileURL:  file.FilePath,
	}

	err = s.fileRepo.Insert(ctx, fileTableSchema)
	if err != nil {
		log.Logger.Error().Err(err).Msg("fileRepo.Insert")
		return "", err
	}

	return key, nil
}
