package models

import (
	"time"

	"gorm.io/gorm"
)

// Activity represents the activity entity in the database
type Activity struct {
	gorm.Model
	ActivityID          string    `json:"activityId" gorm:"uniqueIndex;not null"`
	UserID              uint      `json:"-" gorm:"not null;index"`
	ActivityType        string    `json:"activityType" gorm:"not null" validate:"required,oneof=Walking Yoga Stretching Cycling Swimming Dancing Hiking Running HIIT JumpRope"`
	DoneAt              time.Time `json:"doneAt" gorm:"not null" validate:"required"`
	DurationInMinutes   int       `json:"durationInMinutes" gorm:"not null" validate:"required,min=1"`
	CaloriesBurned      int       `json:"caloriesBurned" gorm:"not null"`
}

// CreateActivityRequest represents the request body for creating an activity
type CreateActivityRequest struct {
	ActivityType      string `json:"activityType" validate:"required,oneof=Walking Yoga Stretching Cycling Swimming Dancing Hiking Running HIIT JumpRope"`
	DoneAt            string `json:"doneAt" validate:"required"`
	DurationInMinutes int    `json:"durationInMinutes" validate:"required,min=1"`
}

// UpdateActivityRequest represents the request body for updating an activity
type UpdateActivityRequest struct {
	ActivityType      *string `json:"activityType,omitempty" validate:"omitempty,oneof=Walking Yoga Stretching Cycling Swimming Dancing Hiking Running HIIT JumpRope"`
	DoneAt            *string `json:"doneAt,omitempty"`
	DurationInMinutes *int    `json:"durationInMinutes,omitempty" validate:"omitempty,min=1"`
}

// ActivityResponse represents the response format for activity operations
type ActivityResponse struct {
	ActivityID        string    `json:"activityId"`
	ActivityType      string    `json:"activityType"`
	DoneAt            time.Time `json:"doneAt"`
	DurationInMinutes int       `json:"durationInMinutes"`
	CaloriesBurned    int       `json:"caloriesBurned"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// ActivityTypeCalories maps activity types to calories burned per minute
var ActivityTypeCalories = map[string]int{
	"Walking":    4,
	"Yoga":       4,
	"Stretching": 4,
	"Cycling":    8,
	"Swimming":   8,
	"Dancing":    8,
	"Hiking":     10,
	"Running":    10,
	"HIIT":       10,
	"JumpRope":   10,
}
