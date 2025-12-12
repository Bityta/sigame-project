package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	authClient "sigame/game/internal/adapter/grpc/auth"
	"sigame/game/internal/infrastructure/logger"
)

const UserIDContextKey = "user_id"

type AuthService interface {
	ValidateToken(ctx context.Context, token string) (*authClient.ValidateTokenResponse, error)
}

var authClientInstance AuthService

func SetAuthClient(client AuthService) {
	authClientInstance = client
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userID uuid.UUID
		var err error

		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				validatedUserID, err := extractUserIDFromToken(c.Request.Context(), token)
				if err == nil && validatedUserID != uuid.Nil {
					userID = validatedUserID
				}
			}
		}

		if userID == uuid.Nil {
			userIDStr := c.GetHeader("X-User-ID")
			if userIDStr != "" {
				userID, err = uuid.Parse(userIDStr)
				if err != nil {
					logger.Warnf(c.Request.Context(), "Invalid X-User-ID header: %v", err)
				}
			}
		}

		if userID == uuid.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
			c.Abort()
			return
		}

		c.Set(UserIDContextKey, userID)
		c.Next()
	}
}

func extractUserIDFromToken(ctx context.Context, token string) (uuid.UUID, error) {
	if authClientInstance == nil {
		logger.Warnf(ctx, "Auth client not initialized, skipping token validation")
		return uuid.Nil, nil
	}

	resp, err := authClientInstance.ValidateToken(ctx, token)
	if err != nil {
		logger.Warnf(ctx, "Token validation failed: %v", err)
		return uuid.Nil, err
	}

	if !resp.Valid {
		logger.Warnf(ctx, "Token is invalid: %s", resp.Error)
		return uuid.Nil, nil
	}

	return resp.UserID, nil
}

