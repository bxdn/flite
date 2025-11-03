package server

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

func (f *F[T]) PrepareAsSSEHandler() error {
	f.res.Header().Set("Content-Type", "text/event-stream")
	f.res.Header().Set("Cache-Control", "no-cache")
	f.res.Header().Set("Connection", "keep-alive")
	if _, e := fmt.Fprintf(f.res, ": connected\n\n"); e != nil {
		f.ReturnError("Error writing to response writer!", http.StatusInternalServerError)
		return fmt.Errorf("Error: Could not write to response writer?")
	}
	f.res.WriteHeader(200)
	flusher, ok := f.Res().(http.Flusher)
	if !ok {
		f.ReturnError("Streaming unsupported!", http.StatusInternalServerError)
		return fmt.Errorf("Error: Response writer is not a flusher?")
	}
	flusher.Flush()
	return nil
}

func (f *F[T]) SendEvent(event SSEEvent) error {
	flusher, ok := f.Res().(http.Flusher)
	if !ok {
		f.ReturnError("Streaming unsupported!", http.StatusInternalServerError)
		return fmt.Errorf("Error: Response writer is not a flusher?")
	}

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
	flusher.Flush()
	return nil
}

func ReceiveEvent(reader *bufio.Reader) (SSEEvent, error) {
	var buffer strings.Builder
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return SSEEvent{}, fmt.Errorf("Error reading event: %w", err)
		}
		if line == "\n" || line == "\r\n" {
			return parseSSEEvent(buffer.String()), nil
		} else {
			buffer.WriteString(line)
		}
	}
}

type SSEEvent struct {
	Event string
	Data  string
	ID    string
}

func parseSSEEvent(raw string) SSEEvent {
	var e SSEEvent
	var dataLines []string

	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "event: "):
			e.Event = strings.TrimPrefix(line, "event: ")
		case strings.HasPrefix(line, "data: "):
			dataLines = append(dataLines, strings.TrimPrefix(line, "data: "))
		case strings.HasPrefix(line, "id: "):
			e.ID = strings.TrimPrefix(line, "id: ")
		}
	}

	if len(dataLines) > 0 {
		e.Data = strings.Join(dataLines, "")
	}

	return e
}
