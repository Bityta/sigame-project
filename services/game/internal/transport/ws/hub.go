package websocket

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sigame/game/internal/logger"
)

type GameManager interface {
	HandleClientMessage(userID uuid.UUID, message *ClientMessage)
	SendStateToClient(client *Client)
	SetPlayerConnected(userID uuid.UUID, connected bool)
}

type ClientMessageWrapper struct {
	client  *Client
	message *ClientMessage
}

type Hub struct {
	games         map[uuid.UUID]map[*Client]bool
	managers      map[uuid.UUID]GameManager
	register      chan *Client
	unregister    chan *Client
	clientMessage chan *ClientMessageWrapper
	broadcast     chan *BroadcastMessage
	mu            sync.RWMutex
}

type BroadcastMessage struct {
	GameID  uuid.UUID
	Message []byte
}

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

func (h *Hub) GetClientRTT(gameID, userID uuid.UUID) time.Duration {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.games[gameID]
	if !ok {
		return 0
	}

	for client := range clients {
		if client.GetUserID() == userID {
			return client.GetRTT()
		}
	}

	return 0
}

func (h *Hub) broadcastToGame(gameID uuid.UUID, message []byte) {
	h.mu.RLock()
	clients := h.games[gameID]
	h.mu.RUnlock()

	for client := range clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			h.mu.Lock()
			delete(h.games[gameID], client)
			h.mu.Unlock()
		}
	}
}

func (h *Hub) RegisterGameManager(gameID uuid.UUID, manager GameManager) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.managers[gameID] = manager
	logger.Infof(nil, "Game manager registered for game %s", gameID)
}

func (h *Hub) UnregisterGameManager(gameID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.managers, gameID)
	logger.Infof(nil, "Game manager unregistered for game %s", gameID)
}

func (h *Hub) GetGameManager(gameID uuid.UUID) (GameManager, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	manager, exists := h.managers[gameID]
	return manager, exists
}

func (h *Hub) Broadcast(gameID uuid.UUID, message []byte) {
	h.broadcast <- &BroadcastMessage{
		GameID:  gameID,
		Message: message,
	}
}

func (h *Hub) GetActiveGamesCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.managers)
}

func (h *Hub) GetGameClientsCount(gameID uuid.UUID) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.games[gameID]; ok {
		return len(clients)
	}
	return 0
}

