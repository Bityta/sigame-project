package hub

import "github.com/google/uuid"

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

func (h *Hub) handleClientMessage(wrapper *ClientMessageWrapper) {
	h.mu.RLock()
	manager, exists := h.managers[wrapper.Client.GetGameID()]
	h.mu.RUnlock()

	if !exists {
		return
	}

	manager.HandleClientMessage(wrapper.Client.GetUserID(), wrapper.Message)
}

