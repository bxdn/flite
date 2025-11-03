package server

import (
	"fmt"
	"strings"

	"github.com/bxdn/flite/shared"
)

func (f *F[T]) PrepareAsSSEHandler() {
	f.res.Header().Set("Content-Type", "text/event-stream")
	f.res.Header().Set("Cache-Control", "no-cache")
	f.res.Header().Set("Connection", "keep-alive")
}

func (f *F[T]) SendEvent(event shared.SSEEvent) error {
	if event.ID != "" {
		if _, e := fmt.Fprintf(f.res, "id: %s\n", event.ID); e != nil {
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
