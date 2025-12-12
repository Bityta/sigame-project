package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/sigame/game/internal/transport/ws/client"
	"github.com/sigame/game/internal/transport/ws/handler"
	"github.com/sigame/game/internal/transport/ws/hub"
	"github.com/sigame/game/internal/transport/ws/message"
)

type Hub = hub.Hub
type Client = client.Client
type ClientMessage = message.ClientMessage
type ServerMessage = message.ServerMessage
type MessageType = message.MessageType

func NewHub() *Hub {
	return hub.New()
}

type Handler struct {
	wsHandler *handler.Handler
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		wsHandler: handler.NewHandler(h),
	}
}

func (h *Handler) HandleWebSocket(c *gin.Context) {
	h.wsHandler.HandleWebSocket(c)
}



