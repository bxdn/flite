package flite

import (
	"fmt"
	"log"
	"net/http"
)

type Endpoint interface {
	Path() string
	Handler() RequestHandler
}

type Never struct{}

type endpoint[T any] struct {
	handlers    		[]RequestNode[T]
	allowedMethod, path string
}

// Creates an endpoint from a given path.
//
// Uses ServeMux path syntax.
func CreateEndpoint(path string) *endpoint[Never] {
	ep := endpoint[Never]{path: path}
	return &ep
}

func CreateJsonEndpoint[T any](path string) *endpoint[T] {
	ep := endpoint[T]{path: path}
	return &ep
}

func (e *endpoint[T]) Path() string {
	return fmt.Sprintf("%s %s", e.allowedMethod, e.path)
}

func (e *endpoint[T]) Handler() RequestHandler {
	return e.handleRequest
}

func GET(path string, handlers ...RequestNode[Never]) *endpoint[Never] {
	e := endpoint[Never]{path: path, handlers: handlers, allowedMethod: "GET"}
	return &e
}

func POST[T any](path string, handlers ...RequestNode[T]) *endpoint[T] {
	e := endpoint[T]{path: path, handlers: handlers, allowedMethod: "POST"}
	return &e
}

func PUT[T any](path string, handlers ...RequestNode[T]) *endpoint[T] {
	e := endpoint[T]{path: path, handlers: handlers, allowedMethod: "PUT"}
	return &e
}

func DELETE(path string, handlers ...RequestNode[Never]) *endpoint[Never] {
	e := endpoint[Never]{path: path, handlers: handlers, allowedMethod: "DELETE"}
	return &e
}

func PATCH[T any](path string, handlers ...RequestNode[T]) *endpoint[T] {
	e := endpoint[T]{path: path, handlers: handlers, allowedMethod: "PATCH"}
	return &e
}

func (e *endpoint[T]) executeEndpointPipeline(w http.ResponseWriter, r *http.Request, handlers []RequestNode[T]) {
	f := &Flite[T]{res: &statusCacheResponseWriter{ResponseWriter: w}, req: r}
	for _, handler := range handlers {
		if e := handler(f); e != nil {
			log.Printf("ERROR: %v\n", e)
		}
		if f.done {
			return
		}
	}
	if !f.done {
		if e := f.Return(); e != nil {
			log.Printf("ERROR: %v\n", e)
		}
	}
}
func (e *endpoint[T]) handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", e.allowedMethod)
	e.executeEndpointPipeline(w, r, e.handlers)
}
