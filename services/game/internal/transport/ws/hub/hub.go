package hub

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type GameManager interface {
	HandleClientMessage(userID uuid.UUID, message interface{})
	SendStateToClient(client interface{})
	SetPlayerConnected(userID uuid.UUID, connected bool)
}

type Client interface {
	GetUserID() uuid.UUID
	GetGameID() uuid.UUID
	GetRTT() time.Duration
	Send(data []byte)
}

type Hub struct {
	games         map[uuid.UUID]map[Client]bool
	managers      map[uuid.UUID]GameManager
	register      chan Client
	unregister    chan Client
	clientMessage chan *ClientMessageWrapper
	broadcast     chan *BroadcastMessage
	mu            sync.RWMutex
}

type ClientMessageWrapper struct {
	Client  Client
	Message interface{}
}

type BroadcastMessage struct {
	GameID  uuid.UUID
	Message []byte
}

func New() *Hub {
	return &Hub{
		games:         make(map[uuid.UUID]map[Client]bool),
		managers:      make(map[uuid.UUID]GameManager),
		register:      make(chan Client),
		unregister:    make(chan Client),
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

func (h *Hub) Register(client Client) {
	h.register <- client
}

func (h *Hub) Unregister(client Client) {
	h.unregister <- client
}

func (h *Hub) HandleMessage(client Client, message interface{}) {
	h.clientMessage <- &ClientMessageWrapper{
		Client:  client,
		Message: message,
	}
}

func (h *Hub) Broadcast(gameID uuid.UUID, message []byte) {
	h.broadcast <- &BroadcastMessage{
		GameID:  gameID,
		Message: message,
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

