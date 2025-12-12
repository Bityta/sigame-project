package hub

import "github.com/google/uuid"

func (h *Hub) broadcastToGame(gameID uuid.UUID, message []byte) {
	h.mu.RLock()
	clients := h.games[gameID]
	h.mu.RUnlock()

	for client := range clients {
		client.Send(message)
	}
}

func (h *Hub) BroadcastToUser(gameID, userID uuid.UUID, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.games[gameID]
	if !ok {
		return
	}

	for client := range clients {
		if client.GetUserID() == userID {
			client.Send(message)
			return
		}
	}
}

func (h *Hub) BroadcastExcept(gameID, exceptUserID uuid.UUID, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.games[gameID]
	if !ok {
		return
	}

	for client := range clients {
		if client.GetUserID() != exceptUserID {
			client.Send(message)
		}
	}
}

