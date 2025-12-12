package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sigame/game/internal/transport/http/handler"
	"github.com/sigame/game/internal/transport/http/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.Logging())
	r.Use(middleware.CORS())

	healthHandler := handler.NewHealthHandler()
	gameHandler := handler.NewGameHandler()

	r.GET("/health", healthHandler.Health)
	r.HEAD("/health", healthHandler.Health)

	api := r.Group("/api/game")
	{
		api.POST("", gameHandler.CreateGame)
		api.GET("/my-active", middleware.Auth(), gameHandler.GetMyActiveGame)
		api.GET("/:id", gameHandler.GetGame)
	}

	return r
}

