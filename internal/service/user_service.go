package service

import (
	"FitByte/configs"
	customErrors "FitByte/internal/errors"
	"FitByte/internal/models"
	"FitByte/internal/repositories"
	"FitByte/pkg/log"
	"FitByte/pkg/token"
	"context"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, user models.User) (models.RegisterResponse, error)
	Login(ctx context.Context, user models.User) (string, error)
}

type userService struct {
	appConfig configs.Config
	userRepo  repositories.UserRepository
}

func NewUserService(appConfig configs.Config, userRepo repositories.UserRepository) UserService {
	return &userService{
		appConfig: appConfig,
		userRepo:  userRepo,
	}
}

func (u *userService) Register(ctx context.Context, user models.User) (models.RegisterResponse, error) {
	isUserExist, err := u.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		log.Logger.Error().Err(err).Msg("error occurred on Register(ctx context.Context, user models.User)")
		return models.RegisterResponse{}, err
	}

	if isUserExist != nil {
		log.Logger.Warn().Str("email", user.Email).Msg("user already exists")
		return models.RegisterResponse{}, customErrors.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Logger.Error().Err(err).Msg("error occurred on Register(ctx context.Context, user models.User)")
		return models.RegisterResponse{}, err
	}

	user.Password = string(hashedPassword)
	err = u.userRepo.CreateUser(ctx, user)
	if err != nil {
		log.Logger.Error().Err(err).Msg("error occurred on Register(ctx context.Context, user models.User)")
		return models.RegisterResponse{}, err
	}

	signedToken, err := token.GenerateJWTToken(user.ID, user.Email, u.appConfig.Secret.JWTSecret)
	if err != nil {
		log.Logger.Error().Err(err).Msg("error occurred on Register(ctx context.Context, user models.User)")
		return models.RegisterResponse{}, err
	}

	return models.RegisterResponse{
		Token: signedToken,
		Email: user.Email,
	}, nil
}

func (u *userService) Login(ctx context.Context, user models.User) (string, error) {
	userDetail, err := u.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		log.Logger.Error().Err(err).Msg("error occurred on Login(ctx context.Context, user models.User)")
		return "", err
	}

	if userDetail == nil {
		log.Logger.Warn().Str("email", user.Email).Msg("user not found")
		return "", customErrors.ErrorUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDetail.Password), []byte(user.Password))
	if err != nil {
		log.Logger.Warn().Str("email", user.Email).Msg("invalid password")
		return "", customErrors.ErrInvalidCredentials
	}

	signedToken, err := token.GenerateJWTToken(userDetail.ID, userDetail.Email, u.appConfig.Secret.JWTSecret)
	if err != nil {
		log.Logger.Error().Err(err).Msg("error occurred on Login(ctx context.Context, user models.User)")
		return "", err
	}

	return signedToken, nil
}
