package repositories

import (
	"FitByte/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"

	"FitByte/pkg/log"
)

type ProfileRepository interface {
	CreateUser(ctx context.Context, profile *models.Profile) error
	GetProfileByEmail(ctx context.Context, email string) (*models.Profile, error)
	UpdateUser(ctx context.Context, userID uint, updates map[string]interface{}) error
	GetProfileByID(ctx context.Context, userID uint) (*models.Profile, error)
}

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) CreateUser(ctx context.Context, profile *models.Profile) error {
	err := r.db.WithContext(ctx).Create(profile).Error
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to create profile")
		return err
	}
	return nil
}

func (r *profileRepository) GetProfileByEmail(ctx context.Context, email string) (*models.Profile, error) {
	var profile models.Profile
	err := r.db.Table("profiles").WithContext(ctx).Where("email = ?", email).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Logger.Error().Err(err).Msg("Failed to get profile by email")
		return nil, err
	}
	return &profile, nil
}

func (r *profileRepository) UpdateUser(ctx context.Context, userID uint, updates map[string]interface{}) error {
	err := r.db.Table("profiles").WithContext(ctx).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to update profile user")
		return err
	}
	return nil
}
func (r *profileRepository) GetProfileByID(ctx context.Context, userID uint ) (*models.Profile, error) {
	var profile models.Profile
	err := r.db.Table("profiles").WithContext(ctx).Where("id = ?", userID).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Logger.Error().Err(err).Msg("Failed to get profile by ID")
		return nil, err
	}
	return &profile, nil
}