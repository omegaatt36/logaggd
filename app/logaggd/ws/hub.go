package ws

import (
	"time"

	"github.com/gorilla/websocket"
)

// Hub defines hub to control clients.
type Hub struct {
	conns map[int64]*Client

	// Register requests from the clients.
	register chan *websocket.Conn

	// Unregister requests from clients.
	unregister chan int64

	broadcast chan func(*Client) ([]byte, bool)
}

// NewHub creates hub.
func NewHub() *Hub {
	h := &Hub{
		conns:      map[int64]*Client{},
		register:   make(chan *websocket.Conn),
		unregister: make(chan int64),
		broadcast:  make(chan func(*Client) ([]byte, bool)),
	}
	go h.run()
	return h
}

func (h *Hub) run() {
	for {
		select {
		case conn := <-h.register:
			h.AddConn(conn)
		case ID := <-h.unregister:
			h.RemoveConn(ID)
		case fn := <-h.broadcast:
			for ID, client := range h.conns {
				msg, ok := fn(client)
				if !ok {
					continue
				}

				select {
				case client.send <- msg:
				default:
					h.RemoveConn(ID)
				}
			}
		}
	}
}

// AddConn adds connection to hub.
func (h *Hub) AddConn(conn *websocket.Conn) {
	ID := time.Now().UnixNano()
	for _, ok := h.conns[ID]; ok; ID++ {
	}

	h.conns[ID] = NewClient(ID, conn, h)
}

// RemoveConn removes connection and closes it.
func (h *Hub) RemoveConn(ID int64) {
	conn, ok := h.conns[ID]
	if !ok {
		return
	}

	close(conn.send)
	delete(h.conns, ID)
}
