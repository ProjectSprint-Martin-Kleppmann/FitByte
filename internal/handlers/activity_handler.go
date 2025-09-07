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

	// "strconv"
	// "time"

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

	protectedRoutes.POST("/activity", h.CreateActivity)
	protectedRoutes.PATCH("/activity/:activityId", h.UpdateActivity)
	protectedRoutes.DELETE("/activity/:activityId", h.DeleteActivity)
	protectedRoutes.GET("/activity", h.GetActivities)
}

func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := uint(userIDInterface.(int64))

	var req models.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the request
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := uint(userIDInterface.(int64))
	activityID := c.Param("activityId")

	var req models.UpdateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the request
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func (h *ActivityHandler) GetActivities(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDInterface.(int64))

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil || limit < 0 {
		limit = 5
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	params := models.GetActivityParams{
		Limit:        limit,
		Offset:       offset,
		ActivityType: c.Query("activityType"),
	}

	if doneAtFromStr := c.Query("doneAtFrom"); doneAtFromStr != "" {
		if t, err := time.Parse(time.RFC3339, doneAtFromStr); err == nil {
			params.DoneAtFrom = t
		}
	}
	if doneAtToStr := c.Query("doneAtTo"); doneAtToStr != "" {
		if t, err := time.Parse(time.RFC3339, doneAtToStr); err == nil {
			params.DoneAtTo = t
		}
	}

	if calMinStr := c.Query("caloriesBurnedMin"); calMinStr != "" {
		if cal, err := strconv.Atoi(calMinStr); err == nil && cal > 0 {
			params.CaloriesBurnedMin = cal
		}
	}
	if calMaxStr := c.Query("caloriesBurnedMax"); calMaxStr != "" {
		if cal, err := strconv.Atoi(calMaxStr); err == nil && cal > 0 {
			params.CaloriesBurnedMax = cal
		}
	}

	ctx := c.Request.Context()
	activities, err := h.ActivitySvc.GetActivities(ctx, userID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve activities"})
		return
	}

	c.JSON(http.StatusOK, activities)
}
