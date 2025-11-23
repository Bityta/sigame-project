package rest

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sigame/game/internal/domain"
	"github.com/sigame/game/internal/game"
	grpcClient "github.com/sigame/game/internal/grpc"
	"github.com/sigame/game/internal/repository/postgres"
	"github.com/sigame/game/internal/repository/redis"
	"github.com/sigame/game/internal/transport/websocket"
)

// Handler handles REST API requests
type Handler struct {
	hub            *websocket.Hub
	packClient     *grpcClient.PackClient
	gameRepo       *postgres.GameRepository
	eventRepo      *postgres.EventRepository
	redisGameRepo  *redis.GameRepository
	redisCacheRepo *redis.CacheRepository
}

// NewHandler creates a new REST handler
func NewHandler(
	hub *websocket.Hub,
	packClient *grpcClient.PackClient,
	gameRepo *postgres.GameRepository,
	eventRepo *postgres.EventRepository,
	redisGameRepo *redis.GameRepository,
	redisCacheRepo *redis.CacheRepository,
) *Handler {
	return &Handler{
		hub:            hub,
		packClient:     packClient,
		gameRepo:       gameRepo,
		eventRepo:      eventRepo,
		redisGameRepo:  redisGameRepo,
		redisCacheRepo: redisCacheRepo,
	}
}

// Health returns service health status
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":       "healthy",
		"service":      "game-service",
		"timestamp":    time.Now().UTC(),
		"active_games": h.hub.GetActiveGamesCount(),
	})
}

// CreateGame creates a new game session
func (h *Handler) CreateGame(c *gin.Context) {
	var req domain.CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// Validate pack exists
	packExists, err := h.packClient.ValidatePackExists(ctx, req.PackID)
	if err != nil {
		log.Printf("Failed to validate pack: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate pack"})
		return
	}

	if !packExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pack not found"})
		return
	}

	// Get pack content (with caching)
	pack, err := h.getPackWithCache(ctx, req.PackID)
	if err != nil {
		log.Printf("Failed to get pack: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load pack"})
		return
	}

	// Create game
	gameID := uuid.New()
	gameSession := &domain.Game{
		ID:           gameID,
		RoomID:       req.RoomID,
		PackID:       req.PackID,
		Status:       domain.GameStatusWaiting,
		Players:      make(map[uuid.UUID]*domain.Player),
		Rounds:       pack.Rounds,
		CurrentRound: 0,
		CurrentPhase: domain.GameStatusWaiting,
		Settings:     req.Settings,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Add players
	for _, playerInfo := range req.Players {
		role := domain.PlayerRole(playerInfo.Role)
		player := domain.NewPlayer(playerInfo.UserID, playerInfo.Username, role)
		gameSession.Players[playerInfo.UserID] = player
	}

	// Save to PostgreSQL
	if err := h.gameRepo.CreateGameSession(ctx, gameSession); err != nil {
		log.Printf("Failed to create game session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create game"})
		return
	}

	// Save to Redis
	if err := h.redisGameRepo.SaveGameState(ctx, gameSession); err != nil {
		log.Printf("Failed to save game state to Redis: %v", err)
	}

	// Create game manager
	manager := game.NewManager(gameSession, pack, h.hub, h.eventRepo)

	// Register manager with hub
	h.hub.RegisterGameManager(gameID, manager)

	// Start game manager
	manager.Start()

	log.Printf("Game created: %s for room %s", gameID, req.RoomID)

	// Return response
	wsURL := "/api/game/" + gameID.String() + "/ws"
	resp := domain.CreateGameResponse{
		GameID:       gameID,
		WebSocketURL: wsURL,
		Status:       "created",
	}

	c.JSON(http.StatusCreated, resp)
}

// GetGame retrieves game information
func (h *Handler) GetGame(c *gin.Context) {
	gameIDStr := c.Param("id")
	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	ctx := c.Request.Context()

	// Try to get from Redis first
	game, err := h.redisGameRepo.LoadGameState(ctx, gameID)
	if err != nil {
		// Fall back to PostgreSQL
		game, err = h.gameRepo.GetGameSession(ctx, gameID)
		if err != nil {
			if err == domain.ErrGameNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get game"})
			}
			return
		}
	}

	// Build response
	players := make([]domain.PlayerState, 0, len(game.Players))
	for _, player := range game.Players {
		players = append(players, player.ToState())
	}

	resp := domain.GetGameResponse{
		GameID:       game.ID,
		RoomID:       game.RoomID,
		PackID:       game.PackID,
		Status:       game.Status,
		CurrentRound: game.CurrentRound,
		Players:      players,
		StartedAt:    game.StartedAt,
		FinishedAt:   game.FinishedAt,
	}

	c.JSON(http.StatusOK, resp)
}

// getPackWithCache retrieves pack with Redis caching
func (h *Handler) getPackWithCache(ctx context.Context, packID uuid.UUID) (*domain.Pack, error) {
	// Try cache first
	pack, err := h.redisCacheRepo.GetCachedPack(ctx, packID)
	if err == nil {
		return pack, nil
	}

	// Cache miss, fetch from Pack Service
	pack, err = h.packClient.GetPackContent(ctx, packID)
	if err != nil {
		return nil, err
	}

	// Cache it
	if err := h.redisCacheRepo.CachePack(ctx, pack); err != nil {
		log.Printf("Failed to cache pack: %v", err)
	}

	return pack, nil
}

