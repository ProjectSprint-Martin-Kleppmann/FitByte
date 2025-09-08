package service

import (
	"FitByte/internal/models"
	"FitByte/internal/repositories"
	"FitByte/pkg/log"
	customErrors "FitByte/internal/errors"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityService interface {
	CreateActivity(ctx context.Context, userID uint, req models.CreateActivityRequest) (*models.ActivityResponse, error)
	GetActivities(ctx context.Context, userID uint, query models.GetActivitiesQuery) ([]models.ActivityResponse, error)
	UpdateActivity(ctx context.Context, userID uint, activityID string, req models.UpdateActivityRequest) (*models.ActivityResponse, error)
	DeleteActivity(ctx context.Context, userID uint, activityID string) error
}

type activityService struct {
	activityRepo repositories.ActivityRepository
}

func NewActivityService(activityRepo repositories.ActivityRepository) ActivityService {
	return &activityService{
		activityRepo: activityRepo,
	}
}

func (s *activityService) CreateActivity(ctx context.Context, userID uint, req models.CreateActivityRequest) (*models.ActivityResponse, error) {
	// Parse the doneAt time
	doneAt, err := time.Parse(time.RFC3339, req.DoneAt)
	if err != nil {
		log.Logger.Error().Err(err).Str("doneAt", req.DoneAt).Msg("Failed to parse doneAt time")
		return nil, err
	}

	// Calculate calories burned
	caloriesPerMinute, exists := models.ActivityTypeCalories[req.ActivityType]
	if !exists {
		log.Logger.Error().Str("activityType", req.ActivityType).Msg("Invalid activity type")
		return nil, err
	}
	caloriesBurned := caloriesPerMinute * req.DurationInMinutes

	// Generate unique activity ID
	activityID := uuid.New().String()

	activity := models.Activity{
		ActivityID:        activityID,
		UserID:            userID,
		ActivityType:      req.ActivityType,
		DoneAt:            doneAt,
		DurationInMinutes: req.DurationInMinutes,
		CaloriesBurned:    caloriesBurned,
	}

	err = s.activityRepo.CreateActivity(ctx, activity)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to create activity")
		return nil, err
	}

	// Return the created activity with timestamps
	now := time.Now()
	return &models.ActivityResponse{
		ActivityID:        activityID,
		ActivityType:      req.ActivityType,
		DoneAt:            doneAt.Format(time.RFC3339),
		DurationInMinutes: req.DurationInMinutes,
		CaloriesBurned:    caloriesBurned,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

func (s *activityService) GetActivities(ctx context.Context, userID uint, query models.GetActivitiesQuery) ([]models.ActivityResponse, error) {
	// Set default pagination if not provided
	if query.Limit <= 0 {
		query.Limit = 5
	}
	if query.Offset < 0 {
		query.Offset = 0
	}

	activities, err := s.activityRepo.GetActivitiesByUserID(ctx, userID, query)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to get activities")
		return nil, err
	}

	responses := make([]models.ActivityResponse, len(activities))
	for i, activity := range activities {
		responses[i] = models.ActivityResponse{
			ActivityID:        activity.ActivityID,
			ActivityType:      activity.ActivityType,
			DoneAt:            activity.DoneAt.Format(time.RFC3339),
			DurationInMinutes: activity.DurationInMinutes,
			CaloriesBurned:    activity.CaloriesBurned,
			CreatedAt:         activity.CreatedAt,
			UpdatedAt:         activity.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *activityService) UpdateActivity(ctx context.Context, userID uint, activityID string, req models.UpdateActivityRequest) (*models.ActivityResponse, error) {
	// Check if activity exists
	existingActivity, err := s.activityRepo.GetActivityByID(ctx, activityID, userID)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to get activity for update")
		return nil, err
	}
	if existingActivity == nil {
		return nil, customErrors.ErrActivityNotFound
	}

	updates := make(map[string]interface{})
	
	// Track what fields are being updated for recalculation
	var newActivityType string = existingActivity.ActivityType
	var newDurationInMinutes int = existingActivity.DurationInMinutes

	if req.ActivityType != nil {
		updates["activity_type"] = *req.ActivityType
		newActivityType = *req.ActivityType
	}

	if req.DoneAt != nil {
		doneAt, err := time.Parse(time.RFC3339, *req.DoneAt)
		if err != nil {
			log.Logger.Error().Err(err).Str("doneAt", *req.DoneAt).Msg("Failed to parse doneAt time")
			return nil, err
		}
		updates["done_at"] = doneAt
	}

	if req.DurationInMinutes != nil {
		updates["duration_in_minutes"] = *req.DurationInMinutes
		newDurationInMinutes = *req.DurationInMinutes
	}

	// Recalculate calories if activity type or duration changed
	if req.ActivityType != nil || req.DurationInMinutes != nil {
		caloriesPerMinute, exists := models.ActivityTypeCalories[newActivityType]
		if !exists {
			log.Logger.Error().Str("activityType", newActivityType).Msg("Invalid activity type")
			return nil, err
		}
		newCaloriesBurned := caloriesPerMinute * newDurationInMinutes
		updates["calories_burned"] = newCaloriesBurned
	}

	// Add updated_at timestamp
	updates["updated_at"] = time.Now()

	err = s.activityRepo.UpdateActivity(ctx, activityID, userID, updates)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, customErrors.ErrActivityNotFound
		}
		log.Logger.Error().Err(err).Msg("Failed to update activity")
		return nil, err
	}

	// Get the updated activity to return
	updatedActivity, err := s.activityRepo.GetActivityByID(ctx, activityID, userID)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to get updated activity")
		return nil, err
	}

	// Use the original request doneAt format if it was provided, otherwise use the stored format
	doneAtString := updatedActivity.DoneAt.Format(time.RFC3339)
	if req.DoneAt != nil {
		// If doneAt was updated, use the original request format
		doneAtString = *req.DoneAt
	}

	return &models.ActivityResponse{
		ActivityID:        updatedActivity.ActivityID,
		ActivityType:      updatedActivity.ActivityType,
		DoneAt:            doneAtString,
		DurationInMinutes: updatedActivity.DurationInMinutes,
		CaloriesBurned:    updatedActivity.CaloriesBurned,
		CreatedAt:         updatedActivity.CreatedAt,
		UpdatedAt:         updatedActivity.UpdatedAt,
	}, nil
}

func (s *activityService) DeleteActivity(ctx context.Context, userID uint, activityID string) error {
	err := s.activityRepo.DeleteActivity(ctx, activityID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return customErrors.ErrActivityNotFound
		}
		log.Logger.Error().Err(err).Msg("Failed to delete activity")
		return err
	}
	return nil
}

