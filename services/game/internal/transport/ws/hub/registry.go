package hub

import "github.com/google/uuid"

func (h *Hub) registerClient(client Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	gameID := client.GetGameID()

	if _, exists := h.games[gameID]; !exists {
		h.games[gameID] = make(map[Client]bool)
	}

	h.games[gameID][client] = true

	if manager, exists := h.managers[gameID]; exists {
		go manager.SendStateToClient(client)
		manager.SetPlayerConnected(client.GetUserID(), true)
	}
}

func (h *Hub) unregisterClient(cl Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	gameID := cl.GetGameID()

	if clients, ok := h.games[gameID]; ok {
		if _, exists := clients[cl]; exists {
			delete(clients, cl)

			if len(clients) == 0 {
				delete(h.games, gameID)
			}

			if manager, exists := h.managers[gameID]; exists {
				manager.SetPlayerConnected(cl.GetUserID(), false)
			}
		}
	}
}

func (h *Hub) GetClient(gameID, userID interface{}) Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	gid, ok := gameID.(uuid.UUID)
	if !ok {
		return nil
	}

	clients, ok := h.games[gid]
	if !ok {
		return nil
	}

	for client := range clients {
		if client.GetUserID() == userID {
			return client
		}
	}

	return nil
}

