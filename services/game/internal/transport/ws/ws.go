package ws

import (
	"github.com/gin-gonic/gin"
	"sigame/game/internal/transport/ws/client"
	"sigame/game/internal/transport/ws/handler"
	"sigame/game/internal/transport/ws/hub"
	"sigame/game/internal/transport/ws/message"
)

type Hub = hub.Hub
type Client = client.Client
type ClientMessage = message.ClientMessage
type ServerMessage = message.ServerMessage
type MessageType = message.MessageType

func NewHub() *Hub {
	return hub.New()
}

func NewErrorMessage(msg, code string) *ServerMessage {
	return message.NewErrorMessage(msg, code)
}

func NewClientMessage(data []byte) (*ClientMessage, error) {
	return message.NewClientMessage(data)
}

type Handler struct {
	wsHandler *handler.Handler
}

func NewHandler(h *Hub, authClient handler.AuthService) *Handler {
	return &Handler{
		wsHandler: handler.NewHandler(h, authClient),
	}
}

func (h *Handler) HandleWebSocket(c *gin.Context) {
	h.wsHandler.HandleWebSocket(c)
}
