package logaggd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/omegaatt36/logaggd/app/logaggd/ws"
	"github.com/omegaatt36/logaggd/cache"

	"github.com/gin-gonic/gin"
)

// Server is a HTTP server.
type Server struct {
	gcid uint64

	messageChan chan []byte
}

// RegisterEndpoint registers HTTP handler.
func (s *Server) RegisterEndpoint(r *gin.Engine) {
	r.POST("/*action", s.action)
	wsX := ws.Controller{}
	wsX.RegisterEndpoint(r)
	go wsX.Start(s.messageChan)
}

// Start starts HTTP server.
func (s *Server) Start(ctx context.Context, addr string) {
	s.messageChan = make(chan []byte, 100)
	engine := gin.Default()
	engine.RedirectTrailingSlash = true
	cache.Initialize("localhost:6379")
	s.RegisterEndpoint(engine)

	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}
	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown: ", err)
		}
	}()

	go func() {
		t := time.NewTicker(time.Minute)
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if err := deleteExpiredMessages(ctx); err != nil {
					log.Println("deleteExpiredMessages()", err)
				}
			}
		}
	}()

	log.Println("starts serving...")
	if err := srv.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}
}

func (s *Server) action(c *gin.Context) {

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		c.String(500, "bad")
		return
	}

	bs := bytes.Split(body, []byte("\n"))
	for _, b := range bs {
		if len(b) <= 2 {
			continue
		}

		var m map[string]interface{}
		err = json.Unmarshal(b, &m)
		if err != nil {
			log.Println(err)
			continue
		}

		err = record(c.Request.Context(), b)
		if err != nil {
			log.Printf("want to store message but redis get error, %v\n", err)
			continue
		}
		s.messageChan <- b
		log.Printf("%+v\n", m)
	}

	c.String(200, "ok")
}
