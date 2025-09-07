package handlers

import (
	"FitByte/configs"
	"FitByte/internal/middleware"
	"FitByte/internal/models"
	"FitByte/internal/service"
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	customErrors "FitByte/internal/errors"
)

type ProfileHandler struct {
	Engine     *gin.Engine
	AppConfig  configs.Config
	ProfileSvc service.ProfileService
}

func NewProfileHandler(engine *gin.Engine, appConfig configs.Config, profileService service.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		Engine:     engine,
		AppConfig:  appConfig,
		ProfileSvc: profileService,
	}
}

func (h *ProfileHandler) SetupRoutes() {
	// Health check
	h.Engine.GET("/ping", h.pong)

	routes := h.Engine.Group("/v1")
	routes.POST("register", h.Register)
	routes.POST("login", h.Login)

	protectedRoutes := h.Engine.Group("/v1")
	protectedRoutes.Use(middleware.AuthMiddleware(h.AppConfig.Secret.JWTSecret))
	protectedRoutes.PATCH("/user", h.UpdateProfile) 

	// Protected routes
	privateRoutes := h.Engine.Group("/health")
	privateRoutes.Use(middleware.AuthMiddleware(h.AppConfig.Secret.JWTSecret))
	privateRoutes.GET("private-ping", h.pong)
}

func (h *ProfileHandler) pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (h *ProfileHandler) Register(c *gin.Context) {
	var model models.AuthRequest
	ctx := c.Request.Context()
	err := c.ShouldBindJSON(&model)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.ProfileSvc.Register(ctx, model)
	if err != nil {
		if errors.Is(customErrors.ErrUserAlreadyExists, err) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"email": response.Email, "token": response.Token})
}

func (h *ProfileHandler) Login(c *gin.Context) {
	var model models.AuthRequest
	ctx := c.Request.Context()

	err := c.ShouldBindJSON(&model)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.ProfileSvc.Login(ctx, model)
	if err != nil {
		if errors.Is(err, customErrors.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		} else if errors.Is(err, customErrors.ErrorUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email": model.Email,
		"token": token,
	})
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := uint(userIDInterface.(int64))

	var req models.PatchProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fieldMapping := map[string]struct {
		dbField string
		value   interface{}
	}{
		"Preference":  {"preference", req.Preference},
		"WeightUnit":  {"weight_unit", req.WeightUnit},
		"HeightUnit":  {"height_unit", req.HeightUnit},
		"Weight":      {"weight", req.Weight},
		"Height":      {"height", req.Height},
		"Name":        {"name", req.Name},
		"ImageUri":    {"image_uri", req.ImageURI},
	}

	updates := make(map[string]interface{})
	for _, mapping := range fieldMapping {
		if mapping.value != nil {
			updates[mapping.dbField] = mapping.value
		}
	}

	ctx := context.Background()
	if err := h.ProfileSvc.UpdateUserProfile(ctx, userID, updates); err != nil {
		if errors.Is(err, customErrors.ErrorUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"preference":  getStringValue(req.Preference),
		"weightUnit":  getStringValue(req.WeightUnit),
		"heightUnit":  getStringValue(req.HeightUnit),
		"weight":      getFloat64Value(req.Weight),
		"height":      getFloat64Value(req.Height),
		"name":        getStringValue(req.Name),
		"imageUri":    getStringValue(req.ImageURI),
	})
}

// helper
func getStringValue(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

func getFloat64Value(ptr *float64) float64 {
	if ptr != nil {
		return *ptr
	}
	return 0
}