package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sigame/auth/internal/metrics"
)

// SetupRouter configures and returns the Gin router with all routes and middleware
func SetupRouter(handler *Handler, jwtMiddleware gin.HandlerFunc, metrics *metrics.Metrics) *gin.Engine {
	router := gin.Default()

	// Apply metrics middleware globally
	router.Use(MetricsMiddleware(metrics))

	// CORS middleware
	router.Use(CORSMiddleware())

	// Health check endpoint (no auth required)
	router.GET("/health", handler.Health)

	// Auth routes
	auth := router.Group("/auth")
	{
		// Public endpoints
		auth.GET("/check-username", handler.CheckUsername)
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.Refresh)

		// Protected endpoints
		auth.POST("/logout", jwtMiddleware, handler.Logout)
		auth.GET("/me", jwtMiddleware, handler.GetMe)
	}

	return router
}

