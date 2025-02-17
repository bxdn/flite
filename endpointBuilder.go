package flite

import (
	"log"
	"net/http"
	"strings"
)

type endpoint struct {
	getHandlers    []RequestNode
	postHandlers   []RequestNode
	deleteHandlers []RequestNode
	putHandlers    []RequestNode
	allowedMethods string
	path           string
}

func CreateEndpoint(path string) *endpoint {
	builder := endpoint{}
	builder.path = path
	return &builder
}

func (e *endpoint) Path() string {
	return e.path
}

func (e *endpoint) Handler() RequestHandler {
	return e.handleRequest
}

func (e *endpoint) GET(handlers ...RequestNode) *endpoint {
	e.getHandlers = handlers
	if !(strings.Contains(e.allowedMethods, "GET")) {
		if e.allowedMethods != "" {
			e.allowedMethods += ", "
		}
		e.allowedMethods += "GET"
	}
	return e
}

func (e *endpoint) POST(handlers ...RequestNode) *endpoint {
	e.postHandlers = handlers
	if !(strings.Contains(e.allowedMethods, "POST")) {
		if e.allowedMethods != "" {
			e.allowedMethods += ", "
		}
		e.allowedMethods += "POST"
	}
	return e
}

func (e *endpoint) DELETE(handlers ...RequestNode) *endpoint {
	e.deleteHandlers = handlers
	if !(strings.Contains(e.allowedMethods, "DELETE")) {
		if e.allowedMethods != "" {
			e.allowedMethods += ", "
		}
		e.allowedMethods += "DELETE"
	}
	return e
}

func (e *endpoint) PUT(handlers ...RequestNode) *endpoint {
	e.putHandlers = handlers
	if !(strings.Contains(e.allowedMethods, "PUT")) {
		if e.allowedMethods != "" {
			e.allowedMethods += ", "
		}
		e.allowedMethods += "PUT"
	}
	return e
}

func (e *endpoint) executeEndpointPipeline(w http.ResponseWriter, r *http.Request, handlers []RequestNode) {
	for _, handler := range handlers {
		ctx, e := handler(w, r)
		if e != nil {
			log.Println(e)
			return
		}
		r = r.WithContext(ctx)
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
