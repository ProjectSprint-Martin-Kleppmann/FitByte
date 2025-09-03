package middleware

import (
	"FitByte/pkg/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

type AppClaims struct {
	UserID float64 `json:"user_id"`
	jwt.RegisteredClaims
}

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		partedHeader := strings.Split(authHeader, " ")
		if len(partedHeader) != 2 || partedHeader[0] != "Bearer" || partedHeader[1] == "" {
			log.Logger.Error().Msg("Unauthorized: Invalid Authorization header")
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		tokenString := partedHeader[1]
		claims := &AppClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			log.Logger.Error().Err(err).Msg("Unauthorized: Invalid token")
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		log.Logger.Info().Float64("user_id", claims.UserID).Msg("Authenticated request")
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
