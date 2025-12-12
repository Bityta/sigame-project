package main

import (
	"github.com/gin-gonic/gin"
	"sigame/game/internal/infrastructure/config"
	grpcClient "sigame/game/internal/adapter/grpc/pack"
	authClient "sigame/game/internal/adapter/grpc/auth"
	"sigame/game/internal/adapter/repository/postgres"
	"sigame/game/internal/adapter/repository/redis"
	"sigame/game/internal/transport/http"
	"sigame/game/internal/transport/http/middleware"
	"sigame/game/internal/transport/ws"
)

func initPostgreSQL(cfg *config.Config) (*postgres.Client, error) {
	return postgres.NewClient(
		cfg.GetPostgresConnectionString(),
		cfg.Database.MaxConns,
		cfg.Database.MaxIdle,
	)
}

func initRedis(cfg *config.Config) (*redis.Client, error) {
	return redis.NewClient(
		cfg.GetRedisAddress(),
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
}

func initPackClient(cfg *config.Config) (*grpcClient.PackClient, error) {
	return grpcClient.NewPackClient(cfg.GetPackServiceAddress())
}

func initAuthClient(cfg *config.Config) (*authClient.AuthServiceClient, error) {
	client, err := authClient.NewAuthClient(cfg.GetAuthServiceAddress())
	if err != nil {
		return nil, err
	}
	// Set the client in middleware
	middleware.SetAuthClient(client)
	return client, nil
}

type Repositories struct {
	GameRepo      *postgres.GameRepository
	EventRepo     *postgres.EventRepository
	RedisGameRepo *redis.GameRepository
	RedisCacheRepo *redis.CacheRepository
}

func initRepositories(pgClient *postgres.Client, redisClient *redis.Client) *Repositories {
	return &Repositories{
		GameRepo:      postgres.NewGameRepository(pgClient.GetDB()),
		EventRepo:     postgres.NewEventRepository(pgClient.GetDB()),
		RedisGameRepo: redis.NewGameRepository(redisClient.GetClient()),
		RedisCacheRepo: redis.NewCacheRepository(redisClient.GetClient()),
	}
}

func initWebSocketHub() *ws.Hub {
	hub := ws.NewHub()
	return hub
}

type Handlers struct {
	HTTPHandler *http.Handler
}

func initHandlers(hub *ws.Hub, packClient *grpcClient.PackClient, repos *Repositories, pgClient *postgres.Client, redisClient *redis.Client) *Handlers {
	return &Handlers{
		HTTPHandler: http.NewHandler(packClient, repos.GameRepo, repos.RedisGameRepo, hub, repos.EventRepo, pgClient, redisClient, packClient),
	}
}

func initWebSocketHandler(hub *ws.Hub, authClient *authClient.AuthServiceClient) *ws.Handler {
	return ws.NewHandler(hub, authClient)
}

func initRouter(handlers *Handlers, wsHandler *ws.Handler) *gin.Engine {
	return http.SetupRouter(handlers.HTTPHandler.Game, handlers.HTTPHandler.Health, wsHandler)
}


