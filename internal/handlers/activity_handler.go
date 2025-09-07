package handlers

import (
	"FitByte/configs"
	customErrors "FitByte/internal/errors"
	"FitByte/internal/middleware"
	"FitByte/internal/models"
	"FitByte/internal/service"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ActivityHandler struct {
	Engine      *gin.Engine
	AppConfig   configs.Config
	ActivitySvc service.ActivityService
	validator   *validator.Validate
}

func NewActivityHandler(engine *gin.Engine, appConfig configs.Config, activityService service.ActivityService) *ActivityHandler {
	return &ActivityHandler{
		Engine:      engine,
		AppConfig:   appConfig,
		ActivitySvc: activityService,
		validator:   validator.New(),
	}
}

func (h *ActivityHandler) SetupRoutes() {
	protectedRoutes := h.Engine.Group("/v1")
	protectedRoutes.Use(middleware.AuthMiddleware(h.AppConfig.Secret.JWTSecret))
	protectedRoutes.Use(middleware.ContentTypeMiddleware())
	protectedRoutes.Use(middleware.ValidationMiddleware())

	// POST activity with null validation for required fields
	protectedRoutes.POST("/activity", 
		middleware.ValidateJSONForNulls([]string{"activityType", "doneAt", "durationInMinutes"}),
		h.CreateActivity)
	
	protectedRoutes.GET("/activity", h.GetActivities)
	
	// PATCH activity with null validation for optional fields that shouldn't be null when provided
	protectedRoutes.PATCH("/activity/:activityId", 
		middleware.ValidateJSONForNulls([]string{"activityType", "doneAt", "durationInMinutes"}),
		h.UpdateActivity)
		
	protectedRoutes.DELETE("/activity/:activityId", h.DeleteActivity)
}

func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := uint(userIDInterface.(int64))

	var req models.CreateActivityRequest
	
	// Handle JSON binding errors (empty body, malformed JSON)
	err := c.ShouldBindJSON(&req)
	if middleware.HandleValidationError(c, err) {
		return
	}

	// Get validator from context
	validate, exists := c.Get("validator")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Validation service unavailable"})
		return
	}

	// Validate the struct
	if validationErrors := middleware.ValidateStruct(validate.(*validator.Validate), req); validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	ctx := c.Request.Context()
	response, err := h.ActivitySvc.CreateActivity(ctx, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *ActivityHandler) GetActivities(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := uint(userIDInterface.(int64))

	query := models.GetActivitiesQuery{}
	
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			query.Limit = limit
		}
	}
	
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			query.Offset = offset
		}
	}

	if activityType := c.Query("activityType"); activityType != "" {
		if _, exists := models.ActivityTypeCalories[activityType]; exists {
			query.ActivityType = activityType
		}
	}

	if doneAtFromStr := c.Query("doneAtFrom"); doneAtFromStr != "" {
		if doneAtFrom, err := time.Parse(time.RFC3339, doneAtFromStr); err == nil {
			query.DoneAtFrom = doneAtFrom
		}
	}

	if doneAtToStr := c.Query("doneAtTo"); doneAtToStr != "" {
		if doneAtTo, err := time.Parse(time.RFC3339, doneAtToStr); err == nil {
			query.DoneAtTo = doneAtTo
		}
	}

	if caloriesBurnedMinStr := c.Query("caloriesBurnedMin"); caloriesBurnedMinStr != "" {
		if caloriesBurnedMin, err := strconv.Atoi(caloriesBurnedMinStr); err == nil && caloriesBurnedMin >= 0 {
			query.CaloriesBurnedMin = caloriesBurnedMin
		}
	}

	if caloriesBurnedMaxStr := c.Query("caloriesBurnedMax"); caloriesBurnedMaxStr != "" {
		if caloriesBurnedMax, err := strconv.Atoi(caloriesBurnedMaxStr); err == nil && caloriesBurnedMax >= 0 {
			query.CaloriesBurnedMax = caloriesBurnedMax
		}
	}

	ctx := c.Request.Context()
	activities, err := h.ActivitySvc.GetActivities(ctx, userID, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get activities"})
		return
	}

	c.JSON(http.StatusOK, activities)
}

func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := uint(userIDInterface.(int64))
	activityID := c.Param("activityId")

	var req models.UpdateActivityRequest
	
	// Handle JSON binding errors (empty body, malformed JSON)
	err := c.ShouldBindJSON(&req)
	if middleware.HandleValidationError(c, err) {
		return
	}

	// Get validator from context
	validate, exists := c.Get("validator")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Validation service unavailable"})
		return
	}

	// Custom validation for update requests
	if err := h.validateUpdateRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the struct
	if validationErrors := middleware.ValidateStruct(validate.(*validator.Validate), req); validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	ctx := c.Request.Context()
	response, err := h.ActivitySvc.UpdateActivity(ctx, userID, activityID, req)
	if err != nil {
		if errors.Is(err, customErrors.ErrActivityNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update activity"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ActivityHandler) validateUpdateRequest(req *models.UpdateActivityRequest) error {
	if req.ActivityType != nil && *req.ActivityType == "" {
		return errors.New("activityType cannot be empty string")
	}
	
	if req.DoneAt != nil {
		if *req.DoneAt == "" {
			return errors.New("doneAt cannot be empty string")
		}
		if _, err := time.Parse(time.RFC3339, *req.DoneAt); err != nil {
			return errors.New("doneAt must be a valid ISO date")
		}
	}
	
	if req.DurationInMinutes != nil && *req.DurationInMinutes <= 0 {
		return errors.New("durationInMinutes must be greater than 0")
	}
	
	return nil
}

func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := uint(userIDInterface.(int64))
	activityID := c.Param("activityId")

	ctx := c.Request.Context()
	err := h.ActivitySvc.DeleteActivity(ctx, userID, activityID)
	if err != nil {
		if errors.Is(err, customErrors.ErrActivityNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete activity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity deleted successfully"})
}
