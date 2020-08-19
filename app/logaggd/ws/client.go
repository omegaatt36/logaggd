package ws

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 256
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client defines websocket client.
type Client struct {
	ID   int64
	conn *websocket.Conn
	hub  *Hub
	send chan []byte
}

// NewClient creates websocket client.
func NewClient(ID int64, conn *websocket.Conn, hub *Hub) *Client {
	c := &Client{
		ID:   ID,
		conn: conn,
		hub:  hub,
		send: make(chan []byte, 64),
	}

	// todo: using websocket chrome extension can't disconnect maybe cause by I delete readPump.
	go c.writePump()
	return c
}

// Request is common request object that must has type field.
type Request struct {
	Type string `json:"type"`
}

// LoginRequest is request to login.
type LoginRequest struct {
	Token string `json:"token"`
}

// writePump write data from websocket. it doesn't support concurrent write.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
