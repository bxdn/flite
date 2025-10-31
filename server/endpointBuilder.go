package flite

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
	return e.handleRequest
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

	// universal middleware
	for _, m := range defaultServer.middlewares {
		e, halt := m(w, r)
		if e != nil {
			log.Printf("ERROR: %v\n", e)
			return
		}
		if halt {
			return
		}
	}

	// flite handlers
	f := &F[T]{res: &statusCacheResponseWriter{ResponseWriter: w}, req: r}
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
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Headers", "*")
	// w.Header().Set("Access-Control-Allow-Methods", e.allowedMethod)
	e.executeEndpointPipeline(w, r, e.handlers)
}
