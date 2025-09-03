package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"FitByte/pkg/log"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		timeoutCtx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		ctx := context.WithValue(timeoutCtx, "requestID", requestID)
		c.Request = c.Request.WithContext(ctx)

		startTime := time.Now()
		c.Next()
		latency := time.Since(startTime)

		status := c.Writer.Status()
		logger := log.Logger.With().
			Str("request_id", requestID).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", status).
			Str("client_ip", c.ClientIP()).
			Dur("latency", latency).
			Logger()

		if userID, exists := c.Get("user_id"); exists {
			logger.Info().Interface("user_id", userID).Msg("request processed")
		}

		if status >= 500 {
			logger.Error().Msg("request completed with server error")
		} else if status >= 400 {
			logger.Warn().Msg("request completed with client error")
		} else {
			logger.Info().Msg("request completed successfully")
		}
	}
}
