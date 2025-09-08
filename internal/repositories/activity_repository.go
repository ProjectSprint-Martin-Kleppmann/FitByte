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
   GetActivitiesByUserID(ctx context.Context, userID uint, query models.GetActivitiesQuery) ([]models.Activity, error)
   GetActivityByID(ctx context.Context, activityID string, userID uint) (*models.Activity, error)
   UpdateActivity(ctx context.Context, activityID string, userID uint, updates map[string]interface{}) error
   DeleteActivity(ctx context.Context, activityID string, userID uint) error
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


func (r *activityRepository) GetActivitiesByUserID(ctx context.Context, userID uint, query models.GetActivitiesQuery) ([]models.Activity, error) {
   var activities []models.Activity


   db := r.db.WithContext(ctx).Where("user_id = ?", userID)


   // Apply filters
   if query.ActivityType != "" {
       db = db.Where("activity_type = ?", query.ActivityType)
   }


   if !query.DoneAtFrom.IsZero() {
       db = db.Where("done_at >= ?", query.DoneAtFrom)
   }


   if !query.DoneAtTo.IsZero() {
       db = db.Where("done_at <= ?", query.DoneAtTo)
   }


   if query.CaloriesBurnedMin > 0 {
       db = db.Where("calories_burned >= ?", query.CaloriesBurnedMin)
   }


   if query.CaloriesBurnedMax > 0 {
       db = db.Where("calories_burned <= ?", query.CaloriesBurnedMax)
   }


   // Apply pagination
   if query.Limit > 0 {
       db = db.Limit(query.Limit)
   }
   if query.Offset > 0 {
       db = db.Offset(query.Offset)
   }


   // Order by done_at descending (most recent first)
   db = db.Order("done_at DESC")


   err := db.Find(&activities).Error
   if err != nil {
       log.Logger.Error().Err(err).Msg("Failed to get activities by user ID")
       return nil, err
   }


   return activities, nil
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
