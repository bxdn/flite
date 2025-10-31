package flite

import (
	"fmt"
	"net/http"
)

var defaultServer server

type server struct {
	m           *http.ServeMux
	endpoints   []Endpoint
	middlewares []middleware
}

func Use(mid ...middleware) {
	defaultServer.middlewares = append(defaultServer.middlewares, mid...)
}

func Serve(port int) error {
	defaultServer.m = http.NewServeMux()
	return serveMux(port)
}

func ServeDebug(port int) error {
	defaultServer.m = http.DefaultServeMux
	return serveMux(port)
}

func serveMux(port int) error {
	for _, endpoint := range defaultServer.endpoints {
		defaultServer.m.HandleFunc(endpoint.Path(), endpoint.Handler())
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", port), getHandler(defaultServer.m))
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
