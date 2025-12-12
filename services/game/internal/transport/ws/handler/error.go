package handler

import (
	"sigame/game/internal/infrastructure/logger"
	"sigame/game/internal/transport/ws/client"
	"sigame/game/internal/transport/ws/message"
)

func sendErrorMessage(cl *client.Client, msg, code string) {
	errorMsg := message.NewErrorMessage(msg, code)
	if data, err := errorMsg.ToJSON(); err == nil {
		cl.Send(data)
	} else {
		logger.Errorf(nil, "Failed to marshal error message: %v", err)
	}
}

