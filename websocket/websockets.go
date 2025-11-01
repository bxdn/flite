package server

import (
	"context"
	"fmt"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

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
