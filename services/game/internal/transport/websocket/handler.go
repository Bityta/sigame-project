package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// TODO: Restrict origins in production
		return true
	},
}

// Handler handles WebSocket connections
type Handler struct {
	hub *Hub
}

// NewHandler creates a new WebSocket handler
func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

// HandleWebSocket handles WebSocket upgrade requests
func (h *Handler) HandleWebSocket(c *gin.Context) {
	log.Printf("[WS] HandleWebSocket called, path=%s", c.Request.URL.Path)
	
	// Get game ID from path
	gameIDStr := c.Param("id")
	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		log.Printf("[WS] Invalid game ID: %s", gameIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	// Get user ID from query parameter (in production, this should come from JWT token)
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		log.Printf("[WS] Missing user_id for game %s", gameID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("[WS] Invalid user ID: %s", userIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	log.Printf("[WS] Checking game manager for game=%s, user=%s", gameID, userID)

	// TODO: Validate JWT token here
	// For now, we trust the user_id from query parameter

	// Check if game exists
	if _, exists := h.hub.GetGameManager(gameID); !exists {
		log.Printf("[WS] Game manager not found for game %s", gameID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found or not started"})
		return
	}

	log.Printf("[WS] Game manager found, upgrading connection for game=%s, user=%s", gameID, userID)

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WS] Failed to upgrade connection: %v", err)
		return
	}

	log.Printf("[WS] Connection upgraded, creating client for game=%s, user=%s", gameID, userID)

	// Create new client
	client := NewClient(h.hub, conn, userID, gameID)

	// Register client with hub
	h.hub.register <- client

	// Start client goroutines
	client.Run()

	log.Printf("[WS] WebSocket connection established: user=%s, game=%s", userID, gameID)
}

