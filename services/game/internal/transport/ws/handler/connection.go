package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sigame/game/internal/infrastructure/logger"
	"sigame/game/internal/transport/ws/client"
	"sigame/game/internal/transport/ws/hub"
)

const (
	QueryParamUserID         = "user_id"
	ErrorInvalidGameID       = "Invalid game ID"
	ErrorUserIDRequired      = "user_id is required"
	ErrorInvalidUserID       = "Invalid user ID"
	ErrorGameNotFound        = "Game not found or not started"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	hub *hub.Hub
}

func NewHandler(h *hub.Hub) *Handler {
	return &Handler{hub: h}
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

	userIDStr := c.Query(QueryParamUserID)
	if userIDStr == "" {
		logger.Errorf(ctx, "[WS] Missing user_id for game %s", gameID)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrorUserIDRequired})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.Errorf(ctx, "[WS] Invalid user ID: %s", userIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrorInvalidUserID})
		return
	}

	logger.Debugf(ctx, "[WS] Checking game manager for game=%s, user=%s", gameID, userID)

	if _, exists := h.hub.GetGameManager(gameID); !exists {
		logger.Errorf(ctx, "[WS] Game manager not found for game %s", gameID)
		c.JSON(http.StatusNotFound, gin.H{"error": ErrorGameNotFound})
		return
	}

	logger.Debugf(ctx, "[WS] Game manager found, upgrading connection for game=%s, user=%s", gameID, userID)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
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

