package server

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bxdn/flite/shared"
)

type Event = shared.SSEEvent

type JsonEvent struct {
	Id, Event string
	Data      any
}

func (f *F[T]) PrepareAsSSEHandler() {
	f.res.Header().Set("Content-Type", "text/event-stream")
	f.res.Header().Set("Cache-Control", "no-cache")
	f.res.Header().Set("Connection", "keep-alive")
}

func (f *F[T]) SendEvent(event Event) error {
	if event.Id != "" {
		if _, e := fmt.Fprintf(f.res, "id: %s\n", event.Id); e != nil {
			return fmt.Errorf("Error writing text event id: %w", e)
		}
	}

	if event.Event != "" {
		if _, e := fmt.Fprintf(f.res, "event: %s\n", event.Event); e != nil {
			return fmt.Errorf("Error writing text event type: %w", e)
		}
	}

	if event.Data != "" {
		lines := strings.Split(event.Data, "\n")
		for _, line := range lines {
			if _, e := fmt.Fprintf(f.res, "data: %s\n", line); e != nil {
				return fmt.Errorf("Error writing text event data: %w", e)
			}
		}
	}
	if _, e := fmt.Fprintf(f.res, "\n"); e != nil {
		return fmt.Errorf("Error writing text event: %w", e)
	}
	f.res.Flush()
	return nil
}

func (f *F[T]) SendJSONEvent(event JsonEvent) error {
	jsonBytes, e := json.Marshal(event.Data)
	if e != nil {
		return fmt.Errorf("Error marshalling json event: %w", e)
	}
	return f.SendEvent(Event{Id: event.Id, Event: event.Event, Data: string(jsonBytes)})
}
