package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sigame/game/internal/infrastructure/logger"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Errorf(c.Request.Context(), "Panic recovered: %v", recovered)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "INTERNAL_SERVER_ERROR",
			Message: "An unexpected error occurred",
		})
		c.Abort()
	})
}

func AbortWithError(c *gin.Context, statusCode int, errorCode, message string) {
	logger.Warnf(c.Request.Context(), "Request aborted: %s - %s", errorCode, message)
	c.JSON(statusCode, ErrorResponse{
		Error:   errorCode,
		Message: message,
	})
	c.Abort()
}

