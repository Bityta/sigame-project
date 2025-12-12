package http

import (
	"github.com/gin-gonic/gin"
	"sigame/game/internal/transport/http/handler"
	"sigame/game/internal/transport/http/middleware"
)

type WSHandler interface {
	HandleWebSocket(c *gin.Context)
}

func SetupRouter(gameHandler *handler.GameHandler, healthHandler *handler.HealthHandler, wsHandler WSHandler) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.Logging())
	r.Use(middleware.CORS())

	r.GET("/health", healthHandler.Health)
	r.HEAD("/health", healthHandler.Health)

	api := r.Group("/api/game")
	{
		api.POST("", gameHandler.CreateGame)
		api.GET("/my-active", middleware.Auth(), gameHandler.GetMyActiveGame)
		api.GET("/:id", gameHandler.GetGame)
		api.GET("/:id/ws", wsHandler.HandleWebSocket)
	}

	return r
}

