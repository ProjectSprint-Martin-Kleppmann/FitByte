package repositories

import (
	"FitByte/internal/models"
	"FitByte/pkg/log"
	"context"
	"errors"

	"gorm.io/gorm"
)

type ActivityRepository interface {
	CreateActivity(ctx context.Context, activity models.Activity) error
	GetActivityByID(ctx context.Context, activityID string, userID uint) (*models.Activity, error)
	UpdateActivity(ctx context.Context, activityID string, userID uint, updates map[string]interface{}) error
	DeleteActivity(ctx context.Context, activityID string, userID uint) error
	GetActivities(ctx context.Context, userID uint, params models.GetActivityParams) ([]models.Activity, error)
}

type activityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) CreateActivity(ctx context.Context, activity models.Activity) error {
	err := r.db.WithContext(ctx).Create(&activity).Error
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to create activity")
		return err
	}
	return nil
}

func (r *activityRepository) GetActivityByID(ctx context.Context, activityID string, userID uint) (*models.Activity, error) {
	var activity models.Activity
	err := r.db.WithContext(ctx).Where("activity_id = ? AND user_id = ?", activityID, userID).First(&activity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Logger.Error().Err(err).Msg("Failed to get activity by ID")
		return nil, err
	}
	return &activity, nil
}

func (r *activityRepository) UpdateActivity(ctx context.Context, activityID string, userID uint, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).
		Model(&models.Activity{}).
		Where("activity_id = ? AND user_id = ?", activityID, userID).
		Updates(updates)

	if result.Error != nil {
		log.Logger.Error().Err(result.Error).Msg("Failed to update activity")
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *activityRepository) DeleteActivity(ctx context.Context, activityID string, userID uint) error {
	result := r.db.WithContext(ctx).
		Where("activity_id = ? AND user_id = ?", activityID, userID).
		Delete(&models.Activity{})

	if result.Error != nil {
		log.Logger.Error().Err(result.Error).Msg("Failed to delete activity")
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
func (r *activityRepository) GetActivities(ctx context.Context, userID uint, params models.GetActivityParams) ([]models.Activity, error) {
	var activities []models.Activity
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if params.ActivityType != "" {
		query = query.Where("activity_type = ?", params.ActivityType)
	}
	if !params.DoneAtFrom.IsZero() {
		query = query.Where("done_at >= ?", params.DoneAtFrom)
	}
	if !params.DoneAtTo.IsZero() {
		query = query.Where("done_at <= ?", params.DoneAtTo)
	}
	if params.CaloriesBurnedMin > 0 {
		query = query.Where("calories_burned >= ?", params.CaloriesBurnedMin)
	}
	if params.CaloriesBurnedMax > 0 {
		query = query.Where("calories_burned <= ?", params.CaloriesBurnedMax)
	}

	err := query.Offset(params.Offset).Limit(params.Limit).Find(&activities).Error
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to get activities")
		return nil, err
	}
	return activities, nil

}
