package main

import (
	"FitByte/configs"
	"FitByte/internal/handlers"
	"FitByte/internal/infrastructure"
	"FitByte/internal/repositories"
	"FitByte/internal/service"

	"github.com/gin-gonic/gin"

	"FitByte/pkg/log"
)

func main() {
	log.InitLogger()

	appConfig := configs.LoadConfig(
		configs.WithConfigFolder([]string{"./configs"}),
		configs.WithConfigFile("config"),
		configs.WithConfigType("yaml"),
	)

	db := infrastructure.InitDB(appConfig)

	r := gin.Default()
	r.Use(gin.Recovery())

	userRepo := repositories.NewUserRepository(db)
	userService := service.NewUserService(appConfig, userRepo)
	userHandler := handlers.NewUserHandler(r, appConfig, userService)
	userHandler.SetupRoutes()

	log.Logger.Info().Str("port", appConfig.App.Port).Msg("Starting server")
	if err := r.Run(":" + appConfig.App.Port); err != nil {
		log.Logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
