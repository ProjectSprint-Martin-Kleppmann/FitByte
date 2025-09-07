package handlers

import (
	"FitByte/configs"
	customErrors "FitByte/internal/errors"
	"FitByte/internal/middleware"
	"FitByte/internal/models"
	"FitByte/internal/service"
	"errors"
	"net/http"
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
