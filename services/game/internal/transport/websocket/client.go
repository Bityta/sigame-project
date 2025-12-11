package websocket

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send JSON PING to peer with this period for RTT measurement (5 seconds as per spec)
	jsonPingPeriod = 5 * time.Second

	// Maximum message size allowed from peer
	maxMessageSize = 8192

	// Maximum number of RTT samples to keep for averaging
	maxRTTSamples = 10
)

// Client represents a WebSocket client connection
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID uuid.UUID
	gameID uuid.UUID

	// RTT tracking for ping compensation
	rttSamples     []time.Duration
	avgRTT         time.Duration
	lastPingSentAt time.Time
	rttMu          sync.RWMutex
}

// NewClient creates a new WebSocket client
func NewClient(hub *Hub, conn *websocket.Conn, userID, gameID uuid.UUID) *Client {
	return &Client{
		hub:        hub,
		conn:       conn,
		send:       make(chan []byte, 256),
		userID:     userID,
		gameID:     gameID,
		rttSamples: make([]time.Duration, 0, maxRTTSamples),
	}
}

// UpdateRTT updates the RTT with a new sample and recalculates average
func (c *Client) UpdateRTT(rtt time.Duration) {
	c.rttMu.Lock()
	defer c.rttMu.Unlock()

	// Add new sample
	c.rttSamples = append(c.rttSamples, rtt)

	// Keep only last maxRTTSamples
	if len(c.rttSamples) > maxRTTSamples {
		c.rttSamples = c.rttSamples[1:]
	}

	// Calculate average RTT
	var total time.Duration
	for _, sample := range c.rttSamples {
		total += sample
	}
	c.avgRTT = total / time.Duration(len(c.rttSamples))

	log.Printf("[RTT] User %s: new sample=%v, avg=%v (samples=%d)",
		c.userID, rtt, c.avgRTT, len(c.rttSamples))
}

// GetRTT returns the average RTT for this client
func (c *Client) GetRTT() time.Duration {
	c.rttMu.RLock()
	defer c.rttMu.RUnlock()
	return c.avgRTT
}

// SetLastPingSentAt records when the last ping was sent
func (c *Client) SetLastPingSentAt(t time.Time) {
	c.rttMu.Lock()
	defer c.rttMu.Unlock()
	c.lastPingSentAt = t
}

// GetLastPingSentAt returns when the last ping was sent
func (c *Client) GetLastPingSentAt() time.Time {
	c.rttMu.RLock()
	defer c.rttMu.RUnlock()
	return c.lastPingSentAt
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse message
		clientMsg, err := NewClientMessage(message)
		if err != nil {
			log.Printf("Failed to parse client message: %v", err)
			// Send error back to client
			errorMsg := NewErrorMessage("Invalid message format", "INVALID_MESSAGE")
			if data, err := errorMsg.ToJSON(); err == nil {
				c.send <- data
			}
			continue
		}

		// Forward message to hub for processing
		c.hub.clientMessage <- &ClientMessageWrapper{
			client:  c,
			message: clientMsg,
		}
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	// JSON PING ticker for RTT measurement (5 seconds)
	jsonPingTicker := time.NewTicker(jsonPingPeriod)
	defer func() {
		jsonPingTicker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-jsonPingTicker.C:
			// Send JSON PING for RTT measurement
			now := time.Now()
			c.SetLastPingSentAt(now)

			pingMsg := NewPingMessage(now.UnixMilli())
			data, err := pingMsg.ToJSON()
			if err != nil {
				log.Printf("[PING] Failed to marshal ping message: %v", err)
				continue
			}

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("[PING] Failed to send ping: %v", err)
				return
			}
		}
	}
}

// Send sends a message to the client
func (c *Client) Send(data []byte) {
	select {
	case c.send <- data:
	default:
		// Client's send channel is full, close it
		close(c.send)
	}
}

// GetUserID returns the client's user ID
func (c *Client) GetUserID() uuid.UUID {
	return c.userID
}

// GetGameID returns the client's game ID
func (c *Client) GetGameID() uuid.UUID {
	return c.gameID
}

// Run starts the client's read and write pumps
func (c *Client) Run() {
	go c.writePump()
	go c.readPump()
}

