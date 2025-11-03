package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/bxdn/flite/shared"
)

type RequestConfig struct {
	Url     string
	Headers map[string]string
	Query   map[string]string
}

type ConfigWithBody struct {
	RequestConfig
	Body []byte
}

type fullConfig struct {
	ConfigWithBody
	Method string
}

type Event = shared.SSEEvent

func FromJson[T any](res *http.Response, bodyBytes []byte, e error) (*http.Response, T, error) {
	ptr := new(T)
	if e != nil {
		return res, *ptr, e
	}
	decoder := json.NewDecoder(bytes.NewBuffer(bodyBytes))
	e = decoder.Decode(ptr)
	return res, *ptr, e
}

func ToJson(config RequestConfig, object any, req func(ConfigWithBody) (*http.Response, []byte, error)) (*http.Response, []byte, error) {
	jsonBytes, e := json.Marshal(object)
	if e != nil {
		return nil, nil, fmt.Errorf("Error Marshalling body to JSON: %w", e)
	}
	return req(ConfigWithBody{RequestConfig: config, Body: jsonBytes})
}

func Subscribe(config ConfigWithBody, method string, onEvent func(Event) error) error {
	u, e := url.Parse(config.Url)
	if e != nil {
		return fmt.Errorf("Error parsing url: %w", e)
	}

	for k, v := range config.Query {
		u.Query().Set(k, v)
	}

	req, e := http.NewRequest(method, u.String(), bytes.NewBuffer(config.Body))
	if e != nil {
		return fmt.Errorf("Error creating request: %w", e)
	}

	for k, v := range config.Headers {
		req.Header.Set(k, v)
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, e := client.Do(req)
	if e != nil {
		return fmt.Errorf("Error executing request: %w", e)
	}

	reader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	for {
		ev, e := receiveEvent(reader)
		if e != nil {
			return fmt.Errorf("Error reading event: %w", e)
		}
		if e := onEvent(ev); e != nil {
			return fmt.Errorf("Error acting on event: %w", e)
		}
	}
}

func receiveEvent(reader *bufio.Reader) (Event, error) {
	var buffer strings.Builder
	for {
		line, e := reader.ReadString('\n')
		if e != nil {
			return Event{}, fmt.Errorf("Error reading event: %w", e)
		}
		if line == "\n" || line == "\r\n" {
			return parseSSEEvent(buffer.String()), nil
		} else {
			buffer.WriteString(line)
		}
	}
}

func parseSSEEvent(raw string) Event {
	var e Event
	var dataLines []string

	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "event: "):
			e.Event = strings.TrimPrefix(line, "event: ")
		case strings.HasPrefix(line, "data: "):
			dataLines = append(dataLines, strings.TrimPrefix(line, "data: "))
		case strings.HasPrefix(line, "id: "):
			e.Id = strings.TrimPrefix(line, "id: ")
		}
	}

	if len(dataLines) > 0 {
		e.Data = strings.Join(dataLines, "")
	}

	return e
}

func Get(config RequestConfig) (*http.Response, []byte, error) {
	return req(fullConfig{ConfigWithBody: ConfigWithBody{RequestConfig: config}, Method: "GET"})
}

func Delete(config RequestConfig) (*http.Response, []byte, error) {
	return req(fullConfig{ConfigWithBody: ConfigWithBody{RequestConfig: config}, Method: "DELETE"})
}

func Post(config ConfigWithBody) (*http.Response, []byte, error) {
	return req(fullConfig{ConfigWithBody: config, Method: "POST"})
}

func Put(config ConfigWithBody) (*http.Response, []byte, error) {
	return req(fullConfig{ConfigWithBody: config, Method: "PUT"})
}

func Patch(config ConfigWithBody) (*http.Response, []byte, error) {
	return req(fullConfig{ConfigWithBody: config, Method: "PATCH"})
}

func req(config fullConfig) (*http.Response, []byte, error) {

	u, e := url.Parse(config.Url)
	if e != nil {
		return nil, nil, fmt.Errorf("Error parsing url: %w", e)
	}

	for k, v := range config.Query {
		u.Query().Set(k, v)
	}

	req, e := http.NewRequest(config.Method, u.String(), bytes.NewBuffer(config.Body))
	if e != nil {
		return nil, nil, fmt.Errorf("Error creating request: %w", e)
	}

	for k, v := range config.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, e := client.Do(req)
	if e != nil {
		return nil, nil, fmt.Errorf("Error executing request: %w", e)
	}

	defer resp.Body.Close()
	resBody, e := io.ReadAll(resp.Body)
	if e != nil {
		return nil, nil, fmt.Errorf("Error reading response body: %w", e)
	}

	return resp, resBody, nil
}
