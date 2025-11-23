package websocket

import (
	"log"
	"sync"

	"github.com/google/uuid"
)

// GameManager interface to avoid circular dependency
type GameManager interface {
	HandleClientMessage(userID uuid.UUID, message *ClientMessage)
}

// ClientMessageWrapper wraps a client message with the client reference
type ClientMessageWrapper struct {
	client  *Client
	message *ClientMessage
}

// Hub maintains active clients and game managers
type Hub struct {
	// Registered clients per game
	games map[uuid.UUID]map[*Client]bool

	// Game managers
	managers map[uuid.UUID]GameManager

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Client messages
	clientMessage chan *ClientMessageWrapper

	// Broadcast messages to specific game
	broadcast chan *BroadcastMessage

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// BroadcastMessage represents a message to broadcast to all clients in a game
type BroadcastMessage struct {
	GameID  uuid.UUID
	Message []byte
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		games:         make(map[uuid.UUID]map[*Client]bool),
		managers:      make(map[uuid.UUID]GameManager),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		clientMessage: make(chan *ClientMessageWrapper),
		broadcast:     make(chan *BroadcastMessage, 256),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case wrapper := <-h.clientMessage:
			h.handleClientMessage(wrapper)

		case broadcast := <-h.broadcast:
			h.broadcastToGame(broadcast.GameID, broadcast.Message)
		}
	}
}

// registerClient registers a new client
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	gameID := client.GetGameID()

	// Create game clients map if doesn't exist
	if h.games[gameID] == nil {
		h.games[gameID] = make(map[*Client]bool)
	}

	h.games[gameID][client] = true
	log.Printf("Client registered for game %s (user %s)", gameID, client.GetUserID())
}

// unregisterClient unregisters a client
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	gameID := client.GetGameID()

	if clients, ok := h.games[gameID]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			close(client.send)

			log.Printf("Client unregistered from game %s (user %s)", gameID, client.GetUserID())

			// If no more clients, consider cleaning up game manager
			if len(clients) == 0 {
				delete(h.games, gameID)
				// Note: Game manager cleanup is handled separately when game ends
			}
		}
	}
}

// handleClientMessage handles a message from a client
func (h *Hub) handleClientMessage(wrapper *ClientMessageWrapper) {
	h.mu.RLock()
	manager, exists := h.managers[wrapper.client.GetGameID()]
	h.mu.RUnlock()

	if !exists {
		log.Printf("No game manager found for game %s", wrapper.client.GetGameID())
		errorMsg := NewErrorMessage("Game not found or not started", "GAME_NOT_FOUND")
		if data, err := errorMsg.ToJSON(); err == nil {
			wrapper.client.Send(data)
		}
		return
	}

	// Forward message to game manager
	manager.HandleClientMessage(wrapper.client.GetUserID(), wrapper.message)
}

// broadcastToGame broadcasts a message to all clients in a game
func (h *Hub) broadcastToGame(gameID uuid.UUID, message []byte) {
	h.mu.RLock()
	clients := h.games[gameID]
	h.mu.RUnlock()

	for client := range clients {
		select {
		case client.send <- message:
		default:
			// Client channel is full, close it
			close(client.send)
			h.mu.Lock()
			delete(h.games[gameID], client)
			h.mu.Unlock()
		}
	}
}

// RegisterGameManager registers a game manager with the hub
func (h *Hub) RegisterGameManager(gameID uuid.UUID, manager GameManager) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.managers[gameID] = manager
	log.Printf("Game manager registered for game %s", gameID)
}

// UnregisterGameManager unregisters a game manager
func (h *Hub) UnregisterGameManager(gameID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.managers, gameID)
	log.Printf("Game manager unregistered for game %s", gameID)
}

// GetGameManager retrieves a game manager
func (h *Hub) GetGameManager(gameID uuid.UUID) (GameManager, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	manager, exists := h.managers[gameID]
	return manager, exists
}

// Broadcast sends a message to all clients in a game
func (h *Hub) Broadcast(gameID uuid.UUID, message []byte) {
	h.broadcast <- &BroadcastMessage{
		GameID:  gameID,
		Message: message,
	}
}

// GetActiveGamesCount returns the number of active games
func (h *Hub) GetActiveGamesCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.managers)
}

// GetGameClientsCount returns the number of clients in a game
func (h *Hub) GetGameClientsCount(gameID uuid.UUID) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.games[gameID]; ok {
		return len(clients)
	}
	return 0
}

