package service

import (
	"FitByte/configs"
	customErrors "FitByte/internal/errors"
	"FitByte/internal/models"
	"FitByte/internal/repositories"
	"FitByte/pkg/log"
	"FitByte/pkg/token"
	"context"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type ProfileService interface {
	Register(ctx context.Context, authRequest models.AuthRequest) (models.RegisterResponse, error)
	Login(ctx context.Context, authRequest models.AuthRequest) (string, error)
	UpdateUserProfile(ctx context.Context, userID uint, updates map[string]interface{}) error
	GetProfile(ctx context.Context, userID uint) (*models.Profile, error)
}

type profileService struct {
	appConfig   configs.Config
	profileRepo repositories.ProfileRepository
}

func NewProfileService(appConfig configs.Config, profileRepo repositories.ProfileRepository) ProfileService {
	return &profileService{
		appConfig:   appConfig,
		profileRepo: profileRepo,
	}
}

func (u *profileService) Register(ctx context.Context, authRequest models.AuthRequest) (models.RegisterResponse, error) {
	isUserExist, err := u.profileRepo.GetProfileByEmail(ctx, authRequest.Email)
	if err != nil {
		log.Logger.Error().Err(err).Msg("error occurred on Register(ctx context.Context, authRequest models.AuthRequest")
		return models.RegisterResponse{}, err
	}

	if isUserExist != nil {
		log.Logger.Warn().Str("email", authRequest.Email).Msg("user already exists")
		return models.RegisterResponse{}, customErrors.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(authRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Logger.Error().Err(err).Msg("error occurred on Register(ctx context.Context, authRequest models.AuthRequest")
		return models.RegisterResponse{}, err
	}

	userProfile := models.Profile{
		Email:    authRequest.Email,
		Password: string(hashedPassword),
	}

	err = u.profileRepo.CreateUser(ctx, userProfile)
	if err != nil {
		log.Logger.Error().Err(err).Msg("error occurred on Register(ctx context.Context, authRequest models.AuthRequest")
		return models.RegisterResponse{}, err
	}

	signedToken, err := token.GenerateJWTToken(userProfile.ID, userProfile.Email, u.appConfig.Secret.JWTSecret)
	if err != nil {
		log.Logger.Error().Err(err).Msg("error occurred on Register(ctx context.Context, authRequest models.AuthRequest")
		return models.RegisterResponse{}, err
	}

	return models.RegisterResponse{
		Token: signedToken,
		Email: userProfile.Email,
	}, nil
}

func (u *profileService) Login(ctx context.Context, authRequest models.AuthRequest) (string, error) {
	userDetail, err := u.profileRepo.GetProfileByEmail(ctx, authRequest.Email)
	if err != nil {
		log.Logger.Error().Err(err).Msg("error occurred on Login(ctx context.Context, authRequest models.AuthRequest")
		return "", err
	}

	if userDetail == nil {
		log.Logger.Warn().Str("email", authRequest.Email).Msg("user not found")
		return "", customErrors.ErrorUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDetail.Password), []byte(authRequest.Password))
	if err != nil {
		log.Logger.Warn().Str("email", authRequest.Email).Msg("invalid password")
		return "", customErrors.ErrInvalidCredentials
	}

	signedToken, err := token.GenerateJWTToken(userDetail.ID, userDetail.Email, u.appConfig.Secret.JWTSecret)
	if err != nil {
		log.Logger.Error().Err(err).Msg("error occurred on Login(ctx context.Context, authRequest models.AuthRequest")
		return "", err
	}

	return signedToken, nil
}

func (u *profileService) UpdateUserProfile(ctx context.Context, userID uint, updates map[string]interface{}) error {
	userProfile, err := u.profileRepo.GetProfileByID(ctx, userID)
	if userProfile == nil {
		log.Logger.Warn().Str("userID", strconv.FormatUint(uint64(userID), 10)).Msg("user not found")
		return customErrors.ErrorUserNotFound
	}

	err = u.profileRepo.UpdateUser(ctx, userID, updates)
	if err != nil {
		log.Logger.Error().Err(err).Msg("error occurred on UpdateUserProfile(ctx context.Context, email string, updates map[string]interface{})")
		return err
	}

	return nil
}

func (u *profileService) GetProfile(ctx context.Context, userID uint) (*models.Profile, error) {
	return u.profileRepo.GetProfileByID(ctx, userID)
}