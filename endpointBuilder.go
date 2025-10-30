package flite

import (
	"log"
	"net/http"
	"strings"
)

type Endpoint interface {
	Path() string
	Handler() RequestHandler
	GET(handlers ...RequestNode) Endpoint
	POST(handlers ...RequestNode) Endpoint
	DELETE(handlers ...RequestNode) Endpoint
	PUT(handlers ...RequestNode) Endpoint
}

type endpoint struct {
	getHandlers    []RequestNode
	postHandlers   []RequestNode
	deleteHandlers []RequestNode
	putHandlers    []RequestNode
	allowedMethods string
	path           string
}

// Creates an endpoint from a given path.
//
// Uses ServeMux path syntax.
func CreateEndpoint(path string) Endpoint {
	ep := endpoint{}
	ep.path = path
	return &ep
}

func (e *endpoint) Path() string {
	return e.path
}

func (e *endpoint) Handler() RequestHandler {
	return e.handleRequest
}

func (e *endpoint) GET(handlers ...RequestNode) Endpoint {
	e.getHandlers = handlers
	e.addAllowedMethod("GET")
	return e
}

func (e *endpoint) POST(handlers ...RequestNode) Endpoint {
	e.postHandlers = handlers
	e.addAllowedMethod("POST")
	return e
}

func (e *endpoint) DELETE(handlers ...RequestNode) Endpoint {
	e.deleteHandlers = handlers
	e.addAllowedMethod("DELETE")
	return e
}

func (e *endpoint) PUT(handlers ...RequestNode) Endpoint {
	e.putHandlers = handlers
	e.addAllowedMethod("PUT")
	return e
}

func (e *endpoint) addAllowedMethod(method string) {
	if !(strings.Contains(e.allowedMethods, method)) {
		if e.allowedMethods != "" {
			e.allowedMethods += ", "
		}
		e.allowedMethods += method
	}
}

func (e *endpoint) executeEndpointPipeline(w http.ResponseWriter, r *http.Request, handlers []RequestNode) {
	f := &Flite{res: &statusCacheResponseWriter{ResponseWriter: w}, req: r}
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

func (e *endpoint) get(w http.ResponseWriter, r *http.Request) {
	if e.getHandlers == nil {
		w.Header().Set("Allow", e.allowedMethods)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	} else {
		e.executeEndpointPipeline(w, r, e.getHandlers)
	}
}

func (e *endpoint) post(w http.ResponseWriter, r *http.Request) {
	if e.postHandlers == nil {
		w.Header().Set("Allow", e.allowedMethods)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	} else {
		e.executeEndpointPipeline(w, r, e.postHandlers)
	}
}

func (e *endpoint) delete(w http.ResponseWriter, r *http.Request) {
	if e.deleteHandlers == nil {
		w.Header().Set("Allow", e.allowedMethods)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	} else {
		e.executeEndpointPipeline(w, r, e.deleteHandlers)
	}
}

func (e *endpoint) put(w http.ResponseWriter, r *http.Request) {
	if e.putHandlers == nil {
		w.Header().Set("Allow", e.allowedMethods)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	} else {
		e.executeEndpointPipeline(w, r, e.putHandlers)
	}
}

func (e *endpoint) handleRequest(w http.ResponseWriter, r *http.Request) {
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
