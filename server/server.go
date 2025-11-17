package server

import (
	"fmt"
	"net/http"
)

var defaultServer server

type server struct {
	endpoints   []Endpoint
	middlewares []middleware
}

func Use(mid ...middleware) {
	defaultServer.middlewares = append(defaultServer.middlewares, mid...)
}

func Serve(port int) error {
	return serveMux(http.NewServeMux(), port)
}

func ServeDebug(port int) error {
	return serveMux(http.DefaultServeMux, port)
}

func serveMux(mux *http.ServeMux, port int) error {
	for _, endpoint := range defaultServer.endpoints {
		mux.HandleFunc(endpoint.Path(), endpoint.Handler())
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", port), getHandler(mux))
}

func getHandler(mux *http.ServeMux) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, PATCH, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		mux.ServeHTTP(w, r)
	})
}
