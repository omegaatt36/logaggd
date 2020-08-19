package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  512,
		WriteBufferSize: 2048 * 2,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
	}
)

// Controller offers ws related endpoint.
type Controller struct {
	// store *tradingsystem.CacheStore
	hub *Hub

	stMsgsLock sync.Mutex
}

// RegisterEndpoint adds endpoint to upgrade to ws.
func (x *Controller) RegisterEndpoint(r *gin.Engine) {
	r.GET("/ws", x.WebSocket)
}

// Start starts ws hub.
func (x *Controller) Start(messageChan <-chan []byte) {
	x.hub = NewHub()

	for {
		select {
		case message, ok := <-messageChan:
			if !ok {
				log.Println("tick chan closed")
			}
			x.hub.broadcast <- func(client *Client) ([]byte, bool) {

				return []byte(message), true
			}
		}
	}
}

// WebSocket is gin handler.
func (x *Controller) WebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	x.hub.register <- conn
}
