package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"sigame/game/internal/adapter/grpc/pack"
	"sigame/game/internal/adapter/repository/postgres"
	"sigame/game/internal/adapter/repository/redis"
)

type HealthHandler struct {
	pgClient   *postgres.Client
	redisClient *redis.Client
	packClient *pack.PackClient
}

func NewHealthHandler(pgClient *postgres.Client, redisClient *redis.Client, packClient *pack.PackClient) *HealthHandler {
	return &HealthHandler{
		pgClient:    pgClient,
		redisClient: redisClient,
		packClient:  packClient,
	}
}

func (h *HealthHandler) Health(c *gin.Context) {
	status := "healthy"
	checks := make(map[string]string)

	if err := h.pgClient.Ping(); err != nil {
		status = "unhealthy"
		checks["postgres"] = "failed"
	} else {
		checks["postgres"] = "ok"
	}

	if err := h.redisClient.Ping(c.Request.Context()); err != nil {
		status = "unhealthy"
		checks["redis"] = "failed"
	} else {
		checks["redis"] = "ok"
	}

	httpStatus := http.StatusOK
	if status == "unhealthy" {
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, gin.H{
		"status":    status,
		"service":   "game",
		"timestamp": time.Now().UTC(),
		"checks":    checks,
	})
}

