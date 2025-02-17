package flite

import (
	"fmt"
	"net/http"
)

type defaultServer struct{}

func (ds *defaultServer) Register(endpoints []Endpoint) {
	for _, endpoint := range endpoints {
		http.DefaultServeMux.HandleFunc(endpoint.path, endpoint.handler)
	}
}

func (ds *defaultServer) Serve(port int) {
	http.ListenAndServe(fmt.Sprintf(":%d", port), http.DefaultServeMux)
}

var DefaultServer = &defaultServer{}
