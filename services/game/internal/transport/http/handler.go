package http

import (
	"sigame/game/internal/adapter/grpc/pack"
	"sigame/game/internal/adapter/repository/postgres"
	"sigame/game/internal/adapter/repository/redis"
	"sigame/game/internal/port"
	"sigame/game/internal/transport/http/handler"
	"sigame/game/internal/transport/ws/hub"
)

type Handler struct {
	Game   *handler.GameHandler
	Health *handler.HealthHandler
}

func NewHandler(packService port.PackService, gameRepository port.GameRepository, gameCache port.GameCache, hub *hub.Hub, eventLogger port.EventLogger, pgClient *postgres.Client, redisClient *redis.Client, packClient *pack.PackClient) *Handler {
	return &Handler{
		Game:   handler.NewGameHandler(packService, gameRepository, gameCache, hub, eventLogger),
		Health: handler.NewHealthHandler(pgClient, redisClient, packClient),
	}
}

