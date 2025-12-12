package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appGame "sigame/game/internal/application/game"
	domainGame "sigame/game/internal/domain/game"
	"sigame/game/internal/domain/player"
	"sigame/game/internal/port"
	"sigame/game/internal/transport/http/middleware"
	"sigame/game/internal/transport/ws/hub"
)

type GameHandler struct {
	packService    port.PackService
	gameRepository port.GameRepository
	gameCache      port.GameCache
	hub            *hub.Hub
	eventLogger    port.EventLogger
}

func NewGameHandler(packService port.PackService, gameRepository port.GameRepository, gameCache port.GameCache, hub *hub.Hub, eventLogger port.EventLogger) *GameHandler {
	return &GameHandler{
		packService:    packService,
		gameRepository: gameRepository,
		gameCache:      gameCache,
		hub:            hub,
		eventLogger:    eventLogger,
	}
}

func (h *GameHandler) CreateGame(c *gin.Context) {
	var req CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hasHost := false
	for _, playerInfo := range req.Players {
		role := player.Role(playerInfo.Role)
		if role != player.RoleHost && role != player.RolePlayer {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrorInvalidRole})
			return
		}
		if role == player.RoleHost {
			hasHost = true
		}
	}

	if !hasHost {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrorHostRequired})
		return
	}

	settings := domainGame.Settings{
		TimeForAnswer: req.Settings.TimeForAnswer,
		TimeForChoice: req.Settings.TimeForChoice,
	}

	if err := settings.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrorInvalidSettings})
		return
	}

	pack, err := h.packService.GetPackContent(c.Request.Context(), req.PackID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrorPackNotFound})
		return
	}

	game := domainGame.New(req.RoomID, req.PackID, settings, pack.Rounds)

	for _, playerInfo := range req.Players {
		role := player.Role(playerInfo.Role)
		p := player.New(playerInfo.UserID, playerInfo.Username, playerInfo.AvatarURL, role)
		if err := game.AddPlayer(p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrorPlayerAlreadyExists})
			return
		}
	}

	if err := h.gameRepository.CreateGameSession(c.Request.Context(), game); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrorFailedToCreateGame})
		return
	}

	if err := h.gameCache.SaveGameState(c.Request.Context(), game); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrorFailedToSaveGameState})
		return
	}

	manager := appGame.New(game, pack, h.hub, h.eventLogger, h.gameRepository, h.gameCache)
	h.hub.RegisterGameManager(game.ID, manager)
	manager.Start()

	gameID := game.ID

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

	game, err := h.gameRepository.GetGameSession(c.Request.Context(), gameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrorGameNotFound})
		return
	}

	players := make([]PlayerState, 0, len(game.Players))
	for _, p := range game.Players {
		players = append(players, PlayerState{
			UserID:      p.UserID,
			Username:    p.Username,
			AvatarURL:   p.AvatarURL,
			Role:        string(p.Role),
			Score:       p.Score,
			IsActive:    p.IsActive,
			IsReady:     p.IsReady,
			IsConnected: p.IsConnected,
		})
	}

	c.JSON(http.StatusOK, GetGameResponse{
		GameID:       game.ID,
		RoomID:       game.RoomID,
		PackID:       game.PackID,
		Status:       string(game.Status),
		CurrentRound: game.CurrentRound,
		Players:      players,
		Settings: GameSettings{
			TimeForAnswer: game.Settings.TimeForAnswer,
			TimeForChoice: game.Settings.TimeForChoice,
		},
	})
}

func (h *GameHandler) GetMyActiveGame(c *gin.Context) {
	userIDValue, exists := c.Get(middleware.UserIDContextKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrorUserIDRequired})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrorInvalidUserID})
		return
	}

	game, err := h.gameRepository.GetActiveGameForUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"hasActiveGame": false,
		})
		return
	}

	players := make([]PlayerState, 0, len(game.Players))
	for _, p := range game.Players {
		players = append(players, PlayerState{
			UserID:      p.UserID,
			Username:    p.Username,
			AvatarURL:   p.AvatarURL,
			Role:        string(p.Role),
			Score:       p.Score,
			IsActive:    p.IsActive,
			IsReady:     p.IsReady,
			IsConnected: p.IsConnected,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"hasActiveGame": true,
		"game": GetGameResponse{
			GameID:       game.ID,
			RoomID:       game.RoomID,
			PackID:       game.PackID,
			Status:       string(game.Status),
			CurrentRound: game.CurrentRound,
			Players:      players,
			Settings: GameSettings{
				TimeForAnswer: game.Settings.TimeForAnswer,
				TimeForChoice: game.Settings.TimeForChoice,
			},
		},
	})
}

