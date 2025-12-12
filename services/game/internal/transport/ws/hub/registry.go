package hub

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

func (h *Hub) unregisterClient(client Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	gameID := client.GetGameID()

	if clients, ok := h.games[gameID]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)

			if len(clients) == 0 {
				delete(h.games, gameID)
			}

			if manager, exists := h.managers[gameID]; exists {
				manager.SetPlayerConnected(client.GetUserID(), false)
			}
		}
	}
}

func (h *Hub) GetClient(gameID, userID interface{}) Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.games[gameID.(interface{})]
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

