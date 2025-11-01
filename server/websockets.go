package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type wsep struct {
	path    string
	handler func(http.ResponseWriter, *http.Request)
}

func (w *wsep) Path() string { return w.path }

func (w *wsep) Handler() func(w http.ResponseWriter, r *http.Request) { return w.handler }

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

func WS[T any](path string, handler func(c *websocket.Conn) error) {
	defaultServer.endpoints = append(defaultServer.endpoints, &wsep{
		path:    fmt.Sprintf("GET %s", path),
		handler: createHandler(handler),
	})
}

func ReadText(c *websocket.Conn) (string, error) {
	_, msgBytes, e := c.Read(context.Background())
	if e != nil {
		return "", fmt.Errorf("Error Reading text from websocket: %w", e)
	}
	return string(msgBytes), nil
}

func ReadObj[T any](c *websocket.Conn) (T, error) {
	ptr := new(T)
	if e := wsjson.Read(context.Background(), c, ptr); e != nil {
		return *ptr, fmt.Errorf("Error Reading obj from websocket: %w", e)
	}
	return *ptr, nil
}

func SendText(c *websocket.Conn, s string) error {
	e := c.Write(context.Background(), websocket.MessageText, []byte(s))
	if e != nil {
		return fmt.Errorf("Error Sending text to websocket: %w", e)
	}
	return nil
}

func SendObj(c *websocket.Conn, obj any) error {
	if e := wsjson.Write(context.Background(), c, obj); e != nil {
		return fmt.Errorf("Error Sending object to websocket: %w", e)
	}
	return nil
}
