package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const UserIDContextKey = "user_id"

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userIDStr string

		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				userIDFromToken, err := extractUserIDFromToken(token)
				if err == nil {
					userIDStr = userIDFromToken
				}
			}
		}

		if userIDStr == "" {
			userIDStr = c.GetHeader("X-User-ID")
		}

		if userIDStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
			c.Abort()
			return
		}

		c.Set(UserIDContextKey, userID)
		c.Next()
	}
}

func extractUserIDFromToken(token string) (string, error) {
	return "", nil
}

