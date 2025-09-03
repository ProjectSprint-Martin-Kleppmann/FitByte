package token

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateJWTToken(userID uint, email, secretKey string) (string, error) {
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
	})

	token, err := tokenString.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}
