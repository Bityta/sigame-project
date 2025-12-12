package http

import (
	"github.com/sigame/game/internal/transport/http/handler"
)

type Handler struct {
	Game   *handler.GameHandler
	Health *handler.HealthHandler
}

func NewHandler() *Handler {
	return &Handler{
		Game:   handler.NewGameHandler(),
		Health: handler.NewHealthHandler(),
	}
}

