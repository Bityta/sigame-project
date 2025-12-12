package websocket

import (
	"github.com/sigame/game/internal/logger"
)

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	gameID := client.GetGameID()
	userID := client.GetUserID()

	if h.games[gameID] == nil {
		h.games[gameID] = make(map[*Client]bool)
	}

	h.games[gameID][client] = true

	manager, exists := h.managers[gameID]
	h.mu.Unlock()

	logger.Infof(nil, "Client registered for game %s (user %s)", gameID, userID)

	if exists && manager != nil {
		manager.SetPlayerConnected(userID, true)

		logger.Debugf(nil, "Sending initial state to client %s for game %s", userID, gameID)
		manager.SendStateToClient(client)
	}
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	gameID := client.GetGameID()
	userID := client.GetUserID()

	var manager GameManager
	var managerExists bool

	if clients, ok := h.games[gameID]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			close(client.send)

			logger.Infof(nil, "Client unregistered from game %s (user %s)", gameID, userID)

			manager, managerExists = h.managers[gameID]

			if len(clients) == 0 {
				delete(h.games, gameID)
			}
		}
	}
	h.mu.Unlock()

	if managerExists && manager != nil {
		manager.SetPlayerConnected(userID, false)
	}
}

