package token

import (
	"FitByte/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateJWTToken(userID uint, email, secretKey string) (string, error) {
	claims := &middleware.AppClaims{
		UserID: int64(userID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		},
	}

	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenString.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

