package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	authClient "sigame/game/internal/adapter/grpc/auth"
	"sigame/game/internal/infrastructure/logger"
	"sigame/game/internal/transport/ws/client"
	"sigame/game/internal/transport/ws/hub"
)

type AuthService interface {
	ValidateToken(ctx context.Context, token string) (*authClient.ValidateTokenResponse, error)
}

type Handler struct {
	hub        *hub.Hub
	authClient AuthService
}

func NewHandler(h *hub.Hub, authClient AuthService) *Handler {
	return &Handler{
		hub:        h,
		authClient: authClient,
	}
}

func (h *Handler) HandleWebSocket(c *gin.Context) {
	ctx := c.Request.Context()
	logger.Debugf(ctx, "[WS] HandleWebSocket called, path=%s", c.Request.URL.Path)

	gameIDStr := c.Param("id")
	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		logger.Errorf(ctx, "[WS] Invalid game ID: %s", gameIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrorInvalidGameID})
		return
	}

	token := c.Query(QueryParamToken)
	var userID uuid.UUID

	if token != "" {
		if h.authClient == nil {
			logger.Warnf(ctx, "[WS] Auth client not initialized, token validation skipped")
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrorTokenRequired})
			return
		}

		resp, err := h.authClient.ValidateToken(ctx, token)
		if err != nil {
			logger.Errorf(ctx, "[WS] Token validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrorInvalidToken})
			return
		}

		if !resp.Valid {
			logger.Warnf(ctx, "[WS] Invalid token: %s", resp.Error)
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrorInvalidToken})
			return
		}

		userID = resp.UserID
		logger.Debugf(ctx, "[WS] Token validated, user_id=%s", userID)
	} else {
		userIDStr := c.Query(QueryParamUserID)
		if userIDStr == "" {
			logger.Errorf(ctx, "[WS] Missing token or user_id for game %s", gameID)
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrorTokenRequired})
			return
		}

		parsedUserID, err := uuid.Parse(userIDStr)
		if err != nil {
			logger.Errorf(ctx, "[WS] Invalid user ID: %s", userIDStr)
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrorInvalidUserID})
			return
		}
		userID = parsedUserID
	}

	logger.Debugf(ctx, "[WS] Checking game manager for game=%s, user=%s", gameID, userID)

	if _, exists := h.hub.GetGameManager(gameID); !exists {
		logger.Errorf(ctx, "[WS] Game manager not found for game %s", gameID)
		c.JSON(http.StatusNotFound, gin.H{"error": ErrorGameNotFound})
		return
	}

	logger.Debugf(ctx, "[WS] Game manager found, upgrading connection for game=%s, user=%s", gameID, userID)

	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf(ctx, "[WS] Failed to upgrade connection: %v", err)
		return
	}

	logger.Debugf(ctx, "[WS] Connection upgraded, creating client for game=%s, user=%s", gameID, userID)

	cl := client.NewClient(h.hub, conn, userID, gameID)
	h.hub.Register(cl)
	cl.Run()

	logger.Infof(ctx, "[WS] WebSocket connection established: user=%s, game=%s", userID, gameID)
}

