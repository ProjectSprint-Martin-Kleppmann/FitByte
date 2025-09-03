package repositories

import (
	"FitByte/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"

	"FitByte/pkg/log"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user models.User) error {
	err := r.db.Table("users").WithContext(ctx).Create(&user).Error
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to create user")
		return err
	}
	return nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.Table("users").WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Logger.Error().Err(err).Msg("Failed to get user by email")
		return nil, err
	}
	return &user, nil
}
