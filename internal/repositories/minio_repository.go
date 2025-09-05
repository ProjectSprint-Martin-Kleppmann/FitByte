package repositories

import (
	"FitByte/internal/models"
	"FitByte/pkg/log"
	"context"

	"github.com/minio/minio-go/v7"
)

type MinioRepository interface {
	UploadFile(ctx context.Context, fileMetadata models.UploadFile) (string, error)
}

type minioRepository struct {
	client     *minio.Client
	bucketName string
}

func NewMinioRepository(storageClient *minio.Client, bucketName string) MinioRepository {
	return &minioRepository{
		client:     storageClient,
		bucketName: bucketName,
	}
}

func (r *minioRepository) UploadFile(ctx context.Context, fileMetadata models.UploadFile) (string, error) {
	info, err := r.client.PutObject(
		ctx,
		r.bucketName,
		fileMetadata.FilePath,
		fileMetadata.FileData,
		fileMetadata.Size,
		minio.PutObjectOptions{
			ContentType: fileMetadata.ContentType,
		},
	)

	if err != nil {
		log.Logger.Error().Err(err).Msg("failed to upload file")
		return "", err
	}

	return info.Key, nil
}
