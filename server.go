package flite

import (
	"fmt"
	"net/http"
)

var defaultServer server

type server struct {
	m *http.ServeMux
	endpoints []Endpoint
}

func Serve(port int) error {
	defaultServer.m = http.NewServeMux()
	for _, endpoint := range defaultServer.endpoints {
		defaultServer.m.HandleFunc(endpoint.Path(), endpoint.Handler())
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", port), defaultServer.m)
}

func ServeDebug(port int) error {
	defaultServer.m = http.DefaultServeMux
	for _, endpoint := range defaultServer.endpoints {
		defaultServer.m.HandleFunc(endpoint.Path(), endpoint.Handler())
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", port), defaultServer.m)
}