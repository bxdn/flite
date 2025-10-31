package server

import (
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type wsep struct {
	path    string
	handler func(http.ResponseWriter, *http.Request)
}

func (w *wsep) Path() string { return w.path }

func (w *wsep) Handler() func(http.ResponseWriter, *http.Request) { return w.handler }

func WEBSOCK(path string, handler func(*websocket.Conn)) {
	wrapper := func(w http.ResponseWriter, r *http.Request) {
		c, e := websocket.Accept(w, r, nil)
		defer c.CloseNow()
		if e != nil {
			log.Printf("ERROR accepting websocket connection: %v\n", e)
			return
		}
		handler(c)
	}
	defaultServer.endpoints = append(defaultServer.endpoints, &wsep{
		path:    path,
		handler: wrapper,
	})
}
