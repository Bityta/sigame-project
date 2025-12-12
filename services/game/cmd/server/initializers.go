package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sigame/game/internal/infrastructure/config"
	grpcClient "github.com/sigame/game/internal/adapter/grpc/pack"
	"github.com/sigame/game/internal/infrastructure/metrics"
	"github.com/sigame/game/internal/adapter/repository/postgres"
	"github.com/sigame/game/internal/adapter/repository/redis"
	"github.com/sigame/game/internal/infrastructure/tracing"
	"github.com/sigame/game/internal/transport/http"
	"github.com/sigame/game/internal/transport/ws"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func initTracer(serviceName string) (*tracing.TracerProvider, error) {
	return tracing.InitTracer(serviceName)
}

func initMetrics() *metrics.Metrics {
	return metrics.NewMetrics()
}

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

func initHandlers(hub *ws.Hub, packClient *grpcClient.PackClient, repos *Repositories) *Handlers {
	return &Handlers{
		HTTPHandler: http.NewHandler(),
	}
}

func initWebSocketHandler(hub *ws.Hub) *ws.Handler {
	return ws.NewHandler(hub)
}

func initRouter(handlers *Handlers, wsHandler *ws.Handler) *gin.Engine {
	router := http.SetupRouter(handlers.HTTPHandler.Game, handlers.HTTPHandler.Health, wsHandler)
	router.Use(otelgin.Middleware(ServiceName))
	return router
}


