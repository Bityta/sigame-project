package message

import (
	"time"

	"github.com/sigame/game/internal/infrastructure/logger"
)

func (h *Hub) handleClientMessage(wrapper *ClientMessageWrapper) {
	if wrapper.message.Type == MessageTypePong {
		h.handlePong(wrapper)
		return
	}

	h.mu.RLock()
	manager, exists := h.managers[wrapper.client.GetGameID()]
	h.mu.RUnlock()

	if !exists {
		logger.Errorf(nil, "No game manager found for game %s", wrapper.client.GetGameID())
		sendErrorMessage(wrapper.client, ErrorGameManagerNotFound, "GAME_NOT_FOUND")
		return
	}

	manager.HandleClientMessage(wrapper.client.GetUserID(), wrapper.message)
}

func (h *Hub) handlePong(wrapper *ClientMessageWrapper) {
	now := time.Now()

	serverTimeRaw, ok := wrapper.message.Payload["server_time"]
	if !ok {
		logger.Errorf(nil, "[PONG] Missing server_time in PONG from user %s", wrapper.client.GetUserID())
		return
	}

	var serverTimeMs int64
	switch v := serverTimeRaw.(type) {
	case float64:
		serverTimeMs = int64(v)
	case int64:
		serverTimeMs = v
	default:
		logger.Errorf(nil, "[PONG] Invalid server_time type from user %s: %T", wrapper.client.GetUserID(), serverTimeRaw)
		return
	}

	serverTime := time.UnixMilli(serverTimeMs)
	rtt := now.Sub(serverTime)

	if rtt < 0 || rtt > MaxRTTDuration {
		logger.Errorf(nil, "[PONG] Invalid RTT %v from user %s, ignoring", rtt, wrapper.client.GetUserID())
		return
	}

	wrapper.client.UpdateRTT(rtt)
}

