package ws

import "github.com/sigame/game/internal/infrastructure/logger"

func sendErrorMessage(client *Client, message, code string) {
	errorMsg := NewErrorMessage(message, code)
	if data, err := errorMsg.ToJSON(); err == nil {
		client.Send(data)
	} else {
		logger.Errorf(nil, "Failed to marshal error message: %v", err)
	}
}

