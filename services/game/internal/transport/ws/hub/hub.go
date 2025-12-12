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
	Stop()
}

type Client interface {
	GetUserID() uuid.UUID
	GetGameID() uuid.UUID
	GetRTT() time.Duration
	Send(data []byte)
}

type ClientMessageWrapper struct {
	Client  Client
	Message interface{}
}

type BroadcastMessage struct {
	GameID  uuid.UUID
	Message []byte
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
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	for {
		select {
		case client := <-h.register:
			func() {
				defer func() {
					if r := recover(); r != nil {
					}
				}()
				h.registerClient(client)
			}()

		case client := <-h.unregister:
			func() {
				defer func() {
					if r := recover(); r != nil {
					}
				}()
				h.unregisterClient(client)
			}()

		case wrapper := <-h.clientMessage:
			func() {
				defer func() {
					if r := recover(); r != nil {
					}
				}()
				h.handleClientMessage(wrapper)
			}()

		case msg := <-h.broadcast:
			func() {
				defer func() {
					if r := recover(); r != nil {
					}
				}()
				h.broadcastToGame(msg.GameID, msg.Message)
			}()
		}
	}
}

func (h *Hub) Register(cl Client) {
	h.register <- cl
}

func (h *Hub) Unregister(cl interface{ GetUserID() uuid.UUID; GetGameID() uuid.UUID; GetRTT() time.Duration; Send([]byte) }) {
	h.unregister <- cl.(Client)
}

func (h *Hub) HandleMessage(cl interface{ GetUserID() uuid.UUID; GetGameID() uuid.UUID; GetRTT() time.Duration; Send([]byte) }, msgData interface{}) {
	h.clientMessage <- &ClientMessageWrapper{
		Client:  cl.(Client),
		Message: msgData,
	}
}

func (h *Hub) Broadcast(gameID uuid.UUID, message []byte) {
	h.broadcast <- &BroadcastMessage{
		GameID:  gameID,
		Message: message,
	}
}

func (h *Hub) handleClientMessage(wrapper *ClientMessageWrapper) {
	gameID := wrapper.Client.GetGameID()
	userID := wrapper.Client.GetUserID()
	
	h.mu.RLock()
	manager, exists := h.managers[gameID]
	h.mu.RUnlock()

	if !exists {
		return
	}

	manager.HandleClientMessage(userID, wrapper.Message)
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

func (h *Hub) RegisterGameManager(gameID uuid.UUID, manager GameManager) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.managers[gameID] = manager
}

func (h *Hub) UnregisterGameManager(gameID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.managers, gameID)
}

func (h *Hub) GetGameManager(gameID uuid.UUID) (GameManager, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	manager, exists := h.managers[gameID]
	return manager, exists
}

func (h *Hub) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, manager := range h.managers {
		manager.Stop()
	}

	h.managers = make(map[uuid.UUID]GameManager)
}
