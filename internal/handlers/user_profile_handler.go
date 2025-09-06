package handlers

import (
	"FitByte/configs"
	"FitByte/internal/middleware"
	"FitByte/internal/models"
	"FitByte/internal/service"
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
