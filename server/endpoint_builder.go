package server

import (
	"fmt"
	"log"
	"net/http"
)

type Endpoint interface {
	Path() string
	Handler() func(http.ResponseWriter, *http.Request)
}

type No struct{}

type endpoint[T any] struct {
	handlers            []func(*F[T]) error
	allowedMethod, path string
}

func (e *endpoint[T]) Path() string {
	return fmt.Sprintf("%s %s", e.allowedMethod, e.path)
}

func (e *endpoint[T]) Handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		e.executeEndpointPipeline(w, r, e.handlers)
	}
}

func GET(path string, handlers ...func(*F[No]) error) {
	e := endpoint[No]{path: path, handlers: injectMiddleware(handlers), allowedMethod: "GET"}
	defaultServer.endpoints = append(defaultServer.endpoints, &e)
}

func POST[T any](path string, handlers ...func(*F[T]) error) {
	e := endpoint[T]{path: path, handlers: injectMiddleware(handlers), allowedMethod: "POST"}
	defaultServer.endpoints = append(defaultServer.endpoints, &e)
}

func PUT[T any](path string, handlers ...func(*F[T]) error) {
	e := endpoint[T]{path: path, handlers: injectMiddleware(handlers), allowedMethod: "PUT"}
	defaultServer.endpoints = append(defaultServer.endpoints, &e)
}

func DELETE(path string, handlers ...func(*F[No]) error) {
	e := endpoint[No]{path: path, handlers: injectMiddleware(handlers), allowedMethod: "DELETE"}
	defaultServer.endpoints = append(defaultServer.endpoints, &e)
}

func PATCH[T any](path string, handlers ...func(*F[T]) error) {
	e := endpoint[T]{path: path, handlers: injectMiddleware(handlers), allowedMethod: "PATCH"}
	defaultServer.endpoints = append(defaultServer.endpoints, &e)
}

func injectMiddleware[T any](rest []func(*F[T]) error) []func(*F[T]) error {
	handlers := []func(*F[T]) error{DeserializeBody[T]()}
	return append(handlers, rest...)
}

type middleware = func(w http.ResponseWriter, r *http.Request) (error, bool)

func (e *endpoint[T]) executeEndpointPipeline(w http.ResponseWriter, r *http.Request, handlers []func(*F[T]) error) {

	// flite handlers
	wf, ok := w.(WriterFlusher)
	if !ok {
		log.Printf("ERROR: Response writer is not a flusher")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f := &F[T]{res: &statusCacheResponseWriter{WriterFlusher: wf}, req: r}

	// universal middleware
	for _, m := range defaultServer.middlewares {
		e, halt := m(w, r)
		if e != nil {
			log.Printf("ERROR: %v\n", e)
			f.ReturnError("Error in middleware", http.StatusInternalServerError)
			return
		}
		if halt {
			return
		}
	}

	for _, handler := range handlers {
		if e := handler(f); e != nil {
			log.Printf("ERROR: %v\n", e)
			f.ReturnError("Error in handler", http.StatusInternalServerError)
		}
		if f.done {
			return
		}
	}
	if !f.done {
		if e := f.Return(); e != nil {
			log.Printf("ERROR: %v\n", e)
			f.ReturnError("Error returning", http.StatusInternalServerError)
		}
	}
}
