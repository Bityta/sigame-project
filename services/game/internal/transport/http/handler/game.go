package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GameHandler struct{}

func NewGameHandler() *GameHandler {
	return &GameHandler{}
}

func (h *GameHandler) CreateGame(c *gin.Context) {
	var req CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gameID := uuid.New()
	
	c.JSON(http.StatusCreated, CreateGameResponse{
		GameID:       gameID,
		WebSocketURL: "/api/game/" + gameID.String() + "/ws",
		Status:       "created",
	})
}

func (h *GameHandler) GetGame(c *gin.Context) {
	gameIDStr := c.Param("id")
	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrorInvalidGameID})
		return
	}

	c.JSON(http.StatusOK, GetGameResponse{
		GameID: gameID,
	})
}

func (h *GameHandler) GetMyActiveGame(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hasActiveGame": false,
	})
}

