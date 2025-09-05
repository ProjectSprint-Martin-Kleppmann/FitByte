package repositories

import (
	"FitByte/internal/models"
	"FitByte/pkg/log"
	"context"

	"gorm.io/gorm"
)

type FileRepository interface {
	Insert(ctx context.Context, file models.File) error
}

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{
		db: db,
	}
}

func (r *fileRepository) Insert(ctx context.Context, file models.File) error {
	err := r.db.Table("files").WithContext(ctx).Create(&file).Error
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to insert file")
		return err
	}
	return nil
}
