package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type wsep struct {
	path    string
	handler func(http.ResponseWriter, *http.Request)
}

func (w *wsep) Path() string { return w.path }

func (w *wsep) Handler() func(w http.ResponseWriter, r *http.Request) { return w.handler }

func WS(path string, handler func(c *websocket.Conn) error) {
	defaultServer.endpoints = append(defaultServer.endpoints, &wsep{
		path:    fmt.Sprintf("GET %s", path),
		handler: createHandler(handler),
	})
}

func createHandler(handler func(c *websocket.Conn) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c, e := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
		defer c.CloseNow()
		if e != nil {
			log.Printf("ERROR accepting websocket connection: %v\n", e)
			return
		}
		if e := handler(c); e != nil {
			log.Printf("ERROR executing websocket logic: %v\n", e)
		}
	}
}
