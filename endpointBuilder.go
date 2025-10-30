package flite

import (
	"log"
	"net/http"
	"strings"
)

type Endpoint interface {
	Path() string
	Handler() RequestHandler
}

type endpoint[T any] struct {
	getHandlers    []RequestNode[T]
	postHandlers   []RequestNode[T]
	deleteHandlers []RequestNode[T]
	putHandlers    []RequestNode[T]
	allowedMethods string
	path           string
}

// Creates an endpoint from a given path.
//
// Uses ServeMux path syntax.
func CreateEndpoint(path string) *endpoint[string] {
	ep := endpoint[string]{}
	ep.path = path
	return &ep
}

func CreateJsonEndpoint[T any](path string) *endpoint[T] {
	ep := endpoint[T]{}
	ep.path = path
	return &ep
}

func (e *endpoint[T]) Path() string {
	return e.path
}

func (e *endpoint[T]) Handler() RequestHandler {
	return e.handleRequest
}

func (e *endpoint[T]) GET(handlers ...RequestNode[T]) *endpoint[T] {
	e.getHandlers = handlers
	e.addAllowedMethod("GET")
	return e
}

func (e *endpoint[T]) POST(handlers ...RequestNode[T]) *endpoint[T] {
	e.postHandlers = handlers
	e.addAllowedMethod("POST")
	return e
}

func (e *endpoint[T]) DELETE(handlers ...RequestNode[T]) *endpoint[T] {
	e.deleteHandlers = handlers
	e.addAllowedMethod("DELETE")
	return e
}

func (e *endpoint[T]) PUT(handlers ...RequestNode[T]) *endpoint[T] {
	e.putHandlers = handlers
	e.addAllowedMethod("PUT")
	return e
}

func (e *endpoint[T]) addAllowedMethod(method string) {
	if !(strings.Contains(e.allowedMethods, method)) {
		if e.allowedMethods != "" {
			e.allowedMethods += ", "
		}
		e.allowedMethods += method
	}
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

func (e *endpoint[T]) get(w http.ResponseWriter, r *http.Request) {
	if e.getHandlers == nil {
		w.Header().Set("Allow", e.allowedMethods)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	} else {
		e.executeEndpointPipeline(w, r, e.getHandlers)
	}
}

func (e *endpoint[T]) post(w http.ResponseWriter, r *http.Request) {
	if e.postHandlers == nil {
		w.Header().Set("Allow", e.allowedMethods)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	} else {
		e.executeEndpointPipeline(w, r, e.postHandlers)
	}
}

func (e *endpoint[T]) delete(w http.ResponseWriter, r *http.Request) {
	if e.deleteHandlers == nil {
		w.Header().Set("Allow", e.allowedMethods)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	} else {
		e.executeEndpointPipeline(w, r, e.deleteHandlers)
	}
}

func (e *endpoint[T]) put(w http.ResponseWriter, r *http.Request) {
	if e.putHandlers == nil {
		w.Header().Set("Allow", e.allowedMethods)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	} else {
		e.executeEndpointPipeline(w, r, e.putHandlers)
	}
}

func (e *endpoint[T]) handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", e.allowedMethods)
	switch r.Method {
	case http.MethodGet:
		e.get(w, r)
	case http.MethodPost:
		e.post(w, r)
	case http.MethodDelete:
		e.delete(w, r)
	case http.MethodPut:
		e.put(w, r)
	case http.MethodOptions:
		w.Header().Set("Allow", e.allowedMethods)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
