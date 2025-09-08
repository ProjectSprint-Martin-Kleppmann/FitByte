package handlers

import (
	"FitByte/configs"
	"FitByte/internal/middleware"
	"FitByte/internal/models"
	"FitByte/internal/service"

	// "context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

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

	// Add validation middleware to all routes
	routes := h.Engine.Group("/v1")
	routes.Use(middleware.ContentTypeMiddleware())
	routes.Use(middleware.ValidationMiddleware())
	routes.POST("register", h.Register)
	routes.POST("login", h.Login)

	protectedRoutes := h.Engine.Group("/v1")
	protectedRoutes.Use(middleware.ContentTypeMiddleware())
	protectedRoutes.Use(middleware.ValidationMiddleware())
	protectedRoutes.Use(middleware.AuthMiddleware(h.AppConfig.Secret.JWTSecret))
	protectedRoutes.GET("/user", h.GetProfile)
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
	if middleware.HandleValidationError(c, err) {
		return
	}

	validate, exists := c.Get("validator")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Validation service unavailable"})
		return
	}

	if validationErrors := middleware.ValidateStruct(validate.(*validator.Validate), model); validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
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
	if middleware.HandleValidationError(c, err) {
		return
	}

	validate, exists := c.Get("validator")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Validation service unavailable"})
		return
	}

	if validationErrors := middleware.ValidateStruct(validate.(*validator.Validate), model); validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
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

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := uint(userIDInterface.(int64))
	ctx := c.Request.Context()

	profile, err := h.ProfileSvc.GetProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, customErrors.ErrorUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get profile"})
		return
	}

	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	response := gin.H{
		"email": profile.Email,
	}

	if profile.Preference == "" {
		response["preference"] = nil
	} else {
		response["preference"] = profile.Preference
	}

	if profile.WeightUnit == "" {
		response["weightUnit"] = nil
	} else {
		response["weightUnit"] = profile.WeightUnit
	}

	if profile.HeightUnit == "" {
		response["heightUnit"] = nil
	} else {
		response["heightUnit"] = profile.HeightUnit
	}

	if profile.Weight == 0 {
		response["weight"] = nil
	} else {
		response["weight"] = profile.Weight
	}

	if profile.Height == 0 {
		response["height"] = nil
	} else {
		response["height"] = profile.Height
	}

	if profile.Name == "" {
		response["name"] = nil
	} else {
		response["name"] = profile.Name
	}

	if profile.ImageURI == "" {
		response["imageUri"] = nil
	} else {
		response["imageUri"] = profile.ImageURI
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := uint(userIDInterface.(int64))

	var req models.PatchProfileRequest

	err := c.ShouldBindJSON(&req)
	if middleware.HandleValidationError(c, err) {
		return
	}

	validate, exists := c.Get("validator")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Validation service unavailable"})
		return
	}

	if validationErrors := middleware.ValidateStruct(validate.(*validator.Validate), req); validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	updates := map[string]interface{}{
		"preference":  req.Preference,
		"weight_unit": req.WeightUnit,
		"height_unit": req.HeightUnit,
		"weight":      req.Weight,
		"height":      req.Height,
		"name":        req.Name,
		"image_uri":   req.ImageURI,
	}

	ctx := c.Request.Context()
	if err := h.ProfileSvc.UpdateUserProfile(ctx, userID, updates); err != nil {
		if errors.Is(err, customErrors.ErrorUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	response := gin.H{
		"preference": req.Preference,
		"weightUnit": req.WeightUnit,
		"heightUnit": req.HeightUnit,
		"weight":     req.Weight,
		"height":     req.Height,
		"name":       req.Name,
		"imageUri":   req.ImageURI,
	}

	c.JSON(http.StatusOK, response)
}
