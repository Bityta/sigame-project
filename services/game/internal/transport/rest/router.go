package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sigame/game/internal/transport/websocket"
)

// SetupRouter sets up the Gin router with all routes
func SetupRouter(handler *Handler, wsHandler *websocket.Handler) *gin.Engine {
	r := gin.New()

	// Middleware
	r.Use(gin.Recovery())
	r.Use(RequestResponseLoggingMiddleware()) // Async request/response logging
	r.Use(CORSMiddleware())

	// Health check
	r.GET("/health", handler.Health)
	r.HEAD("/health", handler.Health)

	// API routes
	api := r.Group("/api/game")
	{
		// Create game (internal call from lobby service, no auth needed)
		api.POST("/create", handler.CreateGame)

		// Get user's active game (requires auth)
		api.GET("/my-active", AuthMiddleware(), handler.GetMyActiveGame)

		// Get game info
		api.GET("/:id", handler.GetGame)

		// WebSocket connection
		api.GET("/:id/ws", wsHandler.HandleWebSocket)
	}

	return r
}

