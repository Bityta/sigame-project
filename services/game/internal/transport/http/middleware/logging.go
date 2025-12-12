package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"sigame/game/internal/infrastructure/logger"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		logger.Infof(c.Request.Context(), "%s %s %d %v", method, path, statusCode, latency)
	}
}

