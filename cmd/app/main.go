package main

import (
	"FitByte/configs"
	"FitByte/internal/handlers"
	"FitByte/internal/infrastructure"
	"FitByte/internal/middleware"
	"FitByte/internal/repositories"
	"FitByte/internal/service"

	"FitByte/pkg/log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.InitLogger()

	appConfig := configs.LoadConfig(
		configs.WithConfigFolder([]string{"./configs"}),
		configs.WithConfigFile("config"),
		configs.WithConfigType("yaml"),
	)

	db := infrastructure.InitDB(appConfig)
	minioClient := infrastructure.InitMinioStorage(appConfig)

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())

	userRepo := repositories.NewUserRepository(db)
	userService := service.NewUserService(appConfig, userRepo)
	userHandler := handlers.NewUserHandler(r, appConfig, userService)
	userHandler.SetupRoutes()

	minioRepo := repositories.NewMinioRepository(minioClient, appConfig.Minio.Bucket)
	fileRepo := repositories.NewFileRepository(db)
	fileService := service.NewFileService(fileRepo, minioRepo)
	fileHandler := handlers.NewFileHandler(r, appConfig, fileService)
	fileHandler.SetupRoutes()

	log.Logger.Info().Str("port", appConfig.App.Port).Msg("Starting server")
	if err := r.Run(":" + appConfig.App.Port); err != nil {
		log.Logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
