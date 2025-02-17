package common

import (
	"log"
	"net/http"
	"strings"
)

type endpointHandler struct {
	getHandlers    []RequestNode
	postHandlers   []RequestNode
	deleteHandlers []RequestNode
	putHandlers    []RequestNode
	allowedMethods string
	path           string
}

func (e *endpointHandler) executeEndpointPipeline(w http.ResponseWriter, r *http.Request, handlers []RequestNode) {
	for _, handler := range handlers {
		ctx, e := handler(w, r)
		if e != nil {
			log.Println(e)
			return
		}
		r = r.WithContext(ctx)
	}
}

func (e *endpointHandler) Get(w http.ResponseWriter, r *http.Request) {
	if e.getHandlers == nil {
		w.Header().Set("Allow", e.allowedMethods)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	} else {
		e.executeEndpointPipeline(w, r, e.getHandlers)
	}
}

func (e *endpointHandler) Post(w http.ResponseWriter, r *http.Request) {
	if e.postHandlers == nil {
		w.Header().Set("Allow", e.allowedMethods)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	} else {
		e.executeEndpointPipeline(w, r, e.postHandlers)
	}
}

func (e *endpointHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if e.deleteHandlers == nil {
		w.Header().Set("Allow", e.allowedMethods)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	} else {
		e.executeEndpointPipeline(w, r, e.deleteHandlers)
	}
}

func (e *endpointHandler) Put(w http.ResponseWriter, r *http.Request) {
	if e.putHandlers == nil {
		w.Header().Set("Allow", e.allowedMethods)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	} else {
		e.executeEndpointPipeline(w, r, e.putHandlers)
	}
}

func (e *endpointHandler) GET(handlers ...RequestNode) *endpointHandler {
	e.getHandlers = handlers
	if !(strings.Contains(e.allowedMethods, "GET")) {
		if e.allowedMethods != "" {
			e.allowedMethods += ", "
		}
		e.allowedMethods += "GET"
	}
	return e
}

func (e *endpointHandler) POST(handlers ...RequestNode) *endpointHandler {
	e.postHandlers = handlers
	if !(strings.Contains(e.allowedMethods, "POST")) {
		if e.allowedMethods != "" {
			e.allowedMethods += ", "
		}
		e.allowedMethods += "POST"
	}
	return e
}

func (e *endpointHandler) DELETE(handlers ...RequestNode) *endpointHandler {
	e.deleteHandlers = handlers
	if !(strings.Contains(e.allowedMethods, "DELETE")) {
		if e.allowedMethods != "" {
			e.allowedMethods += ", "
		}
		e.allowedMethods += "DELETE"
	}
	return e
}

func (e *endpointHandler) PUT(handlers ...RequestNode) *endpointHandler {
	e.putHandlers = handlers
	if !(strings.Contains(e.allowedMethods, "PUT")) {
		if e.allowedMethods != "" {
			e.allowedMethods += ", "
		}
		e.allowedMethods += "PUT"
	}
	return e
}

func (e *endpointHandler) Build() Endpoint {
	return Endpoint{e.path, e.handleRequest}
}

func CreateEndpointBuilder(path string) *endpointHandler {
	builder := endpointHandler{}
	builder.path = path
	return &builder
}

func (e *endpointHandler) handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", e.allowedMethods)
	switch r.Method {
	case http.MethodGet:
		e.Get(w, r)
	case http.MethodPost:
		e.Post(w, r)
	case http.MethodDelete:
		e.Delete(w, r)
	case http.MethodPut:
		e.Put(w, r)
	case http.MethodOptions:
		w.Header().Set("Allow", e.allowedMethods)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
