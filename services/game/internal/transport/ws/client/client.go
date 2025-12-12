package client

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sigame/game/internal/infrastructure/logger"
	"sigame/game/internal/transport/ws/message"
)

type Hub interface {
	Unregister(client interface{ GetUserID() uuid.UUID; GetGameID() uuid.UUID; GetRTT() time.Duration; Send([]byte) })
	HandleMessage(client interface{ GetUserID() uuid.UUID; GetGameID() uuid.UUID; GetRTT() time.Duration; Send([]byte) }, msgData interface{})
}

type Client struct {
	hub    Hub
	conn   *websocket.Conn
	send   chan []byte
	userID uuid.UUID
	gameID uuid.UUID
	rtt    *RTTTracker
}

func NewClient(hub Hub, conn *websocket.Conn, userID, gameID uuid.UUID) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
		gameID: gameID,
		rtt:    newRTTTracker(),
	}
}

func (c *Client) UpdateRTT(rtt time.Duration) {
	c.rtt.UpdateRTT(rtt, c.userID)
}

func (c *Client) GetRTT() time.Duration {
	return c.rtt.GetRTT()
}

func (c *Client) SetLastPingSentAt(t time.Time) {
	c.rtt.SetLastPingSentAt(t)
}

func (c *Client) GetLastPingSentAt() time.Time {
	return c.rtt.GetLastPingSentAt()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(MaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(PongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	for {
		_, msgData, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Errorf(nil, "WebSocket error: %v", err)
			}
			break
		}

		clientMsg, err := message.NewClientMessage(msgData)
		if err != nil {
			logger.Warnf(nil, "Failed to parse client message: %v, data: %s", err, string(msgData))
			continue
		}

		logger.Infof(nil, "[Client] Parsed message: type=%s, user_id=%s, game_id=%s, payload=%v", clientMsg.GetType(), clientMsg.UserID, clientMsg.GameID, clientMsg.GetPayload())
		c.hub.HandleMessage(c, *clientMsg)
	}
}

func (c *Client) writePump() {
	jsonPingTicker := time.NewTicker(JSONPingPeriod)
	defer func() {
		jsonPingTicker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-jsonPingTicker.C:
			now := time.Now()
			c.SetLastPingSentAt(now)

			c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			pingMsg := message.NewPingMessage(now.UnixMilli())
			pingJSON, err := pingMsg.ToJSON()
			if err != nil {
				logger.Errorf(nil, "[PING] Failed to marshal ping: %v", err)
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, pingJSON); err != nil {
				logger.Errorf(nil, "[PING] Failed to send ping: %v", err)
				return
			}
		}
	}
}

func (c *Client) Send(data []byte) {
	select {
	case c.send <- data:
	default:
		close(c.send)
	}
}

func (c *Client) GetUserID() uuid.UUID {
	return c.userID
}

func (c *Client) GetGameID() uuid.UUID {
	return c.gameID
}

func (c *Client) Run() {
	go c.writePump()
	go c.readPump()
}

